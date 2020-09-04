package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func UpdateDeviceToken(userId, deviceToken, deviceType string) error {
	return db.UpdateDeviceToken(userId, deviceToken, deviceType)
}

func ClearUserDeviceToken(userId string) (string, error) {
	return db.ClearUserDeviceToken(userId)
}

func FindUserIdByDeviceToken(deviceToken string) (string, *string, error) {
	return db.FindUserIdByDeviceToken(deviceToken)
}

func UpdateUid(markId, userId, uid string) error {
	return db.UpdateUid(markId, userId, uid)
}

func UpdatePublicKey(userId, publicKey, privateKey string) error {
	return db.UpdatePublicKey(userId, publicKey, privateKey)
}

func UpdateInviteCode(userId, code string) error {
	return db.UpdateInviteCode(userId, code)
}

func SetUserVerifyed(tx types.Tx, userId, vInfo string) error {
	_, _, err := db.SetUserVerifyed(tx.(*mysql.MysqlTx), userId, vInfo)
	return err
}

//设置邀请入群是否需要验证
func SetRoomInviteConfirm(userId string, needConfirm int) error {
	_, _, err := db.SetRoomInviteConfirm(userId, needConfirm)
	return err
}

//设置邀请入群是否需要验证
func RoomInviteConfirm(userId string) (*types.InviteRoomConf, error) {
	maps, err := db.RoomInviteConfirm(userId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	m := maps[0]
	return &types.InviteRoomConf{
		UserId:      m["user_id"],
		NeedConfirm: utility.ToInt(m["need_confirm"]),
	}, nil
}

func InsertUser(markId, uid, appId, username, account, email, area, phone, userLevel, verified, avatar, depositAddress, device, version string, createTime int64) (int64, error) {
	_, id, err := db.InsertUser(markId, uid, appId, username, account, email, area, phone, userLevel, verified, avatar, depositAddress, device, version, createTime)
	return id, err
}

func AddUser(markId, uid, appId, username, phone, userLevel, device, version string, createTime int64) (int64, error) {
	_, id, err := db.AddUser(markId, uid, appId, username, phone, userLevel, device, version, createTime)
	return id, err
}
func AddToken(userid, token string, createTime int64) (int64, error) {
	return db.AddToken(userid, token, createTime)

}

func FindToken(userid string) (map[string]string, error) {
	maps, err := db.FindToken(userid)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return info, nil
}

func InsertLoginLog(userID, deviceType, deviceName, loginType, uuid, version string, loginTime int64) (int64, error) {
	return db.AddUserLoginLog(userID, deviceType, deviceName, loginType, uuid, version, loginTime)
}

func UpdateDepositAddress(userId, address string) error {
	return db.UpdateDepositAddress(userId, address)
}

func UpdateUsername(userId, username string) error {
	return db.UpdateUsername(userId, username)
}

func UpdateUserAvatar(userId, avatar string) error {
	return db.UpdateUserAvatar(userId, avatar)
}

func UpdatePhone(userId, phone string) error {
	return db.UpdatePhone(userId, phone)
}

func UpdateEmail(userId, email string) error {
	return db.UpdateEmail(userId, email)
}

func UpdateNowVersion(userId, version string) error {
	return db.UpdateNowVersion(userId, version)
}

func UpdateIsChain(userId string, ischain int) error {
	return db.UpdateIsChain(userId, ischain)
}

func GetIsChain(userId string) (int64, error) {
	return db.GetIsChain(userId)
}

func UpdateToken(userId, token string, createTime int64) error {
	return db.UpdateToken(userId, token, createTime)
}

func GetCloseUserCountInApp(appId string) (int64, error) {
	return db.GetCloseUserCountInApp(appId)
}

func convertUser(v map[string]string) *types.User {
	return &types.User{
		UserId:             v["user_id"],
		MarkId:             v["mark_id"],
		Uid:                v["uid"],
		AppId:              v["app_id"],
		Username:           v["username"],
		Account:            v["account"],
		UserLevel:          utility.ToInt(v["user_level"]),
		Verified:           utility.ToInt(v["verified"]),
		Avatar:             v["avatar"],
		CompanyId:          v["com_id"],
		Position:           v["position"],
		Sex:                utility.ToInt(v["sex"]),
		Phone:              v["phone"],
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
		IsChain:            utility.ToInt(v["ischain"]),
	}
}

func convertJoinUser(v map[string]string) *types.User {
	return &types.User{
		UserId:             v["U_user_id"],
		MarkId:             v["mark_id"],
		Uid:                v["uid"],
		AppId:              v["app_id"],
		Username:           v["username"],
		Account:            v["account"],
		PublicKey:          v["public_key"],
		PrivateKey:         v["private_key"],
		UserLevel:          utility.ToInt(v["user_level"]),
		Verified:           utility.ToInt(v["verified"]),
		Avatar:             v["avatar"],
		CompanyId:          v["com_id"],
		Position:           v["position"],
		Sex:                utility.ToInt(v["sex"]),
		Phone:              v["phone"],
		Email:              v["email"],
		InviteCode:         v["invite_code"],
		DeviceToken:        v["device_token"],
		DepositAddress:     v["deposit_address"],
		Device:             v["device"],
		CreateTime:         utility.ToInt64(v["create_time"]),
		RegVersion:         v["reg_version"],
		NowVersion:         v["now_version"],
		CloseUntil:         utility.ToInt64(v["close_until"]),
		SuperUserLevel:     utility.ToInt(v["super_user_level"]),
		Identification:     utility.ToInt(v["identification"]),
		IdentificationInfo: v["identification_info"],
	}
}

//查找所有app下用户信息，包括被封用户，模糊查询 uid account
func GetUsersInAppQueryUid(appId, query string) ([]*types.User, error) {
	maps, err := db.GetUsersInAppQueryUid(appId, query)
	if err != nil {
		return nil, err
	}
	list := make([]*types.User, 0)

	for _, info := range maps {
		item := convertUser(info)
		list = append(list, item)
	}
	return list, nil
}

//查找某个app下所有未封禁用户
func GetUsersInAppUnClose(appId string) ([]*types.User, error) {
	maps, err := db.GetUsersInAppUnClose(appId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.User, 0)

	for _, info := range maps {
		item := convertUser(info)
		list = append(list, item)
	}
	return list, nil
}

//查找某个app下所有封禁用户
func GetUsersInAppClosed(appId string) ([]*types.User, error) {
	maps, err := db.GetUsersInAppClosed(appId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.User, 0)

	for _, info := range maps {
		item := convertUser(info)
		list = append(list, item)
	}
	return list, nil
}

func GetUserInfoByUid(appId, uid string) (*types.User, error) {
	maps, err := db.GetUserInfoByUID(appId, uid)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]

	return convertUser(info), nil
}

func GetUserInfoById(id string) (*types.User, error) {
	maps, err := db.GetUserInfoById(id)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertUser(info), nil
}

func GetUserInfoByPhone(appId, phone string) (*types.User, error) {
	maps, err := db.FindUserByPhoneV2(appId, phone)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertUser(info), nil
}

func GetUserInfoByEmail(appId, email string) (*types.User, error) {
	maps, err := db.FindUserByEmail(appId, email)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertUser(info), nil
}

func GetUserInfoByToken(appId, token string) (*types.User, error) {
	maps, err := db.FindUserByToken(appId, token)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertUser(info), nil
}

func GetUserInfoByMarkId(appId, markId string) (*types.User, error) {
	maps, err := db.FindUserByMarkId(appId, markId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertUser(info), nil
}

func GetUserInfoByAccount(appId, account string) (*types.User, error) {
	maps, err := db.GetUserInfoByAccount(appId, account)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]

	return convertUser(info), nil
}

func GetLastUserLoginLog(userId string, device []string) (*types.LoginLog, error) {
	maps, err := db.GetLastUserLoginLog(userId, device)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return &types.LoginLog{
		Id:         info["id"],
		UserId:     info["user_id"],
		Uuid:       info["uuid"],
		Device:     info["device"],
		DeviceName: info["device_name"],
		LoginTime:  utility.ToInt64(info["login_time"]),
		LoginType:  utility.ToInt(info["login_type"]),
		Version:    info["version"],
	}, nil
}

func GetUsers() ([]string, error) {
	maps, err := db.GetUsers()
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	var str []string

	for _, v := range maps {
		str = append(str, v["user_id"])
	}

	return str, nil
}
