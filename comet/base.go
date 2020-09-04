package comet

import (
	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/inconshreveable/log15"
)

var logBase = log15.New("module", "comet/base")

func init() {
	receiver := proto.GetProtocolReceiver()
	receiver[types.EventCommonMsg] = doReceiveNomal
	receiver[types.EventSyncMsg] = doReceiveSyncEvent
	receiver[types.EventUpdateSKey] = doReceiveUpdateSKeyEvent
	//receiver[types.EventUpdatePublicKey] = doReceiveUpdatePubKeyEvent
	receiver[types.EventStartAck] = doReceiveStartAckEvent
	receiver[types.EventStartGetAllSKey] = doGetAllSKeyEvent
}

func CreateProto(conn router.Connection, msg []byte) {
	iProto, err := proto.CreateProto(msg)
	if err != nil {
		return
	}
	if iProto == nil {
		//wsLog.Warn("can not find IBaseProto handler")
	}
	iProto.Receive(conn, msg)
}

//调度器
type dispatcher interface {
	//拦截
	intercept(args ...interface{}) error
	//添加聊天日志
	appendLog() error
	//
	pushMsg() error
}

func createDispatcher(cType int, p *proto.Proto, userId string, time int64) dispatcher {
	switch cType {
	case types.ToRoom:
		return &roomDispatcher{
			msgParse: p,
			userId:   userId,
			msgTime:  time,
		}
	case types.ToUser:
		return &userDispatcher{
			msgParse: p,
			userId:   userId,
			msgTime:  time,
		}
	default:
	}
	return nil
}

type roomDispatcher struct {
	msgParse *proto.Proto
	userId   string
	msgTime  int64
}

type userDispatcher struct {
	msgParse *proto.Proto
	userId   string
	msgTime  int64
}

func (t *roomDispatcher) intercept(args ...interface{}) error {
	msgId := t.msgParse.GetMsgId()
	targetId := t.msgParse.GetTargetId()
	msgType := t.msgParse.GetMsgType()
	userId := t.userId

	//检查用户是否封禁
	u, err := orm.GetUserInfoById(userId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if u.CloseUntil > utility.NowMillionSecond() {
		logBase.Warn("room intercept", "warn", "AccountClosedByAdmin", "userId", userId, "close until", u.CloseUntil)
		return result.NewError(result.AccountClosedByAdmin).JustShowExtMsg().SetExtMessage(model.ConvertUserClosedAlertStr(u.AppId, u.CloseUntil))
	}

	//检查群是否封禁
	r, err := orm.FindRoomById(targetId, types.RoomNotDeleted)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if r.CloseUntil > utility.NowMillionSecond() {
		logBase.Warn("room intercept", "warn", "RoomClosedByAdmin", "userId", userId, "roomId", targetId, "close until", u.CloseUntil)
		return result.NewError(result.RoomClosedByAdmin).JustShowExtMsg().SetExtMessage(model.ConvertRoomClosedAlertStr(u.AppId, r.CloseUntil))
	}

	/*//判断群性质 1加密群 2非加密群
	if r.Encrypt == types.IsEncrypt && t.msgParse.IsEncryptedMsg(){
		return &result.Error{ErrorCode: result.MustEncryptMsgToRoom}
	}*/

	level := orm.GetMemberLevel(targetId, userId, types.RoomUserNotDeleted) //model.GetRoomUserLevel(targetId, userId, types.RoomUserNotDeleted)
	if level == types.RoomLevelNotExist {
		logBase.Warn("room intercept", "warn", "UserIsNotInRoom", "roomId", targetId, "userId", userId)
		return result.NewError(result.UserIsNotInRoom)
	}
	// 普通成员
	if level == types.RoomLevelNomal {
		if msgType == types.System {
			logBase.Warn("room intercept", "warn", "PermissionDeny:normal member can not send system msg", "roomId", targetId, "userId", userId)
			return result.NewError(result.PermissionDeny)
		}
		//判断是否禁言
		isMuted, err := model.CheckMemberMuted(targetId, userId) //model.GetRoomUserMuted(targetId, userId)
		if err != nil {
			return err
		}
		if isMuted {
			logBase.Warn("room intercept", "warn", "UserMuted:user have been muted in room", "roomId", targetId, "userId", userId)
			return result.NewError(result.UserMuted)
		}
	}

	//消息重复拦截
	msg, err := orm.FindRoomChatLogByMsgId(userId, msgId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if msg != nil {
		logBase.Debug("Repeated message", "logId", msg.Id, "msgId", msg.MsgId)
		//重复消息
		userId := t.userId
		logId := msg.Id
		sendTime := msg.Datetime
		t.msgParse.SetLogId(logId)
		//查询用户信息
		u, _ := orm.GetUserInfoById(userId)
		if u == nil {
			return result.NewError(result.DbConnectFail)
		}
		//生成待发送字节消息
		resp, err2 := t.msgParse.WrapResp(u, sendTime)
		if err2 != nil {
			logBase.Error("room intercept: WrapResp json Marshal", "err", err2.Error())
			return result.NewError(result.ServerInterError)
		}
		//将消息发给自己
		if user, _ := router.GetUser(userId); user != nil {
			user.SendToAllClients(resp)
			logBase.Debug("send back private msg success")
		}
		return types.ERR_REPEAT_MSG
	}
	return nil
}

func (t *roomDispatcher) appendLog() error {
	userId := t.userId
	//add room chat log content
	err := t.msgParse.AppendChatLog(userId, types.NotRead, t.msgTime)
	if err == nil && t.msgParse.IsSpecial() {
		//需要存储接收情况
		go func() {
			//logic.AppendReceiveLog(t.msgParse.GetTargetId(), t.msgParse.GetLogId(), types.NotRead) //model.AppendRoomRevLog(t.msgParse.GetTargetId(), t.msgParse.GetLogId(), types.NotRead) u.UserId
			members, err := orm.FindNotDelMembers(t.msgParse.GetTargetId(), types.SearchAll)
			if err != nil {
				return
			}
			for _, u := range members {
				_, err := orm.AppendMemberRevLog(u.RoomMember.UserId, t.msgParse.GetLogId(), types.NotRead)
				if err != nil {
					logBase.Warn("appendLog", "warn", "db err")
				}
			}
		}()
	}
	return err
}

func (t *roomDispatcher) pushMsg() error {
	targetId := t.msgParse.GetTargetId()
	channelId := t.msgParse.GetRouter()
	userId := t.userId

	//查询用户信息
	u, _ := orm.GetUserInfoById(userId)
	if u == nil {
		return result.NewError(result.DbConnectFail)
	}
	//生成待发送字节消息
	resp, err2 := t.msgParse.WrapResp(u, t.msgTime)
	if err2 != nil {
		logBase.Error("room pushMsg: WrapResp json Marshal", "err", err2.Error())
		return result.NewError(result.ServerInterError)
	}
	//个推 推送通知
	go func() {
		err := model.PushToRoom(u.AppId, userId, targetId, t.msgParse)
		if err != nil {
			logBase.Warn("room pushMsg", "warn", "push to room failed", "error", err)
		}
		/*unConnectUsers := model.GetRoomOfflineMember(targetId)                    //unConnectUsers := model.GetNotStayConnectedRoomUsers(targetId)
		model.PushMsgToOfflineMember(u.AppId, userId, t.msgParse, unConnectUsers) //model.SendToRoomUnConnect(userId, t.msgParse, unConnectUsers)*/
	}()
	//do broadcast
	if ch, ok := router.GetChannel(channelId); !ok {
		logBase.Warn("room pushMsg", "warn", "RoomNotExists", "channelId", channelId)
		return result.NewError(result.RoomNotExists)
	} else {
		if !router.IsUserOnline(userId) {
			logBase.Error("room pushMsg", "warn", "UserIsNotOnline", "channelId", channelId, "userId", userId)
			return result.NewError(result.UserIsNotOnline)
		}
		logBase.Debug("Send to room Compose data success")
		ch.Broadcast(resp)
	}
	return nil
}

func (t *userDispatcher) intercept(args ...interface{}) error {
	msgId := t.msgParse.GetMsgId()
	targetId := t.msgParse.GetTargetId()
	userId := t.userId
	appId := args[0].(string)
	//检查用户是否封禁
	u, err := orm.GetUserInfoById(userId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if u.CloseUntil > utility.NowMillionSecond() {
		logBase.Warn("user intercept", "warn", "AccountClosedByAdmin", "userId", userId, "close until", u.CloseUntil)
		return result.NewError(result.AccountClosedByAdmin).JustShowExtMsg().SetExtMessage(model.ConvertUserClosedAlertStr(u.AppId, u.CloseUntil))
	}

	app := app.GetApp(appId)
	if app.IsOtc == types.IsOtc {
		orderId := ""
		isOtcUid := false
		//获取交易订单号
		ext := t.msgParse.GetExt()
		if e, ok := ext.(map[string]interface{}); ok {
			orderId = utility.ToString(e["orderId"])
			isOtcUid = utility.ToBool(e["idType"])
		}

		user, err := orm.GetUserInfoById(userId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if user == nil {
			logBase.Warn("user intercept", "warn", "UserNotExists", "userId", userId)
			return result.NewError(result.UserNotExists)
		}
		userUid := user.Uid
		targetUid := targetId
		if !isOtcUid {
			target, err := orm.GetUserInfoById(targetId)
			if err != nil {
				return result.NewError(result.DbConnectFail)
			}
			if target == nil {
				logBase.Warn("user intercept", "warn", "Target UserNotExists", "userId", userId, "targetId", targetId)
				return result.NewError(result.UserNotExists)
			}
			targetUid = target.Uid
		}
		//检查交易是否存在
		if !model.CheckOrder(app.OtcServer, userUid, targetUid, orderId) {
			logBase.Warn("user intercept", "warn", "TradeRelationshipNotExt", "OtcServer", app.OtcServer, "userUid", userUid, "targetUid", targetUid, "orderId", orderId)
			return result.NewError(result.TradeRelationshipNotExt)
		}

		//如果传的是otc uid
		if isOtcUid {
			//查询用户是否存在
			u, err := orm.GetUserInfoByUid(appId, targetId)
			if err != nil {
				return result.NewError(result.DbConnectFail)
			}
			if u != nil {
				t.msgParse.SetTargetId(u.UserId)
			} else {
				//用户不存在，则临时创建用户
				//新建用户
				username := utility.RandomUsername()
				userLevel := types.LevelMember
				markId := appId + targetId

				id, err := orm.InsertUser(markId, targetId, appId, username, "", "", "", "", utility.ToString(userLevel), "0", "", "", "", "", utility.NowMillionSecond())
				if err != nil {
					return result.NewError(result.TokenLoginFailed)
				}
				userId = utility.ToString(id)
				t.msgParse.SetTargetId(userId)
			}
		}
	}
	//else {
	//	//检查是否是好友
	//	ok, err := orm.CheckIsFriend(userId, targetId, types.FriendIsNotDelete)
	//	if err != nil {
	//		return result.NewError(result.DbConnectFail)
	//	}
	//	if !ok {
	//		logBase.Warn("user intercept", "warn", "IsNotFriend", "you", userId, "friend", targetId)
	//		return result.NewError(result.IsNotFriend)
	//	}
	//
	//	//检查黑名单用户
	//	if ok, err := model.CheckIsBlocked(targetId, userId); ok || err != nil {
	//		logBase.Debug("user intercept", "warn", "IsBlocked", "you", userId, "friend", targetId)
	//		return result.NewError(result.IsBlocked)
	//	}
	//}

	//消息重复拦截
	msg, err := orm.FindPrivateChatLogByMsgId(userId, msgId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if msg != nil {
		logBase.Debug("Repeated message", "logId", msg.Id, "msgId", msg.MsgId)
		//重复消息
		userId := t.userId
		logId := msg.Id
		sendTime := msg.SendTime
		t.msgParse.SetLogId(logId)
		//查询用户信息
		u, err := orm.GetUserInfoById(userId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if u == nil {
			logBase.Warn("user intercept", "warn", "UserNotExists", "userId", userId)
			return result.NewError(result.UserNotExists)
		}
		//生成待发送字节消息
		resp, err2 := t.msgParse.WrapResp(u, sendTime)
		if err2 != nil {
			logBase.Error("user intercept: WrapResp json Marshal", "err", err2.Error())
			return result.NewError(result.ServerInterError)
		}
		//将消息发给自己
		if user, _ := router.GetUser(userId); user != nil {
			user.SendToAllClients(resp)
			logBase.Debug("send back private msg success")
		}
		return types.ERR_REPEAT_MSG
	}
	return nil
}

func (t *userDispatcher) appendLog() error {
	userId := t.userId
	// Append private chat log ,log state is 2(not send success)
	return t.msgParse.AppendChatLog(userId, types.NotRead, t.msgTime)
}

func (t *userDispatcher) pushMsg() error {
	friendId := t.msgParse.GetRouter()
	userId := t.userId
	text := t.msgParse.GetGTMsg(userId)
	//查询用户信息
	u, _ := orm.GetUserInfoById(userId)
	if u == nil {
		return result.NewError(result.DbConnectFail)
	}
	//生成待发送字节消息
	resp, err2 := t.msgParse.WrapResp(u, t.msgTime)
	if err2 != nil {
		logBase.Error("user pushMsg: WrapResp json Marshal", "err", err2.Error())
		return result.NewError(result.ServerInterError)
	}

	if friend, ok := router.GetUser(friendId); ok && friend != nil && len(friend.GetClients()) != 0 {
		logBase.Debug("send online single")
		friend.SendToAllClients(resp)
	} else if text != "" {
		logBase.Debug("send u-push single", "text", text)
		err := model.PushToFriend(u.AppId, userId, friendId, text, t.msgParse)
		if err != nil {
			logBase.Warn("push to user failed", "error", err)
		}
	}
	//将消息发给自己
	if user, _ := router.GetUser(userId); user != nil {
		user.SendToAllClients(resp)
		logBase.Debug("send back private msg success")
	}
	return nil
}
