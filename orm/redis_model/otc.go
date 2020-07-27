package redis_model

import (
	"github.com/33cn/chat33/cache"
	"github.com/33cn/chat33/pkg/otc"
	"github.com/33cn/chat33/types"
)

//交易双方信息
func GetOrderById(serverUrl, orderId string) (*types.Order, error) {
	order, err := cache.Cache.GetOrderById(orderId)
	if err != nil {
		l.Warn("find user info from redis", "err", err)
	}
	if order == nil {
		//从接口获取并更新缓存
		order, err := otc.GetOrderInfoById(serverUrl, orderId)
		if order == nil {
			l.Warn("order info can not find")
			return nil, err
		}
		err = cache.Cache.SaveOrder(order)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return order, nil
	}
	return order, nil
}
