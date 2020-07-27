package db

import (
	"fmt"
	"strings"
	"testing"
)

func Test_GetLastUserLoginLog(t *testing.T) {
	maps, err := GetLastUserLoginLog("1", []string{"Android", "iOS"})
	if err != nil {
		t.Error(err)
	}
	t.Log(maps)
}

func Test_TTT(t *testing.T) {
	s := []string{"foo", "bar", "baz"}
	ss := strings.Join(s, "','")
	sss := fmt.Sprintf("'%s'", ss)
	t.Log(sss)
}

func Test_InsertUser(t *testing.T) {
	num, userId, err := InsertUser("1001123", "123", "1001", "测试名称1", "", "", "86", "15763946518", "1", "1", "", "", "Android", "2.5.1", 0)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(num)
	t.Log(userId)
}

func Test_UpdatePublicKey(t *testing.T) {
	err := UpdatePublicKey("1", "1231231231", "")
	if err != nil {
		t.Error(err)
	}
}
