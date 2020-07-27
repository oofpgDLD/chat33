package redis

import (
	"encoding/json"

	"github.com/33cn/chat33/types"
	. "github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/garyburd/redigo/redis"
	"github.com/inconshreveable/log15"
)

//存储群里userid和cid的关系

var roomLog = log15.New("cache", "cache user")

type roomCaChe struct{}

//**********************群所有成员和cid*******************//
//key = "room-all-cid" + roomId

//储存群所有成员和cid
func (*roomCaChe) SaveAllMemberCid(roomId string, rc *[]*RoomCid) error {
	key := "room-all-cid" + roomId
	con := GetConnByKey(key)
	defer func() {
		err := con.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	for _, v := range *rc {
		_, err := con.Do("hset", key, v.UserId, v.GetuiCid)
		if err != nil {
			roomLog.Error("SaveAllMemberCid hset err", err)
			return err
		}
	}
	_, err := con.Do("EXPIRE", key, EXPIRETimeWeek)
	if err != nil {
		roomLog.Error("SaveFriendInfo EXPIRE err", err)
		return err
	}
	return nil
}

//判断用户是否在群  第一个bool  redis数据是否存在  2 是否在群里
func (r *roomCaChe) UserIsInRoom(roomId, userId string) (bool, bool, error) {
	k := "room-all-cid" + roomId
	conn := GetConnByKey(k)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", k, EXPIRETimeWeek)
	if err != nil {
		roomLog.Error("UserIsInRoom EXPIRE err", err)
		return false, false, err
	}
	lock := RedisLock{lockKey: k}
	for {
		err := lock.Lock(&conn, LockTimeout)
		if err == nil {
			break
		}
	}
	count, err := redis.Int(conn.Do("exists", k))
	if err != nil {
		roomLog.Error("RoomAllCidExists exists err", err)
		lock.Unlock(&conn)
		return false, false, err
	}
	if count == 0 {
		lock.Unlock(&conn)
		return false, false, nil
	}
	b, err := redis.Bool(conn.Do("hexists", k, userId))
	if err != nil {
		roomLog.Error("UserIsInRoom hexists err", err)
		lock.Unlock(&conn)
		return false, false, err
	}
	lock.Unlock(&conn)
	return true, b, nil
}

//获取所有群成员的cid  返回key 用户ID value 个推cid
func (r *roomCaChe) GetRoomUserCid(roomId string) (map[string]string, error) {
	key := "room-all-cid" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	ret, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("GetRoomUserCid hgetall err", err)
		return nil, err
	}
	return ret, nil
}

//更新roomcid
func (r *roomCaChe) UpdateRoomCid(roomId, userId, cid string) error {
	key := "room-all-cid" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeWeek)
	if err != nil {
		roomLog.Error("UpdateRoomCid EXPIRE err", err)
		return err
	}
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	count, err := redis.Int(conn.Do("exists", key))
	if err != nil {
		roomLog.Error("UpdateRoomCid exists err", err)
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("hset", key, userId, cid)
	if err != nil {
		roomLog.Error("UpdateRoomCid hset err", err)
		return err
	}
	return nil
}

//根据userid删除user-cid
func (r *roomCaChe) DeleteRoomCidByUserId(roomId, userId string) error {
	key := "room-all-cid" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeWeek)
	if err != nil {
		roomLog.Error("UpdateRoomCid EXPIRE err", err)
		return err
	}
	count, err := redis.Int(conn.Do("exists", key))
	if err != nil {
		roomLog.Error("DeleteRoomCidByUserId exists err", err)
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("hdel", key, userId)
	if err != nil {
		roomLog.Error("DeleteRoomCidByUserId hdel err", err)
		return err
	}
	return nil
}

//删除roomcid
func (r *roomCaChe) DeleteRoomCid(roomId string) error {
	key := "room-all-cid" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("del", key)
	if err != nil {
		roomLog.Error("DeleteRoomCid del err", err)
		return err
	}
	return nil
}

//**********************群信息*******************//
//key = "room-" + roomId

//保存群信息
func (r *roomCaChe) SaveRoomInfo(info *Room) error {
	key := "room-" + info.Id
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("hmset", key,
		"markId", info.MarkId,
		"name", info.Name,
		"avatar", info.Avatar,
		"masterId", info.MasterId,
		"createTime", info.CreateTime,
		"canAddFriend", info.CanAddFriend,
		"joinPermission", info.JoinPermission,
		"recordPermision", info.RecordPermision,
		"adminMuted", info.AdminMuted,
		"masterMuted", info.MasterMuted,
		"encrypt", info.Encrypt,
		"isDelete", info.IsDelete,
		"roomlevel", info.RoomLevel,
		"closeUntil", info.CloseUntil,
		"recommend", info.Recommend,
		"identification", info.Identification,
		"identificationInfo", info.IdentificationInfo,
	)
	if err != nil {
		roomLog.Error("SaveRoomInfo EXPIRE err", err)
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeWeek)
	if err != nil {
		roomLog.Error("SaveFriendInfo EXPIRE err", err)
		return err
	}
	return nil
}

//查询群详情
func (r *roomCaChe) FindRoomInfo(roomId string) (*Room, error) {
	key := "room-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	ret, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("roomInfoExists hgetall err", err)
		return nil, err
	}
	roomInfo := &Room{
		Id:                 roomId,
		MarkId:             ret["markId"],
		Name:               ret["name"],
		Avatar:             ret["avatar"],
		MasterId:           ret["masterId"],
		CreateTime:         utility.ToInt64(ret["createTime"]),
		CanAddFriend:       utility.ToInt(ret["canAddFriend"]),
		JoinPermission:     utility.ToInt(ret["joinPermission"]),
		RecordPermision:    utility.ToInt(ret["recordPermision"]),
		AdminMuted:         utility.ToInt(ret["adminMuted"]),
		MasterMuted:        utility.ToInt(ret["masterMuted"]),
		Encrypt:            utility.ToInt(ret["encrypt"]),
		IsDelete:           utility.ToInt(ret["isDelete"]),
		RoomLevel:          utility.ToInt(ret["roomlevel"]),
		CloseUntil:         utility.ToInt64(ret["closeUntil"]),
		Recommend:          utility.ToInt(ret["recommend"]),
		Identification:     utility.ToInt(ret["identification"]),
		IdentificationInfo: ret["identificationInfo"],
	}
	return roomInfo, err
}

//更新群信息
//field   avatar   masterMuted   name
func (r *roomCaChe) UpdateRoomInfo(roomId, field, value string) error {
	//if field != "avatar" && field != "masterMuted" && field != "name" {
	//	return errors.New("field err")
	//}
	key := "room-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeWeek)
	if err != nil {
		roomLog.Error("UpdateRoomInfo EXPIRE err", err)
		return err
	}
	count, err := redis.Int(conn.Do("exists", key))
	if err != nil {
		roomLog.Error("UpdateRoomInfo exists err", err)
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("hset", key, field, value)
	if err != nil {
		roomLog.Error("UpdateRoomInfo hset err", err)
		return err
	}
	return nil
}

//删除群信息
func (r *roomCaChe) DeleteRoomInfo(roomId string) error {
	key := "room-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("del", key)
	if err != nil {
		roomLog.Error("DeleteRoomInfo del err", err)
		return err
	}
	return nil
}

//**********************群成员信息*******************//
//key = "r" + roomId + "-" + userId

//保存群成员信息
func (r *roomCaChe) SaveRoomUser(roomId string, infos *[]*RoomMember) error {
	for _, info := range *infos {
		key := "r" + roomId + "-" + info.UserId
		conn := GetConnByKey(key)
		defer func() {
			err := conn.Close()
			if err != nil {
				roomLog.Error("Close err", err)
			}
		}()
		_, err := conn.Do("hmset", key,
			"Id", info.Id,
			"roomId", info.RoomId,
			"userId", info.UserId,
			"userNickname", info.UserNickname,
			"level", info.Level,
			"noDisturbing", info.NoDisturbing,
			"commonUse", info.CommonUse,
			"roomTop", info.RoomTop,
			"createTime", info.CreateTime,
			"isDelete", info.IsDelete,
			"source", info.Source,
		)
		if err != nil {
			roomLog.Error("SaveRoomUser hmset err", err)
			err1 := conn.Close()
			if err1 != nil {
				return err1
			}
			return err
		}
		_, err = conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
		if err != nil {
			roomLog.Error("SaveRoomUser EXPIRE err", err)
			err1 := conn.Close()
			if err1 != nil {
				return err1
			}
			return err
		}
		err1 := conn.Close()
		if err1 != nil {
			return err1
		}
	}
	return nil
}

//查询群成员信息
func (r *roomCaChe) FindRoomMemberInfo(roomId, userId string) (*RoomMember, error) {
	key := "r" + roomId + "-" + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	ret, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("FindRoomMemberInfo hgetall err", err)
		return nil, err
	}
	rInfo := &RoomMember{
		Id:           ret["Id"],
		RoomId:       ret["roomId"],
		UserId:       ret["userId"],
		UserNickname: ret["userNickname"],
		Level:        utility.ToInt(ret["level"]),
		NoDisturbing: utility.ToInt(ret["noDisturbing"]),
		CommonUse:    utility.ToInt(ret["commonUse"]),
		RoomTop:      utility.ToInt(ret["roomTop"]),
		CreateTime:   utility.ToInt64(ret["createTime"]),
		IsDelete:     utility.ToInt(ret["isDelete"]),
		Source:       ret["source"],
	}
	return rInfo, err
}

//更新群成员信息
func (r *roomCaChe) UpdateRoomUserInfo(roomId, userId, field, value string) error {
	//if field != "userNickname" && field != "level" && field != "noDisturbing" {
	//	return errors.New("field err")
	//}
	key := "r" + roomId + "-" + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
	if err != nil {
		roomLog.Error("SaveRoomUser EXPIRE err", err)
		return err
	}
	count, err := redis.Int(conn.Do("exists", key))
	if err != nil {
		roomLog.Error("UpdateRoomUserInfo exists err", err)
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("hset", key, field, value)
	if err != nil {
		roomLog.Error("UpdateRoomUserInfo hset err", err)
		return err
	}
	return nil
}

//删除群成员信息
func (r *roomCaChe) DeleteRoomUserInfo(roomId, userId string) error {
	key := "r" + roomId + "-" + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("del", key)
	if err != nil {
		roomLog.Error("DeleteRoomInfo del err", err)
		return err
	}
	return nil
}

//**********************禁言信息*******************//
//key = "r-m" + roomId

//存禁言用户  1和3存类型  2存时间
func (r *roomCaChe) SaveRoomUserMued(roomId string, muens *[]*RoomUserMuted) error {
	key := "r-m" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	for _, v := range *muens {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}

		_, err = conn.Do("hset", key, v.UserId, string(b))
		if err != nil {
			roomLog.Error("SaveRoomUserMued hset err", err)
			return err
		}

	}
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		roomLog.Error("SaveRoomUserMued EXPIRE err", err)
		return err
	}
	return nil
}

//获取群中禁言详情
func (r *roomCaChe) GetRoomUserMuted(roomId, userId string) (*RoomUserMuted, error) {
	//conn := GetConn()
	key := "r-m" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	reply, err := conn.Do("hget", key, userId)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		roomLog.Error("GetRoomUserMuted hget err", err)
		return nil, err
	}
	var m RoomUserMuted
	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	return &m, err
}

//更新/添加群禁言信息
func (r *roomCaChe) UpdateRoomUserMuted(roomId string, rum *RoomUserMuted) error {
	key := "r-m" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		roomLog.Error("SaveRoomUserMued EXPIRE err", err)
		return err
	}
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	count, err := redis.Int(conn.Do("exists", key))
	if err != nil {
		roomLog.Error("UpdateRoomUserMuted exists err", err)
		return err
	}
	if count == 0 {
		return nil
	}
	//先取出原来的
	reply, err := conn.Do("hget", key, rum.UserId)
	if b, err := IsExists(reply, err); !b {
		return err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		roomLog.Error("GetRoomUserMuted hget err", err)
		return err
	}
	m := &RoomUserMuted{}
	err = json.Unmarshal([]byte(s), m)
	if err != nil {
		return err
	}
	//存成新的
	m.Deadline = rum.Deadline
	m.ListType = rum.ListType

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, rum.UserId, string(b))

	if err != nil {
		roomLog.Error("UpdateRoomUserMuted hset err", err)
		return err
	}
	return nil
}

//删除群禁言信息
func (r *roomCaChe) DeleteRoomUserMuted(roomId string) error {
	key := "r-m" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	_, err := conn.Do("del", key)
	if err != nil {
		roomLog.Error("DeleteRoomInfo del err", err)
		return err
	}
	return nil
}

//**********************群聊聊天相关*******************//
//key = "rmc" + roomId + "-" + userId

//存储群消息记录
func (r *roomCaChe) SaveRoomMsgContent(log *RoomLog) error {
	key := "Room-log"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
	if err != nil {
		roomLog.Error("SaveRoomMsgContent EXPIRE err", err)
	}
	b, err := json.Marshal(log)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, log.Id, string(b))
	return err
}

//根据logId获取消息记录
func (r *roomCaChe) GetRoomMsgContent(id string) (*RoomLog, error) {
	key := "Room-log"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGET", key, id)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		roomLog.Error("GetRoomMsgContent  Do hgetall err", err)
		return nil, err
	}
	var l RoomLog
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

//获取所有消息记录
func (r *roomCaChe) GetRoomMsgContentAll() ([]*RoomLog, error) {
	key := "Room-log"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGETALL", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}

	s, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("GetRoomMsgContentAll  Do hgetall err", err)
		return nil, err
	}
	res := make([]*RoomLog, 0)

	for _, v := range s {
		var log RoomLog
		err = json.Unmarshal([]byte(v), &log)
		if err != nil {
			return nil, err
		}
		res = append(res, &log)
	}
	return res, nil
}

//根据id删除消息记录
func (r *roomCaChe) DeleteRoomMsgContent(logId string) error {
	key := "Room-log"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	_, err := conn.Do("hdel", key, logId)
	if err != nil {
		roomLog.Error("DeleteRoomMsgContent del err", err)
		return err
	}
	return nil
}

//存储群recv消息记录
func (r *roomCaChe) SaveReceiveLog(log *RoomMsgReceive) error {
	key := "Room-ReceiveLog"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
	if err != nil {
		roomLog.Error("SaveReceiveLog EXPIRE err", err)
	}
	b, err := json.Marshal(log)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, log.Id, string(b))
	return err
}

//根据Id获取消息记录
func (r *roomCaChe) GetReceiveLogbyId(id string) (*RoomMsgReceive, error) {
	key := "Room-ReceiveLog"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGET", key, id)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		roomLog.Error("GetReceiveLogbyId  Do hget err", err)
		return nil, err
	}
	var l RoomMsgReceive
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

//获取所有recv消息记录
func (r *roomCaChe) GetReceiveLogAll() ([]*RoomMsgReceive, error) {
	key := "Room-ReceiveLog"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGETALL", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}

	s, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("GetReceiveLogAll  Do hgetall err", err)
		return nil, err
	}
	res := make([]*RoomMsgReceive, 0)

	for _, v := range s {
		var log RoomMsgReceive
		err = json.Unmarshal([]byte(v), &log)
		if err != nil {
			return nil, err
		}
		res = append(res, &log)
	}
	return res, nil
}

//更新recv消息状态
func (r *roomCaChe) UpdateReceiveLog(id string, state int) error {
	key := "Room-ReceiveLog"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGET", key, id)
	if b, err := IsExists(reply, err); !b {
		return err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		roomLog.Error("ReceiveLog  Do HGET err", err)
		return err
	}

	log := &RoomMsgReceive{}
	//先解码修改
	err = json.Unmarshal([]byte(s), log)
	if err != nil {
		return err
	}
	log.State = state
	//重新marshal往里存
	b, err := json.Marshal(log)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, id, string(b))
	return err
}

//根据id删除recv消息记录
func (r *roomCaChe) DeleteReceiveLog(logId string) error {
	key := "Room-ReceiveLog"
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	_, err := conn.Do("hdel", key, logId)
	if err != nil {
		roomLog.Error("ReceiveLog del err", err)
		return err
	}
	return nil
}

//**********************群和群成员的关系*******************//
//群和成员的关系
func (r *roomCaChe) AddRoomUser(roomId string, muens *[]*types.RoomMember) error {
	key := "room_user" + "-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	//存成key  userid   level 以后可以取出作为map[userid]{level}
	for _, v := range *muens {
		_, err := conn.Do("hset", key,
			v.UserId, v.Level)
		if err != nil {
			roomLog.Error("AddRoomUser hset err", err)
			return err
		}
	}
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		roomLog.Error("AddRoomUser EXPIRE err", err)
		return err
	}
	return nil
}

//获取群和成员关系 map[userid]{level}
func (r *roomCaChe) GetRoomUser(roomId string) (map[string]int64, error) {
	key := "room_user" + "-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	ret, err := redis.Int64Map(reply, err)

	if err != nil {
		roomLog.Error("GetRoomUser hgetall err", err)
		return nil, err
	}
	return ret, err
}

//更新群和成员关系
func (r *roomCaChe) UpdateRoomUser(roomId, userId, value string) error {
	key := "room_user" + "-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		roomLog.Error("UpdateRoomUser EXPIRE err", err)
		return err
	}
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	count, err := redis.Int(conn.Do("exists", key))
	if err != nil {
		roomLog.Error("UpdateRoomUser exists err", err)
		return err
	}
	if count == 0 {
		return nil
	}

	_, err = conn.Do("hset", key, userId, value)

	if err != nil {
		roomLog.Error("UpdateRoomUserMuted hset err", err)
		return err
	}
	return nil
}

//删除群和成员关系(群没被删除，只是成员被删除)
func (r *roomCaChe) DelRoomUser(roomId, userId string) error {
	key := "room_user" + "-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	//删除指定userid，也就是这个群没有这个人了
	_, err := conn.Do("hdel", key, userId)
	if err != nil {
		roomLog.Error("DeleteRoomInfo del err", err)
		return err
	}
	return nil
}

//删除群和成员关系(群被删除的情况)
func (r *roomCaChe) DelRoomUserAll(roomId string) error {
	key := "room_user" + "-" + roomId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	//lock := RedisLock{lockKey: key}
	//defer func() {
	//	lock.Unlock(&conn)
	//}()
	//for {
	//	err := lock.Lock(&conn, LockTimeout)
	//	if err == nil {
	//		break
	//	}
	//}
	//删除这个群，关系也不存在了
	_, err := conn.Do("del", key)
	if err != nil {
		roomLog.Error("DeleteRoom del err", err)
		return err
	}
	return nil
}

//**********************room_config 和 user_config*******************//
//新增或更新room_config
func (r *roomCaChe) SaveRoomConfig(appId string, level, limit int) error {
	key := "Room_Config" + "-" + appId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
	if err != nil {
		roomLog.Error("SaveRoomConfig EXPIRE err", err)
	}

	_, err = conn.Do("hmset", key,
		"level", level,
		"limit", limit,
	)
	return err
}

//查找room_config
func (r *roomCaChe) FindRoomConfig(appId string) (*RoomConfig, error) {
	key := "Room_Config" + "-" + appId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	ret, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("FindRoomConfig hget err", err)
		return nil, err
	}
	return &RoomConfig{
		AppId: appId,
		Level: utility.ToInt(ret["level"]),
		Limit: utility.ToInt(ret["limit"]),
	}, err
}

//新增或更新user_config
func (r *roomCaChe) SaveUserConfig(appId string, level, limit int) error {
	key := "User_Config" + "-" + appId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeHalfDay)
	if err != nil {
		roomLog.Error("SaveUserConfig EXPIRE err", err)
	}
	_, err = conn.Do("hmset", key,
		"level", level,
		"limit", limit,
	)
	return err
}

//获取user_config
func (r *roomCaChe) FindUserConfig(appId string) (*UserConfig, error) {
	key := "User_Config" + "-" + appId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			roomLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	ret, err := redis.StringMap(reply, err)
	if err != nil {
		roomLog.Error("FindRoomConfig hget err", err)
		return nil, err
	}
	return &UserConfig{
		AppId: appId,
		Level: utility.ToInt(ret["level"]),
		Limit: utility.ToInt(ret["limit"]),
	}, err
}
