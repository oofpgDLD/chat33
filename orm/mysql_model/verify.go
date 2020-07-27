package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func PersonalVerifyList(appId string, search *string, start, end int64, state *int) (int64, []*types.VerifyApplyJoinUser, error) {
	count, maps, err := db.PersonalVerifyList(appId, search, start, end, state)
	if err != nil {
		return count, nil, err
	}
	list := make([]*types.VerifyApplyJoinUser, 0)
	for _, m := range maps {
		item := &types.VerifyApplyJoinUser{
			VerifyApply: convertVerifyApply(m),
			User:        convertUser(m),
		}
		list = append(list, item)
	}
	return count, list, nil
}

func RoomVerifyList(appId string, search *string, start, end int64, state *int) (int64, []*types.VerifyApplyJoinRoomAndUser, error) {
	count, maps, err := db.RoomVerifyList(appId, search, start, end, state)
	if err != nil {
		return count, nil, err
	}
	list := make([]*types.VerifyApplyJoinRoomAndUser, 0)
	for _, m := range maps {
		item := &types.VerifyApplyJoinRoomAndUser{
			VerifyApply: convertVerifyApply(m),
			Room:        convertRoom(m),
			User:        convertUser(m),
		}
		list = append(list, item)
	}
	return count, list, nil
}

func AddVerifyApply(val *types.VerifyApply) (int64, error) {
	return db.AddVerifyApply(val)
}

func convertVerifyApply(v map[string]string) *types.VerifyApply {
	return &types.VerifyApply{
		Id:          v["id"],
		Type:        utility.ToInt(v["type"]),
		TargetId:    v["target_id"],
		Description: v["description"],
		Amount:      utility.ToFloat64(v["amount"]),
		Currency:    v["currency"],
		State:       utility.ToInt(v["state"]),
		RecordId:    v["record_id"],
		UpdateTime:  utility.ToInt64(v["update_time"]),
		FeeState:    utility.ToInt(v["fee_state"]),
	}
}

func FindVerifyApplyByState(tp int, targetId string, state int) ([]*types.VerifyApply, error) {
	maps, err := db.FindVerifyApplyByState(tp, targetId, state)
	if err != nil {
		return nil, err
	}
	list := make([]*types.VerifyApply, 0)
	for _, m := range maps {
		item := convertVerifyApply(m)
		list = append(list, item)
	}
	return list, nil
}

func FindVerifyApplyById(id string) (*types.VerifyApply, error) {
	maps, err := db.FindVerifyApplyById(id)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertVerifyApply(info), nil
}
