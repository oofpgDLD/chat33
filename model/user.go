package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/33cn/chat33/user"

	"github.com/33cn/chat33/pkg/account"
	cmn "github.com/33cn/chat33/pkg/btrade/common"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	user_model "github.com/33cn/chat33/user/model"
	"github.com/33cn/chat33/utility"
)

var logUser = log15.New("module", "model/user")

func getUsernameById(id string) (string, error) {
	info, err := orm.GetUserInfoById(id)
	if err != nil {
		return "", err
	}
	return info.Username, err
}

/*func PushMsgToUser(appId, userId, friendId, geTuiHeader, geTuiConent string, msg *proto.Proto) error {
	return GTSinglePush(appId, userId, friendId, geTuiHeader, geTuiConent, msg)
}*/

func CheckUserExist(userId string) bool {
	u, err := orm.GetUserInfoById(userId)
	if err != nil || u == nil {
		return false
	}
	return true
}

//---------------------------------------------------------------------------------------//
func SetPayPwd(accountType, token, mode, payPassword, authType, code, oldPayPassword string) error {
	app := app.GetApp(accountType)
	if app == nil {
		return types.ERR_APPNOTFIND
	}
	switch app.IsInner {
	case types.IsInnerAccount:
		return account.BcoinSetPayPwd(app.AccountServer, token, mode, payPassword, authType, code, oldPayPassword)
	default:
		return errors.New("暂不支持修改支付密码")
	}
}

func CheckPayPwd(accountType, token, payPassword string) error {
	app := app.GetApp(accountType)
	if app == nil {
		return types.ERR_APPNOTFIND
	}
	switch app.IsInner {
	case types.IsInnerAccount:
		return account.BcoinCheckPayPwd(app.AccountServer, token, payPassword)
	default:
		return errors.New("暂不支持验证支付密码")
	}
}

func PhoneLogin(appId, phone, deviceType, version string) (map[string]interface{}, error) {
	//通过手机号查询user信息，如果存在，登入；如果不存在，注册。
	ret := make(map[string]interface{})
	var token string
	//1表示登陆
	t := 1
	sysInfo, err := orm.GetUserInfoByPhone(appId, phone)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	createtime := utility.NowMillionSecond()
	//用户不存在的情况 新增用户和token
	if sysInfo == nil {
		//表示注册
		t = 0
		logUser.Debug("create new user")
		Username := utility.RandomUsername()
		UserLevel := types.LevelMember
		//随机生成的token
		token = utility.GetToken()
		//账户号
		var account string
		area := "86"
		if area == "" {
			area = "86"
		}
		account = area + phone

		//TODO 新增用户 第一次登入uid,markId等字段无法生成  verified 默认是0
		id, err := orm.InsertUser("", "", appId, Username, account, "", area, phone,
			cmn.ToString(UserLevel), "0", "", "", deviceType, version, utility.NowMillionSecond())

		//id, err := orm.AddUser("", "", appId, Username, phone, cmn.ToString(UserLevel), deviceType, version, utility.NowMillionSecond())
		if err != nil {
			return nil, result.NewError(result.TokenLoginFailed)
		}
		//存到user和 token的关系表里
		_, err = orm.AddToken(utility.ToString(id), token, createtime)
		if err != nil {
			return nil, result.NewError(result.TokenLoginFailed)
		}
	} else {
		DeviceToken, time, err := orm.FindToken(sysInfo.UserId)
		//用户存在 但token关系不存在的情况
		if DeviceToken == "" {
			//随机生成的token
			token = utility.GetToken()
			_, err = orm.AddToken(sysInfo.UserId, token, createtime)
			if err != nil {
				return nil, result.NewError(result.TokenLoginFailed)
			}
		} else {
			//判断是否过期
			_, isOverdue := utility.CheckToken(time)
			//如果已经过期 新建一个token
			if isOverdue == true {
				token = utility.GetToken()
				_, err = orm.AddToken(sysInfo.UserId, token, createtime)
				if err != nil {
					return nil, result.NewError(result.TokenLoginFailed)
				}
				//未过期
			} else {
				token = DeviceToken
			}
		}
		logUser.Debug("login get check user success")

	}
	ret["type"] = t
	ret["token"] = token

	return ret, nil
}

func EmailLogin(appId, email, deviceType, version string) (map[string]interface{}, error) {
	//通过手机号查询user信息，如果存在，登入；如果不存在，注册。
	ret := make(map[string]interface{})
	var token string
	//1表示登陆
	t := 1
	sysInfo, err := orm.GetUserInfoByEmail(appId, email)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	createtime := utility.NowMillionSecond()
	//用户不存在的情况 新增用户和token
	if sysInfo == nil {
		//表示注册
		t = 0
		logUser.Debug("create new user")
		Username := utility.RandomUsername()
		UserLevel := types.LevelMember
		//随机生成的token
		token = utility.GetToken()
		//账户号
		account := email

		//TODO 新增用户 第一次登入uid,markId等字段无法生成  verified 默认是0
		id, err := orm.InsertUser("", "", appId, Username, account, email, "", "",
			cmn.ToString(UserLevel), "0", "", "", deviceType, version, utility.NowMillionSecond())

		//id, err := orm.AddUser("", "", appId, Username, phone, cmn.ToString(UserLevel), deviceType, version, utility.NowMillionSecond())
		if err != nil {
			return nil, result.NewError(result.TokenLoginFailed)
		}
		//存到user和 token的关系表里
		_, err = orm.AddToken(utility.ToString(id), token, createtime)
		if err != nil {
			return nil, result.NewError(result.TokenLoginFailed)
		}
	} else {
		DeviceToken, time, err := orm.FindToken(sysInfo.UserId)
		//用户存在 但token关系不存在的情况
		if DeviceToken == "" {
			//随机生成的token
			token = utility.GetToken()
			_, err = orm.AddToken(sysInfo.UserId, token, createtime)
			if err != nil {
				return nil, result.NewError(result.TokenLoginFailed)
			}
		} else {
			//判断是否过期
			_, isOverdue := utility.CheckToken(time)
			//如果已经过期 新建一个token
			if isOverdue == true {
				token = utility.GetToken()
				_, err = orm.AddToken(sysInfo.UserId, token, createtime)
				if err != nil {
					return nil, result.NewError(result.TokenLoginFailed)
				}
				//未过期
			} else {
				token = DeviceToken
			}
		}
		logUser.Debug("login get check user success")

	}
	ret["type"] = t
	ret["token"] = token

	return ret, nil
}

func TokenLogin(appId, token, deviceType, deviceName, loginType, uuid, version string, apiTokenAuth bool) (map[string]interface{}, error) {
	//从对应的账户体系获取用户信息
	sysInfo, err := user.GetUserInfoFromAccountSystem(appId, token)
	if err != nil {
		if err == types.ERR_LOGIN_EXPIRED {
			return nil, result.NewError(result.LoginExpired)
		}
		return nil, result.NewError(result.TokenLoginFailed).SetExtMessage(err.Error())
	}
	if sysInfo == nil {
		return nil, result.NewError(result.UserNotExists)
	}
	if apiTokenAuth {
		//更新缓存
		tokenInfo := &user_model.Token{
			Uid:   sysInfo.Uid,
			AppId: appId,
			Token: token,
			Time:  utility.NowMillionSecond(),
		}
		err = orm.SaveToken(tokenInfo)
		if err != nil {
			logUser.Warn("redis can not save user")
		}
	}

	logUser.Debug("login get account success")
	firstLogin := false //是否第一次登陆
	//根据平台uid获取用户信息
	user, err := orm.GetUserInfoByUid(appId, sysInfo.Uid)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	if user == nil {
		//原先为了解决账户体系变更导致uid重复的问题，添加了判断账号是否存在。
		//但是会导致，如果账户体系返回手机号不正确，会导致数据库的uid被脏数据污染。
		//故去掉根据手机号判断用户是否存在 - 2019年8月21日17:55:55 dld
		logUser.Debug("create new user")
		//新建用户
		if sysInfo.Username == "" {
			sysInfo.Username = utility.RandomUsername()
		}
		sysInfo.UserLevel = types.LevelMember
		sysInfo.MarkId = appId + sysInfo.Uid

		id, err := orm.InsertUser(sysInfo.MarkId, sysInfo.Uid, appId, sysInfo.Username, sysInfo.Account, sysInfo.Email, sysInfo.Area, sysInfo.Phone, cmn.ToString(sysInfo.UserLevel), cmn.ToString(sysInfo.Verified), "", sysInfo.DepositAddress, deviceType, version, utility.NowMillionSecond())
		if err != nil {
			return nil, result.NewError(result.TokenLoginFailed)
		}
		firstLogin = true
		sysInfo.UserId = utility.ToString(id)
		/*//查询账号 是否存在
		user, err := orm.GetUserInfoByAccount(appId, sysInfo.Account)
		if user != nil { //账户存在
			//更新uid
			err = orm.UpdateUid(user.UserId, sysInfo.Uid)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}

			if sysInfo.InviteCode != user.InviteCode {
				err = orm.UpdateInviteCode(user.UserId, sysInfo.InviteCode)
				if err != nil {
					return nil, result.NewError(result.DbConnectFail)
				}
			}

			//更新 钱包地址
			if user.DepositAddress != sysInfo.DepositAddress && sysInfo.DepositAddress != "" {
				err := orm.UpdateDepositAddress(user.UserId, sysInfo.DepositAddress)
				if err != nil {
					logUser.Error("UpdateDepositAddress failed", "err", err)
				}
			}
			if utility.Version(user.NowVersion).Compare(utility.Version(version), utility.LT) {
				err := orm.UpdateNowVersion(user.UserId, version)
				if err != nil {
					logUser.Error("UpdateNowVersion failed", "err", err)
				}
			}
			//从账户系统中取出的信息
			user.IsSetPayPwd = sysInfo.IsSetPayPwd
			user.Verified = sysInfo.Verified
			user.Uid = sysInfo.Uid
			user.InviteCode = sysInfo.InviteCode
			sysInfo = user
		} else {
			//新建用户
			if sysInfo.Username == "" {
				sysInfo.Username = utility.RandomUsername()
			}
			sysInfo.UserLevel = types.LevelMember
			sysInfo.MarkId = appId + sysInfo.Uid

			id, err := orm.InsertUser(sysInfo.MarkId, sysInfo.Uid, appId, sysInfo.Username, sysInfo.Account, sysInfo.Email, sysInfo.Area, sysInfo.Phone, cmn.ToString(sysInfo.UserLevel), cmn.ToString(sysInfo.Verified), "", sysInfo.DepositAddress, deviceType, version, utility.NowMillionSecond())
			if err != nil {
				return nil, &result.Error{ErrorCode: result.TokenLoginFailed}
			}
			firstLogin = true
			sysInfo.UserId = utility.ToString(id)
		}*/
	} else {
		logUser.Debug("find user success")
		if user.Phone != sysInfo.Phone && sysInfo.Phone != "" {
			err = orm.UpdatePhone(user.UserId, sysInfo.Phone)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			user.Phone = sysInfo.Phone
		}

		if user.Email != sysInfo.Email && sysInfo.Email != "" {
			err = orm.UpdateEmail(user.UserId, sysInfo.Email)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			user.Email = sysInfo.Email
		}

		if sysInfo.InviteCode != user.InviteCode {
			err = orm.UpdateInviteCode(user.UserId, sysInfo.InviteCode)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			user.InviteCode = sysInfo.InviteCode
		}

		//更新 钱包地址
		if user.DepositAddress != sysInfo.DepositAddress && sysInfo.DepositAddress != "" {
			err := orm.UpdateDepositAddress(user.UserId, sysInfo.DepositAddress)
			if err != nil {
				logUser.Error("UpdateDepositAddress failed", "err", err)
			}
			user.DepositAddress = sysInfo.DepositAddress
		}
		//如果当前版本号大于之前版本号，更新当前使用版本号
		if utility.Version(user.NowVersion).Compare(utility.Version(version), utility.LT) {
			err := orm.UpdateNowVersion(user.UserId, version)
			if err != nil {
				logUser.Error("UpdateNowVersion failed", "err", err)
			}
		}
		//从账户系统中取出的信息
		user.IsSetPayPwd = sysInfo.IsSetPayPwd
		user.Verified = sysInfo.Verified
		sysInfo = user
	}
	logUser.Debug("login get check user success")
	if sysInfo.CloseUntil >= utility.NowMillionSecond() {
		return nil, result.NewError(result.AccountClosedByAdmin).JustShowExtMsg().SetExtMessage(ConvertUserClosedAlertStr(appId, sysInfo.CloseUntil))
	}
	logUser.Info("Token Login", "userID", sysInfo.UserId, "uid", sysInfo.Uid, "account", sysInfo.Account, "level", sysInfo.UserLevel)

	ret := make(map[string]interface{})
	ret["firstLogin"] = firstLogin
	ret["id"] = sysInfo.UserId
	ret["uid"] = sysInfo.Uid // zhaobi uid
	ret["account"] = sysInfo.Account
	ret["verified"] = sysInfo.Verified
	ret["username"] = sysInfo.Username
	ret["avatar"] = sysInfo.Avatar
	ret["code"] = sysInfo.InviteCode
	ret["isSetPayPwd"] = sysInfo.IsSetPayPwd
	ret["phone"] = sysInfo.Phone
	ret["depositAddress"] = sysInfo.DepositAddress
	ret["publicKey"] = sysInfo.PublicKey
	ret["privateKey"] = sysInfo.PrivateKey
	ret["inviteCode"] = sysInfo.InviteCode
	ret["identification"] = sysInfo.Identification
	ret["identificationInfo"] = sysInfo.IdentificationInfo

	if sysInfo.UserLevel == types.LevelAdmin {
		ret["user_level"] = types.LevelCs
	} else {
		ret["user_level"] = sysInfo.UserLevel
	}

	if deviceType == types.DeviceWeb {
		if sysInfo.UserLevel == types.LevelMember {
			return nil, result.NewError(result.UnknowDeviceType)
		}
	}

	_, err = orm.InsertLoginLog(sysInfo.UserId, deviceType, deviceName, loginType, uuid, version, utility.NowMillionSecond())
	if err != nil {
		logUser.Error("token login log user login failed", "err_msg", err)
	}
	return ret, nil
}

//新版token登入使用
func DevTokenLogin(appId, phone, token, deviceType, deviceName, loginType, uuid, version string, apiTokenAuth bool) (map[string]interface{}, error) {
	//根据phone获取用户信息
	sysInfo, err := orm.GetUserInfoByToken(appId, token)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if sysInfo == nil {
		return nil, result.NewError(result.UserNotExists)
	}
	if apiTokenAuth {
		//更新缓存
		tokenInfo := &user_model.Token{
			Uid:   sysInfo.Uid,
			AppId: appId,
			Token: token,
			Time:  utility.NowMillionSecond(),
		}
		err = orm.SaveToken(tokenInfo)
		if err != nil {
			logUser.Warn("redis can not save user")
		}
	}

	logUser.Debug("login get account success")

	//判断是否上链，默认未上链 0
	var ischain bool = false
	if sysInfo.IsChain == 1 {
		ischain = true
	}

	logUser.Debug("login get check user success")
	if sysInfo.CloseUntil >= utility.NowMillionSecond() {
		return nil, result.NewError(result.AccountClosedByAdmin).JustShowExtMsg().SetExtMessage(ConvertUserClosedAlertStr(appId, sysInfo.CloseUntil))
	}
	logUser.Info("Token Login", "userID", sysInfo.UserId, "uid", sysInfo.Uid, "account", sysInfo.Account, "level", sysInfo.UserLevel)

	ret := make(map[string]interface{})

	//是否上链
	ret["isChain"] = ischain
	ret["id"] = sysInfo.UserId
	ret["uid"] = sysInfo.Uid // zhaobi uid
	ret["account"] = sysInfo.Account
	ret["verified"] = sysInfo.Verified
	ret["username"] = sysInfo.Username
	ret["avatar"] = sysInfo.Avatar
	ret["code"] = sysInfo.InviteCode
	ret["isSetPayPwd"] = sysInfo.IsSetPayPwd
	ret["phone"] = sysInfo.Phone
	ret["depositAddress"] = sysInfo.DepositAddress
	ret["publicKey"] = sysInfo.PublicKey
	ret["privateKey"] = sysInfo.PrivateKey
	ret["inviteCode"] = sysInfo.InviteCode
	ret["identification"] = sysInfo.Identification
	ret["identificationInfo"] = sysInfo.IdentificationInfo

	if sysInfo.UserLevel == types.LevelAdmin {
		ret["user_level"] = types.LevelCs
	} else {
		ret["user_level"] = sysInfo.UserLevel
	}

	if deviceType == types.DeviceWeb {
		if sysInfo.UserLevel == types.LevelMember {
			return nil, result.NewError(result.UnknowDeviceType)
		}
	}

	_, err = orm.InsertLoginLog(sysInfo.UserId, deviceType, deviceName, loginType, uuid, version, utility.NowMillionSecond())
	if err != nil {
		logUser.Error("token login log user login failed", "err_msg", err)
	}
	return ret, nil
}

func Payment(appId, userId, logId, token, currency string, amount, fee float64, opp_address, rid, mode, payword, code string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.PaymentFailed).SetExtMessage("暂不支持付款服务")
	}
	//检查是否是公钥
	isPublicKey := utility.CheckPubkey(opp_address)
	recordFrom := ""
	recordTo := ""
	if logId == "" { //付款
		//提币
		var err error
		recordFrom, recordTo, err = account.BcoinWithdraw(app.AccountServer, token, currency, amount, fee, opp_address, rid, mode, payword, code, isPublicKey)
		if err != nil {
			return nil, result.NewError(result.PaymentFailed).SetExtMessage(err.Error())
		}
	} else { //收款
		//获取Log
		log, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		logInfo, err := GetChatLogAsUser(userId, log)
		if err != nil {
			return nil, err
		}

		if logInfo.MsgType != types.Receipt {
			return nil, result.NewError(result.PaymentFailed).SetExtMessage("这不是一条付款消息")
		}

		content := logInfo.GetCon()
		if content["coinName"] != currency {
			//TODO
		}

		if content["amount"] != amount-fee {
			//TODO
		}

		//提币
		recordFrom, recordTo, err = account.BcoinWithdraw(app.AccountServer, token, currency, amount, fee, opp_address, rid, mode, payword, code, isPublicKey)
		if err != nil {
			return nil, result.NewError(result.PaymentFailed).SetExtMessage(err.Error())
		}

		//修改支付信息
		content["recordId"] = recordFrom + "," + recordTo
		data := utility.StructToString(content)
		err = orm.UpdatePrivateLogContentById(logId, data)
		if err != nil {
			logUser.Error("UpdatePrivateLogContentById failed", "err", err)
		}
		//发消息通知
		target := log.SenderId
		operator := userId
		SendAlert(operator, target, types.ToUser, []string{operator, target}, types.Alert, proto.ComposePaymentAlert(operator, logId, recordFrom+","+recordTo))
	}
	ret := make(map[string]string)
	ret["recordId"] = recordFrom + "," + recordTo
	return ret, nil
}

//-------------------------------new end----------------------------//
// 返回：true 通过 false 不通过
func GetLastDeviceLoginInfo(userId, device string) (*types.LastLoginInfo, bool) {
	var info *types.LoginLog
	if device == types.DeviceAndroid || device == types.DeviceIOS {
		info, _ = orm.GetLastUserLoginLog(userId, []string{types.DeviceAndroid, types.DeviceIOS})
	} else {
		info, _ = orm.GetLastUserLoginLog(userId, []string{device})
	}
	if info == nil {
		return &types.LastLoginInfo{}, false
	}
	return &types.LastLoginInfo{
		Uuid:       info.Uuid,
		Device:     info.Device,
		DeviceName: info.DeviceName,
		LoginTime:  info.LoginTime,
		LoginType:  info.LoginType,
	}, true
}

//多端登录提示信息
func GetAlertMsg(info *types.LastLoginInfo) string {
	/*// TODO
	msg := "你的Chat33账号于%s在%s设备上通过%s登录。如果不是你的操作，你的%s已经泄露。"
	timeStr := utility.MillionSecondToTimeString(info.LoginTime)
	switch info.LoginType {
	case types.LoginTypeSms:
		return fmt.Sprintf(msg+"请勿转发验证码，并排查手机是否被植入木马导致短信被转发", timeStr, info.DeviceName, "短信验证码", "短信验证码")
	case types.LoginTypePwd:
		return fmt.Sprintf(msg+"请尽快登录Chat33修改密码", timeStr, info.DeviceName, "密码", "密码")
	case types.LoginTypeEmail:
		return fmt.Sprintf(msg+"请勿转发验证码", timeStr, info.DeviceName, "邮箱验证码", "邮箱验证码")
	case types.LoginTypeEmailPwd:
		return fmt.Sprintf(msg+"请尽快登录Chat33修改密码", timeStr, info.DeviceName, "密码", "密码")
	}
	return ""*/

	// TODO
	msg := "你的账号于%s在%s设备上通过%s登录。"
	timeStr := utility.MillionSecondToTimeString(info.LoginTime)
	switch info.LoginType {
	case types.LoginTypeSms:
		return fmt.Sprintf(msg, timeStr, info.DeviceName, "短信验证码")
	case types.LoginTypePwd:
		return fmt.Sprintf(msg, timeStr, info.DeviceName, "密码")
	case types.LoginTypeEmail:
		return fmt.Sprintf(msg, timeStr, info.DeviceName, "邮箱验证码")
	case types.LoginTypeEmailPwd:
		return fmt.Sprintf(msg, timeStr, info.DeviceName, "密码")
	}
	return ""
}

//多端登录提示信息v2
func GetAlertMsgV2(info *types.LastLoginInfo) string {
	maps := make(map[string]interface{})
	maps["device"] = info.DeviceName
	maps["time"] = info.LoginTime
	maps["way"] = info.LoginType

	b, err := json.Marshal(maps)
	if err != nil {
		return ""
	}
	return string(b)
}

//查看用户配置
func FindUserConf(appId, userId string) (map[string]interface{}, error) {
	needConfirmInvite := types.RoomInviteNotNeedConfirm
	//根据不同的app获取对应的 入群邀请配置
	if v, ok := types.DefaultInviteConfig[appId]; ok {
		needConfirmInvite = v
	}
	inviteConf, err := orm.RoomInviteConfirm(userId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if inviteConf != nil {
		needConfirmInvite = inviteConf.NeedConfirm
	}

	conf, err := GetAddFriendConfByUserId(userId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var confMap = make(map[string]interface{}, 16)
	if conf == nil {
		confMap["needConfirm"] = types.NeedConfirm
		confMap["needAnswer"] = types.NotNeedAnswer
		confMap["question"] = ""
		confMap["answer"] = ""
		confMap["needConfirmInvite"] = needConfirmInvite
	} else {
		confMap["needConfirm"] = utility.ToInt(conf.NeedConfirm)
		confMap["needAnswer"] = utility.ToInt(conf.NeedAnswer)
		confMap["question"] = conf.Question
		confMap["answer"] = conf.Answer
		confMap["needConfirmInvite"] = needConfirmInvite
	}
	return confMap, nil
}

//客户端获取邀请统计信息
func InviteStatistics(appId, token string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		data, err := account.InviteStatisticsFromBcoin(app.AccountServer, token)
		if err != nil {
			return nil, err
		}
		ret := make(map[string]interface{})
		info := data.(map[string]interface{})

		inviteNum := info["invite_num"]
		primary := make(map[string]interface{})
		primary["currency"] = app.MainCoin
		primary["total"] = "0"
		if sList, ok := info["statistics"].([]interface{}); ok {
			for i := 0; i < len(sList); i++ {
				field := sList[i].(map[string]interface{})
				if utility.ToString(field["currency"]) == app.MainCoin {
					primary["currency"] = field["currency"]
					primary["total"] = field["total"]
					sList = append(sList[:i], sList[i+1:]...)
					info["statistics"] = sList
					break
				}
			}
		}

		ret["invite_num"] = inviteNum
		ret["primary"] = primary
		ret["statistics"] = info["statistics"]
		return ret, nil
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

//客户端获取邀请信息
func SingleInviteInfos(appId, token string, page, size int64) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		return account.InviteInfosFromBcoin(app.AccountServer, token, 1, page, size)
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

//客户端获取邀请信息
func AccumulateInviteInfos(appId, token string, page, size int64) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		return account.InviteInfosFromBcoin(app.AccountServer, token, 2, page, size)
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

//检查是否已经认证
func CheckUserVerifyState(userId string) (int, error) {
	user, err := orm.GetUserInfoById(userId)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	if user.Identification == types.Verified {
		return types.VerifyStateAccept, nil
	} else {
		wait, err := orm.FindVerifyApplyByState(types.VerifyForUser, userId, types.VerifyStateWait)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if len(wait) > 0 {
			return types.VerifyStateWait, nil
		}

		reject, err := orm.FindVerifyApplyByState(types.VerifyForUser, userId, types.VerifyStateReject)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if len(reject) > 0 {
			rej := reject[0]
			if utility.NowMillionSecond() < utility.MillionSecondAddDuration(rej.UpdateTime, types.VerifyApplyInterval*time.Minute) {
				return types.VerifyStateReject, nil
			}
			return 0, nil
		}
	}
	return 0, nil
}

//设置入群邀请确认
func SetInviteConfirm(userId string, needConfirm int) error {
	err := orm.SetRoomInviteConfirm(userId, needConfirm)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}
