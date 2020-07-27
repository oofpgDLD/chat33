package db

import (
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

// 获取申请列表
func GetApplyList(userId string, id int64, number int) ([]map[string]string, error) {
	if id == 0 {
		var sqlStr = "SELECT * FROM `apply` where apply_user = ? or (target = ? and `type` = 2) or (target in(select room_id from room_user where user_id = ? and `level` > 1 and is_delete = ?) and `type` = 1) ORDER BY datetime desc limit 0,?"
		return conn.Query(sqlStr, userId, userId, userId, types.RoomUserNotDeleted, number)
	} else {
		var sqlStr = "SELECT * FROM `apply` WHERE (apply_user = ? or (target = ? and `type` = 2) or (target in(select room_id from room_user where user_id = ? and `level` > 1 and is_delete = ?) and `type` = 1)) and datetime <= (SELECT datetime FROM `apply` WHERE id = ?) ORDER BY datetime desc limit 0,?"
		return conn.Query(sqlStr, userId, userId, userId, types.RoomUserNotDeleted, id, number)
	}
}

func GetApplyById(id int64) ([]map[string]string, error) {
	const sqlStr = "SELECT * FROM `apply` WHERE id = ?"
	return conn.Query(sqlStr, id)
}

func GetApplyByUserAndTarget(userId, targetId string, tp int) ([]map[string]string, error) {
	sqlStr := "SELECT * FROM apply WHERE apply_user = ? AND target = ? and type = ? order by id desc"
	return conn.Query(sqlStr, userId, targetId, tp)
}

func GetApplyListNumber(uid string) (int, error) {
	const sqlStr = "SELECT COUNT(*) as count FROM `apply` where apply_user = ? or target = ?"
	maps, err := conn.Query(sqlStr, uid, uid)
	if err != nil {
		return 0, err
	}
	count := utility.ToInt(maps[0]["count"])
	return count, err
}

func GetUnreadApplyNumber(userId string) (int, error) {
	const sqlStr = "SELECT count(*) as count FROM `apply` where ((target = ? and `type` = 2) or (target in(select id from room where id in (select room_id from room_user where user_id = ? and `level` > 1 and is_delete = ?) and room.is_delete = ?) and `type` = 1)) and state = ?"
	maps, err := conn.Query(sqlStr, userId, userId, types.RoomUserNotDeleted, types.RoomNotDeleted, types.AwaitState)
	if len(maps) < 1 || err != nil {
		return 0, err
	}
	return utility.ToInt(maps[0]["count"]), err
}

//添加申请记录 target：roomId 或者friendId
func AppendApplyLog(targetId, userId, applyReason, source, remark string, state, tp int, datetime int64) (int64, int64, error) {
	const sqlStr = "insert into `apply`(id,`type`,apply_user,target,apply_reason,state,remark,datetime,source) values(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE state = ?, datetime = ?, source = ?, apply_reason = ?, remark = ?"
	return conn.Exec(sqlStr, nil, tp, userId, targetId, applyReason, state, remark, datetime, source, state, datetime, source, applyReason, remark)
}

func GetRoomUserApplyInfo(roomId, userId string) ([]map[string]string, error) {
	const sqlFindApply = "select * from `apply` where type = ? and apply_user = ? and target = ?"
	return conn.Query(sqlFindApply, types.IsRoom, userId, roomId)
}

//更新好友申请数据
func UpdateApply(applyUser, target string, tp int, reason, remark, source string, state int) (int64, error) {
	const sqlStr = "update apply set source = ?,apply_reason = ?,state = ?,remark = ?,datetime = ? where apply_user = ? and target = ? and type = ?"
	_, id, err := conn.Exec(sqlStr, source, reason, state, remark, utility.NowMillionSecond(), applyUser, target, tp)
	return id, err
}
