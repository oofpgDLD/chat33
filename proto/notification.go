package proto

import (
	"encoding/json"

	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

/*
	这个文件是关于事件通知的
*/

// 被封禁群通知
type RemoveClosedNotification struct {
	EventType       int    `json:"eventType"`
	RoomId          string `json:"roomId"`
	Content         string `json:"content"`
	DisableDeadline int64  `json:"disableDeadline"`
	Datetime        int64  `json:"datetime"`
}

// 被拉入群通知
type JoinRoomNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	Datetime  int64  `json:"datetime"`
}

// 被解散群通知
type RemoveRoomNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	UserId    string `json:"userId"`
	Datetime  int64  `json:"datetime"`
}

// 退群通知
type LogOutRoomNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	UserId    string `json:"userId"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
}

// 群在线人数通知
type RoomOnlineNumberNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	Number    string `json:"number"`
	Datetime  int64  `json:"datetime"`
}

// 群禁言通知
type RoomMutedNotification struct {
	EventType int    `json:"eventType"`
	RoomId    string `json:"roomId"`
	Type      int    `json:"type"`
	Deadline  int64  `json:"deadline"`
}

// 发送建群通知
func SendCreateRoomNotification(roomId string) {
	var ret JoinRoomNotification
	ret.EventType = types.EventJoinRoom
	ret.RoomId = roomId
	ret.Datetime = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	roomChannelId := types.GetRoomRouteById(roomId)
	//发送给所有人
	if cl, _ := router.GetChannel(roomChannelId); cl != nil {
		cl.Broadcast(data)
	}
}

// 发送入群通知
func SendJoinRoomNotification(roomId string, members []string) {
	var ret JoinRoomNotification
	ret.EventType = types.EventJoinRoom
	ret.RoomId = roomId
	ret.Datetime = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client, _ := router.GetUser(memId)
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 发送解散群通知
func SendRemoveRoomNotification(operator, roomId string) {
	var ret RemoveRoomNotification
	ret.EventType = types.EventRemoveRoom
	ret.RoomId = roomId
	ret.UserId = operator
	ret.Datetime = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	roomChannelId := types.GetRoomRouteById(roomId)
	//发送给所有人
	if cl, _ := router.GetChannel(roomChannelId); cl != nil {
		cl.Broadcast(data)
	}
}

// 发送退出群通知
func SendLogOutRoomNotification(logOutType int, roomId, userId string, content string, members []string) {
	var ret LogOutRoomNotification
	ret.EventType = types.EventLogOutRoom
	ret.RoomId = roomId
	ret.UserId = userId
	ret.Type = logOutType
	ret.Content = content
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client, _ := router.GetUser(memId)
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 发送群中被禁言通知
func SendRoomMutedNotification(roomId string, mutedType int, deadline int64, members []string) {
	var ret RoomMutedNotification
	ret.EventType = types.EventRoomMuted
	ret.RoomId = roomId
	ret.Type = mutedType
	ret.Deadline = deadline
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client, _ := router.GetUser(memId)
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 发送入群请求和回复通知
func SendApplyNotification(applyInfo *types.ApplyList, members []string) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventJoinApply
	if applyInfo.Type == types.IsRoom {
		ret["senderInfo"] = applyInfo.SenderInfo
		ret["receiveInfo"] = applyInfo.ReceiveInfo
		ret["id"] = applyInfo.Id
		ret["type"] = applyInfo.Type
		ret["applyReason"] = applyInfo.ApplyReason
		ret["status"] = applyInfo.Status
		ret["source"] = applyInfo.Source
		ret["source2"] = applyInfo.Source2
		ret["datetime"] = applyInfo.Datetime
	}
	data, _ := json.Marshal(ret)

	for _, memId := range members {
		client, _ := router.GetUser(memId)
		if client != nil {
			client.SendToAllClients(data)
		}
	}
}

// 其他设备登录
func SendOtherDeviceLogin(alertMsg string, client *router.Client) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventOtherDeviceLogin
	ret["content"] = alertMsg
	ret["datetime"] = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	client.Send(data)
}

// 同步完成通知
func SendBatchSuccessNotification(client *router.Client) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventSyncMsgRlt
	data, _ := json.Marshal(ret)

	client.Send(data)
}

// ack完成通知
func SendAckSuccessNotification(client *router.Client, begin, end int64) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventAckSuccess
	ret["begin"] = begin
	ret["end"] = end
	data, _ := json.Marshal(ret)

	client.Send(data)
}

//群会话秘钥同步完成通知
func SendGetSKeySuccessNotification(client *router.Client) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventGetAllSKeySuccess
	data, _ := json.Marshal(ret)

	client.Send(data)
}

// 发送封群通知
func SendClosedRoomNotification(roomId, content string, deadline int64) {
	var ret RemoveClosedNotification
	ret.EventType = types.EventRoomClosed
	ret.RoomId = roomId
	if deadline != 0 {
		//封群
		ret.Content = content
	}
	ret.DisableDeadline = deadline
	ret.Datetime = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	roomChannelId := types.GetRoomRouteById(roomId)
	//发送给所有人
	if cl, _ := router.GetChannel(roomChannelId); cl != nil {
		cl.Broadcast(data)
	}
}

// 发送封号通知
func SendCloseUserAccountNotification(deadline int64, userId, content string) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventUserClosed
	ret["disableDeadline"] = deadline
	if deadline != 0 {
		ret["content"] = content
	}
	ret["datetime"] = utility.NowMillionSecond()
	data, _ := json.Marshal(ret)

	if user, _ := router.GetUser(userId); user != nil {
		user.SendToAllClients(data)
	}
}

// 发送更新用户公钥通知 S->C
func SendUpdatePubKeyNotification(userId, publicKey, privateKey string, users []string) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventBroadcastPubKey
	ret["userId"] = userId
	ret["publicKey"] = publicKey
	ret["privateKey"] = privateKey

	data, _ := json.Marshal(ret)

	go func() {
		for _, u := range users {
			if user, _ := router.GetUser(u); user != nil {
				user.SendToAllClients(data)
			}
		}
	}()
}

// 发送添加好友通知 S->C
func SendAddFriendNotification(userId, logId, reason, source string, source2 interface{}, status int, senderInfo, receiveInfo map[string]interface{}) {
	var ret = make(map[string]interface{})
	ret["eventType"] = types.EventAddFriend
	ret["id"] = logId
	ret["senderInfo"] = senderInfo
	ret["receiveInfo"] = receiveInfo
	ret["type"] = types.IsFriend
	ret["applyReason"] = reason
	ret["source"] = source
	ret["source2"] = source2
	ret["status"] = status
	ret["datetime"] = utility.NowMillionSecond()

	data, _ := json.Marshal(ret)

	if user, _ := router.GetUser(userId); user != nil {
		user.SendToAllClients(data)
	}
}

func ComposeSystemMsg(msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["content"] = msg
	return content
}

func ComposeAlert(alertType int, operator, name, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["name"] = name
	content["content"] = msg
	return content
}

func ComposeRPReceiveAlert(alertType int, operator, name, owner, ownerName, packetId, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["name"] = name
	content["owner"] = owner
	content["ownerName"] = ownerName
	content["content"] = msg
	content["packetId"] = packetId
	return content
}

func ComposeDelFriendAlert(alertType int, operator, targetId, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["target"] = targetId
	content["content"] = msg
	return content
}

func ComposeAddFriendAlert(alertType int, operator, name, targetId, targetName, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["name"] = name
	content["target"] = targetId
	content["targetName"] = targetName
	content["content"] = msg
	return content
}

func ComposeInviteRoomAlert(alertType int, operator, name string, users, names []string, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["name"] = name
	content["users"] = users
	content["names"] = names
	content["content"] = msg
	return content
}

func ComposeCreateRoomMemberAlert(alertType int, operator, name, targetId, targetName, roomId, roomName, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["name"] = name
	content["target"] = targetId
	content["targetName"] = targetName
	content["roomName"] = roomName
	content["roomId"] = roomId
	content["content"] = msg
	return content
}

func ComposeSetRoomNameAlert(alertType int, operator, name, roomName, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["name"] = name
	content["roomName"] = roomName
	content["content"] = msg
	return content
}

func ComposeSetMemberLevelAlert(alertType int, operator, targetId, targetName, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["target"] = targetId
	content["targetName"] = targetName
	content["content"] = msg
	return content
}

func ComposeKickOutAlert(alertType int, operator, targetId, targetName, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = alertType
	content["operator"] = operator
	content["target"] = targetId
	content["targetName"] = targetName
	content["content"] = msg
	return content
}

//撤回消息通知
func ComposeRevokeAlert(logId, operator, name, msg string, level int) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertRevokeMsg
	content["logId"] = logId
	content["operator"] = operator
	content["name"] = name
	content["level"] = level
	content["content"] = msg
	return content
}

//焚毁消息通知
func ComposeHadBurntAlert(channelType int, logId string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertHadBurntMsg
	content["channelType"] = channelType
	content["logId"] = logId
	return content
}

//群禁言消息通知
func ComposeRoomMutedAlert(operator string, mutedType, level, opt int, names []string, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertRoomMuted
	content["mutedType"] = mutedType
	content["operator"] = operator
	content["level"] = level
	content["opt"] = opt
	content["names"] = names
	content["content"] = msg
	return content
}

//组合成转发消息
func ComposeForwardMsg(channelType int, fromId, fromName, forwardUserName string, data []*types.ChatLog) map[string]interface{} {
	var content = make(map[string]interface{})
	content["forwardType"] = types.MergeForward
	content["channelType"] = channelType
	content["fromId"] = fromId
	content["fromName"] = fromName
	content["forwardUserName"] = forwardUserName
	content["data"] = data
	return content
}

func ComposeSingleForwardMsg(channelType int, fromId, fromName string, data *types.ChatLog) map[string]interface{} {
	content := data.Msg.(map[string]interface{})
	content["forwardType"] = types.SingleForward
	content["channelType"] = channelType
	content["fromId"] = fromId
	content["fromName"] = fromName
	return content
}

func ComposePrintScreen(operator string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertPrintScreen
	content["operator"] = operator
	return content
}

//收款成功消息通知
func ComposePaymentAlert(operator, logId, recordId string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertPayment
	content["operator"] = operator
	content["logId"] = logId
	content["recordId"] = recordId
	return content
}

//更新群秘钥消息
func ComposeUpdateSKeyAlert(roomId, fromKey, kid, key string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertUpdateSKey
	content["roomId"] = roomId
	content["fromKey"] = fromKey
	content["key"] = key
	content["kid"] = kid
	return content
}

func ComposeRoomInviteCard(roomId, markId, roomName, inviter, avatar, identificationInfo string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["roomId"] = roomId
	content["markId"] = markId
	content["roomName"] = roomName
	content["inviterId"] = inviter
	content["avatar"] = avatar
	content["identificationInfo"] = identificationInfo
	return content
}

//拒绝加入群聊
func ComposeRoomInviteReject(operator, targetId, targetName, msg string) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertRoomInviteReject
	content["operator"] = operator
	content["target"] = targetId
	content["targetName"] = targetName
	content["content"] = msg
	return content
}

// 赞赏通知
func ComposePraiseAlert(logId, operator, action string, like, reward int) map[string]interface{} {
	var content = make(map[string]interface{})
	content["type"] = types.AlertPraise
	content["operator"] = operator
	content["logId"] = logId
	content["action"] = action
	content["like"] = like
	content["reward"] = reward
	return content
}
