package orm

import (
	mysql "github.com/33cn/chat33/orm/mysql_model"
	redis "github.com/33cn/chat33/orm/redis_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/user"
	user_model "github.com/33cn/chat33/user/model"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logUser = log15.New("model", "orm/user")

//根据appId获取应用的所有被封的用户
func GetCloseUserCountInApp(appId string) (int64, error) {
	ret, err := mysql.GetCloseUserCountInApp(appId)
	if err != nil {
		logUser.Error("mysql.GetCloseUserCountInApp", "err", err, "appId", appId)
	}
	return ret, err
}

//查找所有app下用户信息，包括被封用户，模糊查询 uid account
func GetUsersInAppQueryUid(appId, queryUid string) ([]*types.User, error) {
	ret, err := mysql.GetUsersInAppQueryUid(appId, queryUid)
	if err != nil {
		logUser.Error("mysql.GetUsersInAppQueryUid", "err", err, "appId", appId, "queryUid", queryUid)
	}
	return ret, err
}

//查找某个app下所有未封禁用户
func GetUsersInAppUnClose(appId string) ([]*types.User, error) {
	ret, err := mysql.GetUsersInAppUnClose(appId)
	if err != nil {
		logUser.Error("mysql.GetUsersInAppUnClose", "err", err, "appId", appId)
	}
	return ret, err
}

//查找某个app下所有封禁用户
func GetUsersInAppClosed(appId string) ([]*types.User, error) {
	ret, err := mysql.GetUsersInAppClosed(appId)
	if err != nil {
		logUser.Error("mysql.GetUsersInAppClosed", "err", err, "appId", appId)
	}
	return ret, err
}

//设置为加v认证用户
func SetUserVerifyed(tx types.Tx, userId, vInfo string) error {
	err := mysql.SetUserVerifyed(tx, userId, vInfo)
	if err != nil {
		logUser.Error("mysql.SetUserVerifyed", "err", err, "userId", userId, "verifyInfo", vInfo)
	}
	return err
}

func GetUserInfoByUid(appId, uid string) (*types.User, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetUserInfoByField(appId, "uid", uid)
		if err != nil {
			logUser.Error("redis.GetUserInfoByField", "err", err, "uid", uid)
		}
		return ret, err
	}
	ret, err := mysql.GetUserInfoByUid(appId, uid)
	if err != nil {
		logUser.Error("redis.GetUserInfoByUid", "err", err, "appId", appId, "uid", uid)
	}
	return ret, err
}

func GetUserInfoById(id string) (*types.User, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetUserInfoById(id)
		if err != nil {
			logUser.Error("redis.GetUserInfoById", "err", err, "userId", id)
		}
		return ret, err
	}
	ret, err := mysql.GetUserInfoById(id)
	if err != nil {
		logUser.Error("mysql.GetUserInfoById", "err", err, "userId", id)
	}
	return ret, err
}

func GetUserInfoByPhone(appId, phone string) (*types.User, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetUserInfoByField(appId, "phone", phone)
		if err != nil {
			logUser.Error("redis.GetUserInfoByField", "err", err, "appId", appId, "phone", phone)
		}
		return ret, err
	}
	ret, err := mysql.GetUserInfoByPhone(appId, phone)
	if err != nil {
		logUser.Error("mysql.GetUserInfoByPhone", "err", err, "appId", appId, "phone", phone)
	}
	return ret, err
}

func GetUserInfoByEmail(appId, email string) (*types.User, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetUserInfoByField(appId, "email", email)
		if err != nil {
			logUser.Error("redis.GetUserInfoByField", "err", err, "appId", appId, "email", email)
		}
		return ret, err
	}
	ret, err := mysql.GetUserInfoByEmail(appId, email)
	if err != nil {
		logUser.Error("mysql.GetUserInfoByEmail", "err", err, "appId", appId, "phone", email)
	}
	return ret, err
}

func GetUserInfoByMarkId(appId, markId string) (*types.User, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetUserInfoByField(appId, "uid", markId)
		if err != nil {
			logUser.Error("redis.GetUserInfoByField", "err", err, "appId", appId, "uid", markId)
		}
		return ret, err
	}
	ret, err := mysql.GetUserInfoByMarkId(appId, markId)
	if err != nil {
		logUser.Error("mysql.GetUserInfoByMarkId", "err", err, "appId", appId, "uid", markId)
	}
	return ret, err
}

func GetUserInfoByAccount(appId, account string) (*types.User, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetUserInfoByField(appId, "account", account)
		if err != nil {
			logUser.Error("redis.GetUserInfoByField", "err", err, "appId", appId, "account", account)
		}
		return ret, err
	}
	ret, err := mysql.GetUserInfoByAccount(appId, account)
	if err != nil {
		logUser.Error("mysql.GetUserInfoByAccount", "err", err, "appId", appId, "account", account)
	}
	return ret, err
}

func UpdateUid(markId, userId, uid string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateUid(markId, userId, uid)
		if err != nil {
			logUser.Error("redis.UpdateUid", "err", err, "userId", userId, "uid", uid)
		}
		return err
	}
	err := mysql.UpdateUid(markId, userId, uid)
	if err != nil {
		logUser.Error("mysql.UpdateUid", "err", err, "userId", userId, "uid", uid)
	}
	return err
}

func InsertUser(markId, uid, appId, username, account, email, area, phone, userLevel, verified, avatar, depositAddress, device, version string, createTime int64) (int64, error) {
	ret, err := mysql.InsertUser(markId, uid, appId, username, account, email, area, phone, userLevel, verified, avatar, depositAddress, device, version, createTime)
	if err != nil {
		logUser.Error("mysql.InsertUser", "err", err)
	}
	return ret, err
}

func AddUser(markId, uid, appId, username, phone, userLevel, device, version string, createTime int64) (int64, error) {
	ret, err := mysql.AddUser(markId, uid, appId, username, phone, userLevel, device, version, createTime)
	if err != nil {
		logUser.Error("mysql.InsertUser", "err", err)
	}
	return ret, err
}

func AddToken(userid, token string, createTime int64) (int64, error) {
	ret, err := mysql.AddToken(userid, token, createTime)
	if err != nil {
		logUser.Error("mysql.InsertToken", "err", err)
	}
	return ret, err
}

//查询根据userid查询Token
func FindToken(userid string) (string, int64, error) {
	res, err := mysql.FindToken(userid)
	var token string
	var time int64
	if res != nil {
		token = res["token"]
		time = utility.ToInt64(res["time"])
	}
	if err != nil {
		logUser.Error("mysql.GetToken", "err", err, "Token", token, "time", time)
	}
	return token, time, err
}

func InsertLoginLog(userID, deviceType, deviceName, loginType, uuid, version string, loginTime int64) (int64, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.InsertLoginLog(userID, deviceType, deviceName, loginType, uuid, version, loginTime)
		if err != nil {
			logUser.Error("redis.InsertLoginLog", "err", err)
		}
		return ret, err
	}
	ret, err := mysql.InsertLoginLog(userID, deviceType, deviceName, loginType, uuid, version, loginTime)
	if err != nil {
		logUser.Error("mysql.InsertLoginLog", "err", err)
	}
	return ret, err
}

func GetLastUserLoginLog(userId string, device []string) (*types.LoginLog, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetLoginLog(userId, device)
		if err != nil {
			logUser.Error("redis.GetLoginLog", "err", err, "userId", userId, "device", device)
		}
		return ret, err
	}
	ret, err := mysql.GetLastUserLoginLog(userId, device)
	if err != nil {
		logUser.Error("mysql.GetLastUserLoginLog", "err", err, "userId", userId, "device", device)
	}
	return ret, err
}

func UpdateDepositAddress(userId, address string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateDepositAddress(userId, address)
		if err != nil {
			logUser.Error("redis.UpdateDepositAddress", "err", err, "userId", userId, "address", address)
		}
		return err
	}
	err := mysql.UpdateDepositAddress(userId, address)
	if err != nil {
		logUser.Error("mysql.UpdateDepositAddress", "err", err, "userId", userId, "address", address)
	}
	return err
}

func UpdateUsername(userId, username string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateUsername(userId, username)
		if err != nil {
			logUser.Error("redis.UpdateUsername", "err", err, "userId", userId, "username", username)
		}
		return err
	}
	err := mysql.UpdateUsername(userId, username)
	if err != nil {
		logUser.Error("mysql.UpdateUsername", "err", err, "userId", userId, "username", username)
	}
	return err
}

func UpdateUserAvatar(userId, avatar string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateUserAvatar(userId, avatar)
		if err != nil {
			logUser.Error("redis.UpdateUserAvatar", "err", err, "userId", userId, "avatar", avatar)
		}
		return err
	}
	err := mysql.UpdateUserAvatar(userId, avatar)
	if err != nil {
		logUser.Error("mysql.UpdateUserAvatar", "err", err, "userId", userId, "avatar", avatar)
	}
	return err
}

func UpdatePhone(userId, phone string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdatePhone(userId, phone)
		if err != nil {
			logUser.Error("redis.UpdatePhone", "err", err, "userId", userId, "phone", phone)
		}
		return err
	}
	err := mysql.UpdatePhone(userId, phone)
	if err != nil {
		logUser.Error("mysql.UpdatePhone", "err", err, "userId", userId, "phone", phone)
	}
	return err
}

func UpdateEmail(userId, email string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateEmail(userId, email)
		if err != nil {
			logUser.Error("redis.UpdateEmail", "err", err, "userId", userId, "email", email)
		}
		return err
	}
	err := mysql.UpdateEmail(userId, email)
	if err != nil {
		logUser.Error("mysql.UpdateEmail", "err", err, "userId", userId, "email", email)
	}
	return err
}

func UpdateNowVersion(userId, version string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateNowVersion(userId, version)
		if err != nil {
			logUser.Error("redis.UpdateNowVersion", "err", err, "userId", userId, "version", version)
		}
		return err
	}
	err := mysql.UpdateNowVersion(userId, version)
	if err != nil {
		logUser.Error("mysql.UpdateNowVersion", "err", err, "userId", userId, "version", version)
	}
	return err
}

func UpdatePublicKey(userId, publicKey, privateKey string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdatePublicKey(userId, publicKey, privateKey)
		if err != nil {
			logUser.Error("redis.UpdatePublicKey", "err", err, "userId", userId, "publicKey", publicKey, "privateKey", privateKey)
		}
		return err
	}
	err := mysql.UpdatePublicKey(userId, publicKey, privateKey)
	if err != nil {
		logUser.Error("mysql.UpdatePublicKey", "err", err, "userId", userId, "publicKey", publicKey, "privateKey", privateKey)
	}
	return err
}

func UpdateInviteCode(userId, code string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateInviteCode(userId, code)
		if err != nil {
			logUser.Error("redis.UpdateInviteCode", "err", err, "userId", userId, "code", code)
		}
		return err
	}
	err := mysql.UpdateInviteCode(userId, code)
	if err != nil {
		logUser.Error("mysql.UpdateInviteCode", "err", err, "userId", userId, "code", code)
	}
	return err
}

func UpdateIsChain(userId string, ischain int) error {
	//if cfg.CacheType.Enable {
	//	err := redis.UpdateUsername(userId, username)
	//	if err != nil {
	//		logUser.Error("redis.UpdateUsername", "err", err, "userId", userId, "username", username)
	//	}
	//	return err
	//}
	err := mysql.UpdateIsChain(userId, ischain)
	if err != nil {
		logUser.Error("mysql.UpdateUsername", "err", err, "userId", userId, "ischain", ischain)
	}
	return err
}

//获取上链状态，直接查数据库
func GetIsChain(userId string) (int64, error) {
	ischain, err := mysql.GetIsChain(userId)
	if err != nil {
		logUser.Error("mysql.GetIsChain", "err", err, "userId", userId, "ischain", ischain)
	}
	return ischain, nil
}

func UpdateToken(userId, token string, createTime int64) error {

	err := mysql.UpdateToken(userId, token, createTime)
	if err != nil {
		logUser.Error("mysql.UpdateUsername", "err", err, "userId", userId, "token", token)
	}
	return err
}

//设置邀请入群是否需要验证
func SetRoomInviteConfirm(userId string, needConfirm int) error {
	if cfg.CacheType.Enable {
		err := redis.SetRoomInviteConfirm(userId, needConfirm)
		if err != nil {
			logUser.Error("redis.SetRoomInviteConfirm", "err", err, "userId", userId, "needConfirm", needConfirm)
		}
		return err
	}
	err := mysql.SetRoomInviteConfirm(userId, needConfirm)
	if err != nil {
		logUser.Error("mysql.SetRoomInviteConfirm", "err", err, "userId", userId, "needConfirm", needConfirm)
	}
	return err
}

func RoomInviteConfirm(userId string) (*types.InviteRoomConf, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.RoomInviteConfirm(userId)
		if err != nil {
			logUser.Error("redis.RoomInviteConfirm", "err", err, "userId", userId)
		}
		return ret, err
	}
	ret, err := mysql.RoomInviteConfirm(userId)
	if err != nil {
		logUser.Error("mysql.RoomInviteConfirm", "err", err, "userId", userId)
	}
	return ret, err
}

//从托管取数据(弃用)
func GetToken(appId, token string) (*user_model.Token, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetToken(appId, token)
		if err != nil {
			logUser.Error("redis.GetToken", "err", err, "appId", appId, "token", token)
		}
		return ret, err
	}
	//从接口获取并更新缓存
	user, err := user.GetUserInfoFromAccountSystem(appId, token)
	if user == nil {
		return nil, err
	}
	return &user_model.Token{
		Uid:   user.Uid,
		AppId: appId,
		Token: token,
		Time:  utility.NowMillionSecond(),
	}, nil
}

//通过Token查询用户信息
func GetUserInfoByToken(appId, token string) (*types.User, error) {

	ret, err := mysql.GetUserInfoByToken(appId, token)
	if err != nil {
		logUser.Error("mysql.GetUserInfoByToken", "err", err, "appId", appId, "Token", token)
	}
	return ret, err
}

func SaveToken(tokenInfo *user_model.Token) error {
	if cfg.CacheType.Enable {
		err := redis.SaveToken(tokenInfo)
		if err != nil {
			logUser.Error("redis.GetToken", "tokenInfo", tokenInfo)
		}
		return err
	}
	return nil
}

//获取所有userId
func GetUsers() ([]string, error) {
	ret, err := mysql.GetUsers()
	if err != nil {
		logUser.Error("mysql.GetUsers")
	}
	return ret, err
}
