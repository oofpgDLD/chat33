package redis_model

import (
	"testing"
)

func TestGetRoomUserLevel(t *testing.T) {
	level := GetRoomUserLevel("19", "1", 1)
	t.Log("level:", level)
}

/*func TestAddMember(t *testing.T) {
	err := AddMember("19", "5", 1)
	if err != nil {
		t.Error(err)
		return
	}
}*/
func TestGetRoomMemberName(t *testing.T) {
	s := GetRoomMemberName("20", "4")
	t.Log(s)
}

//set三种permission,写在一起
func TestSetPermission(t *testing.T) {
	err := SetPermission("19", 1, 1, 2, 99)
	if err != nil {
		t.Error(err)
		return
	}
}

//设置群名
func TestSetRoomName(t *testing.T) {
	_, err := SetRoomName("19", "cs", 99)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetAvatar(t *testing.T) {
	_, err := SetAvatar("19", "https://zb-chat.oss-cn-shanghai.aliyuncs.com/chatList/picture/20181101/20181101143106226_8.jpg", 99)
	if err != nil {
		t.Error(err)
		return
	}
}

//设置成员等级
func TestSetMemberLevel(t *testing.T) {
	err := SetMemberLevel("1", "19", 1, 2)
	if err != nil {
		t.Error(err)
		return
	}
}

//转让群主  master原群主id userid新群主id
func TestSetNewMaster(t *testing.T) {
	err := SetNewMaster("1", "4", "19", 3)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetNoDisturbing(t *testing.T) {
	err := SetNoDisturbing("1", "19", 2, 99)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetOnTop(t *testing.T) {
	err := SetOnTop("1", "19", 2, 99)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetMemberNickname(t *testing.T) {
	err := SetMemberNickname("1", "19", "测试群1", 99)
	if err != nil {
		t.Error(err)
		return
	}
}

//新增成员禁言
func TestAddMutedMember(t *testing.T) {
	err := AddMutedMember(1, "194", "15", 2, 0)
	if err != nil {
		t.Error(err)
		return
	}
}

//取消某个成员禁言，其实是更新那个成员禁言状态
func TestDelMemberMuted(t *testing.T) {
	err := DelMemberMuted("124", "1")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetRoomMutedType(t *testing.T) {
	info, err := GetRoomMutedType("19", 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}

func TestSetRoomMutedType(t *testing.T) {
	err := SetRoomMutedType("21", 1, 1)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetRoomUserMuted(t *testing.T) {
	mutedType, deadline := GetRoomUserMuted("19", "1")
	t.Log(mutedType, deadline)
}

//获取某种禁言信息
func TestGetMutedListByType(t *testing.T) {
	info, err := GetMutedListByType("194", 1)
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range info {
		t.Log(k, v)
	}

}

//查找recv群消息
func TestFindReceiveLogById(t *testing.T) {
	info, err := FindReceiveLogById("123", "7")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info.Id, info.RoomMsgId, info.ReceiveId)
}

//添加群聊接收日志
func TestAppendMemberRevLog(t *testing.T) {
	err := AppendMemberRevLog(1, "7", "123", 1)
	if err != nil {
		t.Error(err)
		return
	}
}

//获取群中所有管理员
func TestGetRoomManagerAndMaster(t *testing.T) {
	info, err := GetRoomManagerAndMaster("19")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range info {
		t.Log(v.RoomMember)
		t.Log(v.User)
	}

}

//根据id查询群成员信息
func TestFindRoomMemberById(t *testing.T) {
	info, err := FindRoomMemberById("19", "1", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info.User, info.RoomMember)
}

//查找群中消息免打扰的成员
func TestFindSetNoDisturbingMembers(t *testing.T) {
	info, err := FindSetNoDisturbingMembers("19")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range info {
		t.Log(v)
	}
}

/*func TestCreateNewRoom(t *testing.T) {
	err := CreateNewRoom("1", 20, 2, 1)
	if err != nil {
		t.Error(err)
		return
	}
}*/

//根据room id获取群信息
func TestSearchRoomInfo(t *testing.T) {
	info, err := SearchRoomInfo("19", 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}

//删除群
func TestDeleteRoomInfoById(t *testing.T) {
	err := DeleteRoomInfoById("19")
	if err != nil {
		t.Error(err)
		return
	}
}

//删除群成员
func TestDeleteRoomMemberById(t *testing.T) {
	_, err := DeleteRoomMemberById("1", "20")
	if err != nil {
		t.Error(err)
		return
	}
}

//根据logid查群聊天信息
func TestFindRoomChatLogByContentId(t *testing.T) {
	info, err := FindRoomChatLogByContentId("118")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info.User)
	t.Log(info.RoomLog)
}

//根据senderid和msgid查群聊天信息
func TestFindRoomChatLogByMsgId(t *testing.T) {
	info, err := FindRoomChatLogByMsgId("4", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info.User)
	t.Log(info.RoomLog)
}

//根据id删除群聊天信息
func TestDelRoomChatLogById(t *testing.T) {
	err := DelRoomChatLogById("118")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAlertRoomRevStateByRevId(t *testing.T) {
	err := AlertRoomRevStateByRevId("1", 1)
	if err != nil {
		t.Error(err)
		return
	}
}

// 添加群聊聊天日志
func TestAppendRoomChatLog(t *testing.T) {
	err := AppendRoomChatLog(
		118,
		"4",
		"16",
		"1",
		1,
		2,
		"content:msg",
		"",
		1539585971497)
	if err != nil {
		t.Error(err)
		return
	}
}

//获取用户创建群个数上限
func TestGetCreateRoomsLimit(t *testing.T) {
	info, err := GetCreateRoomsLimit("1004", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}

//获取群的成员数上限
func TestGetRoomMembersLimit(t *testing.T) {
	info, err := GetRoomMembersLimit("1004", 1)

	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}

//设置用户创建群个数上限
func TestSetCreateRoomsLimit(t *testing.T) {
	err := SetCreateRoomsLimit("1001", 1, 100)
	if err != nil {
		t.Error(err)
		return
	}
}

//设置群的成员数上限
func TestSetRoomMembersLimit(t *testing.T) {
	err := SetRoomMembersLimit("1001", 1, 100)
	if err != nil {
		t.Error(err)
		return
	}
}

//设置为加v认证群
func TestSetRoomVerifyed(t *testing.T) {
	err := SetRoomVerifyed("19", "1")
	if err != nil {
		t.Error(err)
		return
	}
}
