package redis

import (
	"encoding/json"
	"errors"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/garyburd/redigo/redis"
)

type praiseCaChe struct{}

const (
	praiseTokenKey = "praise"
)

//储存账户服务器的token和用户的对应关系
func (c *praiseCaChe) SaveLeaderBoard(tp int, records map[string]*types.RankingItem, startTime, endTime int64) error {
	typeStr := ""
	switch tp {
	case types.Like:
		typeStr = "Like"
	case types.Reward:
		typeStr = "Reward"
	default:
		return errors.New("err praise type")
	}
	key := praiseTokenKey + "-" + typeStr + "-" + utility.ToString(startTime) + "-" + utility.ToString(endTime)
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()

	for _, v := range records {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = conn.Do("hset", key, v.UserId, string(b))
		if err != nil {
			roomLog.Error("SaveAllMemberCid hset err", err)
			return err
		}
	}
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	return err
}

//通过账户体系token获取用户Id
func (*praiseCaChe) GetPraiseStatic(tp int, startTime, endTime int64) (map[string]*types.RankingItem, error) {
	typeStr := ""
	switch tp {
	case types.Like:
		typeStr = "Like"
	case types.Reward:
		typeStr = "Reward"
	default:
		return nil, errors.New("err praise type")
	}
	key := praiseTokenKey + "-" + typeStr + "-" + utility.ToString(startTime) + "-" + utility.ToString(endTime)
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	list, err := redis.StringMap(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return nil, err
	}
	ret := make(map[string]*types.RankingItem)
	for k, v := range list {
		var item types.RankingItem
		err := json.Unmarshal([]byte(v), &item)
		if err != nil {
			continue
		}
		ret[k] = &item
	}
	return ret, err
}
