package model

import (
	"testing"
	//"github.com/33cn/chat33/types"
	//"github.com/33cn/chat33/utility"
	"fmt"
	"time"
)

func Test_LeaderBoard(t *testing.T) {
	ret, err := LeaderBoard("6", 1, 1576389262000, 1576389262000, 1, 20)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
func Test_statisticsBoard(t *testing.T) {
	ret, err := statisticsBoard(1, 1574611200000, 1575129600000, []string{"8", "6", "9", "12"})

	for _, v := range ret {
		t.Log(
			v.Type,
			v.UserId,
			v.Price,
			v.Number,
		)
	}

	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func Test_redisDaily(t *testing.T) {
	redisDaily()
}

func TestSaveDaily(t *testing.T) {
	i := 1
	j := 1
	fmt.Println(i, j)

	redisDaily()

	for range time.Tick(5 * time.Second) {

		redisDaily()
		j++
		fmt.Println(j)

	}

}
