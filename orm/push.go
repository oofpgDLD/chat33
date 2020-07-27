package orm

import (
	mysql "github.com/33cn/chat33/orm/mysql_model"
	redis "github.com/33cn/chat33/orm/redis_model"
	"github.com/inconshreveable/log15"
)

var logPush = log15.New("module", "model/push")

//deviceType userId err
func FindUserIdByDeviceToken(deviceToken string) (string, string, error) {
	if cfg.CacheType.Enable {
		ret, ret2, err := redis.FindUserIdByDeviceToken(deviceToken)
		if err != nil {
			logPush.Error("redis.FindUserIdByDeviceToken", "err", err, "deviceToken", deviceToken)
		}
		return ret, ret2, err
	}
	deviceType, userId, err := mysql.FindUserIdByDeviceToken(deviceToken)
	if userId == nil {
		if err != nil {
			logPush.Error("mysql.FindUserIdByDeviceToken", "err", err, "deviceToken", deviceToken)
		}
		return "", "", err
	}
	return deviceType, *userId, err
}

func SetUserDeviceToken(userId, deviceToken, deviceType string) error {
	if cfg.CacheType.Enable {
		err := redis.UpdateDeviceToken(userId, deviceToken, deviceType)
		if err != nil {
			logPush.Error("redis.UpdateDeviceToken", "err", err, "deviceToken", deviceToken, "deviceType", deviceType)
		}
		return err
	}
	err := mysql.UpdateDeviceToken(userId, deviceToken, deviceType)
	if err != nil {
		logPush.Error("mysql.UpdateDeviceToken", "err", err, "deviceToken", deviceToken, "deviceType", deviceType)
	}
	return err
}

func ClearUserDeviceToken(userId string) error {
	if cfg.CacheType.Enable {
		err := redis.ClearDeviceToken(userId)
		if err != nil {
			logPush.Error("redis.ClearDeviceToken", "err", err, "userId", userId)
		}
		return err
	}
	_, err := mysql.ClearUserDeviceToken(userId)
	if err != nil {
		logPush.Error("mysql.ClearUserDeviceToken", "err", err, "userId", userId)
	}
	return err
}
