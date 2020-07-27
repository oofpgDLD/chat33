package db

import (
	"fmt"
	"testing"

	"github.com/33cn/chat33/types"
)

func Test_ReadSnapMsg(t *testing.T) {
	_, _, err := AlertRoomRevStateByRevId("121113", types.HadBurnt)
	fmt.Println("err:", err)
	fmt.Println("ok")
}
