package excRate

var Slice = map[string]string{
	"BTC":  "BTCUSDT",
	"ETH":  "ETHUSDT",
	"ETC":  "ETCUSDT",
	"ZEC":  "ZECUSDT",
	"LTC":  "LTCUSDT",
	"BCC":  "BCCUSDT",
	"BTY":  "BTYUSDT",
	"SC":   "SCUSDT",
	"BTS":  "BTSUSDT",
	"BNT":  "BNTUSDT",
	"DCR":  "DCRUSDT",
	"SCTC": "SCTCUSDT",
	"YCC":  "YCCUSDT",
	"JB":   "JBUSDT",
	"OPTC": "OPTCUSDT",
	"ITVB": "ITVBUSDT",
	"HT":   "HTUSDT",
	"FH":   "FHUSDT",
	"CWV":  "CWVUSDT",
	"FANS": "FANSUSDT",
	"CNS":  "CNSUSDT",
	"HBT":  "HBTUSDT",
	"BNB":  "BNBUSDT",
	"SCT":  "SCTUSDT",
	"WTC":  "WTCUSDT",
	"EOS":  "EOSUSDT",
	"SFT":  "SFTUSDT",
	"CCNY": "CCNYUSDT",
	"CNY":  "CNYUSDT",
}

// Json解析JSON "https://api.biqianbao.top/api/data/rate"
type UrlPostUsdbuy struct {
	Code    int                    `json:"code"`
	Ecode   string                 `json:"ecode"`
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// Json解析data "http://122.224.124.250:10084/h24data?callback=jsonp1421378514626&type=bty&symbol=***"
type UrlPostParam struct {
	Ts     string `json:"ts"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Last   string `json:"last"`
	Volume string `json:"volume"`
}
