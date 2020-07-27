package model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	var cfg types.Config
	if _, err := toml.DecodeFile("../etc/config.toml", &cfg); err != nil {
		panic(err)
	}
	db.InitDB(&cfg)
}

/*func Test_ChatLogBatch(t *testing.T) {
	userId := "1"
	var time int64
	logs := types.NewchatLogs()
	if time == 0 {
		//若时间是0 则推送成已读格式
		time = utility.MillionSecondAddDate(utility.NowMillionSecond(), 0, 0, -types.BatchDayAgo)
	}
	rooms, roomTime, err := getNeedBatchRooms(userId)
	if err != nil {
		batchLog.Error("getNeedBatchRooms err", "err", err.Error())
	}
	maps, err := db.GetRoomMsgContentAfter(rooms, time, types.RoomMsgNotDelete)
	if err != nil {
		batchLog.Error("GetRoomMsgContentAfter err", "err", err.Error())
	}
	for _, v := range maps {
		log, err := GetChatLogAsRoom(userId, v)
		if err != nil {
			batchLog.Error("GetChatLogAsRoom err", "err", err.ErrorCode)
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
			maps2, err2 := db.GetRoomRevMsgByLogId(log.LogId, userId)
			if err2 != nil {
				batchLog.Error("GetRoomRevMsgByLogId err", "err", err2)
				continue
			}
			if len(maps2) > 0 {
				revLog := maps2[0]
				if utility.ToInt(revLog["state"]) == types.HadBurnt {
					continue
				}
			} else {
				continue
			}
		}
		if log.MsgType == types.Alert {
			ck := proto.NewMsg(log.GetCon())
			if ck.CheckSpecialMsg() {
				maps2, err2 := db.GetRoomRevMsgByLogId(log.LogId, userId)
				if err2 != nil {
					batchLog.Error("GetRoomRevMsgByLogId err", "err", err2)
					continue
				}
				if len(maps2) < 1 {
					continue
				}
			}
		}
		logs = append(logs, log)
	}
	roomLength := len(logs)
	userMaps, err := db.FindNotBurndLogAfter(userId, types.FriendMsgIsNotDelete, time)
	if err != nil {
		batchLog.Error("FindNotBurndLogAfter err", "err", err.Error())
		return
	}

	for _, v := range userMaps {
		log, err := GetChatLogAsUser(userId, v)
		if err != nil {
			batchLog.Error("GetChatLogAsUser err", "err", err.ErrorCode)
			continue
		}
		logs = append(logs, log)
	}
	userLength := len(logs) - roomLength
	batchLog.Debug("Push Room Recent Number:", "Number", roomLength)
	batchLog.Debug("Push User Recent Number:", "Number", userLength)

	fmt.Println("before")
	for _, v := range logs {
		fmt.Println(v.Datetime)
	}

	sort.Sort(logs)
	fmt.Println("after")
	for _, v := range logs {
		fmt.Println(v.Datetime)
	}
	fmt.Println("ok")
}*/
