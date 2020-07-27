package model

import (
	"github.com/33cn/chat33/orm"
)

//查询otc交易关系
func CheckOrder(serverUrl, userUid, targetUid, orderId string) bool {
	//根据订单号 查询交易双方成员
	if orderId == "" {
		return false
	}
	order, err := orm.GetOrderById(serverUrl, orderId)
	if order == nil || err != nil {
		return false
	}

	if (userUid == order.Uid && targetUid == order.Opposite) || (userUid == order.Opposite && targetUid == order.Uid) {
		return true
	}
	return false
}
