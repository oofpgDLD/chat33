package db

import (
	"fmt"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

//获取指定好友或者群中所有人对你的奖赏列表
func GetPraiseListByTarget(msgSenderId string, channelType int, targetId string, startId int64, number int) ([]map[string]string, int64, error) {
	var nextId int64 = -1
	condition := ""
	if startId != 0 {
		condition = fmt.Sprintf("and id <= %v", startId)
	}
	sqlStr := fmt.Sprintf("select * from praise where is_delete = ? and channel_type = ? and target_id = ? and sender_id = ? %s order by id desc limit 0,?", condition)
	maps, err := conn.Query(sqlStr, types.PraiseIsNotDelete, channelType, targetId, msgSenderId, number+1)
	if err != nil {
		return nil, nextId, err
	}
	if len(maps) > number {
		nextId = utility.ToInt64(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextId, err
}

//查询某个人对某条消息的赞赏情况
func GetPraiseByLogIdAndOptId(channelType int, logId, optId string) ([]map[string]string, error) {
	const sqlStr = "select * from praise where is_delete = ? and channel_type = ? and log_id = ? and opt_id = ?"
	return conn.Query(sqlStr, types.PraiseIsNotDelete, channelType, logId, optId)
}

func GetPraiseByLogId(channelType int, logId string) ([]map[string]string, error) {
	const sqlStr = "select * from praise where is_delete = ? and channel_type = ? and log_id = ?"
	return conn.Query(sqlStr, types.PraiseIsNotDelete, channelType, logId)
}

func GetPraiseByLogIdLimit(channelType int, logId string, startId int64, number int) ([]map[string]string, int64, error) {
	var nextId int64 = -1
	condition := ""
	if startId != 0 {
		condition = fmt.Sprintf("and id <= %v", startId)
	}
	sqlStr := fmt.Sprintf("select * from praise where is_delete = ? and channel_type = ? and log_id = ? %s order by id desc limit 0,?", condition)
	maps, err := conn.Query(sqlStr, types.PraiseIsNotDelete, channelType, logId, number+1)
	if err != nil {
		return nil, nextId, err
	}
	if len(maps) > number {
		nextId = utility.ToInt64(maps[len(maps)-1]["id"])
		maps = maps[:len(maps)-1]
	}
	return maps, nextId, err
}

func LikeOrRewardCount(channelType int, logId string, tp int) ([]map[string]string, error) {
	const sqlStr = "select count(*) as count from praise where is_delete = ? and channel_type = ? and log_id = ? and type = ?"
	return conn.Query(sqlStr, types.PraiseIsNotDelete, channelType, logId, tp)
}

func InsertPraise(val *types.Praise) (int64, int64, error) {
	const sqlStr = "insert into praise(channel_type,target_id,log_id,sender_id,opt_id,type,record_id,coin_id,coin_name,amount,create_time,is_delete) values(?,?,?,?,?,?,?,?,?,?,?,?)"
	return conn.Exec(sqlStr,
		val.ChannelType,
		val.TargetId,
		val.LogId,
		val.SenderId,
		val.OptId,
		val.Type,
		val.RecordId,
		val.CoinId,
		val.CoinName,
		val.Amount,
		val.CreateTime,
		types.PraiseIsNotDelete,
	)
}

func InsertPraiseUser(val *types.PraiseUser) (int64, int64, error) {
	const sqlStr = "insert into praise_user(target_id,opt_id,record_id,coin_id,coin_name,amount,create_time,is_delete) values(?,?,?,?,?,?,?,?)"
	return conn.Exec(sqlStr,
		val.TargetId,
		val.OptId,
		val.RecordId,
		val.CoinId,
		val.CoinName,
		val.Amount,
		val.CreateTime,
		types.PraiseIsNotDelete,
	)
}

func DeletePraise(id string) (int64, int64, error) {
	const sqlStr = "update praise set is_delete = ? where id = ?"
	return conn.Exec(sqlStr, types.PraiseIsDelete, id)
}

func GetPraiseTodayLimit(optId string) ([]map[string]string, error) {
	const sqlStr = "SELECT amount,coin_name FROM `praise` WHERE opt_id = ? AND create_time > UNIX_TIMESTAMP(CAST(SYSDATE() AS DATE))*1000 AND amount > 0"
	return conn.Query(sqlStr, optId)
}

/*func GetLeaderBoardAsLike(Type int, startTime, endTime int64) ([]map[string]string, error) {
	const sqlStr = "SELECT amount,coin_name,sender_id,type FROM `praise` WHERE type = ? AND create_time >= ? AND create_time <= ?"
	return conn.Query(sqlStr, Type, startTime, endTime)
}*/

func GetLeaderBoardAsLike(startTime, endTime int64) ([]map[string]string, error) {
	const sqlStr = "SELECT sender_id, count(*) as number FROM `praise` WHERE type = ? and is_delete = ? AND create_time >= ? AND create_time <= ? GROUP BY sender_id"
	return conn.Query(sqlStr, types.Like, types.PraiseIsNotDelete, startTime, endTime)
}

func GetLeaderBoardAsReward(startTime, endTime int64) ([]map[string]string, error) {
	const sqlStr = "SELECT sender_id, coin_name, SUM(amount) as amount FROM `praise` WHERE type = ? and is_delete = ? AND create_time >= ? AND create_time <= ? GROUP BY sender_id,coin_name "
	return conn.Query(sqlStr, types.Reward, types.PraiseIsNotDelete, startTime, endTime)
}
