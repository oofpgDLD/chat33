package types

var ExcelAddr = ""

// chat appId  --> red packet Platform
var CustomServiceInfo = map[string]string{
	AppIdZhaobi:    "微信号zhaobikefu", //币钱包
	AppIdFunnyChat: "QQ号762389302",  //趣聊 --> 币钱包
	AppIdTSC:       "微信号17621543259",
}

// request param  --> db tag
var QueryAdminOptParamToDb = map[int]int{
	QueryAdminOptLogBanUser:       BanUser,       //封号
	QueryAdminOptLogBanRoom:       BanRoom,       //封群
	QueryAdminOptLogCancleBanUser: BanUserCancel, //解封号
	QueryAdminOptLogCancleBanRoom: BanRoomCancel, //解封群
}

// request param  --> db tag
var QueryAdminOptDbToParam = map[int]int{
	BanUser:       QueryAdminOptLogBanUser,       //封号
	BanRoom:       QueryAdminOptLogBanRoom,       //封群
	BanUserCancel: QueryAdminOptLogCancleBanUser, //解封号
	BanRoomCancel: QueryAdminOptLogCancleBanRoom, //解封群
}

const (
	CountNew    = 0
	CountActive = 1
	CountOpen   = 2
	CountAll    = 3

	BanUser       = 1
	BanRoom       = 2
	BanUserCancel = 3
	BanRoomCancel = 4

	QueryUsersAll     = 99
	QueryUsersUnClose = 0
	QueryUsersClosed  = 1

	QueryRoomsAll       = 99
	QueryRoomsUnClose   = 0
	QueryRoomsClosed    = 1
	QueryRoomsRecommend = 2

	QueryAdminOptLogAll           = 99
	QueryAdminOptLogBanUser       = 0
	QueryAdminOptLogBanRoom       = 1
	QueryAdminOptLogCancleBanUser = 2
	QueryAdminOptLogCancleBanRoom = 3

	AdIsNotDelete = 0
	AdIsDelete    = 1
	MaxAdNumbers  = 5
	AdIsActive    = 1
	AdIsNotActive = 0
)
