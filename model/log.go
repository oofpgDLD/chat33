package model

import (
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

//从消息接收者的视角，返回这条消息的点赞详情
func getPraise(caller string, channelType int, logId string) (interface{}, error) {
	likeNumber, err := orm.LikeOrRewardCount(channelType, logId, types.Like)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	rewardNumber, err := orm.LikeOrRewardCount(channelType, logId, types.Reward)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	state := 0
	//查询自己是否点赞
	rs, err := orm.GetPraiseByLogIdAndOptId(channelType, logId, caller)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	for _, r := range rs {
		if r.Type == types.Like {
			state |= types.PS_Like
		}
		if r.Type == types.Reward {
			state |= types.PS_Reward
		}
	}

	var ret = make(map[string]interface{})
	ret["like"] = likeNumber
	ret["reward"] = rewardNumber
	ret["state"] = state
	return ret, nil
}

func GetChatLogAsRoom(userId string, ru *types.RoomLogJoinUser) (*types.ChatLog, error) {
	info := &types.ChatLog{
		LogId:       ru.Id,
		MsgId:       ru.MsgId,
		ChannelType: types.ToRoom,
		IsSnap:      ru.IsSnap,
		FromId:      ru.SenderId,
		TargetId:    ru.RoomId,
		MsgType:     ru.MsgType,
		Datetime:    ru.Datetime,
	}
	info.SetIsDel(ru.IsDelete)
	var senderInfo = make(map[string]interface{})
	senderInfo["nickname"] = ru.Username
	senderInfo["avatar"] = ru.Avatar
	info.SenderInfo = senderInfo

	con := utility.StringToJobj(ru.Content)
	if info.MsgType == types.RedPack {
		err := ConvertRedPackInfoToSend(userId, con)
		if err != nil {
			logRoom.Error("GetChatLog ConvertRedPackInfoToSend", "err", err.Error())
			return nil, result.NewError(result.MsgFormatError).SetExtMessage(err.Error())
		}
	}

	//转账和收款也添加上打赏消息（需求） 2019年12月16日11:05:43 —— dld
	if info.MsgType != types.Alert /*&& info.MsgType != types.Transfer && info.MsgType != types.Receipt*/ {
		praise, err := getPraise(userId, info.ChannelType, info.LogId)
		if err != nil {
			return nil, result.NewError(result.MsgFormatError).SetExtMessage(err.Error())
		}
		info.Praise = praise
	}

	ext := utility.StringToJobj(ru.Ext)
	info.SetCon(con)
	info.Msg = con
	info.Ext = ext
	return info, nil
}

func GetChatLogAsUser(userId string, pu *types.PrivateLogJoinUser) (*types.ChatLog, error) {
	info := &types.ChatLog{
		LogId:       pu.Id,
		MsgId:       pu.MsgId,
		ChannelType: types.ToUser,
		IsSnap:      pu.IsSnap,
		FromId:      pu.SenderId,
		TargetId:    pu.ReceiveId,
		MsgType:     pu.MsgType,
		Datetime:    pu.SendTime,
	}
	var senderInfo = make(map[string]interface{})
	senderInfo["nickname"] = pu.Username
	senderInfo["avatar"] = pu.Avatar
	if pu.Remark != "" {
		senderInfo["remark"] = pu.Remark
	}
	info.SenderInfo = senderInfo

	con := utility.StringToJobj(pu.Content)
	if info.MsgType == types.RedPack {
		err := ConvertRedPackInfoToSend(userId, con)
		if err != nil {
			logRoom.Error("GetChatLog ConvertRedPackInfoToSend", "err", err)
			return nil, result.NewError(result.MsgFormatError).SetExtMessage(err.Error())
		}
	}

	//转账和收款也添加上打赏消息（需求） 2019年12月16日11:05:43 —— dld
	if info.MsgType != types.Alert /*&& info.MsgType != types.Transfer && info.MsgType != types.Receipt*/ {
		praise, err := getPraise(userId, info.ChannelType, info.LogId)
		if err != nil {
			return nil, result.NewError(result.MsgFormatError).SetExtMessage(err.Error())
		}
		info.Praise = praise
	}

	ext := utility.StringToJobj(pu.Ext)
	info.SetCon(con)
	info.Msg = con
	info.Ext = ext
	return info, nil
}
