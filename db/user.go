package db

import (
	"fmt"
	"strings"

	"github.com/33cn/chat33/pkg/btrade/common/mysql"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func GetUserInfoById(userId string) ([]map[string]string, error) {
	const sqlStr = "SELECT * from `user` where user_id= ?"
	return conn.Query(sqlStr, userId)
}

func UpdateUserAvatar(id, avatar string) error {
	const sqlStr = "update `user` set avatar = ? where user_id = ?"
	_, _, err := conn.Exec(sqlStr, avatar, id)
	return err
}

func UpdatePhone(id, phone string) error {
	const sqlStr = "update `user` set phone = ? where user_id = ?"
	_, _, err := conn.Exec(sqlStr, phone, id)
	return err
}

func UpdateEmail(id, email string) error {
	const sqlStr = "update `user` set email = ? where user_id = ?"
	_, _, err := conn.Exec(sqlStr, email, id)
	return err
}

// GetUserInfo 获取用户信息 根据uid
func GetUserInfoByUID(appId, uid string) ([]map[string]string, error) {
	const sqlStr = `select * from user where app_id = ? and uid = ? limit 1`
	return conn.Query(sqlStr, appId, uid)
}

// GetUserInfo 获取用户信息 根据account
func GetUserInfoByAccount(appId, account string) ([]map[string]string, error) {
	const sqlStr = `select * from user where app_id = ? and account = ? limit 1`
	return conn.Query(sqlStr, appId, account)
}

// GetUserInfo 获取用户信息 根据device_token
func FindUserByDeviceToken(appId, deviceToken string) ([]map[string]string, error) {
	const sqlStr = `select * from user where app_id = ? and device_token = ? limit 1`
	return conn.Query(sqlStr, appId, deviceToken)
}

// insert user info
func InsertUser(markId, uid, appId, userName, account, email, area, phone, userLevel, verified, avatar, depositAddress, device, version string, createTime int64) (num int64, userId int64, err error) {
	if len(avatar) == 0 {
		const sqlStr = `INSERT IGNORE INTO user(mark_id,uid,app_id,username,account,email,area,phone,user_level,verified,deposit_address,device,create_time,reg_version,now_version) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		num, userId, err = conn.Exec(sqlStr, markId, uid, appId, userName, account, email, area, phone, userLevel, verified, depositAddress, device, createTime, version, version)
	} else {
		const sqlStr = `INSERT IGNORE INTO user(mark_id,uid,app_id,username,account,email,area,phone,user_level,verified,avatar,deposit_address,device,create_time,reg_version,now_version) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		num, userId, err = conn.Exec(sqlStr, markId, uid, appId, userName, account, email, area, phone, userLevel, verified, avatar, depositAddress, device, createTime, version, version)
	}
	return
}

// insert user info
func AddUser(markId, uid, appId, username, phone, userLevel, device, version string, createTime int64) (num int64, userId int64, err error) {
	const sqlStr = `INSERT IGNORE INTO user(mark_id,uid,app_id,username,phone,user_level,device,create_time,reg_version,now_version) VALUES(?,?,?,?,?,?,?,?,?,?)`
	num, userId, err = conn.Exec(sqlStr, markId, uid, appId, username, phone, userLevel, device, createTime, version, version)
	return
}
func AddToken(userid, token string, createTime int64) (int64, error) {
	const sqlStr = `INSERT IGNORE INTO token(user_id,token,time) VALUES(?,?,?) on duplicate key update token = ?, time = ?`
	_, id, err := conn.Exec(sqlStr, userid, token, createTime, token, createTime)
	return id, err
}

func FindToken(userid string) ([]map[string]string, error) {
	const sqlStr = "select * from `token` where user_id = ?"
	return conn.Query(sqlStr, userid)
}

func UpdateToken(userId, token string, createTime int64) error {
	const sqlStr = "UPDATE token SET token = ?,time = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sqlStr, token, createTime, userId)
	return err
}

func AddUserLoginLog(userID, deviceType, deviceName, loginType, uuid, version string, loginTime int64) (int64, error) {
	if deviceType == "" {
		deviceType = "Unknown"
	}
	if deviceName == "" {
		deviceName = "Unknown"
	}
	const sqlStr = "insert into login_log (user_id, device, device_name, login_type, login_time, uuid, version) values(?,?,?,?,?,?,?)"
	_, id, err := conn.Exec(sqlStr, userID, deviceType, deviceName, loginType, loginTime, uuid, version)
	return id, err
}

func GetLastUserLoginLog(userID string, deviceType []string) ([]map[string]string, error) {
	caseType := ""
	if len(deviceType) > 0 {
		caseType = fmt.Sprintf("and `device` in(%s)", fmt.Sprintf("'%s'", strings.Join(deviceType, "','")))
	}
	sqlStr := fmt.Sprintf("select * from login_log where user_id = ? %s order by login_time desc limit 0,1", caseType)
	return conn.Query(sqlStr, userID)
}

func UpdateDeviceToken(userId, deviceToken, deviceType string) error {
	tx, err := conn.NewTx()
	if err != nil {
		return err
	}

	sql := "update user set device_token = ? where user_id = ?"
	_, _, err = tx.Exec(sql, deviceToken, userId)
	if err != nil {
		tx.RollBack()
		return err
	}
	sql = "insert into `push`(device_token,user_id,device_type) values(?,?,?) ON DUPLICATE KEY UPDATE user_id = ?,device_type = ?"
	_, _, err = tx.Exec(sql, deviceToken, userId, deviceType, userId, deviceType)
	if err != nil {
		tx.RollBack()
		return err
	}
	return tx.Commit()
}

//return deviceToken and error
func ClearUserDeviceToken(userId string) (string, error) {
	sql := "select * from `user` where user_id =?"
	maps, err := conn.Query(sql, userId)
	if err != nil {
		return "", err
	}
	if len(maps) < 1 {
		return "", nil
	}
	deviceToken := maps[0]["device_token"]

	sql = "update user set device_token = ? where user_id = ?"
	_, _, err = conn.Exec(sql, "", userId)
	if err != nil {
		return "", err
	}

	sql = "select * from push where device_token = ?"
	maps, err = conn.Query(sql, deviceToken)
	if err != nil {
		return "", err
	}
	if len(maps) < 1 {
		return deviceToken, nil
	}
	nowUser := maps[0]["user_id"]

	if nowUser == userId {
		sql = "UPDATE `push` SET user_id = ? WHERE device_token = ? "
		_, _, err = conn.Exec(sql, "", deviceToken)
		if err != nil {
			return deviceToken, err
		}
	}
	return "", nil
}

func FindUserIdByDeviceToken(deviceToken string) (string, *string, error) {
	sql := "select * from `push` where device_token = ?"
	maps, err := conn.Query(sql, deviceToken)
	if err != nil {
		return "", nil, err
	}
	if len(maps) <= 0 {
		return "", nil, nil
	}
	userId := maps[0]["user_id"]
	deviceType := maps[0]["device_type"]
	return deviceType, &userId, nil
}

func UpdateCodeByUid(uid, code string) error {
	sql := "UPDATE user SET invite_code  = ? WHERE uid = ?"
	_, _, err := conn.Exec(sql, code, uid)
	return err
}

func UpdateUid(markId, userId, uid string) error {
	sql := "UPDATE user SET uid = ?,mark_id = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sql, uid, markId, userId)
	return err
}

func UpdatePublicKey(userId, publicKey, privateKey string) error {
	const sql = "UPDATE user SET public_key = ?, private_key = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sql, publicKey, privateKey, userId)
	return err
}

func UpdateInviteCode(userId, code string) error {
	const sql = "UPDATE user SET invite_code = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sql, code, userId)
	return err
}

func UpdateDepositAddress(userId, address string) error {
	const sqlStr = "UPDATE user SET deposit_address = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sqlStr, address, userId)
	return err
}

func UpdateUsername(userId, username string) error {
	const sqlStr = "update `user` set username=? where user_id=?"
	_, _, err := conn.Exec(sqlStr, username, userId)
	return err
}

func UpdateNowVersion(userId, version string) error {
	const sqlStr = "UPDATE user SET now_version = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sqlStr, version, userId)
	return err
}

func UpdateIsChain(userId string, ischain int) error {
	const sqlStr = "UPDATE user SET ischain = ? WHERE user_id = ?"
	_, _, err := conn.Exec(sqlStr, ischain, userId)
	return err
}

func GetIsChain(userId string) (int64, error) {
	const sqlStr = "SELECT ischain FROM `user` where user_id = ?"
	maps, err := conn.Query(sqlStr, userId)
	if err != nil || len(maps) == 0 {
		return -1, err
	}
	return utility.ToInt64(maps[0]["ischain"]), nil
}

func GetCloseUserCountInApp(appId string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and close_until > ?"
	maps, err := conn.Query(sqlStr, appId, utility.NowMillionSecond())
	if err != nil || len(maps) == 0 {
		return 0, err
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

//查找所有app下用户信息，包括被封用户，模糊查询 uid account
func GetUsersInAppQueryUid(appId, query string) ([]map[string]string, error) {
	sqlStr := ""
	if query != "" {
		query = QueryStr(query)
		//sqlStr = "SELECT * from `user` where app_id = ? and (uid LIKE '%" + query + "%' or account LIKE '%" + query + "%') order by create_time desc"
		sqlStr = fmt.Sprintf("SELECT * from `user` where app_id = ? and (uid LIKE '%%%s%%' or account LIKE '%%%s%%') order by create_time desc", query, query)
	} else {
		sqlStr = "SELECT * from `user` where app_id = ? order by create_time desc"
	}
	return conn.Query(sqlStr, appId)
}

//查找某个app下所有未封禁用户
func GetUsersInAppUnClose(appId string) ([]map[string]string, error) {
	sqlStr := "SELECT * from `user` where app_id = ? and close_until <= ? order by create_time desc"
	return conn.Query(sqlStr, appId, utility.NowMillionSecond())
}

//查找某个app下所有封禁用户
func GetUsersInAppClosed(appId string) ([]map[string]string, error) {
	sqlStr := "SELECT * from `user` where app_id = ? and close_until > ? order by create_time desc"
	return conn.Query(sqlStr, appId, utility.NowMillionSecond())
}

//
func RoomInviteConfirm(userId string) ([]map[string]string, error) {
	const sqlStr = "select * from `invite_room_conf` where user_id = ?"
	return conn.Query(sqlStr, userId)
}

//设置邀请入群是否需要验证
func SetRoomInviteConfirm(userId string, needConfirm int) (int64, int64, error) {
	const sqlStr = "insert into `invite_room_conf`(user_id,need_confirm) values(?,?) ON DUPLICATE KEY UPDATE need_confirm = ?"
	return conn.Exec(sqlStr, userId, needConfirm, needConfirm)
}

//设置为认证用户
func SetUserVerifyed(tx *mysql.MysqlTx, userId, identificationInfo string) (int64, int64, error) {
	const sqlStr = "UPDATE `user` SET identification = ?, identification_info = ? WHERE user_id = ?"
	return tx.Exec(sqlStr, types.Verified, identificationInfo, userId)
}

//获取所有UserId
func GetUsers() ([]map[string]string, error) {
	const sqlStr = "SELECT user_id FROM `user`"
	return conn.Query(sqlStr)
}
