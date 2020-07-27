package api

import (
	"encoding/hex"
	"encoding/json"

	"github.com/33cn/chat33/comet"
	"github.com/33cn/chat33/comet/ws"
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/pkg/address"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logStatic = log15.New("module", "api/statistics")

//弃用 储存个推cid
func SaveGTCid(c *gin.Context) {
	c.Set(ReqError, nil)
}

func Push(c *gin.Context) {
	var params struct {
		Data interface{} `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//验证用户
	token := c.MustGet(Token).(string)
	appId := c.MustGet(AppId).(string)
	device := ""
	d, _ := c.Get(DeviceType)
	if d != nil {
		device = d.(string)
	}
	msg, err := json.Marshal(params.Data)
	if err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//鉴权
	tokenInfo, err := orm.GetToken(appId, token)
	if err != nil {
		if err == types.ERR_LOGIN_EXPIRED {
			c.Set(ReqError, result.NewError(result.LoginExpired))
			return
		}
		c.Set(ReqError, result.NewError(result.TokenLoginFailed).SetExtMessage(err.Error()))
		return
	}

	user, err := orm.GetUserInfoByUid(appId, tokenInfo.Uid)
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	if user == nil {
		logStatic.Warn("Push", "warn", "UserNotReg", "appId", appId, "uid", tokenInfo.Uid)
		c.Set(ReqError, result.NewError(result.UserNotReg).SetExtMessage(err.Error()))
		return
	}
	userId := user.UserId
	//创建客户端
	err = comet.HttpPush(appId, userId, device, msg)
	c.Set(ReqError, err)
}

// 精确搜索用户或群
func ClearlySearch(c *gin.Context) {
	type requestParams struct {
		Id string `json:"markId" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	appId := c.MustGet(AppId)
	if appId == nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetExtMessage("app_id 未获取"))
		return
	}

	userId := ""
	hUserId := c.MustGet(UserId)
	if hUserId != nil {
		userId = hUserId.(string)
	}

	info, err := model.ClearlySearch(appId.(string), userId, params.Id)
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 获取群会话秘钥
func RoomSessionKey(c *gin.Context) {
	type requestParams struct {
		Datetime *int64 `json:"datetime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	info, err := model.RoomSessionKey(c.MustGet(UserId).(string), *params.Datetime)
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 获取请求列表
func GetApplyList(c *gin.Context) {
	type requestParams struct {
		Id     string `json:"id"`
		Number int    `json:"number" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	info, err := model.GetApplyList(c.MustGet(UserId).(string), params.Id, params.Number)
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 获取请求列表
func GetUnreadApplyNumber(c *gin.Context) {
	info, err := orm.GetUnreadApplyNumber(c.MustGet(UserId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 撤回消息记录
func RevokeMsg(c *gin.Context) {
	type requestParams struct {
		LogId string `json:"logId" binding:"required"`
		Type  int    `json:"type"  binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.RevokeMsg(c.MustGet(UserId).(string), params.LogId, params.Type)
	c.Set(ReqError, err)
}

// 批量撤回文件消息
func RevokeFiles(c *gin.Context) {
	type requestParams struct {
		Logs []string `json:"logs" binding:"required"`
		Type int      `json:"type"  binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	info, err := model.BatchRevokeFiles(c.MustGet(UserId).(string), params.Logs, params.Type)
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 阅读指定一条阅后即焚消息
func ReadSnapMsg(c *gin.Context) {
	type requestParams struct {
		LogId string `json:"logId" binding:"required"`
		Type  int    `json:"type"  binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.ReadSnapMsg(c.MustGet(UserId).(string), params.LogId, params.Type)
	c.Set(ReqError, err)
}

// 消息转发
func ForwardMsg(c *gin.Context) {
	type requestParams struct {
		SourceId    string   `json:"sourceId" binding:"required"`
		Type        int      `json:"type"  binding:"required"`
		ForwardType int      `json:"forwardType"  binding:"required"`
		LogArray    []string `json:"logArray"  binding:"required"`
		TargetRooms []string `json:"targetRooms"  binding:"required"`
		TargetUsers []string `json:"targetUsers"  binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	rlt, err := model.ForwardMsg(c.MustGet(UserId).(string), params.SourceId, params.Type, params.ForwardType, params.LogArray, params.TargetRooms, params.TargetUsers)
	if rlt == nil {
		rlt = gin.H{}
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, rlt)
}

func ForwardEncryptMsg(c *gin.Context) {
	type requestParams struct {
		Type     int                      `json:"type"  binding:"required"`
		RoomLogs []map[string]interface{} `json:"roomLogs"  binding:"required"`
		UserLogs []map[string]interface{} `json:"userLogs"  binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	rlt, err := model.ForwardEncryptMsg(c.MustGet(UserId).(string), params.Type, params.RoomLogs, params.UserLogs)
	if rlt == nil {
		rlt = gin.H{}
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, rlt)
}

// 获取开屏广告
func Advertisement(c *gin.Context) {
	info, err := model.Advertisement(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

func UploadSecretKey(c *gin.Context) {
	var params struct {
		PublicKey  string `json:"publicKey"   binding:"required"`
		PrivateKey string `json:"privateKey"  binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	userId := c.MustGet(UserId).(string)
	device := c.MustGet(DeviceType).(string)
	uuid := c.MustGet(Uuid).(string)
	//TODO 通过好友关系获取userId
	friends, _ := orm.FindFriendsId(userId)

	var fs []string
	fs = append(fs, friends...)

	//去除字符串首的 "0x"
	str := params.PublicKey
	in, _ := hex.DecodeString(str[2:])
	//publickey转成地址
	uid := address.PublicKeyToAddress(address.NormalVer, in)

	MarkId := c.MustGet(AppId).(string) + uid
	err := orm.UpdateUid(MarkId, c.MustGet(UserId).(string), uid)
	if err != nil {
		c.Set(ReqError, result.NewError(result.WriteDbFailed))
		return
	}
	err = orm.UpdatePublicKey(userId, params.PublicKey, params.PrivateKey)
	if err != nil {
		c.Set(ReqError, result.NewError(result.WriteDbFailed))
		return
	}
	proto.SendUpdatePubKeyNotification(userId, params.PublicKey, params.PrivateKey, fs)

	ws.CloseOther(userId, device, uuid)

	c.Set(ReqError, err)
	c.Set(ReqResult, nil)
}
