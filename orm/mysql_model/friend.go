package mysql_model

import (
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func convertFriend(v map[string]string) *types.Friend {
	return &types.Friend{
		UserId:    v["F_user_id"],
		FriendId:  v["friend_id"],
		Remark:    v["remark"],
		AddTime:   utility.ToInt64(v["add_time"]),
		DND:       utility.ToInt(v["DND"]),
		Top:       utility.ToInt(v["top"]),
		Type:      utility.ToInt(v["type"]),
		IsDelete:  utility.ToInt(v["is_delete"]),
		Source:    v["source"],
		ExtRemark: v["ext_remark"],
		IsBlocked: utility.ToInt(v["is_blocked"]),
	}
}

func FindFriendById(userID, friendID string) (*types.FriendJoinUser, error) {
	maps, err := db.FindFriend(userID, friendID)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return &types.FriendJoinUser{
		Friend: convertFriend(info),
		User:   convertJoinUser(info),
	}, nil
}

func FindFriendsId(userId string) ([]string, error) {
	maps, err := db.FindFriendIdByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	friends := make([]string, 0)
	for _, v := range maps {
		friends = append(friends, v["friend_id"])
	}
	return friends, nil
}

func FindFriendsById(userId string, commonUse, isDel int) ([]*types.FriendJoinUser, error) {
	maps, err := db.GetFriendList(userId, commonUse, isDel)
	if err != nil {
		return nil, err
	}
	friends := make([]*types.FriendJoinUser, 0)
	for _, v := range maps {
		friend := &types.FriendJoinUser{
			Friend: convertFriend(v),
			User:   convertJoinUser(v),
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func FindFriendsAfterTime(userId string, commonUse int, time int64) ([]*types.FriendJoinUser, error) {
	maps, err := db.FindFriendsAfterTime(userId, commonUse, time)
	if err != nil {
		return nil, err
	}
	friends := make([]*types.FriendJoinUser, 0)
	for _, v := range maps {
		friend := &types.FriendJoinUser{
			Friend: convertFriend(v),
			User:   convertJoinUser(v),
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func FindPrivateChatLogById(id string) (*types.PrivateLogJoinUser, error) {
	maps, err := db.FindPrivateChatLogById(id)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	v := maps[0]
	return &types.PrivateLogJoinUser{
		PrivateLog: &types.PrivateLog{
			Id:        v["id"],
			MsgId:     v["msg_id"],
			ReceiveId: v["receive_id"],
			SenderId:  v["sender_id"],
			IsSnap:    utility.ToInt(v["is_snap"]),
			MsgType:   utility.ToInt(v["msg_type"]),
			Content:   v["content"],
			Status:    utility.ToInt(v["status"]),
			SendTime:  utility.ToInt64(v["send_time"]),
			Ext:       v["ext"],
			IsDelete:  utility.ToInt(v["is_delete"]),
		},
		User: convertUser(v),
	}, nil
}

func FindPrivateChatLogByMsgId(senderId, msgId string) (*types.PrivateLogJoinUser, error) {
	maps, err := db.FindPrivateChatLogByMsgId(senderId, msgId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	v := maps[0]
	return &types.PrivateLogJoinUser{
		PrivateLog: &types.PrivateLog{
			Id:        v["id"],
			MsgId:     v["msg_id"],
			ReceiveId: v["receive_id"],
			SenderId:  v["sender_id"],
			IsSnap:    utility.ToInt(v["is_snap"]),
			MsgType:   utility.ToInt(v["msg_type"]),
			Content:   v["content"],
			Status:    utility.ToInt(v["status"]),
			SendTime:  utility.ToInt64(v["send_time"]),
			Ext:       v["ext"],
			IsDelete:  utility.ToInt(v["is_delete"]),
		},
		User: convertUser(v),
	}, nil
}

func UpdatePrivateLogContentById(logId, content string) error {
	return db.UpdatePrivateLogContentById(logId, content)
}

func FindNotBurnedLogsAfter(userId string, isDel int, time int64) ([]*types.PrivateLogJoinUser, error) {
	maps, err := db.FindNotBurndLogAfter(userId, isDel, time)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, v := range maps {
		log := &types.PrivateLogJoinUser{
			PrivateLog: &types.PrivateLog{
				Id:        v["id"],
				MsgId:     v["msg_id"],
				ReceiveId: v["receive_id"],
				SenderId:  v["sender_id"],
				IsSnap:    utility.ToInt(v["is_snap"]),
				MsgType:   utility.ToInt(v["msg_type"]),
				Content:   v["content"],
				Status:    utility.ToInt(v["status"]),
				SendTime:  utility.ToInt64(v["send_time"]),
				Ext:       v["ext"],
				IsDelete:  utility.ToInt(v["is_delete"]),
			},
			User: convertUser(v),
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindNotBurnedLogsBetween(userId string, isDel int, begin, end int64) ([]*types.PrivateLogJoinUser, error) {
	maps, err := db.FindNotBurndLogBetween(userId, isDel, begin, end)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, v := range maps {
		log := &types.PrivateLogJoinUser{
			PrivateLog: &types.PrivateLog{
				Id:        v["id"],
				MsgId:     v["msg_id"],
				ReceiveId: v["receive_id"],
				SenderId:  v["sender_id"],
				IsSnap:    utility.ToInt(v["is_snap"]),
				MsgType:   utility.ToInt(v["msg_type"]),
				Content:   v["content"],
				Status:    utility.ToInt(v["status"]),
				SendTime:  utility.ToInt64(v["send_time"]),
				Ext:       v["ext"],
				IsDelete:  utility.ToInt(v["is_delete"]),
			},
			User: convertUser(v),
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindPrivateChatLogsNumberBetween(userId string, isDel int, begin, end int64) (int64, error) {
	maps, err := db.FindChatLogsNumberBetween(userId, isDel, begin, end)
	if err != nil {
		return 0, err
	}
	if len(maps) < 1 {
		return 0, nil
	}
	return utility.ToInt64(maps[0]["count"]), nil
}

func FindAllPrivateLogs() ([]*types.PrivateLog, error) {
	maps, err := db.FindAllPrivateLogs()
	if err != nil {
		return nil, err
	}
	logs := make([]*types.PrivateLog, 0)
	for _, v := range maps {
		log := &types.PrivateLog{
			Id:        v["id"],
			MsgId:     v["msg_id"],
			ReceiveId: v["receive_id"],
			SenderId:  v["sender_id"],
			IsSnap:    utility.ToInt(v["is_snap"]),
			MsgType:   utility.ToInt(v["msg_type"]),
			Content:   v["content"],
			Status:    utility.ToInt(v["status"]),
			SendTime:  utility.ToInt64(v["send_time"]),
			Ext:       v["ext"],
			IsDelete:  utility.ToInt(v["is_delete"]),
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindPrivateChatLogs(userId, friendId string, start int64, number int) (int64, []*types.PrivateLogJoinUser, error) {
	//maps, nextLog, err := db.FindCatLog(userId, friendId, start, number)
	maps, nextLog, err := db.FindChatLogV2(userId, friendId, start, number)
	if err != nil {
		return -1, nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, v := range maps {
		log := &types.PrivateLogJoinUser{
			PrivateLog: &types.PrivateLog{
				Id:        v["id"],
				MsgId:     v["msg_id"],
				ReceiveId: v["receive_id"],
				SenderId:  v["sender_id"],
				IsSnap:    utility.ToInt(v["is_snap"]),
				MsgType:   utility.ToInt(v["msg_type"]),
				Content:   v["content"],
				Status:    utility.ToInt(v["status"]),
				SendTime:  utility.ToInt64(v["send_time"]),
				Ext:       v["ext"],
				IsDelete:  utility.ToInt(v["is_delete"]),
			},
			User: convertUser(v),
		}
		logs = append(logs, log)
	}
	return utility.ToInt64(nextLog), logs, nil
}

func FindSessionKeyAlert(userId string, endTime int64) ([]*types.PrivateLogJoinUser, error) {
	maps, err := db.FindSessionKeyAlert(userId, endTime)
	if err != nil {
		return nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, v := range maps {
		log := &types.PrivateLogJoinUser{
			PrivateLog: &types.PrivateLog{
				Id:        v["id"],
				MsgId:     v["msg_id"],
				ReceiveId: v["receive_id"],
				SenderId:  v["sender_id"],
				IsSnap:    utility.ToInt(v["is_snap"]),
				MsgType:   utility.ToInt(v["msg_type"]),
				Content:   v["content"],
				Status:    utility.ToInt(v["status"]),
				SendTime:  utility.ToInt64(v["send_time"]),
				Ext:       v["ext"],
				IsDelete:  utility.ToInt(v["is_delete"]),
			},
			User: convertUser(v),
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func FindTypicalPrivateChatLogs(userId, friendId, owner string, startId int64, number int, queryType []string) (int64, []*types.PrivateLogJoinUser, error) {
	maps, nextLog, err := db.FindTypicalChatLogs(userId, friendId, owner, startId, number, queryType)
	if err != nil {
		return -1, nil, err
	}
	logs := make([]*types.PrivateLogJoinUser, 0)
	for _, v := range maps {
		log := &types.PrivateLogJoinUser{
			PrivateLog: &types.PrivateLog{
				Id:        v["id"],
				MsgId:     v["msg_id"],
				ReceiveId: v["receive_id"],
				SenderId:  v["sender_id"],
				IsSnap:    utility.ToInt(v["is_snap"]),
				MsgType:   utility.ToInt(v["msg_type"]),
				Content:   v["content"],
				Status:    utility.ToInt(v["status"]),
				SendTime:  utility.ToInt64(v["send_time"]),
				Ext:       v["ext"],
				IsDelete:  utility.ToInt(v["is_delete"]),
			},
			User: convertUser(v),
		}
		logs = append(logs, log)
	}
	return utility.ToInt64(nextLog), logs, nil
}

func FindFirstPrivateMsg(userId, friendId string) (*types.PrivateLog, error) {
	maps, err := db.FindFirstMsg(userId, friendId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	v := maps[0]
	return &types.PrivateLog{
		Id:        v["id"],
		MsgId:     v["msg_id"],
		ReceiveId: v["receive_id"],
		SenderId:  v["sender_id"],
		IsSnap:    utility.ToInt(v["is_snap"]),
		MsgType:   utility.ToInt(v["msg_type"]),
		Content:   v["content"],
		Status:    utility.ToInt(v["status"]),
		SendTime:  utility.ToInt64(v["send_time"]),
		Ext:       v["ext"],
		IsDelete:  utility.ToInt(v["is_delete"]),
	}, nil
}

func FindAddFriendConfByUserId(userId string) (*types.AddFriendConf, error) {
	maps, err := db.FindAddFriendConfByUserId(userId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	v := maps[0]
	return &types.AddFriendConf{
		Id:          v["id"],
		UserId:      v["user_id"],
		NeedConfirm: v["need_confirm"],
		NeedAnswer:  v["need_answer"],
		Question:    v["question"],
		Answer:      v["answer"],
	}, nil
}

func FindLastCatLogId(userId, friendId string) (int64, error) {
	res, err := db.FindLastCatLogId(userId, friendId)
	if err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, nil
	}
	return utility.ToInt64(res[0]["MAX(`id`)"]), nil
}

func DelPrivateChatLog(logId string) (int, error) {
	return db.DeleteCatLog(logId)
}

func UpdatePrivateLogStateById(revId string, state int) error {
	_, _, err := db.UpdatePrivateLogStateById(revId, state)
	return err
}

func CheckIsFriend(userId, friendId string, isDel int) (bool, error) {
	return db.CheckFriend(userId, friendId, isDel)
}

func FindFriendApplyCount(userId, friendId string) (int32, error) {
	return db.FindApplyCount(userId, friendId)
}

func InsertFriend(userID, friendID, remark, extRemark string, dnd, top int, addTime int64) error {
	return db.InsertFriend(userID, friendID, remark, extRemark, dnd, top, addTime)
}

func SetFriendRemark(userId, friendId, remark string) error {
	return db.SetFriendRemark(userId, friendId, remark)
}

func SetFriendExtRemark(userId, friendId, remark string) error {
	return db.SetFriendExtRemark(userId, friendId, remark)
}

func SetFriendDND(userId, friendId string, DND int) error {
	return db.SetFriendDND(userId, friendId, DND)
}

func SetFriendIsTop(userId, friendId string, isTop int) error {
	return db.SetFriendTop(userId, friendId, isTop)
}

func SetQuestionandAnswer(userId, question, answer string) error {
	return db.SetQuestionandAnswer(userId, question, answer)
}

func SetNeedAnswer(userId, question, answer string) error {
	return db.NeedAnswer(userId, question, answer)
}

func SetNotNeedAnswer(userId string) error {
	return db.NotNeedAnswer(userId)
}

func IsNeedConfirm(userId string, state int) error {
	return db.IsNeedConfirm(userId, state)
}

func DeleteFriend(userId, friendId string, alterTime int64) error {
	return db.DeleteFriend(userId, friendId, alterTime)
}

func GetUnReadNumber(userId, friendId string) (int32, error) {
	return db.FindUnReadNum(userId, friendId)
}

func AcceptFriend(userId, friendId, source string, addTime int64) error {
	return db.AcceptFriend(userId, friendId, source, addTime)
}

func AddPrivateChatLog(senderId, receiveId, msgId string, msgType, status, isSnap int, content, ext string, time int64) (int64, int64, error) {
	return db.AddPrivateChatLog(senderId, receiveId, msgId, msgType, status, isSnap, content, ext, time)
}

func SetFriendIsBlock(userId, friendId string, state int, alterTime int64) error {
	return db.SetFriendIsBlock(userId, friendId, state, alterTime)
}

func FindBlockedFriends(userId string) ([]*types.FriendJoinUser, error) {
	maps, err := db.FindBlockedList(userId)
	if err != nil {
		return nil, err
	}
	friends := make([]*types.FriendJoinUser, 0)
	for _, v := range maps {
		friend := &types.FriendJoinUser{
			Friend: convertFriend(v),
			User:   convertJoinUser(v),
		}
		friends = append(friends, friend)
	}
	return friends, nil
}
