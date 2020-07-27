package app

import (
	"github.com/33cn/chat33/app/dao"
	"github.com/33cn/chat33/types"
)

var appList = make(map[string]*types.App)

var versionList = make(map[string]*types.Version)

func UpdateAppsConfig() {
	appList = make(map[string]*types.App)
	apps, err := dao.FindAllAppInfo()
	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		appList[a.AppId] = a
	}
}

func GetApps() map[string]*types.App {
	return appList
}

func GetApp(appId string) *types.App {
	//如果不存在
	if a, ok := appList[appId]; !ok {
		//从数据库查询
		app, err := dao.FindAppInfoByAppId(appId)
		if err != nil {
			return nil
		}
		if app != nil {
			appList[appId] = app
		}
		return app
	} else {
		return a
	}
}

func UpdateVersionConfig() {
	versionList = make(map[string]*types.Version)
	versions, err := dao.FindAllVersionConfig()
	if err != nil {
		panic(err)
	}

	for _, v := range versions {
		versionList[v.AppId+v.Name] = v
	}
}

func GetVersions() map[string]*types.Version {
	return versionList
}

func GetVersion(appId string, devType string) *types.Version {
	return versionList[appId+devType]
}

// oss
func GetOss(appId string) *types.OssAddress {
	app := GetApp(appId)
	if app == nil {
		return nil
	}
	return app.OssConfig
}

func UpdateOss(appId string, oss *types.OssAddress) error {
	return dao.UpdateOssConfig(appId, oss)
}
