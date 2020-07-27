package types

//群 个推
type RoomCid struct {
	UserId   string
	GetuiCid string
}

//群成员信息
type RoomMemberInfo RoomMember

//群禁言
type RoomUserMued struct {
	UserId   string
	ListType int
	Deadline int64
}
