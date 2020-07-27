package utility_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/33cn/chat33/utility"
)

func TestUUID(t *testing.T) {
	v := utility.RandomID()
	t.Log(v)
}

func Test_CheckPhoneNumber(t *testing.T) {
	v := utility.CheckPhoneNumber("19857036394")
	fmt.Println(v)
}

func Test_PadLeft(t *testing.T) {
	v := utility.PadLeft("123", 5, '0')
	fmt.Println(v)
}

/*
func Test_GetMonthStartAndEnd(t *testing.T) {
	start, end := utility.GetMonthStartAndEnd("2019-01")
	fmt.Println(start, "    ", end)
}*/

func Test_ToString(t *testing.T) {
	var v float64 = 0
	s := strconv.FormatFloat(v, 'f', -1, 64)
	fmt.Println(s)
}

func Test_NowDayStartAndEnd(t *testing.T) {
	tm := utility.MillionSecondToDateTime(1554825600000)
	t1, t2 := utility.DayStartAndEndNowMillionSecond(tm)
	fmt.Println(t1)
	fmt.Println(t2)
}

func Test_MillionSecondToDateTime(t *testing.T) {
	tm := utility.MillionSecondToDateTime(1554912000000)  //1554911999000)
	tm2 := utility.MillionSecondToDateTime(1554048000000) //1554048000000)
	s := tm.Sub(tm2)
	subDay := s / (time.Hour * 24)

	if s%(time.Hour*24) > 0 {
		subDay += 1
	}

	/*fmt.Printf("yearday:%d\n" ,tm.YearDay())
	fmt.Printf("weekday:%d\n" ,tm.Weekday())*/
	fmt.Printf("subDay:%d\n", subDay)
	fmt.Println(tm.Format("2006-01-02 15:04"))
	fmt.Println(tm2.Format("2006-01-02 15:04"))
}

func Test_CreateSalt(t *testing.T) {
	fmt.Printf("%x", utility.CreateSalt())
}

func Test_Dirc(t *testing.T) {
	str, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(str)
}

func Test_CheckPubkey(t *testing.T) {
	ret := utility.CheckPubkey("c2694f92c3a77996d6d7658adef36a6ade881eef5ec8f70346d00882726950ae") //1554048000000)
	t.Log(ret)
}

func Test_GetWeekStartAndEnd(t *testing.T) {
	tm := time.Now()
	tm = tm.AddDate(0, 0, -7)
	start, end := utility.GetWeekStartAndEnd(tm)
	t.Log(start, "-----", end)
	t.Log(time.Unix(start/1000, start%1000).Format("2006-01-02 15:04:05"))
	t.Log(time.Unix(end/1000, end%1000).Format("2006-01-02 15:04:05"))
}

func Test_GetWeekStartAndEnd2(t *testing.T) {
	startTime := int64(1576389262000)
	endTime := int64(0)
	sec := int64(1576389262000 / 1000)
	nsec := int64(1576389262000 % 1000)
	tm := time.Unix(sec, nsec)
	nowStart, nowEnd := utility.GetWeekStartAndEnd(tm)

	t.Log(nowStart, "-----", nowEnd)
	t.Log(time.Unix(nowStart/1000, nowStart%1000).Format("2006-01-02 15:04:05"))
	t.Log(time.Unix(nowEnd/1000, nowEnd%1000).Format("2006-01-02 15:04:05"))

	if nowStart <= startTime && nowEnd >= endTime {
		startTime = nowStart
		endTime = nowEnd
	}

	t.Log(startTime, "-----", endTime)
	t.Log(time.Unix(startTime/1000, startTime%1000).Format("2006-01-02 15:04:05"))
	t.Log(time.Unix(endTime/1000, endTime%1000).Format("2006-01-02 15:04:05"))
}

func Test_PageWeek(t *testing.T) {
	page := 1
	number := 2

	baseTime := time.Now()
	for i := (page-1)*number + 1; i <= page*number; i++ {
		tm := baseTime.AddDate(0, 0, -(i * 7))
		start, end := utility.GetWeekStartAndEnd(tm)
		t.Log(start, "-----", end)
		t.Log(time.Unix(start/1000, start%1000).Format("2006-01-02 15:04:05"))
		t.Log(time.Unix(end/1000, end%1000).Format("2006-01-02 15:04:05"))
		t.Log("---------------", "number:", i, "-------------------")
	}
}
