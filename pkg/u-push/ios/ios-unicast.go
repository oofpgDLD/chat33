package ios

import u_push "github.com/33cn/chat33/pkg/u-push"

type IOSUnicast struct {
	*u_push.IOSNotification
}

func NewIOSUnicast(appkey, appMasterSecret string) IOSUnicast {
	var t IOSUnicast
	t.IOSNotification = u_push.NewIOSNotification()
	t.SetAppMasterSecret(appMasterSecret)
	t.SetPredefinedKeyValue("appkey", appkey)
	t.SetPredefinedKeyValue("type", "unicast")
	return t
}

func (t *IOSUnicast) SetDeviceToken(token string) {
	t.SetPredefinedKeyValue("device_tokens", token)
}
