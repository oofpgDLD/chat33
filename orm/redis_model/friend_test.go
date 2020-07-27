package redis_model

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/cache"

	mysql "github.com/33cn/chat33/orm/mysql_model"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func Test_FindFriendById(t *testing.T) {
	info, err := FindFriendById("1", "8")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(fmt.Sprintf("friend info %v", info.Friend))
	t.Log(fmt.Sprintf("friend info detail %v", info.User))
}

func Test_FindFriends(t *testing.T) {
	friendsId, err := FindFriends("1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(fmt.Sprintf("friends id %v", friendsId))
}

func Test_DeleteFriend(t *testing.T) {
	err := DeleteFriend("1", "6", utility.NowMillionSecond())
	if err != nil {
		t.Error(err)
	}
}

func Test_AcceptFriend(t *testing.T) {
	err := AcceptFriend("8", "1", `0`, utility.NowMillionSecond())
	if err != nil {
		t.Error(err)
	}
}

func Test_FindFriendsFilterByCommonUse(t *testing.T) {
	list, err := FindFriendsFilterByCommonUse("1", types.CommonUseFindAll)
	if err != nil {
		t.Error(err)
		return
	}
	for _, l := range list {
		t.Log(fmt.Sprintf("%v %v", l.Friend, l.User))
	}
}

/*func Test_FindFriendsAfterTime(t *testing.T) {
	list, err := FindFriendsAfterTime("1", types.UncommonUse, utility.TimeStrToMillionSecond("2021-01-01 00:00:01"))
	if err != nil {
		t.Error(err)
		return
	}
	for _, l := range list {
		t.Log(fmt.Sprintf("%v %v", l.Friend, l.User))
	}
}*/

func Test_CheckIsFriend(t *testing.T) {
	b, err := CheckIsFriend("1", "8")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(b)
}

func Test_SetFriendRemark(t *testing.T) {
	//err := SetFriendRemark("1", "8", "test remark")
	err := SetFriendRemark("1", "8", "")
	if err != nil {
		t.Error(err)
	}
}

func Test_SetFriendExtRemark(t *testing.T) {
	//err := SetFriendExtRemark("1", "8", "test ext remark")
	err := SetFriendExtRemark("1", "8", `{"telephones":[],"description":"123","pictures":["p1","p2","p3"]}`)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetFriendDND(t *testing.T) {
	err := SetFriendDND("1", "8", types.NoDisturbingOn)
	//err := SetFriendDND("1", "8", types.NoDisturbingOff)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetFriendIsTop(t *testing.T) {
	//err := SetFriendIsTop("1", "8", types.NotOnTop)
	err := SetFriendIsTop("1", "8", types.OnTop)
	if err != nil {
		t.Error(err)
	}
}

func Test_FindAddFriendConfByUserId(t *testing.T) {
	conf, err := FindAddFriendConfByUserId("8")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(conf)
}

func Test_SetQuestionandAnswer(t *testing.T) {
	err := SetQuestionandAnswer("8", "我叫什么名字？", "郑家烨")
	//err := SetQuestionandAnswer("8", "我是谁？", "猥琐男")
	if err != nil {
		t.Error(err)
	}
}

func Test_SetNeedAnswer(t *testing.T) {
	err := SetNeedAnswer("8", "我叫什么名字？", "郑家烨")
	//err := SetQuestionandAnswer("8", "我是谁？", "猥琐男")
	if err != nil {
		t.Error(err)
	}
}

func Test_SetNotNeedAnswer(t *testing.T) {
	err := SetNotNeedAnswer("8")
	if err != nil {
		t.Error(err)
	}
}

func Test_IsNeedConfirm(t *testing.T) {
	err := IsNeedConfirm("8", types.NeedConfirm)
	//err := IsNeedConfirm("8", types.NotNeedConfirm)
	if err != nil {
		t.Error(err)
	}
}

func loadAllPrivteLog() error {
	logs, err := mysql.FindAllPrivateLogs()
	if err != nil {
		return err
	}

	err = cache.Cache.SavePrivateChatLogs(logs)
	return err
}

func Test_LoadPrivateLog(t *testing.T) {
	err := loadAllPrivteLog()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("load success")
}

func Test_FindPrivateChatLogById(t *testing.T) {
	ret, err := FindPrivateChatLogById("44639")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("private chat info %v, %v", ret.PrivateLog, ret.User)
}

func Test_UpdatePrivateLogContentById(t *testing.T) {
	//err := UpdatePrivateLogContentById("44639", `{"content":"hahaha"}`)
	err := UpdatePrivateLogContentById("44639", `{"content":"服用"}`)
	if err != nil {
		t.Error(err)
	}
}

func Test_FindNotBurnedLogsAfter(t *testing.T) {
	err := loadAllPrivteLog()
	if err != nil {
		t.Error(err)
		return
	}

	logs, err := FindNotBurnedLogsAfter("1", types.FriendMsgIsNotDelete, 1564628844000) // 2019-08-01 11:07:24
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("length :%v", len(logs))
}

func Test_FindNotBurnedLogsBetween(t *testing.T) {
	err := loadAllPrivteLog()
	if err != nil {
		t.Error(err)
		return
	}

	logs, err := FindNotBurnedLogsBetween("1", types.FriendMsgIsNotDelete, 1565924948760, 1565924949393) // 2019-08-01 11:07:24
	if err != nil {
		t.Error(err)
		return
	}
	for i, l := range logs {
		t.Logf("%v : %v ,%v \n", i, l.PrivateLog, l.User)
	}
	t.Logf("length :%v", len(logs))
}

func Test_UpdatePrivateLogStateById(t *testing.T) {
	err := UpdatePrivateLogStateById("44583", types.NotRead) // 2019-08-01 11:07:24
	if err != nil {
		t.Error(err)
	}
}

func Test_FindSessionKeyAlert(t *testing.T) {
	/*err := loadAllPrivteLog()
	if err != nil {
		t.Error(err)
		return
	}*/

	logs, err := FindSessionKeyAlert("1", 1566200773000) // 2019-08-01 11:07:24
	if err != nil {
		t.Error(err)
		return
	}
	/*for i,l := range logs{
		t.Logf("%v : %v ,%v \n", i, l.PrivateLog, l.User)
	}*/
	t.Logf("length :%v", len(logs))
}

/*func Test_Regexp(t *testing.T) {
	src := `{"fromKey":"03e178ca30651d206b4afcaeb57737c9b90567d0792b9e03e0103ff3b3424fa325","key":"5e7e4497dfef95f1dac2a53976c91c744e986ba7f1e0f1067cbca67d7d269a3e4b514ed176536ac580e5618e667d57fdbc9637655fe751f43be57f7008d6eb72312e5e1e19f3241bdc29c8e6265040fd7c5d8bfe704204ff1481f7d6","kid":"1561546733952","roomId":"878","type":19}`
	b, err := regexp.MatchString(`"type":19`, src)
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(b)
	return
}*/
