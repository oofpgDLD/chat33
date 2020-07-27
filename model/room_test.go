package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	var cfg types.Config
	if _, err := toml.DecodeFile("../etc/config.toml", &cfg); err != nil {
		panic(err)
	}
	Init(&cfg)
}

//---------------------------------------------------------------------//
func Test_getRoomMembers(t *testing.T) {
	maps, errMsg := getRoomMembers("22", -1)
	fmt.Println("err:", errMsg)
	for i, v := range maps {
		fmt.Println("numnber:", i)
		fmt.Println("room info :", v)
	}
	fmt.Println("ok")
}

func Test_SetSystemMsg(t *testing.T) {
	errMsg := SetSystemMsg("8", "22", "123")
	fmt.Println("err:", errMsg)
	fmt.Println("ok")
}

func Test_SetRoomMutedSingle(t *testing.T) {
	errMsg := SetRoomMutedSingle("7", "194", "13", 0)
	fmt.Println("err:", errMsg)
	fmt.Println("ok")
}

func Test_AdminSetPermission(t *testing.T) {
	errMsg := AdminSetPermission("256", 0, 0, 2)
	fmt.Println("err:", errMsg)
	fmt.Println("ok")
}

func Test_JoinRoomApply(t *testing.T) {
	_, errMsg := JoinRoomApply("1001", "1", "49", "", 1, "")
	fmt.Println("err:", errMsg)
	fmt.Println("ok")
}

func Test_UserIsInRoom(t *testing.T) {
	info, errMsg := UserIsInRoom("1", "49")
	fmt.Println("info", info, "err:", errMsg)
	fmt.Println("ok")
}

func Test_GetRoomUserInfo(t *testing.T) {
	ret, err := GetRoomUserInfo("22", "1")
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func Test_GetRoomInfo(t *testing.T) {
	ret, err := GetRoomInfo("62", "596")
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func Test_GetLastDeviceLoginInfo(t *testing.T) {
	ret, ok := GetLastDeviceLoginInfo("1", "Android")
	t.Log(ok)
	t.Log(ret)
}

func testViewRoomChatLog(t *testing.T, caller, roomId, nextLog string, number int) {
	ret, err := GetRoomChatLog(caller, roomId, nextLog, number)
	if err != nil {
		t.Error(err)
		return
	}
	if ret != nil {
		f := ret.(map[string]interface{})
		logs := f["logs"].([]*types.ChatLog)
		for _, r := range logs {
			b, err := json.Marshal(r)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(string(b))
		}
	}
}

func Test_CreateRoom(t *testing.T) {
	info, err := CreateRoom(
		"1001",
		"1",
		"",
		"shanghai.aliyuncs.com/chatList/picture/20181101/20181101143106226_8.jpg",
		1,
		1,
		1,
		1,
		1,
		1,
		[]string{"1", "4"})
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}

func Test_GetRoomChatLog(t *testing.T) {
	testViewRoomChatLog(t, "4831", "1838", "70867", 1)
	testViewRoomChatLog(t, "14104", "2093", "70918", 1)
	testViewRoomChatLog(t, "2166", "1446", "70948", 1)
	testViewRoomChatLog(t, "4831", "662", "70999", 1)
	testViewRoomChatLog(t, "4831", "1838", "71052", 1)
}
