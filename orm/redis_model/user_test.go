package redis_model

import (
	"testing"

	"github.com/33cn/chat33/cache"
	"github.com/33cn/chat33/cache/redis"
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

	switch cfg.CacheType.CacheType {
	case "redis":
		redis.RegistRedis(&cfg)
	}
	cache.Cache, _ = cache.GetInstance(cfg.CacheType.CacheType)

}

func Test_GetUserInfoById(t *testing.T) {
	user, err := GetUserInfoById("1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func Test_GetUserInfoByField(t *testing.T) {
	user, err := GetUserInfoByField("1001", "uid", "171")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func Test_UpdateUid(t *testing.T) {
	err := UpdateUid("1", "34")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_UpdateDepositAddress(t *testing.T) {
	err := UpdateDepositAddress("1", "1111111")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_UpdateNowVersion(t *testing.T) {
	//err := UpdateNowVersion("1", "2.7.0")
	err := UpdateNowVersion("1", "2.7.1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_UpdatePublicKey(t *testing.T) {
	//err := UpdatePublicKey("1", "5555555")
	err := UpdatePublicKey("1", "03618b4c65758591046ee637fcab290469e5b24ae3e76afd4cedd41963666dfa8a", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_UpdateInviteCode(t *testing.T) {
	//err := UpdateInviteCode("1", "test invite code")
	err := UpdateInviteCode("1", "GvtLvd")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func Test_GetLoginLog(t *testing.T) {
	info, err := GetLoginLog("1", []string{"Android"})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}

func Test_RoomInviteConfirm(t *testing.T) {
	info, err := RoomInviteConfirm("15")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info)
}
