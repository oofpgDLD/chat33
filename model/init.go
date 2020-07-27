package model

import (
	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/cache"
	"github.com/33cn/chat33/cache/redis"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"

	logic "github.com/33cn/chat33/router"
)

var cfg *types.Config

func Init(c *types.Config) {
	cfg = c
	loadRoom()
	initCache(cfg)
	orm.Init(c)
	app.UpdateAppsConfig()
	app.UpdateVersionConfig()
	//init batch config
	batchConfig = &cfg.BatchConfig
	StartServe()
}

func SetCfg(c *types.Config) {
	cfg = c
}

func loadRoom() {
	rooms, err := orm.GetEnableRoomsId()
	if err != nil {
		panic(err)
	}
	//load room
	for _, v := range rooms {
		channelId := types.GetRoomRouteById(utility.ToString(v))
		logic.AppendChannel(channelId, logic.NewChannel(channelId))
	}
}

func initCache(cfg *types.Config) {
	switch cfg.CacheType.CacheType {
	case "redis":
		redis.RegistRedis(cfg)
	}
	cache.Cache, _ = cache.GetInstance(cfg.CacheType.CacheType)
}
