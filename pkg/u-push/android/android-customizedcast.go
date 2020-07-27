package android

import u_push "github.com/33cn/chat33/pkg/u-push"

type AndroidCustomizedcast struct {
	*u_push.AndroidNotification
}

func NewAndroidCustomizedcast(appkey, appMasterSecret string) AndroidCustomizedcast {
	var t AndroidCustomizedcast
	t.AndroidNotification = u_push.NewAndroidNotification()
	t.SetAppMasterSecret(appMasterSecret)
	t.SetPredefinedKeyValue("appkey", appkey)
	t.SetPredefinedKeyValue("type", "customizedcast")
	return t
}

func (t *AndroidCustomizedcast) SetAlias(alias, aliasType string) {
	t.SetPredefinedKeyValue("alias", alias)
	t.SetPredefinedKeyValue("alias_type", aliasType)
}

func (t *AndroidCustomizedcast) SetFileId(fileId, aliasType string) {
	t.SetPredefinedKeyValue("file_id", fileId)
	t.SetPredefinedKeyValue("alias_type", aliasType)
}
