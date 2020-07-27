package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func AdminCheckLogin(appId, account string) (*types.Admin, error) {
	maps, err := db.AdminCheckLogin(appId, account)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, err
	}
	info := maps[0]
	return &types.Admin{
		Id:       info["id"],
		AppId:    info["app_id"],
		Account:  info["account"],
		Password: info["password"],
		Salt:     info["salt"],
	}, nil
}

// 获取管理员信息 根据id
func FindAdminById(id string) (*types.Admin, error) {
	maps, err := db.FindAdminById(id)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, err
	}
	info := maps[0]
	return &types.Admin{
		Id:       info["id"],
		AppId:    info["app_id"],
		Account:  info["account"],
		Password: info["password"],
		Salt:     info["salt"],
	}, nil
}

//用户被封号次数
func FindUserBanTimes(userId string) (int, error) {
	return db.FindUserBanTimes(userId)
}

func CloseUser(tx types.Tx, userId string, closeUntil int64) error {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	return db.CloseUser(tx.(*mysql.MysqlTx), userId, closeUntil)
}

//群被封群次数
func FindRoomBanTimes(roomId string) (int, error) {
	return db.FindRoomBanTimes(roomId)
}

func CloseRoom(tx types.Tx, roomId string, closeUntil int64) error {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	return db.CloseRoom(tx.(*mysql.MysqlTx), roomId, closeUntil)
}

func GetVersionList(appId string) ([]string, error) {
	maps, err := db.GetVersionList(appId)
	if err != nil {
		return nil, err
	}
	versions := make([]string, 0)
	for _, v := range maps {
		if v["now_version"] == "" {
			continue
		}
		versions = append(versions, v["now_version"])
	}
	return versions, nil
}

func UpdateVersionUsersNumber(appId, nowVersion string) (int64, error) {
	return db.UpdateVersionUsersNumber(appId, nowVersion)
}

func IncreaseUsers(appId string, startTime, endTime int64) (int64, error) {
	return db.IncreaseUsers(appId, startTime, endTime)
}

func IncreaseUsersAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	return db.IncreaseUsersAsPlatform(appId, startTime, endTime, device)
}

func IncreaseUsersAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	return db.IncreaseUsersAsVersion(appId, startTime, endTime, version)
}

func IncreaseUsersInfoAsVersion(appId string, startTime, endTime int64, version string) ([]*types.User, error) {
	maps, err := db.IncreaseUsersInfoAsVersion(appId, startTime, endTime, version)
	if err != nil {
		return nil, err
	}

	list := make([]*types.User, 0)
	for _, m := range maps {
		item := convertUser(m)
		list = append(list, item)
	}
	return list, nil
}

func IncreaseUsersAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	return db.IncreaseUsersAsVersionAndPlatform(appId, startTime, endTime, version, platform)
}

func IncreaseUsersInfoAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) ([]*types.User, error) {
	maps, err := db.IncreaseUsersInfoAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		return nil, err
	}

	list := make([]*types.User, 0)
	for _, m := range maps {
		item := convertUser(m)
		list = append(list, item)
	}
	return list, nil
}

func GrandTotalUsers(appId string, time int64) (int64, error) {
	return db.GrandTotalUsers(appId, time)
}

func GrandTotalUsersAsPlatform(appId string, time int64, device string) (int64, error) {
	return db.GrandTotalUsersAsPlatform(appId, time, device)
}

func GrandTotalUsersAsVersion(appId string, time int64, version string) (int64, error) {
	return db.GrandTotalUsersAsVersion(appId, time, version)
}

//打开应用统计
func GetOpenNumber(appId string, startTime, endTime int64) (int64, error) {
	return db.GetOpenNumber(appId, startTime, endTime)
}

//打开应用统计 筛选平台
func GetOpenNumberAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	return db.GetOpenNumberAsPlatform(appId, startTime, endTime, device)
}

//打开应用统计 筛选平台
func GetOpenNumberAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	return db.GetOpenNumberAsVersion(appId, startTime, endTime, version)
}

//打开应用统计 筛选版本和平台
func GetOpenNumberAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	return db.GetOpenNumberAsVersionAndPlatform(appId, startTime, endTime, version, platform)
}

//活跃用户统计
func GetActiveNumber(appId string, startTime, endTime int64) (int64, error) {
	return db.GetActiveNumber(appId, startTime, endTime)
}

//活跃用户统计
func GetActiveNumberAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	return db.GetActiveNumberAsPlatform(appId, startTime, endTime, device)
}

//活跃用户统计
func GetActiveNumberAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	return db.GetActiveNumberAsVersion(appId, startTime, endTime, version)
}

//活跃用户统计信息
func GetActiveUsersInfoAsVersion(appId string, startTime, endTime int64, version string) ([]*types.ActiveUsersView, error) {
	maps, err := db.GetActiveUsersInfoAsVersion(appId, startTime, endTime, version)
	if err != nil {
		return nil, err
	}
	list := make([]*types.ActiveUsersView, 0)
	for _, m := range maps {
		item := &types.ActiveUsersView{
			UserId: m["user_id"],
			AppId:  m["app_id"],
			Count:  m["count"],
		}
		list = append(list, item)
	}
	return list, nil
}

//活跃用户统计 根据版本和平台
func GetActiveNumberAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	return db.GetActiveNumberAsVersionAndPlatform(appId, startTime, endTime, version, platform)
}

//活跃用户统计信息 根据版本和平台
func GetActiveUsersInfoAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) ([]*types.ActiveUsersView, error) {
	maps, err := db.GetActiveUsersInfoAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		return nil, err
	}
	list := make([]*types.ActiveUsersView, 0)
	for _, m := range maps {
		item := &types.ActiveUsersView{
			UserId: m["user_id"],
			AppId:  m["app_id"],
			Count:  m["count"],
		}
		list = append(list, item)
	}
	return list, nil
}

func InsertAdminOperateLog(tx types.Tx, admin, target string, targetType, optType int, reason string, createTime, effectiveTime int64) error {
	if _, b := tx.(*mysql.MysqlTx); !b {
		panic("types tx to sql tx failed")
	}
	return db.InsertAdminOperateLog(tx.(*mysql.MysqlTx), admin, target, targetType, optType, reason, createTime, effectiveTime)
}

func InsertOpenLog(userId, appId, device, version string, createTime int64) error {
	return db.InsertOpenLog(userId, appId, device, version, createTime)
}

//查询操作日志列表，根据操作类型
func FindAdminOperateLogByOptType(appId string, optType int, startId int64, number int) ([]*types.AdminOptLog, error) {
	maps, err := db.FindAdminOperateLogByOptType(appId, optType, startId, number)
	if err != nil {
		return nil, err
	}

	list := make([]*types.AdminOptLog, 0)
	for _, info := range maps {
		item := &types.AdminOptLog{
			Id:            info["id"],
			Operator:      info["operator"],
			Type:          utility.ToInt(info["type"]),
			Target:        info["target"],
			OperateType:   utility.ToInt(info["operate_type"]),
			Reason:        info["reason"],
			CreateTime:    utility.ToInt64(info["create_time"]),
			EffectiveTime: utility.ToInt64(info["effective_time"]),
		}

		list = append(list, item)
	}
	return list, nil
}

//查询操作日志列表
func FindAdminOperateLog(appId, query string, optType *int, startId int64, number int) (int64, []*types.AdminOptLogView, error) {
	total, maps, err := db.FindAdminOperateLog(appId, query, optType, startId, number)
	if err != nil {
		return total, nil, err
	}

	list := make([]*types.AdminOptLogView, 0)
	for _, info := range maps {
		item := &types.AdminOptLogView{
			Id:            info["id"],
			AppId:         info["appId"],
			Operator:      info["operator"],
			Type:          utility.ToInt(info["type"]),
			Target:        info["target"],
			MarkId:        info["mark_id"],
			Name:          info["name"],
			OperateType:   utility.ToInt(info["operate_type"]),
			Reason:        info["reason"],
			CreateTime:    utility.ToInt64(info["create_time"]),
			EffectiveTime: utility.ToInt64(info["effective_time"]),
		}

		list = append(list, item)
	}
	return total, list, nil
}

func InsertAd(appId, name, url string, duration int, link string, isActive int) (int64, error) {
	return db.InsertAd(appId, name, url, duration, link, isActive)
}

func GetAdNumbers(appId string) (int64, error) {
	maps, err := db.GetAdNumbers(appId)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	info := maps[0]
	return utility.ToInt64(info["count"]), nil
}

func GetAllAd(appId string) ([]*types.Advertisement, error) {
	maps, err := db.GetAllAd(appId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}

	list := make([]*types.Advertisement, 0)
	for _, m := range maps {
		item := &types.Advertisement{
			Id:       m["id"],
			AppId:    m["app_id"],
			Name:     m["name"],
			Url:      m["url"],
			Duration: utility.ToInt(m["duration"]),
			Link:     m["link"],
			IsActive: utility.ToInt(m["is_active"]),
			IsDelete: utility.ToInt(m["is_delete"]),
		}
		list = append(list, item)
	}
	return list, nil
}

//设置广告名称
func SetAdName(id, name string) error {
	return db.SetAdName(id, name)
}

func ActiveAd(id string, isActive int) error {
	return db.ActiveAd(id, isActive)
}

func DeleteAd(id string) error {
	return db.DeleteAd(id)
}
