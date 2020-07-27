package model

import (
	"time"

	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logStatistics = log15.New("module", "model/statistics")

type SearchInfo struct {
	Type     int         `json:"type"`
	RoomInfo interface{} `json:"roomInfo"`
	UserInfo interface{} `json:"userInfo"`
}

func applyInfo(caller string, info *types.Apply) (*types.ApplyList, error) {
	var applyInfo types.ApplyList
	applyInfo.Id = info.Id
	applyInfo.Type = info.Type
	applyInfo.ApplyReason = info.ApplyReason
	applyInfo.Status = info.State
	applyInfo.Datetime = info.Datetime

	targetId := info.Target
	applyUser := info.ApplyUser
	sourceFromDb := info.Source
	//get sender Info
	senderInfo, err := ApplyFriendInfo(caller, applyUser)
	if err != nil {
		return &applyInfo, err
	}
	applyInfo.SenderInfo = senderInfo

	//get target Info
	if applyInfo.Type == types.IsRoom {
		targetInfo, err := ApplyRoomInfo(targetId)
		if err != nil {
			return &applyInfo, err
		}
		applyInfo.ReceiveInfo = targetInfo

		uInfo, err := orm.GetUserInfoById(caller)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		// 当type 为room时候targetInfo表示Room有关信息 申请加入的群为targetInfo.Id
		source, _ := getJoinSource(uInfo.AppId, targetInfo.Id, sourceFromDb)
		source2, _ := getJoinSourceV2(uInfo.AppId, targetInfo.Id, sourceFromDb)
		applyInfo.Source = source
		applyInfo.Source2 = source2
	} else if applyInfo.Type == types.IsFriend {
		targetInfo, err := ApplyFriendInfo(caller, targetId)
		if err != nil {
			return &applyInfo, err
		}
		applyInfo.ReceiveInfo = targetInfo
		source, err := ConverFriendSource(sourceFromDb, caller, applyInfo.SenderInfo.Id)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		source2, err := ConverFriendSourceV2(sourceFromDb, applyInfo.SenderInfo.Id)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		applyInfo.Source = source
		applyInfo.Source2 = source2
	}
	return &applyInfo, nil
}

func GetApplyInfoByLogId(caller string, logId int64) (*types.ApplyList, error) {
	info, err := orm.FindApplyLogById(logId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if info == nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	return applyInfo(caller, info)
}

func GetApplyList(caller, id string, number int) (interface{}, error) {
	nextId, logs, err := orm.FindApplyLogs(caller, utility.ToInt64(id), number)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var applyInfoList = make([]*types.ApplyList, 0)
	for _, v := range logs {
		info, err := applyInfo(caller, v)
		if err != nil {
			continue
		}
		applyInfoList = append(applyInfoList, info)
	}
	var ret = make(map[string]interface{})
	count, _ := orm.GetApplyListNumber(caller)

	ret["applyList"] = applyInfoList
	ret["nextId"] = utility.ToString(nextId)
	ret["totalNumber"] = count
	return ret, nil
}

func ClearlySearch(appId, caller, searchId string) (interface{}, error) {
	r, err := orm.FindRoomByMarkId(searchId, types.RoomNotDeleted)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var ret SearchInfo
	if r != nil {
		var room RoomBaseInfo
		{
			masterId := utility.ToString(r.MasterId)
			//查看appId
			user, err := orm.GetUserInfoById(masterId)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			roomAppId := ""
			if user != nil {
				roomAppId = utility.ToString(user.AppId)
			}
			if appId != roomAppId {
				var empty struct{}
				//return empty, &result.Error{ErrorCode: result.CodeOK, Message: ""}
				return empty, nil
			}
		}
		room.Id = r.Id
		room.MarkId = r.MarkId
		room.Name = r.Name
		room.Avatar = r.Avatar
		room.Encrypt = r.Encrypt
		room.CanAddFriend = r.CanAddFriend
		room.JoinPermission = r.JoinPermission
		room.RecordPermission = r.RecordPermision
		room.DisableDeadline = r.CloseUntil
		room.MemberNumber, _ = orm.GetMemberNumber(room.Id)
		room.ManagerNumber, _ = orm.GetRoomMasterNumber(room.Id)
		room.SystemMsg, _ = GetSystemMsg(room.Id, 0, 1)
		room.Identification = r.Identification
		room.IdentificationInfo = r.IdentificationInfo

		ret.Type = 1
		ret.RoomInfo = &room
		var empty struct{}
		ret.UserInfo = empty
		return &ret, nil
	}

	//find friend
	//find user by uid
	var user *types.User
	/*isPhone := utility.CheckPhoneNumber(searchId)
	if isPhone {
		user, err = orm.GetUserInfoByPhone(appId, searchId)
	} else {
		user, err = orm.GetUserInfoByMarkId(appId, searchId)
	}*/

	user, err = orm.GetUserInfoByPhone(appId, searchId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if user == nil {
		user, err = orm.GetUserInfoByMarkId(appId, searchId)
	}
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if user != nil {
		userId := user.UserId
		var friendInfo interface{}
		if caller == "" {
			friendInfo, err = UserInfo(userId)
		} else {
			friendInfo, err = FriendInfo(caller, userId)
		}
		if err != nil {
			return nil, err
		}
		if friendInfo != nil {
			ret.Type = 2
			ret.UserInfo = friendInfo
			var empty struct{}
			ret.RoomInfo = empty
			return &ret, nil
		}
	}
	var empty struct{}
	return empty, nil
}

//获取群会话秘钥
func RoomSessionKey(userId string, datetime int64) (interface{}, error) {
	logs, err := orm.FindSessionKeyAlert(userId, datetime)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var list = make([]*types.ChatLog, 0)
	for _, l := range logs {
		info, err := GetChatLogAsUser(userId, l)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		list = append(list, info)
	}

	var ret = make(map[string]interface{})
	ret["logs"] = list
	return ret, err
}

//群消息是否被焚毁
//返回：false:未焚毁;true 焚毁
func getRoomMsgHadBurnt(logId string) (bool, error) {
	count, err := orm.GetRoomMsgBurntNumber(logId)
	if err != nil {
		return false, result.NewError(result.DbConnectFail)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// 焚毁指定的一条消息
func ReadSnapMsg(userId, logId string, cType int) error {
	switch cType {
	case types.IsRoom:
		log, err := orm.FindRoomChatLogByContentId(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if log == nil {
			return nil //result.NewError(result.ParamsError)
		}
		isSnap := log.IsSnap
		roomId := log.RoomId
		senderId := log.SenderId
		//检查是否是 阅后即焚消息
		if isSnap == types.IsNotSnap {
			logStatistics.Warn("ReadSnapMsg", "warn", "ParamsError:room log is not snap log", "logId", logId)
			return result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "is not snap message")
		}
		//查出接收状态
		revInfo, err := orm.FindReceiveLogById(logId, userId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if revInfo == nil {
			return nil //result.NewError(result.ParamsError)
		}
		revLog := revInfo.Id
		isBurnt, err := getRoomMsgHadBurnt(logId)
		if err != nil {
			return err
		}
		err = orm.AlertRoomRevStateByRevId(revLog, types.HadBurnt)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if !isBurnt {
			//设置发送消息者已焚毁
			sRevInfo, err := orm.FindReceiveLogById(logId, senderId)
			if err != nil {
				return result.NewError(result.DbConnectFail)
			}
			if sRevInfo == nil {
				return nil //result.NewError(result.ParamsError)
			}
			senderRevLog := sRevInfo.Id
			err = orm.AlertRoomRevStateByRevId(senderRevLog, types.HadBurnt)
			if err != nil {
				return result.NewError(result.DbConnectFail)
			}
			defer SendAlert(userId, roomId, types.ToRoom, []string{senderId, userId}, types.Alert, proto.ComposeHadBurntAlert(types.ToRoom, logId))
		}
	case types.IsFriend:
		log, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if log == nil {
			return nil //result.NewError(result.ParamsError)
		}
		senderId := log.SenderId
		revId := log.ReceiveId
		isSnap := log.IsSnap
		//isDelete := utility.ToInt(info["is_delete"])
		//检查是否是 阅后即焚消息
		if isSnap == types.IsNotSnap {
			logStatistics.Warn("ReadSnapMsg", "warn", "ParamsError:friend log is not snap log", "logId", logId)
			return result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "is not snap message")
		}
		//如果消息接收对象不是该用户
		if revId != userId {
			return nil //result.NewError(result.ParamsError)
		}
		err = orm.AlertPrivateRevStateByRevId(logId, types.HadBurnt)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		defer SendAlert(userId, senderId, types.ToUser, []string{senderId, userId}, types.Alert, proto.ComposeHadBurntAlert(types.ToUser, logId))
	default:
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

// 转发消息
func ForwardMsg(callerId, sourceId string, channelType, forwardType int, logArray, targetRooms, targetUsers []string) (interface{}, error) {
	//找到源消息
	var fromName string
	//转发者信息
	var forwardUserName string
	var list = make([]*types.ChatLog, 0)
	var chType int
	switch channelType {
	case types.IsRoom:
		chType = types.ToRoom
		//查出消息来源群的信息
		room, err := orm.FindRoomById(sourceId, types.RoomDeletedOrNot)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if room != nil {
			fromName = room.Name
		}

		for _, v := range logArray {
			log, err := orm.FindRoomChatLogByContentId(v)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			info, err := GetChatLogAsRoom(callerId, log)
			if err != nil {
				return nil, err
			}
			//去掉合并转发中的Msg内容 防止无限转发导致的内容增长 2019年1月7日14:40:38 dld
			if info.MsgType == types.Forward {
				info.Msg = make(map[string]interface{})
			}
			list = append(list, info)
		}
	case types.IsFriend:
		chType = types.ToUser
		//查出消息来源者的信息
		user, err := orm.GetUserInfoById(sourceId)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if user != nil {
			fromName = user.Username
		}

		for _, v := range logArray {
			log, err := orm.FindPrivateChatLogById(v)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			info, err := GetChatLogAsUser(callerId, log)
			if err != nil {
				return nil, err
			}
			//去掉合并转发中的Msg内容 防止无限转发导致的内容增长 2019年1月7日14:40:38 dld
			if info.MsgType == types.Forward {
				info.Msg = make(map[string]interface{})
			}
			list = append(list, info)
		}
	default:
		return nil, result.NewError(result.ParamsError)
	}

	//查出转发者信息
	user, err := orm.GetUserInfoById(callerId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if user != nil {
		forwardUserName = user.Username
	}

	logStatistics.Debug("[Forward] start send forward", "log number", len(list), "rooms number", len(targetRooms), "users number", len(targetUsers))
	var roomFails = make([]string, 0)
	var userFails = make([]string, 0)
	switch forwardType {
	case types.SingleForward:
		//发送到群中
		for _, v := range targetRooms {
			//判断是否禁言
			isMuted, err := CheckMemberMuted(v, callerId)
			if err != nil {
				logStatistics.Warn("[forward] user muted db err", "roomId", v, "userId", callerId)
				roomFails = append(roomFails, v)
				continue
			}
			if isMuted {
				logStatistics.Debug("[forward] user have been muted", "roomId", v, "userId", callerId)
				roomFails = append(roomFails, v)
				continue
			}
			var msgs = make([]*proto.Proto, 0)
			for _, msg := range list {
				msgType := msg.MsgType
				//将系统消息转为普通消息
				if msgType == types.System {
					msgType = types.Text
				}
				msgTime := utility.NowMillionSecond()
				p, err := SendAlertCreate(callerId, v, types.ToRoom, nil, msgType, proto.ComposeSingleForwardMsg(chType, sourceId, fromName, msg), msgTime)
				if err != nil {
					roomFails = append(roomFails, v)
					continue
				}
				msgs = append(msgs, p)
			}
			SendBatch(callerId, v, types.ToRoom, nil, msgs)
		}
		//发送给用户
		for _, v := range targetUsers {
			////判断是否为好友
			//if ok, err := orm.CheckIsFriend(callerId, v, types.FriendIsNotDelete); !ok || err != nil {
			//	logStatistics.Debug("[receive] is not friend", "you", callerId, "friend", v)
			//	userFails = append(userFails, v)
			//	continue
			//}
			//判断是否是黑名单用户
			if ok, err := CheckIsBlocked(v, callerId); ok || err != nil {
				if err != nil {
					logStatistics.Error("[receive] check is blocked", "err", err.Error())
				}
				logStatistics.Debug("[receive] you are in the blocked", "you", callerId, "friend", v)
				userFails = append(userFails, v)
				continue
			}
			var msgs = make([]*proto.Proto, 0)
			for _, msg := range list {
				msgType := msg.MsgType
				//将系统消息转为普通消息
				if msgType == types.System {
					msgType = types.Text
				}
				msgTime := utility.NowMillionSecond()
				p, err := SendAlertCreate(callerId, v, types.ToUser, []string{v, callerId}, msgType, proto.ComposeSingleForwardMsg(chType, sourceId, fromName, msg), msgTime)
				if err != nil {
					userFails = append(userFails, v)
					continue
				}
				msgs = append(msgs, p)
			}
			SendBatch(callerId, v, types.ToUser, []string{v, callerId}, msgs)
		}
	case types.MergeForward:
		//发送到群中
		for _, v := range targetRooms {
			//判断是否禁言
			isMuted, err := CheckMemberMuted(v, callerId)
			if err != nil {
				logStatistics.Warn("[forward] user muted db err", "roomId", v, "userId", callerId)
				roomFails = append(roomFails, v)
				continue
			}
			if isMuted {
				logStatistics.Debug("[forward] user have been muted", "roomId", v, "userId", callerId)
				roomFails = append(roomFails, v)
				continue
			}
			SendAlert(callerId, v, types.ToRoom, nil, types.Forward, proto.ComposeForwardMsg(chType, sourceId, fromName, forwardUserName, list))
		}
		//发送给用户
		for _, v := range targetUsers {
			////判断是否为好友
			//if ok, err := orm.CheckIsFriend(callerId, v, types.FriendIsNotDelete); !ok || err != nil {
			//	logStatistics.Debug("[receive] is not friend", "you", callerId, "friend", v)
			//	userFails = append(userFails, v)
			//	continue
			//}
			//判断是否是黑名单用户
			if ok, err := CheckIsBlocked(v, callerId); ok || err != nil {
				if err != nil {
					logStatistics.Error("[receive] check is blocked", "err", err.Error())
				}
				logStatistics.Debug("[receive] you are in the blocked", "you", callerId, "friend", v)
				userFails = append(userFails, v)
				continue
			}
			SendAlert(callerId, v, types.ToUser, []string{v, callerId}, types.Forward, proto.ComposeForwardMsg(chType, sourceId, fromName, forwardUserName, list))
		}
	default:
		return nil, result.NewError(result.ParamsError)
	}
	logStatistics.Debug("[Forward] success send forward", "log number", len(list), "rooms number", len(targetRooms), "users number", len(targetUsers))

	ret := make(map[string]interface{})
	ret["failsNumber"] = len(roomFails) + len(userFails)
	ret["roomFails"] = roomFails
	ret["userFails"] = userFails
	return ret, nil
}

//转发消息(加密)
func ForwardEncryptMsg(callerId string, forwardType int, roomLogs, usersLogs []map[string]interface{}) (interface{}, error) {
	var list = make([]*types.ChatLog, 0)
	logStatistics.Debug("[Forward] start send forward", "log number", len(list), "rooms number", len(roomLogs), "users number", len(usersLogs))
	var roomFails = make([]string, 0)
	var userFails = make([]string, 0)
	switch forwardType {
	case types.SingleForward:
		//发送到群中
		for _, v := range roomLogs {
			//判断是否禁言
			isMuted, err := CheckMemberMuted(v["targetId"].(string), callerId)
			if err != nil {
				logStatistics.Warn("[forward] user muted db err", "roomId", v["targetId"].(string), "userId", callerId)
				roomFails = append(roomFails, v["targetId"].(string))
				continue
			}
			if isMuted {
				logStatistics.Debug("[forward] user have been muted", "roomId", v["targetId"].(string), "userId", callerId)
				roomFails = append(roomFails, v["targetId"].(string))
				continue
			}

			var msgs = make([]*proto.Proto, 0)

			messages := v["messages"].([]interface{})
			for _, kv := range messages {
				r := kv.(map[string]interface{})
				//messages 包含msgType和加密信息msg
				msgType := utility.ToInt(r["msgType"])
				//将系统消息转为普通消息
				if msgType == types.System {
					msgType = types.Text
				}
				msg := r["msg"].(map[string]interface{})

				msgTime := utility.NowMillionSecond()
				p, err := SendAlertCreate(callerId, v["targetId"].(string), types.ToRoom, nil, msgType, msg, msgTime)
				if err != nil {
					roomFails = append(roomFails, v["targetId"].(string))
					continue
				}
				msgs = append(msgs, p)
			}

			SendBatch(callerId, v["targetId"].(string), types.ToRoom, nil, msgs)
		}
		//发送给用户
		for _, v := range usersLogs {
			////判断是否为好友
			//if ok, err := orm.CheckIsFriend(callerId, v["targetId"].(string), types.FriendIsNotDelete); !ok || err != nil {
			//	logStatistics.Debug("[receive] is not friend", "you", callerId, "friend", v["targetId"].(string))
			//	userFails = append(userFails, v["targetId"].(string))
			//	continue
			//}
			//判断是否是黑名单用户
			if ok, err := CheckIsBlocked(v["targetId"].(string), callerId); ok || err != nil {
				if err != nil {
					logStatistics.Error("[receive] check is blocked", "err", err.Error())
				}
				logStatistics.Debug("[receive] you are in the blocked", "you", callerId, "friend", v["targetId"].(string))
				userFails = append(userFails, v["targetId"].(string))
				continue
			}
			var msgs = make([]*proto.Proto, 0)
			messages := v["messages"].([]interface{})
			for _, kv := range messages {
				r := kv.(map[string]interface{})
				//messages 包含msgType和加密信息msg
				msgType := utility.ToInt(r["msgType"])
				//将系统消息转为普通消息
				if msgType == types.System {
					msgType = types.Text
				}
				msg := r["msg"].(map[string]interface{})

				msgTime := utility.NowMillionSecond()
				p, err := SendAlertCreate(callerId, v["targetId"].(string), types.ToUser, []string{v["targetId"].(string), callerId}, msgType, msg, msgTime)
				if err != nil {
					userFails = append(userFails, v["targetId"].(string))
					continue
				}
				msgs = append(msgs, p)
			}
			SendBatch(callerId, v["targetId"].(string), types.ToUser, []string{v["targetId"].(string), callerId}, msgs)
		}
	case types.MergeForward:
		//发送到群中
		for _, v := range roomLogs {
			//判断是否禁言
			isMuted, err := CheckMemberMuted(v["targetId"].(string), callerId)
			if err != nil {
				logStatistics.Warn("[forward] user muted db err", "roomId", v["targetId"].(string), "userId", callerId)
				roomFails = append(roomFails, v["targetId"].(string))
				continue
			}
			if isMuted {
				logStatistics.Debug("[forward] user have been muted", "roomId", v["targetId"].(string), "userId", callerId)
				roomFails = append(roomFails, v["targetId"].(string))
				continue
			}

			messages := v["messages"].([]interface{})
			for _, kv := range messages {
				r := kv.(map[string]interface{})
				msg := r["msg"].(map[string]interface{})

				SendAlert(callerId, v["targetId"].(string), types.ToRoom, nil, types.Forward, msg)
			}
		}
		//发送给用户
		for _, v := range usersLogs {
			////判断是否为好友
			//if ok, err := orm.CheckIsFriend(callerId, v["targetId"].(string), types.FriendIsNotDelete); !ok || err != nil {
			//	logStatistics.Debug("[receive] is not friend", "you", callerId, "friend", v)
			//	userFails = append(userFails, v["targetId"].(string))
			//	continue
			//}
			//判断是否是黑名单用户
			if ok, err := CheckIsBlocked(v["targetId"].(string), callerId); ok || err != nil {
				if err != nil {
					logStatistics.Error("[receive] check is blocked", "err", err.Error())
				}
				logStatistics.Debug("[receive] you are in the blocked", "you", callerId, "friend", v["targetId"].(string))
				userFails = append(userFails, v["targetId"].(string))
				continue
			}
			messages := v["messages"].([]interface{})

			for _, kv := range messages {
				r := kv.(map[string]interface{})
				msg := r["msg"].(map[string]interface{})

				SendAlert(callerId, v["targetId"].(string), types.ToUser, []string{v["targetId"].(string), callerId}, types.Forward, msg)
			}
		}
	default:
		return nil, result.NewError(result.ParamsError)
	}
	logStatistics.Debug("[Forward] success send forward", "log number", len(list), "rooms number", len(roomLogs), "users number", len(usersLogs))

	ret := make(map[string]interface{})
	ret["failsNumber"] = len(roomFails) + len(userFails)
	ret["roomFails"] = roomFails
	ret["userFails"] = userFails
	return ret, nil
}

// 撤回指定的一条消息
func RevokeMsg(userId, logId string, cType int) error {
	if cType == types.IsRoom {
		log, err := orm.FindRoomChatLogByContentId(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if log == nil {
			logStatistics.Warn("RevokeMsg", "warn", "CanNotFindRoomMsg", "logId", logId)
			return result.NewError(result.CanNotFindRoomMsg)
		}
		senderId := log.SenderId
		sendTime := log.Datetime
		msgType := log.MsgType
		roomId := log.RoomId

		// 红包和提示消息不可撤回
		if msgType == types.RedPack || msgType == types.Alert {
			return result.NewError(result.DeleteMsgFailed).SetExtMessage("该消息不可撤回").SetChildErr(result.ServiceChat, nil, "message type err")
		}

		operatorLevel := orm.GetMemberLevel(roomId, userId, types.RoomUserDeletedOrNot)
		if operatorLevel == types.RoomLevelNotExist {
			logStatistics.Warn("RevokeMsg", "warn", "UserIsNotInRoom", "userId", userId, "roomId", roomId)
			return result.NewError(result.UserIsNotInRoom)
		}
		//判断操作者是否为普通用户
		if operatorLevel != types.RoomLevelMaster && operatorLevel != types.RoomLevelManager {
			//普通用户：只能撤回自己10分钟内发送的消息
			if senderId != userId {
				logStatistics.Warn("RevokeMsg", "warn", "DeleteMsgFailed: normal member can not revoke other members message", "userId", userId, "senderId", senderId)
				return result.NewError(result.DeleteMsgFailed)
			}
			lastTime := utility.MillionSecondAddDuration(utility.NowMillionSecond(), -(10 * time.Minute))
			if sendTime < lastTime {
				logStatistics.Warn("RevokeMsg", "warn", "CanNotDelMsgOverTime", "logId", logId, "sendTime", sendTime, "lastTime", lastTime)
				return result.NewError(result.CanNotDelMsgOverTime)
			}
		} else if senderId != userId {
			//群主和管理员
			senderLevel := orm.GetMemberLevel(roomId, senderId, types.RoomUserDeletedOrNot)
			if senderLevel == types.RoomLevelNotExist {
				logStatistics.Warn("RevokeMsg", "warn", "UserIsNotInRoom", "roomId", roomId, "userId", senderId)
				return result.NewError(result.UserIsNotInRoom)
			}
			//非公告消息：管理员不可撤回其他管理员和群主的消息
			if msgType != types.System && operatorLevel <= senderLevel {
				logStatistics.Warn("RevokeMsg", "warn", "PermissionDeny", "logId", logId, "userId", userId, "operator level", operatorLevel, "msg owner level", senderLevel)
				return result.NewError(result.PermissionDeny)
			}
		}

		count, err := orm.DelRoomChatLogById(logId)
		if count != 1 {
			return result.NewError(result.DeleteMsgFailed)
		}

		// get user info in the room
		userRoomNickname := orm.GetMemberName(roomId, userId)
		//send alert
		var msg string
		if senderId != userId {
			if operatorLevel == types.RoomLevelManager {
				msg = "管理员撤回了一条成员消息"
			} else if operatorLevel == types.RoomLevelMaster {
				msg = "群主撤回了一条成员消息"
			}
		} else {
			msg = userRoomNickname + "撤回了一条消息"
		}
		SendAlert(userId, roomId, types.ToRoom, nil, types.Alert, proto.ComposeRevokeAlert(logId, userId, userRoomNickname, msg, operatorLevel))
		return nil
	}
	if cType == types.IsFriend {
		log, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if log == nil {
			logStatistics.Warn("RevokeMsg", "warn", "CanNotFindFriendMsg", "logId", logId)
			return result.NewError(result.CanNotFindFriendMsg)
		}

		senderId := log.SenderId
		sendTime := log.SendTime
		msgType := log.MsgType
		receiveId := log.ReceiveId

		if msgType == types.RedPack || msgType == types.Alert {
			return result.NewError(result.DeleteMsgFailed).SetExtMessage("该消息不可撤回").SetChildErr(result.ServiceChat, nil, "message type err")
		}

		if senderId != userId {
			logStatistics.Warn("RevokeMsg", "warn", "DeleteMsgFailed: can not revoke friend message", "logId", logId, "userId", userId, "senderId", senderId)
			return result.NewError(result.DeleteMsgFailed)
		}

		lastTime := utility.MillionSecondAddDuration(utility.NowMillionSecond(), -(10 * time.Minute))
		if sendTime < lastTime {
			logStatistics.Warn("RevokeMsg", "warn", "CanNotDelMsgOverTime", "logId", logId, "sendTime", sendTime, "lastTime", lastTime)
			return result.NewError(result.CanNotDelMsgOverTime)
		}

		count, err := orm.DelPrivateChatLog(logId)
		if count < 1 {
			return result.NewError(result.DeleteMsgFailed)
		}
		// send alert
		//查找该显示的名称
		nameSM, err := UserInfo(userId)
		if err != nil {
			return err
		}
		receiveName := utility.ToString(nameSM["name"])

		msg := receiveName + "撤回了一条消息"
		SendAlert(userId, receiveId, types.ToUser, []string{userId, receiveId}, types.Alert, proto.ComposeRevokeAlert(logId, userId, receiveName, msg, types.RoomLevelNomal))
		return nil
	}
	return result.NewError(result.ParamsError)
}

func revokeMultimedia(userId, logId string, cType int) error {
	if cType == types.IsRoom {
		log, err := orm.FindRoomChatLogByContentId(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if log == nil {
			logStatistics.Warn("revokeMultimedia", "warn", "CanNotFindRoomMsg", "logId", logId)
			return result.NewError(result.CanNotFindRoomMsg)
		}
		senderId := log.SenderId
		//sendTime := utility.ToInt64(maps[0]["datetime"])
		msgType := log.MsgType
		roomId := log.RoomId

		// 只可批量撤回 文件、图片、视频消息
		if msgType != types.File {
			return result.NewError(result.DeleteMsgFailed).SetChildErr(result.ServiceChat, nil, "msg type is not File")
		}

		operatorLevel := orm.GetMemberLevel(roomId, userId, types.RoomUserDeletedOrNot)
		if operatorLevel == types.RoomLevelNotExist {
			logStatistics.Warn("revokeMultimedia", "warn", "UserIsNotInRoom", "roomId", roomId, "userId", userId)
			return result.NewError(result.UserIsNotInRoom)
		}
		//判断操作者是否为普通用户
		if operatorLevel != types.RoomLevelMaster && operatorLevel != types.RoomLevelManager {
			//普通用户：只能撤回自己的消息
			if senderId != userId {
				logStatistics.Warn("revokeMultimedia", "warn", "DeleteMsgFailed: normal member can not revoke other members message", "logId", logId, "userId", userId, "senderId", senderId)
				return result.NewError(result.DeleteMsgFailed)
			}
			//lastTime := utility.MillionSecondAddDuration(utility.NowMillionSecond(), -(10 * time.Minute))
			//if sendTime < lastTime {
			//	return &result.Error{ErrorCode: result.CanNotDelMsgOverTime, Message: ""}
			//}
		} else if senderId != userId {
			//群主和管理员
			senderLevel := orm.GetMemberLevel(roomId, senderId, types.RoomUserDeletedOrNot)
			if senderLevel == types.RoomLevelNotExist {
				logStatistics.Warn("revokeMultimedia", "warn", "UserIsNotInRoom", "roomId", roomId, "userId", senderId)
				return result.NewError(result.UserIsNotInRoom)
			}
			//非公告消息：管理员不可撤回其他管理员和群主的消息
			if msgType != types.System && operatorLevel <= senderLevel {
				logStatistics.Warn("revokeMultimedia", "warn", "PermissionDeny", "logId", logId, "userId", userId, "operator level", operatorLevel, "msg owner level", senderLevel)
				return result.NewError(result.PermissionDeny)
			}
		}

		//删除消息
		count, err := orm.DelRoomChatLogById(logId)
		if count != 1 {
			return result.NewError(result.DeleteMsgFailed)
		}

		// get user info in the room
		userRoomNickname := orm.GetMemberName(roomId, userId)
		//send alert
		var msg string
		if senderId != userId {
			if operatorLevel == types.RoomLevelManager {
				msg = "管理员撤回了一条成员消息"
			} else if operatorLevel == types.RoomLevelMaster {
				msg = "群主撤回了一条成员消息"
			}
		} else {
			msg = userRoomNickname + "撤回了一条消息"
		}
		SendAlert(userId, roomId, types.ToRoom, nil, types.Alert, proto.ComposeRevokeAlert(logId, userId, userRoomNickname, msg, operatorLevel))
		return nil
	}
	if cType == types.IsFriend {
		log, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if log == nil {
			logStatistics.Warn("revokeMultimedia", "warn", "CanNotFindFriendMsg", "logId", logId)
			return result.NewError(result.CanNotFindFriendMsg)
		}

		senderId := log.SenderId
		//sendTime := utility.ToInt64(maps[0]["send_time"])
		msgType := log.MsgType
		receiveId := log.ReceiveId

		// 只可批量撤回 文件、图片、视频消息
		if msgType != types.File {
			return result.NewError(result.DeleteMsgFailed).SetChildErr(result.ServiceChat, nil, "msg type is not File")
		}

		if senderId != userId {
			logStatistics.Warn("revokeMultimedia", "warn", "DeleteMsgFailed: can not revoke friend message", "logId", logId, "userId", userId, "senderId", senderId)
			return result.NewError(result.DeleteMsgFailed)
		}

		//lastTime := utility.MillionSecondAddDuration(utility.NowMillionSecond(), -(10 * time.Minute))
		//if sendTime < lastTime {
		//	return &result.Error{ErrorCode: result.CanNotDelMsgOverTime, Message: ""}
		//}

		count, err := orm.DelPrivateChatLog(logId)
		if count < 1 {
			return result.NewError(result.DeleteMsgFailed)
		}
		// send alert
		//查找该显示的名称
		nameSM, err := UserInfo(userId)
		if err != nil {
			return err
		}
		receiveName := utility.ToString(nameSM["name"])

		msg := receiveName + "撤回了一条消息"
		SendAlert(userId, receiveId, types.ToUser, []string{userId, receiveId}, types.Alert, proto.ComposeRevokeAlert(logId, userId, receiveName, msg, types.RoomLevelNomal))
		return nil
	}
	return result.NewError(result.ParamsError)
}

// 批量撤回文件消息
func BatchRevokeFiles(userId string, logs []string, cType int) (interface{}, error) {
	failed := make([]string, 0)
	for _, v := range logs {
		err := revokeMultimedia(userId, v, cType)
		if err != nil {
			failed = append(failed, v)
		}
	}
	ret := make(map[string]interface{})
	ret["failsNumber"] = len(failed)
	ret["fails"] = failed
	return ret, nil
}
