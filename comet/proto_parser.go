package comet

import (
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	l "github.com/inconshreveable/log15"
)

var logParser = l.New("module", "comet/proto_parser")

func doReceiveNomal(ibp proto.IBaseProto, wsc router.Connection, msg []byte) {
	msgParse := ibp.(*proto.Proto)
	var err error
	defer func() {
		if err != nil {
			switch e := err.(type) {
			case *result.Error:
				err = &result.WsError{
					E:         e,
					EventCode: types.EventCommonMsg,
					MsgId:     msgParse.GetMsgId(),
				}
			default:
				return
			}
			ret, err := result.ComposeWsError(err)
			if err == nil {
				wsc.WriteResponse(ret)
			} else {
				logParser.Error(err.Error())
			}
		}
	}()

	err = msgParse.FromBytes(msg)
	if err != nil {
		return
	}

	msgTime := utility.NowMillionSecond()

	args := wsc.Args().(map[string]string)
	userId := args["userId"]
	device := args["device"]
	appId := args["appId"]

	dispatcher := createDispatcher(msgParse.GetChannelType(), msgParse, userId, msgTime)
	if dispatcher == nil {
		//错误
		err = result.NewError(result.MsgFormatError)
		return
	}

	err = dispatcher.intercept(appId)
	if err != nil {
		if err == types.ERR_REPEAT_MSG {
			err = nil
		}
		return
	}

	err = dispatcher.appendLog()
	if err != nil {
		return
	}

	err = dispatcher.pushMsg()
	if err != nil {
		return
	}
	logParser.Info("rev msg", "raw", msgParse.GetMsg(), "userId", userId, "device", device)
}

func doReceiveSyncEvent(ibp proto.IBaseProto, wsc router.Connection, msg []byte) {
	eveParse := ibp.(*proto.EventSyncMsg)
	var err error
	defer func() {
		if err != nil {
			switch e := err.(type) {
			case *result.Error:
				err = &result.WsError{
					E:         e,
					EventCode: types.EventSyncMsg,
					MsgId:     "",
				}
			default:
				return
			}
			ret, err := result.ComposeWsError(err)
			if err == nil {
				wsc.WriteResponse(ret)
			} else {
				logParser.Error(err.Error())
			}
		}
	}()

	err = eveParse.FromBytes(msg)
	if err != nil {
		return
	}

	//开始同步
	args := wsc.Args().(map[string]string)
	clientId := args["clientId"]
	userId := args["userId"]
	u, ok := router.GetUser(userId)
	if !ok {
		logParser.Error("doReceiveSyncEvent", "err", "UserNotExists", "userId", userId)
		err = result.NewError(result.UserNotExists)
		return
	}
	client, ok := u.GetClient(clientId)

	logParser.Info("[batch] start push msg", "userId", u.Id)
	go model.BatchPushRecentMsg(client, eveParse.GetTime())
}

//更新群会话秘钥
func doReceiveUpdateSKeyEvent(ibp proto.IBaseProto, wsc router.Connection, msg []byte) {
	eveParse := ibp.(*proto.EventUpdateSKey)
	var err error
	defer func() {
		if err != nil {
			switch e := err.(type) {
			case *result.Error:
				err = &result.WsError{
					E:         e,
					EventCode: types.EventUpdateSKey,
					MsgId:     "",
				}
			default:
				return
			}
			ret, err := result.ComposeWsError(err)
			if err == nil {
				wsc.WriteResponse(ret)
			} else {
				logParser.Error(err.Error())
			}
		}
	}()

	err = eveParse.FromBytes(msg)
	if err != nil {
		logParser.Debug("消息格式错误：", "msg", string(msg))
		return
	}

	args := wsc.Args().(map[string]string)
	userId := args["userId"]
	clientId := args["clientId"]
	device := args["device"]
	kid := utility.ToString(utility.NowMillionSecond())

	//非加密群不更新群会话秘钥 dld-v2.8.4 2019年7月3日14:36:09
	roomId := eveParse.GetRoomId()
	secrets := eveParse.GetSecret()
	fromKey := eveParse.GetFromKey()

	room, err2 := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err2 != nil {
		logParser.Error("更新群会话秘钥 查找群信息失败")
		return
	}
	if room.Encrypt == types.IsNotEncrypt {
		return
	}

	//TODO 用于测试
	number, err := orm.GetMemberNumber(roomId)
	if err != nil {
	}
	if int(number) != len(secrets) {
		logParser.Info("Test Log", "roomId", roomId, "members", number, "secrets", len(secrets), "update key user", userId, "device", device, "clientId", clientId)
	}

	go func() {
		if fromKey == "" {
			//获取公钥
			user, _ := orm.GetUserInfoById(userId)
			fromKey = user.PublicKey
		}

		for _, s := range secrets {
			model.SendAlert(userId, s.GetUserId(), types.ToUser, []string{s.GetUserId()}, types.Alert, proto.ComposeUpdateSKeyAlert(roomId, fromKey, kid, s.GetKey()))
		}
	}()
}

//同步群会话秘钥
func doGetAllSKeyEvent(ibp proto.IBaseProto, wsc router.Connection, msg []byte) {
	eveParse := ibp.(*proto.EventGetSKey)
	var err error
	defer func() {
		if err != nil {
			switch e := err.(type) {
			case *result.Error:
				err = &result.WsError{
					E:         e,
					EventCode: types.EventStartGetAllSKey,
					MsgId:     "",
				}
			default:
				return
			}
			ret, err := result.ComposeWsError(err)
			if err == nil {
				wsc.WriteResponse(ret)
			} else {
				logParser.Error(err.Error())
			}
		}
	}()

	err = eveParse.FromBytes(msg)
	if err != nil {
		return
	}

	args := wsc.Args().(map[string]string)
	userId := args["userId"]
	clientId := args["clientId"]

	u, ok := router.GetUser(userId)
	if !ok {
		logParser.Error("doGetAllSKeyEvent", "err", "UserNotExists", "userId", userId)
		err = result.NewError(result.UserNotExists)
		return
	}
	client, ok := u.GetClient(clientId)

	go model.BatchSKeyMsg(client, eveParse.GetDatetime())
}

//更新用户公钥
func doReceiveUpdatePubKeyEvent(ibp proto.IBaseProto, wsc router.Connection, msg []byte) {
	eveParse := ibp.(*proto.EventUpdatePubKey)
	var err error
	defer func() {
		if err != nil {
			switch e := err.(type) {
			case *result.Error:
				err = &result.WsError{
					E:         e,
					EventCode: types.EventUpdatePublicKey,
					MsgId:     "",
				}
			default:
				return
			}
			ret, err := result.ComposeWsError(err)
			if err == nil {
				wsc.WriteResponse(ret)
			} else {
				logParser.Error(err.Error())
			}
		}
	}()

	err = eveParse.FromBytes(msg)
	if err != nil {
		return
	}

	args := wsc.Args().(map[string]string)
	userId := args["userId"]

	friends, _ := orm.FindFriendsId(userId)

	var fs []string
	fs = append(fs, friends...)

	if eveParse.PublicKey != "" && eveParse.PrivateKey != "" {
		err := orm.UpdatePublicKey(userId, eveParse.PublicKey, eveParse.PrivateKey)
		if err != nil {
			logParser.Warn("update user public key filed", "userId", userId)
		}
		proto.SendUpdatePubKeyNotification(userId, eveParse.PublicKey, eveParse.PrivateKey, fs)
	} else {
		logParser.Warn("can not update empty public key")
	}
}

//开始确认
func doReceiveStartAckEvent(ibp proto.IBaseProto, wsc router.Connection, msg []byte) {
	eveParse := ibp.(*proto.EventAck)
	var err error
	defer func() {
		if err != nil {
			switch e := err.(type) {
			case *result.Error:
				err = &result.WsError{
					E:         e,
					EventCode: types.EventStartAck,
					MsgId:     "",
				}
			default:
				return
			}
			ret, err := result.ComposeWsError(err)
			if err == nil {
				wsc.WriteResponse(ret)
			} else {
				logParser.Error(err.Error())
			}
		}
	}()

	err = eveParse.FromBytes(msg)
	if err != nil {
		return
	}

	args := wsc.Args().(map[string]string)
	userId := args["userId"]
	clientId := args["clientId"]

	u, ok := router.GetUser(userId)
	if !ok {
		logParser.Error("doReceiveStartAckEvent", "err", "UserNotExists", "userId", userId)
		err = result.NewError(result.UserNotExists)
		return
	}
	client, ok := u.GetClient(clientId)
	go model.BatchPushAckMsg(client, eveParse.GetBegin(), eveParse.GetEnd(), eveParse.GetTotal())
}
