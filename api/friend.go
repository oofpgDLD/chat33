package api

import (
	"encoding/json"
	"fmt"

	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/gin-gonic/gin"
)

//获取好友列表
func FriendList(c *gin.Context) {
	type friendListParams struct {
		Type int   `json:"type"`
		Time int64 `json:"time"`
	}

	var params friendListParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	friends, err := model.FriendList(c.MustGet(UserId).(string), params.Type, params.Time)
	c.Set(ReqError, err)
	c.Set(ReqResult, friends)
}

//申请添加好友
func AddFriend(c *gin.Context) {
	type addFriendParams struct {
		Id         string `json:"id" binding:"required"`
		Remark     string `json:"remark"`
		Reason     string `json:"reason"`
		SourceId   string `json:"sourceId"`
		SourceType int    `json:"sourceType"`
		Answer     string `json:"answer"`
	}
	var params addFriendParams
	if err := c.ShouldBindJSON(&params); err != nil || (params.SourceType != 1 && params.SourceType != 2 && params.SourceType != 3 && params.SourceType != 4 && params.SourceType != 5) {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	state, err := model.AddFriend(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.Id, params.Remark, params.Reason, params.SourceId, params.SourceType, params.Answer)
	var retMsg = make(map[string]int)
	retMsg["state"] = state

	c.Set(ReqError, err)
	c.Set(ReqResult, retMsg)
}

//处理好友申请
func HandleFriendRequest(c *gin.Context) {
	type addFriendParams struct {
		Id    string `json:"id" binding:"required"`
		Agree *int   `json:"agree" binding:"required"`
	}

	var params addFriendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	if (*params.Agree != types.AcceptRequest) && (*params.Agree != types.RejectRequest) {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, fmt.Sprintf("unrecognized type :%d", *params.Agree)))
		return
	}
	err := model.HandleFriendRequest(c.MustGet(UserId).(string), params.Id, *params.Agree)
	c.Set(ReqError, err)
}

//设置备注名
func FriendSetRemark(c *gin.Context) {
	type setRemarkParams struct {
		Id     string  `json:"id" binding:"required"`
		Remark *string `json:"remark" binding:"required"`
	}
	var params setRemarkParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetFriendRemark(c.MustGet(UserId).(string), params.Id, *params.Remark)
	c.Set(ReqError, err)
}

//设置详细备注
func FriendSetExtRemark(c *gin.Context) {
	type params struct {
		Id         string  `json:"id" binding:"required"`
		Remark     *string `json:"remark" binding:"required"`
		Telephones []struct {
			Remark *string `json:"remark"`
			Phone  string  `json:"phone"`
		} `json:"telephones"`
		Description *string  `json:"description"`
		Pictures    []string `json:"pictures"`
		Encrypt     string   `json:"encrypt"`
	}
	var p params
	if err := c.ShouldBindJSON(&p); err != nil {
		logUser.Debug("FriendSetExtRemark Params warn", "err", err)
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	remark := &types.ExtRemark{
		Telephones:  p.Telephones,
		Description: p.Description,
		Pictures:    p.Pictures,
		Encrypt:     p.Encrypt,
	}
	data, err := json.Marshal(remark)
	if err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	err = model.SetFriendRemark(c.MustGet(UserId).(string), p.Id, *p.Remark)
	if err != nil {
		c.Set(ReqError, err)
		return
	}
	err = model.SetFriendExtRemark(c.MustGet(UserId).(string), p.Id, string(data))
	logUser.Debug("test data", "data", string(data))
	c.Set(ReqError, err)
}

/*
	设置好友免打扰
*/
func FriendSetDND(c *gin.Context) {
	type setDNDParams struct {
		Id              string `json:"id" binding:"required"`
		SetNoDisturbing *int   `json:"setNoDisturbing" binding:"required"`
	}
	var params setDNDParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetFriendDND(c.MustGet(UserId).(string), params.Id, *params.SetNoDisturbing)
	c.Set(ReqError, err)
}

/*
	设置好友置顶
*/
func FriendSetTop(c *gin.Context) {
	type setTopParams struct {
		Id  string `json:"id" binding:"required"`
		Top *int   `json:"stickyOnTop" binding:"required"`
	}
	var params setTopParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetFriendTop(c.MustGet(UserId).(string), params.Id, *params.Top)
	c.Set(ReqError, err)
}

/*
	删除好友
*/
func DeleteFriend(c *gin.Context) {
	type deleteFriendParams struct {
		Id string `json:"id" binding:"required"`
	}
	var params deleteFriendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.DeleteFriend(c.MustGet(UserId).(string), params.Id)
	c.Set(ReqError, err)
}

//查看好友详情
func FriendInfo(c *gin.Context) {
	var params struct {
		Id string `json:"id"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	if params.Id == "" {
		data, err := model.FriendInfo(c.MustGet(UserId).(string), c.MustGet(UserId).(string))
		c.Set(ReqError, err)
		c.Set(ReqResult, data)
		return
	} else {
		data, err := model.FriendInfo(c.MustGet(UserId).(string), params.Id)
		c.Set(ReqError, err)
		c.Set(ReqResult, data)
		return
	}
}

func UserInfo(c *gin.Context) {
	var params struct {
		Uids []string `json:"uids"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	data, err := model.UserListByUids(c.MustGet(UserId).(string), c.MustGet(AppId).(string), params.Uids)
	c.Set(ReqError, err)
	c.Set(ReqResult, data)
	return
}

//通过UID查询好友详情
func UserListByUid(c *gin.Context) {
	var params struct {
		Uid string `json:"uid"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	data, err := model.UserListByUid(c.MustGet(UserId).(string), c.MustGet(AppId).(string), params.Uid)
	c.Set(ReqError, err)
	c.Set(ReqResult, data)
	return
}

//获取好友消息记录
func FriendPhotosLog(c *gin.Context) {
	type CatLogPara struct {
		Id      string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number" binding:"required"`
	}
	var params CatLogPara
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	data, err := model.FindTypicalChatLog(c.MustGet(UserId).(string), params.Id, "", "", params.StartId, params.Number, []string{utility.ToString(types.Photo), utility.ToString(types.Video)})
	if data == nil {
		var res = make(map[string]interface{})
		var logs = make([]string, 0)
		res["logs"] = logs
		res["nextLog"] = "-1"
		data = res
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, data)
}

//获取好友消息记录
func FriendFliesLog(c *gin.Context) {
	type CatLogPara struct {
		Id      string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number" binding:"required"`
		Query   string `json:"query"`
		Owner   string `json:"owner"`
	}
	var params CatLogPara
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	data, err := model.FindTypicalChatLog(c.MustGet(UserId).(string), params.Id, params.Owner, params.Query, params.StartId, params.Number, []string{utility.ToString(types.File)})
	if data == nil {
		var res = make(map[string]interface{})
		var logs = make([]string, 0)
		res["logs"] = logs
		res["nextLog"] = "-1"
		data = res
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, data)
}

//获取好友消息记录
func FriendChatLog(c *gin.Context) {
	type ChatLogPara struct {
		Id      string `json:"id" binding:"required"`
		StartId string `json:"startId"`
		Number  int    `json:"number" binding:"required"`
	}
	var params ChatLogPara
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	data, err := model.FindCatLog(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.Id, params.StartId, params.Number)
	if data == nil {
		var res = make(map[string]interface{})
		var logs = make([]string, 0)
		res["logs"] = logs
		res["nextLog"] = "-1"
		data = res
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, data)
}

//GetAllFriendUnreadMsg获取所有好友未读消息统计
func GetAllFriendUnreadMsg(c *gin.Context) {
	data, err := model.GetAllFriendUnreadMsg1(c.MustGet(UserId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, data)
}

//判断是否是好友
func IsFriend(c *gin.Context) {
	type friendParams struct {
		FriendId string `json:"friendId" binding:"required"`
	}
	var params friendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.ParamsError)
		return
	}

	isFreind, err := model.IsFriend(c.MustGet(UserId).(string), params.FriendId)
	c.Set(ReqError, err)
	c.Set(ReqResult, isFreind)
}

//
func PrintScreen(c *gin.Context) {
	type requestParams struct {
		FriendId string `json:"id" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.ParamsError)
		return
	}

	err := model.SendPrintScreen(c.MustGet(UserId).(string), params.FriendId)
	c.Set(ReqError, err)
}

//设置答案验证
func Question(c *gin.Context) {
	//tp = 0  修改问题答案  需要同时传入问题答案
	//tp = 1  设置需要回答问题  需要同时传入问题答案
	//tp = 2  设置不需要回答问题
	type requestParams struct {
		Tp       int    `json:"tp"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
	var params requestParams
	err := c.ShouldBindJSON(&params)
	if err != nil || ((params.Tp == 0 || params.Tp == 1) && (params.Answer == "" || params.Question == "")) || (params.Tp != 0 && params.Tp != 1 && params.Tp != 2) {
		c.Set(ReqError, result.ParamsError)
		return
	}
	err = model.Question(c.MustGet(UserId).(string), params.Tp, params.Question, params.Answer)
	c.Set(ReqError, err)
}

//是否需要验证
func Confirm(c *gin.Context) {
	type requestParams struct {
		Tp int `json:"tp" binding:"required"`
	}
	var params requestParams
	err := c.ShouldBindJSON(&params)
	if err != nil || (params.Tp != 1 && params.Tp != 2) {
		c.Set(ReqError, result.ParamsError)
		return
	}
	err = model.Confirm(c.MustGet(UserId).(string), params.Tp)
	c.Set(ReqError, err)
}

//验证答案是否正确
func CheckAnswer(c *gin.Context) {
	type requestParams struct {
		FriendId string `json:"friendId" binding:"required"`
		Answer   string `json:"answer"`
	}
	var params requestParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.Set(ReqError, result.ParamsError)
		return
	}
	retBool, err := model.CheckAnswer(params.FriendId, params.Answer)
	var boolMap = make(map[string]bool, 1)
	boolMap["success"] = retBool
	c.Set(ReqError, err)
	c.Set(ReqResult, boolMap)
}

/*
	申请添加好友
*/
func AddFriendNotConfirm(c *gin.Context) {
	type addFriendParams struct {
		Id         string `json:"id" binding:"required"`
		Remark     string `json:"remark"`
		Reason     string `json:"reason"`
		SourceId   string `json:"sourceId"`
		SourceType int    `json:"sourceType"`
		Answer     string `json:"answer"`
	}
	var params addFriendParams
	if err := c.ShouldBindJSON(&params); err != nil || (params.SourceType != 1 && params.SourceType != 2 && params.SourceType != 3 && params.SourceType != 4 && params.SourceType != 5) {
		c.Set(ReqError, result.ParamsError)
		return
	}

	state, err := model.AddFriendNotConfirm(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.Id, params.Remark, params.Reason, params.SourceId, params.SourceType, params.Answer)
	var retMsg = make(map[string]int)
	retMsg["state"] = state

	c.Set(ReqError, err)
	c.Set(ReqResult, retMsg)
}

//加入黑名单
func BlockFriend(c *gin.Context) {
	var params struct {
		UserId string `json:"userId" binding:"required"`
	}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.Set(ReqError, result.ParamsError)
		return
	}
	err = model.BlockFriend(c.MustGet(UserId).(string), params.UserId)
	c.Set(ReqError, err)
}

//移出黑名单
func UnblockFriend(c *gin.Context) {
	var params struct {
		UserId string `json:"userId" binding:"required"`
	}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.Set(ReqError, result.ParamsError)
		return
	}
	err = model.UnblockFriend(c.MustGet(UserId).(string), params.UserId)
	c.Set(ReqError, err)
}

//获取黑名单列表
func FriendBlockList(c *gin.Context) {
	friends, err := model.BlockedFriends(c.MustGet(UserId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, friends)
}
