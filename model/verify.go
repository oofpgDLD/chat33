package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/pkg/account"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logVerify = log15.New("model", "model/admin")

//---------------------认证-------------------------//
//退回手续费 id:请求记录id ， recordId 托管打币记录id
func feeBack(appId, id, recordId string) error {
	//调用打币接口
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	feeState := types.VerifyFeeStateBacking
	switch app.IsInner {
	case types.IsInnerAccount:
		if recordId != "" {
			_, err := account.BcoinReturnCost(app.AccountServer, recordId, types.VerifyRecordRefund)
			if err != nil {
				return err
			}
		} else {
			feeState = types.VerifyFeeStateBackSuccess
		}
	default:
		return types.ERR_NOTSUPPORT
	}
	//更改手续费到账状态:等待到账
	err := orm.UpdateVerifyFeeState(id, feeState)
	if err != nil {
		return errors.New("数据库访问错误")
	}
	return nil
}

func VerifyApprove(appId, id string, accept int) error {
	//查找对应记录
	apply, err := orm.FindVerifyApplyById(id)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if apply != nil {
		if apply.State != types.AwaitState {
			//退回手续费
			if apply.State == types.RejectState && accept == types.VerifyApproveFeeBack {
				err := feeBack(appId, id, apply.RecordId)
				if err != nil {
					return result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error())
				}
			}
			//返回成功
			return nil
		}

		tx, err := orm.GetTx()
		if err != nil {
			logVerify.Error("VerifyApprove tx GetTx", "err", err.Error())
			return result.NewError(result.DbConnectFail)
		}

		if accept == types.VerifyApproveAccept {
			switch apply.Type {
			case types.VerifyForUser:
				err := orm.SetUserVerifyed(tx, apply.TargetId, apply.Description)
				if err != nil {
					tx.RollBack()
					return result.NewError(result.DbConnectFail)
				}
			case types.VerifyForRoom:
				err := orm.SetRoomVerifyed(tx, apply.TargetId, apply.Description)
				if err != nil {
					tx.RollBack()
					return result.NewError(result.DbConnectFail)
				}
			default:
				tx.RollBack()
				logVerify.Error("apply type from not match", "type", apply.Type)
				return result.NewError(result.DbConnectFail).SetChildErr(result.ServiceChat, nil, fmt.Sprintf("apply tpye %d not match", apply.Type))
			}
		} else if accept == types.VerifyApproveReject {
			//退回手续费
			err := feeBack(appId, id, apply.RecordId)
			if err != nil {
				tx.RollBack()
				return result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error())
			}
		}

		err = orm.SetVerifyState(tx, id, types.VerifyApproveToState[accept])
		if err != nil {
			tx.RollBack()
			return result.NewError(result.DbConnectFail)
		}
		err = tx.Commit()
		if err != nil {
			logVerify.Error("VerifyApprove tx Commit", "err", err.Error())
			return result.NewError(result.DbConnectFail)
		}
	}
	return nil
}

func PersonalVerifyList(appId string, search *string, page, size int, state *int) (int64, []*types.VerifyApplyJoinUser, error) {
	start := int64((page - 1) * size)
	end := int64(page * size)

	app := app.GetApp(appId)
	if app == nil {
		return 0, nil, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	count, apply, err := orm.PersonalVerifyList(appId, search, start, end, state)
	if err != nil {
		return 0, nil, result.NewError(result.DbConnectFail)
	}
	for _, v := range apply {
		if v.FeeState == types.VerifyFeeStateCosting || v.FeeState == types.VerifyFeeStateBacking {
			//查询记录 检查
			switch app.IsInner {
			case types.IsInnerAccount:
				feeState, err := account.BcoinCheckTrans(app.AccountServer, v.RecordId, types.VerifyFeeRecord[v.FeeState])
				if err != nil {
					return 0, nil, result.NewError(result.VisitAccountSystemFailed).SetExtMessage(err.Error())
				}

				//更新手续费入账状态
				err = orm.UpdateVerifyFeeState(v.VerifyApply.Id, feeState)
				if err != nil {
					return 0, nil, result.NewError(result.DbConnectFail)
				}
				v.FeeState = feeState
			default:
				return 0, nil, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(types.ERR_NOTSUPPORT.Error())
			}
		}
	}

	return count, apply, err
}

func RoomVerifyList(appId string, search *string, page, size int, state *int) (int64, []*types.VerifyApplyJoinRoomAndUser, error) {
	start := int64((page - 1) * size)
	end := int64(page * size)

	app := app.GetApp(appId)
	if app == nil {
		return 0, nil, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	count, apply, err := orm.RoomVerifyList(appId, search, start, end, state)
	if err != nil {
		return 0, nil, result.NewError(result.DbConnectFail)
	}
	for _, v := range apply {
		if v.FeeState == types.VerifyFeeStateCosting || v.FeeState == types.VerifyFeeStateBacking {
			//查询记录 检查
			switch app.IsInner {
			case types.IsInnerAccount:
				feeState, err := account.BcoinCheckTrans(app.AccountServer, v.RecordId, types.VerifyFeeRecord[v.FeeState])
				if err != nil {
					return 0, nil, result.NewError(result.VisitAccountSystemFailed).SetExtMessage(err.Error())
				}
				//更新手续费入账状态
				err = orm.UpdateVerifyFeeState(v.VerifyApply.Id, feeState)
				if err != nil {
					return 0, nil, result.NewError(result.DbConnectFail)
				}
				v.FeeState = feeState
			default:
				return 0, nil, result.NewError(result.DbConnectFail).SetExtMessage(types.ERR_NOTSUPPORT.Error())
			}
		}
	}

	return count, apply, err
}

func VerifyGetConfig(appId string) ([]*types.VerifyFee, error) {
	fee, err := orm.VerifyGetFee(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	return fee, nil
}

func VerifySetFee(fee []*types.VerifyFee) error {
	err := orm.VerifySetFee(fee)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func VerifyFeeStatistics(appId string) (interface{}, error) {
	app := app.GetApp(appId)
	if app == nil {
		return nil, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error())
	}

	//查询认证手续费
	sum, err := orm.FindFeeSum(appId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	currency := app.MainCoin
	amount := "0"
	for k, v := range sum {
		if app.MainCoin == k {
			amount = utility.ToString(v)
			break
		}
	}

	return map[string]interface{}{
		"currency": currency,
		"amount":   amount,
	}, nil
}

//检查是否可发起认证请求
func checkVerifyApply(tp int, targetId string) (bool, error) {
	//是否存在 待审核或是已认证的请求
	accepted, err := orm.FindVerifyApplyByState(tp, targetId, types.VerifyStateAccept)
	if err != nil {
		return false, result.NewError(result.DbConnectFail)
	}
	if len(accepted) > 0 {
		return false, nil
	}

	waiting, err := orm.FindVerifyApplyByState(tp, targetId, types.VerifyStateWait)
	if err != nil {
		return false, result.NewError(result.DbConnectFail)
	}
	if len(waiting) > 0 {
		return false, nil
	}

	//是否在 上次被拒绝 1周时间内 发起请求
	rejected, err := orm.FindVerifyApplyByState(tp, targetId, types.VerifyStateReject)
	if err != nil {
		return false, result.NewError(result.DbConnectFail)
	}

	if len(rejected) > 0 {
		rej := rejected[0]
		if utility.NowMillionSecond() < utility.MillionSecondAddDuration(rej.UpdateTime, types.VerifyApplyInterval*time.Minute) {
			//拒绝请求
			return false, result.NewError(result.VerifyLimit)
		}
	}
	return true, nil
}

//认证申请
func VerifyApply(appId, token string, tp int, targetId, description, payPassword, code string) (int64, error) {
	//检查
	ok, err := checkVerifyApply(tp, targetId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, nil
	}

	//开始认证
	//查询手续费
	currency := ""
	var amount float64 = 0
	fee, err := orm.VerifyGetFee(appId)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}

	for _, v := range fee {
		if v.Type == tp {
			currency = v.Currency
			amount = v.Amount
			break
		}
	}

	//调用打币接口
	app := app.GetApp(appId)
	if app == nil {
		return 0, types.ERR_APPNOTFIND
	}

	recordId := ""
	feeState := types.VerifyFeeStateCosting
	switch app.IsInner {
	case types.IsInnerAccount:
		if utility.ToFloat64(amount) != 0 {
			recordId, err = account.BcoinTakeCost(app.AccountServer, token, currency, utility.ToString(amount), payPassword, code, utility.RandomID(), "vip_auth")
			if err != nil {
				return 0, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error())
			}
		} else {
			feeState = types.VerifyFeeStateSuccess
		}
	default:
		return 0, types.ERR_NOTSUPPORT
	}

	//添加请求信息
	info := &types.VerifyApply{
		Type:        tp,
		AppId:       appId,
		TargetId:    targetId,
		Description: description,
		Amount:      amount,
		Currency:    currency,
		State:       types.VerifyStateWait,
		RecordId:    recordId,
		FeeState:    feeState,
		UpdateTime:  utility.NowMillionSecond(),
	}
	id, err := orm.AddVerifyApply(info)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	return id, nil
}
