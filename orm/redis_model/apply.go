package redis_model

import (
	"github.com/33cn/chat33/cache"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func UpdateApply(applyUser, target string, tp int, reason, remark, source string, state int) (int64, error) {
	id, err := mysql.UpdateApply(applyUser, target, tp, reason, remark, source, state)
	if err != nil {
		return 0, err
	}

	apply, err := FindApplyLogByUserAndTarget(applyUser, target, tp)
	if apply != nil {
		apply.ApplyUser = applyUser
		apply.Target = target
		apply.Type = tp
		apply.ApplyReason = reason
		apply.State = state
		apply.Remark = remark
		apply.Source = source

		err = cache.Cache.SaveApplyLog(apply)
		if err != nil {
			//TODO 是否删除该记录 使之再次读取的时候更新？
		}
	}
	return id, nil
}

func AcceptFriendApply(userId, friendId string) (int64, error) {
	num, err := mysql.AcceptFriendApply(userId, friendId)
	if err != nil {
		return 0, err
	}

	{
		apply, err := FindApplyLogByUserAndTarget(userId, friendId, types.IsFriend)
		if err != nil {
			//TODO 回滚
		}
		if apply != nil {
			apply.State = types.AcceptState
			err = cache.Cache.SaveApplyLog(apply)
			if err != nil {
				//TODO 回滚
			}
		}
	}

	{
		//获取 保存
		apply, err := FindApplyLogByUserAndTarget(friendId, userId, types.IsFriend)
		if err != nil {
			//TODO 回滚
		}
		if apply != nil {
			apply.State = types.AcceptState
			err = cache.Cache.SaveApplyLog(apply)
			if err != nil {
				//TODO 回滚
			}
		}
	}
	return num, nil
}

func AppendApplyLog(targetId, userId, applyReason, source, remark string, state, tp int, datetime int64) (int64, error) {
	id, err := mysql.AppendApplyLog(targetId, userId, applyReason, source, remark, state, tp, datetime)
	if err != nil {
		return 0, err
	}
	err = cache.Cache.SaveApplyLog(&types.Apply{
		Id:          utility.ToString(id),
		Target:      targetId,
		ApplyUser:   userId,
		ApplyReason: applyReason,
		Source:      source,
		Remark:      remark,
		State:       state,
		Type:        tp,
		Datetime:    datetime,
	})
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return id, nil
}

//根据id获取加好友/群请求记录
func FindApplyLogById(logId int64) (*types.Apply, error) {
	apply, err := cache.Cache.GetApplyLogById(utility.ToString(logId))
	if err != nil {
		l.Warn("get apply from redis failed", "err", err)
	}
	if apply == nil {
		//从mysql 获取并更新缓存
		apply, err := mysql.FindApplyLogById(logId)
		if apply == nil {
			l.Warn("apply can not find")
			return nil, err
		}
		err = cache.Cache.SaveApplyLog(apply)
		if err != nil {
			l.Warn("redis can not save apply")
		}
		return apply, nil
	}
	return apply, nil
}

func FindApplyLogByUserAndTarget(applyUser, target string, tp int) (*types.Apply, error) {
	apply, err := cache.Cache.GetApplyLogByUserAndTarget(applyUser, target, tp)
	if err != nil {
		l.Warn("get apply from redis failed", "err", err)
	}
	if apply == nil {
		//从mysql 获取并更新缓存
		apply, err := mysql.FindApplyLogByUserAndTarget(applyUser, target, tp)
		if apply == nil {
			l.Warn("apply can not find")
			return nil, err
		}
		err = cache.Cache.SaveApplyLog(apply)
		if err != nil {
			l.Warn("redis can not save apply")
		}
		return apply, nil
	}
	return apply, nil
}

func GetAddFriendApplyLog(userId, friendId string) (*types.Apply, error) {
	var apply *types.Apply
	a, err := FindApplyLogByUserAndTarget(userId, friendId, types.IsFriend)
	if err != nil {
		return nil, err
	}
	a2, err := FindApplyLogByUserAndTarget(friendId, userId, types.IsFriend)
	if err != nil {
		return nil, err
	}

	if a != nil && a2 != nil {
		if a.Datetime > a2.Datetime {
			apply = a
		} else {
			apply = a2
		}
	} else {
		if a != nil {
			apply = a
		} else if a2 != nil {
			apply = a2
		}
	}
	return apply, nil
}

func CheckFriendApplyExists(userId, friendId string) (bool, error) {
	apply, err := FindApplyLogByUserAndTarget(userId, friendId, types.IsFriend)
	if err != nil {
		return false, err
	}
	if apply == nil {
		return false, nil
	}
	return true, nil
}

func ApproveChangeStateStep(tx types.Tx, logId int64, status int) (bool, error) {
	ok, err := mysql.ApproveChangeStateStep(tx, logId, status)
	if err != nil {
		return ok, err
	}
	err = cache.Cache.UpdateApplyStateById(utility.ToString(logId), status)
	if err != nil {
		//TODO 是否删除该记录 使之再次读取的时候更新？
	}
	return ok, nil
}
