package u_push

type IOSNotification struct {
	*UmengNotification
}

func NewIOSNotification() *IOSNotification {
	var t IOSNotification
	t.UmengNotification = NewUmengNotification(&t)
	return &t
}

func (t *IOSNotification) SetPredefinedKeyValue(key string, value interface{}) bool {
	if _, ok := ROOT_KEYS[key]; ok {
		// This key should be in the root level
		t.rootJson[key] = value
	} else if ok := APS_KEYS[key]; ok {
		// This key should be in the aps level
		var aps map[string]interface{}
		var payload map[string]interface{}
		if p, ok := t.rootJson["payload"]; ok {
			payload = p.(map[string]interface{})
		} else {
			payload = make(map[string]interface{})
		}

		if a, ok := payload["aps"]; ok {
			aps = a.(map[string]interface{})
		} else {
			aps = make(map[string]interface{})
		}

		aps[key] = value
		payload["aps"] = aps
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
		if key == "payload" || key == "aps" || key == "policy" {
			panic("You don't need to set value for " + key + " , just set values for the sub keys in it.")
		} else {
			panic("Unknown key: " + key)
		}
	}
	return true
}

// Set customized key/value for IOS notification
func (t *IOSNotification) SetCustomizedField(key string, value interface{}) bool {
	// This key should be in the body level
	var payload map[string]interface{}
	// 'body' is under 'payload', so build a payload if it doesn't exist
	if p, ok := t.rootJson["payload"]; ok {
		payload = p.(map[string]interface{})
	} else {
		payload = make(map[string]interface{})
	}

	payload[key] = value
	t.rootJson["payload"] = payload
	return true
}

type IOSAlert struct {
	// 可为JSON类型
	Title    string `json:"title",omitempty`
	Subtitle string `json:"subtitle",omitempty`
	Body     string `json:"body",omitempty`
}

func (t *IOSNotification) SetAlertJson(alert IOSAlert) {
	t.SetPredefinedKeyValue("alert", alert)
}

func (t *IOSNotification) SetAlert(token string) {
	t.SetPredefinedKeyValue("alert", token)
}

func (t *IOSNotification) SetBadge(badge int) {
	t.SetPredefinedKeyValue("badge", badge)
}

func (t *IOSNotification) SetSound(sound string) {
	t.SetPredefinedKeyValue("sound", sound)
}

func (t *IOSNotification) SetContentAvailable(contentAvailable int) {
	t.SetPredefinedKeyValue("content-available", contentAvailable)
}
