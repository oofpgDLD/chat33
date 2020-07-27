package orm

import (
	mysql "github.com/33cn/chat33/orm/mysql_model"
	redis "github.com/33cn/chat33/orm/redis_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logStatic = log15.New("model", "orm/statistics")

func FindApplyLogById(logId int64) (*types.Apply, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindApplyLogById(logId)
		if err != nil {
			logStatic.Error("redis.FindApplyLogById", "err", err, "logId", logId)
		}
		return ret, err
	}
	ret, err := mysql.FindApplyLogById(logId)
	if err != nil {
		logStatic.Error("mysql.FindApplyLogById", "err", err, "logId", logId)
	}
	return ret, err
}

func FindAddFriendApplyLog(userId, friendId string) (*types.Apply, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.GetAddFriendApplyLog(userId, friendId)
		if err != nil {
			logStatic.Error("redis.GetAddFriendApplyLog", "err", err, "userId", userId, "friendId", friendId)
		}
		return ret, err
	}
	ret, err := mysql.GetAddFriendApplyLog(userId, friendId)
	if err != nil {
		logStatic.Error("mysql.GetAddFriendApplyLog", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

func FindApplyLogByUserAndTarget(userId, targetId string, tp int) (*types.Apply, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.FindApplyLogByUserAndTarget(userId, targetId, tp)
		if err != nil {
			logStatic.Error("redis.FindApplyLogByUserAndTarget", "err", err, "userId", userId, "targetId", targetId, "tp", tp)
		}
		return ret, err
	}
	ret, err := mysql.FindApplyLogByUserAndTarget(userId, targetId, tp)
	if err != nil {
		logStatic.Error("mysql.FindApplyLogByUserAndTarget", "err", err, "userId", userId, "targetId", targetId, "tp", tp)
	}
	return ret, err
}

func FindApplyLogs(userId string, id int64, number int) (int64, []*types.Apply, error) {
	ret, ret2, err := mysql.FindApplyLogs(userId, id, number)
	if err != nil {
		logStatic.Error("mysql.FindApplyLogs", "err", err, "userId", userId, "id", id, "number", number)
	}
	return ret, ret2, err
}

//返回：更改记录的Id，error
func UpdateApply(applyUser, target string, tp int, reason, remark, source string, state int) (int64, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.UpdateApply(applyUser, target, tp, reason, remark, source, state)
		if err != nil {
			logStatic.Error("redis.UpdateApply", "err", err, "applyUser", applyUser, "target", target, "tp", tp, "reason", reason, "remark", remark, "source", source, "state", state)
		}
		return ret, err
	}
	ret, err := mysql.UpdateApply(applyUser, target, tp, reason, remark, source, state)
	if err != nil {
		logStatic.Error("mysql.UpdateApply", "err", err, "applyUser", applyUser, "target", target, "tp", tp, "reason", reason, "remark", remark, "source", source, "state", state)
	}
	return ret, err
}

//已经是好友更新所有请求
func AcceptFriendApply(userId, friendId string) (int64, error) {
	if cfg.CacheType.Enable {
		ret, err := redis.AcceptFriendApply(userId, friendId)
		if err != nil {
			logStatic.Error("redis.AcceptFriendApply", "err", err, "userId", userId, "friendId", friendId)
		}
		return ret, err
	}
	ret, err := mysql.AcceptFriendApply(userId, friendId)
	if err != nil {
		logStatic.Error("mysql.AcceptFriendApply", "err", err, "userId", userId, "friendId", friendId)
	}
	return ret, err
}

//添加请求记录
func AppendApplyLog(targetId, userId, applyReason, source, remark string, state, tp int) (int64, error) {
	datetime := utility.NowMillionSecond()
	if cfg.CacheType.Enable {
		ret, err := redis.AppendApplyLog(targetId, userId, applyReason, source, remark, state, tp, datetime)
		if err != nil {
			logStatic.Error("redis.AppendApplyLog", "err", err, "userId", userId, "targetId", targetId)
		}
		return ret, err
	}
	ret, err := mysql.AppendApplyLog(targetId, userId, applyReason, source, remark, state, tp, datetime)
	if err != nil {
		logStatic.Error("mysql.AppendApplyLog", "err", err, "userId", userId, "targetId", targetId)
	}
	return ret, err
}

func GetApplyListNumber(id string) (int, error) {
	ret, err := mysql.GetApplyListNumber(id)
	if err != nil {
		logStatic.Error("mysql.GetApplyListNumber", "err", err, "id", id)
	}
	return ret, err
}

func GetUnreadApplyNumber(caller string) (interface{}, error) {
	number, err := mysql.GetUnreadApplyNumber(caller)
	if err != nil {
		logStatic.Error("mysql.GetUnreadApplyNumber", "err", err, "caller", caller)
		return nil, err
	}
	var ret = make(map[string]interface{})
	ret["number"] = number
	return ret, nil
}
