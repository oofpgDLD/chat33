package mysql_model

/*import (
	cmn "github.com/33cn/chat33/pkg/btrade/common"
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/33cn/chat33/cache"
	"github.com/33cn/chat33/cache/redis"
	"github.com/33cn/chat33/router"
)

var cfg *types.Config
var coinList []*types.Coin
var appList []*types.App

var Rc cache.Cache

func Init(c *types.Config) {
	cfg = c
	loadRoom()
	loadCoins()
	initRedis(cfg)
}

func SetCfg(c *types.Config) {
	cfg = c
}

func loadRoom() {
	rooms, err := GetEnableRoomsId()
	if err != nil {
		panic(err)
	}
	//load room
	for _, v := range rooms {
		channelId := types.GetRoomRouteById(utility.ToString(v))
		router.AppendChannel(channelId, router.NewChannel(channelId))
	}
}

func loadCoins() {
	coins, err := db.GetAllCoins()
	if err != nil {
		panic(err)
	}

	for _, c := range coins {
		one := &types.Coin{
			CoinId:   cmn.ToInt(c["coin_id"]),
			CoinName: c["coin_name"],
		}
		coinList = append(coinList, one)
	}
}

func initRedis(cfg *types.Config) {
	redis.RegistRedis(cfg)
	Rc, _ = cache.New("redis")
}*/
