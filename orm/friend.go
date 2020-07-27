package orm

import (
	mysql "github.com/33cn/chat33/orm/mysql_model"
	redis "github.com/33cn/chat33/orm/redis_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logFriend = log15.New("module", "model/friend")

//查询好友信息
func FindFriendById(userID, friendID string) (*types.FriendJoinUser, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindFriendById(userID, friendID)
		if err != nil {
			logFriend.Error("redis.FindFriendById", "err", err, "userID", userID, "friendID", friendID)
		}
		return ret, err
	}
	ret, err := mysql.FindFriendById(userID, friendID)
	if err != nil {
		logFriend.Error("mysql.FindFriendById", "err", err, "userID", userID, "friendID", friendID)
	}
	return ret, err
}

func AcceptFriend(userId, friendId, source string, addTime int64) error {
	if cfg.CacheType.Enable {
		err := redis.AcceptFriend(userId, friendId, source, addTime)
		if err != nil {
			logFriend.Error("redis.AcceptFriend", "err", err, "userId", userId, "friendId", friendId, "source", source, "addTime", addTime)
		}
		return err
	}
	err := mysql.AcceptFriend(userId, friendId, source, addTime)
	if err != nil {
		logFriend.Error("mysql.AcceptFriend", "err", err, "userId", userId, "friendId", friendId, "source", source, "addTime", addTime)
	}
	return err
}

func FindPrivateChatLogById(id string) (*types.PrivateLogJoinUser, error) {
	/*if cfg.CacheType.Enable {
		return redis.FindPrivateChatLogById(id)
	}*/
	ret, err := mysql.FindPrivateChatLogById(id)
	if err != nil {
		logFriend.Error("mysql.FindPrivateChatLogById", "err", err, "logId", id)
	}
	return ret, err
}

func FindPrivateChatLogByMsgId(senderId, msgId string) (*types.PrivateLogJoinUser, error) {
	ret, err := mysql.FindPrivateChatLogByMsgId(senderId, msgId)
	if err != nil {
		logFriend.Error("mysql.FindPrivateChatLogById", "err", err, "senderId", senderId, "msgId", msgId)
	}
	return ret, err
}

func UpdatePrivateLogContentById(logId, content string) error {
	/*if cfg.CacheType.Enable {
		return redis.UpdatePrivateLogContentById(logId, content)
	}*/
	err := mysql.UpdatePrivateLogContentById(logId, content)
	if err != nil {
		logFriend.Error("mysql.FindPrivateChatLogById", "err", err, "logId", logId, "content", content)
	}
	return err
}

func FindPrivateChatLogs(userId, friendId string, start int64, number int) (int64, []*types.PrivateLogJoinUser, error) {
	/*if cfg.CacheType.Enable {
		return redis.FindPrivateChatLogs(userId, friendId, start, number)
	}*/
	ret, ret2, err := mysql.FindPrivateChatLogs(userId, friendId, start, number)
	if err != nil {
		logFriend.Error("mysql.FindPrivateChatLogs", "err", err, "userId", userId, "friendId", friendId, "start", start, "number", number)
	}
	return ret, ret2, err
}

func FindSessionKeyAlert(userId string, endTime int64) ([]*types.PrivateLogJoinUser, error) {
	/*if cfg.CacheType.Enable {
		return redis.FindSessionKeyAlert(userId, endTime)
	}*/
	ret, err := mysql.FindSessionKeyAlert(userId, endTime)
	if err != nil {
		logFriend.Error("mysql.FindSessionKeyAlert", "err", err, "userId", userId, "endTime", endTime)
	}
	return ret, err
}

func FindNotBurnedLogsAfter(userId string, isDel int, time int64) ([]*types.PrivateLogJoinUser, error) {
	/*if cfg.CacheType.Enable {
		return redis.FindNotBurnedLogsAfter(userId, isDel, time)
	}*/
	ret, err := mysql.FindNotBurnedLogsAfter(userId, isDel, time)
	if err != nil {
		logFriend.Error("mysql.FindNotBurnedLogsAfter", "err", err, "userId", userId, "isDel", isDel, "time", time)
	}
	return ret, err
}

func FindNotBurnedLogsBetween(userId string, isDel int, begin, end int64) ([]*types.PrivateLogJoinUser, error) {
	/*if cfg.CacheType.Enable {
		return redis.FindNotBurnedLogsBetween(userId, isDel, begin, end)
	}*/
	ret, err := mysql.FindNotBurnedLogsBetween(userId, isDel, begin, end)
	if err != nil {
		logFriend.Error("mysql.FindNotBurnedLogsBetween", "err", err, "userId", userId, "isDel", isDel, "begin", begin, "end", end)
	}
	return ret, err
}

//查询[begin,end]间的消息数量
func FindPrivateChatLogsNumberBetween(userId string, isDel int, begin, end int64) (int64, error) {
	ret, err := mysql.FindPrivateChatLogsNumberBetween(userId, isDel, begin, end)
	if err != nil {
		logFriend.Error("mysql.FindPrivateChatLogsNumberBetween", "err", err, "userId", userId, "isDel", isDel, "begin", begin, "end", end)
	}
	return ret, err
}

func DelPrivateChatLog(logId string) (int, error) {
	/*if cfg.CacheType.Enable {
		return redis.DelPrivateChatLog(logId)
	}*/
	ret, err := mysql.DelPrivateChatLog(logId)
	if err != nil {
		logFriend.Error("mysql.DelPrivateChatLog", "err", err, "logId", logId)
	}
	return ret, err
}

func AlertPrivateRevStateByRevId(revId string, state int) error {
	/*if cfg.CacheType.Enable {
		return redis.UpdatePrivateLogStateById(revId, state)
	}*/
	err := mysql.UpdatePrivateLogStateById(revId, state)
	if err != nil {
		logFriend.Error("mysql.UpdatePrivateLogStateById", "err", err, "revId", revId, "state", state)
	}
	return err
}

func CheckIsFriend(userId, friendId string, isDel int) (bool, error) {
	//if cfg.CacheType.Enable {
	//	ret, err := redis.CheckIsFriend(userId, friendId)
	//	if err != nil {
	//		logFriend.Error("redis.CheckIsFriend", "err", err, "userId", userId, "friendId", friendId)
	//	}
	//	return ret, err
	//}
	ret, err := mysql.CheckIsFriend(userId, friendId, isDel)
	if err != nil {
		logFriend.Error("mysql.CheckIsFriend", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

func FindFriendsById(userId string, commonUse, isDel int) ([]*types.FriendJoinUser, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindFriendsFilterByCommonUse(userId, commonUse)
		if err != nil {
			logFriend.Error("redis.FindFriendsFilterByCommonUse", "err", err, "userId", userId, "commonUse", commonUse, "isDel", isDel)
		}
		return ret, err
	}
	ret, err := mysql.FindFriendsById(userId, commonUse, isDel)
	if err != nil {
		logFriend.Error("mysql.FindFriendsFilterByCommonUse", "err", err, "userId", userId, "commonUse", commonUse, "isDel", isDel)
	}
	return ret, err
}

func FindFriendsAfterTime(userId string, commonUse int, time int64) ([]*types.FriendJoinUser, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindFriendsAfterTime(userId, commonUse, time)
		if err != nil {
			logFriend.Error("redis.FindFriendsAfterTime", "err", err, "userId", userId, "commonUse", commonUse, "time", time)
		}
		return ret, err
	}
	ret, err := mysql.FindFriendsAfterTime(userId, commonUse, time)
	if err != nil {
		logFriend.Error("mysql.FindFriendsAfterTime", "err", err, "userId", userId, "commonUse", commonUse, "time", time)
	}
	return ret, err
}

func GetFriendApplyCount(userId, friendId string) (int32, error) {
	if cfg.CacheType.Enable {
		b, err := redis.CheckFriendApplyExists(userId, friendId)
		if err != nil {
			logFriend.Error("redis.CheckFriendApplyExists", "err", err, "userId", userId, "friendId", friendId)
			return 0, err
		}
		return int32(utility.BoolToInt(b)), nil
	}
	ret, err := mysql.FindFriendApplyCount(userId, friendId)
	if err != nil {
		logFriend.Error("mysql.FindFriendApplyCount", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

func InsertFriend(userID, friendID, remark, extRemark string, dnd, top int, addTime int64) error {
	err := mysql.InsertFriend(userID, friendID, remark, extRemark, dnd, top, addTime)
	if err != nil {
		logFriend.Error("mysql.SetFriendRemark", "err", err, "userId", userID, "friendId", friendID, "remark", remark, "DND", dnd, "top", top, "addtime", addTime)
	}
	return err
}

func SetFriendRemark(userId, friendId, remark string) error {
	if cfg.CacheType.Enable {
		err := redis.SetFriendRemark(userId, friendId, remark)
		if err != nil {
			logFriend.Error("redis.SetFriendRemark", "err", err, "userId", userId, "friendId", friendId, "remark", remark)
		}
		return err
	}
	err := mysql.SetFriendRemark(userId, friendId, remark)
	if err != nil {
		logFriend.Error("mysql.SetFriendRemark", "err", err, "userId", userId, "friendId", friendId, "remark", remark)
	}
	return err
}

func SetFriendExtRemark(userId, friendId, remark string) error {
	if cfg.CacheType.Enable {
		err := redis.SetFriendExtRemark(userId, friendId, remark)
		if err != nil {
			logFriend.Error("redis.SetFriendExtRemark", "err", err, "userId", userId, "friendId", friendId, "remark", remark)
		}
		return err
	}
	err := mysql.SetFriendExtRemark(userId, friendId, remark)
	if err != nil {
		logFriend.Error("mysql.SetFriendExtRemark", "err", err, "userId", userId, "friendId", friendId, "remark", remark)
	}
	return err
}

func SetFriendDND(userId, friendId string, DND int) error {
	if cfg.CacheType.Enable {
		err := redis.SetFriendDND(userId, friendId, DND)
		if err != nil {
			logFriend.Error("redis.SetFriendDND", "err", err, "userId", userId, "friendId", friendId, "DND", DND)
		}
		return err
	}
	err := mysql.SetFriendDND(userId, friendId, DND)
	if err != nil {
		logFriend.Error("mysql.SetFriendDND", "err", err, "userId", userId, "friendId", friendId, "DND", DND)
	}
	return err
}

func SetFriendIsTop(userId, friendId string, isTop int) error {
	if cfg.CacheType.Enable {
		err := redis.SetFriendIsTop(userId, friendId, isTop)
		if err != nil {
			logFriend.Error("redis.SetFriendIsTop", "err", err, "userId", userId, "friendId", friendId, "isTop", isTop)
		}
		return err
	}
	err := mysql.SetFriendIsTop(userId, friendId, isTop)
	if err != nil {
		logFriend.Error("mysql.SetFriendIsTop", "err", err, "userId", userId, "friendId", friendId, "isTop", isTop)
	}
	return err
}

func SetQuestionandAnswer(userId, question, answer string) error {
	if cfg.CacheType.Enable {
		err := redis.SetQuestionandAnswer(userId, question, answer)
		if err != nil {
			logFriend.Error("redis.SetQuestionandAnswer", "err", err, "userId", userId, "question", question, "answer", answer)
		}
		return err
	}
	err := mysql.SetQuestionandAnswer(userId, question, answer)
	if err != nil {
		logFriend.Error("mysql.SetQuestionandAnswer", "err", err, "userId", userId, "question", question, "answer", answer)
	}
	return err
}

func SetNeedAnswer(userId, question, answer string) error {
	if cfg.CacheType.Enable {
		err := redis.SetNeedAnswer(userId, question, answer)
		if err != nil {
			logFriend.Error("redis.SetNeedAnswer", "err", err, "userId", userId, "question", question, "answer", answer)
		}
		return err
	}
	err := mysql.SetNeedAnswer(userId, question, answer)
	if err != nil {
		logFriend.Error("mysql.SetNeedAnswer", "err", err, "userId", userId, "question", question, "answer", answer)
	}
	return err
}

func SetNotNeedAnswer(userId string) error {
	if cfg.CacheType.Enable {
		err := redis.SetNotNeedAnswer(userId)
		if err != nil {
			logFriend.Error("mysql.SetNotNeedAnswer", "err", err, "userId", userId)
		}
		return err
	}
	err := mysql.SetNotNeedAnswer(userId)
	if err != nil {
		logFriend.Error("mysql.SetNotNeedAnswer", "err", err, "userId", userId)
	}
	return err
}

func IsNeedConfirm(userId string, state int) error {
	if cfg.CacheType.Enable {
		err := redis.IsNeedConfirm(userId, state)
		if err != nil {
			logFriend.Error("redis.IsNeedConfirm", "err", err, "userId", userId, "state", state)
		}
		return err
	}
	err := mysql.IsNeedConfirm(userId, state)
	if err != nil {
		logFriend.Error("mysql.IsNeedConfirm", "err", err, "userId", userId, "state", state)
	}
	return err
}

func DeleteFriend(userId, friendId string, alterTime int64) error {
	if cfg.CacheType.Enable {
		err := redis.DeleteFriend(userId, friendId, alterTime)
		if err != nil {
			logFriend.Error("redis.DeleteFriend", "err", err, "userId", userId, "friendId", friendId, "alterTime", alterTime)
		}
		return err
	}
	err := mysql.DeleteFriend(userId, friendId, alterTime)
	if err != nil {
		logFriend.Error("mysql.DeleteFriend", "err", err, "userId", userId, "friendId", friendId, "alterTime", alterTime)
	}
	return err
}

func FindTypicalPrivateLog(userId, friendId, owner string, start int64, number int, queryType []string) (int64, []*types.PrivateLogJoinUser, error) {
	ret, ret2, err := mysql.FindTypicalPrivateChatLogs(userId, friendId, owner, start, number, queryType)
	if err != nil {
		logFriend.Error("mysql.FindTypicalPrivateChatLogs", "err", err, "userId", userId, "friendId", friendId, "owner", owner, "start", start, "number", number, "queryType", queryType)
	}
	return ret, ret2, err
}

func FindFirstPrivateMsg(userId, friendId string) (*types.PrivateLog, error) {
	ret, err := mysql.FindFirstPrivateMsg(userId, friendId)
	if err != nil {
		logFriend.Error("mysql.FindFirstPrivateMsg", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

func FindLastCatLogId(userId, friendId string) (int64, error) {
	ret, err := mysql.FindLastCatLogId(userId, friendId)
	if err != nil {
		logFriend.Error("mysql.FindLastCatLogId", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

func FindFriendsId(userId string) ([]string, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindFriends(userId)
		if err != nil {
			logFriend.Error("redis.FindFriends", "err", err, "userId", userId)
		}
		return ret, err
	}
	ret, err := mysql.FindFriendsId(userId)
	if err != nil {
		logFriend.Error("mysql.FindFriendsId", "err", err, "userId", userId)
	}
	return ret, err
}

func GetUnReadNumber(userId, friendId string) (int32, error) {
	ret, err := mysql.GetUnReadNumber(userId, friendId)
	if err != nil {
		logFriend.Error("mysql.GetUnReadNumber", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

func FindAddFriendConfByUserId(userId string) (*types.AddFriendConf, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindAddFriendConfByUserId(userId)
		if err != nil {
			logFriend.Error("redis.FindAddFriendConfByUserId", "err", err, "userId", userId)
		}
		return ret, err
	}
	ret, err := mysql.FindAddFriendConfByUserId(userId)
	if err != nil {
		logFriend.Error("mysql.FindAddFriendConfByUserId", "err", err, "userId", userId)
	}
	return ret, err
}

func AddPrivateChatLog(senderId, receiveId, msgId string, msgType, status, isSnap int, content, ext string, time int64) (int64, error) {
	/*if cfg.CacheType.Enable {
		param := types.PrivateLog{
			SendTime:  time,
			SenderId:  senderId,
			ReceiveId: receiveId,
			MsgType:   msgType,
			MsgId:     msgId,
			Status:    status,
			IsSnap:    isSnap,
			Content:   content,
			IsDelete:  types.FriendMsgIsNotDelete,
		}
		return redis.AddPrivateChatLog(&param)
	}*/
	_, logId, err := mysql.AddPrivateChatLog(senderId, receiveId, msgId, msgType, status, isSnap, content, ext, time)
	if err != nil {
		logFriend.Error("mysql.AddPrivateChatLog", "err", err)
	}
	return logId, err
}

//设置是否进入黑名单
func SetFriendIsBlock(userId, friendId string, state int) error {
	now := utility.NowMillionSecond()
	if cfg.CacheType.Enable {
		err := redis.SetFriendIsBlock(userId, friendId, state, now)
		if err != nil {
			logFriend.Error("redis.SetFriendIsBlock", "err", err, "userId", userId, "friendId", friendId, "state", state)
		}
		return err
	}
	err := mysql.SetFriendIsBlock(userId, friendId, state, now)
	if err != nil {
		logFriend.Error("mysql.SetFriendIsBlock", "err", err, "userId", userId, "friendId", friendId, "state", state)
	}
	return err
}

//黑名单列表，包含已经删除的好友
func BlockedFriends(userId string) ([]*types.FriendJoinUser, error) {
	if cfg.CacheType.Enable {
		//return redis.FindFriendsFilterByCommonUse(userId, commonUse)
	}
	ret, err := mysql.FindBlockedFriends(userId)
	if err != nil {
		logFriend.Error("mysql.FindBlockedFriends", "err", err, "userId", userId)
	}
	return ret, err
}
