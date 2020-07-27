package orm

import (
	"errors"

	mysql "github.com/33cn/chat33/orm/mysql_model"
	redis "github.com/33cn/chat33/orm/redis_model"
	"github.com/33cn/chat33/pkg/excRate"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logPraise = log15.New("module", "model/praise")

func GetPraiseListByTarget(msgSenderId string, channelType int, targetId string, startId int64, number int) ([]*types.Praise, int64, error) {
	ret, ret2, err := mysql.GetPraiseListByTarget(msgSenderId, channelType, targetId, startId, number)
	if err != nil {
		logPraise.Error("mysql.GetPraiseListByTarget", "err", err, "msgSenderId", msgSenderId, "channelType", channelType, "targetId", targetId, "startId", startId, "number", number)
	}
	return ret, ret2, err
}

func GetPraiseByLogIdAndOptId(channelType int, logId, optId string) ([]*types.Praise, error) {
	ret, err := mysql.GetPraiseByLogIdAndOptId(channelType, logId, optId)
	if err != nil {
		logPraise.Error("mysql.GetPraiseByLogIdAndOptId", "err", err, "channelType", channelType, "logId", logId, "optId", optId)
	}
	return ret, err
}

func GetPraiseByLogId(channelType int, logId string) ([]*types.Praise, error) {
	ret, err := mysql.GetPraiseByLogId(channelType, logId)
	if err != nil {
		logPraise.Error("mysql.GetPraiseByLogId", "err", err, "channelType", channelType, "logId", logId)
	}
	return ret, err
}

func GetPraiseByLogIdLimit(channelType int, logId string, startId int64, number int) ([]*types.Praise, int64, error) {
	ret, ret2, err := mysql.GetPraiseByLogIdLimit(channelType, logId, startId, number)
	if err != nil {
		logPraise.Error("mysql.GetPraiseByLogIdLimit", "err", err, "channelType", channelType, "logId", logId, "startId", startId, "number", number)
	}
	return ret, ret2, err
}

func LikeOrRewardCount(channelType int, logId string, tp int) (int, error) {
	ret, err := mysql.LikeOrRewardCount(channelType, logId, tp)
	if err != nil {
		logPraise.Error("mysql.LikeOrRewardCount", "err", err, "channelType", channelType, "logId", logId, "tp", tp)
	}
	return ret, err
}

func InsertPraise(val *types.Praise) error {
	err := mysql.InsertPraise(val)
	if err != nil {
		logPraise.Error("mysql.InsertPraise", "err", err)
	}
	return err
}

func InsertPraiseUser(val *types.PraiseUser) error {
	err := mysql.InsertPraiseUser(val)
	if err != nil {
		logPraise.Error("mysql.InsertPraiseUser", "err", err)
	}
	return err
}

func DeletePraise(id string) error {
	err := mysql.DeletePraise(id)
	if err != nil {
		logPraise.Error("mysql.DeletePraise", "err", err, "id", id)
	}
	return err
}

//获取今日币种和数量
func GetPraiseTodayLimit(optId string) ([]map[string]string, error) {
	ret, err := mysql.GetPraiseTodayLimit(optId)
	if err != nil {
		logPraise.Error("mysql.GetPraiseTodayLimit", "err", err, "optId", optId)
	}
	return ret, err
}

//获取 赞赏统计
func GetPraiseStatic(tp int, startTime, endTime int64) (map[string]*types.RankingItem, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetPraiseStatic(tp, startTime, endTime)
		if err != nil {
			logPraise.Error("redis.GetPraiseStatic", "err", err, "tp", tp, "startTime", startTime, "endTime", endTime)
		}
		return ret, err
	}

	var items map[string]*types.RankingItem
	switch tp {
	case types.Like:
		maps, err := mysql.GetPraiseStaticAsLike(startTime, endTime)
		if err != nil {
			logPraise.Error("mysql.GetPraiseStaticAsLike", "err", err, "startTime", startTime, "endTime", endTime)
			return nil, err
		}
		items = make(map[string]*types.RankingItem)
		for _, m := range maps {
			number := utility.ToInt(m["number"])
			items[m["sender_id"]] = &types.RankingItem{
				UserId: m["sender_id"],
				Type:   types.Like,
				Number: number,
			}
		}
	case types.Reward:
		maps, err := mysql.GetPraiseStaticAsReward(startTime, endTime)
		if err != nil {
			logPraise.Error("mysql.GetPraiseStaticAsReward", "err", err, "startTime", startTime, "endTime", endTime)
			return nil, err
		}
		items = make(map[string]*types.RankingItem)
		for _, m := range maps {
			if item, ok := items[m["sender_id"]]; ok {
				coinName := m["coin_name"]
				amount := m["amount"]
				price := excRate.Price(coinName, utility.ToFloat64(amount)) + item.Price
				item.Price = price
			} else {
				coinName := m["coin_name"]
				amount := m["amount"]
				price := excRate.Price(coinName, utility.ToFloat64(amount))
				items[m["sender_id"]] = &types.RankingItem{
					UserId: m["sender_id"],
					Type:   types.Reward,
					Price:  price,
				}
			}
		}
	default:
		return nil, errors.New("不是点赞或打赏类型")
	}
	return items, nil
}

//不通过redis
func GetPraiseStaticBySql(tp int, startTime, endTime int64) (map[string]*types.RankingItem, error) {
	var items map[string]*types.RankingItem
	switch tp {
	case types.Like:
		maps, err := mysql.GetPraiseStaticAsLike(startTime, endTime)
		if err != nil {
			logPraise.Error("mysql.GetPraiseStaticAsLike", "err", err, "tp", tp, "startTime", startTime, "endTime", endTime)
			return nil, err
		}
		items = make(map[string]*types.RankingItem)
		for _, m := range maps {
			number := utility.ToInt(m["number"])
			items[m["sender_id"]] = &types.RankingItem{
				UserId: m["sender_id"],
				Type:   types.Like,
				Number: number,
			}
		}
	case types.Reward:
		maps, err := mysql.GetPraiseStaticAsReward(startTime, endTime)
		if err != nil {
			logPraise.Error("mysql.GetPraiseStaticAsReward", "err", err, "startTime", startTime, "endTime", endTime)
			return nil, err
		}
		items = make(map[string]*types.RankingItem)
		for _, m := range maps {
			if item, ok := items[m["sender_id"]]; ok {
				coinName := m["coin_name"]
				amount := m["amount"]
				price := excRate.Price(coinName, utility.ToFloat64(amount)) + item.Price
				item.Price = price
			} else {
				coinName := m["coin_name"]
				amount := m["amount"]
				price := excRate.Price(coinName, utility.ToFloat64(amount))
				items[m["sender_id"]] = &types.RankingItem{
					UserId: m["sender_id"],
					Type:   types.Reward,
					Price:  price,
				}
			}
		}
	default:
		return nil, errors.New("不是点赞或打赏类型")
	}
	return items, nil
}

func SavePraiseStatic(tp int, items map[string]*types.RankingItem, startTime, endTime int64) error {
	err := redis.SavePraiseStatic(tp, items, startTime, endTime)
	if err != nil {
		logPraise.Error("mysql.SavePraiseStatic", "err", err, "tp", tp, "startTime", startTime, "endTime", endTime)
	}
	return err
}
