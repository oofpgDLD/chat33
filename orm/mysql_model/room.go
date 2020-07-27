package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

// 获取用户加入的所有群id
func GetUserJoinedRooms(userId string) ([]string, error) {
	rooms, err := db.GetRoomsById(userId)
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	for _, v := range rooms {
		roomId := v["room_id"]
		data = append(data, roomId)
	}
	return data, nil
}

//获取所有群
func GetEnableRoomsId() ([]string, error) {
	var rlt []string
	list, err := db.GetEnabledRooms()
	if err != nil {
		return rlt, err
	}
	for _, v := range list {
		rlt = append(rlt, utility.ToString(v["id"]))
	}
	return rlt, nil
}

func GetMemberNumber(roomId string) (int64, error) {
	return db.GetRoomMemberNumber(roomId)
}

func GetRoomMemberNumberByLevel(roomId string, level int) (int64, error) {
	return db.GetRoomMemberNumberByLevel(roomId, level)
}

//返回群中名称
func GetRoomMemberName(roomId, userId string) (string, error) {
	maps, err := db.FindRoomMemberName(roomId, userId)
	if len(maps) < 1 {
		return "", err
	}
	return maps[0]["name"], nil
}

// 获取群成员身份等级(普通用户 管理员 群主)
func GetRoomUserLevel(roomId, userId string, isDel int) (int, error) {
	maps, err := db.GetRoomMemberInfo(roomId, userId, isDel)
	if err != nil || len(maps) < 1 {
		return types.RoomLevelNotExist, err
	}
	userInfo := maps[0]
	return utility.ToInt(userInfo["level"]), nil
}

func FindRoomMemberById(roomId, userId string, isDel int) (*types.MemberJoinUser, error) {
	maps, err := db.GetRoomMemberInfo(roomId, userId, isDel)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return &types.MemberJoinUser{
		RoomMember: &types.RoomMember{
			Id:           info["id"],
			RoomId:       info["room_id"],
			UserId:       info["user_id"],
			UserNickname: info["user_nickname"],
			Level:        utility.ToInt(info["level"]),
			NoDisturbing: utility.ToInt(info["no_disturbing"]),
			CommonUse:    utility.ToInt(info["common_use"]),
			RoomTop:      utility.ToInt(info["room_top"]),
			CreateTime:   utility.ToInt64(info["create_time"]),
			IsDelete:     utility.ToInt(info["is_delete"]),
			Source:       info["source"],
		},
		User: convertUser(info),
	}, nil
}

func FindJoinedRooms(userId string) ([]*types.RoomMember, error) {
	maps, err := db.GetJoinedRooms(userId)
	if err != nil {
		return nil, err
	}
	var members = make([]*types.RoomMember, 0)
	for _, info := range maps {
		m := &types.RoomMember{
			Id:           info["id"],
			RoomId:       info["room_id"],
			UserId:       info["user_id"],
			UserNickname: info["user_nickname"],
			Level:        utility.ToInt(info["level"]),
			NoDisturbing: utility.ToInt(info["no_disturbing"]),
			CommonUse:    utility.ToInt(info["common_use"]),
			RoomTop:      utility.ToInt(info["room_top"]),
			CreateTime:   utility.ToInt64(info["create_time"]),
			IsDelete:     utility.ToInt(info["is_delete"]),
			Source:       info["source"],
		}
		members = append(members, m)
	}
	return members, nil
}

func FindRoomMemberByName(roomId, name string) ([]*types.MemberJoinUser, error) {
	maps, err := db.GetRoomMemberInfoByName(roomId, name)
	if err != nil {
		return nil, err
	}
	var members = make([]*types.MemberJoinUser, 0)
	for _, info := range maps {
		m := &types.MemberJoinUser{
			RoomMember: &types.RoomMember{
				Id:           info["id"],
				RoomId:       info["room_id"],
				UserId:       info["user_id"],
				UserNickname: info["user_nickname"],
				Level:        utility.ToInt(info["level"]),
				NoDisturbing: utility.ToInt(info["no_disturbing"]),
				CommonUse:    utility.ToInt(info["common_use"]),
				RoomTop:      utility.ToInt(info["room_top"]),
				CreateTime:   utility.ToInt64(info["create_time"]),
				IsDelete:     utility.ToInt(info["is_delete"]),
				Source:       info["source"],
			},
			User: &types.User{
				Username: info["username"],
				Avatar:   info["avatar"],
			},
		}
		members = append(members, m)
	}
	return members, nil
}

func CheckRoomMarkIdExist(randomId string) (bool, error) {
	return db.CheckRoomMarkIdExist(randomId)
}

func convertRoom(v map[string]string) *types.Room {
	return &types.Room{
		Id:                 v["id"],
		MarkId:             v["mark_id"],
		Name:               v["name"],
		Avatar:             v["avatar"],
		MasterId:           v["master_id"],
		CreateTime:         utility.ToInt64(v["create_time"]),
		CanAddFriend:       utility.ToInt(v["can_add_friend"]),
		JoinPermission:     utility.ToInt(v["join_permission"]),
		RecordPermision:    utility.ToInt(v["record_permission"]),
		AdminMuted:         utility.ToInt(v["admin_muted"]),
		MasterMuted:        utility.ToInt(v["master_muted"]),
		Encrypt:            utility.ToInt(v["encrypt"]),
		IsDelete:           utility.ToInt(v["is_delete"]),
		RoomLevel:          utility.ToInt(v["room_level"]),
		CloseUntil:         utility.ToInt64(v["close_until"]),
		Recommend:          utility.ToInt(v["recommend"]),
		Identification:     utility.ToInt(v["identification"]),
		IdentificationInfo: v["identification_info"],
	}
}

func SearchRoomInfo(id string, isDel int) (*types.Room, error) {
	maps, err := db.GetRoomsInfo(id, isDel)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertRoom(info), nil
}

func SearchRoomInfoByMarkId(markId string, isDel int) (*types.Room, error) {
	maps, err := db.GetRoomsInfoByMarkId(markId, isDel)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertRoom(info), nil
}

// 获取群中所有管理员 返回：key 用户ID value level
func GetRoomManagerAndMaster(roomId string) ([]*types.MemberJoinUser, error) {
	var users = make([]*types.MemberJoinUser, 0)
	maps, err := db.GetRoomManagerAndMaster(roomId)
	if err != nil {
		return nil, err
	}
	for _, info := range maps {
		u := &types.MemberJoinUser{
			RoomMember: &types.RoomMember{
				Id:           info["id"],
				RoomId:       info["room_id"],
				UserId:       info["user_id"],
				UserNickname: info["user_nickname"],
				Level:        utility.ToInt(info["level"]),
				NoDisturbing: utility.ToInt(info["no_disturbing"]),
				CommonUse:    utility.ToInt(info["common_use"]),
				RoomTop:      utility.ToInt(info["room_top"]),
				CreateTime:   utility.ToInt64(info["create_time"]),
				IsDelete:     utility.ToInt(info["is_delete"]),
				Source:       info["source"],
			},
			User: &types.User{
				Username:    info["username"],
				Avatar:      info["avatar"],
				DeviceToken: info["device_token"],
			},
		}
		users = append(users, u)
	}
	return users, nil
}

func FindRoomMembers(roomId string, searchNumber int) ([]*types.MemberJoinUser, error) {
	var users = make([]*types.MemberJoinUser, 0)
	maps, err := db.GetRoomMembers(roomId, searchNumber)
	if err != nil {
		return nil, err
	}
	for _, info := range maps {
		u := &types.MemberJoinUser{
			RoomMember: &types.RoomMember{
				Id:           info["id"],
				RoomId:       info["room_id"],
				UserId:       info["user_id"],
				UserNickname: info["user_nickname"],
				Level:        utility.ToInt(info["level"]),
				NoDisturbing: utility.ToInt(info["no_disturbing"]),
				CommonUse:    utility.ToInt(info["common_use"]),
				RoomTop:      utility.ToInt(info["room_top"]),
				CreateTime:   utility.ToInt64(info["create_time"]),
				IsDelete:     utility.ToInt(info["is_delete"]),
				Source:       info["source"],
			},
			User: convertUser(info),
		}
		users = append(users, u)
	}
	return users, nil
}

func FindSetNoDisturbingMembers(roomId string) ([]*types.RoomMember, error) {
	var users = make([]*types.RoomMember, 0)
	maps, err := db.FindRoomMemberSetNoDisturbing(roomId, types.NoDisturbingOn)
	if err != nil {
		return nil, err
	}
	for _, info := range maps {
		u := &types.RoomMember{
			Id:           info["id"],
			RoomId:       info["room_id"],
			UserId:       info["user_id"],
			UserNickname: info["user_nickname"],
			Level:        utility.ToInt(info["level"]),
			NoDisturbing: utility.ToInt(info["no_disturbing"]),
			CommonUse:    utility.ToInt(info["common_use"]),
			RoomTop:      utility.ToInt(info["room_top"]),
			CreateTime:   utility.ToInt64(info["create_time"]),
			IsDelete:     utility.ToInt(info["is_delete"]),
			Source:       info["source"],
		}
		users = append(users, u)
	}
	return users, nil
}

func DeleteRoomById(roomId string) (bool, error) {
	num, _, err := db.DeleteRoomById(roomId)
	if err != nil {
		return false, err
	}
	if num < 1 {
		return false, nil
	}
	return true, nil
}

func DeleteRoomMemberById(userId, roomId string, time int64) (bool, error) {
	num, _, err := db.DeleteRoomMemberById(userId, roomId, time)
	if err != nil {
		return false, err
	}
	if num < 1 {
		return false, nil
	}
	return true, nil
}

//获取群中成员禁言详情
func GetRoomUserMuted(roomId, userId string) (mutedType int, deadline int64, err error) {
	mutedMaps, err := db.GetRoomUserMuted(roomId, userId)
	if err == nil && len(mutedMaps) > 0 {
		mutedInfo := mutedMaps[0]
		mutedType = utility.ToInt(mutedInfo["list_type"])
		deadline = utility.ToInt64(mutedInfo["deadline"])
	} else {
		return 1, 0, err
	}
	return
}

//获取群的某种禁言信息列表
func GetMutedListByType(roomId string, mutedType int) ([]*types.RoomUserMuted, error) {
	maps, err := db.GetRoomUsersMutedInfo(roomId, mutedType)
	if err != nil {
		return nil, err
	}
	mutedInfos := make([]*types.RoomUserMuted, 0)
	for _, v := range maps {
		info := &types.RoomUserMuted{
			Id:       v["id"],
			RoomId:   v["room_id"],
			UserId:   v["user_id"],
			ListType: utility.ToInt(v["list_type"]),
			Deadline: utility.ToInt64(v["deadline"]),
		}
		mutedInfos = append(mutedInfos, info)
	}
	return mutedInfos, nil
}

// 设置群禁言类型
func SetRoomMutedType(tx types.Tx, roomId string, mutedType int) error {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	_, _, err := db.SetRoomMutedType(tx.(*mysql.MysqlTx), roomId, mutedType)
	return err
}

func ClearMutedList(tx types.Tx, roomId string) error {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	_, _, err := db.ClearRoomMutedList(tx.(*mysql.MysqlTx), roomId)
	return err
}

func AddMutedMember(tx types.Tx, roomId, userId string, mutedType int, deadline int64) (int64, error) {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	_, id, err := db.AddRoomUserMuted(tx.(*mysql.MysqlTx), roomId, userId, mutedType, deadline)
	return id, err
}

func DelMemberMuted(tx types.Tx, roomId, userId string) error {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	_, _, err := db.DelRoomUserMuted(tx.(*mysql.MysqlTx), roomId, userId)
	return err
}

// 获取群禁言类型
func GetRoomMutedType(roomId string) (int, error) {
	maps, err := db.GetRoomMutedType(roomId)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	return utility.ToInt(maps[0]["master_muted"]), nil
}

func GetJoinedRoom(userId string, Type int) ([]*types.MemberJoinRoom, error) {
	maps, err := db.GetRoomList(userId, Type)
	if err != nil {
		return nil, err
	}

	var roomList = make([]*types.MemberJoinRoom, 0)
	for _, info := range maps {
		roomMember := &types.MemberJoinRoom{
			Room: convertRoom(info),
			RoomMember: &types.RoomMember{
				RoomId:       info["room_id"],
				UserId:       info["user_id"],
				UserNickname: info["user_nickname"],
				Level:        utility.ToInt(info["level"]),
				NoDisturbing: utility.ToInt(info["no_disturbing"]),
				CommonUse:    utility.ToInt(info["common_use"]),
				RoomTop:      utility.ToInt(info["room_top"]),
				CreateTime:   utility.ToInt64(info["create_time"]),
			},
		}
		roomList = append(roomList, roomMember)
	}
	return roomList, nil
}

func SetCanAddFriendPermission(roomId string, permisson int) error {
	return db.AlterRoomCanAddFriendPermission(roomId, permisson)
}

func SetJoinPermission(roomId string, permisson int) error {
	return db.AlterRoomJoinPermission(roomId, permisson)
}

func SetRecordPermission(roomId string, permisson int) error {
	return db.AlterRoomRecordPermission(roomId, permisson)
}

func SetRoomName(roomId, name string) (bool, error) {
	num, _, err := db.AlterRoomName(roomId, name)
	if err != nil {
		return false, err
	}
	if num < 1 {
		return false, err
	}
	return true, nil
}

func SetAvatar(roomId, avatar string) (bool, error) {
	num, _, err := db.AlterRoomAvatar(roomId, avatar)
	if err != nil {
		return false, err
	}
	if num < 1 {
		return false, err
	}
	return true, nil
}

func SetMemberLevel(userId, roomId string, level int) error {
	_, _, err := db.SetRoomMemberLevel(userId, roomId, level)
	return err
}

func SetNewMaster(master, userId, roomId string, level int) error {
	return db.SetNewMaster(master, userId, roomId, level)
}

func SetNoDisturbing(userId, roomId string, noDisturbing int) error {
	_, _, err := db.SetRoomNoDisturbing(userId, roomId, noDisturbing)
	return err
}

func SetOnTop(userId, roomId string, onTop int) error {
	_, _, err := db.SetRoomOnTop(userId, roomId, onTop)
	return err
}

func SetMemberNickname(userId, roomId string, nickname string) error {
	_, _, err := db.SetMemberNickname(userId, roomId, nickname)
	return err
}

func AddMember(tx types.Tx, userId, roomId string, memberLevel int, createTime int64, source string) (int64, error) {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	_, id, err := db.RoomAddMember(tx.(*mysql.MysqlTx), userId, roomId, memberLevel, createTime, source)
	return id, err
}

func FindRoomChatLogByContentId(logId string) (*types.RoomLogJoinUser, error) {
	maps, err := db.GetRoomMsgContent(logId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	v := maps[0]
	return &types.RoomLogJoinUser{
		RoomLog: &types.RoomLog{
			Id:       v["id"],
			MsgId:    v["msg_id"],
			RoomId:   v["room_id"],
			SenderId: v["sender_id"],
			IsSnap:   utility.ToInt(v["is_snap"]),
			MsgType:  utility.ToInt(v["msg_type"]),
			Content:  v["content"],
			Datetime: utility.ToInt64(v["datetime"]),
			Ext:      v["ext"],
			IsDelete: utility.ToInt(v["is_delete"]),
		},
		User: &types.User{
			Username: v["username"],
			Avatar:   v["avatar"],
			Uid:      v["uid"],
		},
	}, nil
}

func FindRoomChatLogByMsgId(senderId, msgId string) (*types.RoomLogJoinUser, error) {
	maps, err := db.GetRoomMsgContentByMsgId(senderId, msgId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	v := maps[0]
	return &types.RoomLogJoinUser{
		RoomLog: &types.RoomLog{
			Id:       v["id"],
			MsgId:    v["msg_id"],
			RoomId:   v["room_id"],
			SenderId: v["sender_id"],
			IsSnap:   utility.ToInt(v["is_snap"]),
			MsgType:  utility.ToInt(v["msg_type"]),
			Content:  v["content"],
			Datetime: utility.ToInt64(v["datetime"]),
			Ext:      v["ext"],
			IsDelete: utility.ToInt(v["is_delete"]),
		},
		User: &types.User{
			Username: v["username"],
			Avatar:   v["avatar"],
		},
	}, nil
}

func FindRoomChatLogsAfterTime(rooms []string, time int64) ([]*types.RoomLogJoinUser, error) {
	maps, err := db.GetRoomMsgContentAfter(rooms, time, types.RoomMsgNotDelete)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.RoomLogJoinUser, 0)
	for _, v := range maps {
		log := &types.RoomLogJoinUser{
			RoomLog: &types.RoomLog{
				Id:       v["id"],
				MsgId:    v["msg_id"],
				RoomId:   v["room_id"],
				SenderId: v["sender_id"],
				IsSnap:   utility.ToInt(v["is_snap"]),
				MsgType:  utility.ToInt(v["msg_type"]),
				Content:  v["content"],
				Datetime: utility.ToInt64(v["datetime"]),
				Ext:      v["ext"],
				IsDelete: utility.ToInt(v["is_delete"]),
			},
			User: &types.User{
				Username: v["username"],
				Avatar:   v["avatar"],
			},
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindRoomChatLogsNumberBetween(rooms []string, begin, end int64) (int64, error) {
	maps, err := db.GetRoomMsgContentNumberBetween(rooms, begin, end, types.RoomMsgNotDelete)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

func FindRoomChatLogsBetweenTime(rooms []string, begin, end int64) ([]*types.RoomLogJoinUser, error) {
	maps, err := db.GetRoomMsgContentBetween(rooms, begin, end, types.RoomMsgNotDelete)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.RoomLogJoinUser, 0)
	for _, v := range maps {
		log := &types.RoomLogJoinUser{
			RoomLog: &types.RoomLog{
				Id:       v["id"],
				MsgId:    v["msg_id"],
				RoomId:   v["room_id"],
				SenderId: v["sender_id"],
				IsSnap:   utility.ToInt(v["is_snap"]),
				MsgType:  utility.ToInt(v["msg_type"]),
				Content:  v["content"],
				Datetime: utility.ToInt64(v["datetime"]),
				Ext:      v["ext"],
				IsDelete: utility.ToInt(v["is_delete"]),
			},
			User: &types.User{
				Username: v["username"],
				Avatar:   v["avatar"],
			},
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindRoomLogs(roomId string, startId, joinTime int64, number int) (int64, []*types.RoomLogJoinUser, error) {
	maps, nextLog, err := db.GetChatlog(roomId, startId, joinTime, number)
	if err != nil {
		return -1, nil, err
	}
	logs := make([]*types.RoomLogJoinUser, 0)
	for _, v := range maps {
		log := &types.RoomLogJoinUser{
			RoomLog: &types.RoomLog{
				Id:       v["id"],
				MsgId:    v["msg_id"],
				RoomId:   v["room_id"],
				SenderId: v["sender_id"],
				IsSnap:   utility.ToInt(v["is_snap"]),
				MsgType:  utility.ToInt(v["msg_type"]),
				Content:  v["content"],
				Datetime: utility.ToInt64(v["datetime"]),
				Ext:      v["ext"],
				IsDelete: utility.ToInt(v["is_delete"]),
			},
			User: &types.User{
				Username: v["username"],
				Avatar:   v["avatar"],
			},
		}
		logs = append(logs, log)
	}
	return utility.ToInt64(nextLog), logs, nil
}

func FindRoomLogsById(roomId, owner string, startId, joinTime int64, number int, queryType []string) (int64, []*types.RoomLogJoinUser, error) {
	maps, nextLog, err := db.GetRoomChatLogsByUserId(roomId, owner, startId, joinTime, number, queryType)
	if err != nil {
		return -1, nil, err
	}
	logs := make([]*types.RoomLogJoinUser, 0)
	for _, v := range maps {
		log := &types.RoomLogJoinUser{
			RoomLog: &types.RoomLog{
				Id:       v["id"],
				MsgId:    v["msg_id"],
				RoomId:   v["room_id"],
				SenderId: v["sender_id"],
				IsSnap:   utility.ToInt(v["is_snap"]),
				MsgType:  utility.ToInt(v["msg_type"]),
				Content:  v["content"],
				Datetime: utility.ToInt64(v["datetime"]),
				Ext:      v["ext"],
				IsDelete: utility.ToInt(v["is_delete"]),
			},
			User: &types.User{
				Username: v["username"],
				Avatar:   v["avatar"],
			},
		}
		logs = append(logs, log)
	}
	return utility.ToInt64(nextLog), logs, nil
}

func DelRoomChatLogById(logId string) (int, error) {
	return db.DeleteRoomMsgContent(logId)
}

func AlertRoomRevStateByRevId(revId string, state int) error {
	_, _, err := db.AlertRoomRevStateByRevId(revId, state)
	return err
}

func AppendRoomChatLog(userId, roomId, msgId string, msgType, isSnap int, content, ext string, time int64) (int64, int64, error) {
	return db.AppendRoomChatLog(userId, roomId, msgId, msgType, isSnap, content, ext, time)
}

func FindReceiveLogById(logId, userId string) (*types.RoomMsgReceive, error) {
	maps, err := db.GetRoomRevMsgByLogId(logId, userId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	log := maps[0]
	return &types.RoomMsgReceive{
		Id:        log["id"],
		RoomMsgId: log["room_msg_id"],
		ReceiveId: log["receive_id"],
		State:     utility.ToInt(log["state"]),
	}, nil
}

func FindSystemMsg(roomId string, startId int64, number int) (int64, []*types.RoomLog, error) {
	maps, nextId, err := db.GetSystemMsg(roomId, startId, number)
	if err != nil {
		return -1, nil, err
	}

	list := make([]*types.RoomLog, 0)
	for _, v := range maps {
		info := &types.RoomLog{
			Id:       v["id"],
			MsgId:    v["msg_id"],
			RoomId:   v["room_id"],
			SenderId: v["sender_id"],
			IsSnap:   utility.ToInt(v["is_snap"]),
			MsgType:  utility.ToInt(v["msg_type"]),
			Content:  v["content"],
			Datetime: utility.ToInt64(v["datetime"]),
			Ext:      v["ext"],
			IsDelete: utility.ToInt(v["is_delete"]),
		}
		/*info.Id = r["id"]
		info.SenderId = r["sender_id"]
		info.Content = utility.ToString(utility.StringToJobj(r["content"])["content"])
		info.Datetime = utility.ToInt64(r["datetime"])*/
		list = append(list, info)
	}
	return utility.ToInt64(nextId), list, nil
}

func GetSystemMsgNumber(roomId string) (int64, error) {
	return db.GetRoomSystemMsgNumber(roomId)
}

func GetRoomMsgBurntNumber(logId string) (int, error) {
	return db.GetRoomMsgBurntNumber(logId)
}

func CreateNewRoom(creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, members []string, randomRoomId string, createTime int64) (int64, error) {
	return db.CreateRoom(creater, roomName, roomAvatar, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, members, randomRoomId, createTime)
}

func CreateNewRoomV2(tx types.Tx, creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, randomRoomId string, createTime int64) (int64, error) {
	return db.CreateRoomV2(tx.(*mysql.MysqlTx), creater, roomName, roomAvatar, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, randomRoomId, createTime)
}

func FindUserCreateRoomsNumber(userId string) (int, error) {
	return db.FindUserCreateRoomsNumber(userId)
}

func GetMemberApplyInfo(roomId, userId string) (*types.Apply, error) {
	maps, err := db.GetRoomUserApplyInfo(roomId, userId)
	if err != nil || len(maps) < 1 {
		return nil, err
	}
	info := maps[0]

	return &types.Apply{
		Id:          info["id"],
		Type:        utility.ToInt(info["type"]),
		ApplyUser:   info["apply_user"],
		Target:      info["target"],
		ApplyReason: info["apply_reason"],
		State:       utility.ToInt(info["state"]),
		Remark:      info["remark"],
		Datetime:    utility.ToInt64(info["datetime"]),
		Source:      info["source"],
	}, nil
}

func ApproveInsertMemberStep(tx types.Tx, roomId, userId, source string) (bool, int64, error) {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	num, id, err := db.JoinRoomApproveStepInsert(tx.(*mysql.MysqlTx), roomId, userId, source)
	if err != nil {
		return false, 0, err
	}
	if num < 1 {
		return false, 0, nil
	}
	return true, id, nil
}

func ApproveChangeStateStep(tx types.Tx, logId int64, status int) (bool, error) {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	num, _, err := db.JoinRoomApproveStepChangeState(tx.(*mysql.MysqlTx), logId, status)
	if err != nil {
		return false, err
	}
	if num < 1 {
		return false, nil
	}
	return true, nil
}

//添加群成员消息接收日志
func AddMemberReceiveLog(logId, receiver string, state int) (int64, error) {
	_, id, err := db.AppendRoomMemberReceiveLog(logId, receiver, state)
	return id, err
}

// 获取群禁言数量
func GetRoomMutedNumber(roomId string) (int64, error) {
	maps, err := db.GetRoomMutedListNumber(roomId)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

// 获取群禁言数量 事务
func GetRoomMutedNumberTx(tx types.Tx, roomId string) (int64, error) {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	maps, err := db.GetRoomMutedListNumberByTx(tx.(*mysql.MysqlTx), roomId)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

func GetTx() (*mysql.MysqlTx, error) {
	return db.GetNewTx()
}

//根据appId获取应用的所有未解散的群聊，包括被封群
func GetRoomCountInApp(appId string) (int64, error) {
	return db.GetRoomCountInApp(appId)
}

func GetCloseRoomCountInApp(appId string) (int64, error) {
	return db.GetCloseRoomCountInApp(appId)
}

//获取用户创建群个数上限
func GetCreateRoomsLimit(appId string, level int) (int, error) {
	return db.GetCreateRoomsLimit(appId, level)
}

//获取群的成员数上限
func GetRoomMembersLimit(appId string, level int) (int, error) {
	return db.GetRoomMembersLimit(appId, level)
}

//设置用户创建群个数上限
func SetCreateRoomsLimit(appId string, level, limit int) error {
	return db.SetCreateRoomsLimit(appId, level, limit)
}

//设置群的成员数上限
func SetRoomMembersLimit(appId string, level, limit int) error {
	return db.SetRoomMembersLimit(appId, level, limit)
}

//创建群个数 未解散
func FindCreateRoomNumbers(masterId string) (int, error) {
	return db.FindCreateRoomNumbers(masterId)
}

//查询某个app的所有群，模糊查询 markId 群名
func FindRoomsInAppQueryMarkId(appId, query string) ([]*types.RoomJoinUser, error) {
	maps, err := db.FindRoomsInAppQueryMarkId(appId, query)
	if err != nil {
		return nil, err
	}
	list := make([]*types.RoomJoinUser, 0)
	for _, info := range maps {
		item := &types.RoomJoinUser{
			Room: convertRoom(info),
			User: convertUser(info),
		}
		list = append(list, item)
	}
	return list, nil
}

//查找某个app所有下未封禁群
func FindRoomsInAppUnClose(appId string) ([]*types.RoomJoinUser, error) {
	maps, err := db.FindRoomsInAppUnClose(appId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.RoomJoinUser, 0)
	for _, info := range maps {
		item := &types.RoomJoinUser{
			Room: convertRoom(info),
			User: convertUser(info),
		}
		list = append(list, item)
	}
	return list, nil
}

//查找某个app所有下封禁群
func FindRoomsInAppClosed(appId string) ([]*types.RoomJoinUser, error) {
	maps, err := db.FindRoomsInAppClosed(appId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.RoomJoinUser, 0)
	for _, info := range maps {
		item := &types.RoomJoinUser{
			Room: convertRoom(info),
			User: convertUser(info),
		}
		list = append(list, item)
	}
	return list, nil
}

//查找某个app所有下推荐群（带 群主、人数、封群次数等信息）
func FindRoomsInAppRecommend(appId string) ([]*types.RoomJoinUser, error) {
	maps, err := db.FindAllRecommendRooms(appId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.RoomJoinUser, 0)
	for _, info := range maps {
		item := &types.RoomJoinUser{
			Room: convertRoom(info),
			User: convertUser(info),
		}
		list = append(list, item)
	}
	return list, nil
}

//根据群的发言人数递减排序
func RoomsOrderActiveMember(appId string, datetime int64) ([]*types.Room, error) {
	maps, err := db.RoomsOrderActiveMember(appId, datetime)
	if err != nil {
		return nil, err
	}
	list := make([]*types.Room, 0)
	for _, m := range maps {
		item := convertRoom(m)
		list = append(list, item)
	}

	return list, nil
}

//根据群的发言条数递减排序
func RoomsOrderActiveMsg(appId string, datetime int64) ([]*types.Room, error) {
	maps, err := db.RoomsOrderActiveMsg(appId, datetime)
	if err != nil {
		return nil, err
	}
	list := make([]*types.Room, 0)
	for _, m := range maps {
		item := convertRoom(m)
		list = append(list, item)
	}

	return list, nil
}

func FindAllRecommendRooms(appId string) ([]*types.Room, error) {
	maps, err := db.FindAllRecommendRooms(appId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.Room, 0)
	for _, m := range maps {
		item := convertRoom(m)
		list = append(list, item)
	}

	return list, nil
}
