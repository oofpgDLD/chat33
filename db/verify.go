package db

import (
	"fmt"

	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

//----------------------------认证-----------------------//
func VerifySetFee(tx *mysql.MysqlTx, appId string, tp int, currency string, amount float64) error {
	const sqlStr = "INSERT INTO verify_fee(app_id,`type`,currency,amount) values(?,?,?,?) ON DUPLICATE KEY UPDATE currency = ?, amount = ?"
	_, _, err := tx.Exec(sqlStr, appId, tp, currency, amount, currency, amount)
	return err
}

func VerifyGetFee(appId string) ([]map[string]string, error) {
	const sqlStr = "SELECT * FROM verify_fee WHERE app_id = ? ORDER BY `type` ASC"
	return conn.Query(sqlStr, appId)
}

//更新 手续费到账信息
func UpdateVerifyFeeState(id string, feeState int) error {
	const sqlStr = "UPDATE verify_apply SET fee_state = ? WHERE id = ?"
	_, _, err := conn.Exec(sqlStr, feeState, id)
	return err
}

func PersonalVerifyList(appId string, search *string, start, end int64, state *int) (int64, []map[string]string, error) {
	stateQuery := ""
	if state != nil {
		stateQuery = fmt.Sprintf(" AND state = %v", *state)
	}

	searchQuery := ""
	if search != nil && *search != "" {
		*search = QueryStr(*search)
		searchQuery = fmt.Sprintf(" AND (u.uid LIKE '%%%v%%' OR u.phone LIKE '%%%v%%' OR u.username LIKE '%%%v%%')", *search, *search, *search)
	}
	sqlStr := fmt.Sprintf("SELECT vf.*,u.avatar,u.uid,u.phone,u.username FROM verify_apply AS vf LEFT JOIN `user` AS u ON vf.target_id = u.user_id"+
		" WHERE vf.type = ? AND u.app_id = ? %s ORDER BY update_time DESC LIMIT ?,?", stateQuery+searchQuery)

	countSqlStr := fmt.Sprintf("SELECT COUNT(*) AS count FROM verify_apply AS vf LEFT JOIN `user` AS u ON vf.target_id = u.user_id"+
		" WHERE vf.type = ? AND u.app_id = ? %s ORDER BY update_time DESC LIMIT ?,?", stateQuery+searchQuery)

	maps, err := conn.Query(countSqlStr, types.VerifyForUser, appId, start, end-start)
	if err != nil || len(maps) < 1 {
		return 0, nil, err
	}
	count := utility.ToInt64(maps[0]["count"])

	maps, err = conn.Query(sqlStr, types.VerifyForUser, appId, start, end-start)
	return count, maps, err
}

func RoomVerifyList(appId string, search *string, start, end int64, state *int) (int64, []map[string]string, error) {
	stateQuery := ""
	if state != nil {
		stateQuery = fmt.Sprintf(" AND state = %v", *state)
	}

	searchQuery := ""
	if search != nil && *search != "" {
		*search = QueryStr(*search)
		searchQuery = fmt.Sprintf(" AND (r.mark_id LIKE '%%%v%%' OR r.name LIKE '%%%v%%')", *search, *search)
	}
	sqlStr := fmt.Sprintf("SELECT vf.*,r.avatar,r.mark_id,r.name,u.phone FROM verify_apply AS vf LEFT JOIN `room` AS r ON vf.target_id = r.id LEFT JOIN `user` AS u ON r.master_id = u.user_id"+
		" WHERE vf.type = ? AND u.app_id = ? %s ORDER BY update_time DESC LIMIT ?,?", stateQuery+searchQuery)

	countSqlStr := fmt.Sprintf("SELECT COUNT(*) AS count FROM verify_apply AS vf LEFT JOIN `room` AS r ON vf.target_id = r.id LEFT JOIN `user` AS u ON r.master_id = u.user_id"+
		" WHERE vf.type = ? AND u.app_id = ? %s ORDER BY update_time DESC LIMIT ?,?", stateQuery+searchQuery)

	maps, err := conn.Query(countSqlStr, types.VerifyForRoom, appId, start, end-start)
	if err != nil || len(maps) < 1 {
		return 0, nil, err
	}
	count := utility.ToInt64(maps[0]["count"])

	maps, err = conn.Query(sqlStr, types.VerifyForRoom, appId, start, end-start)
	return count, maps, err
}

/*func QueryVerifyList(appId string, tp ) ([]map[string]string, error) {
	const sqlStr = "SELECT * FROM verify_apply WHERE app_id = ? ORDER BY `type` ASC"
	return conn.Query(sqlStr, appId)
}*/

//添加认证请求
func AddVerifyApply(val *types.VerifyApply) (int64, error) {
	const sqlStr = "insert into verify_apply(app_id,`type`,target_id,description,amount,currency,state,record_id,update_time) values(?,?,?,?,?,?,?,?,?)"
	_, id, err := conn.Exec(sqlStr, val.AppId, val.Type, val.TargetId, val.Description, val.Amount, val.Currency, val.State, val.RecordId, val.UpdateTime)
	return id, err
}

//
func FindVerifyApplyByState(tp int, targetId string, state int) ([]map[string]string, error) {
	const sqlStr = "select * from verify_apply where `type` = ? and target_id = ? and state = ? order by update_time desc"
	return conn.Query(sqlStr, tp, targetId, state)
}

//
func FindVerifyApplyById(id string) ([]map[string]string, error) {
	const sqlStr = "select * from verify_apply where id = ?"
	return conn.Query(sqlStr, id)
}

//设置状态
func SetVerifyState(tx *mysql.MysqlTx, id string, state int) error {
	const sqlStr = "update verify_apply set state = ? where id = ?"
	_, _, err := conn.Exec(sqlStr, state, id)
	return err
}

//查询手续费总额
func FindFeeSum(appId string) ([]map[string]string, error) {
	const sqlStr = "select currency,sum(amount) as sum_amount from verify_apply where state = ? and fee_state = ? and app_id = ? GROUP BY currency"
	return conn.Query(sqlStr, types.VerifyStateAccept, types.VerifyFeeStateSuccess, appId)
}
