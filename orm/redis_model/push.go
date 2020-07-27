package redis_model

import (
	"github.com/33cn/chat33/cache"
	mysql "github.com/33cn/chat33/orm/mysql_model"
)

func UpdateDeviceToken(userId, deviceToken, deviceType string) error {
	err := mysql.UpdateDeviceToken(userId, deviceToken, deviceType)
	if err != nil {
		return err
	}

	err = cache.Cache.UpdateUserInfo(userId, "device_token", deviceToken)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.SaveDeviceToken(deviceType, deviceToken, userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func FindUserIdByDeviceToken(deviceToken string) (string, string, error) {
	deviceType, userId, err := cache.Cache.GetUserIdByDeviceToken(deviceToken)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	if userId == nil {
		deviceType, userId, err = mysql.FindUserIdByDeviceToken(deviceToken)
		if err != nil {
			return deviceType, "", err
		}
		if userId == nil {
			return deviceType, "", err
		}
		err = cache.Cache.SaveDeviceToken(deviceType, deviceToken, *userId)
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
		return deviceType, *userId, nil
	}
	return deviceType, *userId, err
}

func ClearDeviceToken(userId string) error {
	deviceToken, err := mysql.ClearUserDeviceToken(userId)
	if err != nil {
		return err
	}

	err = cache.Cache.ClearDeviceToken(userId, deviceToken)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.UpdateUserInfo(userId, "device_token", "")
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}
