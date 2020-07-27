package api

import (
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/result"
	"github.com/gin-gonic/gin"
)

// 获取奖赏记录
func PraiseList(c *gin.Context) {
	var params struct {
		ChannelType int    `json:"channelType" binding:"required"`
		TargetId    string `json:"targetId" binding:"required"`
		StartId     string `json:"startId"`
		Number      int    `json:"number"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	logs, err := model.PraiseList(c.MustGet(UserId).(string), params.ChannelType, params.TargetId, params.StartId, params.Number)
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

// 获取奖赏详情
func PraiseDetails(c *gin.Context) {
	var params struct {
		ChannelType int    `json:"channelType" binding:"required"`
		LogId       string `json:"logId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	logs, err := model.PraiseDetails(c.MustGet(UserId).(string), params.ChannelType, params.LogId)
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

// 获取奖赏详情
func PraiseDetailList(c *gin.Context) {
	var params struct {
		ChannelType int    `json:"channelType" binding:"required"`
		LogId       string `json:"logId" binding:"required"`
		StartId     string `json:"startId"`
		Number      int    `json:"number"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	logs, err := model.PraiseDetailList(c.MustGet(UserId).(string), params.ChannelType, params.LogId, params.StartId, params.Number)
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

// 点赞
func PraiseLike(c *gin.Context) {
	var params struct {
		ChannelType int    `json:"channelType" binding:"required"`
		LogId       string `json:"logId" binding:"required"`
		Action      string `json:"action"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.PraiseLike(c.MustGet(UserId).(string), params.ChannelType, params.LogId, params.Action)
	c.Set(ReqError, err)
}

// 打赏
func PraiseReward(c *gin.Context) {
	var params struct {
		ChannelType int     `json:"channelType" binding:"required"`
		LogId       string  `json:"logId" binding:"required"`
		Currency    string  `json:"currency" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Password    string  `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	err := model.PraiseReward(c.MustGet(UserId).(string), c.MustGet(AppId).(string), c.MustGet(Token).(string),
		params.Currency, params.Amount, params.Password, params.ChannelType, params.LogId)
	c.Set(ReqError, err)
}

//打赏用户
func PraiseRewardUser(c *gin.Context) {
	var params struct {
		UserId   string  `json:"userId" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
		Password string  `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	err := model.PraiseRewardUser(c.MustGet(UserId).(string), c.MustGet(AppId).(string), c.MustGet(Token).(string),
		params.Currency, params.Amount, params.Password, params.UserId)
	c.Set(ReqError, err)
}

//赞赏榜单
func LeaderBoard(c *gin.Context) {
	var params struct {
		Type      int   `json:"type" binding:"required"`
		StartId   int   `json:"startId" binding:"required"`
		Number    int   `json:"number" binding:"required"`
		StartTime int64 `json:"startTime" binding:"required"`
		EndTime   int64 `json:"endTime" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	if params.StartId != 0 && params.Number < 1 {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "number 不得小于1"))
		return
	}
	logs, err := model.LeaderBoard(c.MustGet(UserId).(string), params.Type, params.StartTime, params.EndTime, params.StartId, params.Number)
	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}

//历史榜单
func LeaderBoardHistory(c *gin.Context) {
	var params struct {
		Page int `json:"page" binding:"required"`
		Size int `json:"size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//page number的特殊值判断
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "size 不得小于1"))
		return
	}
	logs, err := model.LeaderBoardHistory(c.MustGet(UserId).(string), params.Page, params.Size)

	c.Set(ReqError, err)
	c.Set(ReqResult, logs)
}
