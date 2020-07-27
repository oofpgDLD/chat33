package model

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"os"
	"time"

	"github.com/33cn/chat33/pkg/work"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/pkg/account"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
	"github.com/tealeg/xlsx"
)

var logAdmin = log15.New("model", "model/admin")

//组消息体
func ConvertUserClosedAlertStr(appId string, endTime int64) string {
	cus := types.CustomServiceInfo[appId]
	timeStr := ""
	if endTime != types.MutedForvevr {
		/*deadTime := utility.MillionSecondToDateTime(endTime)
		//dua, _:= time.ParseDuration(time.Now().Format("2006-01-02 15:04"))
		subDua := deadTime.Sub(time.Now())
		s := fmt.Sprintf("%v", subDua)
		index := strings.LastIndex(s, "m")
		r := strings.NewReplacer("h", "小时", "m", "分钟")
		timeStr = "该群聊已被查封至" + r.Replace(s[0:index+1]) + "，"*/

		deadTime := utility.MillionSecondToDateTime(endTime)
		timeStr = "你的账号已被查封至" + deadTime.Format("2006-01-02 15:04") + "，"
	} else {
		timeStr = "你的账号被永久查封，"
	}
	str := fmt.Sprintf("%s如需解封可联系客服：%s", timeStr, cus)
	return str
}

func ConvertRoomClosedAlertStr(appId string, endTime int64) string {
	cus := types.CustomServiceInfo[appId]
	timeStr := ""
	if endTime != types.MutedForvevr {
		/*deadTime := utility.MillionSecondToDateTime(endTime)
		//dua, _:= time.ParseDuration(time.Now().Format("2006-01-02 15:04"))
		subDua := deadTime.Sub(time.Now())
		s := fmt.Sprintf("%v", subDua)
		index := strings.LastIndex(s, "m")
		r := strings.NewReplacer("h", "小时", "m", "分钟")
		timeStr = "该群聊已被查封至" + r.Replace(s[0:index+1]) + "，"*/

		deadTime := utility.MillionSecondToDateTime(endTime)
		timeStr = "该群聊已被查封至" + deadTime.Format("2006-01-02 15:04") + "，"
	} else {
		timeStr = "该群聊已被永久查封，"
	}
	str := fmt.Sprintf("%s如需解封可联系客服：%s", timeStr, cus)
	return str
}

func ConvertCreateRoomsOutOfLimit(number int) string {
	return fmt.Sprintf("群聊已达上限，你最多可创建%v个群聊", number)
}

func ConvertRoomMembersLimit(number int) string {
	return fmt.Sprintf("群聊人数最多%v人", number)
}

//应用埋点
func AppOpen(userId, appId, device, version string) error {
	err := orm.InsertOpenLog(userId, appId, device, version, utility.NowMillionSecond())
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//管理员登录
func AdminLogin(appId, id, password string) (string, error) {
	adminInfo, err := orm.AdminCheckLogin(appId, id)
	if err != nil {
		return "", result.NewError(result.DbConnectFail)
	}
	if adminInfo == nil {
		logAdmin.Debug("can not find admin", "appId", appId, "id", id, "password", password)
		return "", result.NewError(result.AdminLoginFailed)
	}
	pwd := utility.EncodeSha256WithSalt(password, adminInfo.Salt)

	if pwd != adminInfo.Password {
		logAdmin.Debug("can find admin,but password not match", "appId", appId, "id", id, "password", password)
		return "", result.NewError(result.AdminLoginFailed)
	}
	return adminInfo.Id, nil
}

//管理员账号
func AdminAccount(id string) (string, error) {
	adminInfo, err := orm.AdminAccountInfo(id)
	if err != nil {
		return "", result.NewError(result.DbConnectFail)
	}
	if adminInfo == nil {
		logAdmin.Debug("can not find admin by id", "id", id)
		return "", result.NewError(result.QueryDbFailed)
	}
	return adminInfo.Account, nil
}

//用户管理 数量统计
func AdminUsersCount(appId string) (interface{}, error) {
	totalUser, _ := orm.GrandTotalUsers(appId, utility.NowMillionSecond())
	totalRoom, _ := orm.GetRoomCountInApp(appId)
	banUser, _ := orm.GetCloseUserCountInApp(appId)
	banRoom, _ := orm.GetCloseRoomCountInApp(appId)

	ret := make(map[string]interface{})
	ret["totalUser"] = totalUser
	ret["totalRoom"] = totalRoom
	ret["banUser"] = banUser
	ret["banRoom"] = banRoom

	return ret, nil
}

//用户管理 用户详情列表
func AdminUsersList(appId string, cTypes []int, query string, page, count int64) (interface{}, error) {
	tp := 0
	if len(cTypes) == 0 {
		tp = types.QueryUsersAll
	} else {
		tp = cTypes[0]
	}

	if tp != types.QueryUsersAll && query != "" {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "lack query param")
	}

	var users []*types.User
	switch tp {
	case types.QueryUsersAll:
		users, _ = orm.GetUsersInAppQueryUid(appId, query)
	case types.QueryUsersUnClose:
		users, _ = orm.GetUsersInAppUnClose(appId)
	case types.QueryUsersClosed:
		users, _ = orm.GetUsersInAppClosed(appId)
	default:
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "type not exist")
	}

	totalCount := len(users)
	list := make([]interface{}, 0)
	if page*count > utility.ToInt64(totalCount) {
		users = users[(page-1)*count:]
	} else {
		users = users[(page-1)*count : page*count]
	}
	for _, u := range users {
		item := make(map[string]interface{})
		item["id"] = u.UserId
		item["uid"] = u.Uid
		item["account"] = u.Account
		item["createTime"] = u.CreateTime
		item["roomCreatedCount"], _ = orm.FindCreateRoomNumbers(u.UserId)
		item["banTimes"], _ = orm.FindUserBanTimes(u.UserId)
		item["banEndTime"] = u.CloseUntil
		list = append(list, item)
	}

	ret := make(map[string]interface{})
	ret["total"] = totalCount
	ret["list"] = list
	return ret, nil
}

//用户管理 封用户
func BanUser(appId, admin, id, reason string, endTime int64) error {
	user, err := orm.GetUserInfoById(id)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if user.AppId != appId {
		return result.NewError(result.UserNotExists)
	}

	tx, err := orm.GetTx()
	if err != nil {
		logAdmin.Error("BanUser tx GetTx", "err", err.Error())
		return result.NewError(result.ServerInterError)
	}
	err = orm.CloseUser(tx, id, endTime)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	createTime := utility.NowMillionSecond()
	err = orm.InsertAdminOperateLog(tx, admin, id, types.IsFriend, types.BanUser, reason, createTime, endTime)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = tx.Commit()
	if err != nil {
		logAdmin.Error("BanUser tx Commit", "err", err.Error())
		return result.NewError(result.DbConnectFail)
	}
	proto.SendCloseUserAccountNotification(endTime, id, ConvertUserClosedAlertStr(appId, endTime))
	return nil
}

//用户管理 解封用户
func BanUserCancel(appId, admin, id string) error {
	user, err := orm.GetUserInfoById(id)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if user.AppId != appId {
		return result.NewError(result.UserNotExists)
	}

	tx, err := orm.GetTx()
	if err != nil {
		logAdmin.Error("BanUserCancel tx GetTx", "err", err.Error())
		return result.NewError(result.ServerInterError)
	}
	err = orm.CloseUser(tx, id, 0)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = orm.InsertAdminOperateLog(tx, admin, id, types.IsFriend, types.BanUserCancel, "", utility.NowMillionSecond(), 0)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = tx.Commit()
	if err != nil {
		logAdmin.Error("BanUserCancel tx Commit", "err", err.Error())
		return result.NewError(result.DbConnectFail)
	}
	proto.SendCloseUserAccountNotification(0, id, "")
	return nil
}

//群管理 设置上限
func SetLimit(appId string, memberLimit, roomCreateLimit int) error {
	err := orm.SetCreateRoomsLimit(appId, 1, roomCreateLimit)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	err = orm.SetRoomMembersLimit(appId, 1, memberLimit)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//群管理 获取上限
func GetLimit(appId string) (interface{}, error) {
	roomCreateLimit, err := orm.GetCreateRoomsLimit(appId, 1)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	memberLimit, err := orm.GetRoomMembersLimit(appId, 1)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	ret := make(map[string]interface{})
	ret["memberLimit"] = memberLimit
	ret["roomCreateLimit"] = roomCreateLimit
	return ret, nil
}

//群管理 群数量统计
func AdminRoomsCount(appId string) (interface{}, error) {
	totalRoom, _ := orm.GetRoomCountInApp(appId)
	banRoom, _ := orm.GetCloseRoomCountInApp(appId)

	ret := make(map[string]interface{})
	ret["total"] = totalRoom
	ret["ban"] = banRoom
	return ret, nil
}

//群管理 群详情列表
func AdminRoomsList(appId string, cTypes []int, query string, page, count int64) (interface{}, error) {
	tp := 0
	if len(cTypes) == 0 {
		tp = types.QueryRoomsAll
	} else {
		tp = cTypes[0]
	}

	if tp != types.QueryRoomsAll && query != "" {
		return nil, result.NewError(result.ParamsError)
	}

	var rooms []*types.RoomJoinUser
	switch tp {
	case types.QueryRoomsAll:
		rooms, _ = orm.FindRoomsInAppQueryMarkId(appId, query)
	case types.QueryRoomsUnClose:
		rooms, _ = orm.FindRoomsInAppUnClose(appId)
	case types.QueryRoomsClosed:
		rooms, _ = orm.FindRoomsInAppClosed(appId)
	case types.QueryRoomsRecommend:
		rooms, _ = orm.FindRoomsInAppRecommend(appId)
	default:
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, fmt.Sprintf("type %d not exists", tp))
	}

	totalCount := len(rooms)
	list := make([]interface{}, 0)
	if page*count > utility.ToInt64(totalCount) {
		rooms = rooms[(page-1)*count:]
	} else {
		rooms = rooms[(page-1)*count : page*count]
	}
	for _, r := range rooms {
		item := make(map[string]interface{})
		item["id"] = r.Id
		item["name"] = r.Name
		item["markId"] = r.Room.MarkId
		owner := make(map[string]interface{})
		owner["id"] = r.UserId
		owner["account"] = r.Account
		owner["uid"] = r.Uid
		item["owner"] = owner

		item["memberCount"], _ = orm.GetMemberNumber(r.Id)
		item["banTimes"], _ = orm.FindRoomBanTimes(r.Id)
		item["banEndTime"] = r.Room.CloseUntil
		item["recommend"] = r.Room.Recommend
		list = append(list, item)
	}

	ret := make(map[string]interface{})
	ret["total"] = totalCount
	ret["list"] = list
	return ret, nil
}

//群管理 封禁群
func BanRoom(appId, admin, id, reason string, endTime int64) error {
	room, err := orm.FindRoomById(id, types.RoomNotDeleted)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if room == nil {
		logAdmin.Debug("BanRoom room not find", "appId", appId, "admin", admin, "roomId", id)
		return result.NewError(result.RoomNotExists)
	}
	user, err := orm.GetUserInfoById(room.MasterId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if user.AppId != appId {
		logAdmin.Debug("BanRoom master appId not match", "appId", appId, "master appId", user.AppId, "roomId", id, "masterId", room.MasterId)
		return result.NewError(result.RoomNotExists)
	}

	tx, err := orm.GetTx()
	if err != nil {
		logAdmin.Error("BanRoom tx GetTx", "err", err.Error())
		return result.NewError(result.DbConnectFail)
	}
	err = orm.CloseRoom(tx, id, endTime)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = orm.InsertAdminOperateLog(tx, admin, id, types.IsRoom, types.BanRoom, reason, utility.NowMillionSecond(), endTime)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = tx.Commit()
	if err != nil {
		logAdmin.Error("BanRoom tx Commit", "err", err.Error())
		return result.NewError(result.DbConnectFail)
	}
	proto.SendClosedRoomNotification(id, ConvertRoomClosedAlertStr(appId, endTime), endTime)
	return nil
}

//群管理 解封群
func BanRoomCancel(appId, admin, id string) error {
	room, err := orm.FindRoomById(id, types.RoomNotDeleted)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if room == nil {
		logAdmin.Debug("BanRoomCancel room not find", "appId", appId, "admin", admin, "roomId", id)
		return result.NewError(result.RoomNotExists)
	}
	user, err := orm.GetUserInfoById(room.MasterId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if user.AppId != appId {
		logAdmin.Debug("BanRoomCancel master appId not match", "appId", appId, "master appId", user.AppId, "roomId", id, "masterId", room.MasterId)
		return result.NewError(result.RoomNotExists)
	}

	tx, err := orm.GetTx()
	if err != nil {
		logAdmin.Error("BanRoomCancel tx GetTx", "err", err.Error())
		return result.NewError(result.DbConnectFail)
	}
	err = orm.CloseRoom(tx, id, 0)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = orm.InsertAdminOperateLog(tx, admin, id, types.IsRoom, types.BanRoomCancel, "", utility.NowMillionSecond(), 0)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = tx.Commit()
	if err != nil {
		logAdmin.Error("BanRoom tx Commit", "err", err.Error())
		return result.NewError(result.DbConnectFail)
	}
	proto.SendClosedRoomNotification(id, "", 0)
	return nil
}

func SetRecommend(appId, id string, recommend int) error {
	room, err := orm.FindRoomById(id, types.RoomNotDeleted)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	_, _, err = orm.SetRecommendRoom(id, recommend)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	//更新缓存
	re := GetRecommendRoom(appId)
	if recommend == types.RoomRecommend {
		re.AddData(room)
	} else {
		re.DelData(room)
	}
	return nil
}

func AdminOperateLog(appId, query string, cTypes []int, page, count int) (interface{}, error) {
	tp := 0
	if len(cTypes) == 0 {
		tp = types.QueryAdminOptLogAll
	} else {
		tp = cTypes[0]
	}

	if tp != types.QueryAdminOptLogAll && query != "" {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "lack query param")
	}

	var optType *int
	if tp == types.QueryAdminOptLogAll {
		optType = nil
	} else {
		v := types.QueryAdminOptParamToDb[tp]
		optType = &v
	}

	total, logs, err := orm.FindAdminOperateLog(appId, query, optType, utility.ToInt64((page-1)*count), count)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	list := make([]interface{}, 0)
	for _, v := range logs {
		item := make(map[string]interface{})
		item["datetime"] = v.CreateTime
		item["type"] = types.QueryAdminOptDbToParam[v.OperateType]

		info := make(map[string]interface{})
		info["id"] = v.Id
		if v.Type == types.IsRoom {
			info["markId"] = v.MarkId
			info["roomName"] = v.Name
		} else if v.Type == types.IsFriend {
			info["uid"] = v.MarkId
			info["account"] = v.Name
		}
		item["targetInfo"] = info

		list = append(list, item)
	}
	ret := make(map[string]interface{})
	ret["total"] = total
	ret["list"] = list

	return ret, nil
}

//---------------------统计数据---------------------//
//近期统计数据
func LatestData(appId string) (interface{}, error) {
	todayStart, todayEnd := utility.DayStartAndEndNowMillionSecond(time.Now())
	todayNew, _ := orm.IncreaseUsers(appId, todayStart, todayEnd)
	todayTotal, _ := orm.GrandTotalUsers(appId, todayEnd)
	todayOpen, _ := orm.GetOpenNumber(appId, todayStart, todayEnd)
	todayActive, _ := orm.GetActiveNumber(appId, todayStart, todayEnd)

	today := make(map[string]interface{})
	today["new"] = todayNew
	today["active"] = todayActive
	today["open"] = todayOpen
	today["total"] = todayTotal

	yesterdayStart, yesterdayEnd := utility.DayStartAndEndNowMillionSecond(time.Now().AddDate(0, 0, -1))
	yesterdayNew, _ := orm.IncreaseUsers(appId, yesterdayStart, yesterdayEnd)
	yesterdayTotal, _ := orm.GrandTotalUsers(appId, yesterdayEnd)
	yesterdayOpen, _ := orm.GetOpenNumber(appId, yesterdayStart, yesterdayEnd)
	yesterdayActive, _ := orm.GetActiveNumber(appId, yesterdayStart, yesterdayEnd)

	yesterday := make(map[string]interface{})
	yesterday["new"] = yesterdayNew
	yesterday["active"] = yesterdayActive
	yesterday["open"] = yesterdayOpen
	yesterday["total"] = yesterdayTotal

	ret := make(map[string]interface{})
	ret["today"] = today
	ret["yesterday"] = yesterday
	return ret, nil
}

//近期统计数据 根据平台分
func LatestDataPlatform(appId string) (interface{}, error) {
	list := make([]interface{}, 0)
	devices := []string{types.DeviceAndroid, types.DeviceIOS, types.DevicePC}
	for _, dev := range devices {
		todayStart, todayEnd := utility.DayStartAndEndNowMillionSecond(time.Now())
		yesterdayStart, yesterdayEnd := utility.DayStartAndEndNowMillionSecond(time.Now().AddDate(0, 0, -1))

		todayNew, _ := orm.IncreaseUsersAsPlatform(appId, todayStart, todayEnd, dev)
		todayTotal, _ := orm.GrandTotalUsersAsPlatform(appId, todayEnd, dev)
		todayOpen, _ := orm.GetOpenNumberAsPlatform(appId, todayStart, todayEnd, dev)
		todayActive, _ := orm.GetActiveNumberAsPlatform(appId, todayStart, todayEnd, dev)

		yesterdayNew, _ := orm.IncreaseUsersAsPlatform(appId, yesterdayStart, yesterdayEnd, dev)
		yesterdayOpen, _ := orm.GetOpenNumberAsPlatform(appId, yesterdayStart, yesterdayEnd, dev)
		yesterdayActive, _ := orm.GetActiveNumberAsPlatform(appId, yesterdayStart, yesterdayEnd, dev)

		one := make(map[string]interface{})
		one["platform"] = dev
		one["todayNew"] = todayNew
		one["yesterdayNew"] = yesterdayNew
		one["todayActive"] = todayActive
		one["yesterdayActive"] = yesterdayActive
		one["todayOpen"] = todayOpen
		one["yesterdayOpen"] = yesterdayOpen
		one["todayTotal"] = todayTotal

		list = append(list, one)
	}

	ret := make(map[string]interface{})
	ret["list"] = list
	return ret, nil
}

type dayTs struct {
	Start int64
	End   int64
}

func SplitDay(startTime, endTime int64) []*dayTs {
	sTime := utility.MillionSecondToDateTime(startTime)
	eTime := utility.MillionSecondToDateTime(endTime)
	/*s := eTime.Sub(sTime)
	countDay := s / (time.Hour * 24)

	if s % (time.Hour * 24) > 0{
		countDay += 1
	}*/

	list := make([]*dayTs, 0)
	for begin := sTime; eTime.Sub(begin) > 0; {
		dayStart, dayEnd := utility.DayStartAndEndNowMillionSecond(begin)
		begin = begin.AddDate(0, 0, 1)

		list = append(list, &dayTs{
			Start: dayStart,
			End:   dayEnd,
		})
	}
	return list
}

//统计数据 详情
func SumDetails(appId string, startTime, endTime int64, count, page int) (interface{}, error) {
	list := make([]interface{}, 0)

	splitList := SplitDay(startTime, endTime)
	total := len(splitList)

	for i := count * (page - 1); i < count*page && total > i; i++ {
		oneDay := splitList[len(splitList)-i-1] //倒序
		if oneDay == nil {
			continue
		}
		dayStart, dayEnd := oneDay.Start, oneDay.End

		new, _ := orm.IncreaseUsers(appId, dayStart, dayEnd)
		total, _ := orm.GrandTotalUsers(appId, dayEnd)
		open, _ := orm.GetOpenNumber(appId, dayStart, dayEnd)
		active, _ := orm.GetActiveNumber(appId, dayStart, dayEnd)

		one := make(map[string]interface{})
		one["datetime"] = dayStart
		one["new"] = new
		one["active"] = active
		one["open"] = open
		one["total"] = total

		list = append(list, one)
	}

	ret := make(map[string]interface{})
	ret["total"] = total
	ret["list"] = list
	return ret, nil
}

//导出统计数据 详情
func ExportSumDetails(appId string, startTime, endTime int64) (string, error) {
	splitList := SplitDay(startTime, endTime)

	content := make([][]interface{}, 0)
	title := make([]interface{}, 5)
	title[0] = "时间"
	title[1] = "新增用户"
	title[2] = "活跃用户"
	title[3] = "启动次数"
	title[4] = "累计用户"
	content = append(content, title)

	for i := 0; i < len(splitList); i++ {
		oneDay := splitList[len(splitList)-i-1] //倒序
		if oneDay == nil {
			continue
		}
		dayStart, dayEnd := oneDay.Start, oneDay.End

		new, _ := orm.IncreaseUsers(appId, dayStart, dayEnd)
		total, _ := orm.GrandTotalUsers(appId, dayEnd)
		open, _ := orm.GetOpenNumber(appId, dayStart, dayEnd)
		active, _ := orm.GetActiveNumber(appId, dayStart, dayEnd)

		item := make([]interface{}, 5)
		item[0] = utility.MillionSecondToTimeString(dayStart)
		item[1] = new
		item[2] = active
		item[3] = open
		item[4] = total
		content = append(content, item)
	}

	filename := fmt.Sprintf("%s-sum-%v-%v", appId, startTime, endTime)
	url, err := WriteExcel(types.ExcelAddr, filename, content)
	if err != nil {
		logAdmin.Error("write excel failed", "err", err.Error(), "ExcelAddr", types.ExcelAddr, "filename", filename)
		return "", result.NewError(result.AdminExportFailed)
	}
	return url, err
}

func SumGraph(appId string, startTime, endTime int64, typeList []int) (interface{}, error) {
	dataMapArry := make([]interface{}, 0)
	for _, v := range typeList {
		splitList := SplitDay(startTime, endTime)
		data := make([]int64, 0)

		for _, oneDay := range splitList {
			if oneDay == nil {
				continue
			}
			dayStart, dayEnd := oneDay.Start, oneDay.End

			var number int64 = 0
			switch v {
			case types.CountNew:
				number, _ = orm.IncreaseUsers(appId, dayStart, dayEnd)
			case types.CountActive:
				number, _ = orm.GetActiveNumber(appId, dayStart, dayEnd)
			case types.CountOpen:
				number, _ = orm.GetOpenNumber(appId, dayStart, dayEnd)
			case types.CountAll:
				number, _ = orm.GrandTotalUsers(appId, dayEnd)
			default:
				continue
			}
			data = append(data, number)
		}

		item := make(map[string]interface{})
		item["type"] = v
		item["data"] = data
		dataMapArry = append(dataMapArry, item)
	}
	dataMap := make(map[string]interface{})
	dataMap["dataMap"] = dataMapArry

	return dataMap, nil
}

func ModuleEnable(appId, userId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		logAdmin.Debug("ModuleEnable", "err", types.ERR_APPNOTFIND.Error())
		return nil, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	models := make(map[int]types.Module)
	for _, ms := range app.Modules {
		models[ms.Type] = types.Module{
			Type:   ms.Type,
			Name:   ms.Name,
			Enable: ms.Enable,
		}
	}

	//新增 打卡模块 判断
	switch appId {
	case "1001":
		u, err := orm.GetUserInfoById(userId)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		uInfo, err := work.GetUser(appId, u.Uid)
		if err != nil {
			logAdmin.Warn("GetUser from work service", "err", err.Error())
		}
		if err == nil && uInfo != nil && uInfo.EnterpriseCode == "FZM0001" {
			models[2] = types.Module{
				Type:   2,
				Name:   "work",
				Enable: true,
			}
		}
	}
	ret := make([]*types.Module, 0)
	for _, ms := range models {
		ret = append(ret, &ms)
	}
	return map[string]interface{}{
		"modules": ret,
	}, nil
}

//----------------版本数据--------------//
func GetAppVersions(appId string) (interface{}, error) {
	versions, err := orm.GetVersions(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	ret := make(map[string]interface{})
	ret["list"] = versions
	return ret, nil
}

func VersionGraph(appId string, startTime, endTime int64, typeList []int, versions []string) (interface{}, error) {
	typeMapArry := make([]interface{}, 0)
	for _, tp := range typeList {
		versionMapArry := make([]interface{}, 0)
		for _, version := range versions {
			splitList := SplitDay(startTime, endTime)
			data := make([]int64, 0)

			for _, oneDay := range splitList {
				if oneDay == nil {
					continue
				}
				dayStart, dayEnd := oneDay.Start, oneDay.End

				var number int64 = 0
				switch tp {
				case types.CountNew:
					number, _ = orm.IncreaseUsersAsVersion(appId, dayStart, dayEnd, version)
				case types.CountActive:
					number, _ = orm.GetActiveNumberAsVersion(appId, dayStart, dayEnd, version)
				case types.CountOpen:
					number, _ = orm.GetOpenNumberAsVersion(appId, dayStart, dayEnd, version)
				default:
					continue
				}
				data = append(data, number)
			}

			versionItem := make(map[string]interface{})
			versionItem["version"] = version
			versionItem["data"] = data

			versionMapArry = append(versionMapArry, versionItem)
		}
		typeItem := make(map[string]interface{})
		typeItem["type"] = tp
		typeItem["versionMap"] = versionMapArry

		typeMapArry = append(typeMapArry, typeItem)
	}
	dataMap := make(map[string]interface{})
	dataMap["typeMap"] = typeMapArry

	return dataMap, nil
}

//单个版本明细
func AppSpecificVersionDetails(appId, version string, startTime, endTime int64, count, page int) (interface{}, error) {
	splitList := SplitDay(startTime, endTime)
	total := len(splitList)

	list := make([]interface{}, 0)
	for i := count * (page - 1); i < count*page && total > i; i++ {
		oneDay := splitList[len(splitList)-i-1] //倒序
		dayStart, dayEnd := oneDay.Start, oneDay.End

		new, _ := orm.IncreaseUsersAsVersion(appId, dayStart, dayEnd, version)
		newUsers, _ := orm.IncreaseUsersInfoAsVersion(appId, dayStart, dayEnd, version)
		open, _ := orm.GetOpenNumberAsVersion(appId, dayStart, dayEnd, version)
		active, _ := orm.GetActiveNumberAsVersion(appId, dayStart, dayEnd, version)
		activeUsers, _ := orm.GetActiveUsersInfoAsVersion(appId, dayStart, dayEnd, version)

		newActive := 0
		newUsermap := make(map[string]bool)
		for _, v := range newUsers {
			newUsermap[v.UserId] = true
		}
		for _, v := range activeUsers {
			if _, ok := newUsermap[v.UserId]; ok {
				newActive++
			}
		}
		var activeNewPercent float64 = 0
		if active != 0 {
			activeNewPercent = float64(newActive) / float64(active)
		}

		item := make(map[string]interface{})
		item["datetime"] = dayStart
		item["new"] = new
		item["open"] = open
		item["active"] = active
		item["activeNewPercent"] = activeNewPercent

		list = append(list, item)
	}

	ret := make(map[string]interface{})
	ret["total"] = total
	ret["list"] = list
	return ret, nil
}

//导出单个版本明细
func ExportAppSpecificVersionDetails(appId, version string, startTime, endTime int64) (string, error) {
	splitList := SplitDay(startTime, endTime)

	content := make([][]interface{}, 0)
	title := make([]interface{}, 4)
	title[0] = "时间"
	title[1] = "新增用户"
	title[2] = "活跃用户(新用户占比)"
	title[3] = "启动次数"
	content = append(content, title)

	for i := 0; i < len(splitList); i++ {
		oneDay := splitList[len(splitList)-i-1] //倒序
		dayStart, dayEnd := oneDay.Start, oneDay.End

		new, _ := orm.IncreaseUsersAsVersion(appId, dayStart, dayEnd, version)
		newUsers, _ := orm.IncreaseUsersInfoAsVersion(appId, dayStart, dayEnd, version)
		open, _ := orm.GetOpenNumberAsVersion(appId, dayStart, dayEnd, version)
		active, _ := orm.GetActiveNumberAsVersion(appId, dayStart, dayEnd, version)
		activeUsers, _ := orm.GetActiveUsersInfoAsVersion(appId, dayStart, dayEnd, version)

		newActive := 0
		newUsermap := make(map[string]bool)
		for _, v := range newUsers {
			newUsermap[v.UserId] = true
		}
		for _, v := range activeUsers {
			if _, ok := newUsermap[v.UserId]; ok {
				newActive++
			}
		}
		var activeNewPercent float64 = 0
		if active != 0 {
			activeNewPercent = float64(newActive) / float64(active)
		}

		item := make([]interface{}, 4)
		item[0] = utility.MillionSecondToTimeString(dayStart)
		item[1] = new
		item[2] = fmt.Sprintf("%v(%v%%)", active, math.Floor(activeNewPercent*100+0.5))
		item[3] = open
		content = append(content, item)

	}

	filename := fmt.Sprintf("%s-%vVersion-%v-%v", appId, version, startTime, endTime)

	url, err := WriteExcel(types.ExcelAddr, filename, content)
	if err != nil {
		logAdmin.Error("write excel failed", "err", err.Error(), "ExcelAddr", types.ExcelAddr, "filename", filename)
		return "", result.NewError(result.AdminExportFailed)
	}
	return url, err
}

func AppVersionDetails(appId string, startTime, endTime int64, count, currentPage int, sortField string, sortRule int) (interface{}, error) {
	versions, err := orm.GetVersions(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	_, todayEnd := utility.DayStartAndEndNowMillionSecond(time.Now())
	totalLogs := len(versions)
	if count*currentPage > totalLogs {
		versions = versions[(currentPage-1)*count:]
	} else {
		versions = versions[(currentPage-1)*count : currentPage*count]
	}

	list := make([]interface{}, 0)
	for i := 0; i < len(versions); i++ {
		version := versions[len(versions)-i-1] //倒序
		//所有版本累计用户 用于计算百分比
		allVersionTotal, _ := orm.GrandTotalUsers(appId, todayEnd)
		total, _ := orm.GrandTotalUsersAsVersion(appId, todayEnd, version)
		var totalPercent float64 = 0
		if allVersionTotal != 0 {
			totalPercent = float64(total) / float64(allVersionTotal)
		}
		new, _ := orm.IncreaseUsersAsVersion(appId, startTime, endTime, version)
		newUsers, _ := orm.IncreaseUsersInfoAsVersion(appId, startTime, endTime, version)
		open, _ := orm.GetOpenNumberAsVersion(appId, startTime, endTime, version)
		active, _ := orm.GetActiveNumberAsVersion(appId, startTime, endTime, version)
		activeUsers, _ := orm.GetActiveUsersInfoAsVersion(appId, startTime, endTime, version)
		update, _ := orm.UpdateVersionUsersNumber(appId, version)

		newActive := 0
		newUsermap := make(map[string]bool)
		for _, v := range newUsers {
			newUsermap[v.UserId] = true
		}
		for _, v := range activeUsers {
			if _, ok := newUsermap[v.UserId]; ok {
				newActive++
			}
		}
		var activeNewPercent float64 = 0
		if active != 0 {
			activeNewPercent = float64(newActive) / float64(active)
		}

		item := make(map[string]interface{})
		item["version"] = version
		item["total"] = total
		item["new"] = new
		item["open"] = open
		item["active"] = active
		item["update"] = update
		item["totalPercent"] = totalPercent
		item["activeNewPercent"] = activeNewPercent

		list = append(list, item)
	}

	if sortField != "" {
		desc := false
		if sortRule == 1 {
			desc = true
		}
		//排序
		utility.OrderedBy(versionInfoSort[sortField], desc).Sort(list)
	}

	ret := make(map[string]interface{})
	ret["total"] = totalLogs
	ret["list"] = list
	return ret, nil
}

//导出版本统计明细
func ExportAppVersionDetails(appId string, startTime, endTime int64, sortField string, sortRule int) (string, error) {
	versions, err := orm.GetVersions(appId)
	if err != nil {
		return "", result.NewError(result.DbConnectFail)
	}
	_, todayEnd := utility.DayStartAndEndNowMillionSecond(time.Now())

	content := make([][]interface{}, 0)
	title := make([]interface{}, 6)
	title[0] = "版本"
	title[1] = "截至今日版本累计用户(%)"
	title[2] = "新增用户"
	title[3] = "活跃用户(新用户占比)"
	title[4] = "启动次数"
	title[5] = "升级用户"
	content = append(content, title)

	list := make([]interface{}, 0)
	for i := 0; i < len(versions); i++ {
		version := versions[len(versions)-i-1] //倒序
		//所有版本累计用户 用于计算百分比
		allVersionTotal, _ := orm.GrandTotalUsers(appId, todayEnd)
		total, _ := orm.GrandTotalUsersAsVersion(appId, todayEnd, version)
		var totalPercent float64 = 0
		if allVersionTotal != 0 {
			totalPercent = float64(total) / float64(allVersionTotal)
		}
		new, _ := orm.IncreaseUsersAsVersion(appId, startTime, endTime, version)
		newUsers, _ := orm.IncreaseUsersInfoAsVersion(appId, startTime, endTime, version)
		open, _ := orm.GetOpenNumberAsVersion(appId, startTime, endTime, version)
		active, _ := orm.GetActiveNumberAsVersion(appId, startTime, endTime, version)
		activeUsers, _ := orm.GetActiveUsersInfoAsVersion(appId, startTime, endTime, version)
		update, _ := orm.UpdateVersionUsersNumber(appId, version)

		newActive := 0
		newUsermap := make(map[string]bool)
		for _, v := range newUsers {
			newUsermap[v.UserId] = true
		}
		for _, v := range activeUsers {
			if _, ok := newUsermap[v.UserId]; ok {
				newActive++
			}
		}
		var activeNewPercent float64 = 0
		if active != 0 {
			activeNewPercent = float64(newActive) / float64(active)
		}

		item := make(map[string]interface{})
		item["version"] = version
		item["total"] = total
		item["new"] = new
		item["open"] = open
		item["active"] = active
		item["update"] = update
		item["totalPercent"] = totalPercent
		item["activeNewPercent"] = activeNewPercent

		list = append(list, item)
	}
	sortName := ""
	if sortField != "" {
		desc := false
		sortName = "aes"
		if sortRule == 1 {
			desc = true
			sortName = "desc"
		}
		//排序
		utility.OrderedBy(versionInfoSort[sortField], desc).Sort(list)
	}

	for _, v := range list {
		val := v.(map[string]interface{})
		item := make([]interface{}, 6)
		item[0] = val["version"]
		item[1] = fmt.Sprintf("%v(%v%%)", val["total"], math.Floor(utility.ToFloat64(val["totalPercent"])*100+0.5))
		item[2] = val["new"]
		item[3] = fmt.Sprintf("%v(%v%%)", val["active"], math.Floor(utility.ToFloat64(val["activeNewPercent"])*100+0.5))
		item[4] = val["open"]
		item[5] = val["update"]
		content = append(content, item)
	}

	filename := fmt.Sprintf("%s-allVersion-%v-%v-%v-%v", appId, sortField, sortName, startTime, endTime)

	url, err := WriteExcel(types.ExcelAddr, filename, content)
	if err != nil {
		logAdmin.Error("write excel failed", "err", err.Error(), "ExcelAddr", types.ExcelAddr, "filename", filename)
		return "", result.NewError(result.AdminExportFailed)
	}
	return url, err
}

//单个版本明细 根据平台
func AppSpecificVersionAsPlatform(appId, version string, startTime, endTime int64) (interface{}, error) {
	devices := []string{types.DeviceAndroid, types.DeviceIOS, types.DevicePC}
	list := make([]interface{}, 0)

	for _, device := range devices {
		new, _ := orm.IncreaseUsersAsVersionAndPlatform(appId, startTime, endTime, version, device)
		newUsers, _ := orm.IncreaseUsersInfoAsVersionAndPlatform(appId, startTime, endTime, version, device)
		open, _ := orm.GetOpenNumberAsVersionAndPlatform(appId, startTime, endTime, version, device)
		active, _ := orm.GetActiveNumberAsVersionAndPlatform(appId, startTime, endTime, version, device)
		activeUsers, _ := orm.GetActiveUsersInfoAsVersionAndPlatform(appId, startTime, endTime, version, device)

		newActive := 0
		newUsermap := make(map[string]bool)
		for _, v := range newUsers {
			newUsermap[v.UserId] = true
		}
		for _, v := range activeUsers {
			if _, ok := newUsermap[v.UserId]; ok {
				newActive++
			}
		}
		var activeNewPercent float64 = 0
		if active != 0 {
			activeNewPercent = float64(newActive) / float64(active)
		}

		item := make(map[string]interface{})
		item["platformName"] = device
		item["new"] = new
		item["open"] = open
		item["active"] = active
		item["activeNewPercent"] = activeNewPercent

		list = append(list, item)
	}

	ret := make(map[string]interface{})
	ret["platformMap"] = list
	return ret, nil
}

var versionInfoSort = map[string]func(c1, c2 interface{}) bool{
	"total":  sortTotal,
	"new":    sortNew,
	"active": sortActive,
	"open":   sortOpen,
	"update": sortUpdate,
}

// Closures that order the Change structure.
var sortTotal = func(c1, c2 interface{}) bool {
	t1 := c1.(map[string]interface{})
	t2 := c2.(map[string]interface{})
	return utility.ToInt64(t1["total"]) < utility.ToInt64(t2["total"])
}
var sortNew = func(c1, c2 interface{}) bool {
	t1 := c1.(map[string]interface{})
	t2 := c2.(map[string]interface{})
	return utility.ToInt64(t1["new"]) < utility.ToInt64(t2["new"])
}
var sortActive = func(c1, c2 interface{}) bool {
	t1 := c1.(map[string]interface{})
	t2 := c2.(map[string]interface{})
	return utility.ToInt64(t1["active"]) < utility.ToInt64(t2["active"])
}
var sortOpen = func(c1, c2 interface{}) bool {
	t1 := c1.(map[string]interface{})
	t2 := c2.(map[string]interface{})
	return utility.ToInt64(t1["open"]) < utility.ToInt64(t2["open"])
}
var sortUpdate = func(c1, c2 interface{}) bool {
	t1 := c1.(map[string]interface{})
	t2 := c2.(map[string]interface{})
	return utility.ToInt64(t1["update"]) < utility.ToInt64(t2["update"])
}

func FindExcel(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteExcel(path, filename string, content [][]interface{}) (string, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}

	for row := 0; row < len(content); row++ {
		for col := 0; col < len(content[row]); col++ {
			cell := sheet.Cell(row, col)
			cell.Value = utility.ToString(content[row][col])
		}
	}

	err = file.Save(path + filename + ".xlsx")
	if err != nil {
		return "", err
	}
	return filename + ".xlsx", nil
}

//-------------------开屏广告------------------//
func CreateAd(appId, name, url string, duration int, link string, isActive int) (interface{}, error) {
	count, err := orm.GetAdNumbers(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if count >= types.MaxAdNumbers {
		logAdmin.Debug("CreateAd max number limited", "max", types.MaxAdNumbers, "current", count)
		return nil, result.NewError(result.AdminAdNumbersLimit)
	}

	id, err := orm.InsertAd(appId, name, url, duration, link, isActive)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	ret := make(map[string]string)
	ret["id"] = utility.ToString(id)
	return ret, err
}

//客户端获取开屏广告
func Advertisement(appId string) (interface{}, error) {
	rlt, err := orm.GetAllAd(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	ads := make([]*types.Advertisement, 0)
	for _, a := range rlt {
		if a.IsActive == types.AdIsActive {
			ads = append(ads, a)
		}
	}

	index := 0
	if len(ads) > 0 {
		//以时间作为初始化种子
		rand.Seed(time.Now().UnixNano())
		index = rand.Intn(len(ads))
	}

	ret := make(map[string]interface{})
	others := make([]string, 0)
	for i, a := range ads {
		if i == index {
			ret = utility.StringToJobj(a)
		} else {
			others = append(others, a.Url)
		}
	}
	ret["others"] = others
	return ret, nil
}

func GetAllAd(appId string) (interface{}, error) {
	ads, err := orm.GetAllAd(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	ret := make(map[string]interface{})
	if len(ads) < 1 {
		ret["list"] = make([]string, 0)
	} else {
		ret["list"] = ads
	}
	return ret, nil
}

func SetAdName(id, name string) error {
	err := orm.SetAdName(id, name)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func ActiveAd(id string, isActive int) error {
	err := orm.ActiveAd(id, isActive)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func DeleteAd(id string) error {
	err := orm.DeleteAd(id)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//创建文件，返回文件路径
func CreateAdFile(filename string, file *multipart.File) (string, error) {
	//写入文件
	path := cfg.FileStore.FilePath
	out, err := os.Create(path + "/" + filename)
	if err != nil {
		logAdmin.Error("CreateAdFile:can not create file", "err", err.Error(), "filename", filename, "path", path)
		return "", result.NewError(result.AdUploadFailed)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			logAdmin.Error("CreateAdFile:can not close file", "err", err.Error(), "filename", filename, "path", path)
		}
	}()

	_, err = io.Copy(out, *file)
	if err != nil {
		logAdmin.Error("CreateAdFile:can not copy file", "err", err.Error(), "filename", filename, "path", path)
		return "", result.NewError(result.AdUploadFailed)
	}

	url := fmt.Sprintf("%s/%s", cfg.FileStore.Path, filename)
	return url, nil
}

//---------------------------资金管理---------------//

func EditReward(appId string, params *account.EditRewardParam) error {
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		return account.EditRealBcoin(app.AppKey, app.AppSecret, app.AccountServer, params)
	default:
		return types.ERR_NOTSUPPORT
	}
}

func ShowReward(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		return account.ShowRealBcoin(app.AppKey, app.AppSecret, app.AccountServer)
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

func RewardRule(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		rlt, err := account.ShowRealBcoin(app.AppKey, app.AppSecret, app.AccountServer)
		if err != nil {
			return nil, err
		}
		ret := make(map[string]interface{})
		list := make([]interface{}, 0)

		data := rlt.(map[string]interface{})
		if base, ok := data["base"].(map[string]interface{}); ok {
			item := make(map[string]interface{})
			item["type"] = 1
			item["ext"] = map[string]string{}
			item["newUser"] = map[string]string{
				"amount":   utility.ToString(base["rewardForUser"]),
				"currency": utility.ToString(base["currency"]),
			}
			inviterAmount := "0"
			if rw, ok := base["rewardForInviter"].([]interface{}); ok && len(rw) > 0 {
				inviterAmount = utility.ToString(rw[0])
			}
			item["inviter"] = map[string]string{
				"amount":   inviterAmount,
				"currency": utility.ToString(base["currency"]),
			}
			list = append(list, item)
		}

		if advance, ok := data["advance"].(map[string]interface{}); ok {
			item := make(map[string]interface{})
			item["type"] = 2
			item["ext"] = map[string]string{
				"reachNum": utility.ToString(advance["reachNum"]),
			}
			item["newUser"] = map[string]string{
				"amount":   utility.ToString(advance["rewardForUser"]),
				"currency": utility.ToString(advance["currency"]),
			}
			item["inviter"] = map[string]string{
				"amount":   utility.ToString(advance["rewardForNum"]),
				"currency": utility.ToString(advance["currency"]),
			}
			list = append(list, item)
		}
		ret["base"] = list
		return ret, nil
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

func CoinSupport(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		return account.CoinSupport(app.AppKey, app.AppSecret, app.AccountServer)
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

func RewardStatistics(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		data, err := account.RewardList(app.AppKey, app.AppSecret, app.AccountServer, nil)
		if err != nil {
			return nil, err
		}
		info := data.(map[string]interface{})
		ret := make(map[string]interface{})
		ret["statistics"] = info["statistics"]
		return ret, err
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

func RewardList(appId string, query *account.RewardListParam) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		data, err := account.RewardList(app.AppKey, app.AppSecret, app.AccountServer, query)
		if err != nil {
			return nil, err
		}
		info := data.(map[string]interface{})
		ret := make(map[string]interface{})
		ret["count"] = info["count"]
		ret["list"] = info["data"]
		return ret, err
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

//红包手续费配置信息
func RPFeeConfig(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		data, err := account.RPFeeConfigFromBcoin(app.AppKey, app.AppSecret, app.AccountServer)
		if err != nil {
			return nil, err
		}
		ret := make(map[string]interface{})
		ret["config"] = data
		return ret, err
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

//设置红包手续费
func SetRPFeeConfig(appId string, params *account.SetRPFeeParam) error {
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		err := account.SetRPFeeConfigFromBcoin(app.AppKey, app.AppSecret, app.AccountServer, params)
		if err != nil {
			return err
		}
		return nil
	default:
		return types.ERR_NOTSUPPORT
	}
}

//设置红包手续费
func RPFeeStatistics(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, types.ERR_APPNOTFIND
	}

	switch app.IsInner {
	case types.IsInnerAccount:
		data, err := account.RPFeeStatistics(app.AppKey, app.AppSecret, app.AccountServer)
		if err != nil {
			return nil, err
		}
		ret := map[string]interface{}{
			"statistics": []struct{}{},
		}
		if data != nil {
			ret["statistics"] = data
		}
		return ret, nil
	default:
		return nil, types.ERR_NOTSUPPORT
	}
}

func UserOnlineNums() int {
	return router.GetUserNums()
}

func UserClientNums() int {
	return router.GetUserClientNums()
}
