package api

import (
	"io"

	"github.com/33cn/chat33/types"

	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/pkg/sms"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/user"
	"github.com/33cn/chat33/utility"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logUser = log15.New("module", "api/user")

//查询用户是否设置支付密码
func UserIsSetPayPwd(c *gin.Context) {
	user, err := user.GetUserInfoFromAccountSystem(c.MustGet(AppId).(string), c.MustGet(Token).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.VisitAccountSystemFailed).SetExtMessage(err.Error()))
		return
	}

	ret := struct {
		IsSetPayPwd int `json:"IsSetPayPwd"`
	}{
		IsSetPayPwd: user.IsSetPayPwd,
	}

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

func UserSetPayPwd(c *gin.Context) {
	type requestParams struct {
		Mode           string `json:"mode" binding:"required"`
		PayPassword    string `json:"payPassword" binding:"required"`
		Type           string `json:"type"`
		Code           string `json:"code"`
		OldPayPassword string `json:"oldPayPassword"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetPayPwd(c.MustGet(AppId).(string), c.MustGet(Token).(string), params.Mode, params.PayPassword, params.Type, params.Code, params.OldPayPassword)
	if err != nil {
		c.Set(ReqError, result.NewError(result.SetPayPwdError).SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
}

func UserCheckPayPwd(c *gin.Context) {
	type requestParams struct {
		PayPassword string `json:"payPassword" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.CheckPayPwd(c.MustGet(AppId).(string), c.MustGet(Token).(string), params.PayPassword)
	if err != nil {
		c.Set(ReqError, result.NewError(result.CheckPayPwdError).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
}

// 根据token登录
func UserTokenLogin(context *gin.Context) {
	type requestParams struct {
		Type     int    `json:"type" binding:"required"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}

	token := context.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-AUTH-TOKEN"))
		return
	}

	deviceType := context.GetHeader("FZM-DEVICE")
	if deviceType == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-DEVICE"))
		return
	}

	deviceName := context.GetHeader("FZM-DEVICE-NAME")
	if deviceName == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-DEVICE-NAME"))
		return
	}

	appId := context.GetHeader("FZM-APP-ID")
	if appId == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-APP-ID"))
		return
	}

	//非强制性
	uuid := context.GetHeader("FZM-UUID")
	//非强制性
	version := context.GetHeader("FZM-VERSION")

	var params requestParams
	if err := context.ShouldBindJSON(&params); err != nil {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	var data map[string]interface{}
	if pwd, ok := cfg.Env.Super[params.Phone]; ok && cfg.Env.Env == "debug" && pwd == params.Password {
		//超级用户
		logUser.Info("UserTokenLogin : super user login")
		u, err := orm.GetUserInfoByPhone(appId, params.Phone)
		if err != nil {
			context.Set(ReqError, result.NewError(result.DbConnectFail))
			return
		}
		if u == nil {
			context.Set(ReqError, result.NewError(result.UserNotExists))
			return
		}
		ret := make(map[string]interface{})
		ret["id"] = u.UserId
		ret["uid"] = u.Uid // zhaobi uid
		ret["account"] = u.Account
		ret["verified"] = u.Verified
		ret["username"] = u.Username
		ret["avatar"] = u.Avatar
		ret["code"] = u.InviteCode
		ret["isSetPayPwd"] = u.IsSetPayPwd
		ret["phone"] = u.Phone
		ret["depositAddress"] = u.DepositAddress
		ret["publicKey"] = u.PublicKey
		ret["privateKey"] = u.PrivateKey
		ret["inviteCode"] = u.InviteCode
		ret["identification"] = u.Identification
		ret["identificationInfo"] = u.IdentificationInfo

		data = ret
	} else {
		var err error
		data, err = model.DevTokenLogin(appId, params.Phone, token, deviceType, deviceName, utility.ToString(params.Type), uuid, version, apiAuthTypeToken)
		if err != nil {
			context.Set(ReqError, err)
			return
		}
	}

	userId := data["id"]
	session, err := store.Get(context.Request, sessionName)
	if session == nil {
		logUser.Error("UserTokenLogin: get session err", "err", err.Error())
		context.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}

	session.Values["user_id"] = userId
	session.Values["token"] = token
	session.Values["devtype"] = deviceType
	session.Values["uuid"] = uuid
	session.Values["appId"] = appId
	session.Values["time"] = utility.NowMillionSecond()
	err = session.Save(context.Request, context.Writer)
	if err != nil {
		logUser.Error("UserTokenLogin: save session err", "err", err.Error())
		context.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}

	logUser.Info("token login success", "appId", appId, "userId", userId, "token", token)
	context.Set(ReqError, nil)
	context.Set(ReqResult, data)
}

// 验证验证码，返回token
func PhoneLogin(context *gin.Context) {
	type requestParams struct {
		Phone string `json:"phone"`
		Code  string `json:"Code"`
	}
	deviceType := context.GetHeader("FZM-DEVICE")
	if deviceType == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-DEVICE"))
		return
	}

	appId := context.GetHeader("FZM-APP-ID")
	if appId == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-APP-ID"))
		return
	}
	//非强制性
	version := context.GetHeader("FZM-VERSION")

	var params requestParams
	if err := context.ShouldBindJSON(&params); err != nil {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//TODO 暂时写死 测试用
	if cfg.Env.Env == "debug" && params.Code != "111111" {
		//验证验证码是否正确
		err := sms.ValidateCode(cfg.SMS.Curl, cfg.SMS.CodeType, params.Phone, params.Code)
		if err != nil {
			context.Set(ReqError, result.NewError(result.VerifyCodeError))
			return
		}
	}

	//验证登录/注册
	data, err := model.PhoneLogin(appId, params.Phone, deviceType, version)
	if err != nil {
		context.Set(ReqError, err)
		return
	}

	context.Set(ReqError, nil)
	context.Set(ReqResult, data)
}

// 发送验证码
func SendCode(context *gin.Context) {
	type requestParams struct {
		Phone string `json:"phone"`
		//图形验证码所需
		//Ticket     string `json:"ticket"`
		//BusinessId string `json:"businessId"`
	}
	var params requestParams
	if err := context.ShouldBindJSON(&params); err != nil {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	//TODO 发送验证码方法
	_, err := sms.Send(cfg.SMS.Surl, cfg.SMS.CodeType, params.Phone, cfg.SMS.Msg, "", "")
	if err != nil {
		context.Set(ReqError, nil)
		return
	}
	context.Set(ReqError, nil)
}

func UserEditAvatar(c *gin.Context) {
	type editAvatarParams struct {
		Avatar string `json:"avatar" binding:"required"`
	}
	var params editAvatarParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := orm.UpdateUserAvatar(c.MustGet(UserId).(string), params.Avatar)
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	c.Set(ReqError, nil)
}

func UserEditNickname(c *gin.Context) {
	type editNickNameParams struct {
		Nickname string `json:"nickname" binding:"required"`
	}
	var params editNickNameParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := orm.UpdateUsername(c.MustGet(UserId).(string), params.Nickname)
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	c.Set(ReqError, nil)
}

//更新推送deviceToken
func UpdatePushDeviceToken(c *gin.Context) {
	type requestParams struct {
		DeviceToken string `json:"deviceToken" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	//存个推cid
	err := orm.SetUserDeviceToken(c.MustGet(UserId).(string), params.DeviceToken, c.MustGet(DeviceType).(string))
	if err != nil {
		logUser.Error("SetUserDeviceToken failed", "err", err)
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	c.Set(ReqError, nil)
}

//删除session  删除deviceToken
func UserLoginOut(c *gin.Context) {

	session, err := store.Get(c.Request, sessionName)
	if session == nil {
		logUser.Error("UserLogOut session get store", "err", err.Error())
		c.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}
	if session.Values["user_id"] != nil {
		//删除deviceToken
		err := orm.ClearUserDeviceToken(utility.ToString(session.Values["user_id"]))
		if err != nil {
			logUser.Error("ClearUserDeviceToken failed", "err", err)
		}
	}
	session.Values["user_id"] = nil
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logUser.Error("UserLoginOut: save session err", "err", err.Error())
		c.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}
	c.Set(ReqError, nil)
}

//查看用户设置
func UserConf(c *gin.Context) {
	conf, err := model.FindUserConf(c.MustGet(AppId).(string), c.MustGet(UserId).(string))
	if err != nil {
		c.Set(ReqError, err)
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, conf)
}

//付款
func Payment(c *gin.Context) {
	type editNickNameParams struct {
		LogId      string  `json:"logId"`
		Currency   *string `json:"currency" binding:"required"`
		Amount     *string `json:"amount" binding:"required"`
		Fee        *string `json:"fee" binding:"required"`
		OppAddress *string `json:"opp_address" binding:"required"`
		Rid        *string `json:"rid" binding:"required"`
		Mode       *string `json:"mode" binding:"required"`
		Payword    *string `json:"payword" binding:"required"`
		Code       *string `json:"code" binding:"required"`
	}
	var params editNickNameParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.Payment(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.LogId, c.MustGet(Token).(string),
		*params.Currency, utility.ToFloat64(*params.Amount), utility.ToFloat64(*params.Fee), *params.OppAddress, *params.Rid, *params.Mode, *params.Payword, *params.Code)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//邀请奖励统计信息
func InviteStatistics(c *gin.Context) {
	ret, err := model.InviteStatistics(c.MustGet(AppId).(string), c.MustGet(Token).(string))
	if err != nil {
		logUser.Warn(" InviteStatistics", "warn", err.Error())
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

// 单次邀请奖励列表
func SingleInviteInfos(c *gin.Context) {
	var params struct {
		Page int64 `json:"page" binding:"required"`
		Size int64 `json:"size"`
	}
	params.Page = 1
	params.Size = 5
	if err := c.ShouldBindJSON(&params); err != nil && err != io.EOF {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.SingleInviteInfos(c.MustGet(AppId).(string), c.MustGet(Token).(string), params.Page, params.Size)
	if err != nil {
		logUser.Warn(" SingleInviteInfos", "warn", err.Error())
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

// 累加邀请奖励列表
func AccumulateInviteInfos(c *gin.Context) {
	var params struct {
		Page int64 `json:"page" binding:"required"`
		Size int64 `json:"size"`
	}
	params.Page = 1
	params.Size = 5
	if err := c.ShouldBindJSON(&params); err != nil && err != io.EOF {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AccumulateInviteInfos(c.MustGet(AppId).(string), c.MustGet(Token).(string), params.Page, params.Size)
	if err != nil {
		logUser.Warn(" AccumulateInviteInfos", "err", err.Error(), "appId", c.MustGet(AppId).(string), "page", params.Page, "size", params.Size)
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//设置入群邀请确认
func SetInviteConfirm(c *gin.Context) {
	var params struct {
		NeedConfirmInvite *int `json:"needConfirmInvite" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	switch *params.NeedConfirmInvite {
	case types.RoomInviteNotNeedConfirm:
	case types.RoomInviteNeedConfirm:
	default:
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized needConfirmInvite type "))
		return
	}

	err := model.SetInviteConfirm(c.MustGet(UserId).(string), *params.NeedConfirmInvite)
	c.Set(ReqError, err)
}

//认证申请
func VerifyApplyForUser(c *gin.Context) {
	var params struct {
		Description string `json:"description" binding:"required"`
		PayPassword string `json:"payPassword" binding:"required"`
		Code        string `json:"code"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	_, err := model.VerifyApply(c.MustGet(AppId).(string), c.MustGet(Token).(string), types.VerifyForUser, c.MustGet(UserId).(string), params.Description, params.PayPassword, params.Code)
	c.Set(ReqError, err)
	//c.Set(ReqResult, ret)
}

//认证审核状态
func UserVerificationState(c *gin.Context) {
	state, err := model.CheckUserVerifyState(c.MustGet(UserId).(string))
	ret := map[string]interface{}{
		"state": state,
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

// 合约登陆接口
func UserLogin(context *gin.Context) {
	type requestParams struct {
		Type int `json:"type" binding:"required"`
	}

	token := context.GetHeader("FZM-AUTH-TOKEN")
	if token == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-AUTH-TOKEN"))
		return
	}

	deviceType := context.GetHeader("FZM-DEVICE")
	if deviceType == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-DEVICE"))
		return
	}

	deviceName := context.GetHeader("FZM-DEVICE-NAME")
	if deviceName == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-DEVICE-NAME"))
		return
	}

	appId := context.GetHeader("FZM-APP-ID")
	if appId == "" {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "Header:FZM-APP-ID"))
		return
	}

	//非强制性
	uuid := context.GetHeader("FZM-UUID")
	//非强制性
	version := context.GetHeader("FZM-VERSION")

	var params requestParams
	if err := context.ShouldBindJSON(&params); err != nil {
		context.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	data, err := model.TokenLogin(appId, token, deviceType, deviceName, utility.ToString(params.Type), uuid, version, apiAuthTypeToken)
	if err != nil {
		context.Set(ReqError, err)
		return
	}
	userId := data["id"]
	session, err := store.Get(context.Request, sessionName)
	if session == nil {
		logUser.Error("tokenLogin session get store", "err", err.Error())
		context.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}

	session.Values["user_id"] = userId
	session.Values["token"] = token
	session.Values["devtype"] = deviceType
	session.Values["uuid"] = uuid
	session.Values["appId"] = appId
	session.Values["time"] = utility.NowMillionSecond()
	err = session.Save(context.Request, context.Writer)
	if err != nil {
		logUser.Error("tokenLogin session save", "err", err.Error())
		context.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}

	logUser.Info("token login success", "appId", appId, "userId", userId, "token", token)
	context.Set(ReqError, nil)
	context.Set(ReqResult, data)
}

//TODO 上链
func IsChain(c *gin.Context) {
	//1代表上链，上了就上了，没有下链的说法
	IsChain := 1
	err := orm.UpdateIsChain(c.MustGet(UserId).(string), IsChain)
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	c.Set(ReqError, nil)
}
