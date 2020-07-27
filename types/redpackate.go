package types

const RemarkLengthLimit = 20

const (
	AuthType = "Bearer"
)

type ReqBase struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type SendParams struct {
	Amount float64     `json:"amount"`
	CoinId int         `json:"coin_id"`
	Size   int         `json:"size"`
	Type   int         `json:"type"`
	Remark string      `json:"remark"`
	To     string      `json:"to"`
	Ext    interface{} `json:"ext"`
}

type ReceiveParams struct {
	PacketId string `json:"packet_id"`
}

type InfoParams struct {
	PacketId string `json:"packet_id"`
}

type ReceiveDetailParams struct {
	PacketId string `json:"packet_id"`
}

type StatisticParams struct {
	CoinId    int   `json:"coin_id"`
	Operation int   `json:"operation"`
	Type      int   `json:"type"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	PageNumer int   `json:"page_num"`
	PageSize  int   `json:"page_size"`
}

type CoinInfoParams struct {
	PlatformId string `json:"platform_id"`
}

//-----api result---------//
type RPReceiveInfo struct {
	UserId      string  `json:"userId"`
	UserName    string  `json:"userName"`
	UserAvatar  string  `json:"userAvatar"`
	CoinId      int     `json:"coinId"`
	CoinName    string  `json:"coinName"`
	Amount      float64 `json:"amount"`
	CreatedAt   int64   `json:"createdAt"`
	Status      int     `json:"status"`
	FailMessage string  `json:"failMessage"`
}

type RedPacketInfo struct {
	SenderId           string         `json:"senderId"`
	SenderUid          string         `json:"senderUid"`
	SenderAvatar       string         `json:"senderAvatar"`
	SenderName         string         `json:"senderName"`
	Identification     int            `json:"identification"`
	IdentificationInfo string         `json:"identificationInfo"`
	Type               int            `json:"type"`
	CoinId             int            `json:"coinId"`
	CoinName           string         `json:"coinName"`
	Amount             float64        `json:"amount"`
	Size               int            `json:"size"`
	ToUsers            string         `json:"toUsers"`
	Remain             int            `json:"remain"`
	Status             int            `json:"status"`
	Remark             string         `json:"remark"`
	CreatedAt          int64          `json:"createdAt"`
	PacketId           string         `json:"packetId"`
	PacketUrl          string         `json:"packetUrl"`
	RevInfo            *RPReceiveInfo `json:"revInfo"`
}

type RPStatisticInfo struct {
	Count      int              `json:"count"`
	Sum        float64          `json:"sum"`
	CoinId     int              `json:"coinId"`
	CoinName   string           `json:"coinName"`
	RedPackets []*RedPacketInfo `json:"redPackets"`
}

type RedPacket struct {
	PacketId   string `json:"packetId"`
	PacketUrl  string `json:"packetUrl"`
	PacketType int    `json:"packetType"`
	Coin       int    `json:"coin"`
	Remark     string `json:"remark"`
}
