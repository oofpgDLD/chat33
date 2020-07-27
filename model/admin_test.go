package model

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_SplitDay(t *testing.T) {
	ret := SplitDay(1554048000000, 1554912000000)
	fmt.Println(len(ret))
	fmt.Println(ret)
}

func Test_BanUser(t *testing.T) {
	err := BanUser("1001", "1", "1", "测试封用户", 7258089600000)
	t.Log(err)
}

func Test_BanUserCancel(t *testing.T) {
	err := BanUserCancel("1001", "1", "1")
	t.Log(err)
}

func Test_BanRoom(t *testing.T) {
	err := BanRoom("1001", "1", "21", "测试封群", 7258089600000)
	t.Log(err)
}

func Test_BanRoomCancel(t *testing.T) {
	err := BanRoomCancel("1001", "1", "21")
	t.Log(err)
}

func Test_ConvertClosedAlertStr(t *testing.T) {
	str := ConvertUserClosedAlertStr("1001", 1555752222000)
	t.Log(str)
}

func Test_ExportSumDetails(t *testing.T) {
	ret, err := ExportSumDetails("1001", 1553529600000, 1556121599000)
	t.Log(err)
	t.Log(ret)
}

func Test_CreateAd(t *testing.T) {
	ret, err := CreateAd("1001", "test1", "hhhhh", 3, "aaaa", 0)
	t.Log(err)
	t.Log(ret)
}

func Test_Advertisement(t *testing.T) {
	ret, err := Advertisement("1001")
	t.Log(ret)
	t.Log(err)
}

func Test_HHHAAA(t *testing.T) {
	params := url.Values{}
	/*	params.Set("search","1")
		params.Set("page","1")
		params.Set("size","1")
		params.Set("rewardType","1")*/

	t.Log(params.Encode())
}

func Test_AdminRoomsList(t *testing.T) {
	ret, err := AdminRoomsList("1006", []int{2}, "", 1, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}
