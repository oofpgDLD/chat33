package middle_ware

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/33cn/chat33/utility"
)

func BcoinBackend(appKey, appSecret, time string, params map[string]string) string {
	//1.将参数按照 键名 的 字典序升序 排序
	ordered := make([]string, 0)
	for k := range params {
		ordered = append(ordered, k)
	}
	sort.Strings(ordered)
	//2.将参数按照 k1=v1&k2=v2.... 的形式组合成字符串，假定名为 params
	u := url.Values{}
	for _, k := range ordered {
		v := params[k]
		u.Set(k, v)
	}
	str := u.Encode()
	//3.将 appKey  params  appSecret  time  按顺序拼接在一起，并作md5加密，然后转大写(其中 appKey appSecret都是事先交换的，time 是头部传递过来的时间)
	secretStr := appKey + str + appSecret + time
	w := md5.New()
	_, err := io.WriteString(w, secretStr)
	if err != nil {

	}
	return strings.ToUpper(hex.EncodeToString(w.Sum(nil)))
}

func GetParamsMap(i interface{}) map[string]string {
	ret := make(map[string]string)

	st := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	ss := st.Elem()
	vv := v.Elem()
	for i := 0; i < ss.NumField(); i++ {
		field := ss.Field(i)
		val := vv.Field(i).Interface()

		ret[utility.ToString(field.Tag.Get("json"))] = utility.ToString(val)
	}

	return ret
}
