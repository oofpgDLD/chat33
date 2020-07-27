package types

const (
	DbMsgLength = 7168 //

	AppIdZhaobi      = "1001" //找币
	AppIdPWalelt     = "1002" //敢么
	AppIdYuanScraper = "1003" //
	AppIdFunnyChat   = "1004" //趣信
	AppIdTwitter     = "1005" //推特链
	AppIdTSC         = "1006" //TS钱包

	DeviceWeb     = "Web"
	DeviceAndroid = "Android"
	DeviceIOS     = "iOS"
	DevicePC      = "PC"

	//用户登录身份
	LevelVisitor = 0
	LevelMember  = 1
	LevelCs      = 2
	LevelAdmin   = 3

	//登录时的验证方式
	LoginTypeSms      = 1
	LoginTypePwd      = 2
	LoginTypeEmail    = 3
	LoginTypeEmailPwd = 4

	// 获取聊天历史记录条数上限
	ChatHistoryLimit = 1000

	//分页搜索
	SearchAll = -1
)

const (
	IsRoom   = 1
	IsFriend = 2
)

// message about
const (
	// channelType
	ToGroup = 1
	ToRoom  = 2
	ToUser  = 3

	IsSnap    = 1 //阅后即焚
	IsNotSnap = 2

	//message state
	HadRead  = 1 //已读
	NotRead  = 2 //未读
	HadBurnt = 3 //已经被焚毁
)

const (
	//event type
	EventCommonMsg        = 0 //普通消息
	EventOtherDeviceLogin = 9
	EventUserClosed       = 10
	EventRoomClosed       = 11

	EventJoinRoom         = 20
	EventLogOutRoom       = 21
	EventRemoveRoom       = 22
	EventJoinApply        = 23
	EventRoomOnlineNumber = 24
	EventRoomMuted        = 25
	//更新加密群会话秘钥
	EventUpdateSKey = 26

	//同步群会话秘钥
	EventStartGetAllSKey   = 27
	EventGetAllSKey        = 28
	EventGetAllSKeySuccess = 29

	EventAddFriend = 31
	//C->S 用户更新公钥
	EventUpdatePublicKey = 33
	//S->C 广播用户更新的公钥
	EventBroadcastPubKey = 34

	EventBatchPush       = 40
	EventBatchPushUnread = 41
	EventSyncMsg         = 42
	EventSyncMsgRlt      = 43
	//批量推送
	EventBatchCustomize = 44

	//开始消息确认
	EventStartAck   = 45
	EventBatchAck   = 46
	EventAckSuccess = 47

	// msgType
	System     = 0
	Text       = 1
	Audio      = 2
	Photo      = 3
	RedPack    = 4
	Video      = 5
	Alert      = 6
	Forward    = 8
	File       = 9
	Transfer   = 10 //转账
	Receipt    = 11 //收款
	RoomInvite = 12 //入群邀请

	//alert type
	//ws消息通知类型
	AlertRevokeMsg         = 1
	AlertCreateRoom        = 2
	AlertLoginOutRoom      = 3
	AlertKickOutRoom       = 4
	AlertJoinInRoom        = 5
	AlertRemoveRoom        = 6
	AlertAddFriendInRoom   = 7
	AlertDeleteFriend      = 8
	AlertSetAsMaster       = 9
	AlertSetAsManager      = 10
	AlertRenameRoom        = 11
	AlertReceiveRedpackage = 12
	AlertAddFriend         = 13
	AlertRoomMuted         = 14
	AlertHadBurntMsg       = 15
	AlertPrintScreen       = 16
	AlertInviteJoinRoom    = 17
	//收款成功
	AlertPayment = 18
	//更新群秘钥
	AlertUpdateSKey = 19
	//拒绝加入群聊
	AlertRoomInviteReject = 20
	//赞赏通知
	AlertPraise = 22
)

//Personal set
const (
	//是否消息免打扰 1免打扰 2关闭
	NoDisturbingOn  = 1
	NoDisturbingOff = 2

	//是否是常用类型 1 普通 2 常用
	UncommonUse      = 1
	CommonUse        = 2
	CommonUseFindAll = 3

	//是否置顶 1置顶 2不置顶
	OnTop    = 1
	NotOnTop = 2

	//是否可查看入群之前历史记录
	CanVisitLog    = 1
	CanNotVisitLog = 2

	IsEncrypt    = 1
	IsNotEncrypt = 2
)

//forward
const (
	SingleForward = 1
	MergeForward  = 2
)

// Apply
const (
	//http request 参数
	AcceptRequest = 1
	RejectRequest = 2

	//apply ack
	AwaitState  = 1
	RejectState = 2
	AcceptState = 3
)

//red packet
const (
	//红包toid   1群 0好友
	RPToUser = 0
	RPToRoom = 1

	PacketTypeLucky = 1 //拼手气红包
	PacketTypeAdv   = 2 //推广红包

	RecvForOldCustomer = 1 //老用户领取
	RecvForNewCustomer = 2 //未注册用户领取

	PageNo    = 1
	PageLimit = 15
)

const (
	//消息是否同步推送
	RecentUnpush = 1
	RecentPushed = 2
)

//个推安卓ios
const (
	DevAndroid = 1
	DevIOS     = 2
)

const (
	//加群/好友来源
	Search = 1
	Scan   = 2
	Group  = 4 //通过群
	//
	Share  = 3
	Invite = 4
	Unknow = 5
)

const (
	//认证对象
	VerifyForUser = 1
	VerifyForRoom = 2

	//认证状态
	//1：待审核；2已认证；3未通过
	VerifyStateWait   = 1
	VerifyStateAccept = 2
	VerifyStateReject = 3

	//认证处理请求参数
	VerifyApproveReject  = 0
	VerifyApproveAccept  = 1
	VerifyApproveFeeBack = 2

	//认证费用去向
	VerifyRecordCost   = 1
	VerifyRecordRefund = 2

	//认证费用到账状态
	VerifyFeeStateCosting     = 0
	VerifyFeeStateSuccess     = 1
	VerifyFeeStateFailed      = 2
	VerifyFeeStateBacking     = 3
	VerifyFeeStateBackSuccess = 4
	VerifyFeeStateBackFailed  = 5

	//群和用户状态 是否认证
	Unverified = 0
	Verified   = 1

	//群 认证人数限制
	VerifiedLimitMembers = 5000

	//加v申请时间间隔 单位分钟
	//7天
	VerifyApplyInterval = 10080
)

var VerifyFeeRecord = map[int]int{
	VerifyFeeStateCosting: VerifyRecordCost,
	VerifyFeeStateBacking: VerifyRecordRefund,
}

var VerifyApproveToState = map[int]int{
	VerifyApproveReject: VerifyStateReject,
	VerifyApproveAccept: VerifyStateAccept,
}
