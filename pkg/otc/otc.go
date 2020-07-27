package otc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/33cn/chat33/pkg/http"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var log = log15.New("otc server")

//检查交易对双方的
func GetOrderInfoById(serverUrl, orderId string) (*types.Order, error) {
	params := url.Values{}
	params.Add("order_num", orderId)

	byte, err := http.HTTPPostForm(serverUrl+"/receive/trade-opponent", nil, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/config/edit-fee failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
		{
			"code": 200,
			"message": "请求成功",
			"data": null
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return nil, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return nil, errors.New("no 'data' info")
	}

	var order *types.Order
	if d, ok := data.(map[string]interface{}); ok {
		user_id := d["user_id"].(string)
		merchant_id := d["merchant_id"].(string)

		order = &types.Order{
			OrderId:  orderId,
			Uid:      user_id,
			Opposite: merchant_id,
		}
	}
	return order, nil
}
