package model

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	l "github.com/inconshreveable/log15"
)

var batchLog = l.New("module", "model/batch")

var batchConfig *types.BatchConfig

type roomBatchInfo struct {
	IsDel      int
	CreateTime int64
}

type BatchSingle struct {
	P        *proto.Proto
	FromId   string
	SendTime int64
}

func FormatBatchTargetMsg(batchList []*BatchSingle, eventType int) []byte {
	var bData = make(map[string]interface{})
	bData["eventType"] = eventType
	var list = make([]interface{}, 0)
	for _, v := range batchList {
		p := v.P
		userId := v.FromId
		sendTime := v.SendTime
		//获取用户信息
		var data = make(map[string]interface{})
		data["msgId"] = p.GetMsgId()
		data["logId"] = p.GetLogId()
		data["fromId"] = userId
		data["channelType"] = p.GetChannelType()
		data["isSnap"] = p.GetIsSnap()
		data["targetId"] = p.GetTargetId()
		data["msgType"] = p.GetMsgType()
		data["msg"] = p.GetMsg()
		data["datetime"] = sendTime
		if p.GetPraise() != nil {
			data["praise"] = p.GetPraise()
		}
		/*if p.GetMsgType() == types.RedPack {
			data["msg"] = redpacketLengthCheck(p.GetMsg())
		}*/

		userInfo, err := orm.GetUserInfoById(userId)
		if err != nil {
			batchLog.Error("FormatBatchTargetMsg", "err", err.Error())
			continue
		}
		var senderInfo = make(map[string]interface{})
		if userInfo == nil {
			//未找到用户
			senderInfo["nickname"] = utility.VisitorNameSplit(userId)
			senderInfo["avatar"] = ""
		} else {
			senderInfo["nickname"] = userInfo.Username
			senderInfo["avatar"] = userInfo.Avatar
		}
		data["senderInfo"] = senderInfo

		list = append(list, data)
	}
	bData["list"] = list

	targetData, _ := json.Marshal(bData)
	return targetData
}

func SendBatch(fromId, targetId string, target int, members []string, list []*proto.Proto) {
	var batchList = make([]*BatchSingle, 0)
	for i, p := range list {
		batchList = append(batchList, &BatchSingle{P: p, FromId: fromId, SendTime: p.GetSendTime()})
		if i+1 == len(list) || (i+1)%batchConfig.BatchPackLength == 0 {
			targetData := FormatBatchTargetMsg(batchList, types.EventBatchCustomize)
			if targetData != nil {
				time.Sleep(time.Duration(batchConfig.BatchInterval) * time.Millisecond)

				if target == types.ToRoom {
					if members == nil {
						clId := types.GetRoomRouteById(targetId)
						cl, _ := router.GetChannel(clId)
						if cl != nil {
							cl.Broadcast(targetData)
						}
					} else {
						for _, memId := range members {
							if client, ok := router.GetUser(memId); ok && client != nil {
								client.SendToAllClients(targetData)
							}
						}
					}
				}
				if target == types.ToUser {
					//text := "收到一条好友消息"
					for _, memId := range members {
						if u, ok := router.GetUser(memId); ok && u != nil {
							u.SendToAllClients(targetData)
						}
					}
				}
			}
			batchList = make([]*BatchSingle, 0)
		}
	}
}

//log已经经过解析成info
func batchPusherAfterParse(client *router.Client, logs []*types.ChatLog, eventType int) error {
	var batchList = make([]*BatchSingle, 0)
	for i, log := range logs {
		p, err := proto.NewCmnProto(log.LogId, log.MsgId, log.TargetId, log.ChannelType, log.MsgType, log.IsSnap, log.GetCon(), log.Datetime, log.Praise)
		if err != nil {
			batchLog.Error("create cmn proto", "err", err.Error())
			return errors.New("error create cmn proto1")
		}

		batchList = append(batchList, &BatchSingle{P: p, FromId: log.FromId, SendTime: log.Datetime})
		if i+1 == len(logs) || (i+1)%batchConfig.BatchPackLength == 0 {
			targetData := FormatBatchTargetMsg(batchList, eventType)
			if targetData != nil {
				time.Sleep(time.Duration(batchConfig.BatchInterval) * time.Millisecond)
				client.Send(targetData)
			}
			batchList = make([]*BatchSingle, 0)
		}
	}
	return nil
}

func getNeedBatchRooms(userId string) ([]string, map[string]*roomBatchInfo, error) {
	joinedRooms, err := orm.FindJoinedRooms(userId)
	if err != nil {
		return nil, nil, err
	}
	var roomBatch = make(map[string]*roomBatchInfo)
	var rooms = make([]string, len(joinedRooms))
	for i, v := range joinedRooms {
		id := utility.ToString(v.RoomId)
		rooms[i] = id
		roomBatch[id] = &roomBatchInfo{
			CreateTime: utility.ToInt64(v.CreateTime),
			IsDel:      utility.ToInt(v.IsDelete),
		}
	}
	return rooms, roomBatch, nil
}

// 批量推送断开连接后的最近消息
func BatchPushRecentMsg(client *router.Client, time int64) {
	user := client.GetUser()
	if user == nil {
		return
	}
	userId := user.Id
	pushType := types.EventBatchPushUnread //默认推送成未读格式
	//var logs = make([]*types.ChatLog, 0)
	logs := types.NewchatLogs()
	batchLog.Debug("start recent", "time", utility.MillionSecondToTimeString(time))
	//如果是0表示是新设备
	if time == 0 {
		//若时间是0 则推送成已读格式
		pushType = types.EventBatchPush
		time = utility.MillionSecondAddDate(utility.NowMillionSecond(), 0, 0, -batchConfig.BatchDayAgo)
	}
	rooms, roomTime, err := getNeedBatchRooms(userId)
	if err != nil {
		batchLog.Error("getNeedBatchRooms err", "err", err.Error())
		return
	}
	dbLogs, err := orm.FindRoomChatLogsAfterTime(rooms, time)
	if err != nil {
		batchLog.Error("GetRoomMsgContentAfter err", "err", err.Error())
		return
	}
	for _, v := range dbLogs {
		log, err := GetChatLogAsRoom(userId, v)
		if err != nil {
			batchLog.Error("GetChatLogAsRoom err", "err", err.Error())
			continue
		}
		if rt, ok := roomTime[log.TargetId]; !ok {
			continue
		} else {
			//过滤掉已经退出群聊后的消息
			if rt.IsDel == types.RoomUserDeleted && log.Datetime > rt.CreateTime {
				continue
			}
		}
		//过滤掉已焚毁和不属于自己的消息
		if log.IsSnap == types.IsSnap {
			revLog, err := orm.FindReceiveLogById(log.LogId, userId)
			if err != nil {
				batchLog.Error("GetRoomRevMsgByLogId err", "err", err)
				continue
			}
			if revLog != nil {
				if revLog.State == types.HadBurnt {
					continue
				}
			} else {
				continue
			}
		}
		if log.MsgType == types.Alert {
			ck := proto.NewMsg(log.GetCon())

			//过滤掉赞赏通知
			if pushType == types.EventBatchPush && ck.GetTypeCode() == types.AlertPraise {
				continue
			}

			if ck.CheckSpecialMsg() {
				revLog, err := orm.FindReceiveLogById(log.LogId, userId)
				if err != nil {
					batchLog.Error("GetRoomRevMsgByLogId err", "err", err)
					continue
				}
				if revLog == nil {
					continue
				}
			}
		}
		logs = append(logs, log)
	}
	roomLength := len(logs)
	nbLogs, err := orm.FindNotBurnedLogsAfter(userId, types.FriendMsgIsNotDelete, time)
	if err != nil {
		batchLog.Error("FindNotBurndLogAfter err", "err", err.Error())
		return
	}

	for _, v := range nbLogs {
		log, err := GetChatLogAsUser(userId, v)
		if err != nil {
			batchLog.Error("GetChatLogAsUser err", "err", err.Error())
			continue
		}
		if log.MsgType == types.Alert {
			msg := proto.NewMsg(log.GetCon())

			//过滤掉赞赏通知
			if pushType == types.EventBatchPush && msg.GetTypeCode() == types.AlertPraise {
				continue
			}

			if msg.GetTypeCode() == types.AlertDeleteFriend {
				if opt, ok := msg["operator"].(string); ok && opt == userId {
					continue
				}
			}

			//过滤掉不属于自己的红包通知消息
			if msg.GetTypeCode() == types.AlertReceiveRedpackage {
				if opt, ok := msg["operator"].(string); ok && opt != userId {
					if owner, ok := msg["owner"].(string); ok && owner != userId {
						continue
					}
				}
			}
		}
		logs = append(logs, log)
	}
	userLength := len(logs) - roomLength
	batchLog.Debug("Push Room Recent Number:", "Number", roomLength)
	batchLog.Debug("Push User Recent Number:", "Number", userLength)

	sort.Sort(logs)
	err = batchPusherAfterParse(client, logs, pushType)
	if err == nil {
		//发送推送成功通知
		proto.SendBatchSuccessNotification(client)
		batchLog.Debug("stop recent: success")
	} else {
		batchLog.Debug("stop recent: fail", "err", err.Error())
	}
}

//确认消息批量发送
func BatchPushAckMsg(client *router.Client, begin, end, total int64) {
	loger := l.New("func", "batch Ack", "logId", utility.RandInt(1, 1000))

	userId := client.GetUser().Id
	pushType := types.EventBatchAck
	logs := types.NewackChatLogs()

	//如果无穷大
	if end == -1 {
		end = math.MaxInt64
	}
	loger.Debug("start ack", "userId", userId, "begin time", utility.MillionSecondToTimeString2(begin), "begin timestamp", begin, "end time", utility.MillionSecondToTimeString2(end), "end timestamp", end)
	now := utility.NowMillionSecond()

	rooms, roomTime, err := getNeedBatchRooms(userId)
	if err != nil {
		loger.Error("getNeedBatchRooms err", "err", err.Error())
		return
	}

	rLogs, err := orm.FindRoomChatLogsBetweenTime(rooms, begin, end)
	if err != nil {
		loger.Error("GetRoomMsgContentAfter err", "err", err.Error())
		return
	}
	for _, v := range rLogs {
		log, err := GetChatLogAsRoom(userId, v)
		if err != nil {
			loger.Error("GetChatLogAsRoom err", "err", err.Error())
			continue
		}
		if rt, ok := roomTime[log.TargetId]; !ok {
			continue
		} else {
			//过滤掉已经退出群聊后的消息
			if rt.IsDel == types.RoomUserDeleted && log.Datetime > rt.CreateTime {
				continue
			}
		}

		//过滤掉已焚毁和不属于自己的消息
		if log.IsSnap == types.IsSnap {
			revLog, err := orm.FindReceiveLogById(log.LogId, userId)
			if err != nil {
				loger.Error("FindReceiveLogById err", "err", err)
				continue
			}
			if revLog != nil {
				if revLog.State == types.HadBurnt {
					continue
				}
			} else {
				continue
			}
		}

		if log.MsgType == types.Alert {
			ck := proto.NewMsg(log.GetCon())

			//过滤掉赞赏通知
			if ck.GetTypeCode() == types.AlertPraise {
				continue
			}

			if ck.CheckSpecialMsg() {
				if ck.GetTypeCode() == types.AlertHadBurntMsg {
					continue
				}
				revLog, err := orm.FindReceiveLogById(log.LogId, userId)
				if err != nil {
					loger.Error("FindReceiveLogById err", "err", err)
					continue
				}
				if revLog == nil {
					continue
				}
			}
		}
		logs = append(logs, log)
	}
	roomLength := len(logs)
	nbLogs, err := orm.FindNotBurnedLogsBetween(userId, types.FriendMsgIsNotDelete, begin, end)
	if err != nil {
		loger.Error("FindNotBurnedLogsBetween err", "err", err.Error())
		return
	}

	for _, v := range nbLogs {
		log, err := GetChatLogAsUser(userId, v)
		if err != nil {
			loger.Error("GetChatLogAsUser err", "err", err.Error())
			continue
		}
		if log.MsgType == types.Alert {
			msg := proto.NewMsg(log.GetCon())

			if msg.GetTypeCode() == types.AlertDeleteFriend {
				if opt, ok := msg["operator"].(string); ok && opt == userId {
					continue
				}
			}
			//过滤掉赞赏通知
			if msg.GetTypeCode() == types.AlertPraise {
				continue
			}
			//过滤掉不属于自己的红包通知消息
			if msg.GetTypeCode() == types.AlertReceiveRedpackage {
				if opt, ok := msg["operator"].(string); ok && opt != userId {
					if owner, ok := msg["owner"].(string); ok && owner != userId {
						continue
					}
				}
			}
		}
		logs = append(logs, log)
	}

	sort.Sort(logs)
	if end == math.MaxInt64 {
		end = now
	}

	if utility.ToInt64(len(logs)) != total {
		userLength := len(logs) - roomLength
		loger.Debug("Ack Room Recent Number:", "Number", roomLength, "userId", userId, "begin time", utility.MillionSecondToTimeString2(begin), "begin timestamp", begin, "end time", utility.MillionSecondToTimeString2(end), "end timestamp", end)
		loger.Debug("Ack User Recent Number:", "Number", userLength, "userId", userId, "begin time", utility.MillionSecondToTimeString2(begin), "begin timestamp", begin, "end time", utility.MillionSecondToTimeString2(end), "end timestamp", end)

		err = batchPusherAfterParse(client, logs, pushType)
	}
	if err == nil {
		//发送推送成功通知
		proto.SendAckSuccessNotification(client, begin, end)
		loger.Debug("ack success")
	} else {
		loger.Debug("ack fail", "err", err.Error())
	}
}

//批量发送群会话秘钥
func BatchSKeyMsg(client *router.Client, datetime int64) {
	userId := client.GetUser().Id
	pushType := types.EventGetAllSKey
	logs := types.NewackChatLogs()
	batchLog.Debug("start get SKey")

	nbLogs, err := orm.FindSessionKeyAlert(userId, datetime)
	if err != nil {
		batchLog.Error("FindSessionKeyAlert err", "err", err.Error())
		return
	}

	for _, v := range nbLogs {
		log, err := GetChatLogAsUser(userId, v)
		if err != nil {
			batchLog.Error("GetChatLogAsUser err", "err", err.Error())
			continue
		}
		logs = append(logs, log)
	}

	batchLog.Debug("Room Session key Recent Number:", "Number", len(logs))

	sort.Sort(logs)
	err = batchPusherAfterParse(client, logs, pushType)
	if err == nil {
		//发送推送成功通知
		proto.SendGetSKeySuccessNotification(client)
		batchLog.Debug("stop recent: success")
	} else {
		batchLog.Debug("stop recent: fail", "err", err.Error())
	}
}
