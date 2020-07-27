package api

import (
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/gin-gonic/gin"
	l "github.com/inconshreveable/log15"
)

var redPacketLog = l.New("module", "chat/api/red_envelope")

// 红包统计信息
func RPStatistics(c *gin.Context) {
	type requestParams struct {
		Operation int   `json:"operation"`
		CoinId    int   `json:"coinId"`
		Type      int   `json:"type"`
		StartTime int64 `json:"startTime"`
		EndTime   int64 `json:"endTime"`
		PageNum   *int  `json:"pageNum" binding:"required"`
		PageSize  int   `json:"pageSize" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	if c.MustGet(AppId).(string) == "" {
		c.Set(ReqError, result.NewError(result.LackHeader).SetExtMessage("FZM-APP-ID"))
		return
	}

	resp, err := model.RPStatistic(c.MustGet(AppId).(string), c.MustGet(UserId).(string), c.MustGet(Token).(string), params.Operation, params.CoinId, params.Type, *params.PageNum, params.PageSize, params.StartTime/1000, params.EndTime/1000)
	if err != nil {
		c.Set(ReqError, err)
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, resp)
}

// 红包查看详情
func RedEnvelopeDetail(c *gin.Context) {
	type requestParams struct {
		PacketId string `json:"packetId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	resp, err := model.RedEnvelopeDetail(c.MustGet(AppId).(string), c.MustGet(UserId).(string), c.MustGet(Token).(string), params.PacketId)
	if err != nil {
		redPacketLog.Error("red packet detail", "err", err)
		c.Set(ReqError, err)
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, resp)
}

func RedEnvelopeReceiveDetail(c *gin.Context) {
	type requestParams struct {
		PacketId string `json:"packetId" binding:"required"`
	}

	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	resp, err := model.GetPacketRecvDetails(c.MustGet(AppId).(string), c.MustGet(UserId).(string), params.PacketId)
	if err != nil {
		redPacketLog.Error("red packet detail", "err", err)
		c.Set(ReqError, err)
		return
	}

	c.Set(ReqError, nil)
	c.Set(ReqResult, &struct {
		Rows interface{} `json:"rows"`
	}{
		Rows: resp,
	})
}

// 查询账户余额
func Balance(c *gin.Context) {
	resp, err := model.Balance(c.MustGet(AppId).(string), c.MustGet(Token).(string))
	if err != nil {
		redPacketLog.Error("query balance", "err", err.Error())
		c.Set(ReqError, err)
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, resp)
}

//（登录用户）领取红包
func ReceiveEntry(c *gin.Context) {
	type receiveEntryParams struct {
		PacketID string `json:"packetId" binding:"required"`
	}
	var params receiveEntryParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	resp, err := model.ReceiveEntry(c.MustGet(AppId).(string), c.MustGet(Token).(string), params.PacketID, c.MustGet(UserId).(string))
	if err != nil {
		redPacketLog.Info("receive entry", "err", err.Error(), "userId", c.MustGet(UserId).(string), "packetId", params.PacketID)
		c.Set(ReqError, err)
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, resp)
}

// 发红包
func Send(c *gin.Context) {
	type sendParams struct {
		Ext     interface{} `json:"ext"`
		CType   *int        `json:"cType" binding:"required"`
		ToId    string      `json:"toId" binding:"required"`
		Coin    int         `json:"coin" binding:"required"`
		Type    int         `json:"type" binding:"required"`
		Amount  float64     `json:"amount" binding:"required"`
		Size    int         `json:"size" binding:"required"`
		ToUsers string      `json:"toUsers"`
		Remark  string      `json:"remark"`
	}

	var params sendParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	if !model.CheckUserExist(c.MustGet(UserId).(string)) {
		c.Set(ReqError, result.NewError(result.UserNotExists))
		return
	}

	req := &types.SendParams{
		Amount: params.Amount,
		Size:   params.Size,
		Remark: params.Remark,
		Type:   params.Type,
		CoinId: params.Coin,
		To:     params.ToUsers,
		Ext:    params.Ext,
	}

	packet, err := model.Send(c.MustGet(AppId).(string), c.MustGet(Token).(string), req, c.MustGet(UserId).(string), params.ToId, *params.CType)
	if err != nil {
		redPacketLog.Error("send packet", "err", err.Error())
		c.Set(ReqError, err)
		return
	}

	err = model.InsertPacketInfo(packet.PacketId, c.MustGet(UserId).(string), params.ToId, params.Type, params.Size, params.Amount,
		packet.Remark, *params.CType, params.Coin, utility.NowMillionSecond())
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	//里面做了截断，所以返回的时候再原样赋值回去
	packet.Remark = params.Remark
	c.Set(ReqError, nil)
	c.Set(ReqResult, packet)
}

//获取代币信息
func GetCoinInfo(c *gin.Context) {
	coinInfo, err := model.GetCoinInfo(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, coinInfo)
}
