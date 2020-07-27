package orm

import (
	redis "github.com/33cn/chat33/orm/redis_model"
	"github.com/33cn/chat33/pkg/otc"
	"github.com/33cn/chat33/types"
	"github.com/inconshreveable/log15"
)

var logOtc = log15.New("module", "model/otc")

func GetOrderById(serverUrl, orderId string) (*types.Order, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetOrderById(serverUrl, orderId)
		if err != nil {
			logFriend.Error("redis.GetOrderById", "err", err, "orderId", orderId)
		}
		return ret, err
	}
	ret, err := otc.GetOrderInfoById(serverUrl, orderId)
	if err != nil {
		logFriend.Error("otc.GetOrderInfoById", "err", err, "orderId", orderId)
	}
	return ret, err
}
