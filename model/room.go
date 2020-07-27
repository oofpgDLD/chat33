package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var logRoom = log15.New("logic", "chat/logic")

type RoomBaseInfo struct {
	Id                 string      `json:"id"`
	MarkId             string      `json:"markId"`
	Name               string      `json:"name"`
	Avatar             string      `json:"avatar"`
	Encrypt            int         `json:"encrypt"`
	CanAddFriend       int         `json:"canAddFriend"`
	JoinPermission     int         `json:"joinPermission"`
	RecordPermission   int         `json:"recordPermission"`
	MemberNumber       int64       `json:"memberNumber"`
	ManagerNumber      int64       `json:"managerNumber"`
	DisableDeadline    int64       `json:"disableDeadline"`
	SystemMsg          interface{} `json:"systemMsg"`
	Identification     int         `json:"identification"`
	IdentificationInfo string      `json:"identificationInfo"`
}

type RoomInfo struct {
	*RoomBaseInfo
	OnlineNumber  int             `json:"onlineNumber"`
	NoDisturbing  int             `json:"noDisturbing"`
	OnTop         int             `json:"onTop"`
	MemberLevel   int             `json:"memberLevel"`
	RoomNickname  string          `json:"roomNickname"`
	RoomMutedType int             `json:"roomMutedType"`
	MutedNumber   int64           `json:"mutedNumber"`
	MutedType     int             `json:"mutedType"`
	Deadline      int64           `json:"deadline"`
	Users         []*RoomUserList `json:"users"`
}

type RoomListInfo struct {
	Id                 string `json:"id"`
	MarkId             string `json:"markId"`
	Name               string `json:"name"`
	Avatar             string `json:"avatar"`
	NoDisturbing       int    `json:"noDisturbing"`
	CommonlyUsed       int    `json:"commonlyUsed"`
	OnTop              int    `json:"onTop"`
	Encrypt            int    `json:"encrypt"`
	DisableDeadline    int64  `json:"disableDeadline"`
	Identification     int    `json:"identification"`
	IdentificationInfo string `json:"identificationInfo"`
}

type RoomUserList struct {
	Id                 string      `json:"id"`
	Nickname           string      `json:"nickname"`
	PublicKey          string      `json:"publicKey"`
	RoomNickname       string      `json:"roomNickname"`
	Avatar             string      `json:"avatar"`
	MemberLevel        int         `json:"memberLevel"`
	RoomMutedType      int         `json:"roomMutedType"`
	MutedType          int         `json:"mutedType"`
	Source             string      `json:"source"`
	Source2            interface{} `json:"source2"`
	Deadline           int64       `json:"deadline"`
	Identification     int         `json:"identification"`
	IdentificationInfo string      `json:"identificationInfo"`
}

type SystemMessage struct {
	Id         string `json:"logId"`
	SenderName string `json:"senderName"`
	Content    string `json:"content"`
	Datetime   int64  `json:"datetime"`
}

func CheckUserInRoom(userId, roomId string) (bool, error) {
	roomInfo, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err != nil {
		return false, err
	}
	if roomInfo == nil {
		return false, nil
	}

	member, err := orm.FindNotDelMember(roomId, userId)
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return true, nil
}

//普通用户id
func normalMembersId(roomId string) []string {
	members, err := orm.FindNotDelMembers(roomId, types.SearchAll)
	if err != nil || len(members) < 1 {
		return nil
	}
	var ret = make([]string, 0)
	for _, m := range members {
		if m.Level == types.RoomLevelNomal {
			ret = append(ret, m.RoomMember.UserId)
		}
	}
	return ret
}

func getMemberJoinTime(roomId, userId string) int64 {
	info, err := orm.FindNotDelMember(roomId, userId)
	if err != nil || info == nil {
		return 0
	}
	return info.RoomMember.CreateTime
}

//加入群
func JoinInRoom(userId, roomId string) {
	if user, ok := router.GetUser(userId); ok {
		roomChannelId := types.GetRoomRouteById(roomId)
		if channel, ok := router.GetChannel(roomChannelId); !ok || channel == nil {
			logRoom.Error("channel room not exist", "userID", userId, "roomId", roomId)
		} else {
			user.Subscribe(channel)
		}
	}
}

// 根据群号获取在线人数
func GetOnlineNumber(roomId string) int {
	channelId := types.GetRoomRouteById(roomId)
	if channel, _ := router.GetChannel(channelId); channel != nil {
		return channel.GetRegisterNumber()
	}
	return 0
}

// 初始化群的channel
func initRoomChannel(roomId string) {
	channelId := types.GetRoomRouteById(roomId)
	router.AppendChannel(channelId, router.NewChannel(channelId))
}

func roomMemberSubscribe(roomId string, members []string) {
	roomChannelId := types.GetRoomRouteById(roomId)
	cl, _ := router.GetChannel(roomChannelId)
	if cl != nil {
		for _, v := range members {
			if user, ok := router.GetUser(v); ok {
				user.Subscribe(cl)
			}
		}
	}
}

func RemoveRoomChannel(roomId string) {
	channelId := types.GetRoomRouteById(roomId)
	router.DeleteChannel(channelId)
}

func RoomMemberUnSubscribe(roomId string, members []string) {
	roomChannelId := types.GetRoomRouteById(roomId)
	cl, _ := router.GetChannel(roomChannelId)
	if cl != nil {
		for _, v := range members {
			if user, ok := router.GetUser(v); ok {
				user.UnSubscribe(cl)
			}
		}
	}
}

//检查用户在指定群中是否被禁言
func CheckMemberMuted(roomId, userId string) (bool, error) {
	mutedType, err := orm.GetRoomMutedType(roomId)
	if err != nil {
		return true, result.NewError(result.DbConnectFail)
	}
	switch mutedType {
	case types.AllMuted:
		return true, nil
	case types.AllSpeak:
		return false, nil
	case types.Whitelist:
		mutedType, _ := orm.GetRoomUserMuted(roomId, userId)
		if mutedType != types.Whitelist {
			return true, nil
		}
	case types.Blacklist:
		mutedType, deadline := orm.GetRoomUserMuted(roomId, userId)
		if mutedType == types.Blacklist && deadline > utility.NowMillionSecond() {
			return true, nil
		}
	}
	return false, nil
}

//将sourceType sourceId 封装成json 格式
func ConvertSourece(sourceType int, sourceId string) string {
	var m = make(map[string]interface{})
	switch sourceType {
	case types.Search:
	case types.Scan:
	case types.Share:
		m["sourceId"] = sourceId
	case types.Invite:
		//根据id 获取uid
		user, err := orm.GetUserInfoById(sourceId)
		if err != nil {
			return ""
		}
		if user != nil {
			m["sourceId"] = user.Uid
		}
		//
		//maps, err := db.GetUserInfoByID(sourceId)
		//if err != nil || len(maps) < 1 {
		//	return ""
		//}
		//info := maps[0]["uid"]
		//m["sourceId"] = info
	default:
		return ""
	}
	m["sourceType"] = sourceType
	return utility.StructToString(m)
}

//通过source json 获取source内容
func getJoinSource(appId, roomId, source string) (string, error) {
	roomInfo, err := orm.FindRoomById(roomId, types.RoomDeletedOrNot)
	if err != nil {
		return "", err
	}
	roomName := roomInfo.Name
	m := utility.StringToJobj(source)
	if _, ok := m["sourceType"]; !ok {
		str := "申请加入%s"
		return fmt.Sprintf(str, roomName), nil
		//return "", fmt.Errorf("param error")
	}
	sourceType := utility.ToInt(m["sourceType"])

	switch sourceType {
	case types.Search:
		str := "通过搜索申请加入%s"
		return fmt.Sprintf(str, roomName), nil
	case types.Scan:
		str := "通过扫一扫申请加入%s"
		return fmt.Sprintf(str, roomName), nil
	case types.Share:
		fallthrough
	case types.Invite:
		var str string
		if sourceType == types.Share {
			str = "通过%s分享申请加入%s"
		} else if sourceType == types.Invite {
			str = "通过%s邀请进入%s"
		}

		if _, ok := m["sourceId"]; !ok {
			return "", fmt.Errorf("param error")
		}
		userId := utility.ToString(m["sourceId"])
		userInfo, err := orm.GetUserInfoByUid(appId, utility.ToString(userId))
		if err != nil {
			return "", err
		}
		userName := ""
		if userInfo != nil {
			userName = userInfo.Username
		}

		return fmt.Sprintf(str, userName, roomName), nil
	default:
		return "", fmt.Errorf("param error")
	}
}

//通过source json 获取source内容
func getJoinSourceV2(appId, roomId, source string) (interface{}, error) {
	roomInfo, err := orm.FindRoomById(roomId, types.RoomDeletedOrNot)
	if err != nil {
		return "", err
	}
	roomName := roomInfo.Name

	m := utility.StringToJobj(source)
	sourceType := utility.ToInt(m["sourceType"])
	uid := utility.ToString(m["sourceId"])
	sourceName := ""
	sourceId := ""
	if uid != "" {
		userInfo, err := orm.GetUserInfoByUid(appId, uid)
		if err != nil {
			return "", err
		}
		if userInfo != nil {
			sourceName = userInfo.Username
			sourceId = userInfo.UserId
		}
	}

	maps := make(map[string]interface{})
	maps["sourceType"] = sourceType
	maps["sourceId"] = sourceId
	maps["sourceName"] = sourceName
	maps["roomName"] = roomName

	return maps, nil
}

//随机MarkId 9位长度数字 10次尝试
func randomRoomMarkId() (string, error) {
	for count := 10; count > 0; count-- {
		randomId := utility.RandomRoomId()
		isExist, err := orm.CheckRoomMarkIdExist(randomId)
		if err != nil {
			return "", result.NewError(result.DbConnectFail)
		}
		if !isExist {
			return randomId, nil
		}
	}
	return "", result.NewError(result.NetWorkError)
}

//根据成员转成群名
func convertRoomName(creater string, members []string) (string, error) {
	roomName := ""
	name, err := getUsernameById(creater)
	if err != nil {
		return roomName, err
	}
	roomName += name
	if len(members) < 1 {
		return roomName + "创建的群聊", nil
	} else {
		for _, v := range members {
			name, err := getUsernameById(v)
			if err != nil {
				return roomName, err
			}
			roomName += "、" + name
			arry := []rune(roomName)
			if len(arry) > 20 {
				roomName = string(arry[:16]) + "...等"
			}
		}
		return roomName, nil
	}
}

//获取群成员详细信息
func getRoomMembers(roomId string, searchNumber int) ([]*RoomUserList, error) {
	var users = make([]*RoomUserList, 0)
	infos, err := orm.FindNotDelMembers(roomId, searchNumber)
	if err != nil {
		return nil, err
	}
	for _, r := range infos {
		var user RoomUserList
		user.Id = r.RoomMember.UserId
		user.PublicKey = r.PublicKey
		user.RoomNickname = r.UserNickname
		user.Nickname = r.Username
		user.MemberLevel = r.Level
		user.Avatar = r.Avatar
		user.RoomMutedType, _ = orm.GetRoomMutedType(roomId)
		user.MutedType, user.Deadline = orm.GetRoomUserMuted(roomId, user.Id)
		user.Source, _ = getJoinSource(r.AppId, roomId, r.Source)
		user.Source2, _ = getJoinSourceV2(r.AppId, roomId, r.Source)
		user.Identification = r.Identification
		user.IdentificationInfo = r.IdentificationInfo
		users = append(users, &user)
	}
	return users, nil
}

func filterMember(members []string) []string {
	//判断群成员是否有公钥
	ret := make([]string, 0)
	for _, v := range members {
		u, err := orm.GetUserInfoById(v)
		if err != nil {
			continue
		}
		if u.PublicKey != "" {
			ret = append(ret, v)
		}
	}
	return ret
}

//encrypt: 1加密群，2普通群
func CreateRoom(appId, creator, roomName, roomAvatar string, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted int, members []string) (interface{}, error) {
	//特殊判断 如果 是Tsc不使用加密群
	if appId == "1006" {
		encrypt = types.IsNotEncrypt
	}

	//判断创建的群聊个数
	createLimit, _ := orm.GetCreateRoomsLimit(appId, 1)
	createNumber, _ := orm.FindUserCreateRoomsNumber(creator)
	if createNumber >= createLimit && createLimit != 0 {
		logRoom.Warn("CreateRoom", "warn", "CreateRoomsOutOfLimit", "current number", createNumber, "limit number", createLimit)
		return nil, result.NewError(result.CreateRoomsOutOfLimit).SetExtMessage(ConvertCreateRoomsOutOfLimit(createLimit))
	}

	//如果是加密群 筛选掉没有公钥的成员
	if encrypt == types.IsEncrypt {
		members = filterMember(members)
	}

	//如果对方将你拉入黑名单，则无法邀请
	old := make([]string, 0)
	rejectUsers := make([]string, 0)
	old = append(old, members...)
	members = make([]string, 0)
	for _, u := range old {
		ok, err := CheckIsBlocked(u, creator)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if ok {
			rejectUsers = append(rejectUsers, u)
			continue
		}
		members = append(members, u)
	}

	//判断最大成员数量
	limitNumber, err := memberLimit(appId, nil)
	if err != nil {
		return nil, err
	}
	if limitNumber != 0 && len(members)+1 > limitNumber {
		logRoom.Warn("CreateRoom", "warn", "MembersOutOfLimit", "current number", members, "limit number", limitNumber)
		return nil, result.NewError(result.MembersOutOfLimit).SetExtMessage(ConvertRoomMembersLimit(limitNumber))
	}
	randomId, err := randomRoomMarkId()
	if err != nil {
		return nil, err
	}
	//聊天加密版
	roomName = "群聊"
	if roomName == "" {
		name, err := convertRoomName(creator, members)
		if err != nil {
			return nil, err
		}
		roomName = name
	}
	createTime := utility.NowMillionSecond()
	tx, err := orm.GetTx()
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	roomId, err := orm.CreateNewRoomV2(tx, creator, roomName, roomAvatar, encrypt, canAddFriend, joinPermission, recordPermission, adminMuted, masterMuted, randomId, createTime)
	if err != nil {
		tx.RollBack()
		return nil, result.NewError(result.DbConnectFail)
	}
	//添加群主
	err = orm.AddMember(tx, creator, utility.ToString(roomId), types.RoomLevelMaster, createTime, "")
	if err != nil {
		tx.RollBack()
		return nil, result.NewError(result.DbConnectFail)
	}
	//添加成员
	sendInvite := make([]string, 0)
	joined := make([]string, 0)
	info := &types.Room{
		Id:     utility.ToString(roomId),
		MarkId: randomId,
		Name:   roomName,
		Avatar: roomAvatar,
	}
	for _, u := range members {
		if roomInviteCheck(appId, creator, u, info) {
			err = orm.AddMember(tx, u, utility.ToString(roomId), types.RoomLevelNomal, createTime, "")
			if err != nil {
				tx.RollBack()
				return nil, result.NewError(result.DbConnectFail)
			}
			joined = append(joined, u)
		} else {
			sendInvite = append(sendInvite, u)
		}
	}
	members = joined
	err = tx.Commit()
	if err != nil {
		tx.RollBack()
		return nil, result.NewError(result.DbConnectFail)
	}

	for _, u := range sendInvite {
		sendInviteCard(creator, u, info)
	}

	room, err := roomInfoOfApiResult(creator, utility.ToString(roomId))
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if room == nil {
		return nil, result.NewError(result.RoomNotExists)
	}
	// add creater
	members = append(members, creator)
	//init room channel
	initRoomChannel(utility.ToString(roomId))
	logRoom.Debug("InitRoomChannel pass")
	//member subscribe
	roomMemberSubscribe(utility.ToString(roomId), members)
	logRoom.Debug("RoomMemberSubscribe pass")
	//send notification to all member
	proto.SendCreateRoomNotification(utility.ToString(roomId))
	logRoom.Debug("SendJoinRoomNotification pass")

	// get user info in the room
	userRoomNickname := orm.GetMemberName(utility.ToString(roomId), creator)
	var msg = userRoomNickname + "创建了群聊"
	SendAlert(creator, utility.ToString(roomId), types.ToRoom, nil, types.Alert, proto.ComposeAlert(types.AlertCreateRoom, creator, userRoomNickname, msg))
	//给被拉群的每个人发送消息通知
	go func() {
		for _, u := range rejectUsers {
			//发送通知 对方拒绝入群
			targetName, err := FriendRemarkOrName(creator, u)
			if err != nil {
				continue
			}
			msg := fmt.Sprintf("%s拒绝加入群聊", targetName)
			SendAlert(u, utility.ToString(roomId), types.ToRoom, []string{creator}, types.Alert, proto.ComposeRoomInviteReject(creator, u, targetName, msg))
		}
		for _, v := range members {
			if v == creator {
				continue
			}
			msg := fmt.Sprintf(" %s 已将 %s 拉入群聊[%s]", userRoomNickname, orm.GetMemberName(utility.ToString(roomId), v), roomName)
			SendAlert(creator, utility.ToString(v), types.ToUser, []string{creator, v}, types.Alert, proto.ComposeCreateRoomMemberAlert(types.AlertCreateRoom, creator, userRoomNickname, v, orm.GetMemberName(utility.ToString(roomId), v), utility.ToString(roomId), roomName, msg))
		}
	}()
	return &room, nil
}

func RemoveRoom(operator, roomId string) error {
	//get all member
	infos, _ := orm.FindNotDelMembers(roomId, types.SearchAll)
	var members = make([]string, 0)
	for _, info := range infos {
		if info.RoomMember.UserId != operator {
			members = append(members, info.RoomMember.UserId)
		}
	}

	rlt, err := orm.DelRoomById(roomId)
	if err != nil {
		return err
	}
	if rlt {
		//send remove room alert message
		msg := "该群已经被解散"
		SendAlert(operator, utility.ToString(roomId), types.ToRoom, nil, types.Alert, proto.ComposeAlert(types.AlertRemoveRoom, operator, "", msg))
		//send remove room notification to all member
		proto.SendRemoveRoomNotification(operator, roomId)
		//member unsubscribe
		RoomMemberUnSubscribe(roomId, append(members, operator))
		//delete room channel
		RemoveRoomChannel(roomId)
	}
	return nil
}

func LoginOutRoom(operator, roomId string) error {
	// get user info in the room
	userRoomNickname := orm.GetMemberName(roomId, operator)
	time := utility.NowMillionSecond()
	rlt, err := orm.DelRoomMemberById(operator, roomId, time)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if !rlt {
		logRoom.Warn("LoginOutRoom", "warn", "user is not in room")
		return result.NewError(result.AlreadyJoinOutRoom)
	}
	//current member unsubscribe
	RoomMemberUnSubscribe(roomId, []string{operator})
	//send logOut room notification to all member
	proto.SendLogOutRoomNotification(1, roomId, operator, "", []string{operator})

	var managers = make([]string, 0)
	managerInfos, _ := orm.GetRoomManagerAndMaster(roomId)
	for _, v := range managerInfos {
		userId := v.RoomMember.UserId
		if operator != userId {
			managers = append(managers, userId)
		}
	}
	managers = append(managers, operator)
	var msg = userRoomNickname + "退出群聊"
	SendAlert(operator, roomId, types.ToRoom, managers, types.Alert, proto.ComposeAlert(types.AlertLoginOutRoom, operator, userRoomNickname, msg))
	return nil
}

func KickOutRoom(caller, roomId string, userId string) error {
	// get user info in the room
	userRoomNickname := orm.GetMemberName(roomId, userId)
	time := utility.NowMillionSecond()
	rlt, err := orm.DelRoomMemberById(userId, roomId, time)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if !rlt {
		logRoom.Warn("KickOutRoom", "warn", "user is not in room")
		return result.NewError(result.AlreadyJoinOutRoom)
	}
	var msg = userRoomNickname + "被移出群聊"
	SendAlert(caller, roomId, types.ToRoom, nil, types.Alert, proto.ComposeKickOutAlert(types.AlertKickOutRoom, caller, userId, userRoomNickname, msg))
	//send logOut room notification to all member
	proto.SendLogOutRoomNotification(2, roomId, userId, "", []string{userId})
	//current member unsubscribe
	RoomMemberUnSubscribe(roomId, []string{userId})
	return nil
}

func roomInfoOfApiResult(queryUser, roomId string) (interface{}, error) {
	var room RoomBaseInfo
	roomInfo, err := orm.FindRoomById(roomId, types.RoomDeletedOrNot)
	if err != nil {
		return nil, err
	}
	if roomInfo == nil {
		logRoom.Warn("roomInfoOfApiResult", "warn", "room not exists", "roomId", roomId, "queryUser", queryUser)
		return nil, nil
	}
	room.Id = roomInfo.Id
	room.MarkId = roomInfo.MarkId
	room.Name = roomInfo.Name
	room.Avatar = roomInfo.Avatar
	room.Encrypt = roomInfo.Encrypt
	room.CanAddFriend = roomInfo.CanAddFriend
	room.JoinPermission = roomInfo.JoinPermission
	room.RecordPermission = roomInfo.RecordPermision
	room.DisableDeadline = roomInfo.CloseUntil
	room.MemberNumber, _ = orm.GetMemberNumber(roomId)
	room.ManagerNumber, _ = orm.GetRoomMasterNumber(roomId)
	room.SystemMsg, _ = GetSystemMsg(roomId, 0, 1)
	room.Identification = roomInfo.Identification
	room.IdentificationInfo = roomInfo.IdentificationInfo

	//如果群未被删除,则可查看和成员相关信息
	if roomInfo.IsDelete != types.RoomDeleted {
		memberInfo, err := orm.FindNotDelMember(room.Id, queryUser)
		if memberInfo != nil {
			var room = &RoomInfo{
				RoomBaseInfo: &room,
			}
			room.NoDisturbing = memberInfo.NoDisturbing
			room.OnTop = memberInfo.RoomTop
			room.MemberLevel = memberInfo.Level
			room.RoomNickname = memberInfo.UserNickname
			room.RoomMutedType, _ = orm.GetRoomMutedType(roomId)
			room.MutedNumber, _ = orm.GetMutedCount(roomId)
			room.MutedType, room.Deadline = orm.GetRoomUserMuted(roomId, queryUser)
			room.OnlineNumber = GetOnlineNumber(roomId)
			room.Users, err = getRoomMembers(roomId, 16)
			if err != nil {
				return nil, err
			}
			return room, nil
		}
	}
	return &room, nil
}

func GetRoomList(queryUser string, Type int) (interface{}, error) {
	infos, err := orm.GetJoinedRooms(queryUser, Type)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var roomList = make([]*RoomListInfo, 0)
	for _, info := range infos {
		var room RoomListInfo
		room.Id = info.Room.Id
		room.MarkId = info.MarkId
		room.Name = info.Name
		room.Avatar = info.Avatar
		room.DisableDeadline = info.CloseUntil
		room.NoDisturbing = info.NoDisturbing
		room.CommonlyUsed = info.CommonUse
		room.OnTop = info.RoomTop
		room.Encrypt = info.Encrypt
		room.Identification = info.Identification
		room.IdentificationInfo = info.IdentificationInfo
		roomList = append(roomList, &room)
	}
	var ret = make(map[string]interface{})
	ret["roomList"] = roomList
	return ret, nil
}

// 获取群信息
func GetRoomInfo(queryUser, roomId string) (interface{}, error) {
	info, err := roomInfoOfApiResult(queryUser, roomId)
	if err != nil {
		return info, result.NewError(result.DbConnectFail)
	}
	if info == nil {
		return nil, result.NewError(result.RoomNotExists)
	}
	return info, nil
}

func GetRoomUserList(roomId string) (interface{}, error) {
	users, err := getRoomMembers(roomId, -1)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var ret = make(map[string]interface{})
	ret["userList"] = users
	return ret, nil
}

func GetRoomUserInfo(roomId, userId string) (interface{}, error) {
	memberInfo, err := orm.FindNotDelMember(roomId, userId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var empty struct{}
	mutedType, err := orm.GetRoomMutedType(roomId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if memberInfo == nil {
		logRoom.Warn("GetRoomUserInfo", "warn", "user is not in room")
		return empty, result.NewError(result.UserIsNotInRoom)
	}
	var user RoomUserList
	user.Id = memberInfo.RoomMember.UserId
	user.PublicKey = memberInfo.PublicKey
	user.RoomNickname = memberInfo.UserNickname
	user.Nickname = memberInfo.Username
	user.MemberLevel = memberInfo.Level
	user.Avatar = memberInfo.Avatar
	user.RoomMutedType = mutedType
	user.MutedType, user.Deadline = orm.GetRoomUserMuted(roomId, user.Id)
	user.Source, _ = getJoinSource(memberInfo.AppId, roomId, memberInfo.Source)
	user.Source2, _ = getJoinSourceV2(memberInfo.AppId, roomId, memberInfo.Source)
	user.Identification = memberInfo.Identification
	user.IdentificationInfo = memberInfo.IdentificationInfo
	return &user, nil
}

func GetRoomSearchMember(roomId, name string) (interface{}, error) {
	infos, err := orm.FindMemberByName(roomId, name)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	mutedType, err := orm.GetRoomMutedType(roomId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var infoList = make([]*RoomUserList, 0)
	for _, info := range infos {
		var user RoomUserList
		user.Id = info.RoomMember.UserId
		user.RoomNickname = info.UserNickname
		user.Nickname = info.Username
		user.MemberLevel = info.Level
		user.Avatar = info.Avatar
		user.RoomMutedType = mutedType
		user.MutedType, user.Deadline = orm.GetRoomUserMuted(roomId, user.Id)
		user.Source, _ = getJoinSource(info.AppId, roomId, info.Source)
		user.Source2, _ = getJoinSourceV2(info.AppId, roomId, info.Source)
		user.Identification = info.Identification
		user.IdentificationInfo = info.IdentificationInfo
		infoList = append(infoList, &user)
	}
	var ret = make(map[string]interface{})
	ret["userList"] = infoList
	return &ret, nil
}

func AdminSetPermission(roomId string, canAddFriend, joinPermission, recordPermission int) error {
	err := orm.SetCanAddFriendPermission(roomId, canAddFriend)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	err = orm.SetJoinPermission(roomId, joinPermission)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	err = orm.SetRecordPermission(roomId, recordPermission)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func SetRoomName(operator, roomId, name string) error {
	isSuccess, err := orm.SetRoomName(roomId, name)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	if isSuccess {
		// get user info in the room
		userRoomNickname := orm.GetMemberName(roomId, operator)
		var msg = userRoomNickname + `将群名更改为"` + name + `"`
		//notify
		SendAlert(operator, roomId, types.ToRoom, nil, types.Alert, proto.ComposeSetRoomNameAlert(types.AlertRenameRoom, operator, userRoomNickname, name, msg))
	}
	return nil
}

func SetRoomAvatar(roomId, avatar string) error {
	_, err := orm.SetAvatar(roomId, avatar)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func SetLevel(master, userId, roomId string, level int) error {
	switch level {
	case types.RoomLevelNomal:
		fallthrough
	case types.RoomLevelManager:
		err := orm.SetMemberLevel(userId, roomId, level)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	case types.RoomLevelMaster:
		err := orm.SetNewMaster(master, userId, roomId, level)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
	default:
		return result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, fmt.Sprintf("unrecognized level %d", level))
	}

	// get user info in the room
	userRoomNickname := orm.GetMemberName(roomId, userId)
	var content map[string]interface{}
	var msg string
	if level == types.RoomLevelMaster {
		msg = userRoomNickname + "成为新的群主"
		content = proto.ComposeSetMemberLevelAlert(types.AlertSetAsMaster, master, userId, userRoomNickname, msg)
	}
	if level == types.RoomLevelManager {
		msg = userRoomNickname + "已被群主设为管理员"
		content = proto.ComposeSetMemberLevelAlert(types.AlertSetAsManager, master, userId, userRoomNickname, msg)
	}

	if msg != "" {
		SendAlert(master, roomId, types.ToRoom, nil, types.Alert, content)
	}
	return nil
}

func SetNoDisturbing(caller, roomId string, noDisturbing int) error {
	err := orm.SetNoDisturbing(caller, roomId, noDisturbing)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func SetStickyOnTop(caller, roomId string, onTop int) error {
	err := orm.SetOnTop(caller, roomId, onTop)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

func SetMemberNickname(caller, roomId string, nickname string) error {
	err := orm.SetMemberNickname(caller, roomId, nickname)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	return nil
}

//是否设置入群邀请需审核
func roomInviteCheck(appId, caller, userId string, roomInfo *types.Room) bool {
	//查看邀请配置
	inviteConfrim, err := orm.RoomInviteConfirm(userId)
	if err != nil {
		logRoom.Error("roomInviteCheck failed", "err", err)
	}
	needConfirmInvite := types.RoomInviteNotNeedConfirm
	if inviteConfrim == nil {
		//根据不同的app获取对应的 入群邀请配置
		if v, ok := types.DefaultInviteConfig[appId]; ok {
			needConfirmInvite = v
		}
	} else {
		needConfirmInvite = inviteConfrim.NeedConfirm
	}

	if needConfirmInvite == types.RoomInviteNotNeedConfirm {
		return true
	} else {
		return false
	}
}

//发送邀请名片
func sendInviteCard(caller, userId string, roomInfo *types.Room) {
	//发送邀请名片
	roomId := roomInfo.Id
	markId := roomInfo.MarkId
	roomName := roomInfo.Name
	avatar := roomInfo.Avatar
	identificationInfo := roomInfo.IdentificationInfo
	SendAlert(caller, userId, types.ToUser, []string{userId, caller}, types.RoomInvite, proto.ComposeRoomInviteCard(roomId, markId, roomName, caller, avatar, identificationInfo))
}

//bool:是否成功修改 error:错误
func JoinInRoomImiditly(tx types.Tx, caller, userId string, roomInfo *types.Room, source string) (bool, error) {
	roomId := roomInfo.Id
	createTime := utility.NowMillionSecond()
	//确保还不是群成员
	member, err := orm.FindNotDelMember(roomId, userId)
	if member != nil {
		return false, nil
	}

	err = orm.AddMember(tx, userId, roomId, types.RoomLevelNomal, createTime, source)
	if err != nil {
		return false, result.NewError(result.DbConnectFail)
	}
	return true, nil
}

//入群订阅
func joinInRoomSubscrib(userId string, roomId string) {
	//只能在当前同一协程，确保channel级别入群
	JoinInRoom(userId, roomId)
	go func() {
		//send notification to all member
		proto.SendJoinRoomNotification(roomId, []string{userId})
	}()
}

//入群 无需审批
func joinRoomNotApply(operator, roomId string, success bool, sourceType int, sourceId string) (bool, error) {
	//如果是加密群 不能同意没有公钥的人入群
	roomInfo, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err != nil {
		return success, result.NewError(result.DbConnectFail)
	}
	if roomInfo.Encrypt == types.IsEncrypt {
		user, err := orm.GetUserInfoById(operator)
		if err != nil {
			return success, result.NewError(result.DbConnectFail)
		}
		if user.PublicKey == "" {
			logRoom.Warn("joinRoomNotApply: user not have PublicKey", "warn", "CanNotJoinEncryptRoom", "PublicKey", user.PublicKey)
			return success, result.NewError(result.CanNotJoinEncryptRoom)
		}
	}
	//join in room Imiditly
	tx, err := orm.GetTx()
	if err != nil {
		tx.RollBack()
		return success, result.NewError(result.DbConnectFail)
	}
	_, err = JoinInRoomImiditly(tx, operator, operator, roomInfo, ConvertSourece(sourceType, sourceId))
	if err != nil {
		tx.RollBack()
		return success, err
	}
	err = tx.Commit()
	if err != nil {
		return success, result.NewError(result.DbConnectFail)
	}
	joinInRoomSubscrib(operator, roomId)
	success = true
	// get user info in the room
	//userRoomNickname := getUserName(userId)
	userRoomNickname := orm.GetMemberName(roomId, operator)
	msg := userRoomNickname + "加入群聊"
	SendAlert(operator, roomId, types.ToRoom, nil, types.Alert, proto.ComposeAlert(types.AlertJoinInRoom, operator, userRoomNickname, msg))
	return success, nil
}

// 返回： bool 是否直接加群成功
func JoinRoomApply(appId, operator, roomId, applyReason string, sourceType int, sourceId string) (bool, error) {
	//get room configuration
	success := false
	//判断群人数上限
	limitNumber, err := memberLimit(appId, &roomId)
	if err != nil {
		return success, err
	}
	currentNumber, _ := orm.GetMemberNumber(roomId)
	if utility.ToInt64(currentNumber) >= utility.ToInt64(limitNumber) && limitNumber != 0 {
		logRoom.Warn("JoinRoomApply", "warn", "MembersOutOfLimit", "current number", currentNumber, "limit number", limitNumber)
		return success, result.NewError(result.MembersOutOfLimit).SetExtMessage(ConvertRoomMembersLimit(limitNumber))
	}

	info, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err != nil {
		return success, result.NewError(result.DbConnectFail)
	}
	if info == nil {
		logRoom.Warn("JoinRoomApply", "warn", "RoomNotExists", "roomId", roomId)
		return success, result.NewError(result.RoomNotExists)
	}
	joinPermission := info.JoinPermission
	//check operator is room member
	if level := orm.GetMemberLevel(roomId, operator, types.RoomUserNotDeleted); level >= types.RoomLevelNomal {
		logRoom.Warn("JoinRoomApply", "warn", "IsRoomMemberAlready", "userId", operator, "roomId", roomId, "member level", level)
		return success, result.NewError(result.IsRoomMemberAlready)
	}

	switch joinPermission {
	case types.CanNotJoinRoom:
		inviterLevel := orm.GetMemberLevel(roomId, sourceId, types.RoomUserNotDeleted)
		if inviterLevel == types.RoomLevelMaster || inviterLevel == types.RoomLevelManager {
			return joinRoomNotApply(operator, roomId, success, sourceType, sourceId)
		}
		logRoom.Warn("JoinRoomApply: room setting is CanNotJoinRoom", "warn", "CanNotJoinRoom", "roomId", roomId)
		return success, result.NewError(result.CanNotJoinRoom)
	case types.ShouldApproval:
		logId, err := orm.AppendApplyLog(roomId, operator, applyReason, ConvertSourece(sourceType, sourceId), "", types.AwaitState, types.IsRoom)
		if err != nil {
			return success, result.NewError(result.DbConnectFail)
		}
		// get all manager and master
		managers, err := orm.GetRoomManagerAndMaster(roomId)
		if err == nil {
			var arry []string
			for _, v := range managers {
				userId := v.RoomMember.UserId
				if userId == operator {
					continue
				}
				arry = append(arry, userId)
			}
			applyInfo, err := GetApplyInfoByLogId(operator, logId)
			if err == nil {
				proto.SendApplyNotification(applyInfo, arry)
			}
		}
	case types.ShouldNotApproval:
		return joinRoomNotApply(operator, roomId, success, sourceType, sourceId)
	default:
		return success, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, fmt.Sprintf("unrecognized permission : %d", joinPermission))
	}
	return success, nil
}

//入群人数限制
func memberLimit(appId string, roomId *string) (int, error) {
	limitNumber := 0
	isVerified := false
	if roomId != nil {
		//判断是否加v
		room, err := orm.FindRoomById(*roomId, types.RoomNotDeleted)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if room == nil {
			logRoom.Warn("roomInfoOfApiResult", "warn", "room not exists", "roomId", *roomId, "appId", appId)
			return 0, result.NewError(result.RoomNotExists)
		}
		if room.Identification == types.Verified {
			limitNumber = types.VerifiedLimitMembers
			isVerified = true
		}
	}

	if !isVerified {
		//判断群人数是否超出
		numb, err := orm.GetRoomMembersLimit(appId, 1)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		limitNumber = numb
	}
	return limitNumber, nil
}

//返回state
func JoinRoomInvite(appId, operator, roomId string, users []string) (interface{}, error) {
	ret := make(map[string]interface{})
	ret["state"] = 1

	//get room configuration
	//判断群人数是否超出
	limitNumber, err := memberLimit(appId, &roomId)
	if err != nil {
		return nil, err
	}
	currentNumber, _ := orm.GetMemberNumber(roomId)
	if utility.ToInt64(currentNumber)+utility.ToInt64(len(users)) > utility.ToInt64(limitNumber) && limitNumber != 0 {
		logRoom.Warn("JoinRoomInvite", "warn", "MembersOutOfLimit", "current number", currentNumber, "invite number", len(users), "limit number", limitNumber)
		return nil, result.NewError(result.MembersOutOfLimit).SetExtMessage(ConvertRoomMembersLimit(limitNumber))
	}
	info, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err != nil {
		if err.Error() == "room not exisist" {
			logRoom.Warn("JoinRoomInvite", "warn", "RoomNotExists", "roomId", roomId)
			return nil, result.NewError(result.RoomNotExists)
		}
		return nil, result.NewError(result.DbConnectFail)
	}

	//如果是加密群 不能邀请没有公钥的人
	if info.Encrypt == types.IsEncrypt {
		old := users
		users = make([]string, 0)
		for _, u := range old {
			user, err := orm.GetUserInfoById(u)
			if err != nil {
				continue
			}
			if user.PublicKey != "" {
				users = append(users, u)
			}
		}
	}
	//如果对方将你拉入黑名单，则无法邀请
	old := make([]string, 0)
	old = append(old, users...)
	users = make([]string, 0)
	for _, u := range old {
		ok, err := CheckIsBlocked(u, operator)
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		if ok {
			//发送通知 对方拒绝入群
			targetName, err := FriendRemarkOrName(operator, u)
			if err != nil {
				return nil, result.NewError(result.DbConnectFail)
			}
			msg := fmt.Sprintf("%s拒绝加入群聊", targetName)
			SendAlert(u, roomId, types.ToRoom, []string{operator}, types.Alert, proto.ComposeRoomInviteReject(operator, u, targetName, msg))
			continue
		}
		users = append(users, u)
	}

	sendInvite := make([]string, 0)
	newJoin := make([]string, 0)
	joinPermission := info.JoinPermission
	roomName := info.Name
	switch joinPermission {
	case types.CanNotJoinRoom:
		//check operator is master or admin
		level := orm.GetMemberLevel(roomId, operator, types.RoomUserNotDeleted)
		if level != types.RoomLevelMaster && level != types.RoomLevelManager {
			logRoom.Warn("JoinRoomInvite", "warn", "PermissionDeny", "operator level", level, "must level", fmt.Sprintf("%d or %d", types.RoomLevelMaster, types.RoomLevelManager))
			return nil, result.NewError(result.PermissionDeny)
		}
		//join in room Imiditly
		tx, err := orm.GetTx()
		if err != nil {
			tx.RollBack()
			return nil, result.NewError(result.DbConnectFail)
		}
		source := ConvertSourece(types.Invite, operator)
		for _, userId := range users {
			//如果
			if roomInviteCheck(appId, operator, userId, info) {
				b, err := JoinInRoomImiditly(tx, operator, userId, info, source)
				if err != nil {
					tx.RollBack()
					return nil, err
				}
				if b {
					newJoin = append(newJoin, userId)
				}
			} else {
				sendInvite = append(sendInvite, userId)
			}
		}
		err = tx.Commit()
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}
		//发送邀请卡片
		for _, u := range sendInvite {
			sendInviteCard(operator, u, info)
		}
		//订阅channel
		for _, u := range newJoin {
			joinInRoomSubscrib(u, roomId)
		}

	case types.ShouldApproval:
		//如果不是管理员
		if orm.GetMemberLevel(roomId, operator, types.RoomUserNotDeleted) < types.RoomLevelManager {
			for _, v := range users {
				_, err := JoinRoomApply(appId, v, roomId, "", types.Invite, operator)
				if err != nil {
					logRoom.Warn("JoinRoomApply Warn", "warn", err)
				}
			}
			ret["state"] = 0
			return ret, nil
		}
		//join in room Imiditly
		tx, err := orm.GetTx()
		if err != nil {
			tx.RollBack()
			return nil, result.NewError(result.DbConnectFail)
		}
		source := ConvertSourece(types.Invite, operator)
		for _, userId := range users {
			if roomInviteCheck(appId, operator, userId, info) {
				b, err := JoinInRoomImiditly(tx, operator, userId, info, source)
				if err != nil {
					tx.RollBack()
					return nil, err
				}
				if b {
					newJoin = append(newJoin, userId)
				}
			} else {
				sendInvite = append(sendInvite, userId)
			}
		}
		err = tx.Commit()
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}

		//发送邀请卡片
		for _, u := range sendInvite {
			sendInviteCard(operator, u, info)
		}
		//订阅channel
		for _, u := range newJoin {
			joinInRoomSubscrib(u, roomId)
		}
	case types.ShouldNotApproval:
		tx, err := orm.GetTx()
		if err != nil {
			tx.RollBack()
			return nil, result.NewError(result.DbConnectFail)
		}
		source := ConvertSourece(types.Invite, operator)
		for _, userId := range users {
			if roomInviteCheck(appId, operator, userId, info) {
				b, err := JoinInRoomImiditly(tx, operator, userId, info, source)
				if err != nil {
					tx.RollBack()
					return nil, err
				}
				if b {
					newJoin = append(newJoin, userId)
				}
			} else {
				sendInvite = append(sendInvite, userId)
			}
		}
		err = tx.Commit()
		if err != nil {
			return nil, result.NewError(result.DbConnectFail)
		}

		//发送邀请卡片
		for _, u := range sendInvite {
			sendInviteCard(operator, u, info)
		}
		//订阅channel
		for _, u := range newJoin {
			joinInRoomSubscrib(u, roomId)
		}
	default:
		logRoom.Warn("JoinRoomInvite", "warn", "room JoinPermission type is unrecognized")
		return nil, result.NewError(result.ParamsError)
	}
	//给被邀请的每个人发送消息通知
	go func() {
		userRoomNickname := orm.GetMemberName(utility.ToString(roomId), operator)
		for _, v := range newJoin {
			msg := fmt.Sprintf(" %s 已将 %s 拉入群聊[%s]", userRoomNickname, orm.GetMemberName(utility.ToString(roomId), v), roomName)
			SendAlert(operator, utility.ToString(v), types.ToUser, []string{operator, v}, types.Alert, proto.ComposeCreateRoomMemberAlert(types.AlertCreateRoom, operator, userRoomNickname, v, orm.GetMemberName(utility.ToString(roomId), v), utility.ToString(roomId), roomName, msg))
		}
	}()

	//fix:空邀请时提示入群消息去掉 dld-v2.8.4 2019年7月2日16:39:08
	if len(newJoin) > 0 {
		//发送消息通知
		names, operName, msg := ConvertInviteAlertContent(roomId, operator, newJoin)
		SendAlert(operator, roomId, types.ToRoom, nil, types.Alert, proto.ComposeInviteRoomAlert(types.AlertInviteJoinRoom, operator, operName, newJoin, names, msg))
	}
	return ret, nil
}

//已经入群：如果已经入群，修改入群请求为接受状态。如果管理员点击的是拒绝则返回已经入群错误；如果点击的是同意则成功。
func JoinRoomApprove(appId, operator, roomId, userId string, aggre int) error {
	tx, err := orm.GetTx()
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	//获取入群请求信息
	info, err := orm.GetMemberApplyInfo(roomId, userId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}

	//是否早已加入群
	joined := false
	switch utility.ToInt(info.State) {
	case types.AcceptState:
		return result.NewError(result.ApplyAlreadyDeal)
	case types.RejectState:
		return result.NewError(result.ApplyAlreadyDeal)
	case types.AwaitState:
		//check userId is room member
		if orm.GetMemberLevel(roomId, userId, types.RoomUserNotDeleted) >= types.RoomLevelNomal {
			joined = true
		}
	default:
		logRoom.Warn("JoinRoomApprove", "warn", "unrecognized state type", "roomId", roomId, "userId", userId, "state", info.State)
		return result.NewError(result.ParamsError)
	}
	applyId := utility.ToInt64(info.Id)
	source := info.Source
	sObj := utility.StringToJobj(source)
	sType := sObj["sourceType"]
	sId := sObj["sourceId"]
	var status int
	isOk := false
	if joined {
		status = types.AcceptState
	} else if aggre == types.AcceptRequest {
		//判断群人数上限
		limitNumber, err := memberLimit(appId, &roomId)
		if err != nil {
			return err
		}
		currentNumber, _ := orm.GetMemberNumber(roomId)
		if utility.ToInt64(currentNumber) >= utility.ToInt64(limitNumber) && limitNumber != 0 {
			logRoom.Warn("JoinRoomApprove", "warn", "MembersOutOfLimit", "current number", currentNumber, "limit number", limitNumber)
			return result.NewError(result.MembersOutOfLimit).SetExtMessage(ConvertRoomMembersLimit(limitNumber))
		}

		//如果是加密群 不能同意没有公钥的人入群
		roomInfo, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
		if err != nil {
			return result.NewError(result.DbConnectFail)
		}
		if roomInfo.Encrypt == types.IsEncrypt {
			user, err := orm.GetUserInfoById(userId)
			if err != nil {
				return result.NewError(result.DbConnectFail)
			}
			if user.PublicKey == "" {
				logRoom.Warn("JoinRoomApprove", "warn", "CanNotJoinEncryptRoom:public key is empty", "userId", userId, "publicKey", user.PublicKey)
				return result.NewError(result.CanNotJoinEncryptRoom)
			}
		}

		isSuccess, err := orm.ApproveInsertMemberStep(tx, roomId, userId, source)
		if err != nil {
			tx.RollBack()
			return result.NewError(result.DbConnectFail)
		}
		isOk = isSuccess
		JoinInRoom(userId, roomId)
		//send notification to all member
		proto.SendJoinRoomNotification(roomId, []string{userId})
		status = types.AcceptState
	} else {
		status = types.RejectState
	}
	//修改请求信息状态
	_, err = orm.ApproveChangeStateStep(tx, applyId, status)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	err = tx.Commit()
	if err != nil {
		logRoom.Error("JoinRoomApprove tx commit err", "err", err)
		return result.NewError(result.DbConnectFail)
	}
	if joined && aggre == types.RejectRequest {
		return result.NewError(result.IsRoomMemberAlready)
	}
	if isOk {
		// get all manager and master 给所有管理员和群主发送事件通知
		managers, err := orm.GetRoomManagerAndMaster(roomId)
		if err == nil {
			var array []string
			for _, v := range managers {
				userId := v.RoomMember.UserId
				array = append(array, userId)
			}
			array = append(array, userId)
			applyInfo, err := GetApplyInfoByLogId(operator, applyId)
			if err == nil {
				proto.SendApplyNotification(applyInfo, array)
			}
			//SendApplyNotification(operator, applyId, array)
		}

		switch utility.ToInt(sType) {
		case types.Invite:
			//根据邀请者uid获取id
			inviterInfo, _ := orm.GetUserInfoByUid(appId, utility.ToString(sId))
			if inviterInfo != nil {
				inviterId := utility.ToString(inviterInfo.UserId)
				users := []string{userId}
				names, operName, msg := ConvertInviteAlertContent(roomId, inviterId, users)
				//给被邀请的每个人发送消息通知
				go func() {
					roomInfo, _ := orm.FindRoomById(roomId, types.RoomNotDeleted)
					if roomInfo != nil {
						roomName := roomInfo.Name
						userRoomNickname := orm.GetMemberName(utility.ToString(roomId), inviterId)
						targetName := orm.GetMemberName(utility.ToString(roomId), userId)
						msg := fmt.Sprintf(" %s 已将 %s 拉入群聊[%s]", userRoomNickname, targetName, roomName)
						SendAlert(inviterId, utility.ToString(userId), types.ToUser, []string{inviterId, userId}, types.Alert, proto.ComposeCreateRoomMemberAlert(types.AlertCreateRoom, inviterId, userRoomNickname, userId, targetName, utility.ToString(roomId), roomName, msg))
					}
				}()
				SendAlert(operator, roomId, types.ToRoom, nil, types.Alert, proto.ComposeInviteRoomAlert(types.AlertInviteJoinRoom, operator, operName, users, names, msg))
			}
		default:
			// get user info in the room
			userRoomNickname := orm.GetMemberName(roomId, userId)
			msg := userRoomNickname + "加入群聊"
			SendAlert(operator, roomId, types.ToRoom, nil, types.Alert, proto.ComposeAlert(types.AlertJoinInRoom, userId, userRoomNickname, msg))
		}
	}
	return nil
}

func ApplyRoomInfo(roomId string) (*types.ApplyInfo, error) {
	var targetInfo types.ApplyInfo
	room, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if room == nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	targetInfo.Id = room.Id
	targetInfo.MarkId = room.MarkId
	targetInfo.Name = room.Name
	targetInfo.Avatar = room.Avatar
	targetInfo.Position = ""
	targetInfo.Identification = room.Identification
	targetInfo.IdentificationInfo = room.IdentificationInfo
	return &targetInfo, nil
}

//startTime:是入群时间或者是0（表示可随意查询历史记录）
func msgLogs(list []*types.ChatLog, callerId, roomId string, number int, startId, startTime int64) ([]*types.ChatLog, int64, error) {
	nextLog, logs, err := orm.FindRoomLogs(roomId, startId, startTime, number-len(list))
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	if logs == nil {
		return nil, -1, nil
	}

	for _, log := range logs {
		info, err := GetChatLogAsRoom(callerId, log)
		if err != nil {
			logRoom.Warn("GetChatLogAsRoom error", "log info", info)
			continue
		}

		if info.IsSnap == types.IsSnap {
			continue
		}
		if info.MsgType == types.Alert {
			ck := proto.NewMsg(info.GetCon())
			//过滤掉赞赏通知
			if ck.GetTypeCode() == types.AlertPraise {
				continue
			}
			if ck.CheckSpecialMsg() {
				revLog, err := orm.FindReceiveLogById(info.LogId, callerId)
				if err != nil {
					return nil, -1, result.NewError(result.DbConnectFail)
				}
				if revLog == nil {
					continue
				}
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
		return msgLogs(list, callerId, roomId, number, utility.ToInt64(nextLog), startTime)
	}
	return list, utility.ToInt64(nextLog), nil
}

func GetRoomChatLog(callerId, roomId, startId string, number int) (interface{}, error) {
	if startId != "" && number < 1 {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "number less than 1")
	}
	startLogId := utility.ToInt64(startId)
	if startLogId == 0 && number <= 0 {
		number = 20
	}

	//获取群配置
	roomInfo, err := orm.FindRoomById(roomId, types.RoomDeletedOrNot)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	if roomInfo == nil {
		logRoom.Warn("GetRoomChatLog", "warn", "RoomNotExists", "roomId", roomId)
		return nil, result.NewError(result.RoomNotExists)
	}
	rp := roomInfo.RecordPermision
	var joinTime int64
	if rp == types.CanReadAllLog {
		joinTime = 0
	} else if rp == types.CanNotReadAllLog {
		joinTime = getMemberJoinTime(roomId, callerId)
	}

	list, nextLog, err := msgLogs([]*types.ChatLog{}, callerId, roomId, number, startLogId, joinTime)
	if list == nil {
		list = []*types.ChatLog{}
	}

	var ret = make(map[string]interface{})
	ret["logs"] = list
	ret["nextLog"] = utility.ToString(nextLog)
	return ret, nil
}

//
func typicalMsgLogs(list []*types.ChatLog, callerId, roomId, owner, query string, number int, queryType []string, startId int64) ([]*types.ChatLog, int64, error) {
	nextLog, logs, err := orm.FindRoomLogsByUserId(roomId, owner, startId, 0, number-len(list), queryType)
	if err != nil {
		return nil, -1, result.NewError(result.DbConnectFail)
	}
	if logs == nil {
		return nil, -1, nil
	}

	for _, log := range logs {
		info, err := GetChatLogAsRoom(callerId, log)
		if err != nil {
			logRoom.Warn("GetChatLogAsRoom error", "log info", info)
			continue
		}

		if info.IsSnap == types.IsSnap {
			continue
		}
		if info.MsgType == types.Alert {
			ck := proto.NewMsg(info.GetCon())
			if ck.CheckSpecialMsg() {
				revLog, err := orm.FindReceiveLogById(info.LogId, callerId)
				if err != nil {
					return nil, -1, result.NewError(result.DbConnectFail)
				}
				if revLog == nil {
					continue
				}
			}
		}
		isAppend := true
		//根据query筛选文件名,上传者群名/备注/昵称
		if query != "" {
			isAppend = false
			if !isAppend {
				//昵称
				senderInfo := info.SenderInfo.(map[string]interface{})
				username := utility.ToString(senderInfo["nickname"])
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
				if nickname == "" {
					continue
				}
				if strings.Contains(nickname, query) {
					isAppend = true
				}
			}
			if !isAppend && callerId != info.FromId {
				//不匹配 查询备注
				friendInfo, _ := orm.FindFriendById(callerId, info.FromId)
				if friendInfo == nil {
					continue
				}
				remark := friendInfo.Remark
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
		return typicalMsgLogs(list, callerId, roomId, owner, query, number, queryType, utility.ToInt64(nextLog))
	}
	return list, utility.ToInt64(nextLog), nil
}

//获取指定消息类型的聊天记录
func GetTypicalMsgLogs(callerId, roomId, startId, owner, query string, number int, queryType []string) (interface{}, error) {
	if startId != "" && number < 1 {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "number less than 1")
	}
	startLogId := utility.ToInt64(startId)
	if startLogId == 0 && number <= 0 {
		number = 20
	}

	list := make([]*types.ChatLog, 0)
	list, nextLog, err := typicalMsgLogs(list, callerId, roomId, owner, query, number, queryType, startLogId)

	var ret = make(map[string]interface{})
	ret["logs"] = list
	ret["nextLog"] = utility.ToString(nextLog)
	return ret, err
}

// 获取群公告
func GetSystemMsg(roomId string, startId int64, number int) (interface{}, error) {
	if startId != 0 && number < 1 {
		return nil, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "number less than 1")
	}
	if startId == 0 && number == 0 {
		number = 20
	}

	nextId, logs, err := orm.FindSystemMsg(roomId, startId, number)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}

	var list = make([]*SystemMessage, 0)
	for _, r := range logs {
		var info SystemMessage
		info.Id = r.Id
		info.SenderName = orm.GetMemberName(roomId, r.SenderId)
		info.Content = utility.ToString(utility.StringToJobj(r.Content)["content"])
		info.Datetime = r.Datetime
		list = append(list, &info)
	}

	totalNumber, err := orm.GetSystemMsgNumber(roomId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	var ret = make(map[string]interface{})
	ret["number"] = totalNumber
	ret["list"] = list
	ret["nextLog"] = utility.ToString(nextId)
	return ret, nil
}

// 发布群公告
func SetSystemMsg(caller, roomId, content string) error {
	SendAlert(caller, roomId, types.ToRoom, nil, types.System, proto.ComposeSystemMsg(content))
	return nil
}

func GetRoomOnlineNumber(roomId string) (interface{}, error) {
	var ret = make(map[string]int)
	ret["onlineNumber"] = GetOnlineNumber(roomId)
	return &ret, nil
}

//返回 content 和用户名称列表、群内等级
func getMutedAlertMsg(caller string, mutedType int, roomId string, users []string, deadline int64) (int, []string, string) {
	names := make([]string, 0)
	var msg string
	var callerName string
	level := orm.GetMemberLevel(roomId, caller, types.RoomUserNotDeleted)
	if level == types.RoomLevelMaster {
		callerName = "群主"
	} else if level == types.RoomLevelManager {
		callerName = "管理员"
	}
	switch mutedType {
	case types.AllSpeak:
		if callerName != "" {
			msg = callerName + " 设置全员可发言"
		}
	case types.AllMuted:
		if callerName != "" {
			msg = callerName + " 设置全员禁言"
		}
	case types.Blacklist:
		fallthrough
	case types.Whitelist:
		for i, v := range users {
			//more than seven
			if i > 6 {
				break
			}
			memName := orm.GetMemberName(roomId, v)
			if i < 6 {
				if i != 0 {
					msg += "、"
				}
				msg += memName
			}
			names = append(names, memName)
		}
		if len(msg) > 0 {
			if mutedType == types.Blacklist {
				if deadline != 0 {
					if len(users) > 6 {
						msg += " 等" + utility.ToString(len(users)) + "人已禁言"
					} else {
						msg += " 已禁言"
					}
				} else {
					if len(users) > 6 {
						msg += " 等 " + utility.ToString(len(users)) + "人可发言"
					} else {
						msg += " 可发言"
					}
				}
			} else {
				if len(users) > 6 {
					msg += " 等 " + utility.ToString(len(users)) + "人可发言"
				} else {
					msg += " 可发言"
				}
			}
		}
	}
	return level, names, msg
}

//设置成员禁言
func SetRoomMuted(operator, roomId string, mutedType int, userMap map[string]bool, deadline int64) error {
	oldType, _ := orm.GetRoomMutedType(roomId)
	tx, err := orm.GetTx()
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	//设置群中所采用的发言方式
	err = orm.SetRoomMutedType(tx, roomId, mutedType)
	if err != nil {
		tx.RollBack()
		return result.NewError(result.DbConnectFail)
	}
	var users []string
	for userId := range userMap {
		users = append(users, userId)
	}
	switch mutedType {
	case types.Blacklist:
		fallthrough
	case types.Whitelist:
		var oldWhiteList []*types.RoomUserMuted
		if mutedType == types.Whitelist {
			//获取原先白名单成员,清空白名单
			oldWhiteList, err = orm.GetMutedListByType(roomId, types.Whitelist)
			if err != nil {
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
			err := orm.ClearMutedList(tx, roomId)
			if err != nil {
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
		}
		//添加黑、白名单成员
		for userId := range userMap {
			_, err := orm.AddMutedMember(tx, roomId, userId, mutedType, deadline)
			if err != nil {
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
		}
		go func() {
			if mutedType == types.Blacklist {
				//发送【禁言通知】
				proto.SendRoomMutedNotification(roomId, types.MutedEnable, deadline, users)
				//全部禁言、白名单——>黑名单 解禁非黑名单成员
				if oldType == types.AllMuted || oldType == types.Whitelist {
					var members []string
					users := normalMembersId(roomId)
					for _, v := range users {
						if _, ok := userMap[v]; !ok {
							members = append(members, v)
						}
					}
					proto.SendRoomMutedNotification(roomId, types.MutedDisable, 0, members)
				}
			}

			if mutedType == types.Whitelist {
				//发送【解禁通知】
				proto.SendRoomMutedNotification(roomId, types.MutedDisable, 0, users)
				//全员发言、黑名单——>白名单 则【禁言】非白名单成员
				if oldType == types.AllSpeak || oldType == types.Blacklist {
					users := normalMembersId(roomId)
					var members []string
					for _, v := range users {
						if _, ok := userMap[v]; !ok {
							members = append(members, v)
						}
					}
					proto.SendRoomMutedNotification(roomId, types.MutedEnable, types.MutedForvevr, members)
				} else {
					//【禁言】之前的白名单成员（不在现有名单中）
					var members []string
					for _, v := range oldWhiteList {
						if _, ok := userMap[v.UserId]; !ok {
							members = append(members, v.UserId)
						}
					}
					proto.SendRoomMutedNotification(roomId, types.MutedEnable, types.MutedForvevr, members)
				}
			}
		}()
	case types.AllMuted:
		users := normalMembersId(roomId)
		proto.SendRoomMutedNotification(roomId, types.MutedEnable, types.MutedForvevr, users)
	case types.AllSpeak:
		users := normalMembersId(roomId)
		proto.SendRoomMutedNotification(roomId, types.MutedDisable, 0, users)
	}
	err = tx.Commit()
	if err != nil {
		logRoom.Error("SetRoomMuted tx commit err", "err", err)
		return result.NewError(result.DbConnectFail)
	}
	level, names, msg := getMutedAlertMsg(operator, mutedType, roomId, users, deadline)
	opt := 1
	switch deadline {
	case 0:
		opt = 2
	default:
		opt = 1
	}
	SendAlert(operator, utility.ToString(roomId), types.ToRoom, nil, types.Alert, proto.ComposeRoomMutedAlert(operator, mutedType, level, opt, names, msg))
	return nil
}

func SetRoomMutedSingle(operator, roomId, userId string, deadline int64) error {
	//获取群采用的发言类型
	mutedType, err := orm.GetRoomMutedType(roomId)
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	tx, err := orm.GetTx()
	if err != nil {
		return result.NewError(result.DbConnectFail)
	}
	mType := mutedType
	//解除禁言
	if deadline == 0 {
		switch mutedType {
		case types.AllSpeak:
			return nil
		case types.Blacklist:
			//remove from blacklist
			err := orm.DelMemberMuted(tx, roomId, userId)
			if err != nil {
				tx.RollBack()
				logRoom.Error("set blacklist DelRoomUserMuted err", "err", err)
				return result.NewError(result.DbConnectFail)
			}
			//如果人数为0,设为全员发言
			number, err := orm.GetMutedCountTx(tx, roomId)
			if err == nil && number == 0 {
				err := orm.SetRoomMutedType(tx, roomId, types.AllSpeak)
				if err != nil {
					logRoom.Error("single set blacklist err", "err", err)
					tx.RollBack()
					return result.NewError(result.DbConnectFail)
				}
				mType = types.AllSpeak
			}
		case types.AllMuted:
			//设置群中所采用的发言方式为【白名单】模式
			err := orm.SetRoomMutedType(tx, roomId, types.Whitelist)
			if err != nil {
				logRoom.Error("single set all muted err", "err", err)
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
			mType = types.Whitelist
			fallthrough
		case types.Whitelist:
			//append user
			_, err := orm.AddMutedMember(tx, roomId, userId, types.Whitelist, deadline)
			if err != nil {
				logRoom.Error("single set whitelist err", "err", err)
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
		default:
			return nil
		}
	} else { //设置禁言时长
		switch mutedType {
		case types.AllSpeak:
			//设置群中所采用的发言方式为【黑名单】模式
			err := orm.SetRoomMutedType(tx, roomId, types.Blacklist)
			if err != nil {
				logRoom.Error("single set allspeak err", "err", err)
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
			mType = types.Blacklist
			fallthrough
		case types.Blacklist:
			//append user
			_, err := orm.AddMutedMember(tx, roomId, userId, types.Blacklist, deadline)
			if err != nil {
				logRoom.Error("single set blacklist err", "err", err)
				tx.RollBack()
				return result.NewError(result.DbConnectFail)
			}
		case types.AllMuted:
			return nil
		case types.Whitelist:
			//remove from whitelist
			err := orm.DelMemberMuted(tx, roomId, userId)
			if err != nil {
				tx.RollBack()
				logRoom.Error("set blacklist DelRoomUserMuted err", "err", err)
				return result.NewError(result.DbConnectFail)
			}
			//如果人数为0,设为全员禁言
			number, err := orm.GetMutedCountTx(tx, roomId)
			if err == nil && number == 0 {
				err := orm.SetRoomMutedType(tx, roomId, types.AllMuted)
				if err != nil {
					logRoom.Error("single set whitelist err", "err", err)
					tx.RollBack()
					return result.NewError(result.DbConnectFail)
				}
				mType = types.AllMuted
			}
		default:
			return nil
		}
	}
	if deadline == 0 {
		proto.SendRoomMutedNotification(roomId, types.MutedDisable, deadline, []string{userId})
	} else {
		proto.SendRoomMutedNotification(roomId, types.MutedEnable, deadline, []string{userId})
	}
	err = tx.Commit()
	if err != nil {
		logRoom.Error("SetRoomMutedSingle tx commit err", "err", err)
		return result.NewError(result.DbConnectFail)
	}
	level, names, msg := getMutedAlertMsg(operator, mType, roomId, []string{userId}, deadline)
	opt := 1
	switch deadline {
	case 0:
		opt = 2
	default:
		opt = 1
	}
	SendAlert(operator, utility.ToString(roomId), types.ToRoom, nil, types.Alert, proto.ComposeRoomMutedAlert(operator, mType, level, opt, names, msg))
	return nil
}

//判断是否在群内
func UserIsInRoom(userId, roomId string) (map[string]bool, error) {
	var isInRoom = make(map[string]bool)
	ret, err := CheckUserInRoom(userId, roomId)
	if err != nil {
		return nil, result.NewError(result.DbConnectFail)
	}
	isInRoom["isInRoom"] = ret
	return isInRoom, nil
}

//获取推荐群
func GetRecommendRooms(appId string, number int, times int) (interface{}, error) {
	re := GetRecommendRoom(appId)

	nextTimes := times + 1
	times = times - 1

	r1 := re.GetCopyData()
	r2 := re.GetCopyDataMem()
	r3 := re.GetCopyDataMsg()

	ck := make(map[string]struct{})
	rooms := make([]*types.Room, 0)

	if number*(times+1) < len(r1) {
		rooms = r1[times*number : (times+1)*number]
	} else {
		if number*times < len(r1) {
			rooms = r1[times*number:]
			for _, r := range rooms {
				ck[r.Id] = struct{}{}
			}
		}
		num := number - len(rooms)
		r23 := append(r2, r3...)
		//数据从r2和r3中随机选择
		for i := 0; i < num; {
			index := utility.RandInt(0, len(r23)-1)
			if index < 0 {
				break
			}

			r := r23[index]
			r23 = append(r23[:index], r23[index+1:]...)

			if _, ok := ck[r.Id]; !ok {
				i++
				rooms = append(rooms, r)
				ck[r.Id] = struct{}{}
			}
		}
	}
	ret := make(map[string]interface{})
	ret["roomList"] = rooms
	ret["nextTimes"] = nextTimes
	return ret, nil
}

//检查是否已经认证
func CheckRoomVerifyState(userId, roomId string) (int, error) {
	room, err := orm.FindRoomById(roomId, types.RoomNotDeleted)
	if err != nil {
		return 0, result.NewError(result.DbConnectFail)
	}
	if room.Identification == types.Verified {
		return types.VerifyStateAccept, nil
	} else {
		wait, err := orm.FindVerifyApplyByState(types.VerifyForRoom, roomId, types.VerifyStateWait)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if len(wait) > 0 {
			return types.VerifyStateWait, nil
		}

		reject, err := orm.FindVerifyApplyByState(types.VerifyForRoom, roomId, types.VerifyStateReject)
		if err != nil {
			return 0, result.NewError(result.DbConnectFail)
		}
		if len(reject) > 0 {
			rej := reject[0]
			if utility.NowMillionSecond() < utility.MillionSecondAddDuration(rej.UpdateTime, types.VerifyApplyInterval*time.Minute) {
				return types.VerifyStateReject, nil
			}
			return 0, nil
		}
	}
	return 0, nil
}
