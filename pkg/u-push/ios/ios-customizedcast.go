package ios

import u_push "github.com/33cn/chat33/pkg/u-push"

type IOSCustomizedcast struct {
	*u_push.IOSNotification
}

func NewIOSCustomizedcast(appkey, appMasterSecret string) IOSCustomizedcast {
	var t IOSCustomizedcast
	t.IOSNotification = u_push.NewIOSNotification()
	t.SetAppMasterSecret(appMasterSecret)
	t.SetPredefinedKeyValue("appkey", appkey)
	t.SetPredefinedKeyValue("type", "customizedcast")
	return t
}

func (t *IOSCustomizedcast) SetAlias(alias, aliasType string) {
	t.SetPredefinedKeyValue("alias", alias)
	t.SetPredefinedKeyValue("alias_type", aliasType)
}

func (t *IOSCustomizedcast) SetFileId(fileId, aliasType string) {
	t.SetPredefinedKeyValue("file_id", fileId)
	t.SetPredefinedKeyValue("alias_type", aliasType)
}
