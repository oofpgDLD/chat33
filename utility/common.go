package utility

import (
	"crypto/md5"
	realRand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"reflect"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/33cn/chat33/pkg/btrade/crypto/ed25519"
	l "github.com/inconshreveable/log15"
	uuid "github.com/satori/go.uuid"
)

var utility_log = l.New("module", "chat/utility/common")

/*var loc = local()

func local() *time.Location {
	loc, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		utility_log.Warn("LoadLocation err", "err", err)
		panic(err)
	}
	return loc
}*/

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const RoomRandomId = "0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// generate random string; fast
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandStringBytesMaskImprSrc(n int, lib string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(lib) {
			b[i] = lib[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandomRoomId() string {
	return RandStringBytesMaskImprSrc(9, RoomRandomId)
}

func RandomUsername() string {
	return "chat" + RandStringBytesMaskImprSrc(10, letterBytes)
}

func NowMillionSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func NowSecond() int64 {
	return time.Now().Unix()
}

//时间戳 毫秒
func DayStartAndEndNowMillionSecond(t time.Time) (int64, int64) {
	now := t
	year := now.Year()
	month := now.Month()
	day := now.Day()

	start := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, 1).Add(-time.Second)
	return start.UnixNano() / 1e6, end.UnixNano() / 1e6
}

// 2006-01-02 15:04
func MillionSecondToDateTime(ts int64) time.Time {
	sec := ts / 1000
	nsec := ts % 1000
	return time.Unix(sec, nsec)
}

func MillionSecondAddDate(tm int64, years, months, days int) int64 {
	sec := tm / 1000
	nsec := tm % 1000
	t := time.Unix(sec, nsec)
	t = t.AddDate(years, months, days)
	return t.UnixNano() / 1e6
}

func MillionSecondAddDuration(tm int64, duration time.Duration) int64 {
	sec := tm / 1000
	nsec := tm % 1000
	t := time.Unix(sec, nsec)
	t = t.Add(duration)
	return t.UnixNano() / 1e6
}

func MillionSecondToTimeString(tm int64) string {
	sec := tm / 1000
	nsec := tm % 1000
	return time.Unix(sec, nsec).Format("2006-01-02 15:04")
}

func MillionSecondToTimeString2(tm int64) string {
	sec := tm / 1000
	nsec := tm % 1000
	return time.Unix(sec, nsec).Format("2006-01-02 15:04:05")
}

func getTimeByMillionSecond(tm int64) time.Time {
	sec := tm / 1000
	nsec := tm % 1000
	return time.Unix(sec, nsec)
}

/*func TimeStrToMillionSecond(tm string) int64 {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", tm, loc)
	if err != nil {
		panic(err)
	}
	return t.UnixNano() / 1e6
}*/

// true: first > second + day
func CompareTimeInterval(first, second int64, day int) bool {
	firstTime := getTimeByMillionSecond(first)
	secondTime := getTimeByMillionSecond(second)
	return firstTime.After(secondTime.AddDate(0, 0, day))
}

/**
	随机生成ID
**/
func RandomID() string {
	_uuid := uuid.NewV4()
	//rlt := fmt.Sprintf("%v", _uuid)
	rlt := _uuid.String()
	return strings.Replace(rlt, "-", "", -1)
}

func RandInt(min, max int) int {
	if min >= max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

/*
func TimeStampIntToTimeStr(ts int64) string {
	tm := time.Unix(ts/1000, ts%1000*1000000)
	return tm.Format(TimeLayoutMillionSecond)
}

func TimeStampToTimeStr(ts string) string {
	val, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ts
	}
	return TimeStampIntToTimeStr(val)
}

const (
	TimeLayoutMillionSecond = "2006-01-02 15:04:05.000"
	TimeLayoutSecond        = "2006-01-02 15:04:05"
)

func TimeStrToTimeStamp(tm string) string {
	if tm == "" {
		return "0"
	}
	ts, err := time.ParseInLocation(TimeLayoutMillionSecond, tm, loc)
	if err != nil {
		ts, err = time.ParseInLocation(TimeLayoutSecond, tm, loc)
		fmt.Println("TimeStrToTimeStamp: Warn time layout format")
	}
	return strconv.FormatInt(ts.UnixNano()/1000000, 10)
}*/

func RFC3339ToTimeStampMillionSecond(rfc string) int64 {
	const layout = "2006-01-02T15:04:05.000Z"
	if rfc == "" {
		return 0
	}
	ts, err := time.Parse(layout, rfc)
	if err != nil {
		l.Error("RFC3339ToTimeStampMillionSecond", "err_msg", err, "rfc", rfc)
		return 0
	}
	return ts.UnixNano() / 1000000
}

func ParseString(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

func ToBool(val interface{}) bool {
	ret := ToInt(val)
	return ret > 0
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ToInt(val interface{}) int {
	return int(ToInt32(val))
}

func ToInt32(o interface{}) int32 {
	if o == nil {
		return 0
	}
	switch t := o.(type) {
	case int:
		return int32(t)
	case int32:
		return t
	case int64:
		return int32(t)
	case float64:
		return int32(t)
	case string:
		if o == "" {
			return 0
		}
		temp, err := strconv.ParseInt(o.(string), 10, 32)
		if err != nil {
			return 0
		}
		return int32(temp)
	default:
		return 0
	}
}

func ToInt64(val interface{}) int64 {
	if val == nil {
		return 0
	}
	switch val.(type) {
	case int:
		return int64(val.(int))
	case string:
		if val.(string) == "" {
			return 0
		}
		ret, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			utility_log.Error("func ToInt64 error")
			debug.PrintStack()
			return 0
		}
		return ret
	case float64:
		return int64(val.(float64))
	case int64:
		return val.(int64)
	case json.Number:
		v := val.(json.Number)
		ret, err := v.Int64()
		if err != nil {
			return 0
		}
		return ret
	default:
		utility_log.Error("unknow type", "type", fmt.Sprintf("%T", val))
		return 0
	}
}

func ToString(val interface{}) string {
	if val == nil {
		return ""
	}
	switch val.(type) {
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case int64:
		return strconv.FormatInt(val.(int64), 10)
	}
	return fmt.Sprintf("%v", val)
}

func ToFloat32(v interface{}) float32 {
	if v == nil {
		return 0
	}

	switch t := v.(type) {
	case int:
		return float32(t)
	case int32:
		return float32(t)
	case int64:
		return float32(t)
	case float32:
		return t
	case float64:
		return float32(t)
	case string:
		ret, _ := strconv.ParseFloat(t, 32)
		return float32(ret)
	default:
		panic(reflect.TypeOf(t).String())
	}
}

func ToFloat64(val interface{}) float64 {
	if val == nil {
		return 0
	}
	switch val.(type) {
	case string:
		ret, _ := strconv.ParseFloat(val.(string), 64)
		return ret
	default:
		if v, ok := val.(float64); ok {
			return v
		}
		return 0
	}
}

func StructToString(val interface{}) string {
	if val == nil {
		return ""
	}

	switch val.(type) {
	case interface{}:
		bytes, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		return *(*string)(unsafe.Pointer(&bytes))
	default:
		return ""
	}
}

func StringToJobj(val interface{}) map[string]interface{} {
	var rlt = make(map[string]interface{})
	switch val.(type) {
	case string:
		err := json.Unmarshal([]byte(val.(string)), &rlt)
		if err != nil {
			return nil
		}
		return rlt
	default:
		bytes, err := json.Marshal(val)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, &rlt)
		if err != nil {
			panic(err)
		}
		return rlt
	}
}

func VisitorNameSplit(src string) string {
	if len(src) >= 8 {
		return src[0:8]
	} else {
		return src[0:]
	}
}

func CheckPhoneNumber(phone string) bool {
	reg := `^1\d{10}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func CheckYearly(date string) bool {
	reg := `^\d{4}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(date)
}

func PadLeft(src string, num int, s byte) string {
	if len(src) < num {
		len := num - len(src)
		data := make([]byte, len)
		for i := 0; i < len; i++ {
			data[i] = s
		}
		return fmt.Sprintf(string(data[:]) + src)
	}
	return src
}

/*func GetMonthStartAndEnd(date string) (int64, int64) {
	const layout = "2006-01"
	if date == "" {
		return 0, 0
	}
	ts, err := time.ParseInLocation("2006-01", date, loc)
	if err != nil {
		l.Error("RFC3339ToTimeStampMillionSecond", "err_msg", err, "date", date)
		return 0, 0
	}
	endTs := ts.AddDate(0, 1, 0).Add(-time.Second)
	return ts.Unix(), endTs.Unix()
}*/

func GetWeekStartAndEnd(tm time.Time) (int64, int64) {
	tmWeekday := tm.Weekday()
	if tmWeekday == time.Sunday {
		tmWeekday += 7
	}
	sub := tmWeekday - time.Monday
	add := (time.Sunday + 7) - tmWeekday
	start := tm.AddDate(0, 0, -int(sub))
	end := tm.AddDate(0, 0, int(add))

	startTs, _ := DayStartAndEndNowMillionSecond(start)
	_, endTs := DayStartAndEndNowMillionSecond(end)
	return startTs, endTs
}

//---------------版本比较-------------//
const (
	EQ = 0
	NE = 1
	GT = 2
	GE = 3
	LT = 4
	LE = 5
)

type Version string

func (v Version) Compare(t Version, tp int) bool {
	switch tp {
	case EQ:
		return v.toInt() == t.toInt()
	case NE:
		return v.toInt() != t.toInt()
	case GT:
		return v.toInt() > t.toInt()
	case GE:
		return v.toInt() >= t.toInt()
	case LT:
		return v.toInt() < t.toInt()
	case LE:
		return v.toInt() <= t.toInt()
	}
	return v.toInt() > t.toInt()
}

func (v Version) toInt() int {
	ls := strings.Split(string(v), ".")

	if len(ls) != 3 {
		return 0
	}

	return ToInt(ls[0])*10000 + ToInt(ls[1])*100 + ToInt(ls[2])
}

//------------------排序--------------//
type lessFunc func(p1, p2 interface{}) bool

// multiSorter implements the Sort interface, sorting the changes within.
type AutoSorter struct {
	changes []interface{}
	less    lessFunc
	desc    bool
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *AutoSorter) Sort(changes []interface{}) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less lessFunc, desc bool) *AutoSorter {
	return &AutoSorter{
		less: less,
		desc: desc,
	}
}

// Len is part of sort.Interface.
func (ms *AutoSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *AutoSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *AutoSorter) Less(i, j int) bool {
	p, q := ms.changes[i], ms.changes[j]
	if ms.desc {
		return ms.less(q, p)
	} else {
		return ms.less(p, q)
	}
}

//----------加盐加密----------------//
func EncodeSha256WithSalt(src string, salt string) string {
	h := sha256.New()
	h.Write([]byte(src + salt))
	//return h.Sum(nil)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func EncodeSha256(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CreateSalt() string {
	//key := make([]byte, 10)  realRand.Read(key)
	r, _ := realRand.Int(realRand.Reader, big.NewInt(10000))
	rnum := r.Int64()
	h := md5.New()
	h.Write([]byte(ToString(rnum)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//---------------检查公钥 ed25519-------------//
func CheckPubkey(key string) bool {
	data, err := hex.DecodeString(key)
	if err != nil {
		return false
	}
	driver := ed25519.Ed25519Driver{}
	_, err = driver.PubKeyFromBytes(data)
	return err == nil
}

func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

//随机生成Token (新用户注册)
func GetToken() string {
	// 获取当前时间的时间戳
	t := time.Now().Unix()

	// 生成一个MD5的哈希
	h := md5.New()

	// 将时间戳转换为byte，并写入哈希
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(t))
	h.Write([]byte(b))

	// 将字节流转化为16进制的字符串
	return hex.EncodeToString(h.Sum(nil))
}

//判断token是否过期，按30天过期
func CheckToken(startTime int64) (int64, bool) {
	// 精确到毫秒 拿到当前时间
	time := NowMillionSecond()
	// 求相差天数
	date := (time - startTime) / 86400000

	if date >= 30 {
		return date, true
	}
	return date, false
}
