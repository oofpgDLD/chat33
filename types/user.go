package types

const (
	ValidateCodeTypeSMS = 0
	ValidateCodeTypeEmail = 1
)

//邀请入群
// 1需要验证 0 不需要验证
const RoomInviteNeedConfirm = 1
const RoomInviteNotNeedConfirm = 2

var DefaultInviteConfig = map[string]int{
	"1001": RoomInviteNotNeedConfirm,
	"1006": RoomInviteNeedConfirm,
}

type Sender struct {
	UserId   string
	NickName string
	Avatar   string
}
