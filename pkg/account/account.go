package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/33cn/chat33/pkg/http"
	middle_ware "github.com/33cn/chat33/pkg/middle-ware"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var log = log15.New("module", "library/bcoin")

// 从币钱包获取用户信息
func UserInfoFromBcoin(serverUrl, token string) (*types.User, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	byte, err := http.HTTPPostForm(serverUrl+"/v1/user/info", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /v1/user/info failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
				{
		    		"code": 200,
		   		 	"message": "请求成功",
		    		"data": {
		        		"id": 34,
		        		"mobile": ,
		        		"avatar": "",
		        		"is_setpaypasswd": 1,
		        		"is_real": 0
		    		}
				}

				{
		    		"code": 401,
		    		"message": "token无效或过期",
		    		"data": null,
		    		"info": {
		        		"name": "Unauthorized",
		        		"message": "token无效或过期",
		        		"code": 0,
		        		"status": 401,
		        		"type": "yii\\web\\UnauthorizedHttpException"
		    		}
				}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, errors.New("请求账户信息超时")
	}

	if utility.ToInt(resp["code"]) != 200 {
		if utility.ToInt(resp["code"]) == 401 {
			return nil, types.ERR_LOGIN_EXPIRED
		}
		return nil, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return nil, errors.New("no 'data' info")
	}

	info := data.(map[string]interface{})
	mobile := utility.ToString(info["mobile"])
	area := utility.ToString(info["area"])
	email := utility.ToString(info["email"])
	isSetPayPwd := utility.ToInt(info["is_setpaypasswd"])
	//账户号
	var account string
	if mobile != "" {
		if area == "" {
			area = "86"
		}
		account = area + mobile
	} else {
		account = email
	}

	return &types.User{
		Uid:         utility.ToString(info["id"]),
		Account:     account,
		Email:       email,
		Area:        area,
		Phone:       mobile,
		Username:    utility.ToString(info["username"]),
		Avatar:      utility.ToString(info["avatar"]),
		IsSetPayPwd: isSetPayPwd,
		Verified:    utility.ToInt(info["is_real"]),
		InviteCode:  utility.ToString(info["invite_code"]),
	}, nil
}

//币钱包 获取用户公钥
func PublicKeyFromBcoin(serverUrl, token string) (string, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	byte, err := http.HTTPPostForm(serverUrl+"/v1/user/public-key", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:08
		log.Warn(fmt.Sprintf("http client request %v /v1/user/public-key failed", serverUrl), "err", err.Error())
		return "", errors.New("请求账户信息超时")
	}

	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return "", errors.New("请求账户信息超时")
	}

	if utility.ToInt(resp["code"]) != 200 {
		return "", errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return "", errors.New("no 'data' info")
	}

	info := data.(map[string]interface{})
	publicKey := utility.ToString(info["public_key"])
	return publicKey, nil
}

func BcoinSetPayPwd(serverUrl, token, mode, payPassword, authType, code, oldPayPassword string) error {
	params := url.Values{}
	params.Set("mode", mode)
	params.Set("type", authType)
	params.Set("code", code)
	params.Set("old_pay_password", oldPayPassword)
	params.Set("pay_password", payPassword)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	headers["Fzm-Request-Source"] = "chat"
	byte, err := http.HTTPPostForm(serverUrl+"/v1/user/set-pay-password", headers, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:28:44
		log.Warn(fmt.Sprintf("http client request %v /v1/user/set-pay-password failed", serverUrl), "err", err.Error())
		return errors.New("请求账户信息超时")
	}

	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return errors.New(resp["message"].(string))
	}

	_, ok := resp["data"]
	if !ok {
		return errors.New("no 'data' info")
	}

	//info := data.(map[string]interface{})
	//uid := utility.ToString(info["uid"])

	return nil
}

func BcoinCheckPayPwd(serverUrl, token, payPassword string) error {
	params := url.Values{}
	params.Set("pay_password", payPassword)

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	headers["Fzm-Request-Source"] = "chat"
	byte, err := http.HTTPPostForm(serverUrl+"/v1/user/validate-pay-password", headers, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:29:21
		log.Warn(fmt.Sprintf("http client request %v /v1/user/validate-pay-password failed", serverUrl), "err", err.Error())
		return errors.New("请求账户信息超时")
	}

	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return errors.New(resp["message"].(string))
	}

	_, ok := resp["data"]
	if !ok {
		return errors.New("no 'data' info")
	}

	//info := data.(map[string]interface{})
	//uid := utility.ToString(info["uid"])

	return nil
}

//代币划转
func BcoinWithdraw(serverUrl, token, currency string, amount, fee float64, opp_address, rid, mode, payword, code string, isPublickey bool) (string, string, error) {
	params := url.Values{}
	params.Set("currency", currency)
	params.Set("amount", utility.ToString(amount))
	params.Set("fee", utility.ToString(fee))
	params.Set("opp_address", opp_address)
	params.Set("rid", rid)
	params.Set("mode", mode)
	params.Set("payword", payword)
	params.Set("code", code)
	if isPublickey {
		params.Set("is_publickey", "1")
	} else {
		params.Set("is_publickey", "0")
	}

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	headers["Fzm-Request-Source"] = "chat"
	byte, err := http.HTTPPostForm(serverUrl+"/v1/coin/withdraw", headers, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:30:49
		log.Warn(fmt.Sprintf("http client request %v /v1/coin/withdraw failed", serverUrl), "err", err.Error())
		return "", "", errors.New("请求账户信息超时")
	}

	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return "", "", err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return "", "", errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return "", "", errors.New("no 'data' info")
	}

	info := data.(map[string]interface{})
	recordTo := utility.ToString(info["record_to_id"])
	recordFrom := utility.ToString(info["record_from_id"])
	if recordTo == "" || recordFrom == "" {
		log.Warn("call account system can not get record id", "recordTo", recordTo, "recordFrom", recordFrom)
	}
	return recordFrom, recordTo, nil
}

//从钱包获取 推荐统计信息
func InviteStatisticsFromBcoin(serverUrl, token string) (interface{}, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	byte, err := http.HTTPRequest("GET", serverUrl+"/v1/user/invite-info", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:29:21
		log.Warn(fmt.Sprintf("http client request %v /v1/user/invite-info failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}

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
	return data, nil
}

//从钱包获取 推荐信息
func InviteInfosFromBcoin(serverUrl, token string, rewardType, page, size int64) (interface{}, error) {
	params := url.Values{}
	params.Set("reward_type", utility.ToString(rewardType))
	params.Set("page", utility.ToString(page))
	params.Set("size", utility.ToString(size))

	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	byte, err := http.HTTPRequest("GET", serverUrl+"/v1/user/invite-reward-list?"+params.Encode(), headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:29:21
		log.Warn(fmt.Sprintf("http client request %v /v1/user/invite-reward-list failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}

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
	return data, nil
}

//------------------------------------币钱包后台-------------------------------------//
//后台鉴权
func backendAuthHeader(appKey, appSecret string, params map[string]string) map[string]string {
	time := utility.ToString(utility.NowSecond())
	//time := "1562930763"
	signature := middle_ware.BcoinBackend(appKey, appSecret, time, params)

	headers := make(map[string]string)
	headers["FZM-Wallet-Check"] = "true"
	//分配的key
	headers["FZM-Wallet-AppKey"] = appKey
	//当前时间戳
	headers["FZM-Wallet-Timestamp"] = time
	//签名
	headers["FZM-Wallet-Signature"] = signature

	return headers
}

// 从币钱包 设置奖励金额
func EditRealBcoin(appKey, appSecret, serverUrl string, params *EditRewardParam) error {
	pMap := middle_ware.GetParamsMap(params)
	headers := backendAuthHeader(appKey, appSecret, pMap)

	body, err := json.Marshal(params)
	if err != nil {
		log.Error("json can not Marshal", "err", err)
		return errors.New("参数错误")
	}
	byte, err := http.HTTPPostJSON(serverUrl+"/backend/config/edit-real", headers, bytes.NewBuffer(body))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/config/edit-real failed", serverUrl), "err", err.Error())
		return errors.New("请求账户信息超时")
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
		return err
	}

	if utility.ToInt(resp["code"]) != 200 {
		log.Error("bcoin system return err", "req body", string(body))
		return errors.New(resp["message"].(string))
	}

	_, ok := resp["data"]
	if !ok {
		return errors.New("接口返回解析错误: `data` not find")
	}

	return nil
}

// 从币钱包 查询奖励金额
func ShowRealBcoin(appKey, appSecret, serverUrl string) (interface{}, error) {
	headers := backendAuthHeader(appKey, appSecret, nil)
	byte, err := http.HTTPRequest("GET", serverUrl+"/backend/config/real-config", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/config/real-config failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
		{
			"code":200,
			"message":"请求成功",
			"data":{
				"base":{
					"open":"1",
					"level":"3",
					"currency":"YCC",
					"rewardForUser":"2",
					"rewardForInviter":[
						"3",
						"2",
						"0"
					]
				},
				"advance":{
					"open":"1",
					"currency":"BTY",
					"reachNum":"2",
					"rewardForNum":"3"
				}
			}
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

	return data, nil
}

// 从币钱包 获取支持币种 返回：[]string
func CoinSupport(appKey, appSecret, serverUrl string) (interface{}, error) {
	headers := backendAuthHeader(appKey, appSecret, nil)
	byte, err := http.HTTPRequest("GET", serverUrl+"/backend/coin/coin-list", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/coin/coin-list failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
		{
			"code":200,
			"message":"请求成功",
			"data":[
				"BTY",
				"ETH"
			]
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

	return data, nil
}

// 从币钱包 获取支持币种 返回：[]string
func RewardList(appKey, appSecret, serverUrl string, query *RewardListParam) (interface{}, error) {
	pMap := make(map[string]string)
	params := url.Values{}
	if query != nil {
		if query.Search != "" {
			params.Set("search", query.Search)
			pMap["search"] = query.Search
		}
		if query.RewardType != nil {
			params.Set("rewardType", utility.ToString(*query.RewardType))
			pMap["rewardType"] = utility.ToString(*query.RewardType)

		}
		params.Set("page", utility.ToString(query.Page))
		params.Set("size", utility.ToString(query.Size))
		pMap["page"] = utility.ToString(query.Page)
		pMap["size"] = utility.ToString(query.Size)
	}
	headers := backendAuthHeader(appKey, appSecret, pMap)
	byte, err := http.HTTPRequest("GET", serverUrl+"/backend/coin/real-reward-list?"+params.Encode(), headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/coin/real-reward-list failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
		{
			"code":200,
			"message":"请求成功",
			"data":[
				"BTY",
				"ETH"
			]
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

	return data, nil
}

//红包手续费配置信息
func RPFeeConfigFromBcoin(appKey, appSecret, serverUrl string) (interface{}, error) {
	headers := backendAuthHeader(appKey, appSecret, nil)
	byte, err := http.HTTPRequest("GET", serverUrl+"/backend/config/fee-config", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/config/fee-config failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
		{
			"code": 200,
			"message": "请求成功",
			"data": [
				{
					"amount": "2",
					"currency": "TSC"
				},
				{
					"amount": "1",
					"currency": "HT"
				}
			]
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

	return data, nil
}

//设置红包手续费
func SetRPFeeConfigFromBcoin(appKey, appSecret, serverUrl string, params *SetRPFeeParam) error {
	pMap := middle_ware.GetParamsMap(params)
	headers := backendAuthHeader(appKey, appSecret, pMap)
	body, err := json.Marshal(params)
	if err != nil {
		return errors.New("参数解失败")
	}
	byte, err := http.HTTPPostJSON(serverUrl+"/backend/config/edit-fee", headers, bytes.NewBuffer(body))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/config/edit-fee failed", serverUrl), "err", err.Error())
		return errors.New("请求账户信息超时")
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
		return err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return errors.New(resp["message"].(string))
	}

	_, ok := resp["data"]
	if !ok {
		return errors.New("no 'data' info")
	}

	return nil
}

//红包手续费统计信息
func RPFeeStatistics(appKey, appSecret, serverUrl string) (interface{}, error) {
	headers := backendAuthHeader(appKey, appSecret, nil)
	byte, err := http.HTTPRequest("GET", serverUrl+"/backend/config/fee-profit", headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/config/fee-profit failed", serverUrl), "err", err.Error())
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

	return data, nil
}

//红包手续费配置信息 公共数据接口
func RPFeeInfoFromBcoin(serverUrl string) (interface{}, error) {
	//headers := make(map[string]string)
	//headers["Authorization"] = "Bearer " + token
	byte, err := http.HttpGet(serverUrl + "/v1/data/red-packet-fee")
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /v1/data/red-packet-fee failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
		{
			"code": 200,
			"message": "请求成功",
			"data": [
				{
					"amount": "2",
					"currency": "TSC"
				},
				{
					"amount": "1",
					"currency": "HT"
				}
			]
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

	return data, nil
}

//币钱包 转账
func BcoinTakeCost(serverUrl, token, currency, amount, payPassword, code, uniqueId string, costCategory string) (string, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	headers["Fzm-Request-Source"] = "chat"

	params := url.Values{}
	params.Set("currency", currency)
	params.Set("amount", amount)
	params.Set("pay_password", payPassword)
	if code != "" {
		params.Set("code", code)
	}
	params.Set("cost_category", costCategory)
	params.Set("unique_id", uniqueId)

	byte, err := http.HTTPPostForm(serverUrl+"/v1/coin/take-cost", headers, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /v1/coin/take-cost failed", serverUrl), "err", err.Error())
		return "", errors.New("请求账户信息超时")
	}
	/*
		{
			"code": 200,
			"message": "请求成功",
			"data": {
				 "unique_id":"NQrqWISRghmMB9DkVHRjPNFnBYWvZIJh"
			}
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return "", err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return "", errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return "", errors.New("no 'data' info")
	}

	recordId := ""
	if d, ok := data.(map[string]interface{}); !ok {
		return "", errors.New("account result err")
	} else {
		recordId = utility.ToString(d["unique_id"])
	}

	return recordId, nil
}

//查询转账记录
func BcoinCheckTrans(serverUrl, uniqueId string, costType int) (int, error) {
	costCategory := ""
	switch costType {
	case types.VerifyRecordCost:
		costCategory = "vip_auth"
	case types.VerifyRecordRefund:
		costCategory = "vip_auth_return"
	default:
		return 0, errors.New("请求参数错误")
	}

	params := url.Values{}
	params.Set("unique_id", uniqueId)
	params.Set("cost_category", costCategory)

	byte, err := http.HTTPPostForm(serverUrl+"/backend/coin/check-transaction", nil, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/coin/check-transaction", serverUrl), "err", err.Error())
		return 0, errors.New("请求账户信息超时")
	}
	/*
		{
			"code":200,
			"message":"操作成功",
			"data":{
				"status":-1
			}
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return 0, err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return 0, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return 0, errors.New("no 'data' info")
	}

	ret := 0
	if d, ok := data.(map[string]interface{}); !ok {
		return 0, errors.New("account result err")
	} else {
		if costType == types.VerifyRecordCost {
			switch utility.ToInt(d["status"]) {
			case 1:
				ret = types.VerifyFeeStateSuccess
			case 3:
				ret = types.VerifyFeeStateCosting
			default:
				ret = types.VerifyFeeStateFailed
			}
		} else if costType == types.VerifyRecordRefund {
			switch utility.ToInt(d["status"]) {
			case 1:
				ret = types.VerifyFeeStateBackSuccess
			case 3:
				ret = types.VerifyFeeStateBacking
			default:
				ret = types.VerifyFeeStateBackFailed
			}
		}
	}

	return ret, nil
}

//退回转账金额
func BcoinReturnCost(serverUrl, uniqueId string, costType int) (string, error) {
	costCategory := ""
	switch costType {
	case types.VerifyRecordCost:
		costCategory = "vip_auth"
	case types.VerifyRecordRefund:
		costCategory = "vip_auth_return"
	default:
		panic("错误的参数类型")
	}

	headers := make(map[string]string)
	headers["FZM-Wallet-Source"] = "chat"

	params := url.Values{}
	params.Set("unique_id", uniqueId)
	params.Set("cost_category", costCategory)

	byte, err := http.HTTPPostForm(serverUrl+"/backend/coin/return-cost", headers, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /backend/coin/return-cost", serverUrl), "err", err.Error())
		return "", errors.New("请求账户信息超时")
	}
	/*
		{
			"code": 200,
			"message": "请求成功",
			"data": {
				 "unique_id":"NQrqWISRghmMB9DkVHRjPNFnBYWvZIJh"
			}
		}
	*/
	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return "", err
	}

	if utility.ToInt(resp["code"]) != 200 {
		return "", errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return "", errors.New("no 'data' info")
	}

	recordId := ""
	if d, ok := data.(map[string]interface{}); !ok {
		return "", errors.New("account result err")
	} else {
		recordId = utility.ToString(d["unique_id"])
	}

	return recordId, nil
}

//打赏转账
//币钱包 转账
func BcoinPraise(serverUrl, token, toUid, currency, amount, payPassword, uniqueId string, costCategory string) (interface{}, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + token
	headers["FZM-Wallet-Source"] = "chat"
	headers["Fzm-Request-Source"] = "chat"

	params := url.Values{}
	params.Set("currency", currency)
	params.Set("amount", amount)
	params.Set("pay_password", payPassword)
	params.Set("to_uid", toUid)
	params.Set("cost_category", costCategory)
	params.Set("unique_id", uniqueId)

	byte, err := http.HTTPPostForm(serverUrl+"/v1/coin/user-transfer", headers, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		log.Warn(fmt.Sprintf("http client request %v /v1/coin/user-transfer failed", serverUrl), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	/*
				{
					"code": 200,
					"message": "获取成功",
					"data": {
						 "unique_id":"SyEn2W4rqbVuZKtuzryKqiMyOEfjwEkd",
		                 "currency":"BTY",
		                 "amount":10
					}
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

	return data, nil
}
