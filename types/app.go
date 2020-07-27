package types

import (
	"errors"
)

const (
	IsInnerAccount    = 1
	IsNotInnerAccount = 0

	IsOtc = 1
)

var ERR_APPNOTFIND = errors.New("app未认证")
var ERR_NOTSUPPORT = errors.New("暂不支持该功能")
var ERR_SYS_RESULT = errors.New("该功能无法使用")

var ERR_REPEAT_MSG = errors.New("重复发送的消息")

var ERR_LOGIN_EXPIRED = errors.New("登录过期")
