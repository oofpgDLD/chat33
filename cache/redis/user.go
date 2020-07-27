package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	. "github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/garyburd/redigo/redis"
	"github.com/inconshreveable/log15"
)

//存储用户和用户的一些常用信息

var userLog = log15.New("cache", "cache user")

const (
	userKey        = "user-"
	userLoginKey   = "user-login"
	userInviteKey  = "user-invite-conf"
	deviceTokenKey = "device-token"
)

type userCaChe struct{}

//储存用户信息
func (c *userCaChe) SaveUserInfo(u *User) error {
	key := userKey + u.UserId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("HMSET", key,
		"user_id", u.UserId,
		"mark_id", u.MarkId,
		"uid", u.Uid,
		"app_id", u.AppId,
		"username", u.Username,
		"account", u.Account,
		"user_level", u.UserLevel,
		"verified  ", u.Verified, //是否实名制
		"avatar", u.Avatar,
		"area", u.Area,
		"phone", u.Phone,
		"sex", u.Sex,
		"email", u.Email,
		"invite_code", u.InviteCode,
		"device_token", u.DeviceToken,
		"deposit_address", u.DepositAddress,
		"device", u.Device,
		"create_time", u.CreateTime,
		"reg_version", u.RegVersion,
		"now_version", u.NowVersion,
		"close_until", u.CloseUntil,
		"super_user_level", u.SuperUserLevel,
		"public_key", u.PublicKey,
		"private_key", u.PrivateKey,
		"identification_info", u.IdentificationInfo,
		"identification", u.Identification,
		"device_token", u.DeviceToken,
	)
	if err != nil {
		userLog.Error("HMSET EXPIRE err", err)
		return err
	}
	//存储其他映射
	err = c.SaveUserInfoByField("uid", u.AppId, u.Uid, u.UserId)
	err = c.SaveUserInfoByField("phone", u.AppId, u.Phone, u.UserId)
	err = c.SaveUserInfoByField("mark_id", u.AppId, u.MarkId, u.UserId)
	err = c.SaveUserInfoByField("account", u.AppId, u.Account, u.UserId)
	//-------------//
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("SaveUserInfo EXPIRE err", err)
	}
	return err
}

//通过userid获取userinfo
func (*userCaChe) GetUserInfoById(userId string) (*User, error) {
	key := userKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("hgetall", key)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	v, err := redis.StringMap(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return nil, err
	}
	userInfo := &User{
		UserId:             v["user_id"],
		MarkId:             v["mark_id"],
		Uid:                v["uid"],
		AppId:              v["app_id"],
		Username:           v["username"],
		Account:            v["account"],
		UserLevel:          utility.ToInt(v["user_level"]),
		Verified:           utility.ToInt(v["verified"]),
		Avatar:             v["avatar"],
		Area:               v["area"],
		Phone:              v["phone"],
		Sex:                utility.ToInt(v["sex"]),
		Email:              v["email"],
		InviteCode:         v["invite_code"],
		DeviceToken:        v["device_token"],
		DepositAddress:     v["deposit_address"],
		PublicKey:          v["public_key"],
		PrivateKey:         v["private_key"],
		Device:             v["device"],
		CreateTime:         utility.ToInt64(v["create_time"]),
		RegVersion:         v["reg_version"],
		NowVersion:         v["now_version"],
		CloseUntil:         utility.ToInt64(v["close_until"]),
		SuperUserLevel:     utility.ToInt(v["super_user_level"]),
		Identification:     utility.ToInt(v["identification"]),
		IdentificationInfo: v["identification_info"],
	}
	return userInfo, err
}

func (*userCaChe) getUserField(userId, field string) (string, string, error) {
	key := userKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			friendLog.Error("Close err", err)
		}
	}()
	appId := ""
	value := ""
	{
		reply, err := conn.Do("hget", key, "app_id")
		if b, err := IsExists(reply, err); !b {
			return "", "", err
		}
		appId, err = redis.String(reply, err)
		if err != nil {
			userLog.Error("GetUserInfoById  Do hgetall err", err)
			return "", "", err
		}
	}

	{
		reply, err := conn.Do("hget", key, field)
		if b, err := IsExists(reply, err); !b {
			return "", "", err
		}
		value, err = redis.String(reply, err)
		if err != nil {
			userLog.Error("GetUserInfoById  Do hgetall err", err)
			return "", "", err
		}
	}

	return appId, value, nil
}

func (*userCaChe) delUserField(field, appId, value string) error {
	key := userKey + field
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	k := appId + "-" + value
	_, err = conn.Do("HDEL", key, k)
	return err
}

func (*userCaChe) getUserIdByField(field, appId, value string) (string, error) {
	key := userKey + field
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	k := appId + "-" + value
	reply, err := conn.Do("HGET", key, k)
	if b, err := IsExists(reply, err); !b {
		return "", err
	}
	id, err := redis.String(reply, err)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (*userCaChe) checkUserField(field string) (bool, error) {
	key := userKey + field
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("EXISTS", key)
	if b, err := IsExists(reply, err); !b {
		return false, err
	}
	v, err := redis.Bool(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return false, err
	}
	return v, err
}

//更新用户信息
//field userId uid username verified avatar getuiCid
func (u *userCaChe) UpdateUserInfo(userId, field, value string) error {
	//TODO 这里应该是需要事务的
	ok, err := u.checkUserField(field)
	if err != nil {
		return err
	}
	if ok {
		//获取该字段原值 和 appId
		appId, oldVal, err := u.getUserField(userId, field)
		if err != nil {
			return err
		}
		//查看原字段值是否属于该用户
		id, err := u.getUserIdByField(field, appId, oldVal)
		if err != nil {
			return err
		}
		if id == userId {
			//删除原值
			err = u.delUserField(field, appId, oldVal)
			if err != nil {
				return err
			}
		}
		//添加新值
		err = u.SaveUserInfoByField(field, appId, value, userId)
	}

	key := userKey + userId
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	_, err = conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	count, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		userLog.Error("UpdateUserInfo hexists err ", err)
		return err
	}
	if count == 0 {
		return nil
	}
	_, err = conn.Do("hset", key, field, value)
	return err
}

func (*userCaChe) SaveUserInfoByField(field, appId, value, id string) error {
	if field == "" || appId == "" || value == "" || id == "" {
		return errors.New("字段不可为空")
	}
	key := userKey + field
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	k := appId + "-" + value
	_, err = conn.Do("hset", key, k, id)
	return err
}

func (c *userCaChe) GetUserInfoByField(field, appId, value string) (*User, error) {
	id, err := c.getUserIdByField(field, appId, value)
	if err != nil {
		return nil, err
	}
	return c.GetUserInfoById(id)
}

func mutilDeviceKey(deviceType string) string {
	switch deviceType {
	case DeviceAndroid, DeviceIOS:
		return "mobile"
	default:
		return deviceType
	}
}

//保存用户登录信息
func (c *userCaChe) SaveUserLoginInfo(l *LoginLog) error {
	key := userLoginKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	k := mutilDeviceKey(l.Device) + "-" + l.UserId
	_, err = conn.Do("hset", key, k, string(b))
	return err
}

//获取用户登录信息
func (c *userCaChe) GetUserLoginInfo(userId, device string) (*LoginLog, error) {
	key := userLoginKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	k := mutilDeviceKey(device) + "-" + userId
	reply, err := conn.Do("HGET", key, k)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return nil, err
	}
	var l LoginLog
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

//保存拉群 是否需验证配置
func (c *userCaChe) SaveUserInviteRoomConf(conf *InviteRoomConf) error {
	key := userInviteKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	b, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	_, err = conn.Do("hset", key, conf.UserId, string(b))
	return err
}

//获取拉群 是否需验证配置
func (c *userCaChe) GetUserInviteRoomConf(userId string) (*InviteRoomConf, error) {
	key := userInviteKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGET", key, userId)
	if b, err := IsExists(reply, err); !b {
		return nil, err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return nil, err
	}
	var l InviteRoomConf
	err = json.Unmarshal([]byte(s), &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (c *userCaChe) SaveDeviceToken(deviceType, deviceToken, userId string) error {
	key := deviceTokenKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	_, err := conn.Do("EXPIRE", key, EXPIRETimeDay)
	if err != nil {
		userLog.Error("UpdateUserInfo EXPIRE err", err)
	}
	_, err = conn.Do("hset", key, deviceToken, fmt.Sprintf("%s:%s", deviceType, userId))
	return err
}

//获取deviceToken
func (c *userCaChe) GetUserIdByDeviceToken(deviceToken string) (string, *string, error) {
	key := deviceTokenKey
	conn := GetConnByKey(key)
	defer func() {
		err := conn.Close()
		if err != nil {
			userLog.Error("Close err", err)
		}
	}()
	reply, err := conn.Do("HGET", key, deviceToken)
	if b, err := IsExists(reply, err); !b {
		return "", nil, err
	}
	userId, err := redis.String(reply, err)
	if err != nil {
		userLog.Error("GetUserInfoById  Do hgetall err", err)
		return "", nil, err
	}
	ret := strings.Split(userId, ":")
	if len(ret) != 2 {
		return "", nil, errors.New("GetUserIdByDeviceToken split err")
	}
	return ret[0], &ret[1], nil
}

//获取拉群 是否需验证配置
func (c *userCaChe) ClearDeviceToken(userId, deviceToken string) error {
	deviceType, nowId, err := c.GetUserIdByDeviceToken(deviceToken)
	if err != nil {
		return err
	}
	if nowId == nil || *nowId != userId {
		return nil
	}

	return c.SaveDeviceToken(deviceType, deviceToken, "")
}
