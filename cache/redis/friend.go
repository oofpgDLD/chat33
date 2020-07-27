package redis

import (
	"encoding/json"

	"github.com/33cn/chat33/utility"

	. "github.com/33cn/chat33/types"
	"github.com/garyburd/redigo/redis"
	"github.com/inconshreveable/log15"
)

const (
	friendsKey         = "friends"
	friendsDelKey      = "del-f"
	friendInfoKey      = "f"
	addFriendConfKey   = "f-conf"
	privateLogKey      = "f-log"
	privateLogIndexKey = "f-log-index"
)

var friendLog = log15.New("cache", "cache friend")

type friendCache struct{}

//获取好友列表
func (f *friendCache) GetFriends(userId string) ([]string, error) {
	key := friendsKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("SMEMBERS", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	return redis.Strings(reply, err)
}

//储存好友列表
func (f *friendCache) SaveFriends(userId string, friendIds []string) error {
	key := friendsKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	for _, v := range friendIds {
		_, err := conn.Do("sadd", key, v)
		if err != nil {
			friendLog.Error("SaveFriendList sadd err", err)
			return err
		}
	}
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	return err
}

//添加好友
func (f *friendCache) AddFriend(userId, friendId string) error {
	key := friendsKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		return err
	}
	//必须判断键是否存在，防止获取好友列表时，拉取的好友数量错误
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		friendLog.Error("AddFriend  EXISTS  err", err)
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("sadd", key, friendId)
	if err != nil {
		friendLog.Error("AddFriend  sadd  err", err)
	}
	return err
}

//获取已删除的好友
func (f *friendCache) GetDelFriends(userId string) ([]string, error) {
	key := friendsDelKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("SMEMBERS", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	return redis.Strings(reply, err)
}

//存储已删除好友列表
func (f *friendCache) SaveDelFriends(userId string, friendIds []string) error {
	key := friendsDelKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	for _, v := range friendIds {
		_, err := conn.Do("sadd", key, v)
		if err != nil {
			friendLog.Error("SaveFriendList sadd err", err)
			return err
		}
	}
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	return err
}

func (f *friendCache) addDelFriend(userId, friendId string) error {
	key := friendsDelKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		return err
	}
	//必须判断键是否存在，防止获取好友列表时，拉取的好友数量错误
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("sadd", key, friendId)
	if err != nil {
		friendLog.Error("savet deleted friend sadd err", err)
	}
	return err
}

//删除好友
func (f *friendCache) DeleteFriend(userId, friendId string) error {
	key := friendsKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		return err
	}
	//必须判断键是否存在，防止获取好友列表时，拉取的好友数量错误
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("srem", key, friendId)
	if err != nil {
		friendLog.Error("DeleteFriend  srem  err", err)
		return err
	}
	return f.addDelFriend(userId, friendId)
}

//是否是好友  第一个bool redis数据是否存在  第二个 是否是好友
func (f *friendCache) IsFriend(userId, friendId string) (*bool, error) {
	key := friendsKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	b, err := redis.Bool(conn.Do("SISMEMBER", key, friendId))
	if err != nil {
		return nil, err
	}
	return &b, nil
}

//储存好友备注(名称)等信息
func (f *friendCache) SaveFriendInfo(userId string, info *Friend) error {
	key := friendInfoKey + userId + "-" + info.FriendId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("HMSET", key,
		"user_id", userId,
		"friend_id", info.FriendId,
		"remark", info.Remark,
		"add_time", info.AddTime,
		"DND", info.DND,
		"top", info.Top,
		"type", info.Type,
		"is_delete", info.IsDelete,
		"source", info.Source,
		"ext_remark", info.ExtRemark,
		"is_blocked", info.IsBlocked,
	)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	return err
}

//获取好友信息
func (f *friendCache) GetFriendInfo(userId, friendId string) (*Friend, error) {
	key := friendInfoKey + userId + "-" + friendId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGETALL", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	v, err := redis.StringMap(reply, err)
	if err != nil {
		return nil, err
	}
	friendInfo := &Friend{
		UserId:    v["user_id"],
		FriendId:  v["friend_id"],
		Remark:    v["remark"],
		AddTime:   utility.ToInt64(v["add_time"]),
		DND:       utility.ToInt(v["DND"]),
		Top:       utility.ToInt(v["top"]),
		Type:      utility.ToInt(v["type"]),
		IsDelete:  utility.ToInt(v["is_delete"]),
		Source:    v["source"],
		ExtRemark: v["ext_remark"],
		IsBlocked: utility.ToInt(v["is_blocked"]),
	}
	return friendInfo, err
}

//修改好友信息
//field  remark  noDisturbing  top
func (f *friendCache) UpdateFriendInfo(userId, friendId, field, value string) error {
	key := friendInfoKey + userId + "-" + friendId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		return err
	}
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("hset", key, field, value)
	if err != nil {
	}
	return err
}

//删除好友信息
func (f *friendCache) DeleteFriendInfo(userId, friendId string) error {
	key := friendInfoKey + userId + "-" + friendId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		friendLog.Error("AddFriend EXPIRE err", err)
		return err
	}

	_, err = conn.Do("del", key)
	return err
}

//存储加好友配置
func (f *friendCache) SaveAddFriendConfig(userId string, conf *AddFriendConf) error {
	key := addFriendConfKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	b, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, userId, string(b))
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	return err
}

//获取加好友配置
func (f *friendCache) GetAddFriendConfig(userId string) (*AddFriendConf, error) {
	key := addFriendConfKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hget", key, userId)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	v, err := redis.String(reply, err)
	if err != nil {
		return nil, err
	}
	var conf AddFriendConf
	err = json.Unmarshal([]byte(v), &conf)
	if err != nil {
		return nil, err
	}
	return &conf, err
}

func (f *friendCache) savePrivateLogIndex(logs []*PrivateLog) error {
	key := privateLogIndexKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	const keyNumb = 1
	args := make([]interface{}, keyNumb+2*len(logs))
	args[keyNumb-1] = key
	for i := 0; i < len(logs); i++ {
		args[keyNumb+2*i] = logs[i].SendTime
		args[keyNumb+(2*i)+1] = logs[i].Id
	}
	//存在时间为永久
	_, err := conn.Do("zadd", args...)
	return err
}

//获取私聊记录
func (f *friendCache) deletePrivateLogIndex(logId string) (int, error) {
	key := privateLogIndexKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	//存在时间为永久
	reply, err := conn.Do("ZREM", key, logId)
	return redis.Int(reply, err)
}

//获取私聊记录
func (f *friendCache) GetPrivateChatLogsIndexByTime(start, end *int64, startEQ, endEQ bool) ([]string, error) {
	key := privateLogIndexKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	const minUnLimit = "-inf"
	const maxUnLimit = "+inf"
	startArg := ""
	endArg := ""
	if start == nil {
		startArg = minUnLimit
	} else {
		if !startEQ {
			startArg = "(" + utility.ToString(*start)
		} else {
			startArg = utility.ToString(*start)
		}
	}

	if end == nil {
		endArg = maxUnLimit
	} else {
		if !endEQ {
			endArg = "(" + utility.ToString(*end)
		} else {
			endArg = utility.ToString(*end)
		}
	}

	reply, err := conn.Do("ZRANGEBYSCORE", key, startArg, endArg)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	return redis.Strings(reply, err)
}

/*//储存私聊记录
func (f *friendCache) SavePrivateChatLog(log *PrivateLog) error {
	key := privateLogKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	b, err := json.Marshal(log)
	if err != nil {
		return err
	}

	_, err = conn.Do("hset", key, log.Id, string(b))
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		return err
	}
	return nil
}*/

//储存私聊记录
func (f *friendCache) SavePrivateChatLogs(logs []*PrivateLog) error {
	key := privateLogKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	const keyNumb = 1
	args := make([]interface{}, keyNumb+2*len(logs))
	args[keyNumb-1] = key
	for i := 0; i < len(logs); i++ {
		b, err := json.Marshal(logs[i])
		if err != nil {
			return err
		}
		args[keyNumb+2*i] = logs[i].Id
		args[keyNumb+(2*i)+1] = string(b)
	}

	_, err := conn.Do("hmset", args...)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		return err
	}
	return f.savePrivateLogIndex(logs)
}

//获取私聊记录
func (f *friendCache) GetPrivateChatLog(logId string) (*PrivateLog, error) {
	key := privateLogKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hget", key, logId)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	v, err := redis.String(reply, err)
	if err != nil {
		return nil, err
	}
	var log PrivateLog
	err = json.Unmarshal([]byte(v), &log)
	if err != nil {
		return nil, err
	}
	return &log, err
}

//删除私聊记录
func (f *friendCache) DeletePrivateChatLog(logId string) (int, error) {
	key := privateLogKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		friendLog.Error("AddFriend EXPIRE err", err)
		return 0, err
	}

	_, err = conn.Do("hdel", key, logId)
	//number, err := redis.Int(reply, err)
	if err != nil {
		return 0, err
	}
	return f.deletePrivateLogIndex(logId)
}
