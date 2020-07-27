package proto

import (
	"bytes"
	"encoding/json"

	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
)

//开始同步消息
type EventSyncMsg struct {
	time int64
}

func (p *EventSyncMsg) Receive(conn router.Connection, msg []byte) {
	receive := proto_receiver[types.EventSyncMsg]
	if receive == nil {
		protoLog.Warn("can not find receive handler")
	}
	receive(p, conn, msg)
}

func (p *EventSyncMsg) FromBytes(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return result.NewError(result.MsgFormatError)
	}

	if !p.parse(m) {
		return result.NewError(result.MsgFormatError)
	}
	return nil
}

func (p *EventSyncMsg) GetTime() int64 {
	return p.time
}

type secretInfo struct {
	userId string
	key    string
}

func (t *secretInfo) GetUserId() string {
	return t.userId
}

func (t *secretInfo) GetKey() string {
	return t.key
}

//更新会话秘钥
type EventUpdateSKey struct {
	roomId  string
	fromKey string
	secret  []secretInfo
}

func (p *EventUpdateSKey) GetRoomId() string {
	return p.roomId
}

func (p *EventUpdateSKey) GetSecret() []secretInfo {
	return p.secret
}

func (p *EventUpdateSKey) GetFromKey() string {
	return p.fromKey
}

func (p *EventUpdateSKey) Receive(conn router.Connection, msg []byte) {
	receive := proto_receiver[types.EventUpdateSKey]
	if receive == nil {
		protoLog.Warn("can not find receive handler")
	}
	receive(p, conn, msg)
}

func (p *EventUpdateSKey) FromBytes(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return result.NewError(result.MsgFormatError)
	}

	if !p.parse(m) {
		return result.NewError(result.MsgFormatError)
	}
	return nil
}

//同步群会话秘钥
type EventGetSKey struct {
	datetime int64
}

func (p *EventGetSKey) Receive(conn router.Connection, msg []byte) {
	receive := proto_receiver[types.EventStartGetAllSKey]
	if receive == nil {
		protoLog.Warn("can not find receive handler")
	}
	receive(p, conn, msg)
}

func (p *EventGetSKey) FromBytes(b []byte) error {
	var m map[string]interface{}

	//使用UseNumber防止反序列化时进度丢失
	d := json.NewDecoder(bytes.NewBuffer(b))
	d.UseNumber()

	err := d.Decode(&m)
	if err != nil {
		return result.NewError(result.MsgFormatError)
	}

	if !p.parse(m) {
		return result.NewError(result.MsgFormatError)
	}
	return nil
}

func (p *EventGetSKey) GetDatetime() int64 {
	return p.datetime
}

//用户更新公钥
type EventUpdatePubKey struct {
	PublicKey  string
	PrivateKey string
}

func (p *EventUpdatePubKey) Receive(conn router.Connection, msg []byte) {
	receive := proto_receiver[types.EventUpdatePublicKey]
	if receive == nil {
		protoLog.Warn("can not find receive handler")
	}
	receive(p, conn, msg)
}

func (p *EventUpdatePubKey) FromBytes(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return result.NewError(result.MsgFormatError)
	}

	if !p.parse(m) {
		return result.NewError(result.MsgFormatError)
	}
	return nil
}

//消息确认机制
type EventAck struct {
	begin int64
	end   int64
	total int64
}

func (p *EventAck) Receive(conn router.Connection, msg []byte) {
	receive := proto_receiver[types.EventStartAck]
	if receive == nil {
		protoLog.Warn("can not find receive handler")
	}
	receive(p, conn, msg)
}

func (p *EventAck) FromBytes(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return result.NewError(result.MsgFormatError)
	}

	if !p.parse(m) {
		return result.NewError(result.MsgFormatError)
	}
	return nil
}

func (p *EventAck) GetBegin() int64 {
	return p.begin
}

func (p *EventAck) GetEnd() int64 {
	return p.end
}

func (p *EventAck) GetTotal() int64 {
	return p.total
}
