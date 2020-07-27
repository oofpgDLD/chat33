package cache

import (
	. "github.com/33cn/chat33/types"
	user_model "github.com/33cn/chat33/user/model"
)

//用户接口
type UserCacheI interface {
	//保存用户信息
	SaveUserInfo(u *User) error
	//获取用户信息
	GetUserInfoById(userId string) (*User, error)
	//更新用户信息
	UpdateUserInfo(userId, field, value string) error

	GetUserInfoByField(field, appId, value string) (*User, error)

	//保存用户登录信息
	SaveUserLoginInfo(*LoginLog) error
	//获取用户登录信息
	GetUserLoginInfo(userId, device string) (*LoginLog, error)

	//保存拉群 是否需验证配置
	SaveUserInviteRoomConf(*InviteRoomConf) error
	//获取拉群 是否需验证配置
	GetUserInviteRoomConf(userId string) (*InviteRoomConf, error)

	//添加推送deviceToken
	SaveDeviceToken(deviceType, deviceToken, userId string) error
	//获取deviceToken
	GetUserIdByDeviceToken(deviceToken string) (string, *string, error)
	//清除用户deviceToken
	ClearDeviceToken(userId, deviceToken string) error
}

//存储好友关系
type FriendCacheI interface {
	//存储加好友配置
	SaveAddFriendConfig(userId string, conf *AddFriendConf) error
	//获取加好友配置
	GetAddFriendConfig(userId string) (*AddFriendConf, error)

	//储存私聊记录
	SavePrivateChatLogs(log []*PrivateLog) error
	//获取私聊记录
	GetPrivateChatLog(logId string) (*PrivateLog, error)
	//删除私聊记录
	DeletePrivateChatLog(logId string) (int, error)

	//
	GetPrivateChatLogsIndexByTime(start, end *int64, startEQ, endEQ bool) ([]string, error)
	//添加索引
	//SavePrivateLogIndex(log []*PrivateLog) error
	//删除索引
	//DeletePrivateLogIndex(logId string) (int, error)

	//储存好友列表
	SaveFriends(userId string, friendIds []string) error
	//判断是否是好友
	IsFriend(userId, friendId string) (*bool, error)
	//删除好友
	DeleteFriend(userId, friendId string) error
	//添加好友 在好友请求处理那边
	AddFriend(userId, friendId string) error
	//获取好友列表
	GetFriends(userId string) ([]string, error)

	//获取已经删除的好友
	GetDelFriends(userId string) ([]string, error)
	//储存已删除好友列表
	SaveDelFriends(userId string, friendIds []string) error

	//储存好友备注(名称)等信息
	SaveFriendInfo(userId string, friendInfo *Friend) error
	//获取好友信息
	GetFriendInfo(userId, friendId string) (*Friend, error)
	//修改好友信息
	UpdateFriendInfo(userId, friendId, field, value string) error
	//删除好友信息
	DeleteFriendInfo(userId, friendId string) error
}

//群接口
type RoomCacheI interface {
	//关于群成员与cid的接口
	//储存群所有成员和cid
	SaveAllMemberCid(roomId string, rc *[]*RoomCid) error
	//判断用户是否在群
	UserIsInRoom(roomId, userId string) (bool, bool, error)
	//获取所有群成员的cid
	GetRoomUserCid(roomId string) (map[string]string, error)
	//更新/添加 cid
	UpdateRoomCid(roomId, userId, cid string) error
	//根据userId删除cid  (踢人)
	DeleteRoomCidByUserId(roomId, userId string) error
	//删除 （解散群）
	DeleteRoomCid(roomId string) error

	//关于群信息
	//保存群信息
	SaveRoomInfo(*Room) error
	//查询群详情
	FindRoomInfo(roomId string) (*Room, error)
	//更新群信息
	UpdateRoomInfo(roomId, field, value string) error
	//删除群信息
	DeleteRoomInfo(roomId string) error

	//关于群成员信息
	//保存群成员信息
	SaveRoomUser(roomId string, infos *[]*RoomMember) error
	//查询群成员信息
	FindRoomMemberInfo(roomId, userId string) (*RoomMember, error)
	//更新群成员信息
	UpdateRoomUserInfo(roomId, userId, field, value string) error
	//删除群成员信息
	DeleteRoomUserInfo(roomId, userId string) error

	//禁言信息
	//保存禁言信息
	SaveRoomUserMued(roomId string, muens *[]*RoomUserMuted) error
	//更新/添加群禁言信息
	UpdateRoomUserMuted(roomId string, rum *RoomUserMuted) error
	//查询群中禁言详情
	GetRoomUserMuted(roomId, userId string) (*RoomUserMuted, error)
	//删除群禁言信息
	DeleteRoomUserMuted(roomId string) error

	//群，群成员，群聊，群申请等信息关系
	//新增群和成员的关系
	AddRoomUser(roomId string, muens *[]*RoomMember) error
	//获取群和成员关系
	GetRoomUser(roomId string) (map[string]int64, error)
	//更新群和成员关系
	UpdateRoomUser(roomId, userId, value string) error
	//删除群和成员关系
	DelRoomUser(roomId, userId string) error
	//删除群和成员关系(群被删除的情况)
	DelRoomUserAll(roomId string) error

	//群聊天信息
	//存储群消息记录
	SaveRoomMsgContent(log *RoomLog) error
	//根据logId获取消息记录
	GetRoomMsgContent(id string) (*RoomLog, error)
	//获取所有消息记录
	GetRoomMsgContentAll() ([]*RoomLog, error)
	//根据logid删除消息记录
	DeleteRoomMsgContent(logId string) error

	//存储群recv消息记录
	SaveReceiveLog(log *RoomMsgReceive) error
	//获取所有recv消息记录 通过id
	GetReceiveLogbyId(id string) (*RoomMsgReceive, error)
	//获取所有recv消息记录
	GetReceiveLogAll() ([]*RoomMsgReceive, error)
	//更新recv消息记录
	UpdateReceiveLog(id string, state int) error
	//根据id删除recv消息记录
	DeleteReceiveLog(logId string) error

	//room_config 和 user_config
	//新增或更新room_config
	SaveRoomConfig(appId string, level, limit int) error
	//查找room_config
	FindRoomConfig(appId string) (*RoomConfig, error)
	//新增或更新user_config
	SaveUserConfig(appId string, level, limit int) error
	//获取user_config
	FindUserConfig(appId string) (*UserConfig, error)
}

type ApplyCacheI interface {
	//获取最新一条请求记录
	GetLastApplyLogId(applyUser, target string, tp int) (string, error)
	//保存加群/好友请求
	SaveApplyLog(*Apply) error
	//获取加群/好友请求
	GetApplyLogById(string) (*Apply, error)
	//获取加群/好友请求
	GetApplyLogByUserAndTarget(applyUser, target string, tp int) (*Apply, error)
	//更新请求信息
	UpdateApplyStateById(logId string, state int) error
}

type OrderCacheI interface {
	GetOrderById(orderId string) (*Order, error)
	SaveOrder(order *Order) error
}

type AccountCacheI interface {
	//保存用户信息
	SaveToken(token *user_model.Token) error
	//获取用户信息
	GetToken(appId, token string) (*user_model.Token, error)
}

//赞赏接口
type PariseCacheI interface {
	SaveLeaderBoard(tp int, l map[string]*RankingItem, startTime, endTime int64) error
	GetPraiseStatic(tp int, startTime, endTime int64) (map[string]*RankingItem, error)
}
