package db

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/types"
)

func Test_FindLastlogByUserId(t *testing.T) {
	maps, err := FindLastlogByUserId("1", types.HadRead)
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println(maps, err)
	fmt.Println("ok")
}

func Test_FindUserByMarkId(t *testing.T) {
	maps, err := FindUserByMarkId("1002", "8613037132696")
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println(maps, err)
	fmt.Println("ok")
}

func Test_FindUserByPhone(t *testing.T) {
	maps, err := FindUserByPhone("1001", "13037132696")
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println(maps, err)
	fmt.Println("ok")
}

func Test_FindUserByPhoneV2(t *testing.T) {
	maps, err := FindUserByPhoneV2("1001", "13037132696")
	if err != nil {
		t.Error("db connect fail")
		return
	}
	fmt.Println(maps, err)
	fmt.Println("ok")
}

func Test_FindSessionKeyAlert(t *testing.T) {
	maps, err := FindSessionKeyAlert("1", 1559719186000)
	if err != nil {
		t.Errorf("db connect fail:%s", err.Error())
		return
	}
	fmt.Println(maps, err)
	fmt.Println("ok")
}
