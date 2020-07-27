package types

const (
	//PS praise state
	PS_Like = 1 << iota
	PS_Reward
)

const (
	Like   = 1
	Reward = 2
)

const (
	PraiseIsDelete    = 1
	PraiseIsNotDelete = 0
)

//praise record api
type PraiseRecord struct {
	RecordId    string `json:"recordId"`
	LogId       string `json:"logId"`
	ChannelType int    `json:"channelType"`
	CreateTime  int64  `json:"createTime"`
	User        struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
	Type     int     `json:"type"`
	CoinName string  `json:"coinName"`
	Amount   float64 `json:"amount"`
}

//赞赏榜单
type Record struct {
	Ranking int `json:"ranking"`
	User    struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
	Number int     `json:"number"`
	Price  float64 `json:"price"`
}

type PraiseBoardHistory struct {
	Like struct {
		Ranking int `json:"ranking"`
		Number  int `json:"number"`
	} `json:"like"`
	Reward struct {
		Ranking int     `json:"ranking"`
		Price   float64 `json:"price"`
	} `json:"reward"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
}

//公司-员工列表
type Clist struct {
	Enterprise struct {
		Name string `json:"name"`
	} `json:"enterprise"`
	Members []*Member
}

type Member struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
