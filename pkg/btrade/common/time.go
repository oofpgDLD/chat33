package common

import (
	"time"
)

var loc = local()

func UnixNow() int64 {
	return time.Now().Unix()
}

func local() *time.Location {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	return loc
}

func ToCstTime(layout string, timeStr string) time.Time {
	t, _ := time.ParseInLocation(layout, timeStr, loc)
	return t
}

func CstTime(timeStr string) time.Time {
	// 2017-10-26T10:02:56.205Z      UTC time
	// 2017-10-26T17:15:46.711+08:00 CST time
	t, _ := time.Parse(time.RFC3339, timeStr) // UTC
	return t.In(loc)
}

func Time2YYMMDDhhmmss(t time.Time) string {
	return t.In(loc).Format("2006-01-02 15:04:05")
}

func Sec2YYMMDDhhmmss(sec int64) string {
	return time.Unix(sec, 0).In(loc).Format("2006-01-02 15:04:05")
}

// TimeNowYYMMDDhhmmss kk
func TimeNowYYMMDDhhmmss() string {
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func TimeNowYYMMDD() string {
	return time.Now().In(loc).Format("2006-01-02")
}

// TimeNowYYMMDDhhmmss kk
func TimeNowUnix() int64 {
	return time.Now().In(loc).Unix()
}

//ns
func TimeNowUnixNano() int64 {
	return time.Now().In(loc).UnixNano()
}
