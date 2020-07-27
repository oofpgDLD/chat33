package orm

import (
	"github.com/33cn/chat33/db"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	redis "github.com/33cn/chat33/orm/redis_model"
	sql "github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logRoom = log15.New("module", "model/room")

//群成员 等级：0 普通用户 1 管理员 2 群主
func GetMemberLevel(roomId, userId string, isDel int) int {
	if cfg.CacheType.Enable {
		ret, err := redis.GetRoomUserLevel(roomId, userId, isDel)
		if err != nil {
			logRoom.Error("redis.GetMemberLevel", "err", err, "roomId", roomId, "userId", userId, "isDel", isDel)
		}
		return ret
	}
	ret, err := mysql.GetRoomUserLevel(roomId, userId, isDel)
	if err != nil {
		logRoom.Error("mysql.GetMemberLevel", "err", err, "roomId", roomId, "userId", userId, "isDel", isDel)
	}
	return ret
}

func GetTx() (types.Tx, error) {
	return mysql.GetTx()
}

//新增群成员，mysql操作成功后进行redis操作
func AddMember(tx types.Tx, userId, roomId string, memberLevel int, createTime int64, source string) error {
	id, err := mysql.AddMember(tx, userId, roomId, memberLevel, createTime, source)
	if err != nil {
		logRoom.Error("mysql.AddMember", "err", err, "roomId", roomId, "userId", userId, "memberLevel", memberLevel, "createTime", createTime, "source", source)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.AddMember(tx, userId, roomId, memberLevel, id, createTime, source)
		if err != nil {
			logRoom.Error("redis.AddMember", "err", err, "roomId", roomId, "userId", userId, "memberLevel", memberLevel, "createTime", createTime, "source", source)
			return err
		}
	}
	return err
}

func FindUserCreateRoomsNumber(userId string) (int, error) {
	ret, err := mysql.FindUserCreateRoomsNumber(userId)
	if err != nil {
		logRoom.Error("mysql.FindUserCreateRoomsNumber", "err", err, "userId", userId)
	}
	return ret, err
}

//群成员 群昵称 > 名称
func GetMemberName(roomId, userId string) string {
	if cfg.CacheType.Enable {
		ret, err := redis.GetRoomMemberName(roomId, userId)
		if err != nil {
			logRoom.Error("redis.GetRoomMemberName", "err", err, "userId", userId)
		}
		return ret
	}
	ret, err := mysql.GetRoomMemberName(roomId, userId)
	if err != nil {
		logRoom.Error("mysql.GetRoomMemberName", "err", err, "userId", userId)
	}
	return ret
}

func SetCanAddFriendPermission(roomId string, permisson int) error {
	err := mysql.SetCanAddFriendPermission(roomId, permisson)
	if err != nil {
		logRoom.Error("mysql.SetCanAddFriendPermission", "err", err, "roomId", roomId, "permisson", permisson)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetPermission(roomId, permisson, 0, 0, types.RoomNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetPermission", "err", err, "roomId", roomId, "permisson", permisson)
			return err
		}
	}
	return nil
}

func SetJoinPermission(roomId string, permisson int) error {
	err := mysql.SetJoinPermission(roomId, permisson)
	if err != nil {
		logRoom.Error("mysql.SetJoinPermission", "err", err, "roomId", roomId, "permisson", permisson)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetPermission(roomId, 0, permisson, 0, types.RoomNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetPermission", "err", err, "roomId", roomId, "permisson", permisson)
			return err
		}
	}

	return nil
}

func SetRecordPermission(roomId string, permisson int) error {
	err := mysql.SetRecordPermission(roomId, permisson)
	if err != nil {
		logRoom.Error("mysql.SetRecordPermission", "err", err, "roomId", roomId, "permisson", permisson)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetPermission(roomId, 0, 0, permisson, types.RoomNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetPermission", "err", err, "roomId", roomId, "permisson", permisson)
			return err
		}
	}
	return err
}

func SetRoomName(roomId, name string) (bool, error) {
	b, err := mysql.SetRoomName(roomId, name)
	if err != nil {
		logRoom.Error("mysql.SetRoomName", "err", err, "roomId", roomId, "name", name)
		return b, err
	}
	if cfg.CacheType.Enable && b {
		d, err := redis.SetRoomName(roomId, name, types.RoomNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetRoomName", "err", err, "roomId", roomId, "name", name)
			return d, err
		}
	}

	return b, err
}

func SetAvatar(roomId, avatar string) (bool, error) {
	b, err := mysql.SetAvatar(roomId, avatar)
	if err != nil {
		logRoom.Error("mysql.SetAvatar", "err", err, "roomId", roomId, "avatar", avatar)
		return b, err
	}
	if cfg.CacheType.Enable && b {
		d, err := redis.SetAvatar(roomId, avatar, types.RoomNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetAvatar", "err", err, "roomId", roomId, "avatar", avatar)
			return d, err
		}
	}
	return b, err
}

func SetMemberLevel(userId, roomId string, level int) error {
	err := mysql.SetMemberLevel(userId, roomId, level)
	if err != nil {
		logRoom.Error("mysql.SetMemberLevel", "err", err, "userId", userId, "roomId", roomId, "level", level)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetMemberLevel(userId, roomId, level, types.RoomUserNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetMemberLevel", "err", err, "userId", userId, "roomId", roomId, "level", level)
			return err
		}
	}
	return nil
}

func SetNewMaster(master, userId, roomId string, level int) error {
	err := mysql.SetNewMaster(master, userId, roomId, level)
	if err != nil {
		logRoom.Error("mysql.SetNewMaster", "err", err, "userId", userId, "roomId", roomId, "level", level)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetNewMaster(master, userId, roomId, level)
		if err != nil {
			logRoom.Error("redis.SetNewMaster", "err", err, "userId", userId, "roomId", roomId, "level", level)
			return err
		}
	}
	return nil
}

func SetNoDisturbing(userId, roomId string, permission int) error {
	err := mysql.SetNoDisturbing(userId, roomId, permission)
	if err != nil {
		logRoom.Error("mysql.SetNoDisturbing", "err", err, "userId", userId, "roomId", roomId, "permission", permission)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetNoDisturbing(userId, roomId, permission, types.RoomUserNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetNoDisturbing", "err", err, "userId", userId, "roomId", roomId, "permission", permission)
			return err
		}
	}

	return nil
}

func SetOnTop(userId, roomId string, permission int) error {
	err := mysql.SetOnTop(userId, roomId, permission)
	if err != nil {
		logRoom.Error("mysql.SetOnTop", "err", err, "userId", userId, "roomId", roomId, "permission", permission)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetOnTop(userId, roomId, permission, types.RoomUserNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetOnTop", "err", err, "userId", userId, "roomId", roomId, "permission", permission)
			return err
		}
	}
	return nil
}

func SetMemberNickname(userId, roomId string, nickname string) error {
	err := mysql.SetMemberNickname(userId, roomId, nickname)
	if err != nil {
		logRoom.Error("mysql.SetMemberNickname", "err", err, "userId", userId, "roomId", roomId, "nickname", nickname)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetMemberNickname(userId, roomId, nickname, types.RoomUserNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetMemberNickname", "err", err, "userId", userId, "roomId", roomId, "nickname", nickname)
			return err
		}
	}
	return nil
}

//清除群禁言
func ClearMutedList(tx types.Tx, roomId string) error {
	err := mysql.ClearMutedList(tx, roomId)
	if err != nil {
		logRoom.Error("mysql.SetOnTop", "err", err, "roomId", roomId)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.ClearMutedList(tx, roomId)
		if err != nil {
			logRoom.Error("mysql.SetOnTop", "err", err, "roomId", roomId)
			return err
		}
	}
	return nil
}

//新增返回最后新增的id
func AddMutedMember(tx types.Tx, roomId, userId string, mutedType int, deadline int64) (int64, error) {
	id, err := mysql.AddMutedMember(tx, roomId, userId, mutedType, deadline)
	if err != nil {
		logRoom.Error("mysql.AddMutedMember", "err", err, "userId", userId, "roomId", roomId, "mutedType", mutedType, "deadline", deadline)
		return id, err
	}
	if cfg.CacheType.Enable {
		err = redis.AddMutedMember(id, roomId, userId, mutedType, deadline)
		if err != nil {
			logRoom.Error("redis.AddMutedMember", "err", err, "id", id, "userId", userId, "roomId", roomId, "mutedType", mutedType, "deadline", deadline)
			return 0, err
		}
	}
	return id, err
}

func DelMemberMuted(tx types.Tx, roomId, userId string) error {
	err := mysql.DelMemberMuted(tx, roomId, userId)
	if err != nil {
		logRoom.Error("mysql.DelMemberMuted", "err", err, "userId", userId, "roomId", roomId)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.DelMemberMuted(roomId, userId)
		if err != nil {
			logRoom.Error("redis.DelMemberMuted", "err", err, "userId", userId, "roomId", roomId)
			return err
		}
	}

	return nil
}

//获取禁言类型
func GetRoomMutedType(roomId string) (int, error) {
	//if cfg.CacheType.Enable {
	//	return redis.GetRoomMutedType(roomId, types.RoomNotDeleted)
	//}
	ret, err := mysql.GetRoomMutedType(roomId)
	if err != nil {
		logRoom.Error("mysql.GetRoomMutedType", "err", err, "roomId", roomId)
	}
	return ret, err
}

func SetRoomMutedType(tx types.Tx, roomId string, mutedType int) error {
	err := mysql.SetRoomMutedType(tx, roomId, mutedType)
	if err != nil {
		logRoom.Error("mysql.SetRoomMutedType", "err", err, "roomId", roomId, "mutedType", mutedType)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetRoomMutedType(roomId, mutedType, types.RoomNotDeleted)
		if err != nil {
			logRoom.Error("redis.SetRoomMutedType", "err", err, "roomId", roomId, "mutedType", mutedType)
			return err
		}
	}

	return nil
}

func GetRoomUserMuted(roomId, userId string) (mutedType int, deadline int64) {
	//if cfg.CacheType.Enable {
	//	return redis.GetRoomUserMuted(roomId, userId)
	//}
	ret, ret2, err := mysql.GetRoomUserMuted(roomId, userId)
	if err != nil {
		logRoom.Error("mysql.GetRoomUserMuted", "err", err, "roomId", roomId, "userId", userId)
	}
	return ret, ret2
}

func GetMutedListByType(roomId string, mutedType int) ([]*types.RoomUserMuted, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetMutedListByType(roomId, mutedType)
		if err != nil {
			logRoom.Error("redis.GetMutedListByType", "err", err, "roomId", roomId, "mutedType", mutedType)
		}
		return ret, err
	}
	ret, err := mysql.GetMutedListByType(roomId, mutedType)
	if err != nil {
		logRoom.Error("mysql.GetMutedListByType", "err", err, "roomId", roomId, "mutedType", mutedType)
	}
	return ret, err
}

//群中被禁言数量
func GetMutedCount(roomId string) (int64, error) {
	ret, err := mysql.GetRoomMutedNumber(roomId)
	if err != nil {
		logRoom.Error("mysql.GetRoomMutedNumber", "err", err, "roomId", roomId)
	}
	return ret, err
}

func GetMutedCountTx(tx types.Tx, roomId string) (int64, error) {
	ret, err := mysql.GetRoomMutedNumberTx(tx, roomId)
	if err != nil {
		logRoom.Error("mysql.GetRoomMutedNumberTx", "err", err, "roomId", roomId)
	}
	return ret, err
}

//参数：room message content log id
func FindReceiveLogById(logId, userId string) (*types.RoomMsgReceive, error) {
	//if cfg.CacheType.Enable {
	//	return redis.FindReceiveLogById(logId, userId)
	//}
	ret, err := mysql.FindReceiveLogById(logId, userId)
	if err != nil {
		logRoom.Error("mysql.FindReceiveLogById", "err", err, "logId", logId, "userId", userId)
	}
	return ret, err
}

func AppendMemberRevLog(userId, logId string, state int) (int64, error) {
	//if cfg.CacheType.Enable {
	//	err = redis.AppendMemberRevLog(id, userId, logId, state)
	//	if err != nil {
	//		return 0, err
	//	}
	//}
	ret, err := mysql.AddMemberReceiveLog(logId, userId, state)
	if err != nil {
		logRoom.Error("mysql.AddMemberReceiveLog", "err", err, "userId", userId, "logId", logId, "state", state)
	}
	return ret, err
}

func GetMemberNumber(roomId string) (int64, error) {
	ret, err := mysql.GetMemberNumber(roomId)
	if err != nil {
		logRoom.Error("mysql.GetMemberNumber", "err", err, "roomId", roomId)
	}
	return ret, err
}

func GetRoomMasterNumber(roomId string) (int64, error) {
	ret, err := mysql.GetRoomMemberNumberByLevel(roomId, types.RoomLevelManager)
	if err != nil {
		logRoom.Error("mysql.GetRoomMemberNumberByLevel", "err", err, "roomId", roomId)
	}
	return ret, err
}

func GetUserJoinedRooms(userId string) ([]string, error) {
	ret, err := mysql.GetUserJoinedRooms(userId)
	if err != nil {
		logRoom.Error("mysql.GetUserJoinedRooms", "err", err, "userId", userId)
	}
	return ret, err
}

func GetEnableRoomsId() ([]string, error) {
	ret, err := mysql.GetEnableRoomsId()
	if err != nil {
		logRoom.Error("mysql.GetEnableRoomsId", "err", err)
	}
	return ret, err
}

func GetRoomManagerAndMaster(roomId string) ([]*types.MemberJoinUser, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetRoomManagerAndMaster(roomId)
		if err != nil {
			logRoom.Error("redis.GetRoomManagerAndMaster", "err", err, "roomId", roomId)
		}
		return ret, err
	}
	ret, err := mysql.GetRoomManagerAndMaster(roomId)
	if err != nil {
		logRoom.Error("mysql.GetRoomManagerAndMaster", "err", err, "roomId", roomId)
	}
	return ret, err
}

func FindNotDelMember(roomId, userId string) (*types.MemberJoinUser, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindRoomMemberById(roomId, userId, types.RoomUserNotDeleted)
		if err != nil {
			logRoom.Error("redis.FindRoomMemberById", "err", err, "roomId", roomId, "userId", userId)
		}
		return ret, err
	}
	ret, err := mysql.FindRoomMemberById(roomId, userId, types.RoomUserNotDeleted)
	if err != nil {
		logRoom.Error("mysql.FindRoomMemberById", "err", err, "roomId", roomId, "userId", userId)
	}
	return ret, err
}

func FindJoinedRooms(userId string) ([]*types.RoomMember, error) {
	ret, err := mysql.FindJoinedRooms(userId)
	if err != nil {
		logRoom.Error("mysql.FindJoinedRooms", "err", err, "userId", userId)
	}
	return ret, err
}

func FindMemberByName(roomId, name string) ([]*types.MemberJoinUser, error) {
	ret, err := mysql.FindRoomMemberByName(roomId, name)
	if err != nil {
		logRoom.Error("mysql.FindRoomMemberByName", "err", err, "roomId", roomId, "name", name)
	}
	return ret, err
}

func FindNotDelMembers(roomId string, searchNumber int) ([]*types.MemberJoinUser, error) {
	ret, err := mysql.FindRoomMembers(roomId, searchNumber)
	if err != nil {
		logRoom.Error("mysql.FindRoomMembers", "err", err, "roomId", roomId, "searchNumber", searchNumber)
	}
	return ret, err
}

func FindSetNoDisturbingMembers(roomId string) ([]*types.RoomMember, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindSetNoDisturbingMembers(roomId)
		if err != nil {
			logRoom.Error("redis.FindSetNoDisturbingMembers", "err", err, "roomId", roomId)
		}
		return ret, err
	}
	ret, err := mysql.FindSetNoDisturbingMembers(roomId)
	if err != nil {
		logRoom.Error("mysql.FindSetNoDisturbingMembers", "err", err, "roomId", roomId)
	}
	return ret, err
}

//检查群markid 是否存在
func CheckRoomMarkIdExist(markId string) (bool, error) {
	ret, err := mysql.CheckRoomMarkIdExist(markId)
	if err != nil {
		logRoom.Error("mysql.CheckRoomMarkIdExist", "err", err, "markId", markId)
	}
	return ret, err
}

//返回：roomId,error
func CreateNewRoom(creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, members []string, randomRoomId string, createTime int64) (int64, error) {
	ret, err := mysql.CreateNewRoom(creater, roomName, roomAvatar, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, members, randomRoomId, createTime)
	if err != nil {
		logRoom.Error("mysql.CreateNewRoom", "err", err)
	}
	return ret, err
}

//返回：roomId,error
func CreateNewRoomV2(tx types.Tx, creater, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, randomRoomId string, createTime int64) (int64, error) {
	roomid, err := mysql.CreateNewRoomV2(tx, creater, roomName, roomAvatar, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, randomRoomId, createTime)
	if err != nil {
		logRoom.Error("mysql.CreateNewRoomV2", "err", err)
		return 0, err
	}
	if cfg.CacheType.Enable {
		err = redis.CreateNewRoom(tx, roomid, creater, roomName, roomAvatar, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, randomRoomId, createTime)
		if err != nil {
			logRoom.Error("redis.CreateNewRoomV2", "err", err)
			return 0, err
		}
	}
	return roomid, err
}

func FindRoomById(id string, isDel int) (*types.Room, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.SearchRoomInfo(id, isDel)
		if err != nil {
			logRoom.Error("redis.SearchRoomInfo", "err", err, "roomId", id)
		}
		return ret, err
	}
	ret, err := mysql.SearchRoomInfo(id, isDel)
	if err != nil {
		logRoom.Error("mysql.SearchRoomInfo", "err", err, "roomId", id)
	}
	return ret, err
}

func FindRoomByMarkId(markId string, isDel int) (*types.Room, error) {
	ret, err := mysql.SearchRoomInfoByMarkId(markId, isDel)
	if err != nil {
		logRoom.Error("mysql.SearchRoomInfoByMarkId", "err", err, "markId", markId)
	}
	return ret, err
}

func DelRoomById(roomId string) (bool, error) {
	rlt, err := mysql.DeleteRoomById(roomId)
	if err != nil {
		logRoom.Error("mysql.DeleteRoomById", "err", err, "roomId", roomId)
		return false, err
	}
	if cfg.CacheType.Enable && rlt {
		err = redis.DeleteRoomInfoById(roomId)
		if err != nil {
			logRoom.Error("redis.DeleteRoomInfoById", "err", err, "roomId", roomId)
			return false, err
		}
	}
	return rlt, err
}

func GetJoinedRooms(user string, Type int) ([]*types.MemberJoinRoom, error) {
	ret, err := mysql.GetJoinedRoom(user, Type)
	if err != nil {
		logRoom.Error("mysql.GetJoinedRoom", "err", err, "userId", user, "Type", Type)
	}
	return ret, err
}

func DelRoomMemberById(userId, roomId string, time int64) (bool, error) {
	rlt, err := mysql.DeleteRoomMemberById(userId, roomId, time)
	if err != nil {
		logRoom.Error("mysql.DelRoomMemberById", "err", err, "userId", userId, "roomId", roomId, "time", time)
		return false, err
	}
	if cfg.CacheType.Enable && rlt {
		_, err := redis.DeleteRoomMemberById(userId, roomId)
		if err != nil {
			logRoom.Error("redis.DeleteRoomMemberById", "err", err, "userId", userId, "roomId", roomId, "time", time)
			return false, err
		}
	}
	return rlt, err
}

func GetMemberApplyInfo(roomId, userId string) (*types.Apply, error) {
	ret, err := mysql.GetMemberApplyInfo(roomId, userId)
	if err != nil {
		logRoom.Error("mysql.GetMemberApplyInfo", "err", err, "roomId", roomId, "userId", userId)
	}
	return ret, err
}

// 入群申请，步骤1 添加user
func ApproveInsertMemberStep(tx types.Tx, roomId, userId, source string) (bool, error) {
	b, id, err := mysql.ApproveInsertMemberStep(tx, roomId, userId, source)
	if cfg.CacheType.Enable && b {
		createTime := utility.NowMillionSecond()
		err = redis.AddMember(tx, userId, roomId, types.RoomLevelNomal, id, createTime, source)
	}
	if err != nil {
		logRoom.Error("mysql.ApproveInsertMemberStep", "err", err, "roomId", roomId, "userId", userId, "source", source)
	}
	return b, err
}

//入群申请，步骤2 更改状态
func ApproveChangeStateStep(tx types.Tx, logId int64, status int) (bool, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.ApproveChangeStateStep(tx, logId, status)
		if err != nil {
			logRoom.Error("redis.ApproveChangeStateStep", "err", err, "logId", logId, "status", status)
		}
		return ret, err
	}
	ret, err := mysql.ApproveChangeStateStep(tx, logId, status)
	if err != nil {
		logRoom.Error("mysql.ApproveChangeStateStep", "err", err, "logId", logId, "status", status)
	}
	return ret, err
}

func GetSystemMsgNumber(roomId string) (int64, error) {
	ret, err := mysql.GetSystemMsgNumber(roomId)
	if err != nil {
		logRoom.Error("mysql.GetSystemMsgNumber", "err", err, "roomId", roomId)
	}
	return ret, err
}

func GetRoomMsgBurntNumber(logId string) (int, error) {
	ret, err := mysql.GetRoomMsgBurntNumber(logId)
	if err != nil {
		logRoom.Error("mysql.GetRoomMsgBurntNumber", "err", err, "logId", logId)
	}
	return ret, err
}

//倒序
func FindSystemMsg(roomId string, startId int64, number int) (int64, []*types.RoomLog, error) {
	ret, ret2, err := mysql.FindSystemMsg(roomId, startId, number)
	if err != nil {
		logRoom.Error("mysql.FindSystemMsg", "err", err, "roomId", roomId, "startId", startId, "number", number)
	}
	return ret, ret2, err
}

func FindRoomChatLogByContentId(logId string) (*types.RoomLogJoinUser, error) {
	//if cfg.CacheType.Enable {
	//	return redis.FindRoomChatLogByContentId(logId)
	//}
	ret, err := mysql.FindRoomChatLogByContentId(logId)
	if err != nil {
		logRoom.Error("mysql.FindRoomChatLogByContentId", "err", err, "logId", logId)
	}
	return ret, err
}

func FindRoomChatLogByMsgId(senderId, msgId string) (*types.RoomLogJoinUser, error) {
	//暂时去掉 ，这边没实现好。会导致一直认为是重复消息
	/*if cfg.CacheType.Enable {
		return redis.FindRoomChatLogByMsgId(senderId, msgId)
	}*/
	ret, err := mysql.FindRoomChatLogByMsgId(senderId, msgId)
	if err != nil {
		logRoom.Error("mysql.FindRoomChatLogByMsgId", "err", err, "senderId", senderId, "msgId", msgId)
	}
	return ret, err
}

func FindRoomChatLogsAfterTime(rooms []string, time int64) ([]*types.RoomLogJoinUser, error) {
	if len(rooms) == 0 {
		ret := make([]*types.RoomLogJoinUser, 0)
		return ret, nil
	}
	ret, err := mysql.FindRoomChatLogsAfterTime(rooms, time)
	if err != nil {
		logRoom.Error("mysql.FindRoomChatLogsAfterTime", "err", err, "rooms", rooms, "time", time)
	}
	return ret, err
}
func FindRoomChatLogsBetweenTime(rooms []string, begin, end int64) ([]*types.RoomLogJoinUser, error) {
	if len(rooms) == 0 {
		ret := make([]*types.RoomLogJoinUser, 0)
		return ret, nil
	}
	ret, err := mysql.FindRoomChatLogsBetweenTime(rooms, begin, end)
	if err != nil {
		logRoom.Error("mysql.FindRoomChatLogsBetweenTime", "err", err, "rooms", rooms, "begin", begin, "end", end)
	}
	return ret, err
}

func FindRoomChatLogsNumberBetween(rooms []string, begin, end int64) (int64, error) {
	ret, err := mysql.FindRoomChatLogsNumberBetween(rooms, begin, end)
	if err != nil {
		logRoom.Error("mysql.FindRoomChatLogsNumberBetween", "err", err, "rooms", rooms, "begin", begin, "end", end)
	}
	return ret, err
}

//return nextLogId , logs ,err
func FindRoomLogs(roomId string, startId, joinTime int64, number int) (int64, []*types.RoomLogJoinUser, error) {
	ret, ret2, err := mysql.FindRoomLogs(roomId, startId, joinTime, number)
	if err != nil {
		logRoom.Error("mysql.FindRoomLogs", "err", err, "roomId", roomId, "startId", startId, "joinTime", joinTime, "number", number)
	}
	return ret, ret2, err
}

func FindRoomLogsByUserId(roomId, owner string, startId, joinTime int64, number int, queryType []string) (int64, []*types.RoomLogJoinUser, error) {
	ret, ret2, err := mysql.FindRoomLogsById(roomId, owner, startId, joinTime, number, queryType)
	if err != nil {
		logRoom.Error("mysql.FindRoomLogsById", "err", err, "roomId", roomId, "startId", startId, "joinTime", joinTime, "number", number, "queryType", queryType)
	}
	return ret, ret2, err
}

func DelRoomChatLogById(logId string) (int, error) {
	//if cfg.CacheType.Enable {
	//	err = redis.DelRoomChatLogById(logId)
	//	if err != nil {
	//		return 0, err
	//	}
	//}
	ret, err := mysql.DelRoomChatLogById(logId)
	if err != nil {
		logRoom.Error("mysql.DelRoomChatLogById", "err", err, "logId", logId)
	}
	return ret, err
}

func AlertRoomRevStateByRevId(revId string, state int) error {
	//if cfg.CacheType.Enable {
	//	err = redis.AlertRoomRevStateByRevId(revId, state)
	//	if err != nil {
	//		return err
	//	}
	//}
	err := mysql.AlertRoomRevStateByRevId(revId, state)
	if err != nil {
		logRoom.Error("mysql.AlertRoomRevStateByRevId", "err", err, "revId", revId, "state", state)
	}
	return err
}

func AppendRoomChatLog(userId, roomId, msgId string, msgType, isSnap int, content, ext string, time int64) (int64, error) {
	_, logId, err := mysql.AppendRoomChatLog(userId, roomId, msgId, msgType, isSnap, content, ext, time)
	//if err == nil {
	//	if cfg.CacheType.Enable {
	//		err = redis.AppendRoomChatLog(logId, userId, roomId, msgId, msgType, isSnap, content, time)
	//		if err != nil {
	//			return 0, err
	//		}
	//	}
	//}
	if err != nil {
		logRoom.Error("mysql.AppendRoomChatLog", "err", err.Error())
	}
	return logId, err
}

//-----------------群成员----------------//
//获取用户创建群个数上限
func GetCreateRoomsLimit(appId string, level int) (int, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetCreateRoomsLimit(appId, level)
		if err != nil {
			logRoom.Error("redis.GetCreateRoomsLimit", "err", err, "appId", appId, "level", level)
		}
		return ret, err
	}
	ret, err := mysql.GetCreateRoomsLimit(appId, level)
	if err != nil {
		logRoom.Error("mysql.GetCreateRoomsLimit", "err", err, "appId", appId, "level", level)
	}
	return ret, err
}

//获取群的成员数上限
func GetRoomMembersLimit(appId string, level int) (int, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetRoomMembersLimit(appId, level)
		if err != nil {
			logRoom.Error("redis.GetRoomMembersLimit", "err", err, "appId", appId, "level", level)
		}
		return ret, err
	}
	ret, err := mysql.GetRoomMembersLimit(appId, level)
	if err != nil {
		logRoom.Error("mysql.GetRoomMembersLimit", "err", err, "appId", appId, "level", level)
	}
	return ret, err
}

//设置用户创建群个数上限
func SetCreateRoomsLimit(appId string, level, limit int) error {
	err := mysql.SetCreateRoomsLimit(appId, level, limit)
	if err != nil {
		logRoom.Error("mysql.SetCreateRoomsLimit", "err", err, "appId", appId, "level", level, "limit", limit)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetCreateRoomsLimit(appId, level, limit)
		if err != nil {
			logRoom.Error("redis.SetCreateRoomsLimit", "err", err, "appId", appId, "level", level, "limit", limit)
			return err
		}
	}
	return nil
}

//设置群的成员数上限
func SetRoomMembersLimit(appId string, level, limit int) error {
	err := mysql.SetRoomMembersLimit(appId, level, limit)
	if err != nil {
		logRoom.Error("mysql.SetRoomMembersLimit", "err", err, "appId", appId, "level", level, "limit", limit)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetRoomMembersLimit(appId, level, limit)
		if err != nil {
			logRoom.Error("redis.SetRoomMembersLimit", "err", err, "appId", appId, "level", level, "limit", limit)
			return err
		}
	}
	return nil
}

//创建群个数 未解散
func FindCreateRoomNumbers(masterId string) (int, error) {
	ret, err := mysql.FindCreateRoomNumbers(masterId)
	if err != nil {
		logRoom.Error("mysql.FindCreateRoomNumbers", "err", err, "masterId", masterId)
	}
	return ret, err
}

//根据appId获取应用的所有未解散的群聊，包括被封群
func GetRoomCountInApp(appId string) (int64, error) {
	ret, err := mysql.GetRoomCountInApp(appId)
	if err != nil {
		logRoom.Error("mysql.GetRoomCountInApp", "err", err, "appId", appId)
	}
	return ret, err
}

//根据appId获取应用的所有被封的群聊，不包括解散的群
func GetCloseRoomCountInApp(appId string) (int64, error) {
	ret, err := mysql.GetCloseRoomCountInApp(appId)
	if err != nil {
		logRoom.Error("mysql.GetCloseRoomCountInApp", "err", err, "appId", appId)
	}
	return ret, err
}

//查询某个app的所有群，模糊查询 markId 群名
func FindRoomsInAppQueryMarkId(appId, query string) ([]*types.RoomJoinUser, error) {
	ret, err := mysql.FindRoomsInAppQueryMarkId(appId, query)
	if err != nil {
		logRoom.Error("mysql.FindRoomsInAppQueryMarkId", "err", err, "appId", appId, "query", query)
	}
	return ret, err
}

//查找某个app所有下未封禁群
func FindRoomsInAppUnClose(appId string) ([]*types.RoomJoinUser, error) {
	ret, err := mysql.FindRoomsInAppUnClose(appId)
	if err != nil {
		logRoom.Error("mysql.FindRoomsInAppUnClose", "err", err, "appId", appId)
	}
	return ret, err
}

//查找某个app所有下封禁群
func FindRoomsInAppClosed(appId string) ([]*types.RoomJoinUser, error) {
	ret, err := mysql.FindRoomsInAppClosed(appId)
	if err != nil {
		logRoom.Error("mysql.FindRoomsInAppClosed", "err", err, "appId", appId)
	}
	return ret, err
}

//查找某个app所有下推荐群（带 群主、人数、封群次数等信息）
func FindRoomsInAppRecommend(appId string) ([]*types.RoomJoinUser, error) {
	ret, err := mysql.FindRoomsInAppRecommend(appId)
	if err != nil {
		logRoom.Error("mysql.FindRoomsInAppRecommend", "err", err, "appId", appId)
	}
	return ret, err
}

//根据群的发言人数递减排序
func RoomsOrderActiveMember(appId string, datetime int64) ([]*types.Room, error) {
	ret, err := mysql.RoomsOrderActiveMember(appId, datetime)
	if err != nil {
		logRoom.Error("mysql.RoomsOrderActiveMember", "err", err, "appId", appId, "datetime", datetime)
	}
	return ret, err
}

//根据群的发言条数递减排序
func RoomsOrderActiveMsg(appId string, datetime int64) ([]*types.Room, error) {
	ret, err := mysql.RoomsOrderActiveMsg(appId, datetime)
	if err != nil {
		logRoom.Error("mysql.RoomsOrderActiveMsg", "err", err, "appId", appId, "datetime", datetime)
	}
	return ret, err
}

//获取所有推荐群
func FindAllRecommendRooms(appId string) ([]*types.Room, error) {
	ret, err := mysql.FindAllRecommendRooms(appId)
	if err != nil {
		logRoom.Error("mysql.FindAllRecommendRooms", "err", err, "appId", appId)
	}
	return ret, err
}

//设置为加v认证群
func SetRoomVerifyed(tx types.Tx, roomId, vInfo string) error {
	_, _, err := db.SetRoomVerifyed(tx.(*sql.MysqlTx), roomId, vInfo)
	if err != nil {
		logRoom.Error("db.SetRoomVerifyed", "err", err, "roomId", roomId, "vInfo", vInfo)
		return err
	}
	if cfg.CacheType.Enable {
		err = redis.SetRoomVerifyed(roomId, vInfo)
		if err != nil {
			logRoom.Error("redis.SetRoomVerifyed", "err", err, "roomId", roomId, "vInfo", vInfo)
			return err
		}
	}
	return nil
}
