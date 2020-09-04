package types

import (
	work_model "github.com/33cn/chat33/pkg/work/model"
)

type Config struct {
	Loglevel        string
	LogFile         string
	Server          Server
	Mysql           Mysql
	Service         Service
	Limit           Limit
	Log             Log
	CacheType       CacheType
	Redis           Redis
	NotConfirmAppId NotConfirmAppId
	FileStore       FileStore
	BatchConfig     BatchConfig
	Work            work_model.Config
	Rate            Rate
	SMS             SMS
	Email           Email
	Env             Env
}

type Mysql struct {
	Host string
	Port int32
	User string
	Pwd  string
	Db   string
}

type Server struct {
	Addr string
}

type Service struct {
	DockingDeadline int64
}

type Limit struct {
	ChatPool int
	ChatRate int
}

type Log struct {
	Level string
}

type Coin struct {
	CoinId   int    `json:"coin_id"`
	CoinName string `json:"coin_name"`
}

type LastLoginInfo struct {
	Uuid       string
	Device     string
	DeviceName string
	LoginType  int
	LoginTime  int64
}

type CacheType struct {
	Enable    bool
	CacheType string
}

//redis的一些配置
type Redis struct {
	Url         string `json:"url"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
}

type NotConfirmAppId struct {
	AppIds string `json:"appIds"`
}

type FileStore struct {
	Path     string
	FilePath string
}

type BatchConfig struct {
	//new device batch day
	BatchDayAgo int
	//pack length
	BatchPackLength int
	//batch push interval: Millisecond
	BatchInterval int64
}

type Rate struct {
	TradeUrl string
	UsdtUrl  string
}

type SMS struct {
	Surl     string
	Curl     string
	CodeType string
	Msg      string
}

type Email struct {
	Surl     string
	Curl     string
	CodeType string
	Msg      string
}

type Env struct {
	Env   string
	Super map[string]string
}
