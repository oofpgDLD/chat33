package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/app"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/pkg/account"
	"github.com/33cn/chat33/pkg/http"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/inconshreveable/log15"
)

var logRedpacket = log15.New("module", "model/red_envelope")

type coinInfo struct {
	CoinId        int     `json:"coinId"`
	CoinName      string  `json:"coinName"`
	DecimalPlaces int     `json:"decimalPlaces"`
	CoinFullName  string  `json:"coinFullName"`
	CoinNickname  string  `json:"coinNickname"`
	IconUrl       string  `json:"iconUrl"`
	SingleMax     float64 `json:"singleMax"`
	SingleMin     float64 `json:"singleMin"`
	DailyMax      float64 `json:"dailyMax"`
}

func InsertPacketInfo(packetID, userID, toID string, toType, size int, amount float64, remark string, cType, coin int, time int64) error {
	return mysql.InsertPacketInfo(packetID, userID, toID, utility.ToString(toType), utility.ToString(size), utility.ToString(amount), remark, cType, coin, time)
}

//从本地数据库中查询红包信息
func getRedPacketInfo(packetId string) (*types.RedPacketLog, error) {
	return mysql.GetRedPacketInfo(packetId)
}

var respErrToSysErr = map[int]int{
	10005: result.RPEmpty,      //红包已经领完
	10006: result.RPExpired,    //红包已过期
	10007: result.RPReceived,   //红包已经领取
	8002:  result.LoginExpired, //用户身份非法（登录过期）
}

type respError struct {
	ErrorCode int
	Message   string
}

func (e *respError) Error() string {
	return e.Message
}

func (e *respError) ToChatErrorCode() int {
	if code, ok := respErrToSysErr[e.ErrorCode]; ok {
		return code
	} else {
		return result.RPError
	}
}

func (e *respError) ToChatExtMsg() string {
	if _, ok := respErrToSysErr[e.ErrorCode]; ok {
		return ""
	} else {
		return e.Error()
	}
}

//规则：红包地址+参数 红包id 平台号码
func getRPH5Url(appId, userId, packetId string) string {
	app := app.GetApp(appId)
	if app == nil {
		return ""
	}

	inviteCode := ""
	user, err := orm.GetUserInfoById(userId)
	if err == nil && user != nil {
		inviteCode = user.InviteCode
	}
	return fmt.Sprintf(app.RPUrl+"?id=%s&platform_id=%s&invite_code=%s&app_type=chat", packetId, app.RPPid, inviteCode)
}

func checkRPResponse(byte []byte) (map[string]interface{}, error) {
	var resp map[string]interface{}
	err := json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, &respError{
			Message: err.Error(),
		}
	}

	/*
		{
			"id": 0,
			"result": null,
			"error": {
				"code": 10003,
				"message": "红包金额非法"
		}

		{
			"id": 0,
			"result": {
				"packet_id": "0fac3a55-9641-4869-bc69-92e9387476f1"
			},
			"error": {
				"code": 0,
				"message": ""
			}
		}
	*/

	if resp["error"] != nil {
		if error, ok := resp["error"].(map[string]interface{}); ok {
			if utility.ToInt(error["code"]) != 0 {
				return nil, &respError{
					ErrorCode: utility.ToInt(error["code"]),
					Message:   error["message"].(string),
				}
			}
			if result, ok := resp["result"].(map[string]interface{}); ok {
				return result, nil
			}
		}
	}
	return nil, &respError{
		Message: "格式错误",
	}
}

func ConvertRedPackInfoToSend(callerId string, con map[string]interface{}) error {
	//获取uid
	var appId string
	user, _ := orm.GetUserInfoById(callerId)
	if user != nil {
		appId = user.AppId
		app := app.GetApp(appId)
		if app == nil {
			return types.ERR_APPNOTFIND
		}
		if app.IsInner != types.IsInnerAccount {
			return nil
		}
	}
	// 是否已领取
	packetId := utility.ToString(con["packetId"])
	/*rpInfo, err := redPacketInfo(packetId)
	if err != nil {
		return nil
	}
	if len([]rune(utility.ToString(rpInfo["remark"]))) > 12 {
		con["remark"] = string([]rune(utility.ToString(rpInfo["remark"])))[:12] + "..."
	} else {
		con["remark"] = utility.ToString(rpInfo["remark"])
	}*/
	var isOpened bool
	revs, _ := GetPacketRecvDetails(appId, callerId, packetId)
	for _, v := range revs {
		if callerId == utility.ToString(v.UserId) {
			isOpened = true
			break
		}
	}
	con["isOpened"] = isOpened
	//con["type"] = utility.ToInt(rpInfo["type"])
	return nil
}

func redPacketInfo(appId, packetId string) (map[string]interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	b, err := json.Marshal(&types.ReqBase{
		Method: "Info",
		Params: &types.InfoParams{
			PacketId: packetId,
		},
	})
	if err != nil {
		logRedpacket.Error("redPacketInfo json Marshal", "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, nil, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}
	/*
		{
			"id": 0,
			"result": {
				"user_id": "34",
				"user_name": "",
				"user_mobile": "157****6517",
				"platform_id": "1000",
				"type": 1,
				"coin_id": 3,
				"coin_name": "BTY",
				"amount": 10,
				"size": 10,
				"to": "",
				"remain": 9,
				"remain_amount": 9.37,
				"status": 1,
				"remark": "恭喜发财，大吉大利",
				"create_at": 1548312248,
				"update_at": 1548318099,
				"expire_at": 1548398648
			},
			"error": {
				"code": 0,
				"message": ""
			}
		}
	*/
	ret, err := checkRPResponse(byte)
	if err != nil {
		return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
			SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
	}
	return ret, nil
}

//单个用户指定红包领取信息
func receiveSingleRecord(appId, token, packetId string) (map[string]interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	headers := make(map[string]string)
	headers["Authorization"] = types.AuthType + " " + token
	headers["Platform"] = app.RPPid

	b, err := json.Marshal(&types.ReqBase{
		Method: "ReceiveSingleRecord",
		Params: &types.ReceiveDetailParams{
			PacketId: packetId,
		},
	})
	if err != nil {
		logRedpacket.Error("receiveSingleRecord json Marshal", "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, headers, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	ret, err := checkRPResponse(byte)
	if err != nil {
		//不是 未查到领取记录
		if err.(*respError).ErrorCode != 10011 {
			return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
				SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
		}
		return ret, nil
	}
	return ret, err
	/*{
		"id": int32,
		"result": {
		"amount": 25,
			"coin_id": 2,
			"coin_name": "BTC",
			"user_id": "158962",
			"user_name": "zhangsan",
			"user_mobile": "132****1234",
			"created_at": 1530003110,
			"status": 1,
			"fail_message":""
	}
	}*/
}

//将红包服务获取的代币信息存储到map中
func coinMapFromRPServer(appId string) (map[string]*coinInfo, error) {
	conInfo, err := getCoinInfoFromRPServer(appId)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*coinInfo)
	for _, c := range conInfo {
		ret[c.CoinName] = c
	}
	return ret, nil
}

func Balance(appId, token string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	headers := make(map[string]string)
	headers["Authorization"] = types.AuthType + " " + token
	headers["Platform"] = app.RPPid

	b, err := json.Marshal(&types.ReqBase{
		Method: "Balance",
	})
	if err != nil {
		logRedpacket.Error("Balance json Marshal", "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, headers, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	//获取余额信息
	balanceInfo, err := checkRPResponse(byte)
	if err != nil {
		return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
			SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
	}

	type Balance struct {
		CoinId        int     `json:"coinId"`
		CoinName      string  `json:"coinName"`
		CoinNickname  string  `json:"coinNickname"`
		CoinFullName  string  `json:"coinFullName"`
		DecimalPlaces int     `json:"decimalPlaces"`
		IconUrl       string  `json:"iconUrl"`
		Amount        float64 `json:"amount"`
		SingleMax     float64 `json:"singleMax"`
		SingleMin     float64 `json:"singleMin"`
		DailyMax      float64 `json:"dailyMax"`
		Fee           float64 `json:"fee"`
	}

	clist, ok := balanceInfo["balance"]
	if !ok {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("no data")
	}

	ret := make(map[string]interface{})
	balances := make([]*Balance, 0)
	search := make([]string, 0)
	for _, v := range clist.([]interface{}) {
		one := v.(map[string]interface{})
		b := &Balance{
			CoinId:   utility.ToInt(one["coin_id"]),
			CoinName: utility.ToString(one["coin_name"]),
			Amount:   utility.ToFloat64(one["active"]),
		}
		search = append(search, b.CoinName)
		balances = append(balances, b)
	}

	/*//获取coin图片
	coinPic, err := getCoinFromBCoin(search)
	if err != nil {
		return nil, err
	}*/
	pInfo, err := coinMapFromRPServer(appId)
	if err != nil {
		return nil, err
	}

	data, err := account.RPFeeInfoFromBcoin(app.AccountServer)
	feeInfo := make(map[string]float64)
	if data != nil {
		for _, v := range data.([]interface{}) {
			info := v.(map[string]interface{})
			feeInfo[utility.ToString(info["currency"])] = utility.ToFloat64(info["amount"])
		}
	}

	for _, b := range balances {
		//获取图片和中文名,精度
		if c, ok := pInfo[b.CoinName]; ok {
			b.IconUrl = c.IconUrl
			b.CoinNickname = c.CoinNickname
			b.DecimalPlaces = c.DecimalPlaces
			b.CoinFullName = c.CoinFullName
			b.SingleMin = c.SingleMin
			b.SingleMax = c.SingleMax
			b.DailyMax = c.DailyMax
		}
		//获取币种手续费信息
		if v, ok := feeInfo[b.CoinName]; ok {
			b.Fee = v
		}
	}
	ret["balances"] = balances
	return ret, nil
}

//从红包服务获取币种信息
func getCoinInfoFromRPServer(appId string) ([]*coinInfo, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	b, err := json.Marshal(&types.ReqBase{
		Method: "CoinInfo",
		Params: &types.CoinInfoParams{
			PlatformId: app.RPPid,
		},
	})
	if err != nil {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, nil, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	ret, err := checkRPResponse(byte)
	if err != nil {
		return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
			SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
	}

	coins := make([]*coinInfo, 0)
	cInfo, ok := ret["coin_info"].([]interface{})
	if !ok {
		logRedpacket.Warn("getCoinInfoFromRPServer", "warn", "can not get coin_info")
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("红包服务返回参数错误")
	}
	for _, c := range cInfo {
		info, ok := c.(map[string]interface{})
		if !ok {
			logRedpacket.Warn("getCoinInfoFromRPServer", "warn", "coin_info is nil")
			return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("红包服务返回参数错误")
		}
		coin := &coinInfo{
			CoinId:        utility.ToInt(info["coin_id"]),
			CoinName:      utility.ToString(info["coin_name"]),
			DecimalPlaces: utility.ToInt(info["decimal_places"]),
			CoinFullName:  utility.ToString(info["full_name"]),
			CoinNickname:  utility.ToString(info["nick_name"]),
			IconUrl:       utility.ToString(info["icon_url"]),
			SingleMax:     utility.ToFloat64(info["single_max"]),
			SingleMin:     utility.ToFloat64(info["single_min"]),
			DailyMax:      utility.ToFloat64(info["daily_max"]),
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

//获取代币信息
func GetCoinInfo(appId string) (interface{}, error) {
	coins, err := getCoinInfoFromRPServer(appId)
	if err != nil {
		return nil, err
	}
	return &struct {
		CoinInfo []*coinInfo `json:"coinInfo"`
	}{
		CoinInfo: coins,
	}, nil
}

// 发送红包
func Send(appId, token string, req *types.SendParams, userId, toId string, ctype int) (*types.RedPacket, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}
	//将ext参数重新组装下： 加上手续费
	ext := req.Ext.(map[string]interface{})
	//发送红包前先验证是否是好友  群成员
	if ctype == types.RPToRoom {
		//群是否存在  发红包者是否在群内
		b, err := CheckUserInRoom(userId, toId)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if !b {
			logRedpacket.Warn("Send", "warn", "UserIsNotInRoom", "userId", userId, "roomId", toId)
			return nil, result.NewError(result.UserIsNotInRoom).SetExtMessage("用户不在群内")
		}
		ext["need_fee"] = true
	} else if ctype == types.RPToUser {
		//用户是否存在  是否是好友
		b, err := orm.CheckIsFriend(userId, toId, types.IsNotDelete)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if !b {
			logRedpacket.Warn("Send", "warn", "IsNotFriend", "userId", userId, "friendId", toId)
			return nil, result.NewError(result.IsNotFriend)
		}
		ext["need_fee"] = false
	}
	req.Ext = ext

	headers := make(map[string]string)
	headers["Authorization"] = types.AuthType + " " + token
	headers["Platform"] = app.RPPid
	//限定remark长度
	arry := []rune(utility.ToString(req.Remark))
	if len(arry) > types.RemarkLengthLimit {
		req.Remark = string(arry[:types.RemarkLengthLimit]) + "..."
	}
	b, err := json.Marshal(&types.ReqBase{
		Method: "Send",
		Params: req,
	})
	if err != nil {
		logRedpacket.Error("Send json Marshal", "err", err)
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, headers, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	ret, err := checkRPResponse(byte)
	if err != nil {
		return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
			SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
	}

	packet := &types.RedPacket{
		PacketId:   utility.ToString(ret["packet_id"]),
		PacketType: req.Type,
		PacketUrl:  getRPH5Url(appId, userId, utility.ToString(ret["packet_id"])),
		Coin:       req.CoinId,
		Remark:     req.Remark,
	}

	return packet, nil
}

//校验是否能够领取
func proReceiveCheck(enable bool, userId string, p *types.RedPacketLog) error {
	if !enable {
		return nil
	}
	if p == nil {
		return result.NewError(result.RPTargetIllegal).SetExtMessage("无法领取外部红包")
	}
	if p.CType == types.RPToRoom {
		//群成员才能领取群内的红包
		b, err := CheckUserInRoom(userId, p.ToId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if !b {
			logRedpacket.Warn("proReceiveCheck", "warn", "UserIsNotInRoom", "userId", userId, "roomId", p.ToId)
			return result.NewError(result.UserIsNotInRoom)
		}
	} else if p.CType == types.RPToUser {
		//目标好友才能领取
		if userId != p.ToId && userId != p.UserId {
			logRedpacket.Warn("proReceiveCheck : redpacket is not for you", "warn", "RPTargetIllegal", "userId", userId, "target receiver", p.ToId)
			return result.NewError(result.RPTargetIllegal)
		}
	}
	return nil
}

// 领取红包
func ReceiveEntry(appId, token, packetId, userId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.RPTargetIllegal).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}
	//查找红包信息
	packetInfo, err := getRedPacketInfo(packetId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	err = proReceiveCheck(false, userId, packetInfo)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	headers["Authorization"] = types.AuthType + " " + token
	headers["Platform"] = app.RPPid

	b, err := json.Marshal(&types.ReqBase{
		Method: "Receive",
		Params: &types.ReceiveParams{
			PacketId: packetId,
		},
	})
	if err != nil {
		logRedpacket.Error("ReceiveEntry json Marshal", "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("返回值解析失败")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, headers, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	rlt, err := checkRPResponse(byte)
	if err != nil {
		return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
			SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
	}

	//再查出红包信息
	rpInfo, err := redPacketInfo(appId, packetId)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]interface{})
	ret["coin"] = utility.ToInt(rlt["coin_id"])
	ret["total"] = utility.ToInt(rpInfo["amount"])
	ret["remain"] = utility.ToInt(rpInfo["remain"])
	ret["amount"] = utility.ToFloat32(rlt["amount"])

	//  红包websocket推送
	if packetInfo != nil && packetInfo.CType == types.RPToRoom {
		//领取的人
		receiverName := orm.GetMemberName(packetInfo.ToId, userId)
		//发红包的人
		senderName := orm.GetMemberName(packetInfo.ToId, packetInfo.UserId)
		message := receiverName + " 领取了" + senderName + " 的红包"
		members := []string{packetInfo.UserId}
		if userId != packetInfo.UserId {
			members = append(members, userId)
		}
		SendAlert(userId, packetInfo.ToId, types.ToRoom, members, types.Alert, proto.ComposeRPReceiveAlert(types.AlertReceiveRedpackage, userId, receiverName, packetInfo.UserId, senderName, packetInfo.PacketId, message))
	} else {
		//领取的人
		receiver, _ := orm.GetUserInfoById(userId)
		receiverName := ""
		if receiver != nil {
			receiverName = receiver.Username
		}
		//发红包的人
		sender, _ := orm.GetUserInfoById(packetInfo.UserId)
		senderName := ""
		if sender != nil {
			senderName = sender.Username
		}
		//
		var alertSender string
		message := receiverName + " 领取了" + senderName + " 的红包"
		members := []string{packetInfo.UserId}
		if userId == packetInfo.UserId {
			alertSender = packetInfo.ToId
		} else {
			alertSender = userId
			members = append(members, userId)
		}
		SendAlert(alertSender, packetInfo.UserId, types.ToUser, members, types.Alert, proto.ComposeRPReceiveAlert(types.AlertReceiveRedpackage, userId, receiverName, packetInfo.UserId, senderName, packetInfo.PacketId, message))
	}

	////  红包websocket推送
	//if packetInfo != nil && packetInfo.CType == types.RPToRoom {
	//	if userId == packetInfo.UserId {
	//		var membersReceiveRP = make([]string, 0) //收红包
	//		membersReceiveRP = append(membersReceiveRP, userId)
	//		message := "你领取了自己的红包"
	//		SendAlert(userId, packetInfo.ToId, types.ToRoom, membersReceiveRP, types.Alert, proto.ComposeAlert(types.AlertReceiveRedpackage, message))
	//	} else {
	//		var membersSendRP = make([]string, 0)    //发红包
	//		var membersReceiveRP = make([]string, 0) //收红包
	//		membersSendRP = append(membersSendRP, packetInfo.UserId)
	//		membersReceiveRP = append(membersReceiveRP, userId)
	//		//领取的人
	//		receiverName := orm.GetMemberName(packetInfo.ToId, userId)
	//		//发红包的人
	//		senderName := orm.GetMemberName(packetInfo.ToId, packetInfo.UserId)
	//		//发红包的人收到的消息
	//		message := receiverName + " 领取了你的红包"
	//		SendAlert(userId, packetInfo.ToId, types.ToRoom, membersSendRP, types.Alert, proto.ComposeAlert(types.AlertReceiveRedpackage, message))
	//		//领红包的人收到的消息
	//		message = "你领取了" + senderName + " 的红包"
	//		SendAlert(userId, packetInfo.ToId, types.ToRoom, membersReceiveRP, types.Alert, proto.ComposeAlert(types.AlertReceiveRedpackage, message))
	//
	//	}
	//}
	//todo   下个版本
	//else {
	//	var membersSendRP = make([]string, 0)//发红包
	//	var membersReceiveRP = make([]string, 0)//收红包
	//	membersSendRP = append(membersSendRP,packetInfo[0]["user_id"])
	//	membersReceiveRP = append(membersReceiveRP,userId)
	//	if userId == packetInfo[0]["user_id"]{
	//		message := "你领取了你的红包"
	//		SendAlert(userId, packetInfo[0]["to_id"], types.ToUser, membersSendRP, types.Alert, ComposeAlert(types.AlertReceiveRedpackage, message))
	//	}else {
	//		UserInfo,err := db.FindFriendInfoByUserId(packetInfo[0]["user_id"],userId)
	//		if err != nil {
	//			logFriend.Error("db FindFriendInfoByUserId", err)
	//			return result.QueryDbFailed, nil, err
	//		}
	//		FriendInfo,err := db.FindFriendInfoByUserId(userId,packetInfo[0]["user_id"])
	//		if err != nil {
	//			logFriend.Error("db FindFriendInfoByUserId", err)
	//			return result.QueryDbFailed, nil, err
	//		}
	//		//发红包的人收到的消息
	//		message := UserInfo[0]["username"] + " 领取了你的红包"
	//		SendAlert(userId, packetInfo[0]["to_id"], types.ToUser, membersSendRP, types.Alert, ComposeAlert(types.AlertReceiveRedpackage, message))
	//		//领红包的人收到的消息
	//		message = "你领取了" + FriendInfo[0]["username"] + " 的红包"
	//		SendAlert(userId, packetInfo[0]["to_id"], types.ToUser, membersReceiveRP, types.Alert, ComposeAlert(types.AlertReceiveRedpackage, message))
	//	}
	//}
	return ret, nil
}

func formatRPUserInfo(appId, caller, uid string, cType int, roomId string) *types.User {
	receiver, _ := orm.GetUserInfoByUid(appId, uid)
	if receiver == nil {
		return nil
	}
	var name = receiver.Username
	switch cType {
	case types.RPToRoom:
		friend, _ := orm.FindFriendById(caller, receiver.UserId)
		if friend != nil && friend.Remark != "" {
			name = friend.Remark
		} else {
			nickname := orm.GetMemberName(roomId, receiver.UserId)
			if nickname != "" {
				name = nickname
			}
		}
	case types.RPToUser:
		friend, _ := orm.FindFriendById(caller, receiver.UserId)
		if friend != nil && friend.Remark != "" {
			name = friend.Remark
		}
	}
	receiver.Username = name
	return receiver
}

// 获取红包详情
func RedEnvelopeDetail(appId, caller, token string, packetId string) (*types.RedPacketInfo, error) {
	//从红包服务获取红包信息
	rpInfo, err := redPacketInfo(appId, packetId)
	if err != nil {
		return nil, err
	}

	revInfo, err := receiveSingleRecord(appId, token, packetId)
	if err != nil {
		return nil, err
	}

	var cType int
	var toId string
	if p, _ := getRedPacketInfo(packetId); p != nil {
		cType = p.CType
		toId = p.ToId
	}
	//接收者信息
	var rev *types.RPReceiveInfo
	if revInfo != nil {
		rev = &types.RPReceiveInfo{}
		receiver := formatRPUserInfo(appId, caller, utility.ToString(revInfo["user_id"]), cType, toId)
		if receiver != nil {
			rev.UserId = receiver.UserId
			rev.UserName = receiver.Username
			rev.UserAvatar = receiver.Avatar
		} else {
			rev.UserId = ""
			rev.UserName = utility.ToString(revInfo["user_name"])
			rev.UserAvatar = ""
		}
		rev.CoinId = utility.ToInt(revInfo["coin_id"])
		rev.CoinName = utility.ToString(revInfo["coin_name"])
		rev.Amount = utility.ToFloat64(revInfo["amount"])
		rev.CreatedAt = utility.ToInt64(revInfo["create_at"]) * 1000
		rev.Status = utility.ToInt(revInfo["status"])
		rev.FailMessage = utility.ToString(revInfo["fail_message"])
	}

	var ret types.RedPacketInfo
	sender := formatRPUserInfo(appId, caller, utility.ToString(rpInfo["user_id"]), cType, toId)
	if sender != nil {
		ret.SenderId = sender.UserId
		ret.SenderAvatar = sender.Avatar
		ret.SenderUid = sender.Uid
		ret.SenderName = sender.Username
		ret.IdentificationInfo = sender.IdentificationInfo
		ret.Identification = sender.Identification
	} else {
		ret.SenderId = ""
		ret.SenderName = utility.ToString(rpInfo["user_name"])
		ret.SenderAvatar = ""
		ret.SenderUid = ""
	}
	ret.Type = utility.ToInt(rpInfo["type"])
	ret.CoinId = utility.ToInt(rpInfo["coin_id"])
	ret.CoinName = utility.ToString(rpInfo["coin_name"])
	ret.Amount = utility.ToFloat64(rpInfo["amount"])
	ret.Size = utility.ToInt(rpInfo["size"])
	ret.ToUsers = utility.ToString(rpInfo["to"])
	ret.Remain = utility.ToInt(rpInfo["remain"])
	ret.Remark = utility.ToString(rpInfo["remark"])
	ret.Status = utility.ToInt(rpInfo["status"])
	ret.CreatedAt = utility.ToInt64(rpInfo["create_at"]) * 1000
	ret.PacketId = packetId
	ret.PacketUrl = getRPH5Url(appId, ret.SenderId, packetId)
	ret.RevInfo = rev

	return &ret, nil
}

//红包接收详情详情
func GetPacketRecvDetails(appId, caller, packetId string) ([]*types.RPReceiveInfo, error) {
	b, err := json.Marshal(&types.ReqBase{
		Method: "ReceiveDetail",
		Params: &types.ReceiveDetailParams{
			PacketId: packetId,
		},
	})
	if err != nil {
		result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}

	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, nil, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	rlt, _ := checkRPResponse(byte)
	var revRows = make([]*types.RPReceiveInfo, 0)
	if rlt == nil {
		return revRows, nil
	}
	if rows, ok := rlt["rows"].([]interface{}); ok {
		//从数据库查询红包信息
		var cType int
		var toId string
		if p, _ := getRedPacketInfo(packetId); p != nil {
			cType = p.CType
			toId = p.ToId
		}
		for _, row := range rows {
			item := row.(map[string]interface{})
			receiver := formatRPUserInfo(appId, caller, utility.ToString(item["user_id"]), cType, toId)
			var rev types.RPReceiveInfo
			if receiver != nil {
				rev.UserId = receiver.UserId
				rev.UserName = receiver.Username
				rev.UserAvatar = receiver.Avatar
			} else {
				rev.UserId = ""
				rev.UserName = utility.ToString(item["user_name"])
				rev.UserAvatar = ""
			}
			rev.CoinId = utility.ToInt(item["coin_id"])
			rev.CoinName = utility.ToString(item["coin_name"])
			rev.Amount = utility.ToFloat64(item["amount"])
			rev.CreatedAt = utility.ToInt64(item["create_at"]) * 1000
			rev.Status = utility.ToInt(item["status"])
			rev.FailMessage = utility.ToString(item["fail_message"])

			revRows = append(revRows, &rev)
		}
	}
	return revRows, nil
	/*
		{
			"id": int,
			"result": {
				"rows" : [
					{
						"platform_id": "1001",
						"user_id": "200601",
						"user_name": "zhangsan",
						"user_mobile": "132****1234",
						"amount": 26,
						"coin_id": 2,
						"coin_name": "BTC",
						"created_at": 1530003110,
						"status": 1,
						"fail_message":""
					},
					...
				]
			}
		}
	*/
}

//从红包服务获取统计信息
func statisticFromRPServer(appId, token string, operation, coin, rpType, pageNum, pageSize int, startTime, endTime int64) (map[string]interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		logRedpacket.Warn("statisticFromRPServer", "warn", types.ERR_APPNOTFIND.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	headers := make(map[string]string)
	headers["Authorization"] = types.AuthType + " " + token
	headers["Platform"] = app.RPPid

	b, err := json.Marshal(&types.ReqBase{
		Method: "Statistic",
		Params: &types.StatisticParams{
			Operation: operation,
			CoinId:    coin,
			Type:      rpType,
			StartTime: startTime,
			EndTime:   endTime,
			PageNumer: pageNum,
			PageSize:  pageSize,
		},
	})
	if err != nil {
		logRedpacket.Error("statisticFromRPServer json Marshal", "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("参数解析错误")
	}
	body := bytes.NewBuffer(b)
	byte, err := http.HTTPPostJSON(app.RPServer, headers, body)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:36:33
		logRedpacket.Warn(fmt.Sprintf("http client request %v failed", app.RPServer), "err", err.Error())
		return nil, result.NewError(result.RPError).JustShowExtMsg().SetExtMessage("访问红包服务失败")
	}

	ret, err := checkRPResponse(byte)
	if err != nil {
		return nil, result.NewError(err.(*respError).ToChatErrorCode()).SetExtMessage(err.(*respError).ToChatExtMsg()).
			SetChildErr(result.ServiceRP, err.(*respError).ErrorCode, err.(*respError).Message)
	}
	return ret, err
	/*
		{
			"id": 0,
			"result": {
				"total_count": 9,
				"sum": -1,
				"count": 9,
				"coin_id": 12,
				"coin_name": "YCC",
				"operation": 0,
				"rows": [{
					"operation": 2,
					"packet_id": "65175a60-a1f9-4716-b3f2-faf48376ca75",
					"platform_id": "1000",
					"coin_id": 12,
					"coin_name": "YCC",
					"sender_id": "34",
					"sender_name": "",
					"sender_mobile": "157****6517",
					"receiver_id": "34",
					"receiver_mobile": "157****6517",
					"amount": 1,
					"status": 2,
					"type": 1,
					"size": 1,
					"remain": 0,
					"create_at": 1552619581,
					"update_at": 1552619582,
					"expire_at": -62135596800
				}]
			},
			"error": {
				"code": 0,
				"message": ""
			}
		}
	*/
}

// 获取红包统计信息
func RPStatistic(appId, caller, token string, operation, coin, rpType, pageNum, pageSize int, startTime, endTime int64) (interface{}, error) {
	sInfo, err := statisticFromRPServer(appId, token, operation, coin, rpType, pageNum, pageSize, startTime, endTime)
	if err != nil {
		return nil, err
	}
	count := utility.ToInt(sInfo["count"])
	sum := utility.ToFloat64(sInfo["sum"])
	coinName := utility.ToString(sInfo["coin_name"])
	coinId := utility.ToInt(sInfo["coin_id"])
	rpInfo := make([]*types.RedPacketInfo, 0)
	for _, row := range sInfo["rows"].([]interface{}) {
		if r, ok := row.(map[string]interface{}); ok {
			//根据红包id从数据库查询
			var cType int
			var toId string
			if p, _ := getRedPacketInfo(utility.ToString(r["packet_id"])); p != nil {
				cType = p.CType
				toId = p.ToId
			}
			//接收者信息
			var rev *types.RPReceiveInfo
			revInfo, _ := receiveSingleRecord(appId, token, utility.ToString(r["packet_id"]))
			if revInfo != nil {
				rev = &types.RPReceiveInfo{}
				receiver := formatRPUserInfo(appId, caller, utility.ToString(revInfo["user_id"]), cType, toId)
				if receiver != nil {
					rev.UserId = receiver.UserId
					rev.UserName = receiver.Username
					rev.UserAvatar = receiver.Avatar
				} else {
					rev.UserId = ""
					rev.UserName = utility.ToString(revInfo["user_name"])
					rev.UserAvatar = ""
				}
				rev.CoinId = utility.ToInt(revInfo["coin_id"])
				rev.CoinName = utility.ToString(revInfo["coin_name"])
				rev.Amount = utility.ToFloat64(revInfo["amount"])
				rev.CreatedAt = utility.ToInt64(revInfo["create_at"]) * 1000
				rev.Status = utility.ToInt(revInfo["status"])
				rev.FailMessage = utility.ToString(revInfo["fail_message"])
			}

			//发送者信息
			var ret types.RedPacketInfo
			sender := formatRPUserInfo(appId, caller, utility.ToString(r["sender_id"]), cType, toId)
			if sender != nil {
				ret.SenderId = sender.UserId
				ret.SenderAvatar = sender.Avatar
				ret.SenderUid = sender.Uid
				ret.SenderName = sender.Username
				ret.IdentificationInfo = sender.IdentificationInfo
				ret.Identification = sender.Identification
			} else {
				ret.SenderId = ""
				ret.SenderName = utility.ToString(r["sender_name"])
				ret.SenderAvatar = ""
				ret.SenderUid = ""
			}
			ret.Type = utility.ToInt(r["type"])
			ret.CoinId = utility.ToInt(r["coin_id"])
			ret.CoinName = utility.ToString(r["coin_name"])
			ret.Amount = utility.ToFloat64(r["amount"])
			ret.Size = utility.ToInt(r["size"])
			ret.Remain = utility.ToInt(r["remain"])
			ret.Remark = utility.ToString(r["remark"])
			ret.Status = utility.ToInt(r["status"])
			ret.CreatedAt = utility.ToInt64(r["create_at"]) * 1000
			ret.PacketId = utility.ToString(r["packet_id"])
			ret.PacketUrl = getRPH5Url(appId, ret.SenderId, utility.ToString(r["packet_id"]))
			ret.RevInfo = rev

			rpInfo = append(rpInfo, &ret)
		}
	}

	var ret types.RPStatisticInfo
	ret.Count = count
	ret.Sum = sum
	ret.CoinId = coinId
	ret.CoinName = coinName
	ret.RedPackets = rpInfo

	return ret, nil
}
