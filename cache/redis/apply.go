package redis

import (
	"encoding/json"

	. "github.com/33cn/chat33/types"
	"github.com/garyburd/redigo/redis"
)

const (
	applyKey   = "apply"
	applyIdKey = "apply-id"
)

type applyCaChe struct{}

var applyType = map[int]string{
	IsRoom:   "room",
	IsFriend: "friend",
}

func (c *applyCaChe) saveLastApplyLogId(a *Apply) error {
	key := applyIdKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	k := a.ApplyUser + "-" + a.Target + "-" + applyType[a.Type]
	_, err := conn.Do("hset", key, k, a.Id)
	if err != nil {
		userLog.Error("HMSET EXPIRE err", err)
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("SaveUserInfo EXPIRE err", err)
	}
	return err
}

//获取请求记录id
func (c *applyCaChe) GetLastApplyLogId(applyUser, target string, tp int) (string, error) {
	key := applyIdKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	k := applyUser + "-" + target + "-" + applyType[tp]
	reply, err := conn.Do("hget", key, k)
	if b, err := IsExists(reply, err); !b {
		return "", err
	}
	id, err := redis.String(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return "", err
	}
	return id, err
}

//储存用户信息
func (c *applyCaChe) SaveApplyLog(a *Apply) error {
	key := applyKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	b, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, a.Id, string(b))
	if err != nil {
		userLog.Error("HMSET EXPIRE err", err)
		return err
	}
	//存储映射关系
	//TODO 应该需要事件
	err = c.saveLastApplyLogId(a)
	//-------------//
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("SaveUserInfo EXPIRE err", err)
	}
	return err
}

//通过userid获取userinfo
func (*applyCaChe) GetApplyLogById(logId string) (*Apply, error) {
	key := applyKey
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
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return nil, err
	}
	var a Apply
	err = json.Unmarshal([]byte(v), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *applyCaChe) GetApplyLogByUserAndTarget(applyUser, target string, tp int) (*Apply, error) {
	id, err := c.GetLastApplyLogId(applyUser, target, tp)
	if err != nil {
		return nil, err
	}
	return c.GetApplyLogById(id)
}

/*func (c *applyCaChe) UpdateApplyInfo(id string, f map[string]interface{}) error{
	key := "apply" + id
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()

	args := make([]interface{},0)
	for k,v := range f{
		args = append(args, k)
		args = append(args, v)
	}

	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	_, err = conn.Do("HMSET", key, args)
	return err
}*/

func (c *applyCaChe) UpdateApplyStateById(logId string, state int) error {
	key := applyKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	apply, err := c.GetApplyLogById(logId)
	if err != nil {
		return err
	}
	if apply == nil {
		friendLog.Error("not exist in redis", "err", err)
		return err
	}
	apply.State = state
	b, err := json.Marshal(&apply)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, logId, string(b))
	if err != nil {
		userLog.Error("HMSET EXPIRE err", err)
		return err
	}
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("SaveUserInfo EXPIRE err", err)
	}
	return err
}
