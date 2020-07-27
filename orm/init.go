package orm

import (
	"github.com/33cn/chat33/types"
)

var cfg *types.Config

func Init(c *types.Config) {
	cfg = c
}
