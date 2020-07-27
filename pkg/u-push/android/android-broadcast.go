package android

import (
	u_push "github.com/33cn/chat33/pkg/u-push"
)

type AndroidBroadcast struct {
	*u_push.AndroidNotification
}

func NewAndroidBroadcast(appkey, appMasterSecret string) AndroidBroadcast {
	var t AndroidBroadcast
	t.AndroidNotification = u_push.NewAndroidNotification()
	t.SetAppMasterSecret(appMasterSecret)
	t.SetPredefinedKeyValue("appkey", appkey)
	t.SetPredefinedKeyValue("type", "broadcast")
	return t
}
