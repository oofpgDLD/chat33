package db

import (
	"fmt"

	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

//添加埋点
func InsertOpenLog(userId, appId, device, version string, createTime int64) error {
	const sqlStr = "insert into `open_log`(user_id,app_id,device,version,create_time) values(?,?,?,?,?)"
	_, _, err := conn.Exec(sqlStr, userId, appId, device, version, createTime)
	return err
}

// GetUserInfo 获取用户信息 根据account
func AdminCheckLogin(appId, account string) ([]map[string]string, error) {
	const sqlStr = "select * from `admin` where app_id = ? and account = ?"
	return conn.Query(sqlStr, appId, account)
}

// 获取管理员信息 根据id
func FindAdminById(id string) ([]map[string]string, error) {
	const sqlStr = "select * from `admin` where id = ?"
	return conn.Query(sqlStr, id)
}

func GetVersionList(appId string) ([]map[string]string, error) {
	const sqlStr = `select now_version from user where app_id = ? group by now_version`
	return conn.Query(sqlStr, appId)
}

func UpdateVersionUsersNumber(appId, nowVersion string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and reg_version != ? and now_version = ?"
	maps, err := conn.Query(sqlStr, appId, nowVersion, nowVersion)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//新增用户
func IncreaseUsers(appId string, startTime, endTime int64) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and create_time >= ? and create_time <= ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//新增用户 根据平台查询
func IncreaseUsersAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and create_time >= ? and create_time <= ? and device = ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime, device)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//新增用户 根据版本号
func IncreaseUsersAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and create_time >= ? and create_time <= ? and reg_version = ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime, version)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//新增用户信息 根据版本号
func IncreaseUsersInfoAsVersion(appId string, startTime, endTime int64, version string) ([]map[string]string, error) {
	const sqlStr = "select * from user where app_id = ? and create_time >= ? and create_time <= ? and reg_version = ?"
	return conn.Query(sqlStr, appId, startTime, endTime, version)
}

//新增用户 根据版本号和平台
func IncreaseUsersAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and create_time >= ? and create_time <= ? and reg_version = ? and device = ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime, version, platform)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//新增用户信息 根据版本号和平台
func IncreaseUsersInfoAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) ([]map[string]string, error) {
	const sqlStr = "select * from user where app_id = ? and create_time >= ? and create_time <= ? and reg_version = ? and device = ?"
	return conn.Query(sqlStr, appId, startTime, endTime, version, platform)
}

//累计用户
func GrandTotalUsers(appId string, time int64) (int64, error) {
	const sqlStr = "select count(*) as count from `user` where app_id = ? and create_time <= ?"
	maps, err := conn.Query(sqlStr, appId, time)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//累计用户根据平台查询
func GrandTotalUsersAsPlatform(appId string, time int64, device string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and create_time <= ? and device = ?"
	maps, err := conn.Query(sqlStr, appId, time, device)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//累计用户 根据版本查询
func GrandTotalUsersAsVersion(appId string, time int64, version string) (int64, error) {
	const sqlStr = "select count(*) as count from user where app_id = ? and create_time <= ? and reg_version = ?"
	maps, err := conn.Query(sqlStr, appId, time, version)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

func GetOpenNumber(appId string, startTime, endTime int64) (int64, error) {
	const sqlStr = "select count(*) as count from open_log where create_time >= ? and create_time <= ? and app_id = ?"
	maps, err := conn.Query(sqlStr, startTime, endTime, appId)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

func GetOpenNumberAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	const sqlStr = "select count(*) as count from open_log where app_id = ? and create_time >= ? and create_time <= ? and device = ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime, device)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

func GetOpenNumberAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	const sqlStr = "select count(*) as count from open_log where app_id = ? and create_time >= ? and create_time <= ? and version = ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime, version)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

func GetOpenNumberAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	const sqlStr = "select count(*) as count from open_log where app_id = ? and create_time >= ? and create_time <= ? and version = ? and device = ?"
	maps, err := conn.Query(sqlStr, appId, startTime, endTime, version, platform)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//活跃用户统计
func GetActiveNumber(appId string, startTime, endTime int64) (int64, error) {
	const sqlStr = `SELECT COUNT(*) AS count
		FROM (
		select user_id,app_id,count(*) as count
		from open_log
		where create_time >= ? and create_time <= ?
		group by open_log.user_id,open_log.app_id,open_log.device
		having count >= 2 and app_id = ?)
		as a`
	maps, err := conn.Query(sqlStr, startTime, endTime, appId)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//活跃用户统计 筛选平台
func GetActiveNumberAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	const sqlStr = `SELECT COUNT(*) AS count
		FROM (
		select user_id,app_id,count(*) as count
		from open_log
		where create_time >= ? and create_time <= ? and device = ?
		group by open_log.user_id,open_log.app_id
		having count >= 2 and app_id = ?)
		as a`
	maps, err := conn.Query(sqlStr, startTime, endTime, device, appId)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//活跃用户统计 筛选版本
func GetActiveNumberAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	const sqlStr = `SELECT COUNT(*) AS count
		FROM (
		select user_id,app_id,count(*) as count
		from open_log
		where create_time >= ? and create_time <= ? and version = ?
		group by open_log.user_id,open_log.app_id,open_log.device
		having count >= 2 and app_id = ?)
		as a`
	maps, err := conn.Query(sqlStr, startTime, endTime, version, appId)
	if err != nil {
		return 0, err
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//活跃用户统计信息 筛选版本
func GetActiveUsersInfoAsVersion(appId string, startTime, endTime int64, version string) ([]map[string]string, error) {
	const sqlStr = `SELECT *
		FROM (
		select user_id,app_id,count(*) as count
		from open_log
		where create_time >= ? and create_time <= ? and version = ?
		group by open_log.user_id,open_log.app_id,open_log.device
		having count >= 2 and app_id = ?)
		as a`
	return conn.Query(sqlStr, startTime, endTime, version, appId)
}

//活跃用户统计 筛选版本和平台
func GetActiveNumberAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	const sqlStr = `SELECT COUNT(*) AS count
		FROM (
		select user_id,app_id,count(*) as count
		from open_log
		where create_time >= ? and create_time <= ? and version = ? and device = ?
		group by open_log.user_id,open_log.app_id,open_log.device
		having count >= 2 and app_id = ?)
		as a`
	maps, err := conn.Query(sqlStr, startTime, endTime, version, platform, appId)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	m := maps[0]
	return utility.ToInt64(m["count"]), nil
}

//活跃用户统计信息 筛选版本和平台
func GetActiveUsersInfoAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) ([]map[string]string, error) {
	const sqlStr = `SELECT *
		FROM (
		select user_id,app_id,count(*) as count
		from open_log
		where create_time >= ? and create_time <= ? and version = ? and device = ?
		group by open_log.user_id,open_log.app_id,open_log.device
		having count >= 2 and app_id = ?)
		as a`
	return conn.Query(sqlStr, startTime, endTime, version, platform, appId)
}

func InsertAdminOperateLog(tx *mysql.MysqlTx, admin, target string, targetType, optType int, reason string, createTime, effectiveTime int64) error {
	const sqlStr = "insert into `admin_operate_log`(operator,`type`,target,operate_type,reason,create_time,effective_time) values(?,?,?,?,?,?,?)"
	_, _, err := tx.Exec(sqlStr, admin, targetType, target, optType, reason, createTime, effectiveTime)
	return err
}

//用户被封号次数
func FindUserBanTimes(userId string) (int, error) {
	const sqlStr = "select count(*) as count from `admin_operate_log` where target = ? and type = ? and operate_type = ?"
	maps, err := conn.Query(sqlStr, userId, types.IsFriend, types.BanUser)
	if err != nil {
		return 0, err
	}
	if len(maps) <= 0 {
		return 0, nil
	}
	info := maps[0]
	return utility.ToInt(info["count"]), nil
}

//群被封号次数
func FindRoomBanTimes(roomId string) (int, error) {
	const sqlStr = "select count(*) as count from `admin_operate_log` where target = ? and type = ? and operate_type = ?"
	maps, err := conn.Query(sqlStr, roomId, types.IsRoom, types.BanRoom)
	if err != nil {
		return 0, err
	}
	if len(maps) <= 0 {
		return 0, nil
	}
	info := maps[0]
	return utility.ToInt(info["count"]), nil
}

//查询操作日志列表，根据操作类型
func FindAdminOperateLogByOptType(appId string, optType int, startId int64, number int) ([]map[string]string, error) {
	const sqlStr = "select * from `admin_operate_log` as opt left join `admin` on opt.operator = `admin`.id where `admin`.app_id = ? and opt.`operate_type` = ? order by opt.create_time desc limit ?,?"
	return conn.Query(sqlStr, appId, optType, startId, number)
}

//查询操作日志列表
func FindAdminOperateLog(appId, query string, optType *int, startId int64, number int) (int64, []map[string]string, error) {
	queryType := ""
	queryName := ""
	if optType != nil {
		queryType = fmt.Sprintf(" and opt_view.`operate_type` = %v", *optType)
	}
	if query != "" {
		query = QueryStr(query)
		queryName = fmt.Sprintf(" and opt_view.mark_id like '%%%s%%' or opt_view.`name` like '%%%s%%' ", query, query)
	}

	totalStr := fmt.Sprintf("SELECT count(*) as count"+
		" from (select opt.id,opt.operator,opt.type,opt.target,opt.operate_type,"+
		" case when opt.type = 1 then"+
		" room.mark_id"+
		" when opt.type = 2 then"+
		" `user`.uid"+
		" end as 'mark_id',"+
		" case when opt.type = 1 then"+
		" room.`name`"+
		" when opt.type = 2 then"+
		" `user`.account"+
		" end as 'name',"+
		" opt.reason, opt.create_time, opt.effective_time, `admin`.app_id"+
		" from admin_operate_log as opt"+
		" left join room on opt.type = 1 and opt.target = room.id"+
		" left join `user` on opt.type = 2 and opt.target = `user`.user_id"+
		" left join `admin` on opt.operator = `admin`.id) as opt_view"+
		" WHERE opt_view.app_id = ? %s %s;", queryType, queryName)
	maps, err := conn.Query(totalStr, appId)
	if err != nil {
		return 0, nil, err
	}
	if len(maps) < 0 {
		return 0, nil, nil
	}
	total := utility.ToInt64(maps[0]["count"])
	sqlStr := fmt.Sprintf("SELECT *"+
		" from (select opt.id,opt.operator,opt.type,opt.target,opt.operate_type,"+
		" case when opt.type = 1 then"+
		" room.mark_id"+
		" when opt.type = 2 then"+
		" `user`.uid"+
		" end as 'mark_id',"+
		" case when opt.type = 1 then"+
		" room.`name`"+
		" when opt.type = 2 then"+
		" `user`.account"+
		" end as 'name',"+
		" opt.reason, opt.create_time, opt.effective_time, `admin`.app_id"+
		" from admin_operate_log as opt"+
		" left join room on opt.type = 1 and opt.target = room.id"+
		" left join `user` on opt.type = 2 and opt.target = `user`.user_id"+
		" left join `admin` on opt.operator = `admin`.id) as opt_view"+
		" WHERE opt_view.app_id = ? %s %s order by opt_view.create_time desc limit ?,?;", queryType, queryName)
	maps, err = conn.Query(sqlStr, appId, startId, number)
	return total, maps, err
}

func CloseUser(tx *mysql.MysqlTx, userId string, closeUntil int64) error {
	const sqlStr = "update `user` set close_until = ? where user_id = ?"
	_, _, err := tx.Exec(sqlStr, closeUntil, userId)
	return err
}

func CloseRoom(tx *mysql.MysqlTx, roomId string, closeUntil int64) error {
	const sqlStr = "update `room` set close_until = ? where id = ?"
	_, _, err := tx.Exec(sqlStr, closeUntil, roomId)
	return err
}

//-----------------------//

func InsertAd(appId, name, url string, duration int, link string, isActive int) (int64, error) {
	const sqlStr = "insert into advertisement(app_id,name,url,duration,link,is_active) values(?,?,?,?,?,?)"
	_, id, err := conn.Exec(sqlStr, appId, name, url, duration, link, isActive)
	return id, err
}

func GetAdNumbers(appId string) ([]map[string]string, error) {
	const sqlStr = "select count(*) as count from advertisement where app_id = ? and is_delete = ?"
	return conn.Query(sqlStr, appId, types.AdIsNotDelete)
}

func GetAllAd(appId string) ([]map[string]string, error) {
	const sqlStr = "select * from advertisement where app_id = ? and is_delete = ?"
	return conn.Query(sqlStr, appId, types.AdIsNotDelete)
}

//设置广告名称
func SetAdName(id, name string) error {
	const sqlStr = "update advertisement set name = ? where id = ?"
	_, _, err := conn.Exec(sqlStr, name, id)
	return err
}

func ActiveAd(id string, isActive int) error {
	const sqlStr = "update advertisement set is_active = ? where id = ?"
	_, _, err := conn.Exec(sqlStr, isActive, id)
	return err
}

func DeleteAd(id string) error {
	const sqlStr = "update advertisement set is_delete = ? where id = ?"
	_, _, err := conn.Exec(sqlStr, types.AdIsDelete, id)
	return err
}
