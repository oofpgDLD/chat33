package api

import (
	"fmt"

	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logRoom = log15.New("model", "api/admin")

// 创建群
func CreateRoom(c *gin.Context) {
	type requestParams struct {
		RoomName   string   `json:"roomName"`
		RoomAvatar string   `json:"roomAvatar"`
		Users      []string `json:"users" binding:"required"`
		Encrypt    int
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	encrypt := types.IsNotEncrypt
	approval := types.ShouldNotApproval
	recordPermission := types.CanReadAllLog
	switch params.Encrypt {
	case types.IsEncrypt:
		encrypt = params.Encrypt
		//approval = types.ShouldApproval
		recordPermission = types.CanNotReadAllLog
	default:
		encrypt = types.IsNotEncrypt
	}
	roomInfo, err := model.CreateRoom(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.RoomName, params.RoomAvatar, encrypt, types.CanAddFriend, approval, recordPermission, types.AdminNotMuted, types.MasterNotMuted, params.Users)
	c.Set(ReqResult, roomInfo)
	c.Set(ReqError, err)
}

// 删除群
func RemoveRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomid" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	if orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted) < types.RoomLevelManager {
		//master cannot logout

		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}
	err := model.RemoveRoom(c.MustGet(UserId).(string), params.RoomId)
	c.Set(ReqError, err)
}

// 退出群
func LoginOutRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomid" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//check user is master of this room
	if orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted) == types.RoomLevelMaster {
		//master cannot logout
		c.Set(ReqError, result.NewError(result.CanNotLoginOut))
		return
	}
	err := model.LoginOutRoom(c.MustGet(UserId).(string), params.RoomId)
	c.Set(ReqError, err)
}

// 踢出群
func KickOutRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string   `json:"roomId" binding:"required"`
		Users  []string `json:"users" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	if len(params.Users) < 1 {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "user number is less than one"))
		return
	}

	optLevel := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted)
	//check user is not master of this room
	if optLevel < 2 {
		//master cannot logout
		logRoom.Warn("KickOutRoom", "warn", "PermissionDeny", " operator Level", optLevel, "min level", 2)
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	for _, v := range params.Users {
		if level := orm.GetMemberLevel(params.RoomId, v, types.RoomUserNotDeleted); level >= optLevel {
			//master cannot logout
			logRoom.Warn("KickOutRoom", "warn", "PermissionDeny", " operator Level", optLevel, "target level", level)
			c.Set(ReqError, result.NewError(result.PermissionDeny))
			return
		}
	}

	for _, v := range params.Users {
		err := model.KickOutRoom(c.MustGet(UserId).(string), params.RoomId, v)
		c.Set(ReqError, err)
	}
}

// 获取群列表
func GetRoomList(c *gin.Context) {
	type requestParams struct {
		Type int `json:"type" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	rooms, err := model.GetRoomList(c.MustGet(UserId).(string), params.Type)
	c.Set(ReqError, err)
	c.Set(ReqResult, rooms)
}

// 获取群信息
func GetRoomInfo(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	roomInfo, err := model.GetRoomInfo(c.MustGet(UserId).(string), params.RoomId)
	c.Set(ReqError, err)
	c.Set(ReqResult, roomInfo)
}

// 获取群成员列表
func GetRoomUserList(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	roomUsers, err := model.GetRoomUserList(params.RoomId)
	c.Set(ReqError, err)
	c.Set(ReqResult, roomUsers)
}

// 获取群成员信息
func GetRoomUserInfo(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	roomUserInfo, err := model.GetRoomUserInfo(params.RoomId, params.UserId)
	c.Set(ReqError, err)
	c.Set(ReqResult, roomUserInfo)
}

// 搜索群成员信息
func GetRoomSearchMember(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		Query  string `json:"query" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	roomUserInfo, err := model.GetRoomSearchMember(params.RoomId, params.Query)
	c.Set(ReqError, err)
	c.Set(ReqResult, roomUserInfo)
}

// 管理员设置群
func AdminSetPermission(c *gin.Context) {
	type requestParams struct {
		RoomId           string `json:"roomId" binding:"required"`
		CanAddFriend     int    `json:"canAddFriend"`
		JoinPermission   int    `json:"joinPermission"`
		RecordPermission int    `json:"recordPermission"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//check caller is manager
	if level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted); level < types.RoomLevelManager {
		logRoom.Warn("AdminSetPermission", "warn", "PermissionDeny", " operator Level", level, "min level", types.RoomLevelManager)
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	err := model.AdminSetPermission(params.RoomId, params.CanAddFriend, params.JoinPermission, params.RecordPermission)
	c.Set(ReqError, err)
}

// 设置群名称
func SetRoomName(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		Name   string `json:"name"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	//check caller is manager
	if level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted); level < types.RoomLevelManager {
		logRoom.Warn("SetRoomName", "warn", "PermissionDeny", " operator Level", level, "min level", types.RoomLevelManager)
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	err := model.SetRoomName(c.MustGet(UserId).(string), params.RoomId, params.Name)
	c.Set(ReqError, err)
}

// 设置群头像
func SetRoomAvatar(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		Avatar string `json:"avatar" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	//check caller is manager
	if level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted); level < types.RoomLevelManager {
		logRoom.Warn("SetRoomAvatar", "warn", "PermissionDeny", " operator Level", level, "min level", types.RoomLevelManager)
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}
	err := model.SetRoomAvatar(params.RoomId, params.Avatar)
	c.Set(ReqError, err)
}

// 群内用户身份设置
func SetLevel(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
		Level  int    `json:"level" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	// check userId is admin
	if level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted); level != types.RoomLevelMaster {
		logRoom.Warn("SetRoomAvatar", "warn", "PermissionDeny", " operator Level", level, "must level", types.RoomLevelManager)
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	err := model.SetLevel(c.MustGet(UserId).(string), params.UserId, params.RoomId, params.Level)
	c.Set(ReqError, err)
}

// 群成员设置免打扰
func SetNoDisturbing(c *gin.Context) {
	type requestParams struct {
		RoomId          string `json:"roomId" binding:"required"`
		SetNoDisturbing int    `json:"setNoDisturbing" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	switch params.SetNoDisturbing {
	case 1:
	case 2:
	default:
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized type"))
		return
	}
	err := model.SetNoDisturbing(c.MustGet(UserId).(string), params.RoomId, params.SetNoDisturbing)
	c.Set(ReqError, err)
}

// 群成员设置消息置顶
func SetStickyOnTop(c *gin.Context) {
	type requestParams struct {
		RoomId      string `json:"roomId" binding:"required"`
		StickyOnTop int    `json:"stickyOnTop" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))

		return
	}

	switch params.StickyOnTop {
	case 1:
	case 2:
	default:
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized type"))
		return
	}
	err := model.SetStickyOnTop(c.MustGet(UserId).(string), params.RoomId, params.StickyOnTop)
	c.Set(ReqError, err)
}

// 群成员设置群内昵称
func SetMemberNickname(c *gin.Context) {
	type requestParams struct {
		RoomId   string `json:"roomId" binding:"required"`
		Nickname string `json:"nickname"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetMemberNickname(c.MustGet(UserId).(string), params.RoomId, params.Nickname)
	c.Set(ReqError, err)
}

// 邀请入群
func JoinRoomInvite(c *gin.Context) {
	type requestParams struct {
		RoomId string   `json:"roomId" binding:"required"`
		Users  []string `json:"users" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.JoinRoomInvite(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.RoomId, params.Users)
	c.Set(ReqResult, ret)
	c.Set(ReqError, err)
}

// 入群申请
func JoinRoomApply(c *gin.Context) {
	type requestParams struct {
		RoomId      string `json:"roomId" binding:"required"`
		ApplyReason string `json:"applyReason"`
		SourceType  int    `json:"sourceType"`
		SourceId    string `json:"sourceId"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	_, err := model.JoinRoomApply(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.RoomId, params.ApplyReason, params.SourceType, params.SourceId)
	c.Set(ReqError, err)
}

// 批量入群申请
func BatchJoinRoomApply(c *gin.Context) {
	var params struct {
		Rooms []string `json:"rooms" binding:"required"`
	}
	SourceType := types.Search
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret := make(map[string]interface{})
	rooms := make([]string, 0)
	for _, roomId := range params.Rooms {
		b, err := model.JoinRoomApply(c.MustGet(AppId).(string), c.MustGet(UserId).(string),
			roomId, "", SourceType, "")
		if err != nil {
			logRoom.Warn("BatchJoinRoomApply err", "err", err.(*result.Error).ErrorCode)
		}
		if b {
			rooms = append(rooms, roomId)
		}
	}
	ret["rooms"] = rooms
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

// 入群申请处理
func JoinRoomApprove(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
		UserId string `json:"userId" binding:"required"`
		Agree  int    `json:"agree" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.JoinRoomApprove(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.RoomId, params.UserId, params.Agree)
	c.Set(ReqError, err)
}

// 获取消息记录
func GetRoomChatLog(c *gin.Context) {
	type requestParams struct {
		RoomId  string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	logs, err := model.GetRoomChatLog(c.MustGet(UserId).(string), params.RoomId, params.StartId, params.Number)
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

// 获取群文件列表
func GetRoomFileList(c *gin.Context) {
	type requestParams struct {
		RoomId  string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number"`
		Query   string `json:"query"`
		Owner   string `json:"owner"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	logs, err := model.GetTypicalMsgLogs(c.MustGet(UserId).(string), params.RoomId, params.StartId, params.Owner, params.Query, params.Number, []string{utility.ToString(types.File)})
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

// 获取群图片列表
func GetRoomPhotoList(c *gin.Context) {
	type requestParams struct {
		RoomId  string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	logs, err := model.GetTypicalMsgLogs(c.MustGet(UserId).(string), params.RoomId, params.StartId, "", "", params.Number, []string{utility.ToString(types.Photo), utility.ToString(types.Video)})
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

// 获取群在线人数
func GetOnlineNumber(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	info, err := model.GetRoomOnlineNumber(params.RoomId)
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 设置群中禁言
func SetRoomMuted(c *gin.Context) {
	type requestParams struct {
		RoomId   string   `json:"roomId" binding:"required"`
		ListType int      `json:"listType" binding:"required"`
		Users    []string `json:"users"`
		Deadline int64    `json:"deadline"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted)
	if level != types.RoomLevelManager && level != types.RoomLevelMaster {
		logRoom.Warn("SetRoomAvatar", "warn", "PermissionDeny", " operator Level", level, "must level", fmt.Sprintf("%d or %d", types.RoomLevelManager, types.RoomLevelMaster))
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	switch params.ListType {
	case types.AllSpeak:
	case types.Blacklist:
		// 黑名单禁言时长不可为空
		if params.Deadline == 0 {
			c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized type"))
			return
		}
	case types.Whitelist:
	case types.AllMuted:
	default:
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized type"))
		return
	}
	var userMap = make(map[string]bool)
	// 管理员和群主不可被加入白名单或者黑名单
	for _, v := range params.Users {
		// check hava master or manager
		level := orm.GetMemberLevel(params.RoomId, v, types.RoomUserNotDeleted)
		switch level {
		case types.RoomLevelNotExist:
			c.Set(ReqError, result.NewError(result.UserIsNotInRoom))
			return
		case types.RoomLevelMaster:
			c.Set(ReqError, result.NewError(result.ObjIsMaster))
			return
		case types.RoomLevelManager:
			c.Set(ReqError, result.NewError(result.ObjIsManager))
			return
		}
		userMap[v] = true
	}

	err := model.SetRoomMuted(c.MustGet(UserId).(string), params.RoomId, params.ListType, userMap, params.Deadline)
	c.Set(ReqError, err)
}

//单个禁言
func SetRoomMutedSingle(c *gin.Context) {
	type requestParams struct {
		RoomId   string `json:"roomId" binding:"required"`
		UserId   string `json:"userId" binding:"required"`
		Deadline int64  `json:"deadline"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted)
	if level != types.RoomLevelManager && level != types.RoomLevelMaster {
		logRoom.Warn("SetRoomAvatar", "warn", "PermissionDeny", " operator Level", level, "must level", types.RoomLevelManager)
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	objLevel := orm.GetMemberLevel(params.RoomId, params.UserId, types.RoomUserNotDeleted)
	switch objLevel {
	case types.RoomLevelNotExist:
		c.Set(ReqError, result.NewError(result.UserIsNotInRoom))
		return
	case types.RoomLevelMaster:
		c.Set(ReqError, result.NewError(result.ObjIsMaster))
		return
	case types.RoomLevelManager:
		c.Set(ReqError, result.NewError(result.ObjIsManager))
		return
	}

	err := model.SetRoomMutedSingle(c.MustGet(UserId).(string), params.RoomId, params.UserId, params.Deadline)
	c.Set(ReqError, err)
}

// 获取群公告
func GetSystemMsg(c *gin.Context) {
	type requestParams struct {
		RoomId  string `json:"roomId" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	info, err := model.GetSystemMsg(params.RoomId, utility.ToInt64(params.StartId), params.Number)
	c.Set(ReqError, err)
	c.Set(ReqResult, info)
}

// 发布群公告
func SetSystemMsg(c *gin.Context) {
	type requestParams struct {
		RoomId  string `json:"roomId" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted)
	if level != types.RoomLevelManager && level != types.RoomLevelMaster {
		logRoom.Warn("SetRoomAvatar", "warn", "PermissionDeny", " operator Level", level, "must level", fmt.Sprintf("%d or %d", types.RoomLevelManager, types.RoomLevelMaster))
		c.Set(ReqError, result.NewError(result.PermissionDeny))
		return
	}

	err := model.SetSystemMsg(c.MustGet(UserId).(string), params.RoomId, params.Content)
	c.Set(ReqError, err)
}

//判断用户是否在群里
func UserIsInRoom(c *gin.Context) {
	type requestParams struct {
		RoomId string `json:"roomId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	isInRoom, err := model.UserIsInRoom(c.MustGet(UserId).(string), params.RoomId)
	c.Set(ReqError, err)
	c.Set(ReqResult, isInRoom)
}

//获取推荐群列表
func RecommendRooms(c *gin.Context) {
	var params struct {
		Number int `json:"number"`
		Times  int `json:"times" binding:"required"`
	}
	params.Number = 6

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.GetRecommendRooms(c.MustGet(AppId).(string), params.Number, params.Times)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//群认证申请
func VerifyApplyForRoom(c *gin.Context) {
	var params struct {
		RoomId      string `json:"roomId" binding:"required"`
		Description string `json:"description" binding:"required"`
		PayPassword string `json:"payPassword" binding:"required"`
		Code        string `json:"code"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//检查是否是群主
	if level := orm.GetMemberLevel(params.RoomId, c.MustGet(UserId).(string), types.RoomUserNotDeleted); level != types.RoomLevelMaster {
		logRoom.Warn("SetRoomAvatar", "warn", "PermissionDeny", " operator Level", level, "must level", types.RoomLevelMaster)
		c.Set(ReqError, result.NewError(result.PermissionDeny).SetExtMessage("仅限群主能认证该群"))
		return
	}

	_, err := model.VerifyApply(c.MustGet(AppId).(string), c.MustGet(Token).(string), types.VerifyForRoom, params.RoomId, params.Description, params.PayPassword, params.Code)
	c.Set(ReqError, err)
	//c.Set(ReqResult, ret)
}

//认证审核状态
func RoomVerificationState(c *gin.Context) {
	var params struct {
		RoomId string `json:"roomId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	state, err := model.CheckRoomVerifyState(c.MustGet(UserId).(string), params.RoomId)
	ret := map[string]interface{}{
		"state": state,
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}
