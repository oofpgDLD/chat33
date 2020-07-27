package orm

import (
	"testing"

	"github.com/33cn/chat33/cache"
	"github.com/33cn/chat33/cache/redis"
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	configPath := "../etc/config.toml"
	var cfg types.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		panic(err)
	}

	Init(&cfg)
	db.InitDB(&cfg)

	switch cfg.CacheType.CacheType {
	case "redis":
		redis.RegistRedis(&cfg)
	}
	cache.Cache, _ = cache.GetInstance(cfg.CacheType.CacheType)
}

func Test_GetUserInfoByUid(t *testing.T) {
	user, err := GetUserInfoByUid("1001", "34")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}
