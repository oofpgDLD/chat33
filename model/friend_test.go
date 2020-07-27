package model

import (
	"testing"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func Test_FindTypicalChatLog(t *testing.T) {
	ret, err := FindTypicalChatLog("8", "8", "", "", "", 20, []string{utility.ToString(types.Photo), utility.ToString(types.Video)})
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func Test_FindChatLog(t *testing.T) {
	ret, err := FindCatLog("1001", "8", "1", "", 0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}

func Test_FriendInfo(t *testing.T) {
	ret, err := FriendInfo("1", "1")
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
