package model

import (
	"encoding/json"
	"strconv"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/proto"

	"github.com/33cn/chat33/router"

	"time"

	"strings"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logFriend = log15.New("module", "model/friend")

func FriendRemarkOrName(userId, friendId string) (string, error) {
	f, err := orm.FindFriendById(userId, friendId)
	if err != nil {
		return "", err
	}
	if f.Remark != "" {
		return f.Remark, nil
	}
	return f.Username, nil
}

//判断用户是否存在  是否是好友  对好友的一些操作都需要先判断
func boolIsExistIsFriend(userID, friendID string) error {
	bool := CheckUserExist(friendID)
	if !bool {
		return result.NewError(result.UserNotExists)
	}
	//不需要判断是不是好友
	//bool, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	//if err != nil {
	//	return result.NewError(result.DbConnectFail)
	//}

	//if !bool {
	//	logFriend.Warn("is not friend", "userId", userID, "friendId", friendID)
	//	return result.NewError(result.IsNotFriend)
	//}
	return nil
}

/*
	好友列表
*/
func FriendList(userID string, tp int, time int64) (interface{}, error) {
	if time == 0 {
		friends, err := orm.FindFriendsById(userID, tp, types.FriendIsNotDelete)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		data := make([]types.FriendInfoApi, 0)

		for _, friend := range friends {
			//全量拉取好友列表时 不包含删除的好友和黑名单好友
			if friend.IsBlocked == types.IsBlocked {
				continue
			}
			one := types.FriendInfoApi{
				Uid:                friend.Uid,
				Id:                 friend.FriendId,
				Name:               friend.Username,
				Avatar:             friend.Avatar,
				Remark:             friend.Remark,
				PublicKey:          friend.PublicKey,
				NoDisturbing:       friend.DND,
				CommonlyUsed:       friend.Type,
				OnTop:              friend.Top,
				IsDelete:           friend.IsDelete,
				AddTime:            friend.AddTime,
				DepositAddress:     friend.DepositAddress,
				Identification:     friend.Identification,
				IdentificationInfo: friend.IdentificationInfo,
				IsBlocked:          friend.IsBlocked,
			}
			data = append(data, one)
		}
		return map[string]interface{}{"userList": data}, nil
	} else {
		friends, err := orm.FindFriendsAfterTime(userID, tp, time)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		data := make([]types.FriendInfoApi, 0)

		for _, friend := range friends {
			one := types.FriendInfoApi{
				Uid:                friend.Uid,
				Id:                 friend.FriendId,
				Name:               friend.Username,
				Avatar:             friend.Avatar,
				Remark:             friend.Remark,
				PublicKey:          friend.PublicKey,
				NoDisturbing:       friend.DND,
				CommonlyUsed:       friend.Type,
				OnTop:              friend.Top,
				IsDelete:           friend.IsDelete,
				AddTime:            friend.AddTime,
				DepositAddress:     friend.DepositAddress,
				Identification:     friend.Identification,
				IdentificationInfo: friend.IdentificationInfo,
				IsBlocked:          friend.IsBlocked,
			}
			data = append(data, one)
		}
		return map[string]interface{}{"userList": data}, nil
	}
}

//添加好友
//sourceType  1  【通过搜索添加】  2  【通过扫一扫添加】  3 通过好友分享(uid)  4通过群
func AddFriend(appId, userID, friendID, remark, reason, sourceId string, sourceType int, answer string) (int, error) {
	//不能对自己操作
	if userID == friendID {
		logFriend.Debug("can not add friend", "warn", "CanNotOperateSelf", "userId", userID, "friendID", friendID)
		return 0, result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在
	isExists := CheckUserExist(friendID)
	if !isExists {
		logFriend.Debug("AddFriend UserNotExists", "friendId", friendID)
		return 0, result.NewError(result.UserNotExists)
	}

	//验证是否是好友
	//isFriend, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	//if isFriend {
	//	if err != nil {
	//		return 0, result.NewError(result.DbConnectFail)
	//	}
	//	return 0, result.NewError(result.IsFriendAlready)
	//}

	//判断对方是否把你加入黑名单
	if ok, err := CheckIsBlocked(friendID, userID); ok || err != nil {
		logFriend.Debug("you are in the blocked", "you", userID, "friend", friendID)
		return 0, result.NewError(result.IsBlocked)
	}

	if sourceType == types.Share {
		user, err := orm.GetUserInfoByMarkId(appId, sourceId)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		sourceId = user.UserId
	}

	isFromRoom := true
	//不是通过群添加好友
	if sourceType != types.Group {
		isFromRoom = false
	}

	if isFromRoom {
		room, err := orm.FindRoomById(sourceId, types.RoomDeletedOrNot)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}

		//判断群是否存在
		if room == nil || room.IsDelete == types.IsDelete {
			return 0, result.NewError(result.RoomNotExists)
		}
		//判断该群是否允许添加好友
		if room.CanAddFriend == types.CanNotAddFriend {
			return 0, result.NewError(result.CanNotAddFriendInRoom)
		}

		//判断双方是否都在群里
		userInRoom1, err := CheckUserInRoom(userID, sourceId)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		userInRoom2, err := CheckUserInRoom(friendID, sourceId)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if !userInRoom1 || !userInRoom2 {
			return 0, result.NewError(result.UserIsNotInRoom)
		}
	}

	/*//查询cid
	uFriend, err := orm.GetUserInfoById(friendID)
	if err != nil {
		logFriend.Warn("FindCid db dailed", "err_msg", err)
		return 0, result.NewError(result.DbConnectFail)
	}*/

	senderInfo, receiveInfo, err := findSendReceiveInfo(userID, friendID)
	if err != nil {
		return 0, err
	}

	//text := senderInfo["name"].(string) + "请求添加您为好友"

	//判断对方加好友设置
	conf, err := GetAddFriendConfByUserId(friendID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}

	//需要回答问题
	if conf != nil && utility.ToInt(conf.NeedAnswer) == types.NeedAnswer {
		if answer != conf.Answer {
			return types.AnswerFalse, nil
		}
	}

	//判断请求是否存在
	cou, err := orm.GetFriendApplyCount(userID, friendID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}

	var logId string
	//如果请求存在，只执行一些更新操作
	if cou > 0 {
		id, err := orm.UpdateApply(userID, friendID, types.IsFriend, reason, remark, converSource(sourceType, sourceId), types.AwaitState)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		logId = utility.ToString(id)
	} else {
		//添加
		id, err := orm.AppendApplyLog(friendID, userID, reason, converSource(sourceType, sourceId), remark, types.AwaitState, types.IsFriend)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		logId = utility.ToString(id)
	}
	//获取来源
	//source, err := GetFriendSource(friendID, userID)
	source, err := ConverFriendSource(converSource(sourceType, sourceId), friendID, userID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}

	source2, err := ConverFriendSourceV2(converSource(sourceType, sourceId), userID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}

	//发送事件通知
	proto.SendAddFriendNotification(friendID, logId, reason, source, source2, types.AwaitState, senderInfo, receiveInfo)

	//不需要验证
	if conf != nil && utility.ToInt(conf.NeedConfirm) == types.NotNeedConfirm {
		err = HandleFriendRequest(friendID, userID, types.AcceptState)
		if err != nil {
			logFriend.Warn("HandleFriendRequest error", "err", err)
		}
		return types.AddSuccess, err
	}
	return types.SendSuccess, nil
}

type Source struct {
	SourceType int    `json:"sourceType"`
	SourceId   string `json:"sourceId"`
}

func converSource(sourceType int, sourceId string) string {
	var source = &Source{SourceType: sourceType, SourceId: sourceId}
	return utility.StructToString(source)
}

//处理好友请求
func HandleFriendRequest(userID, friendID string, agree int) error {
	var err error
	//先确保该请求是否存在
	applyInfo, err := orm.FindApplyLogByUserAndTarget(friendID, userID, types.IsFriend)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if applyInfo == nil {
		logFriend.Debug("apply not find", "userId", userID, "friendId", friendID, "agree", agree)
		return result.NewError(result.NotExistFriendRequest)
	}

	//确保请求状态为未处理
	if applyInfo.State != types.AwaitState {
		logFriend.Debug("applyInfo state is not AwaitState", "state", applyInfo.State)
		return result.NewError(result.FriendRequestHadDeal)
	}
	//验证是否是好友,如果已经是好友则直接更改状态 不更新时间
	isFriend, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if isFriend {
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		_, err = orm.AcceptFriendApply(friendID, userID)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		return result.NewError(result.IsFriendAlready)
	}

	receiveInfo, senderInfo, err := findSendReceiveInfo(userID, friendID)
	if err != nil {
		return err
	}

	source, err := ConverFriendSource(applyInfo.Source, friendID, applyInfo.ApplyUser)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	source2, err := ConverFriendSourceV2(applyInfo.Source, applyInfo.ApplyUser)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	if agree == types.RejectRequest {
		//拒绝好友请求
		_, err = orm.UpdateApply(friendID, userID, types.IsFriend, applyInfo.ApplyReason, applyInfo.Remark, applyInfo.Source, types.RejectState)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		//text := receiveInfo["name"].(string) + "拒绝了您的好友请求"
		proto.SendAddFriendNotification(friendID, applyInfo.Id, applyInfo.ApplyReason, source, source2, types.RejectState, senderInfo, receiveInfo)
		//多端同步
		//查找来源
		source, err := ConverFriendSource(applyInfo.Source, userID, applyInfo.ApplyUser)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		proto.SendAddFriendNotification(userID, applyInfo.Id, applyInfo.ApplyReason, source, source2, types.RejectState, senderInfo, receiveInfo)
	} else {
		now := utility.NowMillionSecond()
		//同意添加好友
		err = orm.AcceptFriend(userID, friendID, applyInfo.Source, now)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		//text := receiveInfo["name"].(string) + "已同意您的好友请求"
		proto.SendAddFriendNotification(friendID, applyInfo.Id, applyInfo.ApplyReason, source, source2, types.AcceptState, senderInfo, receiveInfo)

		//多端同步
		//查找来源
		source, err := ConverFriendSource(applyInfo.Source, userID, applyInfo.ApplyUser)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		proto.SendAddFriendNotification(userID, applyInfo.Id, applyInfo.ApplyReason, source, source2, types.AcceptState, senderInfo, receiveInfo)
		//如果是群里添加的好友，则需要向群成员发送websocket
		sourceType, roomId := unConverSource(applyInfo.Source)
		if sourceType == types.Group {
			//fix:判断群是否存在或者删除 dld-v2.8.4 2019年7月2日16:10:25
			room, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
			if room != nil && err == nil {
				//格式：XX1添加XX1为好友
				//查找群成员
				managers, err := orm.GetRoomManagerAndMaster(roomId)
				if err != nil {
					return result.NewError(result.DbConnectFail)
				}
				var members = make([]string, 0)
				for _, v := range managers {
					userId := v.RoomMember.UserId
					members = append(members, userId)
				}
				//封装消息
				//查找该显示的名称
				receiveName := orm.GetMemberName(roomId, userID)
				senderName := orm.GetMemberName(roomId, friendID)

				msgContext := senderName + " 添加 " + receiveName + " 为好友"
				SendAlert(userID, roomId, types.ToRoom, members, types.Alert, proto.ComposeAddFriendAlert(types.AlertAddFriendInRoom, friendID, senderName, userID, receiveName, msgContext))
			}
		}
		//像双方发送成为好友的通知
		go sendBecomeFriendMsg(userID, friendID)
	}
	return nil
}

func sendBecomeFriendMsg(userID, friendID string) {
	//var member1 = []string{userID}
	var member2 = []string{userID, friendID}
	SendAlert(userID, friendID, types.ToUser, member2, types.Alert, proto.ComposeAddFriendAlert(types.AlertAddFriend, userID, "", friendID, "", "你们已经成为好友了"))
}

//修改备注
func SetFriendRemark(userID, friendID, remark string) error {
	//不能对自己操作
	if userID == friendID {
		logFriend.Warn("can not set friend remark", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID, "remark", remark)
		return result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在
	err := boolIsExistIsFriend(userID, friendID)
	if err != nil {
		return err
	}
	//判断对应关系是否存在
	bool, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	//对应关系存在编辑，不存在就新增
	if bool == true {
		err = orm.SetFriendRemark(userID, friendID, remark)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	} else {
		now := utility.NowMillionSecond()
		err = orm.InsertFriend(userID, friendID, remark, "", types.NoDisturbingOff, types.NotOnTop, now)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	}

	return nil
}

//修改详细备注
func SetFriendExtRemark(userID, friendID, remark string) error {
	//不能对自己操作
	if userID == friendID {
		logFriend.Warn("can not set friend ext-remark", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID, "remark", remark)
		return result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在  是否是好友
	err := boolIsExistIsFriend(userID, friendID)
	if err != nil {
		return err
	}
	//判断对应关系是否存在
	bool, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	//对应关系存在编辑，不存在就新增
	if bool == true {
		err = orm.SetFriendExtRemark(userID, friendID, remark)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	} else {
		now := utility.NowMillionSecond()
		err = orm.InsertFriend(userID, friendID, "", remark, types.NoDisturbingOff, types.NotOnTop, now)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	}
	return nil
}

/*
	设置好友免打扰
*/
func SetFriendDND(userID, friendID string, DND int) error {
	//不能对自己操作
	if userID == friendID {
		logFriend.Warn("can not set friend DND", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID, "DND", DND)
		return result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}
	//判断对应关系是否存在
	bool, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	//对应关系存在编辑，不存在就新增
	if bool == true {
		err := orm.SetFriendDND(userID, friendID, DND)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	} else {
		now := utility.NowMillionSecond()
		err = orm.InsertFriend(userID, friendID, "", "", DND, types.NotOnTop, now)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	}
	return nil
}

/*
	设置好友置顶
*/
func SetFriendTop(userID, friendID string, top int) error {
	//不能对自己操作
	if userID == friendID {
		logFriend.Warn("can not set friend top", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID, "top", top)
		return result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}

	//判断对应关系是否存在
	bool, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	//对应关系存在编辑，不存在就新增
	if bool == true {
		err := orm.SetFriendIsTop(userID, friendID, top)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	} else {
		now := utility.NowMillionSecond()
		err = orm.InsertFriend(userID, friendID, "", "", types.NoDisturbingOff, top, now)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	}
	return nil
}

/*
	删除好友
*/
func DeleteFriend(userID, friendID string) error {
	//不能对自己操作
	if userID == friendID {
		logFriend.Warn("can not delete friend", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID)
		return result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在  是否是好友
	results := boolIsExistIsFriend(userID, friendID)
	if results != nil {
		return results
	}

	alterTime := utility.NowMillionSecond()
	err := orm.DeleteFriend(userID, friendID, alterTime)
	if err != nil {
		logFriend.Warn("Delete Friend failed", "err_msg", err)
		return result.NewError(result.DbConnectFail)
	}
	msgAlert := "对方已经和你解除好友关系"
	SendAlert(userID, friendID, types.ToUser, []string{friendID, userID}, types.Alert, proto.ComposeDelFriendAlert(types.AlertDeleteFriend, userID, friendID, msgAlert))

	//封装websock消息
	var msgMap = make(map[string]interface{})
	msgMap["eventType"] = 32

	senderInfo, receiveInfo, err := findSendReceiveInfo(userID, friendID)
	if err != nil {
		return err
	}
	msgMap["senderInfo"] = senderInfo
	msgMap["receiveInfo"] = receiveInfo
	msgMap["datetime"] = utility.NowMillionSecond()
	source, err := GetFriendSource(friendID, userID)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	source2, err := GetFriendSourceV2(friendID, userID)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	msgMap["source"] = source
	msgMap["source2"] = source2
	msg, err := json.Marshal(msgMap)
	if err != nil {
		return result.NewError(result.ConvFail)
	}
	client, _ := router.GetUser(userID)
	if client != nil {
		client.SendToAllClients(msg)
	}
	client2, _ := router.GetUser(friendID)
	if client2 != nil {
		client2.SendToAllClients(msg)
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func UserInfo(friendID string) (map[string]interface{}, error) {
	//用户是否存在
	bool := CheckUserExist(friendID)
	if !bool {
		logFriend.Warn("UserInfo", "warn", "UserNotExists", "userId", friendID)
		return nil, result.NewError(result.UserNotExists)
	}

	//好友的信息
	friend, err := orm.GetUserInfoById(friendID)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var returnInfo = make(map[string]interface{})
	//所有人都可以查看的信息
	returnInfo["sex"] = friend.Sex
	returnInfo["avatar"] = friend.Avatar
	returnInfo["id"] = friend.UserId
	returnInfo["uid"] = friend.Uid
	returnInfo["name"] = friend.Username
	returnInfo["mark_id"] = friend.MarkId
	returnInfo["depositAddress"] = friend.DepositAddress
	returnInfo["disableDeadline"] = friend.CloseUntil
	returnInfo["identification"] = friend.Identification
	returnInfo["identificationInfo"] = friend.IdentificationInfo
	//不是好友就返回2
	returnInfo["isFriend"] = 2

	//查找加好友限制
	conf, err := GetAddFriendConfByUserId(friendID)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if conf != nil {
		returnInfo["needConfirm"] = utility.ToInt(conf.NeedConfirm)
		returnInfo["needAnswer"] = utility.ToInt(conf.NeedAnswer)
		returnInfo["question"] = conf.Question
	} else {
		//默认配置
		returnInfo["needConfirm"] = types.NeedConfirm
		returnInfo["needAnswer"] = types.NotNeedAnswer
		returnInfo["question"] = ""
	}

	return returnInfo, nil
}

func ApplyFriendInfo(userId, friendId string) (*types.ApplyInfo, error) {
	var targetInfo types.ApplyInfo
	userInfo, err := FriendInfo(userId, friendId)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		logFriend.Warn("ApplyFriendInfo", "warn", "UserNotExists", "userId", userId, "friendId", friendId)
		return &targetInfo, result.NewError(result.UserNotExists)
	}
	userInfoMap := userInfo.(map[string]interface{})
	targetInfo.Id = userInfoMap["id"].(string)
	targetInfo.MarkId = ""
	targetInfo.Name = userInfoMap["name"].(string)
	targetInfo.Avatar = userInfoMap["avatar"].(string)
	targetInfo.Identification = userInfoMap["identification"].(int)
	targetInfo.IdentificationInfo = userInfoMap["identificationInfo"].(string)
	if userInfoMap["position"] != nil {
		targetInfo.Position = userInfoMap["position"].(string)
	}
	return &targetInfo, nil
}

func getFriendExtRemark(src string) (*types.ExtRemark, error) {
	var ret types.ExtRemark
	err := json.Unmarshal([]byte(src), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

//查看好友详情
func FriendInfo(userID, friendID string) (interface{}, error) {
	//用户是否存在
	b := CheckUserExist(friendID)
	if !b {
		logFriend.Warn("FriendInfo", "warn", "UserNotExists", "userId", friendID)
		return nil, result.NewError(result.UserNotExists)
	}

	//好友的信息
	friend, err := orm.GetUserInfoById(friendID)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	//用户是否存在
	b = CheckUserExist(userID)
	if !b {
		logFriend.Warn("FriendInfo", "warn", "UserNotExists", "userId", userID)
		return nil, result.NewError(result.UserNotExists)
	}
	//自己的信息
	user, err := orm.GetUserInfoById(userID)
	if err != nil {
		logFriend.Info("FriendInfo query db failed", "err_msg", err)
		return nil, result.NewError(result.DbConnectFail)
	}

	var returnInfo = make(map[string]interface{})

	//所有人都可以查看的信息
	returnInfo["appId"] = friend.AppId
	returnInfo["sex"] = friend.Sex
	returnInfo["avatar"] = friend.Avatar
	returnInfo["id"] = friend.UserId
	returnInfo["name"] = friend.Username
	returnInfo["mark_id"] = friend.MarkId
	returnInfo["uid"] = friend.Uid
	returnInfo["depositAddress"] = friend.DepositAddress
	returnInfo["disableDeadline"] = friend.CloseUntil
	returnInfo["publicKey"] = friend.PublicKey
	returnInfo["identification"] = friend.Identification
	returnInfo["identificationInfo"] = friend.IdentificationInfo
	returnInfo["isSetPayPwd"] = friend.IsSetPayPwd

	//同一个公司可以查看职位
	if friend.CompanyId == user.CompanyId {
		returnInfo["com_id"] = friend.CompanyId
		returnInfo["position"] = friend.Position
	}

	//好友之间的信息
	friendInfo, err := orm.FindFriendById(userID, friendID)
	if err != nil {
		logFriend.Info("FindFreind query db failed", "err_msg", err)
		return nil, result.NewError(result.DbConnectFail)
	}
	if friendInfo != nil {
		returnInfo["isBlocked"] = friendInfo.IsBlocked
		noDisturbing := friendInfo.DND
		returnInfo["noDisturbing"] = noDisturbing
		stickyOnTop := friendInfo.Top
		returnInfo["stickyOnTop"] = stickyOnTop
		returnInfo["remark"] = friendInfo.Remark
		//TODO 好友详细备注信息
		extRemark, err := getFriendExtRemark(friendInfo.ExtRemark)
		if err == nil {
			returnInfo["extRemark"] = extRemark
			/*returnInfo["telephones"] = extRemark.Telephones
			returnInfo["description"] = extRemark.Description
			returnInfo["pictures"] = extRemark.Pictures*/
		} else {
			returnInfo["extRemark"] = struct {
			}{}
		}
		returnInfo["addTime"] = friendInfo.AddTime
		//是好友就返回1
		returnInfo["isFriend"] = 1
	} else {
		//不是好友就返回2
		returnInfo["isFriend"] = 2
		returnInfo["extRemark"] = struct {
		}{}
	}

	//管理员客服可以查看的信息  或者自己看自己的信息
	if (user.UserLevel == types.RoomLevelManager) || (user.UserLevel == types.RoomLevelMaster) || (userID == friendID) {
		chain, err := orm.GetIsChain(friendID)
		if err != nil {
			logFriend.Info("GetIsChain failed", "err_msg", err)
			return nil, result.NewError(result.DbConnectFail)
		}
		//判断是否上链，默认未上链 0
		var ischain bool = false
		if chain == 1 {
			ischain = true
		}
		returnInfo["isChain"] = ischain
		returnInfo["privateKey"] = friend.PrivateKey
		returnInfo["account"] = friend.Account
		returnInfo["phone"] = friend.Phone
		returnInfo["email"] = friend.Email
		verified := utility.ToInt(friend.Verified)
		returnInfo["verified"] = verified
		returnInfo["userLevel"] = friend.UserLevel
		returnInfo["com_id"] = friend.CompanyId
		returnInfo["position"] = friend.Position
		returnInfo["code"] = friend.InviteCode
	}

	return returnInfo, nil
}

//通过uid查询
func UserListByUid(userID, appId, uid string) (interface{}, error) {
	//用户是否存在
	b := CheckUserExist(userID)
	if !b {
		logFriend.Warn("FriendInfo", "warn", "UserNotExists", "userId", userID)
		return nil, result.NewError(result.UserNotExists)
	}
	//好友的信息
	friend, err := orm.GetUserInfoByUid(appId, uid)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if friend == nil {
		logFriend.Warn("FriendInfo", "warn", "UserNotExists", "uid", uid)
		return nil, result.NewError(result.UserNotExists)
	}

	//自己的信息
	user, err := orm.GetUserInfoById(userID)
	if err != nil {
		logFriend.Info("FriendInfo query db failed", "err_msg", err)
		return nil, result.NewError(result.DbConnectFail)
	}

	var returnInfo = make(map[string]interface{})

	returnInfo["appId"] = friend.AppId
	returnInfo["sex"] = friend.Sex
	returnInfo["avatar"] = friend.Avatar
	returnInfo["id"] = friend.UserId
	returnInfo["name"] = friend.Username
	returnInfo["mark_id"] = friend.MarkId
	returnInfo["uid"] = friend.Uid
	returnInfo["depositAddress"] = friend.DepositAddress
	returnInfo["disableDeadline"] = friend.CloseUntil
	returnInfo["publicKey"] = friend.PublicKey
	returnInfo["identification"] = friend.Identification
	returnInfo["identificationInfo"] = friend.IdentificationInfo
	returnInfo["isSetPayPwd"] = friend.IsSetPayPwd

	//TODO 同一个公司可以查看职位  暂时未用到
	if friend.CompanyId == user.CompanyId {
		returnInfo["com_id"] = friend.CompanyId
		returnInfo["position"] = friend.Position
	}

	//存储的关系
	friendInfo, err := orm.FindFriendById(userID, friend.UserId)
	if err != nil {
		logFriend.Info("FindFreind query db failed", "err_msg", err)
		return nil, result.NewError(result.DbConnectFail)
	}
	if friendInfo != nil {
		returnInfo["isBlocked"] = friendInfo.IsBlocked
		noDisturbing := friendInfo.DND
		returnInfo["noDisturbing"] = noDisturbing
		stickyOnTop := friendInfo.Top
		returnInfo["stickyOnTop"] = stickyOnTop
		returnInfo["remark"] = friendInfo.Remark
		//TODO 好友详细备注信息
		extRemark, err := getFriendExtRemark(friendInfo.ExtRemark)
		if err == nil {
			returnInfo["extRemark"] = extRemark
			/*returnInfo["telephones"] = extRemark.Telephones
			returnInfo["description"] = extRemark.Description
			returnInfo["pictures"] = extRemark.Pictures*/
		} else {
			returnInfo["extRemark"] = struct {
			}{}
		}
		returnInfo["addTime"] = friendInfo.AddTime
		//是好友就返回1
		returnInfo["isFriend"] = 1
	} else {
		//不是好友就返回2
		returnInfo["isFriend"] = 2
		returnInfo["extRemark"] = struct {
		}{}
	}

	//管理员客服可以查看的信息  或者自己看自己的信息
	if (user.UserLevel == types.RoomLevelManager) || (user.UserLevel == types.RoomLevelMaster) || (userID == friend.UserId) {
		chain, err := orm.GetIsChain(friend.UserId)
		if err != nil {
			logFriend.Info("GetIsChain failed", "err_msg", err)
			return nil, result.NewError(result.DbConnectFail)
		}
		//判断是否上链，默认未上链 0
		var ischain bool = false
		if chain == 1 {
			ischain = true
		}
		returnInfo["isChain"] = ischain
		returnInfo["privateKey"] = friend.PrivateKey
		returnInfo["account"] = friend.Account
		returnInfo["phone"] = friend.Phone
		returnInfo["email"] = friend.Email
		verified := utility.ToInt(friend.Verified)
		returnInfo["verified"] = verified
		returnInfo["userLevel"] = friend.UserLevel
		returnInfo["com_id"] = friend.CompanyId
		returnInfo["position"] = friend.Position
		returnInfo["code"] = friend.InviteCode
	}

	return returnInfo, nil
}

//通过[]uid查询
func UserListByUids(userID, appId string, uids []string) (map[string]interface{}, error) {
	returnInfo := make(map[string]interface{})
	var info []interface{}
	for _, uid := range uids {
		res, err := UserListByUid(userID, appId, uid)
		if err != nil {
			//res = nil
			continue
		}
		info = append(info, res)
	}
	returnInfo["userList"] = info
	return returnInfo, nil
}

func typicalCatLog(list []*types.ChatLog, userID, friendID, owner, query string, number int, queryType []string, startId int64) ([]*types.ChatLog, int64, error) {
	//-1代表以及没消息了
	if startId <= 0 {
		return list, -1, nil
	}

	//查找发送者的信息
	userInfo, err := orm.GetUserInfoById(userID)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	friendInfo, err := orm.GetUserInfoById(friendID)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	//查找备注
	remark := ""
	relationInfo, err := orm.FindFriendById(userID, friendID)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	if relationInfo != nil {
		remark = relationInfo.Remark
	}

	//查询消息记录
	nextLog, logs, err := orm.FindTypicalPrivateLog(userID, friendID, owner, startId, number-len(list), queryType)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}

	if len(logs) <= 0 {
		return nil, -1, nil
	}
	//封装返回结果
	for _, r := range logs {
		if r.SenderId == userID {
			if userInfo != nil {
				r.Username = userInfo.Username
				r.Avatar = userInfo.Avatar
			}
		} else {
			if friendInfo != nil {
				r.Username = friendInfo.Username
				r.Avatar = friendInfo.Avatar
			}
		}
		info, err := GetChatLogAsUser(userID, r)
		if err != nil {
			logRoom.Warn("GetChatLogAsRoom error", "log info", info)
			continue
		}

		if info.IsSnap == types.IsSnap {
			continue
		}
		isAppend := true
		//根据query筛选文件名,上传者群名/备注/昵称
		if query != "" {
			isAppend = false
			if !isAppend {
				//昵称
				senderInfo := info.SenderInfo.(map[string]interface{})
				username := utility.ToString(senderInfo["name"])
				if strings.Contains(username, query) {
					isAppend = true
				}
			}
			if !isAppend {
				//不匹配 根据文件名筛选
				msg := info.Msg.(map[string]interface{})
				filename := utility.ToString(msg["name"])
				if strings.Contains(filename, query) {
					isAppend = true
				}
			}
			if !isAppend {
				//不匹配 查询群昵称
				nickname := orm.GetMemberName(info.TargetId, info.FromId)
				if strings.Contains(nickname, query) {
					isAppend = true
				}
			}
			if !isAppend && userID != info.FromId {
				if strings.Contains(remark, query) {
					isAppend = true
				}
			}
		}
		if isAppend {
			list = append(list, info)
		}
	}

	if len(list) < number {
		if utility.ToInt64(nextLog) < 0 {
			//完成全部查找
			return list, utility.ToInt64(nextLog), nil
		}
		//直到找到足够number数量的记录
		return typicalCatLog(list, userID, friendID, owner, query, number, queryType, utility.ToInt64(nextLog))
	}
	return list, utility.ToInt64(nextLog), nil
}

//获取好友的特定消息记录
func FindTypicalChatLog(userID, friendID, owner, query string, startId string, number int, queryType []string) (interface{}, error) {
	//判断用户是否存在
	bool := CheckUserExist(friendID)
	if !bool {
		logFriend.Warn("FindTypicalChatLog", "warn", "UserNotExists", "userId", friendID)
		return nil, result.NewError(result.UserNotExists)
	}
	bool = CheckUserExist(userID)
	if !bool {
		logFriend.Warn("FindTypicalChatLog", "warn", "UserNotExists", "userId", userID)
		return nil, result.NewError(result.UserNotExists)
	}
	if userID == friendID {
		logFriend.Warn("FindTypicalChatLog", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID)
		return nil, result.NewError(result.CanNotOperateSelf)
	}
	//判断是否是好友
	err := boolIsExistIsFriend(userID, friendID)
	if err != nil {
		return nil, err
	}

	start := utility.ToInt64(startId)
	if start == 0 {
		//获取最新消息id
		start, err = orm.FindLastCatLogId(userID, friendID)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if start == 0 {
			return nil, nil
		}
		start += 1
		if err != nil {
			return nil, result.NewError(result.ConvFail)
		}
	}

	list := make([]*types.ChatLog, 0)
	list, nextLog, err := typicalCatLog(list, userID, friendID, owner, query, number, queryType, start)

	var ret = make(map[string]interface{})
	ret["logs"] = list
	ret["nextLog"] = nextLog
	return ret, err
}

func privateChatLog(list []*types.ChatLog, userID, friendID string, number int, startId int64) ([]*types.ChatLog, int64, error) {
	//-1代表以及没消息了
	if startId <= 0 {
		return list, -1, nil
	}

	//查找发送者的信息
	userInfo, err := orm.GetUserInfoById(userID)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	friendInfo, err := orm.GetUserInfoById(friendID)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	//查找备注
	remark := ""
	relationInfo, err := orm.FindFriendById(userID, friendID)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	if relationInfo != nil {
		remark = relationInfo.Remark
	}

	//查询消息记录
	nextLog, logs, err := orm.FindPrivateChatLogs(userID, friendID, startId, number-len(list))
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}

	if len(logs) <= 0 {
		return nil, -1, nil
	}
	//封装返回结果
	for _, r := range logs {
		if r.SenderId == userID {
			if userInfo != nil {
				r.Username = userInfo.Username
				r.Avatar = userInfo.Avatar
			}
		} else {
			if friendInfo != nil {
				r.Username = friendInfo.Username
				r.Avatar = friendInfo.Avatar
				r.Remark = remark
			}
		}
		info, err := GetChatLogAsUser(userID, r)
		if err != nil {
			continue
		}

		if info.IsSnap == types.IsSnap {
			continue
		}

		//去除掉 群秘钥更新通知
		if info.MsgType == types.Alert {
			ck := proto.NewMsg(info.GetCon())

			//过滤掉赞赏通知
			if ck.GetTypeCode() == types.AlertPraise {
				continue
			}

			if ck.GetTypeCode() == types.AlertUpdateSKey {
				continue
			}
		}
		list = append(list, info)
	}

	if len(list) < number {
		if utility.ToInt64(nextLog) < 0 {
			//完成全部查找
			return list, utility.ToInt64(nextLog), nil
		}
		//直到找到足够number数量的记录
		return privateChatLog(list, userID, friendID, number, utility.ToInt64(nextLog))
	}
	return list, utility.ToInt64(nextLog), nil
}

//获取好友消息记录
func FindCatLog(appId, userID, friendID string, startId string, number int) (interface{}, error) {
	//判断用户是否存在
	bool := CheckUserExist(friendID)
	if !bool {
		logFriend.Warn("FindCatLog", "warn", "UserNotExists", "userId", friendID)
		return nil, result.NewError(result.UserNotExists)
	}
	bool = CheckUserExist(userID)
	if !bool {
		logFriend.Warn("FindCatLog", "warn", "UserNotExists", "userId", userID)
		return nil, result.NewError(result.UserNotExists)
	}
	if userID == friendID {
		logFriend.Warn("FindCatLog", "warn", "CanNotOperateSelf", "userId", userID, "friendId", friendID)
		return nil, result.NewError(result.CanNotOperateSelf)
	}

	app := app.GetApp(appId)
	if app.IsOtc == types.IsOtc {
	} else {
		//判断是否是好友
		err := boolIsExistIsFriend(userID, friendID)
		if err != nil {
			return nil, err
		}
	}

	start := utility.ToInt64(startId)
	var err error
	if start <= 0 {
		//获取最新消息id
		start, err = orm.FindLastCatLogId(userID, friendID)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if start == 0 {
			return nil, nil
		}
		start += 1
		if err != nil {
			return nil, result.NewError(result.ConvFail)
		}
	}

	list, nextLog, err := privateChatLog([]*types.ChatLog{}, userID, friendID, number, start)
	if list == nil {
		list = []*types.ChatLog{}
	}

	var ret = make(map[string]interface{})
	ret["logs"] = list
	ret["nextLog"] = utility.ToString(nextLog)
	return ret, err
}

//获取所有好友未读消息统计
func GetAllFriendUnreadMsg1(userId string) (interface{}, error) {
	//查询所有好友
	friends, err := orm.FindFriendsId(userId)
	if err != nil {
		logFriend.Warn("FindFriendByUserId query db failed", "err_msg", err)
		return nil, result.NewError(result.DbConnectFail)
	}
	var data = make(map[string]interface{})
	var infos = make([]map[string]interface{}, 0)
	for _, fid := range friends {
		//查询好友的未读聊天记录数
		count, err := orm.GetUnReadNumber(userId, fid)
		if err != nil {
			logFriend.Warn("FindUnReadNum query db failed", "err_msg", err)
			return nil, result.NewError(result.DbConnectFail)
		}

		if count > 0 {
			var one = make(map[string]interface{})
			//查询好友信息
			info, err := orm.FindFriendById(userId, fid)
			if err != nil {
				logFriend.Warn("FindFriendInfoByUserId query db failed", "err_msg", err)
				return nil, result.NewError(result.DbConnectFail)
			}
			var senderInfo = make(map[string]interface{})
			if info.Remark != "" {
				senderInfo["nickname"] = info.Remark
			} else {
				senderInfo["nickname"] = info.Username
			}
			senderInfo["avatar"] = info.Avatar

			//查询好友第一条未读聊天记录
			firstLog, err := orm.FindFirstPrivateMsg(userId, fid)
			if err != nil {
				logFriend.Warn("FindFirstMsg query db failed", "err_msg", err)
				return nil, result.NewError(result.DbConnectFail)
			}
			var lastLog = make(map[string]interface{})
			if firstLog != nil {
				lastLog["logId"] = firstLog.Id
				lastLog["channelType"] = 3
				one["isSnap"] = firstLog.IsSnap
				lastLog["fromId"] = firstLog.SenderId
				lastLog["targetId"] = firstLog.ReceiveId
				lastLog["msgType"] = firstLog.MsgType
				con := utility.StringToJobj(firstLog.Content)
				if firstLog.MsgType == types.RedPack {
					err := ConvertRedPackInfoToSend(userId, con)
					if err != nil {
						logFriend.Error("GetAllFriendUnreadMsg1 ConvertRedPackInfoToSend", "err", err)
						return nil, result.NewError(result.DbConnectFail)
					}
				}
				lastLog["msg"] = con
				lastLog["datetime"] = firstLog.SendTime
				lastLog["senderInfo"] = senderInfo
			}

			one["id"] = fid
			one["number"] = count
			one["lastLog"] = lastLog
			infos = append(infos, one)
		}
	}
	data["infos"] = infos
	return data, nil
}

func IsFriend(userID, friendID string) (map[string]bool, error) {
	var isFreindMap = make(map[string]bool)
	b := CheckUserExist(friendID)
	if !b {
		logFriend.Warn("IsFriend", "warn", "UserNotExists", "userId", friendID)
		return nil, result.NewError(result.UserNotExists)
	}

	b, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	isFreindMap["isFriend"] = b
	return isFreindMap, nil
}

func findSendReceiveInfo(userID, friendID string) (map[string]interface{}, map[string]interface{}, error) {
	//查询自己的信息
	user, err := orm.GetUserInfoById(userID)
	if err != nil {
		return nil, nil, result.NewError(result.DbConnectFail)
	}
	if user == nil {
		logFriend.Warn("findSendReceiveInfo", "warn", "UserNotExists", "userId", userID)
		return nil, nil, result.NewError(result.UserNotExists)
	}
	//查询对方的信息
	friend, err := orm.GetUserInfoById(friendID)
	if err != nil {
		return nil, nil, result.NewError(result.DbConnectFail)
	}
	if friend == nil {
		logFriend.Warn("findSendReceiveInfo", "warn", "UserNotExists", "friendId", friendID)
		return nil, nil, result.NewError(result.UserNotExists)
	}
	var userInfo = make(map[string]interface{})
	userInfo["id"] = user.UserId
	userInfo["name"] = user.Username
	userInfo["avatar"] = user.Avatar
	userInfo["uid"] = user.Uid
	if friend.CompanyId == user.CompanyId {
		userInfo["position"] = user.Position
	}

	var friendIDInfo = make(map[string]interface{})
	friendIDInfo["id"] = friend.UserId
	friendIDInfo["name"] = friend.Username
	friendIDInfo["avatar"] = friend.Avatar
	friendIDInfo["uid"] = friend.Uid
	if friend.CompanyId == user.CompanyId {
		friendIDInfo["position"] = friend.Position
	}

	return userInfo, friendIDInfo, nil
}

func SendPrintScreen(userID, friendID string) error {
	SendAlert(userID, friendID, types.ToUser, []string{userID, friendID}, types.Alert, proto.ComposePrintScreen(userID))
	return nil
}

//获取来源
func GetFriendSource(userId, friendId string) (string, error) {
	log, err := orm.FindAddFriendApplyLog(userId, friendId)
	if err != nil {
		return "", err
	}
	if log == nil {
		return "未知来源", nil
	}

	return ConverFriendSource(log.Source, userId, log.ApplyUser)
}

//获取来源
func GetFriendSourceV2(userId, friendId string) (interface{}, error) {
	log, err := orm.FindAddFriendApplyLog(userId, friendId)
	if err != nil {
		return "", err
	}
	if log == nil {
		return "未知来源", nil
	}

	return ConverFriendSourceV2(log.Source, log.ApplyUser)
}

//组装消息  applyUserId发起申请人的id
func ConverFriendSource(source string, userId, applyUserId string) (string, error) {
	sourceType, sourceId := unConverSource(source)
	if userId == applyUserId {
		if sourceType == types.Search {
			return "通过搜索添加对方", nil
		} else if sourceType == types.Scan {
			return "通过扫一扫添加对方", nil
		} else if sourceType == types.Group {
			room, err := orm.FindRoomById(sourceId, types.RoomDeletedOrNot)
			if err != nil {
				return "", err
			}
			if room == nil {
				return "未知来源", nil
			}
			str := "通过【" + room.Name + "】添加对方"
			return str, nil
		} else if sourceType == types.Share {
			infos, err := orm.GetUserInfoById(sourceId)
			if err != nil {
				return "", err
			}
			if infos == nil {
				return "未知来源", nil
			}
			str := "通过【" + infos.Username + "】分享添加对方"
			return str, nil
		} else {
			return "你添加对方为好友", nil
		}
	} else {
		if sourceType == types.Search {
			return "对方通过搜索添加", nil
		} else if sourceType == types.Scan {
			return "对方通过扫一扫添加", nil
		} else if sourceType == types.Group {
			room, err := orm.FindRoomById(sourceId, types.RoomDeletedOrNot)
			if err != nil {
				return "", err
			}
			if room == nil {
				return "未知来源", nil
			}
			str := "对方通过【" + room.Name + "】添加"
			return str, nil
		} else if sourceType == types.Share {
			infos, err := orm.GetUserInfoById(sourceId)
			if err != nil {
				return "", err
			}
			if infos == nil {
				return "未知来源", nil
			}
			str := "对方通过【" + infos.Username + "】分享添加"
			return str, nil
		} else {
			return "对方添加你为好友", nil
		}
	}
}

//组装消息  applyUserId发起申请人的id
func ConverFriendSourceV2(source string, applyUserId string) (interface{}, error) {
	sourceType, sourceId := unConverSource(source)
	sourceName := ""

	if sourceType == types.Group {
		room, err := orm.FindRoomById(sourceId, types.RoomDeletedOrNot)
		if err != nil {
			return "", err
		}
		if room != nil {
			sourceName = room.Name
		}
	} else if sourceType == types.Share {
		info, err := orm.GetUserInfoById(sourceId)
		if err != nil {
			return "", err
		}
		if info != nil {
			sourceName = info.Username
		}
	}

	maps := make(map[string]interface{})
	maps["sourceType"] = sourceType
	maps["sourceId"] = sourceId
	maps["sourceName"] = sourceName
	maps["applyUser"] = applyUserId

	return maps, nil
}

func unConverSource(source string) (sourceType int, sourceId string) {
	ret := utility.StringToJobj(source)
	if len(ret) > 0 {
		return utility.ToInt(ret["sourceType"].(float64)), ret["sourceId"].(string)
	} else {
		//为了处理以前的数据
		id, err := strconv.Atoi(source)
		if err != nil || id <= 0 {
			return types.Unknow, ""
		} else {
			return types.Group, strconv.Itoa(id)
		}
	}
}

//tp = 0  修改问题答案  需要同时传入问题答案
//tp = 1  设置需要回答问题  需要同时传入问题答案
//tp = 2  设置不需要回答问题
func Question(userID string, tp int, question, answer string) error {
	var err error
	switch tp {
	case 0:
		err = orm.SetQuestionandAnswer(userID, question, answer)
	case types.NeedAnswer:
		err = orm.SetNeedAnswer(userID, question, answer)
	case types.NotNeedAnswer:
		err = orm.SetNotNeedAnswer(userID)
	}
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//是否需要验证
func Confirm(userID string, tp int) error {
	err := orm.IsNeedConfirm(userID, tp)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//验证答案是否正确
func CheckAnswer(friendId string, answer string) (bool, error) {
	conf, err := GetAddFriendConfByUserId(friendId)
	if err != nil {
		return false, result.NewError(result.DbConnectFail)
	}
	if conf == nil {
		return true, nil
	}
	if conf.Answer == answer {
		return true, nil
	}
	return false, nil
}

//查找加好友配置
func GetAddFriendConfByUserId(userId string) (*types.AddFriendConf, error) {
	//查找appid
	infos, err := orm.GetUserInfoById(userId)
	if err != nil {
		return nil, err
	}
	if infos == nil {
		return nil, nil
	}
	appid := infos.AppId

	b := false
	appIds := strings.Split(cfg.NotConfirmAppId.AppIds, ",")
	for _, id := range appIds {
		if id == appid {
			b = true
		}
	}
	if b {
		conf := &types.AddFriendConf{
			UserId:      userId,
			NeedAnswer:  strconv.Itoa(types.NotNeedAnswer),
			NeedConfirm: strconv.Itoa(types.NeedConfirm),
		}
		return conf, nil
	}
	return orm.FindAddFriendConfByUserId(userId)
}

//添加好友
//sourceType  1  【通过搜索添加】  2  【通过扫一扫添加】  3 通过好友分享(uid)  4通过群
func AddFriendNotConfirm(appId, userID, friendID, remark, reason, sourceId string, sourceType int, answer string) (int, error) {
	//不能对自己操作
	if userID == friendID {
		logFriend.Debug("can not add friend confirm", "warn", "CanNotOperateSelf", "userId", userID, "friendID", friendID)
		return 0, result.NewError(result.CanNotOperateSelf)
	}

	//判断用户是否存在
	isExists := CheckUserExist(friendID)
	if !isExists {
		return 0, result.NewError(result.UserNotExists)
	}

	//验证是否是好友
	isFriend, err := orm.CheckIsFriend(userID, friendID, types.FriendIsNotDelete)
	if isFriend {
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		logFriend.Warn("AddFriendNotConfirm", "warn", "IsFriendAlready", "userId", userID, "friendId", friendID)
		return 0, result.NewError(result.IsFriendAlready)
	}

	if sourceType == types.Share {
		user, err := orm.GetUserInfoByMarkId(appId, sourceId)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		sourceId = user.UserId
	}

	isFromRoom := true
	//不是通过群添加好友
	if sourceType != types.Group {
		isFromRoom = false
	}

	if isFromRoom {
		//判断群是否存在
		room, err := orm.FindRoomById(sourceId, types.RoomDeletedOrNot)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if room == nil || room.IsDelete == types.IsDelete {
			logFriend.Warn("AddFriendNotConfirm", "warn", "RoomNotExists", "roomId", sourceId)
			return 0, result.NewError(result.RoomNotExists)
		}
		//判断该群是否允许添加好友
		if room.CanAddFriend == types.CanNotAddFriend {
			logFriend.Warn("AddFriendNotConfirm", "warn", "CanNotAddFriendInRoom")
			return 0, result.NewError(result.CanNotAddFriendInRoom)
		}

		//判断双方是否都在群里
		userInRoom1, err := CheckUserInRoom(userID, sourceId)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		userInRoom2, err := CheckUserInRoom(friendID, sourceId)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if !userInRoom1 || !userInRoom2 {
			return 0, result.NewError(result.UserIsNotInRoom)
		}
	}

	/*	//查询cid
		uFriend, err := orm.GetUserInfoById(friendID)
		if err != nil {
			logFriend.Warn("FindCid db dailed", "err_msg", err)
			return 0, result.NewError(result.DbConnectFail)
		}*/

	senderInfo, receiveInfo, err := findSendReceiveInfo(userID, friendID)
	if err != nil {
		return 0, err
	}

	//text := senderInfo["name"].(string) + "请求添加您为好友"

	//判断用户平台
	info1, err := orm.GetUserInfoById(userID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	info2, err := orm.GetUserInfoById(friendID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	if info1 == nil || info2 == nil {
		logFriend.Warn("AddFriendNotConfirm", "warn", "UserNotExists", "user", info1, "friend", info2)
		return 0, result.NewError(result.UserNotExists)
	}
	if info1.AppId != "1002" || info2.AppId != "1002" {
		logFriend.Warn("AddFriendNotConfirm", "warn", "platform not consistency", "user appId", info1.AppId, "friend appId", info2.AppId)
		return 0, result.NewError(result.PermissionDeny)
	}

	//判断请求是否存在
	cou, err := orm.GetFriendApplyCount(userID, friendID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}

	var logId string
	//如果请求存在，只执行一些更新操作
	if cou > 0 {
		id, err := orm.UpdateApply(userID, friendID, types.IsFriend, reason, remark, converSource(sourceType, sourceId), types.AwaitState)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		logId = utility.ToString(id)
	} else {
		//添加
		id, err := orm.AppendApplyLog(friendID, userID, reason, converSource(sourceType, sourceId), remark, types.AwaitState, types.IsFriend)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		logId = utility.ToString(id)
	}
	//获取来源
	//source, err := GetFriendSource(friendID, userID)
	source, err := ConverFriendSource(converSource(sourceType, sourceId), friendID, userID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	source2, err := ConverFriendSourceV2(converSource(sourceType, sourceId), userID)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	proto.SendAddFriendNotification(friendID, logId, reason, source, source2, types.AwaitState, senderInfo, receiveInfo)

	//不需要验证
	err = HandleFriendRequest(friendID, userID, types.AcceptState)
	if err != nil {
		logFriend.Warn("HandleFriendRequest error", "err", err)
	}
	return types.AddSuccess, nil
}

//加入黑名单
func BlockFriend(userId, friendId string) error {
	err := orm.SetFriendIsBlock(userId, friendId, types.IsBlocked)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//移出黑名单
func UnblockFriend(userId, friendId string) error {
	err := orm.SetFriendIsBlock(userId, friendId, types.IsNotBlocked)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//是否被拉入黑名单:true被拉入 false没被拉入
func CheckIsBlocked(userId, friendId string) (bool, error) {
	info, err := orm.FindFriendById(userId, friendId)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}

	switch info.IsBlocked {
	case types.IsBlocked:
		return true, nil
	default:
		return false, nil
	}
}

//黑名单列表
func BlockedFriends(userId string) (interface{}, error) {
	friends, err := orm.BlockedFriends(userId)
	if err != nil {
		logFriend.Warn("friend list query db failed", "err_msg", err)
		return nil, result.NewError(result.DbConnectFail)
	}
	data := make([]types.FriendInfoApi, 0)

	for _, friend := range friends {
		one := types.FriendInfoApi{
			Uid:                friend.Uid,
			Id:                 friend.FriendId,
			Name:               friend.Username,
			Avatar:             friend.Avatar,
			Remark:             friend.Remark,
			PublicKey:          friend.PublicKey,
			NoDisturbing:       friend.DND,
			CommonlyUsed:       friend.Type,
			OnTop:              friend.Top,
			IsDelete:           friend.IsDelete,
			AddTime:            friend.AddTime,
			DepositAddress:     friend.DepositAddress,
			Identification:     friend.Identification,
			IdentificationInfo: friend.IdentificationInfo,
			IsBlocked:          friend.IsBlocked,
		}
		data = append(data, one)
	}
	return map[string]interface{}{"userList": data}, nil
}
