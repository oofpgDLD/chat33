package mysql_model

import (
	"testing"

	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	configPath := "../../etc/config.toml"
	var cfg types.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		panic(err)
	}

	db.InitDB(&cfg)
}

func Test_FindNotBurnedLogsAfter(t *testing.T) {
	logs, err := FindNotBurnedLogsAfter("1", types.FriendMsgIsNotDelete, 1564628844000) // 2019-08-01 11:07:24
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("length :%v", len(logs))
}

func Test_FindNotBurnedLogsBetween(t *testing.T) {
	logs, err := FindNotBurnedLogsBetween("1", types.FriendMsgIsNotDelete, 1565924948760, 1565924949393)
	if err != nil {
		t.Error(err)
		return
	}
	for i, l := range logs {
		t.Logf("%v : %v ,%v \n", i, l.PrivateLog, l.User)
	}
	t.Logf("length :%v", len(logs))
}

func Test_FindSessionKeyAlert(t *testing.T) {
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

func Test_FindFriendById(t *testing.T) {
	info, err := FindFriendById("142", "8")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info.User, info.Friend)
}
