package u_push

import "encoding/json"

type DisplayType string

const (
	NOTIFICATION DisplayType = "notification"
	MESSAGE      DisplayType = "message"
)

type AfterOpenAction string

const (
	GO_APP      AfterOpenAction = "go_app"      //打开应用
	GO_URL      AfterOpenAction = "go_url"      //跳转到URL
	GO_ACTIVITY AfterOpenAction = "go_activity" //打开特定的activity
	GO_CUSTOM   AfterOpenAction = "go_custom"   //用户自定义内容。
)

type AndroidNotification struct {
	*UmengNotification
}

func NewAndroidNotification() *AndroidNotification {
	var t AndroidNotification
	t.UmengNotification = NewUmengNotification(&t)
	return &t
}

func (t *AndroidNotification) SetPredefinedKeyValue(key string, value interface{}) bool {
	if _, ok := ROOT_KEYS[key]; ok {
		// This key should be in the root level
		t.rootJson[key] = value
	} else if ok := PAYLOAD_KEYS[key]; ok {
		// This key should be in the payload level
		var payload map[string]interface{}
		if p, ok := t.rootJson["payload"]; ok {
			payload = p.(map[string]interface{})
		} else {
			payload = make(map[string]interface{})
		}
		payload[key] = value
		t.rootJson["payload"] = payload
	} else if _, ok := BODY_KEYS[key]; ok {
		// This key should be in the body level
		var payload map[string]interface{}
		var body map[string]interface{}
		// 'body' is under 'payload', so build a payload if it doesn't exist
		if p, ok := t.rootJson["payload"]; ok {
			payload = p.(map[string]interface{})
		} else {
			payload = make(map[string]interface{})
		}

		// Get body JSONObject, generate one if not existed
		if b, ok := payload["body"]; ok {
			body = b.(map[string]interface{})
		} else {
			body = make(map[string]interface{})
		}
		body[key] = value
		//
		payload["body"] = body
		t.rootJson["payload"] = payload
	} else if _, ok := POLICY_KEYS[key]; ok {
		// This key should be in the body level
		var policy map[string]interface{}
		// 'body' is under 'payload', so build a payload if it doesn't exist
		if p, ok := t.rootJson["policy"]; ok {
			policy = p.(map[string]interface{})
		} else {
			policy = make(map[string]interface{})
		}
		policy[key] = value
		t.rootJson["policy"] = policy
	} else {
		if key == "payload" || key == "body" || key == "policy" || key == "extra" {
			panic("You don't need to set value for " + key + " , just set values for the sub keys in it.")
		} else {
			panic("Unknown key: " + key)
		}
	}
	return true
}

// Set extra key/value for Android notification
func (t *AndroidNotification) SetExtraField(key string, value interface{}) bool {
	// This key should be in the body level
	var payload map[string]interface{}
	var extra map[string]interface{}
	// 'body' is under 'payload', so build a payload if it doesn't exist
	if p, ok := t.rootJson["payload"]; ok {
		payload = p.(map[string]interface{})
	} else {
		payload = make(map[string]interface{})
	}

	if e, ok := payload["extra"]; ok {
		extra = e.(map[string]interface{})
	} else {
		extra = make(map[string]interface{})
	}
	extra[key] = value
	//
	payload["extra"] = extra
	t.rootJson["payload"] = payload
	return true
}

//必填，消息类型: notification(通知)、message(消息)
func (t *AndroidNotification) SetDisplayType(d DisplayType) {
	t.SetPredefinedKeyValue("display_type", d)
}

//必填，通知栏提示文字
func (t *AndroidNotification) SetTicker(ticker string) {
	t.SetPredefinedKeyValue("ticker", ticker)
}

//必填，通知标题
func (t *AndroidNotification) SetTitle(title string) {
	t.SetPredefinedKeyValue("title", title)
}

//必填，通知文字描述
func (t *AndroidNotification) SetText(text string) {
	t.SetPredefinedKeyValue("text", text)
}

//可选，默认为0,用于标识该通知采用的样式。使用该参数时, 必须在SDK里面实现自定义通知栏样式。
func (t *AndroidNotification) SetBuilderId(builder_id int) {
	t.SetPredefinedKeyValue("builder_id", builder_id)
}

// 可选，状态栏图标ID，R.drawable.[smallIcon]，
// 如果没有，默认使用应用图标。
// 图片要求为24*24dp的图标，或24*24px放在drawable-mdpi下。
// 注意四周各留1个dp的空白像素
func (t *AndroidNotification) SetIcon(icon string) {
	t.SetPredefinedKeyValue("icon", icon)
}

// 可选，通知栏拉开后左侧图标ID，R.drawable.[largeIcon]，
// 图片要求为64*64dp的图标，
// 可设计一张64*64px放在drawable-mdpi下，
// 注意图片四周留空，不至于显示太拥挤
func (t *AndroidNotification) SetLargeIcon(largeIcon string) {
	t.SetPredefinedKeyValue("largeIcon", largeIcon)
}

//可选，通知栏大图标的URL链接。该字段的优先级大于largeIcon。该字段要求以http或者https开头。
func (t *AndroidNotification) SetImg(img string) {
	t.SetPredefinedKeyValue("img", img)
}

//可选，收到通知是否震动,默认为"true"
func (t *AndroidNotification) SetPlayVibrate(play_vibrate bool) {
	t.SetPredefinedKeyValue("play_vibrate", play_vibrate) //play_vibrate.toString()
}

//可选，收到通知是否闪灯,默认为"true"
func (t *AndroidNotification) SetPlayLights(play_lights bool) {
	t.SetPredefinedKeyValue("play_lights", play_lights) //play_lights.toString()
}

//可选，收到通知是否发出声音,默认为"true"
func (t *AndroidNotification) SetPlaySound(play_sound bool) {
	t.SetPredefinedKeyValue("play_sound", play_sound) //play_sound.toString()
}

//可选，通知声音，R.raw.[sound]. 如果该字段为空，采用SDK默认的声音
// umeng_push_notification_default_sound声音文件。如果
// SDK默认声音文件不存在，则使用系统默认Notification提示音。
func (t *AndroidNotification) SetSound(sound string) {
	t.SetPredefinedKeyValue("sound", sound)
}

///收到通知后播放指定的声音文件
func (t *AndroidNotification) SetPlaySoundAndSound(sound string) {
	t.SetPlaySound(true)
	t.SetSound(sound)
}

///点击"通知"的后续行为，默认为打开app。
func (t *AndroidNotification) GoAppAfterOpen() {
	t.setAfterOpenAction(GO_APP)
}

func (t *AndroidNotification) GoUrlAfterOpen(url string) {
	t.setAfterOpenAction(GO_URL)
	t.setUrl(url)
}

func (t *AndroidNotification) GoActivityAfterOpen(activity string) {
	t.setAfterOpenAction(GO_ACTIVITY)
	t.setActivity(activity)
}

func (t *AndroidNotification) GoCustomAfterOpen(custom string) {
	t.setAfterOpenAction(GO_CUSTOM)
	t.setCustomField(custom)
}

func (t *AndroidNotification) GoCustomAfterOpenJson(custom json.RawMessage) {
	t.setAfterOpenAction(GO_CUSTOM)
	t.setCustomFieldJson(custom)
}

// 可选，默认为false。当为true时，表示MIUI、EMUI、Flyme系统设备离线转为系统下发
// 可选，mipush值为true时生效，表示走系统通道时打开指定页面acitivity的完整包路径。
func (t *AndroidNotification) SetMipush(enable bool, miActive string) {
	t.SetPredefinedKeyValue("mipush", enable)
	t.SetPredefinedKeyValue("mi_activity", miActive)
}

///点击"通知"的后续行为，默认为打开app。原始接口
func (t *AndroidNotification) setAfterOpenAction(action AfterOpenAction) {
	t.SetPredefinedKeyValue("after_open", action)
}

func (t *AndroidNotification) setUrl(url string) {
	t.SetPredefinedKeyValue("url", url)
}

func (t *AndroidNotification) setActivity(activity string) {
	t.SetPredefinedKeyValue("activity", activity)
}

///can be a string of json
func (t *AndroidNotification) setCustomField(custom string) {
	t.SetPredefinedKeyValue("custom", custom)
}

func (t *AndroidNotification) setCustomFieldJson(custom json.RawMessage) {
	t.SetPredefinedKeyValue("custom", string(custom))
}
