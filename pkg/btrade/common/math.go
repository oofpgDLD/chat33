package common

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func Max(first int64, args ...int64) int64 {
	for _, v := range args {
		if first < v {
			first = v
		}
	}
	return first
}

func Min(first int64, args ...int64) int64 {
	for _, v := range args {
		if first > v {
			first = v
		}
	}
	return first
}

func ToInt64(o interface{}) int64 {
	if o == nil {
		return 0
	}
	switch t := o.(type) {
	case float64:
		return int64(t)
	case int:
		return int64(t)
	case int64:
		return t
	case string:
		return StringToInt64(t)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func RoundMin(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc(f*pow10_n) / pow10_n
}

func StringToInt64(str string) int64 {
	temp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return temp
}

func StringToInt32(str string) int32 {
	temp, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(temp)
}

func IntToString(data interface{}) string {
	return ToString(data)
}

func StringConnect(data []string) string {
	buffer := bytes.Buffer{}

	for _, v := range data {
		_, err := buffer.WriteString(v)
		if err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

func StringToFloat64(str string) float64 {
	temp, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return temp
}

func Float64ToString(data float64, num int) string {
	ft := make([]string, 0, 10)
	ft = append(ft, "%0.", IntToString(num), "f")
	result := fmt.Sprintf(StringConnect(ft), data)

	return result
}

func StringToBool(str string) bool {
	b, err := strconv.ParseBool(str)
	if nil != err {
		panic(err)
	}

	return b
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

// 9782.13 --> 978213000000
// cause 9782.13*1e8 --> 978212999999
func String2BaseInt64(num string) int64 {
	var decimals int
	var ret int64

	fnum, _ := strconv.ParseFloat(num, 64)
	ss := strings.Split(num, ".")
	if len(ss) == 2 {
		decimals = len(ss[1])
	}

	if decimals < 8 {
		ret = int64(fnum * 1e8)
		ret = ((ret + 1) / int64(math.Pow10(8-decimals))) * int64(math.Pow10(8-decimals))
	} else if decimals == 8 {
		ret = int64(fnum*1e9) + 1
		ret = ret / 10
	}

	return ret
}
