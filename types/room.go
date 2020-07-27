package types

const (
	RoomLevelNotExist = 0
	RoomLevelNomal    = 1
	RoomLevelManager  = 2
	RoomLevelMaster   = 3

	CanAddFriend      = 1
	CanNotAddFriend   = 2
	ShouldApproval    = 1
	ShouldNotApproval = 2
	CanNotJoinRoom    = 3
	CanReadAllLog     = 1
	CanNotReadAllLog  = 2

	RoomUserNotDeleted   = 1
	RoomUserDeleted      = 2
	RoomUserDeletedOrNot = 99

	RoomMsgNotDelete = 1
	RoomMsgDeleted   = 2

	RoomNotDeleted   = 1
	RoomDeleted      = 2
	RoomDeletedOrNot = 99

	RoomRecommend    = 1
	RoomNotRecommend = 0
)

// Room Muted
const (
	AdminMuted    = 2
	AdminNotMuted = 1

	MasterMuted    = 2
	MasterNotMuted = 1

	AllSpeak  = 1
	Blacklist = 2
	Whitelist = 3
	AllMuted  = 4

	MutedEnable  = 1
	MutedDisable = 2
	MutedForvevr = 7258089600000
)

type ApplyList struct {
	SenderInfo  *ApplyInfo  `json:"senderInfo"`
	ReceiveInfo *ApplyInfo  `json:"receiveInfo"`
	Id          string      `json:"id"`
	Type        int         `json:"type"`
	ApplyReason string      `json:"applyReason"`
	Status      int         `json:"status"`
	Datetime    int64       `json:"datetime"`
	Source      string      `json:"source"`
	Source2     interface{} `json:"source2"`
}

type ApplyInfo struct {
	Id                 string `json:"id"`
	MarkId             string `json:"markId"`
	Name               string `json:"name"`
	Avatar             string `json:"avatar"`
	Position           string `json:"position"`
	Identification     int    `json:"identification"`
	IdentificationInfo string `json:"identificationInfo"`
}

//聊天记录
type ChatLog struct {
	LogId       string      `json:"logId"`
	MsgId       string      `json:"msgId"`
	ChannelType int         `json:"channelType"`
	IsSnap      int         `json:"isSnap"`
	FromId      string      `json:"fromId"`
	TargetId    string      `json:"targetId"`
	MsgType     int         `json:"msgType"`
	Msg         interface{} `json:"msg"`
	Datetime    int64       `json:"datetime"`
	SenderInfo  interface{} `json:"senderInfo"`
	Ext         interface{} `json:"ext"`
	Praise      interface{} `json:"praise",omitempty`
	con         map[string]interface{}
	isDelete    int
}

func (c *ChatLog) GetCon() map[string]interface{} {
	return c.con
}

func (c *ChatLog) SetCon(con map[string]interface{}) {
	c.con = con
}

func (c *ChatLog) GetIsDel() int {
	return c.isDelete
}

func (c *ChatLog) SetIsDel(isDel int) {
	c.isDelete = isDel
}

type chatLogs []*ChatLog

func NewchatLogs() chatLogs {
	return make([]*ChatLog, 0)
}

func (cls chatLogs) Len() int {
	return len(cls)
}

func (cls chatLogs) Less(i, j int) bool {
	return cls[i].Datetime < cls[j].Datetime
}

func (cls chatLogs) Swap(i, j int) {
	cls[i], cls[j] = cls[j], cls[i]
}

type ackChatLogs []*ChatLog

func NewackChatLogs() ackChatLogs {
	return make([]*ChatLog, 0)
}

func (cls ackChatLogs) Len() int {
	return len(cls)
}

func (cls ackChatLogs) Less(i, j int) bool {
	return cls[i].Datetime > cls[j].Datetime
}

func (cls ackChatLogs) Swap(i, j int) {
	cls[i], cls[j] = cls[j], cls[i]
}
