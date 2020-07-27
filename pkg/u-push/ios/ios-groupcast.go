package ios

import (
	"encoding/json"

	u_push "github.com/33cn/chat33/pkg/u-push"
)

type IOSGroupcast struct {
	*u_push.IOSNotification
}

func NewIOSGroupcast(appkey, appMasterSecret string) IOSGroupcast {
	var t IOSGroupcast
	t.IOSNotification = u_push.NewIOSNotification()
	t.SetAppMasterSecret(appMasterSecret)
	t.SetPredefinedKeyValue("appkey", appkey)
	t.SetPredefinedKeyValue("type", "groupcast")
	return t
}

func (t *IOSGroupcast) SetFilter(filter json.RawMessage) {
	t.SetPredefinedKeyValue("filter", filter)
}
