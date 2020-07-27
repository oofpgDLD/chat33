package redis_model

import (
	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/cache"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/user"
	user_model "github.com/33cn/chat33/user/model"
)

var UserFieldMap = map[string]func(string, string) (*types.User, error){
	"uid":     mysql.GetUserInfoByUid,
	"phone":   mysql.GetUserInfoByPhone,
	"mark_id": mysql.GetUserInfoByMarkId,
	"account": mysql.GetUserInfoByAccount,
}

//根据id获取用户信息
func GetUserInfoById(id string) (*types.User, error) {
	userInfo, err := cache.Cache.GetUserInfoById(id)
	if err != nil {
		l.Warn("find user info from redis", "err", err)
	}
	if userInfo == nil {
		//从mysql 获取并更新缓存
		user, err := mysql.GetUserInfoById(id)
		if user == nil {
			return nil, err
		}
		err = cache.Cache.SaveUserInfo(user)
		if err != nil {
			l.Warn("redis can not save user", "err", err)
		}
		return user, nil
	}
	return userInfo, nil
}

//根据id获取用户信息
func GetUserInfoByField(appId, field, value string) (*types.User, error) {
	userInfo, err := cache.Cache.GetUserInfoByField(field, appId, value)
	if err != nil {
		l.Warn("find user info by field from redis failed", "err", err)
	}
	if userInfo == nil {
		//从mysql 获取并更新缓存
		f := UserFieldMap[field]
		user, err := f(appId, value)
		if user == nil {
			l.Warn("user info field can not find", "err", err)
			return nil, err
		}
		err = cache.Cache.SaveUserInfo(user)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return user, nil
	}
	return userInfo, nil
}

func UpdateUid(markId, userId, uid string) error {

	err := mysql.UpdateUid(markId, userId, uid)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "uid", uid)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	err = cache.Cache.UpdateUserInfo(userId, "markId", markId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//更新账户地址
func UpdateDepositAddress(userId, address string) error {
	err := mysql.UpdateDepositAddress(userId, address)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "deposit_address", address)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//更新用户名
func UpdateUsername(userId, username string) error {
	err := mysql.UpdateUsername(userId, username)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "username", username)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//更新用户头像
func UpdateUserAvatar(userId, avatar string) error {
	err := mysql.UpdateUserAvatar(userId, avatar)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "avatar", avatar)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//更新用户手机号
func UpdatePhone(userId, phone string) error {
	err := mysql.UpdatePhone(userId, phone)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "phone", phone)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//更新用户邮箱
func UpdateEmail(userId, email string) error {
	err := mysql.UpdateEmail(userId, email)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "email", email)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func UpdateNowVersion(userId, version string) error {
	err := mysql.UpdateNowVersion(userId, version)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "now_version", version)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func UpdatePublicKey(userId, publicKey, privateKey string) error {
	err := mysql.UpdatePublicKey(userId, publicKey, privateKey)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "public_key", publicKey)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
		return err
	}

	err = cache.Cache.UpdateUserInfo(userId, "private_key", privateKey)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
		return err
	}
	return nil
}

func UpdateInviteCode(userId, code string) error {
	err := mysql.UpdateInviteCode(userId, code)
	if err != nil {
		return err
	}
	err = cache.Cache.UpdateUserInfo(userId, "invite_code", code)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func InsertLoginLog(userID, deviceType, deviceName, loginType, uuid, version string, loginTime int64) (int64, error) {
	id, err := mysql.InsertLoginLog(userID, deviceType, deviceName, loginType, uuid, version, loginTime)
	if err != nil {
		return 0, err
	}
	err = cache.Cache.SaveUserLoginInfo(&types.LoginLog{
		Id:         utility.ToString(id),
		UserId:     userID,
		DeviceName: deviceName,
		Device:     deviceType,
		LoginType:  utility.ToInt(loginType),
		Uuid:       uuid,
		Version:    version,
		LoginTime:  loginTime,
	})
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return id, nil
}

func GetLoginLog(userId string, devices []string) (*types.LoginLog, error) {
	userInfo, err := cache.Cache.GetUserLoginInfo(userId, devices[0])
	if err != nil {
		l.Warn("find login info from redis", "err", err)
	}
	if userInfo == nil {
		//从mysql 获取并更新缓存
		user, err := mysql.GetLastUserLoginLog(userId, devices)
		if user == nil {
			l.Warn("login log can not find")
			return nil, err
		}
		err = cache.Cache.SaveUserLoginInfo(user)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return user, nil
	}
	return userInfo, nil
}

func SetRoomInviteConfirm(userId string, needConfirm int) error {
	err := mysql.SetRoomInviteConfirm(userId, needConfirm)
	if err != nil {
		return err
	}
	err = cache.Cache.SaveUserInviteRoomConf(&types.InviteRoomConf{
		UserId:      userId,
		NeedConfirm: needConfirm,
	})
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func RoomInviteConfirm(userId string) (*types.InviteRoomConf, error) {
	conf, err := cache.Cache.GetUserInviteRoomConf(userId)
	if err != nil {
		l.Warn("find user invite room confirm from redis", "err", err)
	}
	if conf == nil {
		//从mysql 获取并更新缓存
		conf, err = mysql.RoomInviteConfirm(userId)
		if conf == nil {
			l.Warn("room invite confirm can not find")
			return nil, err
		}
		err = cache.Cache.SaveUserInviteRoomConf(conf)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return conf, nil
	}
	return conf, nil
}

//通过token获取用户
func GetToken(appId, token string) (*user_model.Token, error) {
	tokenInfo, err := cache.Cache.GetToken(appId, token)
	if err != nil {
		l.Warn("find user info from redis", "err", err)
	}
	if tokenInfo == nil {
		//从接口获取并更新缓存
		user, err := user.GetUserInfoFromAccountSystem(appId, token)
		if user == nil {
			l.Warn("user can not find")
			return nil, err
		}
		tokenInfo := &user_model.Token{
			Uid:   user.Uid,
			AppId: appId,
			Token: token,
			Time:  utility.NowMillionSecond(),
		}
		err = cache.Cache.SaveToken(tokenInfo)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return tokenInfo, nil
	}
	return tokenInfo, nil
}

//通过token获取用户
func SaveToken(tokenInfo *user_model.Token) error {
	return cache.Cache.SaveToken(tokenInfo)
}
