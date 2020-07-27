package model

import (
	"sort"
	"time"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/pkg/account"
	"github.com/33cn/chat33/pkg/excRate"
	"github.com/33cn/chat33/pkg/work"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

const praiseUpdateInterval = 5 * time.Minute

var logPraise = log15.New("model", "model/praise")

func SaveDaily() {
	go func() {

		redisDaily()

		for range time.Tick(praiseUpdateInterval) {

			redisDaily()
		}
	}()
}

func redisDaily() {
	//获取所有UserId
	member, err := orm.GetUsers()
	if err != nil {
		logPraise.Error("GetUserId error", "log info", err)
	}
	//如果是当前周，转化为这周的起始和当前时间
	nowStart, nowEnd := utility.GetWeekStartAndEnd(time.Now())

	likes := make(map[string]*types.RankingItem)
	rewards := make(map[string]*types.RankingItem)
	like, err := statisticsBySql(types.Like, nowStart, nowEnd, member)
	if err != nil {
		logPraise.Error("get statisticsBoard error", "log info", err)
	}
	for i, item := range like {
		likes[utility.ToString(i)] = item
	}
	reward, err := statisticsBySql(types.Reward, nowStart, nowEnd, member)
	if err != nil {
		logPraise.Error("get statisticsBoard error", "log info", err)
	}
	for i, item := range reward {
		rewards[utility.ToString(i)] = item
	}

	orm.SavePraiseStatic(types.Like, likes, nowStart, nowEnd)
	orm.SavePraiseStatic(types.Reward, rewards, nowStart, nowEnd)
}

func praiseList(caller string, channelType int, targetId string, startId int64, number int) ([]*types.PraiseRecord, int64, error) {
	logs, nextId, err := orm.GetPraiseListByTarget(caller, channelType, targetId, startId, number)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	if logs == nil {
		return nil, -1, nil
	}

	list := make([]*types.PraiseRecord, 0)
	for _, log := range logs {
		user, err := orm.GetUserInfoById(log.OptId)
		if err != nil {
			continue
		}
		info := &types.PraiseRecord{
			RecordId:    log.Id,
			ChannelType: log.ChannelType,
			LogId:       log.LogId,
			CreateTime:  log.CreateTime,
			User: struct {
				Id     string `json:"id"`
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
			}{Id: user.UserId, Name: user.Username, Avatar: user.Avatar},
			Type:     log.Type,
			CoinName: log.CoinName,
			Amount:   log.Amount,
		}
		list = append(list, info)
	}

	return list, nextId, nil
}

func praiseDetailList(caller string, channelType int, logId string, startId int64, number int) ([]*types.PraiseRecord, int64, error) {
	logs, nextId, err := orm.GetPraiseByLogIdLimit(channelType, logId, startId, number)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	if logs == nil {
		return nil, -1, nil
	}

	list := make([]*types.PraiseRecord, 0)
	for _, log := range logs {
		user, err := orm.GetUserInfoById(log.OptId)
		if err != nil {
			continue
		}
		info := &types.PraiseRecord{
			RecordId:    log.Id,
			ChannelType: log.ChannelType,
			LogId:       log.LogId,
			CreateTime:  log.CreateTime,
			User: struct {
				Id     string `json:"id"`
				Name   string `json:"name"`
				Avatar string `json:"avatar"`
			}{Id: user.UserId, Name: user.Username, Avatar: user.Avatar},
			Type:     log.Type,
			CoinName: log.CoinName,
			Amount:   log.Amount,
		}
		list = append(list, info)
	}

	return list, nextId, nil
}

func PraiseList(caller string, channelType int, targetId, startId string, number int) (interface{}, error) {
	if startId != "" && number < 1 {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "number 不得小于1")
	}
	//默认number数量为20
	startLogId := utility.ToInt64(startId)
	if startLogId == 0 && number <= 0 {
		number = 20
	}

	list, nextLog, err := praiseList(caller, channelType, targetId, startLogId, number)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if list == nil {
		list = make([]*types.PraiseRecord, 0)
	}

	var ret = make(map[string]interface{})
	ret["records"] = list
	ret["nextLog"] = utility.ToString(nextLog)
	return ret, nil
}

func PraiseDetails(caller string, channelType int, logId string) (interface{}, error) {
	//消息体
	var log *types.ChatLog
	switch channelType {
	case types.ToRoom:
		l, err := orm.FindRoomChatLogByContentId(logId)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if l == nil {
			logPraise.Warn("PraiseDetails", "warn", "room ChatLogNotFind", "logId", logId)
			return nil, result.NewError(result.ChatLogNotFind)
		}
		log, err = GetChatLogAsRoom(caller, l)
		if err != nil {
			return nil, err
		}
	case types.ToUser:
		l, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if l == nil {
			logPraise.Warn("PraiseDetails", "warn", "private ChatLogNotFind", "logId", logId)
			return nil, result.NewError(result.ChatLogNotFind)
		}
		log, err = GetChatLogAsUser(caller, l)
		if err != nil {
			return nil, err
		}
	default:
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized channelType")
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

	//赞赏人数
	var praiseNumber = 0
	var reward float64 = 0
	list, err := orm.GetPraiseByLogId(channelType, logId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	praiseNumber = len(list)
	for _, item := range list {
		//统计赏金
		if item.Type == types.Reward {
			coinName := item.CoinName
			//币种的个数  比如BTC  3.5
			amount := item.Amount
			//调用汇率转换接口
			price := excRate.Price(coinName, amount)
			reward += price
		}
	}

	var ret = make(map[string]interface{})
	ret["log"] = log
	ret["state"] = state
	ret["praiseNumber"] = praiseNumber
	ret["reward"] = reward
	return ret, nil
}

func PraiseDetailList(caller string, channelType int, logId, startId string, number int) (interface{}, error) {
	if startId != "" && number < 1 {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "number 不得小于1")
	}
	//默认number数量为20
	startLogId := utility.ToInt64(startId)
	if startLogId == 0 && number <= 0 {
		number = 20
	}

	list, nextLog, err := praiseDetailList(caller, channelType, logId, startLogId, number)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if list == nil {
		list = make([]*types.PraiseRecord, 0)
	}

	var ret = make(map[string]interface{})
	ret["records"] = list
	ret["nextLog"] = utility.ToString(nextLog)
	return ret, nil
}

func PraiseLike(caller string, channelType int, logId, action string) error {
	//一条消息一个人只能点一个赞，可打赏多次
	rs, err := orm.GetPraiseByLogIdAndOptId(channelType, logId, caller)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	for _, r := range rs {
		if r.Type == types.Like {
			switch action {
			case "like":
				//已经点过赞
				return nil
			case "cancel_like":
				//取消点赞功能 取消-2020年1月7日16:08:07 dld
				return nil
				err := orm.DeletePraise(r.Id)
				if err != nil {
					return result.NewError(result.DbConnectFail)
				}
			}
		}
	}

	var praise *types.Praise
	//发送通知消息
	senderId := ""
	targetId := ""
	var members []string
	switch channelType {
	case types.ToRoom:
		l, err := orm.FindRoomChatLogByContentId(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if l == nil {
			logPraise.Warn("PraiseLike", "warn", "room ChatLogNotFind", "logId", logId)
			return result.NewError(result.ChatLogNotFind)
		}
		targetId = l.RoomId
		senderId = l.SenderId
		praise = &types.Praise{
			ChannelType: channelType,
			TargetId:    l.RoomId,
			LogId:       logId,
			SenderId:    l.SenderId,
			OptId:       caller,
			Type:        types.Like,
			CreateTime:  utility.NowMillionSecond(),
		}
	case types.ToUser:
		l, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if l == nil {
			logPraise.Warn("PraiseLike", "warn", "private ChatLogNotFind", "logId", logId)
			return result.NewError(result.ChatLogNotFind)
		}
		targetId = l.ReceiveId
		senderId = l.SenderId
		members = []string{l.SenderId, l.ReceiveId}
		praise = &types.Praise{
			ChannelType: channelType,
			TargetId:    l.ReceiveId,
			LogId:       logId,
			SenderId:    l.SenderId,
			OptId:       caller,
			Type:        types.Like,
			CreateTime:  utility.NowMillionSecond(),
		}
	default:
		return result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized channelType")
	}
	if action == "like" {
		//插入记录
		err = orm.InsertPraise(praise)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	}
	likeNumber, err := orm.LikeOrRewardCount(channelType, logId, types.Like)
	if err != nil {
		return nil
	}
	rewardNumber, err := orm.LikeOrRewardCount(channelType, logId, types.Reward)
	if err != nil {
		return nil
	}
	SendAlert(senderId, targetId, channelType, members, types.Alert, proto.ComposePraiseAlert(logId, caller, action, likeNumber, rewardNumber))

	return nil
}

func PraiseReward(caller, appId, token, currency string, amount float64, payPassword string, channelType int, logId string) error {
	//调用打币接口
	app := app.GetApp(appId)
	if app == nil {
		return result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}
	//先判断这一次的打赏是否超过限额,排除一次打赏就超过限额的情况，例如1BTC
	reward := excRate.Price(currency, utility.ToFloat64(amount))
	if reward > 10000 {
		return result.NewError(result.PraiseAmountLimited)
	}
	// 查询转账今日限额  本次打赏加上已有打赏是否大于限额
	maps, err := orm.GetPraiseTodayLimit(caller)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	if maps != nil {
		for _, item := range maps {
			//统计赏金
			coinName := item["coin_name"]
			amount := item["amount"]
			//调用汇率转换接口
			price := excRate.Price(coinName, utility.ToFloat64(amount))
			reward += price
		}
	}
	//TODO 先暂定10000，以后改成可配置
	if reward > 10000 {
		logPraise.Warn("PraiseReward", "warn", "PraiseAmountLimited", "reward", reward)
		return result.NewError(result.PraiseAmountLimited)
	}

	var praise *types.Praise
	//发送通知消息
	senderId := ""
	targetId := ""
	uid := ""
	var members []string
	switch channelType {
	case types.ToRoom:
		l, err := orm.FindRoomChatLogByContentId(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if l == nil {
			logPraise.Warn("PraiseReward", "warn", "room ChatLogNotFind", "logId", logId)
			return result.NewError(result.ChatLogNotFind)
		}
		targetId = l.RoomId
		senderId = l.SenderId
		uid = l.Uid
	case types.ToUser:
		l, err := orm.FindPrivateChatLogById(logId)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if l == nil {
			logPraise.Warn("PraiseReward", "warn", "private ChatLogNotFind", "logId", logId)
			return result.NewError(result.ChatLogNotFind)
		}
		targetId = l.ReceiveId
		senderId = l.SenderId
		uid = l.Uid
		members = []string{l.SenderId, l.ReceiveId}
	default:
		return result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "unrecognized channelType")
	}

	//调用转账接口
	coinId := ""
	bytes, err := account.BcoinPraise(app.AccountServer, token, uid, currency, utility.ToString(amount), payPassword, utility.RandomID(), "gratuity")
	if err != nil {
		return result.NewError(result.PraiseRewardErr).SetExtMessage(err.Error())
	}
	msg := bytes.(map[string]interface{})
	coinName := utility.ToString(msg["currency"])
	recordId := utility.ToString(msg["unique_id"])
	amount = utility.ToFloat64(msg["amount"])

	praise = &types.Praise{
		ChannelType: channelType,
		TargetId:    targetId,
		LogId:       logId,
		SenderId:    senderId,
		OptId:       caller,
		Type:        types.Reward,
		RecordId:    recordId,
		CoinId:      coinId,
		CoinName:    coinName,
		Amount:      amount,
		CreateTime:  utility.NowMillionSecond(),
	}

	//插入记录
	err = orm.InsertPraise(praise)
	if err != nil {
		logPraise.Error("InsertPraise", "err", err)
		return result.NewError(result.DbConnectFail)
	}

	likeNumber, err := orm.LikeOrRewardCount(channelType, logId, types.Like)
	if err != nil {
		logPraise.Error("LikeOrRewardCount Like", "err", err)
		return nil
	}
	rewardNumber, err := orm.LikeOrRewardCount(channelType, logId, types.Reward)
	if err != nil {
		logPraise.Error("LikeOrRewardCount Reward", "err", err)
		return nil
	}
	SendAlert(senderId, targetId, channelType, members, types.Alert, proto.ComposePraiseAlert(logId, caller, "reward", likeNumber, rewardNumber))
	return nil
}

func PraiseRewardUser(caller, appId, token, currency string, amount float64, payPassword string, targetId string) error {
	//调用打币接口
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	//TODO 这些限额的判断还不确定
	//先判断这一次的打赏是否超过限额,排除一次打赏就超过限额的情况，例如1BTC
	reward := excRate.Price(currency, utility.ToFloat64(amount))
	if reward > 10000 {
		logPraise.Warn("PraiseReward", "warn", "PraiseAmountLimited", "reward", reward)
		return result.NewError(result.PraiseAmountLimited)
	}
	// 查询转账今日限额  本次打赏加上已有打赏是否大于限额
	maps, err := orm.GetPraiseTodayLimit(caller)
	if err != nil {
		logPraise.Error("praise Get Praise Today's Limit err", "err", err)
		return result.NewError(result.DbConnectFail)
	}

	if maps != nil {
		for _, item := range maps {
			//统计赏金
			coinName := item["coin_name"]
			amount := item["amount"]
			//调用汇率转换接口
			price := excRate.Price(coinName, utility.ToFloat64(amount))
			reward += price
		}
	}
	//TODO 先暂定10000，以后改成可配置
	if reward > 10000 {
		logPraise.Warn("PraiseReward", "warn", "PraiseAmountLimited", "reward", reward)
		return result.NewError(result.PraiseAmountLimited)
	}

	targetUser, err := orm.GetUserInfoById(targetId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if targetUser == nil {
		return result.NewError(result.UserNotExists)
	}
	uid := targetUser.Uid
	//调用转账接口
	coinId := ""
	bytes, err := account.BcoinPraise(app.AccountServer, token, uid, currency, utility.ToString(amount), payPassword, utility.RandomID(), "gratuity")
	if err != nil {
		return result.NewError(result.PraiseRewardErr).SetExtMessage(err.Error())
	}
	msg := bytes.(map[string]interface{})
	coinName := utility.ToString(msg["currency"])
	recordId := utility.ToString(msg["unique_id"])
	amount = utility.ToFloat64(msg["amount"])

	praise := &types.PraiseUser{
		TargetId:   targetId,
		OptId:      caller,
		RecordId:   recordId,
		CoinId:     coinId,
		CoinName:   coinName,
		Amount:     amount,
		CreateTime: utility.NowMillionSecond(),
	}

	//插入记录
	err = orm.InsertPraiseUser(praise)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//获取一家公司员工列表
func enterpriseMembers(userId string) (string, []string, error) {
	//请求者信息
	u, err := orm.GetUserInfoById(userId)
	if err != nil {
		return "", nil, err
	}
	if u == nil {
		return "", nil, nil
	}

	memInfo, err := work.GetUsers(u.AppId, u.Uid)
	if err != nil {
		return "", nil, err
	}
	if memInfo == nil {
		return "", []string{}, nil
	}
	users := make([]string, 0)
	for _, user := range memInfo.Users {
		u, err := orm.GetUserInfoByUid(user.AppId, user.Uid)
		if err != nil {
			return "", nil, err
		}
		if u != nil {
			users = append(users, u.UserId)
		}
	}
	return memInfo.Enterprise.Name, users, nil
}

//排行榜
func statisticsBoard(tp int, startTime, endTime int64, members []string) ([]*types.RankingItem, error) {
	items, err := orm.GetPraiseStatic(tp, startTime, endTime)
	if err != nil {
		return nil, err
	}
	ret := make([]*types.RankingItem, 0)
	for _, member := range members {
		if m, ok := items[member]; ok {
			ret = append(ret, m)
		} else {
			ret = append(ret, &types.RankingItem{
				UserId: member,
				Type:   tp,
			})
		}
	}

	//排序
	switch tp {
	case types.Like:
		sort.Sort(types.RankingItemWrapper{Items: ret, By: func(p, q *types.RankingItem) bool {
			return q.Number < p.Number // 递减排序
		}})
	case types.Reward:
		sort.Sort(types.RankingItemWrapper{Items: ret, By: func(p, q *types.RankingItem) bool {
			return q.Price < p.Price // 递减排序
		}})
	}
	return ret, nil
}

//只经过数据库查询,定时任务用
func statisticsBySql(tp int, startTime, endTime int64, members []string) ([]*types.RankingItem, error) {
	items, err := orm.GetPraiseStaticBySql(tp, startTime, endTime)
	if err != nil {
		return nil, err
	}
	ret := make([]*types.RankingItem, 0)
	for _, member := range members {
		if m, ok := items[member]; ok {
			ret = append(ret, m)
		} else {
			ret = append(ret, &types.RankingItem{
				UserId: member,
				Type:   tp,
			})
		}
	}
	return ret, nil
}

func LeaderBoard(userId string, tp int, startTime, endTime int64, startId, number int) (interface{}, error) {
	//如果是当前周，转化为这周的起始和当前时间
	nowStart, nowEnd := utility.GetWeekStartAndEnd(time.Now())
	if nowStart <= startTime && nowEnd >= endTime {
		startTime = nowStart
		endTime = nowEnd
	}
	/*//将时间转为一天的开始
	startTime, _ = utility.DayStartAndEndNowMillionSecond(utility.MillionSecondToDateTime(startTime))
	//将时间转为一天的结束
	_, endTime = utility.DayStartAndEndNowMillionSecond(utility.MillionSecondToDateTime(endTime))*/

	startIndex := startId
	endIndex := startId + number

	//获取一家公司员工列表
	epName, members, err := enterpriseMembers(userId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	items, err := statisticsBoard(tp, startTime, endTime, members)
	if err != nil {
		return nil, err
	}
	var mine types.Record
	//榜单的的具体信息
	records := make([]*types.Record, 0)
	for i, item := range items {
		i = i + 1
		//自己的信息
		if item.UserId == userId {
			userName := ""
			avatar := ""
			u, err := orm.GetUserInfoById(item.UserId)
			if err != nil {
			}
			if u != nil {
				userName = u.Username
				avatar = u.Avatar
			}
			mine = types.Record{
				Ranking: i,
				User: struct {
					Id     string `json:"id"`
					Name   string `json:"name"`
					Avatar string `json:"avatar"`
				}{Id: item.UserId, Name: userName, Avatar: avatar},
				Number: item.Number,
				Price:  item.Price,
			}
		}
		//找出在查找区间内的详情
		if i >= startIndex && i < endIndex {
			userName := ""
			avatar := ""
			u, err := orm.GetUserInfoById(item.UserId)
			if err != nil {
			}
			if u != nil {
				userName = u.Username
				avatar = u.Avatar
			}
			record := &types.Record{
				Ranking: i,
				User: struct {
					Id     string `json:"id"`
					Name   string `json:"name"`
					Avatar string `json:"avatar"`
				}{Id: item.UserId, Name: userName, Avatar: avatar},
				Number: item.Number,
				Price:  item.Price,
			}
			records = append(records, record)
		}
	}

	//获取企业名称
	var ret = make(map[string]interface{})
	ret["startTime"] = startTime
	ret["endTime"] = endTime
	ret["enterprise"] = map[string]string{
		"name": epName,
	}
	ret["mine"] = mine
	ret["records"] = records
	if len(items) >= startId+number {
		ret["nextLog"] = startId + number
	} else {
		ret["nextLog"] = -1
	}
	return ret, nil
}

func LeaderBoardHistory(userId string, page, size int) (interface{}, error) {
	baseTime := time.Now()
	//最多拉倒功能开放前 2019-12-2 00:00:00
	openTime := int64(1575216000000)
	sec := openTime / 1000
	nsec := openTime % 1000
	openLimit := time.Unix(sec, nsec)
	//限制最多拉取到半年之前
	limitTime := baseTime.AddDate(0, -6, 0)
	var weeks = make([]*struct {
		Start int64
		End   int64
	}, 0)

	//获取page * size周（排除当周）
	for i := (page-1)*size + 1; i <= page*size; i++ {
		tm := baseTime.AddDate(0, 0, -(i * 7))
		if tm.Before(limitTime) || tm.Before(openLimit) {
			break
		}
		start, end := utility.GetWeekStartAndEnd(tm)
		weeks = append(weeks, &struct {
			Start int64
			End   int64
		}{Start: start, End: end})
	}

	//获取一家公司员工列表
	_, members, err := enterpriseMembers(userId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	records := make([]*types.PraiseBoardHistory, 0)
	//按周划分
	for _, week := range weeks {
		likeRanking := 0
		rewardRanking := 0
		number := 0
		price := float64(0)
		likeList, err := statisticsBoard(types.Like, week.Start, week.End, members)
		if err != nil {
			return nil, err
		}
		rewardList, err := statisticsBoard(types.Reward, week.Start, week.End, members)
		if err != nil {
			return nil, err
		}
		for i, items := range likeList {
			if items.UserId == userId {
				likeRanking = i + 1
				number = items.Number
				break
			}
		}

		for i, items := range rewardList {
			if items.UserId == userId {
				rewardRanking = i + 1
				price = items.Price
				break
			}
		}

		records = append(records, &types.PraiseBoardHistory{
			Like: struct {
				Ranking int `json:"ranking"`
				Number  int `json:"number"`
			}{Ranking: likeRanking, Number: number},
			Reward: struct {
				Ranking int     `json:"ranking"`
				Price   float64 `json:"price"`
			}{Ranking: rewardRanking, Price: price},
			StartTime: week.Start,
			EndTime:   week.End,
		})
	}

	var ret = make(map[string]interface{})
	ret["records"] = records
	return ret, nil
}
