package model

import (
	"testing"

	"github.com/33cn/chat33/db"
	"github.com/33cn/chat33/types"
	"github.com/BurntSushi/toml"
)

func init() {
	var cfg types.Config
	if _, err := toml.DecodeFile("../etc/config.toml", &cfg); err != nil {
		panic(err)
	}
	db.InitDB(&cfg)
}

/*func Test_TokenLogin(t *testing.T) {
	info, err := TokenLogin("1002", "1cd71990d07c43059cb53a1e1e237102", "iOS", "iPhone 6", "2", "", "96F4810A-EBBC-48F1-8FCF-039CB0F4500E20190307171629")
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}*/

func Test_CheckIsBlocked(t *testing.T) {
	b, err := CheckIsBlocked("8", "142")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(b)
}
