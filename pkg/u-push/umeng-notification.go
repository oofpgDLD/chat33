package u_push

import (
	"encoding/json"
)

type fieldSetter interface {
	SetPredefinedKeyValue(key string, value interface{}) bool
	getPostBody() ([]byte, error)
	getAppMasterSecret() string
}

var (
	ROOT_KEYS = map[string]bool{
		"appkey":          true,
		"timestamp":       true,
		"type":            true,
		"device_tokens":   true,
		"alias":           true,
		"alias_type":      true,
		"file_id":         true,
		"filter":          true,
		"production_mode": true,
		"feedback":        true,
		"description":     true,
		"thirdparty_id":   true,
		"mipush":          true,
		"mi_activity":     true,
	}

	POLICY_KEYS = map[string]bool{
		"start_time":   true,
		"expire_time":  true,
		"max_send_num": true,
	}

	PAYLOAD_KEYS = map[string]bool{
		"display_type": true,
	}

	BODY_KEYS = map[string]bool{
		"ticker":       true,
		"title":        true,
		"text":         true,
		"builder_id":   true,
		"icon":         true,
		"largeIcon":    true,
		"img":          true,
		"play_vibrate": true,
		"play_lights":  true,
		"play_sound":   true,
		"sound":        true,
		"after_open":   true,
		"url":          true,
		"activity":     true,
		"custom":       true,
	}

	APS_KEYS = map[string]bool{
		"alert":             true,
		"badge":             true,
		"sound":             true,
		"content-available": true,
	}
)

type UmengNotification struct {
	fieldSetter
	rootJson        map[string]interface{}
	appMasterSecret string
}

func NewUmengNotification(f fieldSetter) *UmengNotification {
	var t UmengNotification
	t.fieldSetter = f
	t.rootJson = make(map[string]interface{})
	return &t
}

func (t *UmengNotification) getPostBody() ([]byte, error) {
	b, err := json.Marshal(t.rootJson)
	return b, err
}

func (t *UmengNotification) getAppMasterSecret() string {
	return t.appMasterSecret
}

func (t *UmengNotification) SetAppMasterSecret(secret string) {
	t.appMasterSecret = secret
}

// 可选，正式/测试模式。默认为true
// 测试模式只对“广播”、“组播”类消息生效，其他类型的消息任务（如“文件播”）不会走测试模式
// 测试模式只会将消息发给测试设备。测试设备需要到web上添加。
// Android: 测试设备属于正式设备的一个子集。
func (t *UmengNotification) setProductionMode(prod bool) {
	t.SetPredefinedKeyValue("production_mode", prod) //prod.toString()
}

///正式模式
func (t *UmengNotification) SetReleaseMode() {
	t.setProductionMode(true)
}

///测试模式
func (t *UmengNotification) SetTestMode() {
	t.setProductionMode(false)
}

// 可选，发送消息描述，建议填写。
func (t *UmengNotification) SetDescription(description string) {
	t.SetPredefinedKeyValue("description", description)
}

// 可选，定时发送时，若不填写表示立即发送。格式: "YYYY-MM-DD hh:mm:ss"。
func (t *UmengNotification) SetStartTime(startTime string) {
	t.SetPredefinedKeyValue("start_time", startTime)
}

// 可选，消息过期时间，其值不可小于发送时间或者start_time(如果填写了的话)，
// 如果不填写此参数，默认为3天后过期。格式同start_time
func (t *UmengNotification) SetExpireTime(expireTime string) {
	t.SetPredefinedKeyValue("expire_time", expireTime)
}

// 可选，发送限速，每秒发送的最大条数。最小值1000
// 开发者发送的消息如果有请求自己服务器的资源，可以考虑此参数。
func (t *UmengNotification) SetMaxSendNum(num int) {
	t.SetPredefinedKeyValue("max_send_num", num)
}
