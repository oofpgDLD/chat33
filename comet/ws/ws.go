package ws

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/33cn/chat33/db"

	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/comet"
	"github.com/33cn/chat33/model"

	"github.com/33cn/chat33/types"

	"github.com/33cn/chat33/router"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	l "github.com/inconshreveable/log15"
)

var wsLog = l.New("module", "chat/comet")

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// allow cross-origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//新客户端登上，挤掉原先的客户端
func closeOld(userId, device, clientId, uuid string, createTime int64) {
	u, ok := router.GetUser(userId)
	if !ok {
		return
	}

	var oldDevice *router.Client
	switch device {
	case types.DeviceAndroid, types.DeviceIOS:
		for _, c := range u.GetClients() {
			if clientId != c.GetId() && (c.GetDevice() == types.DeviceAndroid || c.GetDevice() == types.DeviceIOS) {
				oldDevice = c
				break
			}
		}
	case types.DevicePC:
		for _, c := range u.GetClients() {
			if clientId != c.GetId() && (c.GetDevice() == types.DevicePC) {
				oldDevice = c
				break
			}
		}
	}

	if oldDevice != nil {
		wsLog.Debug("close old client", "old client id", oldDevice.GetId(), "new client id", clientId, "old uuid", oldDevice.GetUUID(), "new uuid", uuid, "old create_time", oldDevice.GetCreateTime(), "new create_time", createTime)
		code := websocket.CloseNormalClosure
		data := "你的账号已在其他端登录"
		if oldDevice.GetUUID() != uuid {
			wsLog.Debug("kick off old connection", "oldClientId", oldDevice.GetId())
			// send kicked off notification to cli
			//获取数据库中上一次该user的连接信息
			lstInfo, _ := model.GetLastDeviceLoginInfo(userId, device)
			//alertMsg := model.GetAlertMsg(lstInfo)
			alertMsg := model.GetAlertMsgV2(lstInfo)
			//proto.SendOtherDeviceLogin(alertMsg, cli)
			code = 4001
			wsLog.Debug("len of ws colse msg", "len", len(alertMsg))
			data = alertMsg
		}
		//老设备断开
		err := oldDevice.Close(code, data)
		if err != nil {
			wsLog.Error("close client failed", "err", err)
		}
	}
}

//ws连接时的过滤器,用于conn、wsconn、client创建之前
//过滤条件：1.未认证用户 2.被挤掉的设备重连
//返回值：bool true：通过 false 过滤
func connFilter(userId, device, clientId string, level int, loginTime int64) bool {
	if level == types.LevelVisitor {
		wsLog.Debug("reject visitor")
		return false
	}

	u, ok := router.GetUser(userId)
	if !ok {
		return true
	}
	//find old client
	var oldClient *router.Client
	switch device {
	case types.DeviceAndroid, types.DeviceIOS:
		for _, c := range u.GetClients() {
			if clientId != c.GetId() && (c.GetDevice() == types.DeviceAndroid || c.GetDevice() == types.DeviceIOS) {
				oldClient = c
				break
			}
		}
	case types.DevicePC:
		for _, c := range u.GetClients() {
			if clientId != c.GetId() && (c.GetDevice() == types.DevicePC) {
				oldClient = c
				break
			}
		}
	}

	if oldClient != nil && loginTime < oldClient.GetLoginTime() {
		wsLog.Debug("reject reconnect")
		return false
	}
	return true
}

type UserInfoHandel func(store *sessions.CookieStore, context *gin.Context) (userId, device, uuid, appId string, level int, loginTime int64)

// ServeWs handles websocket requests from the peer.
func ServeWs(store *sessions.CookieStore, context *gin.Context, callInfo UserInfoHandel) {
	userId, device, uuid, appId, level, loginTime := callInfo(store, context)
	clientId := utility.RandomID()
	createTime := utility.NowMillionSecond()

	//过滤连接
	if !connFilter(userId, device, clientId, level, loginTime) {
		wsLog.Debug("[Process] Ws connFilter reject connection", "userId", userId, "device", device, "appId", appId, "level", level, "client_id", clientId, "uuid", uuid, "create time", createTime, "login time", loginTime)
		/*ret := result.ComposeHttpAck(result.LoginExpired, "", "")
		context.PureJSON(http.StatusOK, ret)*/
		context.String(http.StatusUnauthorized, "check user is login")
		return
	}

	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		wsLog.Error("[Process] Ws Upgrade Failed", "err", err, "appId", appId)
		return
	}
	wsLog.Debug("[Process] Ws Upgrade Success", "userId", userId, "device", device, "appId", appId, "level", level, "client_id", clientId, "uuid", uuid, "create time", createTime, "login time", loginTime)

	//查询所有群
	rooms, err := orm.GetUserJoinedRooms(userId)
	if err != nil {
		wsLog.Error("[Process] get user joined rooms failed", "err", err.Error())
		err := conn.Close()
		if err != nil {
			wsLog.Error("[Process] close client failed", "err", err.Error())
		}
		/*ret := result.ComposeHttpAck(result.ServerInterError, "", "")
		context.PureJSON(http.StatusOK, ret)*/
		context.String(http.StatusInternalServerError, "get user joined rooms error")
		return
	}

	//断开老设备
	closeOld(userId, device, clientId, uuid, createTime)

	args := make(map[string]string)
	args["clientId"] = clientId
	args["userId"] = userId
	args["device"] = device
	args["appId"] = appId
	wsConn := NewWSConnection(conn, args, MsgRecvCb(onRecvMsg), MsgSendCb(onSendMsg), CloseCb(onClose))

	u, ok := router.GetUser(userId)
	if !ok || u == nil {
		u = router.NewUser(userId, level)
		wsLog.Info("add user", "user", u)
		router.AddUser(u)
	}
	client := router.NewClient(clientId, device, uuid, createTime, loginTime, wsConn, u)
	u.AppendClient(client)

	//订阅所有群的消息通道
	for _, v := range rooms {
		chId := types.GetRoomRouteById(v)
		if cl, ok := router.GetChannel(chId); ok {
			u.Subscribe(cl)
		}
	}

	wsLog.Info("[Process] start run", "appId", appId, "userId", userId, "device", device, "level", level)
	go func() {
		err := client.Run()
		if err != nil {
			err := conn.Close()
			if err != nil {
				wsLog.Error("[Process] close client failed where run failed", "err", err.Error())
			}
			wsLog.Error("[Process] Run failed", "appId", appId, "err", err)
		}
	}()
}

func onRecvMsg(wsc *WsConnection, msg []byte) {
	msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
	comet.CreateProto(wsc, msg)
}

func onSendMsg(wsc *WsConnection, message interface{}) ([]byte, error) {
	var err error
	var toSend []byte

	switch msg := message.(type) {
	case []byte:
		toSend = msg
	default:
		err = fmt.Errorf("invalid to send message type: %v", reflect.TypeOf(message))
		wsLog.Error("on send msg", "err", err.Error())
		return nil, err
	}

	return toSend, nil
}

func onClose(wsc *WsConnection) {
	//判断是否写入关闭帧
	if wsc.closeCode != 0 {
		//msg 最大125字节
		err := wsc.baseConn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(wsc.closeCode, wsc.closeMsg), time.Now().Add(wsc.writeWait))
		if err != nil {
			wsLog.Error("Error write close frame", "err", err)
		}
	}

	args := wsc.Args().(map[string]string)
	clientId := args["clientId"]
	userId := args["userId"]
	device := args["device"]
	u, ok := router.GetUser(userId)
	if !ok {
		wsLog.Error("onClose invalid userId", "userId", userId)
		return
	}

	//TODO 采用redis
	maps, _ := db.GetRoomsById(u.Id)
	rooms := make([]string, 0)
	for _, v := range maps {
		rooms = append(rooms, v["room_id"])
	}
	for _, room := range rooms {
		roomChannelKey := types.GetRoomRouteById(room)
		if cl, ok := router.GetChannel(roomChannelKey); ok {
			u.ClientUnRegister(clientId, cl)
		}
	}
	//输出信息
	c, _ := u.GetClient(clientId)
	if c != nil {
		wsLog.Debug("del client", "client_id", clientId, "uuid", c.GetUUID())
	}
	u.DeleteClient(clientId)
	wsLog.Info("client closed", "device", device, "userId", userId)
}

var sessionName = "session-login"

func GetUserInfo(store *sessions.CookieStore, c *gin.Context) (userId, device, uuid, appId string, level int, loginTime int64) {
	/*device = c.GetHeader("FZM-DEVICE")
	uuid = c.GetHeader("FZM-UUID")
	appId = c.GetHeader("FZM-APP-ID")*/

	//获取当前用户id
	session, err := store.Get(c.Request, sessionName)
	if err != nil {
		wsLog.Error("ws can not get session", "err", err.Error())
		return
	}

	userId = utility.ToString(session.Values["user_id"])
	device = utility.ToString(session.Values["devtype"])
	uuid = utility.ToString(session.Values["uuid"])
	appId = utility.ToString(session.Values["appId"])
	loginTime = utility.ToInt64(session.Values["time"])

	switch device {
	case types.DevicePC:
	case types.DeviceAndroid:
	case types.DeviceIOS:
	default:
		device = ""
	}

	if userId == "" || device == "" {
		level = router.VISITOR
	} else {
		level = router.NOMALUSER
	}
	return
}

//挤掉其他所有端
func CloseOther(userId, device, uuid string) {
	u, ok := router.GetUser(userId)
	if !ok {
		return
	}

	var oldDevices = make([]*router.Client, 0)
	for _, c := range u.GetClients() {
		if c.GetDevice() != device || c.GetUUID() != uuid {
			oldDevices = append(oldDevices, c)
		}
	}

	for _, oldDevice := range oldDevices {
		wsLog.Debug("close other client", "old client id", oldDevice.GetId(), "old uuid", oldDevice.GetUUID(), "old create_time", oldDevice.GetCreateTime())
		code := websocket.CloseNormalClosure
		data := "你的账号已在其他端登录"

		wsLog.Debug("kick off old connection", "oldClientId", oldDevice.GetId())
		// send kicked off notification to cli
		//获取数据库中上一次该user的连接信息
		lstInfo, _ := model.GetLastDeviceLoginInfo(userId, device)
		//alertMsg := model.GetAlertMsg(lstInfo)
		alertMsg := model.GetAlertMsgV2(lstInfo)
		//proto.SendOtherDeviceLogin(alertMsg, cli)
		code = 4011
		wsLog.Debug("len of ws colse msg", "len", len(alertMsg))
		data = alertMsg

		//老设备断开
		err := oldDevice.Close(code, data)
		if err != nil {
			wsLog.Error("close client failed", "err", err)
		}
	}
}
