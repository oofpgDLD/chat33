package orm

import (
	"testing"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func Test_InsertPraisePriavte(t *testing.T) {
	caller := "1"
	logId := "1695"
	receiveId := "11"
	senderId := "12"
	praise := &types.Praise{
		ChannelType: types.ToUser,
		TargetId:    receiveId,
		LogId:       logId,
		SenderId:    senderId,
		OptId:       caller,
		Type:        types.Like,
		CreateTime:  utility.NowMillionSecond(),
	}
	InsertPraise(praise)
}

func Test_InsertPraiseRoom(t *testing.T) {
	caller := "1"
	logId := "1563"
	receiveId := "60"
	senderId := "9"
	praise := &types.Praise{
		ChannelType: types.ToRoom,
		TargetId:    receiveId,
		LogId:       logId,
		SenderId:    senderId,
		OptId:       caller,
		Type:        types.Like,
		CreateTime:  utility.NowMillionSecond(),
	}
	InsertPraise(praise)
}

func Test_GetPraiseListByTarget(t *testing.T) {
	caller := "9"
	channelType := types.ToRoom
	targetId := "60"
	startId := 0
	number := 20
	logs, nextId, err := GetPraiseListByTarget(caller, channelType, targetId, int64(startId), number)
	if err != nil {
		t.Error(err)
		return
	}
	for _, log := range logs {
		t.Log(log)
	}
	t.Log(nextId)
}

func Test_GetPraiseByLogIdAndOptId(t *testing.T) {
	channelType := types.ToRoom
	logId := "1563"
	optId := "1"
	logs, err := GetPraiseByLogIdAndOptId(channelType, logId, optId)
	if err != nil {
		t.Error(err)
		return
	}
	for _, log := range logs {
		t.Log(log)
	}
}

func Test_GetPraiseByLogId(t *testing.T) {
	channelType := types.ToRoom
	logId := "1563"
	logs, err := GetPraiseByLogId(channelType, logId)
	if err != nil {
		t.Error(err)
		return
	}
	for _, log := range logs {
		t.Log(log)
	}
}

func Test_GetPraiseByLogIdLimit(t *testing.T) {
	channelType := types.ToRoom
	logId := "1563"
	startId := 0
	number := 20
	logs, nextId, err := GetPraiseByLogIdLimit(channelType, logId, int64(startId), number)
	if err != nil {
		t.Error(err)
		return
	}
	for _, log := range logs {
		t.Log(log)
	}
	t.Log(nextId)
}

func Test_LikeOrRewardCount(t *testing.T) {
	channelType := types.ToRoom
	logId := "1563"
	ret, err := LikeOrRewardCount(channelType, logId, types.Like)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}

func Test_DeletePraise(t *testing.T) {
	err := DeletePraise("1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}
