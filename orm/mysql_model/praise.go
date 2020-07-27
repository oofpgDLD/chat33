package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func convertPraiseRecord(v map[string]string) *types.Praise {
	return &types.Praise{
		Id:          v["id"],
		ChannelType: utility.ToInt(v["channel_type"]),
		TargetId:    v["target_id"],
		LogId:       v["log_id"],
		SenderId:    v["sender_id"],
		OptId:       v["opt_id"],
		Type:        utility.ToInt(v["type"]),
		RecordId:    v["record_id"],
		CoinId:      v["coin_id"],
		CoinName:    v["coin_name"],
		Amount:      utility.ToFloat64(v["amount"]),
		CreateTime:  utility.ToInt64(v["create_time"]),
		IsDelete:    utility.ToInt(v["is_delete"]),
	}
}

func GetPraiseListByTarget(msgSenderId string, channelType int, targetId string, startId int64, number int) ([]*types.Praise, int64, error) {
	maps, nextId, err := db.GetPraiseListByTarget(msgSenderId, channelType, targetId, startId, number)
	if err != nil {
		return nil, nextId, err
	}
	list := make([]*types.Praise, 0)
	for _, info := range maps {
		item := convertPraiseRecord(info)
		list = append(list, item)
	}
	return list, nextId, nil
}

func GetPraiseByLogIdAndOptId(channelType int, logId, optId string) ([]*types.Praise, error) {
	maps, err := db.GetPraiseByLogIdAndOptId(channelType, logId, optId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.Praise, 0)
	for _, info := range maps {
		item := convertPraiseRecord(info)
		list = append(list, item)
	}
	return list, nil
}

func GetPraiseByLogId(channelType int, logId string) ([]*types.Praise, error) {
	maps, err := db.GetPraiseByLogId(channelType, logId)
	if err != nil {
		return nil, err
	}
	list := make([]*types.Praise, 0)
	for _, info := range maps {
		item := convertPraiseRecord(info)
		list = append(list, item)
	}
	return list, nil
}

func GetPraiseByLogIdLimit(channelType int, logId string, startId int64, number int) ([]*types.Praise, int64, error) {
	maps, nextId, err := db.GetPraiseByLogIdLimit(channelType, logId, startId, number)
	if err != nil {
		return nil, nextId, err
	}
	list := make([]*types.Praise, 0)
	for _, info := range maps {
		item := convertPraiseRecord(info)
		list = append(list, item)
	}
	return list, nextId, nil
}

func LikeOrRewardCount(channelType int, logId string, tp int) (int, error) {
	maps, err := db.LikeOrRewardCount(channelType, logId, tp)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	return utility.ToInt(maps[0]["count"]), nil
}

func InsertPraise(val *types.Praise) error {
	_, _, err := db.InsertPraise(val)
	return err
}

func InsertPraiseUser(val *types.PraiseUser) error {
	_, _, err := db.InsertPraiseUser(val)
	return err
}

func DeletePraise(id string) error {
	_, _, err := db.DeletePraise(id)
	return err
}

func GetPraiseTodayLimit(optId string) ([]map[string]string, error) {
	rows, err := db.GetPraiseTodayLimit(optId)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, nil
	}
	return rows, nil
}

func GetPraiseStaticAsLike(startTime, endTime int64) ([]map[string]string, error) {
	rows, err := db.GetLeaderBoardAsLike(startTime, endTime)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, nil
	}
	return rows, nil
}

func GetPraiseStaticAsReward(startTime, endTime int64) ([]map[string]string, error) {
	rows, err := db.GetLeaderBoardAsReward(startTime, endTime)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, nil
	}
	return rows, nil
}
