package redis_model

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/types"
)

func Test_FindApplyLogById(t *testing.T) {
	apply, err := FindApplyLogById(973)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(apply)
}

func Test_FindApplyLogByUserAndTarget(t *testing.T) {
	apply, err := FindApplyLogByUserAndTarget("333", "8", types.IsFriend)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(apply)
}

func Test_GetAddFriendApplyLog(t *testing.T) {
	apply, err := GetAddFriendApplyLog("1", "7")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(apply)
}

func Test_UpdateApply(t *testing.T) {
	id, err := UpdateApply("1", "334", types.IsFriend, "", "test", `{"sourceType":1,"sourceId":""}`, types.AwaitState)
	if err != nil {
		t.Error(err)
	}
	t.Log(fmt.Sprintf("success update log %v", id))
}

func Test_AcceptFriendApply(t *testing.T) {
	num, err := AcceptFriendApply("1", "334")
	if err != nil {
		t.Error(err)
	}
	t.Log(fmt.Sprintf("success update affect numbers %v", num))
}

func Test_AppendApplyLog(t *testing.T) {
	id, err := AppendApplyLog("1", "334", "", `{"sourceType":1,"sourceId":""}`, "", types.AwaitState, types.IsFriend)
	if err != nil {
		t.Error(err)
	}
	t.Log(fmt.Sprintf("success insert id %v", id))
}
