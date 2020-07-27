package proto

import (
	"github.com/33cn/chat33/types"
)

func isCmnEvTypeOk(t int) bool {
	return t == types.EventCommonMsg
}

// TODO
func isAlertEvTypeOk(t int) bool {

	switch t {
	case types.AlertRevokeMsg:
		fallthrough
	case types.AlertCreateRoom:
		fallthrough
	case types.AlertLoginOutRoom:
		fallthrough
	case types.AlertKickOutRoom:
		fallthrough
	case types.AlertJoinInRoom:
		fallthrough
	case types.AlertRemoveRoom:
		fallthrough
	case types.AlertAddFriendInRoom:
		fallthrough
	case types.AlertDeleteFriend:
		fallthrough
	case types.AlertSetAsMaster:
		fallthrough
	case types.AlertSetAsManager:
		fallthrough
	case types.AlertRenameRoom:
		fallthrough
	case types.AlertReceiveRedpackage:
		fallthrough
	case types.AlertAddFriend:
		fallthrough
	case types.AlertRoomMuted:
		fallthrough
	case types.AlertHadBurntMsg:
		fallthrough
	case types.AlertRoomInviteReject:
		return true
	}
	return false
}

func isChTypeOk(t int) bool {
	switch t {
	case types.ToGroup:
	case types.ToRoom:
		fallthrough
	case types.ToUser:
		return true
	}
	return false
}

func isMsgTypeOk(t int) bool {
	switch t {
	case types.System:
		fallthrough
	case types.Text:
		fallthrough
	case types.Audio:
		fallthrough
	case types.Photo:
		fallthrough
	case types.RedPack:
		fallthrough
	case types.Video:
		fallthrough
	case types.Alert:
		fallthrough
	case types.Forward:
		fallthrough
	case types.File:
		fallthrough
	case types.Transfer:
		fallthrough
	case types.Receipt:
		fallthrough
	case types.RoomInvite:
		return true
	}
	return false
}

func isSnapTypeOk(t int) bool {
	switch t {
	case types.IsSnap:
		fallthrough
	case types.IsNotSnap:
		return true
	}
	return false
}

func isTypedMsgOk(msgType int, m msg) bool {
	switch msgType {
	case types.System:
		return m.checkSysMsg()
	case types.Text:
		return m.checkTxtMsg()
	case types.Audio:
		return m.checkAudioMsg()
	case types.Photo:
		return m.checkImgMsg()
	case types.RedPack:
		return m.checkRedPacketMsg()
	case types.Video:
		return m.checkVideoMsg()
	case types.File:
		return m.checkFileMsg()
	case types.Alert:
		// 客户端不会主动发通知类消息
		return false
	case types.Transfer:
		return true
	case types.Receipt:
		return true
	}
	return false
}

func isSpecialType(t int) bool {
	switch t {
	case types.AlertRoomInviteReject:
		fallthrough
	case types.AlertLoginOutRoom:
		fallthrough
	case types.AlertReceiveRedpackage:
		fallthrough
	case types.AlertHadBurntMsg:
		fallthrough
	case types.AlertPrintScreen:
		fallthrough
	case types.AlertAddFriendInRoom:
		return true
	}
	return false
}

func isMsgLengthOk(m msg) bool {
	return m.checkMsgLength()
}
