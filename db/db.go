package db

import (
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/33cn/chat33/app/dao"
	"github.com/33cn/chat33/types"

	"strings"

	"github.com/33cn/chat33/pkg/btrade/common/mysql"
)

var conn *mysql.MysqlConn

func InitDB(cfg *types.Config) {
	dao.InitDB(cfg)

	c, err := mysql.NewMysqlConn(cfg.Mysql.Host, fmt.Sprintf("%v", cfg.Mysql.Port),
		cfg.Mysql.User, cfg.Mysql.Pwd, cfg.Mysql.Db, "UTF8MB4")
	if err != nil {
		log15.Error("mysql init failed", "err", err)
		panic(err)
	}
	conn = c
}

func GetConn() *mysql.MysqlConn {
	return conn
}

func GetNewTx() (*mysql.MysqlTx, error) {
	return conn.NewTx()
}

func QueryStr(src string) string {
	src = strings.Replace(src, `\`, `\\`, -1)
	src = strings.Replace(src, `'`, `\'`, -1)
	src = strings.Replace(src, `%`, `\%`, -1)
	src = strings.Replace(src, `_`, `\_`, -1)
	return src
}
