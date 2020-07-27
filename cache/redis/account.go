package redis

import (
	"encoding/json"

	user_model "github.com/33cn/chat33/user/model"
	"github.com/garyburd/redigo/redis"
)

type accountCaChe struct{}

const (
	accountTokenKey = "account-token"
)

//储存账户服务器的token和用户的对应关系
func (c *accountCaChe) SaveToken(token *user_model.Token) error {
	key := accountTokenKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeMonth)
	if err != nil {
		userLog.Error("SaveUserInfo EXPIRE err", err)
		return err
	}
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", key, token.AppId+"-"+token.Token, string(b))
	if err != nil {
		userLog.Error("HSET EXPIRE err", err)
		return err
	}
	return err
}

//通过账户体系token获取用户Id
func (*accountCaChe) GetToken(appId, token string) (*user_model.Token, error) {
	key := accountTokenKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	k := appId + "-" + token
	reply, err := conn.Do("HGET", key, k)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return nil, err
	}
	var tokenInfo user_model.Token
	err = json.Unmarshal([]byte(s), &tokenInfo)
	if err != nil {
		return nil, err
	}
	return &tokenInfo, nil
}
