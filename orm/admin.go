package orm

import (
	"github.com/33cn/chat33/db"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/types"
	"github.com/inconshreveable/log15"
)

var logAdmin = log15.New("model", "orm/admin")

func AppendAppOpenLog() {

}

//管理员登录 获取id
func AdminCheckLogin(appId, account string) (*types.Admin, error) {
	ret, err := mysql.AdminCheckLogin(appId, account)
	if err != nil {
		logAdmin.Error("mysql.AdminCheckLogin", "err", err, "appId", appId, "account", account)
	}
	return ret, err
}

//管理员信息
func AdminAccountInfo(adminId string) (*types.Admin, error) {
	ret, err := mysql.FindAdminById(adminId)
	if err != nil {
		logAdmin.Error("mysql.FindAdminById", "err", err, "adminId", adminId)
	}
	return ret, err
}

//用户被封号次数
func FindUserBanTimes(userId string) (int, error) {
	ret, err := mysql.FindUserBanTimes(userId)
	if err != nil {
		logAdmin.Error("mysql.FindUserBanTimes", "err", err, "userId", userId)
	}
	return ret, err
}

func CloseUser(tx types.Tx, userId string, closeUntil int64) error {
	err := mysql.CloseUser(tx, userId, closeUntil)
	if err != nil {
		logAdmin.Error("mysql.CloseUser", "err", err, "userId", userId, "closeUntil", closeUntil)
	}
	return err
}

//群被封群次数
func FindRoomBanTimes(roomId string) (int, error) {
	ret, err := mysql.FindRoomBanTimes(roomId)
	if err != nil {
		logAdmin.Error("mysql.FindRoomBanTimes", "err", err, "roomId", roomId)
	}
	return ret, err
}

func CloseRoom(tx types.Tx, roomId string, closeUntil int64) error {
	err := mysql.CloseRoom(tx, roomId, closeUntil)
	if err != nil {
		logAdmin.Error("mysql.CloseRoom", "err", err, "roomId", roomId, "closeUntil", closeUntil)
	}
	return err
}

func GetVersions(appId string) ([]string, error) {
	ret, err := mysql.GetVersionList(appId)
	if err != nil {
		logAdmin.Error("mysql.GetVersionList", "err", err, "appId", appId)
	}
	return ret, err
}

func UpdateVersionUsersNumber(appId, nowVersion string) (int64, error) {
	ret, err := mysql.UpdateVersionUsersNumber(appId, nowVersion)
	if err != nil {
		logAdmin.Error("mysql.UpdateVersionUsersNumber", "err", err, "appId", appId, "nowVersion", nowVersion)
	}
	return ret, err
}

func IncreaseUsers(appId string, startTime, endTime int64) (int64, error) {
	ret, err := mysql.IncreaseUsers(appId, startTime, endTime)
	if err != nil {
		logAdmin.Error("mysql.IncreaseUsers", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime)
	}
	return ret, err
}

func IncreaseUsersAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	ret, err := mysql.IncreaseUsersAsPlatform(appId, startTime, endTime, device)
	if err != nil {
		logAdmin.Error("mysql.IncreaseUsersAsPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "device", device)
	}
	return ret, err
}

func IncreaseUsersAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	ret, err := mysql.IncreaseUsersAsVersion(appId, startTime, endTime, version)
	if err != nil {
		logAdmin.Error("mysql.IncreaseUsersAsVersion", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version)
	}
	return ret, err
}

func IncreaseUsersInfoAsVersion(appId string, startTime, endTime int64, version string) ([]*types.User, error) {
	ret, err := mysql.IncreaseUsersInfoAsVersion(appId, startTime, endTime, version)
	if err != nil {
		logAdmin.Error("mysql.IncreaseUsersInfoAsVersion", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version)
	}
	return ret, err
}

func IncreaseUsersAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	ret, err := mysql.IncreaseUsersAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		logAdmin.Error("mysql.IncreaseUsersAsVersionAndPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version, "platform", platform)
	}
	return ret, err
}

func IncreaseUsersInfoAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) ([]*types.User, error) {
	ret, err := mysql.IncreaseUsersInfoAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		logAdmin.Error("mysql.IncreaseUsersInfoAsVersionAndPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version, "platform", platform)
	}
	return ret, err
}

func GrandTotalUsers(appId string, time int64) (int64, error) {
	ret, err := mysql.GrandTotalUsers(appId, time)
	if err != nil {
		logAdmin.Error("mysql.GrandTotalUsers", "err", err, "appId", appId, "time", time)
	}
	return ret, err
}

func GrandTotalUsersAsPlatform(appId string, time int64, device string) (int64, error) {
	ret, err := mysql.GrandTotalUsersAsPlatform(appId, time, device)
	if err != nil {
		logAdmin.Error("mysql.GrandTotalUsersAsPlatform", "err", err, "appId", appId, "time", time, "device", device)
	}
	return ret, err
}

func GrandTotalUsersAsVersion(appId string, time int64, version string) (int64, error) {
	ret, err := mysql.GrandTotalUsersAsVersion(appId, time, version)
	if err != nil {
		logAdmin.Error("mysql.GrandTotalUsersAsVersion", "err", err, "appId", appId, "time", time, "version", version)
	}
	return ret, err
}

//打开应用统计
func GetOpenNumber(appId string, startTime, endTime int64) (int64, error) {
	ret, err := mysql.GetOpenNumber(appId, startTime, endTime)
	if err != nil {
		logAdmin.Error("mysql.GetOpenNumber", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime)
	}
	return ret, err
}

//打开应用统计 筛选平台
func GetOpenNumberAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	ret, err := mysql.GetOpenNumberAsPlatform(appId, startTime, endTime, device)
	if err != nil {
		logAdmin.Error("mysql.GetOpenNumberAsPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "device", device)
	}
	return ret, err
}

//打开应用统计 筛选版本
func GetOpenNumberAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	ret, err := mysql.GetOpenNumberAsVersion(appId, startTime, endTime, version)
	if err != nil {
		logAdmin.Error("mysql.GetOpenNumberAsVersion", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version)
	}
	return ret, err
}

//打开应用统计 筛选平台和版本
func GetOpenNumberAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	ret, err := mysql.GetOpenNumberAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		logAdmin.Error("mysql.GetOpenNumberAsVersionAndPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version, "platform", platform)
	}
	return ret, err
}

//活跃用户统计
func GetActiveNumber(appId string, startTime, endTime int64) (int64, error) {
	ret, err := mysql.GetActiveNumber(appId, startTime, endTime)
	if err != nil {
		logAdmin.Error("mysql.GetActiveNumber", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime)
	}
	return ret, err
}

//活跃用户统计
func GetActiveNumberAsPlatform(appId string, startTime, endTime int64, device string) (int64, error) {
	ret, err := mysql.GetActiveNumberAsPlatform(appId, startTime, endTime, device)
	if err != nil {
		logAdmin.Error("mysql.GetActiveNumberAsPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "device", device)
	}
	return ret, err
}

//活跃用户统计
func GetActiveNumberAsVersion(appId string, startTime, endTime int64, version string) (int64, error) {
	ret, err := mysql.GetActiveNumberAsVersion(appId, startTime, endTime, version)
	if err != nil {
		logAdmin.Error("mysql.GetActiveNumberAsVersion", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version)
	}
	return ret, err
}

//活跃用户统计
func GetActiveUsersInfoAsVersion(appId string, startTime, endTime int64, version string) ([]*types.ActiveUsersView, error) {
	ret, err := mysql.GetActiveUsersInfoAsVersion(appId, startTime, endTime, version)
	if err != nil {
		logAdmin.Error("mysql.GetActiveUsersInfoAsVersion", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version)
	}
	return ret, err
}

//活跃用户统计 根据版本和平台
func GetActiveNumberAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) (int64, error) {
	ret, err := mysql.GetActiveNumberAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		logAdmin.Error("mysql.GetActiveNumberAsVersionAndPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version, "platform", platform)
	}
	return ret, err
}

//活跃用户统计 根据版本和平台
func GetActiveUsersInfoAsVersionAndPlatform(appId string, startTime, endTime int64, version, platform string) ([]*types.ActiveUsersView, error) {
	ret, err := mysql.GetActiveUsersInfoAsVersionAndPlatform(appId, startTime, endTime, version, platform)
	if err != nil {
		logAdmin.Error("mysql.GetActiveUsersInfoAsVersionAndPlatform", "err", err, "appId", appId, "startTime", startTime, "endTime", endTime, "version", version, "platform", platform)
	}
	return ret, err
}

func InsertAdminOperateLog(tx types.Tx, admin, target string, targetType, optType int, reason string, createTime, effectiveTime int64) error {
	err := mysql.InsertAdminOperateLog(tx, admin, target, targetType, optType, reason, createTime, effectiveTime)
	if err != nil {
		logAdmin.Error("mysql.InsertAdminOperateLog", "err", err, "admin", admin, "target", target, "targetType", targetType, "optType", optType, "reason", reason, "createTime", createTime, "effectiveTime", effectiveTime)
	}
	return err
}

func InsertOpenLog(userId, appId, device, version string, createTime int64) error {
	err := mysql.InsertOpenLog(userId, appId, device, version, createTime)
	if err != nil {
		logAdmin.Error("mysql.InsertOpenLog", "err", err, "userId", userId, "appId", appId, "device", device, "version", version, "createTime", createTime, "createTime", createTime)
	}
	return err
}

func FindAdminOperateLog(appId, query string, optType *int, startId int64, number int) (int64, []*types.AdminOptLogView, error) {
	ret, ret2, err := mysql.FindAdminOperateLog(appId, query, optType, startId, number)
	if err != nil {
		logAdmin.Error("mysql.FindAdminOperateLog", "err", err, "query", query, "optType", optType, "startId", startId, "number", number)
	}
	return ret, ret2, err
}

//添加广告
func InsertAd(appId, name, url string, duration int, link string, isActive int) (int64, error) {
	ret, err := mysql.InsertAd(appId, name, url, duration, link, isActive)
	if err != nil {
		logAdmin.Error("mysql.InsertAd", "err", err, "name", name, "url", url, "duration", duration, "link", link, "isActive", isActive)
	}
	return ret, err
}

//获取广告数量
func GetAdNumbers(appId string) (int64, error) {
	ret, err := mysql.GetAdNumbers(appId)
	if err != nil {
		logAdmin.Error("mysql.GetAdNumbers", "err", err, "appId", appId)
	}
	return ret, err
}

//获取广告
func GetAllAd(appId string) ([]*types.Advertisement, error) {
	ret, err := mysql.GetAllAd(appId)
	if err != nil {
		logAdmin.Error("mysql.GetAllAd", "err", err, "appId", appId)
	}
	return ret, err
}

//设置广告名称
func SetAdName(id, name string) error {
	err := mysql.SetAdName(id, name)
	if err != nil {
		logAdmin.Error("mysql.SetAdName", "err", err, "id", id, "name", name)
	}
	return err
}

func ActiveAd(id string, isActive int) error {
	err := mysql.ActiveAd(id, isActive)
	if err != nil {
		logAdmin.Error("mysql.ActiveAd", "err", err, "id", id, "isActive", isActive)
	}
	return err
}

func DeleteAd(id string) error {
	err := mysql.DeleteAd(id)
	if err != nil {
		logAdmin.Error("mysql.DeleteAd", "err", err, "id", id)
	}
	return err
}

//设置推荐群
func SetRecommendRoom(id string, recommend int) (int64, int64, error) {
	ret, ret2, err := db.SetRecommendRoom(id, recommend)
	if err != nil {
		logAdmin.Error("db.SetRecommendRoom", "err", err, "id", id, "recommend", recommend)
	}
	return ret, ret2, err
}
