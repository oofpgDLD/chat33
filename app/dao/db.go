package dao

import (
	"fmt"

	"github.com/33cn/chat33/pkg/btrade/common/mysql"
	"github.com/33cn/chat33/types"
	"github.com/inconshreveable/log15"
)

var conn *mysql.MysqlConn

func InitDB(cfg *types.Config) {
	c, err := mysql.NewMysqlConn(cfg.Mysql.Host, fmt.Sprintf("%v", cfg.Mysql.Port),
		cfg.Mysql.User, cfg.Mysql.Pwd, cfg.Mysql.Db, "UTF8MB4")
	if err != nil {
		log15.Error("mysql init failed", "err", err)
		panic(err)
	}
	conn = c
}
