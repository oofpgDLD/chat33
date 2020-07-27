package proto

import (
	"encoding/json"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	l "github.com/inconshreveable/log15"

	"github.com/33cn/chat33/result"
)

var protoLog = l.New("module", "proto/proto")

// 客户端发送的消息结构
type Proto struct {
	evType   int // 0: 普通消息, 20:入群通知 ...
	msgId    string
	chType   int // 1: 聊天室； 2：群组；3：好友
	targetId string
	msgType  int // 0：系统消息，1:文字，2:音频，3：图片，4：红包，5：视频, 6 通知消息
	isSnap   int // 是否阅后即焚 1 是 2 否
	msg      msg // 对应 msgType 的具体消息内容
	sendTime int64
	praise   interface{} `json:"praise",omitempty`
	ext      interface{} //额外信息字段

	logId string
}

func (p *Proto) Receive(conn router.Connection, msg []byte) {
	receive := proto_receiver[types.EventCommonMsg]
	if receive == nil {
		protoLog.Error("Receive", "err", "can not find receive handler", "event type", types.EventCommonMsg)
	}
	receive(p, conn, msg)
}

func (p *Proto) FromBytes(b []byte) error {
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

// 普通消息
func NewCmnProto(logId, msgId, targetId string, chType, msgType, isSnap int, msg map[string]interface{}, sendTime int64, praise interface{}) (*Proto, error) {
	if !(isChTypeOk(chType) && isMsgTypeOk(msgType) && isSnapTypeOk(isSnap)) {
		protoLog.Warn("NewCmnProto", "warn", "ParamsError", "chType", chType, "msgType", msgType, "isSnap", isSnap)
		return nil, result.NewError(result.ParamsError)
	}

	m := &Proto{
		evType:   types.EventCommonMsg,
		chType:   chType,
		targetId: targetId,
		msgType:  msgType,
		isSnap:   isSnap,
		msg:      msg,
		sendTime: sendTime,
		logId:    logId,
		msgId:    msgId,
		praise:   praise,
	}
	return m, nil
}

func (p *Proto) GetRouter() string {
	switch p.chType {
	case types.ToGroup:
		return types.GetGroupRouteById(p.targetId)
	case types.ToRoom:
		return types.GetRoomRouteById(p.targetId)
	case types.ToUser:
		return p.targetId
	}
	return ""
}

func (p *Proto) GetChannelType() int {
	return p.chType
}

func (p *Proto) GetTargetId() string {
	return p.targetId
}

func (p *Proto) SetTargetId(targetId string) {
	p.targetId = targetId
}

func (p *Proto) GetMsg() map[string]interface{} {
	return p.msg
}

func (p *Proto) GetMsgId() string {
	return p.msgId
}

func (p *Proto) GetEvType() int {
	return p.evType
}

func (p *Proto) GetMsgType() int {
	return p.msgType
}

func (p *Proto) GetLogId() string {
	return p.logId
}

func (p *Proto) SetLogId(logId string) {
	p.logId = logId
}

func (p *Proto) GetIsSnap() int {
	return p.isSnap
}

func (p *Proto) GetSendTime() int64 {
	return p.sendTime
}

func (p *Proto) GetExt() interface{} {
	return p.ext
}

func (p *Proto) GetPraise() interface{} {
	return p.praise
}

func (p *Proto) IsSpecial() bool {
	if p.isSnap == types.IsSnap {
		return true
	}
	if p.msgType == types.Alert {
		return p.msg.CheckSpecialMsg()
	}
	return false
}

func (p *Proto) IsEncryptedMsg() bool {
	return p.msg.isEncryptedMsg()
}

func (p *Proto) WrapResp(user *types.User, sendTime int64) ([]byte, error) {
	var data = make(map[string]interface{})
	data["eventType"] = p.evType
	data["msgId"] = p.msgId
	data["channelType"] = p.chType
	data["isSnap"] = p.isSnap
	data["targetId"] = p.targetId
	data["msgType"] = p.msgType
	data["msg"] = p.msg
	/*if p.msgType == types.RedPack {
		data["msg"] = redpacketLengthCheck(p.msg)
	}*/
	data["ext"] = p.GetExt()

	data["logId"] = p.logId
	data["fromId"] = user.UserId
	data["datetime"] = sendTime
	var sender = make(map[string]interface{})
	sender["nickname"] = user.Username
	sender["avatar"] = user.Avatar
	data["senderInfo"] = sender

	return json.Marshal(data)
}

// TODO
func (p *Proto) AppendChatLog(userId string, state int, time int64) error {
	var msgStr = utility.StructToString(p.msg)
	if msgStr == "" {
		protoLog.Error("AppendChatLog", "warn", "msg struck err", "user id", userId, "group id", p.targetId)
		return result.NewError(result.ParamsError)
	}
	var ext = utility.StructToString(p.ext)

	var logId int64
	var err error
	switch p.chType {
	case types.ToRoom:
		logId, err = orm.AppendRoomChatLog(userId, p.targetId, p.msgId, p.msgType, p.isSnap, msgStr, ext, time)
		if err != nil {
			protoLog.Error("AppendChatLog", "err", "room log add failed", "user id", userId, "room id", p.targetId)
			return result.NewError(result.DbConnectFail)
		}
	case types.ToUser:
		logId, err = orm.AddPrivateChatLog(userId, p.targetId, p.msgId, p.msgType, state, p.isSnap, msgStr, ext, time)
		if err != nil {
			protoLog.Error("AppendChatLog", "err", "private log add failed", "user id", userId, "group id", p.targetId)
			return result.NewError(result.DbConnectFail)
		}
	}
	p.logId = utility.ToString(logId)
	return nil
}

/*func redpacketLengthCheck(msg map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range msg {
		ret[k] = v
	}
	arry := []rune(utility.ToString(ret["remark"]))
	if len(arry) > 100 {
		ret["remark"] = string(arry[:100]) + "..."
	}
	return ret
}*/
