package android

import u_push "github.com/33cn/chat33/pkg/u-push"

type AndroidUnicast struct {
	*u_push.AndroidNotification
}

func NewAndroidUnicast(appkey, appMasterSecret string) AndroidUnicast {
	var t AndroidUnicast
	t.AndroidNotification = u_push.NewAndroidNotification()
	t.SetAppMasterSecret(appMasterSecret)
	t.SetPredefinedKeyValue("appkey", appkey)
	t.SetPredefinedKeyValue("type", "unicast")
	return t
}

func (t *AndroidUnicast) SetDeviceToken(token string) {
	t.SetPredefinedKeyValue("device_tokens", token)
}
