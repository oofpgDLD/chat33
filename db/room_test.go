package db

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	var cfg types.Config
	if _, err := toml.DecodeFile("../etc/config.toml", &cfg); err != nil {
		panic(err)
	}
	InitDB(&cfg)
}

func Test_GetRoomMembers(t *testing.T) {
	maps, err := GetRoomMembers("22", -1)
	if err != nil {
		t.Error("db connect fail")
		return
	}
	if len(maps) < 1 {
		fmt.Println("not find room info")
		return
	}
	for i, v := range maps {
		fmt.Println("number", i)
		fmt.Println("roomUserInfo:", v)
	}
	fmt.Println("ok")
}

func Test_GetRoomUserMuted(t *testing.T) {
	maps, err := GetRoomUserMuted("22", "16")
	if err != nil {
		t.Error("db connect fail")
		return
	}
	if len(maps) < 1 {
		fmt.Println("not find room info")
		return
	}
	for i, v := range maps {
		fmt.Println("number", i)
		fmt.Println("roomUserMutedInfo:", v)
	}
	fmt.Println("ok")
}

func Test_GetRoomMemberNumberByLevel(t *testing.T) {
	number, err := GetRoomMemberNumberByLevel("22", types.RoomLevelManager)
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println("number", number)
	fmt.Println("ok")
}

func Test_GetSystemMsg(t *testing.T) {
	maps, nextLog, err := GetSystemMsg("186", 2999, 10)
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println("nextLog:", nextLog)
	if len(maps) < 1 {
		fmt.Println("not find room info")
		return
	}
	for i, v := range maps {
		fmt.Println("number", i)
		fmt.Println("roomSystemInfo:", v)
	}
	fmt.Println("ok")
}

func Test_GetRoomMutedListNumber(t *testing.T) {
	maps, err := GetRoomMutedListNumber("1")
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println("number", maps)
	fmt.Println("ok")
}

func Test_AddRoomUserMuted(t *testing.T) {
	tx, err := GetNewTx()
	if err != nil {
		t.Error("db connect fail")
		return
	}
	mun, mmm, err := AddRoomUserMuted(tx, "124", "10", types.AllSpeak, 0)
	if err != nil {
		tx.RollBack()
		t.Error("db connect fail")
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Error("tx commit fail")
		return
	}
	fmt.Println(mun, mmm, err)
	fmt.Println("ok")
}

func Test_GetRoomLastContentByUserId(t *testing.T) {
	args := []string{}
	maps, err := GetRoomMsgContentAfter(args, 0, 1)
	if err != nil {
		t.Error("db connect fail")
		return
	}
	for _, v := range maps {
		fmt.Println(utility.MillionSecondToTimeString(utility.ToInt64(v["datetime"])))
	}
	fmt.Println(len(maps), err)
	fmt.Println("ok")
}

func Test_GetChatlog(t *testing.T) {
	maps, lastLog, err := GetChatlog("43", 0, 0, 10)
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println(maps, lastLog, err)
	fmt.Println("ok")
}

func Test_AddMember(t *testing.T) {
	tx, err := GetNewTx()
	if err != nil {
		panic(err)
	}
	b, _, err := RoomAddMember(tx, "1", "43", types.RoomLevelNomal, utility.NowMillionSecond(), "")
	if err != nil {
		tx.RollBack()
		t.Error(err)
		return
	}
	err = tx.Commit()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(b)
}
