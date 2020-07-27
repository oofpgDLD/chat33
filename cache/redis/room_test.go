package redis

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
	"github.com/garyburd/redigo/redis"
)

var conn redis.Conn
var c *redisCache

func init() {
	configPath := "../../etc/config.toml"
	var cfg types.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		panic(err)
	}
	redisCfg := &nodeConfig{
		Url:         cfg.Redis.Url,
		Password:    cfg.Redis.Password,
		MaxIdle:     cfg.Redis.MaxIdle,
		MaxActive:   cfg.Redis.MaxActive,
		IdleTimeout: cfg.Redis.IdleTimeout,
	}
	InitRedis(redisCfg)
	db.InitDB(&cfg)
	conn = GetConn()
	c = &redisCache{}
}

//
////储存群所有成员和cid
//func TestSaveAllMemberCid(t *testing.T) {
//	err := c.SaveAllMemberCid("43", &[]*types.RoomCid{&types.RoomCid{"1", "1"}, &types.RoomCid{"2", "2"}})
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
////判断用户是否在群
//func TestUserIsInRoom(t *testing.T) {
//	b1, b2, err := c.UserIsInRoom("43", "22")
//	fmt.Println(err)
//	fmt.Println(b1)
//	fmt.Println(b2)
//}
//
////获取所有群成员的cid
//func TestGetRoomUserCid(t *testing.T) {
//	count, err := c.GetRoomUserCid("43")
//	fmt.Println(err)
//	if count == nil {
//		fmt.Println(1111)
//	}
//	fmt.Println(count)
//}
//
////更新roomcid
//func TestUpdateRoomCid(t *testing.T) {
//	err := c.UpdateRoomCid("43", "12", "")
//	fmt.Println(err)
//}
//
////根据userid删除user-cid
//func TestDeleteRoomCidByUserId(t *testing.T) {
//	err := c.DeleteRoomCidByUserId("43", "12")
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
////删除roomcid
//func TestDeleteRoomCid(t *testing.T) {
//	err := c.DeleteRoomCid("43")
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
///////////////////////////////////////////////////////////////////////////
////保存群信息
//func TestSaveRoomInfo(t *testing.T) {
//	err := c.SaveRoomInfo(
//		&types.Room{
//			"19",
//			"137450112"	,
//			"测试群组1",
//			"1539153616197",
//			"1",
//			1,
//			1,
//			1,
//			1,
//			2,
//			2,
//			1,
//			1,
//			0,
//			1,
//			1,
//			1,
//			"",
//	})
//	if err != nil {
//		fmt.Println(err)
//	}
//}
////查询群详情
//func TestFindRoomInfo(t *testing.T) {
//	v, err := c.FindRoomInfo("2")
//	fmt.Println(err)
//	fmt.Println(v)
//}
//
////更新群信息
//func TestUpdateRoomInfo(t *testing.T) {
//	err := c.UpdateRoomInfo("1", "avatar", "2")
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
////删除群信息
//func TestDeleteRoomInfo(t *testing.T) {
//	err := c.DeleteRoomInfo("1")
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//////////////////////////////////////////////////////////
////保存群成员信息
//func TestSaveRoomUser(t *testing.T) {
//	err := c.SaveRoomUser("43", &[]*types.RoomMemberInfo{
//		&types.RoomMemberInfo{"43", "1", "aaa", 1, "1"},
//		&types.RoomMemberInfo{"43", "2", "a2a", 2, "2"}})
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
////查询群成员信息
//func TestFindRoomMemberInfo(t *testing.T) {
//	count, err := c.FindRoomMemberInfo("43", "1")
//	fmt.Println(err)
//	fmt.Println(count)
//}
//
////更新群成员信息
//func TestUpdateRoomUserInfo(t *testing.T) {
//	err := c.UpdateRoomUserInfo("43", "1", "noDisturbing", "11")
//	fmt.Println(err)
//}
//
////删除群成员信息
//func TestDeleteRoomUserInfo(t *testing.T) {
//	err := c.DeleteRoomUserInfo("43", "1")
//	fmt.Println(err)
//}
//
////////////////////////////////////////////
////存禁言用户
//func TestSaveRoomUserMued(t *testing.T) {
//	err := c.SaveRoomUserMued("1", &[]*types.RoomUserMued{&types.RoomUserMued{"1", 1, 100},
//		&types.RoomUserMued{"2", 2, 100},
//		&types.RoomUserMued{"3", 3, 100}})
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
////更新/添加群禁言信息
//func TestUpdateRoomUserMuted(t *testing.T) {
//	err := c.UpdateRoomUserMuted("1", &types.RoomUserMued{"45", 1, 100})
//	fmt.Println(err)
//}
//
////获取群中禁言详情
//func TestGetRoomUserMuted(t *testing.T) {
//	ret, err := c.GetRoomUserMuted("1", "2")
//	fmt.Println(err)
//	fmt.Println(ret)
//}
//
////删除群禁言信息
//func TestDeleteRoomUserMuted(t *testing.T) {
//	err := c.DeleteRoomUserMuted("1")
//	fmt.Println(err)
//}
//
//func TestRedis1(t *testing.T) {
//	conn := GetConn()
//	ret, err := conn.Do("del", "r-")
//	fmt.Println(err)
//	fmt.Println(ret)
//}

//存储群消息记录
func TestSaveRoomMsgContent(t *testing.T) {
	test := &types.RoomLog{
		"119",
		"1",
		"16",
		"4",
		2,
		1,
		"context:msg",
		1539585971497,
		1,
	}
	err := c.SaveRoomMsgContent(test)
	fmt.Println(err)
}

//根据logId获取消息记录
func TestGetRoomMsgContent(t *testing.T) {
	info, err := c.GetRoomMsgContent("118")
	fmt.Println(info)
	fmt.Println(err)
}

//获取所有消息记录
func TestGetRoomMsgContentAll(t *testing.T) {
	info, err := c.GetRoomMsgContentAll()
	for _, v := range info {
		fmt.Println(v)
	}
	fmt.Println(info)
	fmt.Println(err)
}

//根据id删除消息记录
func TestDeleteRoomMsgContent(t *testing.T) {
	err := c.DeleteRoomMsgContent("119")
	fmt.Println(err)
}

//存储群recv消息记录
func TestSaveReceiveLog(t *testing.T) {
	err := c.SaveReceiveLog(&types.RoomMsgReceive{
		"1",
		"123",
		"7",
		1,
	})
	fmt.Println(err)
}

//根据Id获取消息记录
func TestGetReceiveLogbyId(t *testing.T) {
	info, err := c.GetReceiveLogbyId("1")
	fmt.Println(info)
	fmt.Println(err)
}

//获取所有recv消息记录
func TestGetReceiveLogAll(t *testing.T) {
	info, err := c.GetReceiveLogAll()
	for _, v := range info {
		fmt.Println(v)
	}
	fmt.Println(info)
	fmt.Println(err)
}

//更新recv消息状态
func TestUpdateReceiveLog(t *testing.T) {
	err := c.UpdateReceiveLog("1", 2)
	fmt.Println(err)
}

//根据id删除recv消息记录
func TestDeleteReceiveLog(t *testing.T) {
	err := c.DeleteReceiveLog("1")
	fmt.Println(err)
}

////**********************群和群成员的关系*******************//
//群和成员的关系
func TestAddRoomUser(t *testing.T) {
	info := make([]*types.RoomMember, 0)
	member1 := &types.RoomMember{
		"21",
		"19",
		"1",
		"测试群1",
		3,
		1,
		1,
		1,
		1539153616197,
		1,
		"",
	}
	member2 := &types.RoomMember{
		"22",
		"19",
		"4",
		"测试群2",
		2,
		1,
		1,
		1,
		1539153616197,
		1,
		"",
	}
	member3 := &types.RoomMember{
		"23",
		"19",
		"5",
		"测试群3",
		2,
		1,
		1,
		1,
		1539153616197,
		1,
		"",
	}
	info = []*types.RoomMember{member1, member2, member3}
	err := c.AddRoomUser("19", &info)
	fmt.Println(err)
}

//获取群和成员关系 map[userid]{level}
func TestGetRoomUser(t *testing.T) {
	maps, err := c.GetRoomUser("19")
	for k, v := range maps {
		fmt.Println(k)
		fmt.Println(v)
	}
	fmt.Println(maps)
	fmt.Println(err)
}

//更新群和成员关系
func TestUpdateRoomUser(t *testing.T) {
	err := c.UpdateRoomUser("19", "1", "1")
	fmt.Println(err)
}

//删除群和成员关系(群没被删除，只是成员被删除)
func TestDelRoomUser(t *testing.T) {
	err := c.DelRoomUser("19", "1")
	fmt.Println(err)
}

//删除群和成员关系(群被删除的情况)
func TestDelRoomUserAll(t *testing.T) {
	err := c.DelRoomUserAll("19")
	fmt.Println(err)
}

//**********************room_config 和 user_config*******************//
//新增或更新room_config
func TestSaveRoomConfig(t *testing.T) {
	err := c.SaveRoomConfig("1001", 1, 0)
	fmt.Println(err)
}

//查找room_config
func TestFindRoomConfig(t *testing.T) {
	info, err := c.FindRoomConfig("1001")
	fmt.Println(info)
	fmt.Println(err)
}

//新增或更新user_config
func TestSaveUserConfig(t *testing.T) {
	err := c.SaveUserConfig("1001", 1, 0)
	fmt.Println(err)
}

//获取user_config
func TestFindUserConfig(t *testing.T) {
	info, err := c.FindUserConfig("1001")
	fmt.Println(info)
	fmt.Println(err)
}
