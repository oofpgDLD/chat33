package types

type Tx interface {
	RollBack()
	Commit() error
}

//db room table
type Room struct {
	Id                 string `json:"id"`
	MarkId             string `json:"markId"`
	Name               string `json:"name"`
	Avatar             string `json:"avatar"`
	MasterId           string `json:"masterId"`
	CreateTime         int64  `json:"-"`
	CanAddFriend       int    `json:"canAddFriend"`
	JoinPermission     int    `json:"joinPermission"`
	RecordPermision    int    `json:"recordPermission"`
	AdminMuted         int    `json:"-"`
	MasterMuted        int    `json:"-"`
	Encrypt            int    `json:"encrypt"`
	IsDelete           int    `json:"-"`
	RoomLevel          int    `json:"-"`
	CloseUntil         int64  `json:"-"`
	Recommend          int    `json:"recommend"`
	Identification     int    `json:"identification"`
	IdentificationInfo string `json:"identificationInfo"`
}

type RoomMember struct {
	Id           string
	RoomId       string
	UserId       string
	UserNickname string
	Level        int
	NoDisturbing int
	CommonUse    int
	RoomTop      int
	CreateTime   int64
	IsDelete     int
	Source       string
}

type MemberJoinUser struct {
	*RoomMember
	*User
}

type MemberJoinRoom struct {
	*Room
	*RoomMember
}

type RoomJoinUser struct {
	*Room
	*User
}

type RoomLog struct {
	Id       string
	MsgId    string
	RoomId   string
	SenderId string
	IsSnap   int
	MsgType  int
	Content  string
	Datetime int64
	Ext      string
	IsDelete int
}

type RoomLogJoinUser struct {
	*RoomLog
	*User
}

type RoomMsgReceive struct {
	Id        string
	RoomMsgId string
	ReceiveId string
	State     int
}

type RoomUserMuted struct {
	Id       string
	RoomId   string
	UserId   string
	ListType int
	Deadline int64
}

type RoomConfig struct {
	AppId string
	Level int
	Limit int
}

//---------user----------//
type User struct {
	UserId             string
	MarkId             string
	Uid                string
	AppId              string
	Username           string
	Account            string
	UserLevel          int
	Verified           int //是否实名制
	Avatar             string
	CompanyId          string
	Position           string
	Area               string
	Phone              string
	Sex                int
	Email              string
	InviteCode         string
	DeviceToken        string
	DepositAddress     string
	PublicKey          string
	PrivateKey         string
	Device             string
	CreateTime         int64
	RegVersion         string
	NowVersion         string
	CloseUntil         int64
	SuperUserLevel     int
	Identification     int
	IdentificationInfo string
	//
	IsSetPayPwd int
	//是否上链
	IsChain int
}

type Apply struct {
	Id          string
	Type        int
	ApplyUser   string
	Target      string
	ApplyReason string
	State       int
	Remark      string
	Datetime    int64
	Source      string
}

type LoginLog struct {
	Id         string
	UserId     string
	LoginTime  int64
	Device     string
	DeviceName string
	LoginType  int
	Uuid       string
	Version    string
}

type UserConfig struct {
	AppId string
	Level int
	Limit int
}

//--------------friend-------------//
type Friend struct {
	UserId    string
	FriendId  string
	Remark    string
	AddTime   int64
	DND       int
	Top       int
	Type      int
	IsDelete  int
	Source    string
	ExtRemark string
	IsBlocked int
}

type RedPacketLog struct {
	PacketId  string
	CType     int
	UserId    string
	ToId      string
	Coin      int
	Size      int
	Amount    float64
	Remark    string
	Type      int
	CreatedAt int64
}

type PrivateLog struct {
	Id        string
	MsgId     string
	SenderId  string
	ReceiveId string
	IsSnap    int
	MsgType   int
	Content   string
	Status    int
	SendTime  int64
	Ext       string
	IsDelete  int
}

type AddFriendConf struct {
	Id          string
	UserId      string
	NeedConfirm string
	NeedAnswer  string
	Question    string
	Answer      string
}

type FriendJoinUser struct {
	*Friend
	*User
}

type PrivateLogJoinUser struct {
	*PrivateLog
	*User
	Remark string //用于解决，私聊历史记录 返回的好友信息 remark字段
}

//--------------app----------------//
type Module struct {
	Type   int    `json:"type"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

type UPushType struct {
	IOS     string `json:"iOS"`
	Android string `json:"Android"`
}

type App struct {
	AppId         string `json:"app_id"`
	AccountServer string `json:"user_info_url"`
	RPPid         string
	RPServer      string
	RPUrl         string
	IsInner       int
	MainCoin      string

	AppKey    string //`json:"-"`
	AppSecret string //`json:"-"`

	//推送相关
	PushAppKey          UPushType
	PushAppMasterSecret UPushType
	PushMiActive        UPushType //走系统推送时打开的active
	//模块启用状态
	Modules []*Module

	//otc
	IsOtc     int
	OtcServer string

	//oss
	OssConfig *OssAddress
}

type OssAddress struct {
	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
}

//--------------version----------------//
type Version struct {
	AppId          string
	Name           string
	Compatible     bool   `json:"compatible"`
	MinVersionCode int    `json:"minVersionCode"`
	VersionCode    int    `json:"versionCode"`
	VersionName    string `json:"versionName"`
	Description    string `json:"description"`
	Url            string `json:"url"`
	Size           int64  `json:"size"`
	Md5            string `json:"md5"`
	ForceList      string `json:"forceList"`
}

//---------------admin--------//
type Admin struct {
	Id       string
	AppId    string
	Account  string
	Password string
	Salt     string
}

type AdminOptLog struct {
	Id            string
	Operator      string
	Type          int
	Target        string
	OperateType   int
	Reason        string
	CreateTime    int64
	EffectiveTime int64
}

type AdminOptLogView struct {
	Id            string
	AppId         string
	Operator      string
	Type          int
	Target        string
	MarkId        string
	Name          string
	OperateType   int
	Reason        string
	CreateTime    int64
	EffectiveTime int64
}

type ActiveUsersView struct {
	UserId string
	AppId  string
	Count  string
}

type Advertisement struct {
	Id       string `json:"id"`
	AppId    string `json:"-"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Duration int    `json:"duration"`
	Link     string `json:"link"`
	IsActive int    `json:"isActive"`
	IsDelete int    `json:"isDelete"`
}

type InviteRoomConf struct {
	UserId      string
	NeedConfirm int
}

//----------认证手续费表--------------//
type VerifyFee struct {
	AppId    string  `json:"-"`
	Type     int     `json:"type"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount,string"`
}

//--------认证申请-------------------//
type VerifyApply struct {
	Id          string
	AppId       string
	Type        int
	TargetId    string
	Description string
	Amount      float64
	Currency    string
	State       int
	RecordId    string
	UpdateTime  int64
	FeeState    int
}

type VerifyApplyJoinUser struct {
	*VerifyApply
	*User
}

type VerifyApplyJoinRoomAndUser struct {
	*VerifyApply
	*Room
	*User
}

//----------------otc 交易关系---------------------//
type Order struct {
	OrderId  string
	Uid      string //用户uid
	Opposite string //对方uid
}

//--------------Parise------------------------//
type Praise struct {
	Id          string
	ChannelType int
	TargetId    string
	LogId       string
	SenderId    string
	OptId       string
	Type        int
	RecordId    string
	CoinId      string
	CoinName    string
	Amount      float64
	CreateTime  int64
	IsDelete    int
}

type PraiseUser struct {
	Id         string
	TargetId   string
	OptId      string
	RecordId   string
	CoinId     string
	CoinName   string
	Amount     float64
	CreateTime int64
	IsDelete   int
}

//赞赏榜单 redis使用
type RankingItem struct {
	UserId string
	Type   int
	Number int
	Price  float64
}

//排序
type RankingItemWrapper struct {
	Items []*RankingItem
	By    func(p, q *RankingItem) bool
}

func (pw RankingItemWrapper) Len() int { // 重写 Len() 方法
	return len(pw.Items)
}
func (pw RankingItemWrapper) Swap(i, j int) { // 重写 Swap() 方法
	pw.Items[i], pw.Items[j] = pw.Items[j], pw.Items[i]
}
func (pw RankingItemWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return pw.By(pw.Items[i], pw.Items[j])
}
