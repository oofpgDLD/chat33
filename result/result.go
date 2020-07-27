package result

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

const (
//ErrSessionKeyExpired = error()
)

const (
	ServiceChat = "Chat"
	ServiceRP   = "Redpacket"
	ServiceCoin = "Bcoin"
)

const (
	DispMsgExt  = 0
	DispJustExt = 1
)

// TODO
const (
	CodeOK = 0
	//ServerResponseFailed   = -100
	//TokenFailed            = -101
	DbConnectFail            = -1000
	ParamsError              = -1001
	LackParam                = -1002
	SessionError             = -1003
	LoginExpired             = -1004
	MsgFormatError           = -1005
	UnknowDeviceType         = -1006
	UserLoginOtherDevice     = -1007
	LackHeader               = -1008
	ServiceReqFailed         = -1009
	QueryDbFailed            = -1010
	WriteDbFailed            = -1011
	TooManyRequests          = -1012
	DbParamsError            = -1013
	AdminExportFailed        = -1014
	AdminLoginFailed         = -1015
	AdminAdNumbersLimit      = -1016
	AdUploadFailed           = -1017
	AppUpdateFailed          = -1018
	RequestHeaderError       = -1019
	VerifyLimit              = -1020
	UserNotExists            = -2001
	UserExists               = -2002
	CSExists                 = -2003
	TokenLoginFailed         = -2004
	VisitAccountSystemFailed = -2005
	UserSendSysMsg           = -2006
	UserMuted                = -2007
	VisitorSendMsg           = -2008
	SendPrivMsg              = -2009
	UserIsServed             = -2010
	JoinListenFailed         = -2011
	ZhaobiInteractFailed     = -2013
	SendPrivMsgBtwCS         = -2014
	VisitorOffline           = -2015
	NoCSOnline               = -2016
	IsFriendAlready          = -2017
	IsNotFriend              = -2018
	FriendRequestHadDeal     = -2019
	NotExistFriendRequest    = -2020
	CanNotOperateSelf        = -2021
	ConvFail                 = -2022
	DeleteMsgFailed          = -2023
	CanNotDelMsgOverTime     = -2024
	CanNotFindRoomMsg        = -2025
	ObjIsMaster              = -2026
	ObjIsManager             = -2027
	SetPayPwdError           = -2028
	CheckPayPwdError         = -2029
	AccountClosedByAdmin     = -2030
	TradeRelationshipNotExt  = -2031
	IsBlocked                = -2032
	PermissionDeny           = -3000
	CannotJoinGroup          = -3001
	NoCSPermission           = -3002
	NoEditGroupPermission    = -3003
	AlreadyJoinOutRoom       = -3004
	CanNotFindFriendMsg      = -3005
	RPError                  = -4000
	RPEmpty                  = -4001
	UserNotReg               = -4002
	UserHasReg               = -4003
	OnlyForNewUser           = -4004
	RPIdNotMatch             = -4005
	RPIdIllegal              = -4006
	VerifyCodeError          = -4007
	VerifyCodeExpired        = -4008
	RPReceived               = -4009
	CannotSendRP             = -4010
	RPTargetIllegal          = -4011
	ApplyAlreadyDeal         = -4012
	RPExpired                = -4013
	PaymentFailed            = -4014
	PraiseAmountLimited      = -4015
	PraiseRewardErr          = -4016
	GroupNotExists           = -5000
	UserNotEnterGroup        = -5001
	QueryChatLogFailed       = -6000
	RoomClosedByAdmin        = -6001
	IsRoomMemberAlready      = -6100
	RoomNotExists            = -6101
	CanNotInvite             = -6102
	CanNotJoinRoom           = -6103
	CanNotLoginOut           = -6104
	UserIsNotInRoom          = -6105
	CanNotAddFriendInRoom    = -6106
	UserIsNotOnline          = -6107
	MembersOutOfLimit        = -6108
	CreateRoomsOutOfLimit    = -6109
	MustEncryptMsgToRoom     = -6110
	CanNotJoinEncryptRoom    = -6111
	ChatLogNotFind           = -6200
	UserIsNotRoomMember      = -6201
	WSMsgFormatError         = -7000
	GetVerifyCodeFailed      = -8000
	EditRewardFailed         = -8100
	ShowRewardFailed         = -8101
	ServerInterError         = -9000
	NetWorkError             = -9001
	FunctionNotSupport       = -9002
	//....
)

var errorMsg = map[int]string{
	//ServerResponseFailed:   "服务响应失败",
	//TokenFailed:            "token认证失败",
	RequestHeaderError:       "请求头错误",
	CodeOK:                   "操作成功",
	DbConnectFail:            "数据库连接失败",
	ParamsError:              "参数错误",
	DbParamsError:            "数据库字段错误",
	LackParam:                "缺少参数",
	SessionError:             "Session错误",
	LoginExpired:             "登录过期",
	MsgFormatError:           "消息格式错误",
	UnknowDeviceType:         "未知的设备类型",
	VerifyLimit:              "近期内不可申请认证",
	UserLoginOtherDevice:     "账号已经在其他终端登录",
	QueryDbFailed:            "查询数据库失败",
	WriteDbFailed:            "写入数据库失败",
	AppUpdateFailed:          "检查更新失败",
	TooManyRequests:          "发送频率过快，请稍后再试",
	UserNotExists:            "用户不存在",
	UserExists:               "用户已存在",
	CSExists:                 "要添加的客服已存在",
	TokenLoginFailed:         "登录验证失败",
	VisitAccountSystemFailed: "访问账户系统失败",
	UserSendSysMsg:           "用户没有发系统消息权限",
	UserMuted:                "用户被禁言",
	VisitorSendMsg:           "游客没有发消息权限",
	SendPrivMsg:              "没有给普通用户发私聊权限",
	UserIsServed:             "当前用户已经被其他客服接待",
	JoinListenFailed:         "加入旁听失败",
	ZhaobiInteractFailed:     "找币交互失败",
	SendPrivMsgBtwCS:         "客服间不能发送私聊消息",
	VisitorOffline:           "游客已离线",
	NoCSOnline:               "暂无客服在线，请稍后再试，或可登录账号给客服留言!",
	PermissionDeny:           "权限不足",
	CannotJoinGroup:          "用户没有加入聊天群权限",
	NoCSPermission:           "没有客服权限",
	NoEditGroupPermission:    "没有修改聊天室权限",
	AlreadyJoinOutRoom:       "已退出群聊",
	CanNotInvite:             "不可邀请好友",
	CanNotJoinRoom:           "不可加入该群",
	CanNotLoginOut:           "群主不可退出群",
	RoomNotExists:            "群不存在",
	IsRoomMemberAlready:      "已经是群成员",
	MustEncryptMsgToRoom:     "请先设置助记词",
	RPError:                  "红包错误",
	RPEmpty:                  "红包已被领完",
	UserNotReg:               "用户未注册",
	UserHasReg:               "用户已注册",
	OnlyForNewUser:           "仅限新人领取",
	RPIdNotMatch:             "红包标识不匹配",
	RPIdIllegal:              "非法的红包ID",
	VerifyCodeError:          "验证码不正确",
	VerifyCodeExpired:        "验证码已经过期或者已使用",
	RPReceived:               "红包已领取",
	RPExpired:                "红包已过期",
	CannotSendRP:             "用户无发红包权限",
	GroupNotExists:           "聊天室不存在",
	UserNotEnterGroup:        "用户未进入此聊天室",
	QueryChatLogFailed:       "查询聊天记录失败",
	WSMsgFormatError:         "消息格式错误",
	GetVerifyCodeFailed:      "获取手机验证码失败",
	ServerInterError:         "服务端内部错误",

	IsFriendAlready:       "对方已经是您的好友",
	IsNotFriend:           "对方不是您的好友",
	FriendRequestHadDeal:  "好友请求已经被处理",
	NotExistFriendRequest: "好友请求不存在",
	CanNotOperateSelf:     "不能对自己进行操作",
	ConvFail:              "数据转换异常",
	DeleteMsgFailed:       "删除消息失败",
	CanNotDelMsgOverTime:  "消息发出已超过10分钟",
	CanNotFindRoomMsg:     "找不到该聊天记录",
	CanNotFindFriendMsg:   "找不到该聊天记录",
	ObjIsMaster:           "对方是群主",
	ObjIsManager:          "对方是管理员",

	NetWorkError:            "网络错误，请重试",
	UserIsNotInRoom:         "用户不在群中",
	CanNotAddFriendInRoom:   "该群不允许添加好友",
	UserIsNotOnline:         "用户不在线",
	RPTargetIllegal:         "无法领取该红包",
	ApplyAlreadyDeal:        "申请已经被处理",
	SetPayPwdError:          "修改支付密码失败",
	CheckPayPwdError:        "支付密码错误",
	LackHeader:              "缺少Header",
	ServiceReqFailed:        "请求外部服务失败", //只返回拓展错误信息
	PaymentFailed:           "付款失败",
	AccountClosedByAdmin:    "账号被管理员封禁",
	RoomClosedByAdmin:       "群被管理员封禁", //该群已经被管理员封禁
	MembersOutOfLimit:       "群人数已经达到上限",
	CreateRoomsOutOfLimit:   "可创建的群聊数量达到上限",
	AdminExportFailed:       "导出excel失败",
	AdminLoginFailed:        "管理员账号或密码错误",
	AdminAdNumbersLimit:     "可添加广告数达到上限",
	AdUploadFailed:          "上传广告失败",
	CanNotJoinEncryptRoom:   "该用户未设置密聊私钥，无法加入加密群聊",
	EditRewardFailed:        "设置奖励失败",
	ShowRewardFailed:        "查看奖励失败",
	TradeRelationshipNotExt: "交易关系不存在",
	IsBlocked:               "对方拒绝接收你的消息！",
	PraiseAmountLimited:     "已超过当日限额",
	PraiseRewardErr:         "打赏失败",
	ChatLogNotFind:          "消息不存在",
	FunctionNotSupport:      "暂未开放",
}

//错误信息组合
//1. message + : + extMsg
//2. extMsg
func NewError(code int) *Error {
	return &Error{
		ErrorCode: code,
		Message:   ParseError(code, ""),
	}
}

type Error struct {
	ErrorCode int
	//errorCode对应错误
	Message string
	//额外错误信息
	ExtMsg string
	//暴露给接口，但客户端不显示
	Data map[string]interface{}
	//显示方式
	displayType int
}

//策略返回显示的错误信息
func (e *Error) Error() string {
	switch e.displayType {
	case DispJustExt:
		return e.ExtMsg
	default:
		tag := ""
		if e.ExtMsg != "" {
			tag = ":"
		}
		return fmt.Sprintf("%s%s%s", e.Message, tag, e.ExtMsg)
	}
}

func (e *Error) JustShowExtMsg() *Error {
	e.displayType = DispJustExt
	return e
}

func (e *Error) SetExtMessage(extMsg string) *Error {
	e.ExtMsg = extMsg
	return e
}

//服务名 + code + message
func (e *Error) SetChildErr(name string, code interface{}, message interface{}) *Error {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data["service"] = name
	e.WriteCode(code)
	e.WriteMessage(message)
	return e
}

func (e *Error) WriteMessage(msg interface{}) *Error {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data["message"] = msg
	return e
}

func (e *Error) WriteCode(code interface{}) *Error {
	if e.Data == nil {
		e.Data = make(map[string]interface{})
	}
	e.Data["code"] = code
	return e
}

//websocket err
type WsError struct {
	E         *Error
	EventCode int
	MsgId     string
}

func (e *WsError) Error() string {
	return e.E.Error()
}

func ParseError(errcode int, msg string) string {
	errStr, ok := errorMsg[errcode]
	if !ok {
		log15.Warn("ParseError error code not exists", "errcode", errcode)
		return msg
	}
	if errStr != "" {
		errStr += ":"
	}
	return strings.Trim(utility.ParseString(errStr+"%v", msg), " :")
}

func ComposeWsError(err error) ([]byte, error) {
	if e, ok := err.(*WsError); ok {
		type WsAck struct {
			EventType int    `json:"eventType"`
			MsgId     string `json:"msgId"`
			Code      int    `json:"code"`
			Content   string `json:"content"`
		}

		var ret WsAck
		ret.EventType = e.EventCode
		ret.MsgId = e.MsgId
		ret.Code = e.E.ErrorCode
		ret.Content = e.Error()

		return json.Marshal(&ret)
	} else {
		return nil, errors.New("not ws error type")
	}
}

func ComposeHttpAck(code int, msg string, data interface{}) interface{} {
	type HttpAck struct {
		Result  int         `json:"result"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	var ret HttpAck
	ret.Result = code
	ret.Message = msg
	ret.Data = data
	return &ret
}
