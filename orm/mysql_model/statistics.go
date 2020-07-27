package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func FindApplyLogById(logId int64) (*types.Apply, error) {
	maps, err := db.GetApplyById(logId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return &types.Apply{
		Id:          info["id"],
		Type:        utility.ToInt(info["type"]),
		ApplyReason: info["apply_reason"],
		State:       utility.ToInt(info["state"]),
		Datetime:    utility.ToInt64(info["datetime"]),
		Target:      info["target"],
		ApplyUser:   info["apply_user"],
		Source:      info["source"],
	}, nil
}

func FindApplyLogByUserAndTarget(userId, targetId string, tp int) (*types.Apply, error) {
	maps, err := db.GetApplyByUserAndTarget(userId, targetId, tp)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return &types.Apply{
		Id:          info["id"],
		Type:        utility.ToInt(info["type"]),
		ApplyReason: info["apply_reason"],
		State:       utility.ToInt(info["state"]),
		Datetime:    utility.ToInt64(info["datetime"]),
		Target:      info["target"],
		ApplyUser:   info["apply_user"],
		Source:      info["source"],
	}, nil
}

func FindApplyLogs(userId string, id int64, number int) (int64, []*types.Apply, error) {
	//向前分页查询, 找出id号之前的number条记录，并返回下一条记录id
	maps, err := db.GetApplyList(userId, id, number+1)
	if err != nil {
		return -1, nil, err
	}
	if len(maps) < 1 {
		return -1, nil, nil
	}
	var nextId int64 = -1
	if len(maps) == number+1 {
		nextId = utility.ToInt64(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	var logs = make([]*types.Apply, 0)
	for _, v := range maps {
		log := &types.Apply{
			Id:          v["id"],
			Type:        utility.ToInt(v["type"]),
			ApplyReason: v["apply_reason"],
			State:       utility.ToInt(v["state"]),
			Datetime:    utility.ToInt64(v["datetime"]),
			Target:      v["target"],
			ApplyUser:   v["apply_user"],
			Source:      v["source"],
		}
		logs = append(logs, log)
	}
	return nextId, logs, nil
}

func GetAddFriendApplyLog(userId, friendId string) (*types.Apply, error) {
	maps, err := db.FindApplyOrderByTime(userId, friendId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return &types.Apply{
		Id:          info["id"],
		Type:        utility.ToInt(info["type"]),
		ApplyReason: info["apply_reason"],
		State:       utility.ToInt(info["state"]),
		Datetime:    utility.ToInt64(info["datetime"]),
		Target:      info["target"],
		ApplyUser:   info["apply_user"],
		Source:      info["source"],
	}, nil
}

func GetApplyListNumber(id string) (int, error) {
	return db.GetApplyListNumber(id)
}

func GetUnreadApplyNumber(userId string) (int, error) {
	return db.GetUnreadApplyNumber(userId)
}

func UpdateApply(applyUser, target string, tp int, reason, remark, source string, state int) (int64, error) {
	return db.UpdateApply(applyUser, target, tp, reason, remark, source, state)
}

func AcceptFriendApply(userId, friendId string) (int64, error) {
	return db.AcceptFriendApply(userId, friendId)
}

//添加入群申请记录
func AppendApplyLog(targetId, userId, applyReason, source, remark string, state, tp int, datetime int64) (int64, error) {
	_, logId, err := db.AppendApplyLog(targetId, userId, applyReason, source, remark, state, tp, datetime)
	return logId, err
}
