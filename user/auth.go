package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/pkg/account"
	"github.com/33cn/chat33/pkg/http"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var log = log15.New("remote account")

// 根据token返回其他平台uid信息
func getUserInfo(url, zbToken string) (*types.User, error) {
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + zbToken
	byte, err := http.HTTPPostForm(url, headers, nil)
	if err != nil {
		//屏蔽请求超时错误详情，打印详情日志 dld-v2.8.4 2019年7月2日16:26:02
		log.Warn(fmt.Sprintf("http client request %v failed", url), "err", err.Error())
		return nil, errors.New("请求账户信息超时")
	}
	//成功：
	//{
	//	"code": 0,
	//	"data": {
	//		"email": "",
	//		"area": "86",
	//		"mobile": "13858075274",
	//		"username": "86138****5274",
	//		"id": "200093",
	//		"state" : 0,	//0：非实名认证 1：实名认证
	//	}
	//}
	//失败：
	//{
	//	"code": -100,
	//	"error": {
	//		"message":""
	//	}
	//}
	//-100 : 服务响应失败
	//-101 : token认证失败

	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return nil, err
	}

	//失败
	if utility.ToInt(resp["code"]) != 0 {
		if errMsg, ok := resp["error"].(map[string]interface{}); !ok {
			if msg, ok := resp["message"].(string); ok {
				return nil, errors.New(msg)
			}
			return nil, errors.New("获取用户信息返回不合法错误内容")
		} else {
			return nil, errors.New(errMsg["message"].(string))
		}
	}
	data, ok := resp["data"]
	if !ok {
		return nil, errors.New("no 'data' info")
	}

	info := data.(map[string]interface{})

	mobile := utility.ToString(info["mobile"])
	area := utility.ToString(info["area"])
	email := utility.ToString(info["email"])
	//是否实名
	verified := 0
	if utility.ToInt(info["state"]) == 0 {
		verified = 2
	} else {
		verified = 1
	}
	//账户号
	var account string
	if mobile != "" {
		account = area + mobile
	} else {
		account = email
	}

	return &types.User{
		Uid:      utility.ToString(info["id"]),
		Account:  account,
		Verified: verified,
		Email:    email,
		Phone:    mobile,
		Username: utility.ToString(info["username"]),
		Avatar:   utility.ToString(info["avatar"]),
	}, nil
}

func GetUserInfoFromAccountSystem(accountType, token string) (*types.User, error) {
	app := app.GetApp(accountType)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}
	switch app.IsInner {
	case types.IsInnerAccount:
		info, err := account.UserInfoFromBcoin(app.AccountServer, token)
		if err != nil {
			return nil, err
		}
		/*verified, err := verifiedFromBcoin(token)
		if err != nil {
			return nil, err
		}*/
		publicKey, err := account.PublicKeyFromBcoin(app.AccountServer, token)
		if err != nil {
			log.Warn("get remote public key filed", err)
			//return nil, err
		}
		info.DepositAddress = publicKey
		return info, nil
	default:
		return getUserInfo(app.AccountServer, token)
	}
}
