package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/33cn/chat33/api"
	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/pkg/excRate"
	"github.com/33cn/chat33/pkg/work"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/pprof"
	l "github.com/inconshreveable/log15"
)

func initLogLevel(cfg *types.Config) {
	var level l.Lvl
	switch cfg.Log.Level {
	case "debug":
		level = l.LvlDebug
	case "info":
		level = l.LvlInfo
	case "warn":
		level = l.LvlWarn
	case "error":
		level = l.LvlError
	case "crit":
		level = l.LvlCrit
	default:
		level = l.LvlWarn
	}
	l.Root().SetHandler(l.LvlFilterHandler(level, l.StreamHandler(os.Stdout, l.TerminalFormat())))
}

func findConfiger() (string, error) {
	var configPath = ""
	l.Info("runtime:", "os", runtime.GOOS)
	if runtime.GOOS == `windows` {
		configPath = "etc/config.toml"
	} else {
		err := os.Chdir(pwd())
		if err != nil {
			l.Info("get project pwd err", "err", err)
			return configPath, err
		}
		d, _ := os.Getwd()
		l.Info("project info:", "dir", d)
		configPath = d + "/etc/config.toml"
		types.ExcelAddr = d + "/excel/"
	}
	return configPath, nil
}

func main() {
	var cfg types.Config
	configPath, err := findConfiger()
	if err != nil {
		return
	}
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		panic(err)
	}
	l.Info("config info", "cfg", cfg)

	work.Init(&cfg.Work)
	initLogLevel(&cfg)
	db.InitDB(&cfg)
	model.Init(&cfg)
	//redis 每日存储榜单数据
	model.SaveDaily()
	excRate.Init(&cfg)
	r := api.Init(&cfg)
	pprof.Register(r)
	l.Info("init success")
	err = r.Run(cfg.Server.Addr)
	if err != nil {
		l.Error("main run", "err", err)
	}
}

/*
	---workdir/
		| -- bin/
		|     |-- chat(I am here)
		|
		| -- etc/
			  |-- config.toml
			  |-- config.json
*/
func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(filepath.Dir(os.Args[0])))
	if err != nil {
		panic(err)
	}
	return dir
}
