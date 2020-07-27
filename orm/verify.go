package orm

import (
	"github.com/33cn/chat33/db"
	mysql "github.com/33cn/chat33/orm/mysql_model"
	sql "github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logVerify = log15.New("model", "orm/verify")

//---------------------------认证----------------------//

//设置认证手续费
func VerifySetFee(conf []*types.VerifyFee) error {
	tx, err := GetTx()
	if err != nil {
		logVerify.Error("VerifySetFee get tx", "err", err)
		return err
	}
	for _, c := range conf {
		err := db.VerifySetFee(tx.(*sql.MysqlTx), c.AppId, c.Type, c.Currency, c.Amount)
		if err != nil {
			tx.RollBack()
			logVerify.Error("db.VerifySetFee", "err", err, "appId", c.AppId, "type", c.Type, "currency", c.Currency, "amount", c.Amount)
			return err
		}
	}
	return tx.Commit()
}

//获取认证手续费
func VerifyGetFee(appId string) ([]*types.VerifyFee, error) {
	maps, err := db.VerifyGetFee(appId)
	if err != nil {
		logVerify.Error("db.VerifyGetFee", "err", err, "appId", appId)
		return nil, err
	}
	list := make([]*types.VerifyFee, 0)
	for _, m := range maps {
		item := &types.VerifyFee{
			AppId:    appId,
			Type:     utility.ToInt(m["type"]),
			Currency: m["currency"],
			Amount:   utility.ToFloat64(m["amount"]),
		}
		list = append(list, item)
	}
	return list, nil
}

func PersonalVerifyList(appId string, search *string, start, end int64, state *int) (int64, []*types.VerifyApplyJoinUser, error) {
	ret, ret2, err := mysql.PersonalVerifyList(appId, search, start, end, state)
	if err != nil {
		logVerify.Error("mysql.PersonalVerifyList", "err", err, "appId", appId, "search", search, "start", start, "end", end, "state", state)
	}
	return ret, ret2, err
}

func RoomVerifyList(appId string, search *string, start, end int64, state *int) (int64, []*types.VerifyApplyJoinRoomAndUser, error) {
	ret, ret2, err := mysql.RoomVerifyList(appId, search, start, end, state)
	if err != nil {
		logVerify.Error("mysql.RoomVerifyList", "err", err, "appId", appId, "search", search, "start", start, "end", end, "state", state)
	}
	return ret, ret2, err
}

//更新手续费到账状态
func UpdateVerifyFeeState(id string, feeState int) error {
	err := db.UpdateVerifyFeeState(id, feeState)
	if err != nil {
		logVerify.Error("db.UpdateVerifyFeeState", "err", err, "id", id, "feeState", feeState)
	}
	return err
}

//添加认证请求
func AddVerifyApply(val *types.VerifyApply) (int64, error) {
	ret, err := mysql.AddVerifyApply(val)
	if err != nil {
		logVerify.Error("mysql.AddVerifyApply", "err", err)
	}
	return ret, err
}

//根据状态查找认证信息
func FindVerifyApplyByState(tp int, targetId string, state int) ([]*types.VerifyApply, error) {
	ret, err := mysql.FindVerifyApplyByState(tp, targetId, state)
	if err != nil {
		logVerify.Error("mysql.FindVerifyApplyByState", "err", err, "tp", tp, "targetId", targetId, "state", state)
	}
	return ret, err
}

//根据记录id查找认证信息系
func FindVerifyApplyById(id string) (*types.VerifyApply, error) {
	ret, err := mysql.FindVerifyApplyById(id)
	if err != nil {
		logVerify.Error("mysql.FindVerifyApplyById", "err", err, "id", id)
	}
	return ret, err
}

//设置状态
func SetVerifyState(tx types.Tx, id string, state int) error {
	err := db.SetVerifyState(tx.(*sql.MysqlTx), id, state)
	if err != nil {
		logVerify.Error("db.SetVerifyState", "err", err, "id", id, "state", state)
	}
	return err
}

//获取手续费总额
func FindFeeSum(appId string) (map[string]float64, error) {
	maps, err := db.FindFeeSum(appId)
	if err != nil {
		logVerify.Error("db.FindFeeSum", "err", err, "appId", appId)
	}
	list := make(map[string]float64)
	for _, m := range maps {
		list[m["currency"]] = utility.ToFloat64(m["sum_amount"])
	}
	return list, nil
}
