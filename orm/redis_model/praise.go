package redis_model

import (
	"errors"

	"github.com/33cn/chat33/cache"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/pkg/excRate"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

//存储赞赏榜单
func SavePraiseStatic(tp int, items map[string]*types.RankingItem, startTime, endTime int64) error {

	err := cache.Cache.SaveLeaderBoard(tp, items, startTime, endTime)
	if err != nil {
		l.Warn("redis can not save room member")
	}

	return nil
}

//获取赞赏榜单
func GetPraiseStatic(tp int, startTime, endTime int64) (map[string]*types.RankingItem, error) {
	items, err := cache.Cache.GetPraiseStatic(tp, startTime, endTime)
	if err != nil {
		return nil, err
	}

	if items == nil {
		switch tp {
		case types.Like:
			maps, err := mysql.GetPraiseStaticAsLike(startTime, endTime)
			if err != nil {
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

		if items == nil {
			l.Warn("leader board can not find", "time start", startTime, "time end", endTime)
			return nil, err
		}
		err := cache.Cache.SaveLeaderBoard(tp, items, startTime, endTime)
		if err != nil {
			l.Warn("redis can not save room member")
		}
	}
	return items, nil
}
