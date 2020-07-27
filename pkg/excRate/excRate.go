package excRate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	logger "github.com/golang/glog"
)

var Rate *types.Rate
var usdtrmb float64
var tickerSymbol = struct {
	sync.RWMutex
	baseCoinCNYMap map[string]float64 //  "BTC":9000.00
}{baseCoinCNYMap: make(map[string]float64)}

func Init(c *types.Config) {

	Rate = &c.Rate

	go func() {
		//获取usdt的汇率
		getUsdbuy()
		//获取货币的汇率
		getBasicPrice()
		// 每60分钟执行一次

		for range time.Tick(60 * time.Minute) {
			getUsdbuy()
			getBasicPrice()
		}
	}()
}

// 货币价格
func Price(currency string, num float64) float64 {
	rmb, OK := readTickerSymbol(currency)
	if !OK {
		logger.Error("GetMarketInfo err", "error", "基础货币不存在")
		return 0
	}
	return rmb * num
}

func ReadAllTickerSymbol() map[string]float64 {
	ret := make(map[string]float64)
	for k, v := range tickerSymbol.baseCoinCNYMap {
		ret[k] = v
	}
	return ret
}

// 获取usdt汇率
func getUsdbuy() {
	resp, err := postUrl(Rate.UsdtUrl)
	//resp, err := postUrl(UsdtUrl)
	if err != nil {
		logger.Error("UsdtUrl Post err", "error", err.Error())
	}
	//{"code":200,"ecode":"200","error":"OK","message":"OK","data":{"usdbuy":"7.02","usdsell":"6.98","updatetime":1574324043,"from":"usdtCnyByHuobi"}}
	var p UrlPostUsdbuy
	err = json.Unmarshal(resp, &p)
	if err != nil {
		logger.Error("IsSymbol.json.Unmarshal err", "error", err.Error())
	}
	data := p.Data

	usdtrmb = utility.ToFloat64(data["usdbuy"])

}

// 计算所有基础货币价格  tickerSymbol结构体 baseCoinCNYMap   "BTC":9417.45
func getBasicPrice() {
	var rate float64
	var rmb float64
	for k, v := range Slice {

		//打赏的是人民币的情况,为固定值1
		if k == "CCNY" || k == "CNY" {
			writeTickerSymbol(k, 1.00)
			continue
		}
		//每个币种汇率
		rate = isSymbol(v)
		//USTD汇率  实际的人民币价格
		rmb = calculatePrice(rate)
		writeTickerSymbol(k, rmb)
	}
}

func isSymbol(symbol string) float64 {

	//currency币种的汇率  如BTC_USTD
	tradeUrl := Rate.TradeUrl + symbol
	//tradeUrl := TradeUrl + symbol
	resp, err := postUrl(tradeUrl)
	if err != nil {
		logger.Error("TradeUrl http.Post err", "error", err.Error())
	}
	//{"ts":"1574233200","open":"9417.45","high":"9417.45","low":"9417.45","last":"9417.45","volume":"0.0000"}
	var p UrlPostParam
	err = json.Unmarshal(resp, &p)
	if err != nil {
		logger.Error("IsSymbol.json.Unmarshal err", "error", err.Error())
	}

	//特定币种汇率
	last, _ := strconv.ParseFloat(p.Last, 64)

	return last

}

// 写锁 tickerSymbol 储存数据到结构体Map中
func writeTickerSymbol(key string, value float64) {
	tickerSymbol.Lock()
	tickerSymbol.baseCoinCNYMap[key] = value
	tickerSymbol.Unlock()
}

// 读锁 tickerSymbol 读取数据到结构体Map中
func readTickerSymbol(symbol string) (float64, bool) {
	tickerSymbol.RLock()
	rmb, OK := tickerSymbol.baseCoinCNYMap[symbol]
	tickerSymbol.RUnlock()

	return rmb, OK
}

// 计算货币价格(Lastrmb)
func calculatePrice(rate float64) float64 {

	return rate * usdtrmb
}

//获取汇率接口
func postUrl(url string) ([]byte, error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(""))
	if err != nil {
		logger.Error("http.Post err", "error", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(body) < 1 {
		logger.Error("ioutil.ReadAll(resp.Body) err", "error", err.Error(), "len(body)", len(body))
		return nil, err
	}
	return body, nil
}
