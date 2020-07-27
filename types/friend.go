package types

//好友表里的字段
const (
	//是否删除 1 未删除 2 已删除
	FriendIsNotDelete = 1
	FriendIsDelete    = 2
)

//是否删除 1 未删除 2 已删除
const IsNotDelete = 1
const IsDelete = 2

//加好友设置
// 1需要验证 2 不需要验证
const NeedConfirm = 1
const NotNeedConfirm = 2

//1 需要回答问题 2 不需要回答
const NeedAnswer = 1
const NotNeedAnswer = 2

const AnswerFalse = 1 //答案错误
const SendSuccess = 2 //请求发送成功
const AddSuccess = 3  //添加成功

const FriendMsgIsNotDelete = 2
const FriendMsgIsDelete = 1

const IsNotBlocked = 0
const IsBlocked = 1

type ExtRemark struct {
	Telephones []struct {
		Remark *string `json:"remark"`
		Phone  string  `json:"phone"`
	} `json:"telephones"`
	Description *string  `json:"description"`
	Pictures    []string `json:"pictures"`
	Encrypt     string   `json:"encrypt"`
}

type FriendInfoApi struct {
	Uid                string `json:"uid"`
	Id                 string `json:"id"`
	Name               string `json:"name"`
	Avatar             string `json:"avatar"`
	Remark             string `json:"remark"`
	PublicKey          string `json:"publicKey"`
	NoDisturbing       int    `json:"noDisturbing"`
	CommonlyUsed       int    `json:"commonlyUsed"`
	OnTop              int    `json:"onTop"`
	IsDelete           int    `json:"isDelete"`
	AddTime            int64  `json:"addTime"`
	DepositAddress     string `json:"depositAddress"`
	Identification     int    `json:"identification"`
	IdentificationInfo string `json:"identificationInfo"`
	IsBlocked          int    `json:"isBlocked"`
}
