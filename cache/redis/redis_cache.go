package redis

import (
	. "github.com/33cn/chat33/cache"
	"github.com/33cn/chat33/types"
)

type redisCache struct {
	userCaChe
	friendCache
	roomCaChe
	applyCaChe
	orderCache
	accountCaChe
	praiseCaChe
}

/*func init() {
	configPath := "../../etc/config.toml"
	var cfg types.Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		panic(err)
	}
	if cfg.CacheType.CacheType == "redis" {
		Register("redis", &redisCache{})
	}
}*/

func RegistRedis(cfg *types.Config) {
	redisCfg := &nodeConfig{
		Url:         cfg.Redis.Url,
		Password:    cfg.Redis.Password,
		MaxIdle:     cfg.Redis.MaxIdle,
		MaxActive:   cfg.Redis.MaxActive,
		IdleTimeout: cfg.Redis.IdleTimeout,
	}
	Register("redis", &redisCache{})
	InitRedis(redisCfg)
}
