package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/33cn/chat33/pkg/http"
	. "github.com/33cn/chat33/pkg/work/model"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var defaultHost = "http://127.0.0.1:8088"

var l = log15.New("library", "work")

func Init(cfg *Config) {
	if cfg.Host != "" {
		defaultHost = cfg.Host
	}
	l.Info("init work service", "host", defaultHost)
}

func GetUser(appId, uid string) (*UserInfo, error) {
	params := url.Values{}
	params.Set("app_id", appId)
	params.Set("uid", uid)

	byte, err := http.HTTPPostForm(defaultHost+"/work/checkUser", nil, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		l.Warn(fmt.Sprintf("http client request %v /work/checkUser failed", defaultHost), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
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
		return nil, err
	}

	if utility.ToInt(resp["result"]) != 0 {
		return nil, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return nil, errors.New("no 'data' info")
	}

	var ret *UserInfo
	if d, ok := data.(map[string]interface{}); !ok {
		return nil, errors.New("account result err")
	} else {
		ret = &UserInfo{
			AppId:          utility.ToString(d["appId"]),
			Uid:            utility.ToString(d["uid"]),
			Name:           utility.ToString(d["name"]),
			Code:           utility.ToString(d["code"]),
			EnterpriseCode: utility.ToString(d["enterpriseCode"]),
			EnterpriseName: utility.ToString(d["enterpriseName"]),
		}
	}
	return ret, nil
}

func GetUsers(appId, uid string) (*UsersList, error) {
	params := url.Values{}
	params.Set("app_id", appId)
	params.Set("uid", uid)

	byte, err := http.HTTPPostForm(defaultHost+"/work/users", nil, strings.NewReader(params.Encode()))
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:22:53
		l.Warn(fmt.Sprintf("http client request %v /work/users failed", defaultHost), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
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
		return nil, err
	}

	if utility.ToInt(resp["result"]) != 0 {
		return nil, errors.New(resp["message"].(string))
	}

	data, ok := resp["data"]
	if !ok {
		return nil, errors.New("no 'data' info")
	}

	var ret UsersList
	if d, ok := data.(map[string]interface{}); !ok {
		return nil, errors.New("account result err")
	} else {
		if enterprise, ok := d["enterprise"].(map[string]interface{}); ok {
			ret.Enterprise.Code = utility.ToString(enterprise["code"])
			ret.Enterprise.Name = utility.ToString(enterprise["name"])
			ret.Enterprise.Remark = utility.ToString(enterprise["remark"])
		}

		if users, ok := d["users"].([]interface{}); ok {
			us := make([]*UserInfo, 0)
			for _, user := range users {
				if u, ok := user.(map[string]interface{}); ok {
					us = append(us, &UserInfo{
						AppId:          utility.ToString(u["appId"]),
						Uid:            utility.ToString(u["uid"]),
						Name:           utility.ToString(u["name"]),
						Code:           utility.ToString(u["code"]),
						EnterpriseCode: utility.ToString(u["enterpriseCode"]),
						EnterpriseName: utility.ToString(u["enterpriseName"]),
					})
				}
			}
			ret.Users = us
		}
	}
	return &ret, nil
}
