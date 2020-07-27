package excRate

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

// 货币价格
func TestPrice(t *testing.T) {
	i := 1
	getUsdbuy()
	getBasicPrice()
	t.Log("第一次 获取")
	for range time.Tick(10 * time.Second) {
		getUsdbuy()
		getBasicPrice()
		i++
		t.Log(usdtrmb, tickerSymbol)
		currency := "YCC"
		num := 20.0
		p := Price(currency, num)
		t.Log(p)
		t.Logf("连续获取：%v", i)
	}
	t.Log("end")
}

func Test_postUrl(t *testing.T) {
	resp, _ := postUrl("http://127.0.0.1:6062/h24data?callback=jsonp1421378514626&type=bty&symbol=ETHUSDT")

	//{"ts":"1574233200","open":"9417.45","high":"9417.45","low":"9417.45","last":"9417.45","volume":"0.0000"}
	var p UrlPostParam
	json.Unmarshal(resp, &p)

	//特定币种汇率
	last, _ := strconv.ParseFloat(p.Last, 64)
	t.Log(last)
}
