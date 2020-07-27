package redis_model

import (
	"regexp"
	"sort"

	"github.com/33cn/chat33/cache"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

var descById = func(p1, p2 *types.PrivateLogJoinUser) bool {
	return p1.Id > p2.Id
}

/*var ascById = func(p1, p2 *types.PrivateLogJoinUser) bool {
	return p1.Id < p2.Id
}*/

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *types.PrivateLogJoinUser) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(logs []*types.PrivateLogJoinUser) {
	ps := &logsSorter{
		logs: logs,
		by:   by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type logsSorter struct {
	logs []*types.PrivateLogJoinUser
	by   func(p1, p2 *types.PrivateLogJoinUser) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *logsSorter) Len() int {
	return len(s.logs)
}

// Swap is part of sort.Interface.
func (s *logsSorter) Swap(i, j int) {
	s.logs[i], s.logs[j] = s.logs[j], s.logs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *logsSorter) Less(i, j int) bool {
	return s.by(s.logs[i], s.logs[j])
}

func friendJoinUser(f *types.Friend) (*types.FriendJoinUser, error) {
	if f == nil {
		return nil, nil
	}
	u, err := GetUserInfoById(f.FriendId)
	if err != nil {
		return nil, err
	}
	return &types.FriendJoinUser{
		Friend: f,
		User:   u,
	}, nil
}

//获取好友信息
func FindFriendById(userID, friendID string) (*types.FriendJoinUser, error) {
	fInfo, err := cache.Cache.GetFriendInfo(userID, friendID)
	if err != nil {
		l.Warn("find friend by id from redis", "err", err)
	}

	info, err := friendJoinUser(fInfo)
	if err != nil {
		return nil, err
	}

	if info == nil {
		//从mysql 获取并更新缓存
		info, err = mysql.FindFriendById(userID, friendID)
		if info == nil {
			return nil, err
		}
		err = cache.Cache.SaveFriendInfo(userID, info.Friend)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return info, nil
	}
	return info, nil
}

//同意好友请求
func AcceptFriend(userId, friendId, source string, addTime int64) error {
	err := mysql.AcceptFriend(userId, friendId, source, addTime)
	if err != nil {
		return err
	}
	//修改请求记录的状态
	apply, err := FindApplyLogByUserAndTarget(friendId, userId, types.IsFriend)
	if err == nil {
		apply.State = types.AcceptState
		apply.Datetime = utility.NowMillionSecond()
		err = cache.Cache.SaveApplyLog(apply)
		if err != nil {
			l.Warn("SaveApplyLog failed", "err", err)
		}
	} else {
		l.Warn("GetApplyLogByUserAndTarget failed", "err", err)
	}
	//如果缓存中有friend信息 则更新状态
	f1, err := cache.Cache.GetFriendInfo(userId, friendId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	if f1 != nil {
		err = cache.Cache.UpdateFriendInfo(userId, friendId, "is_delete", utility.ToString(types.FriendIsNotDelete))
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
		err = cache.Cache.UpdateFriendInfo(userId, friendId, "add_time", utility.ToString(addTime))
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}

	f2, err := cache.Cache.GetFriendInfo(friendId, userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	if f2 != nil {
		err = cache.Cache.UpdateFriendInfo(friendId, userId, "is_delete", utility.ToString(types.FriendIsNotDelete))
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
		err = cache.Cache.UpdateFriendInfo(friendId, userId, "add_time", utility.ToString(addTime))
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}

	//好友表添加记录
	err = cache.Cache.AddFriend(userId, friendId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.AddFriend(friendId, userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//删除好友
func DeleteFriend(userId, friendId string, alterTime int64) error {
	err := mysql.DeleteFriend(userId, friendId, alterTime)
	if err != nil {
		return err
	}

	err = cache.Cache.DeleteFriend(userId, friendId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	err = cache.Cache.DeleteFriend(friendId, userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "is_delete", utility.ToString(types.FriendIsDelete))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "add_time", utility.ToString(alterTime))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.UpdateFriendInfo(friendId, userId, "is_delete", utility.ToString(types.FriendIsDelete))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.UpdateFriendInfo(friendId, userId, "add_time", utility.ToString(alterTime))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

//获取好友列表
func FindFriends(userId string) ([]string, error) {
	friendsId, err := cache.Cache.GetFriends(userId)
	if err != nil {
		l.Warn("find friends by id info from redis", "err", err)
	}

	if friendsId == nil {
		//从mysql 获取并更新缓存
		friendsId, err = mysql.FindFriendsId(userId)
		if friendsId == nil {
			l.Warn("friends info can not find")
			return nil, err
		}
		err = cache.Cache.SaveFriends(userId, friendsId)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return friendsId, nil
	}
	return friendsId, nil
}

//根据是否是常用群筛选
func FindFriendsFilterByCommonUse(userId string, commonUse int) ([]*types.FriendJoinUser, error) {
	friendsId, err := cache.Cache.GetFriends(userId)
	if err != nil {
		l.Warn("find friends by common use from redis", "err", err)
	}

	friends := make([]*types.FriendJoinUser, 0)
	for _, id := range friendsId {
		f, err := FindFriendById(userId, id)
		if err != nil {
			return nil, err
		}

		switch commonUse {
		case types.CommonUse, types.UncommonUse:
			if f.Type != commonUse {
				continue
			}
		default:
		}

		friends = append(friends, f)
	}

	if friendsId == nil {
		//从mysql 获取
		friends, err = mysql.FindFriendsById(userId, commonUse, types.FriendIsNotDelete)
		if friends == nil {
			l.Warn("friends info can not find")
			return nil, err
		}
		if commonUse == types.CommonUseFindAll {
			friendsId := make([]string, 0)
			for _, f := range friends {
				friendsId = append(friendsId, f.FriendId)
			}
			err = cache.Cache.SaveFriends(userId, friendsId)
			if err != nil {
				l.Warn("redis can not save user")
			}
		}
		return friends, nil
	}
	return friends, nil
}

func FindDeletedFriends(userId string) ([]string, error) {
	friends := make([]string, 0)
	//从mysql 获取
	friendsInfo, err := mysql.FindFriendsById(userId, types.CommonUseFindAll, types.FriendIsDelete)
	if friendsInfo == nil {
		l.Warn("friends info can not find")
		return nil, err
	}
	for _, f := range friendsInfo {
		friends = append(friends, f.FriendId)
	}

	/*friends, err := cache.Cache.GetDelFriends(userId)
	if err != nil {
		l.Warn("find deleted friends after time from redis", "err", err)
	}

	if friends == nil {
		//从mysql 获取
		friendsInfo, err := mysql.FindFriendsById(userId, types.CommonUseFindAll, types.FriendIsDelete)
		if friendsInfo == nil {
			l.Warn("friends info can not find")
			return nil, err
		}
		for _, f := range friendsInfo {
			friends = append(friends, f.FriendId)
		}
		err = cache.Cache.SaveDelFriends(userId, friends)
		if err != nil {
			l.Warn("SaveDelFriends failed", "err", err)
		}
		return friends, nil
	}*/
	return friends, nil
}

//根据是否是常用群筛选
func FindFriendsAfterTime(userId string, commonUse int, time int64) ([]*types.FriendJoinUser, error) {
	friendsId, err := cache.Cache.GetFriends(userId)
	if err != nil {
		l.Warn("find friends after time from redis", "err", err)
	}

	delFriendsId, err := FindDeletedFriends(userId)
	if err != nil {
		return nil, err
	}

	friends := make([]*types.FriendJoinUser, 0)
	fs := append(friendsId, delFriendsId...)
	for _, id := range fs {
		f, err := FindFriendById(userId, id)
		if err != nil {
			return nil, err
		}

		switch commonUse {
		case types.CommonUse, types.UncommonUse:
			if f.Type != commonUse {
				continue
			}
		default:
		}

		if f.AddTime < time {
			continue
		}
		friends = append(friends, f)
	}

	if friendsId == nil {
		//从mysql 获取
		friends, err = mysql.FindFriendsAfterTime(userId, commonUse, time)
		if friends == nil {
			l.Warn("friends info can not find")
			return nil, err
		}
		return friends, nil
	}
	return friends, nil
}

//是否是好友
func CheckIsFriend(userId, friendId string) (bool, error) {
	isFriend, err := cache.Cache.IsFriend(userId, friendId)
	if err != nil {
		l.Warn("check is friend from redis", "err", err)
	}

	if isFriend == nil {
		//从mysql 获取
		return mysql.CheckIsFriend(userId, friendId, types.FriendIsNotDelete)
	}
	return *isFriend, nil
}

func SetFriendRemark(userId, friendId, remark string) error {
	err := mysql.SetFriendRemark(userId, friendId, remark)
	if err != nil {
		return err
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "remark", remark)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func SetFriendExtRemark(userId, friendId, remark string) error {
	err := mysql.SetFriendExtRemark(userId, friendId, remark)
	if err != nil {
		return err
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "ext_remark", remark)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func SetFriendDND(userId, friendId string, DND int) error {
	err := mysql.SetFriendDND(userId, friendId, DND)
	if err != nil {
		return err
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "DND", utility.ToString(DND))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func SetFriendIsTop(userId, friendId string, isTop int) error {
	err := mysql.SetFriendIsTop(userId, friendId, isTop)
	if err != nil {
		return err
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "top", utility.ToString(isTop))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}

func FindAddFriendConfByUserId(userId string) (*types.AddFriendConf, error) {
	conf, err := cache.Cache.GetAddFriendConfig(userId)
	if err != nil {
		l.Warn("find add friend config from redis", "err", err)
	}

	if conf == nil {
		//从mysql 获取并更新缓存
		conf, err = mysql.FindAddFriendConfByUserId(userId)
		if conf == nil {
			l.Warn("add friend config can not find", "userId", userId)
			return nil, err
		}
		err = cache.Cache.SaveAddFriendConfig(userId, conf)
		if err != nil {
			l.Warn("redis can not save user")
		}
		return conf, nil
	}
	return conf, nil
}

func SetQuestionandAnswer(userId, question, answer string) error {
	err := mysql.SetQuestionandAnswer(userId, question, answer)
	if err != nil {
		return err
	}
	conf, err := cache.Cache.GetAddFriendConfig(userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	if conf != nil {
		conf.Question = question
		conf.Answer = answer

		err = cache.Cache.SaveAddFriendConfig(userId, conf)
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}
	return nil
}

func SetNeedAnswer(userId, question, answer string) error {
	err := mysql.SetNeedAnswer(userId, question, answer)
	if err != nil {
		return err
	}
	conf, err := cache.Cache.GetAddFriendConfig(userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	if conf != nil {
		conf.NeedAnswer = utility.ToString(types.NeedAnswer)
		conf.Question = question
		conf.Answer = answer

		err = cache.Cache.SaveAddFriendConfig(userId, conf)
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}
	return nil
}

func SetNotNeedAnswer(userId string) error {
	err := mysql.SetNotNeedAnswer(userId)
	if err != nil {
		return err
	}
	conf, err := cache.Cache.GetAddFriendConfig(userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	if conf != nil {
		conf.NeedAnswer = utility.ToString(types.NotNeedAnswer)

		err = cache.Cache.SaveAddFriendConfig(userId, conf)
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}
	return nil
}

func IsNeedConfirm(userId string, state int) error {
	err := mysql.IsNeedConfirm(userId, state)
	if err != nil {
		return err
	}
	conf, err := cache.Cache.GetAddFriendConfig(userId)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	if conf != nil {
		conf.NeedConfirm = utility.ToString(state)

		err = cache.Cache.SaveAddFriendConfig(userId, conf)
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}
	return nil
}

//TODO 未测试
func AddPrivateChatLog(log *types.PrivateLog) (int64, error) {
	_, logId, err := mysql.AddPrivateChatLog(log.SenderId, log.ReceiveId, log.MsgId, log.MsgType, log.Status, log.IsSnap, log.Content, log.Ext, log.SendTime)
	if err != nil {
		return 0, err
	}
	log.Id = utility.ToString(logId)
	err = cache.Cache.SavePrivateChatLogs([]*types.PrivateLog{log})
	if err != nil {

	}
	return logId, nil
}

func FindPrivateChatLogById(id string) (*types.PrivateLogJoinUser, error) {
	pl, err := cache.Cache.GetPrivateChatLog(id)
	if err != nil {
	}
	var log *types.PrivateLogJoinUser
	if pl != nil {
		user, err := GetUserInfoById(pl.SenderId)
		if err != nil {
			return nil, err
		}
		log = &types.PrivateLogJoinUser{
			PrivateLog: pl,
			User:       user,
		}
	}
	if log == nil {
		//从mysql 获取并更新缓存
		log, err = mysql.FindPrivateChatLogById(id)
		if log == nil {
			l.Warn("private chat log can not find")
			return nil, err
		}
		err = cache.Cache.SavePrivateChatLogs([]*types.PrivateLog{log.PrivateLog})
		if err != nil {
			l.Warn("redis can not save user")
		}
		return log, nil
	}
	return log, nil
}

func UpdatePrivateLogContentById(logId, content string) error {
	err := mysql.UpdatePrivateLogContentById(logId, content)
	if err != nil {
		return err
	}

	pl, err := FindPrivateChatLogById(logId)
	if err != nil {
		return err
	}
	pl.Content = content

	err = cache.Cache.SavePrivateChatLogs([]*types.PrivateLog{pl.PrivateLog})
	if err != nil {
		l.Warn("redis can not change private chat log")
	}
	return nil
}

//TODO 未测试
func DelPrivateChatLog(logId string) (int, error) {
	num, err := mysql.DelPrivateChatLog(logId)
	if err != nil {
		return num, err
	}

	return cache.Cache.DeletePrivateChatLog(logId)
}

//TODO 未测试
func FindPrivateChatLogs(userId, friendId string, start int64, number int) (int64, []*types.PrivateLogJoinUser, error) {

	/*conf, err := cache.Cache.GetPrivateChatLogsIndexByTime(start, nil)

	for range {

	}*/
	return 0, nil, nil
}

func FindNotBurnedLogsAfter(userId string, isDel int, time int64) ([]*types.PrivateLogJoinUser, error) {
	ids, err := cache.Cache.GetPrivateChatLogsIndexByTime(&time, nil, false, false)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, id := range ids {
		log, err := FindPrivateChatLogById(id)
		if err != nil {
			return nil, err
		}
		if log.IsDelete != isDel {
			continue
		}
		if log.Status == types.HadBurnt {
			continue
		}
		if log.ReceiveId != userId && log.SenderId != userId {
			continue
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindNotBurnedLogsBetween(userId string, isDel int, begin, end int64) ([]*types.PrivateLogJoinUser, error) {
	ids, err := cache.Cache.GetPrivateChatLogsIndexByTime(&begin, &end, true, false)
	if err != nil {
		return nil, err
	}

	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, id := range ids {
		log, err := FindPrivateChatLogById(id)
		if err != nil {
			return nil, err
		}
		if log.IsDelete != isDel {
			continue
		}
		if log.Status == types.HadBurnt {
			continue
		}
		if log.ReceiveId != userId && log.SenderId != userId {
			continue
		}
		logs = append(logs, log)
	}
	By(descById).Sort(logs)
	return logs, nil
}

func UpdatePrivateLogStateById(logId string, state int) error {
	err := mysql.UpdatePrivateLogStateById(logId, state)
	if err != nil {
		return err
	}

	pl, err := FindPrivateChatLogById(logId)
	if err != nil {
		return err
	}
	pl.Status = state

	err = cache.Cache.SavePrivateChatLogs([]*types.PrivateLog{pl.PrivateLog})
	if err != nil {
		l.Warn("redis can not change private chat log")
	}
	return nil
}

func FindSessionKeyAlert(userId string, endTime int64) ([]*types.PrivateLogJoinUser, error) {
	ids, err := cache.Cache.GetPrivateChatLogsIndexByTime(nil, &endTime, true, true)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, id := range ids {
		log, err := FindPrivateChatLogById(id)
		if err != nil {
			return nil, err
		}
		if log.ReceiveId != userId {
			continue
		}
		if log.MsgType != types.Alert {
			continue
		}
		//正则表达式
		if b, _ := regexp.MatchString(`"type":19`, log.Content); !b {
			continue
		}
		logs = append(logs, log)
	}
	return logs, nil
}

//设置黑名单
func SetFriendIsBlock(userId, friendId string, state int, alterTime int64) error {
	err := mysql.SetFriendIsBlock(userId, friendId, state, alterTime)
	if err != nil {
		return err
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "is_blocked", utility.ToString(state))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}

	err = cache.Cache.UpdateFriendInfo(userId, friendId, "add_time", utility.ToString(alterTime))
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return nil
}
