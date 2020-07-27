package redis

import (
	"encoding/json"

	. "github.com/33cn/chat33/types"
	"github.com/garyburd/redigo/redis"
)

const (
	ordersKey = "orders"
)

type orderCache struct{}

//获取私聊记录
func (o *orderCache) GetOrderById(orderId string) (*Order, error) {
	key := ordersKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hget", key, orderId)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	v, err := redis.String(reply, err)
	if err != nil {
		return nil, err
	}
	var order Order
	err = json.Unmarshal([]byte(v), &order)
	if err != nil {
		return nil, err
	}
	return &order, err
}

//储存好友列表
func (o *orderCache) SaveOrder(order *Order) error {
	key := ordersKey
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
	b, err := json.Marshal(order)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, order.OrderId, string(b))
	if err != nil {
		friendLog.Error("AddFriend  sadd  err", err)
	}
	return err
}
