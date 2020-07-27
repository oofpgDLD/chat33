package dao

import (
	"encoding/json"

	"github.com/revel/log15"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func convertAppInfo(m map[string]string) *types.App {
	//查询模块启用状态
	module, _ := findAppModuleByAppId(m["app_id"])
	var oss *types.OssAddress
	json.Unmarshal([]byte(m["oss_config"]), oss)

	//解析推送
	var pushAppKey types.UPushType
	var pushAppMs types.UPushType
	var pushMiActive types.UPushType
	err := json.Unmarshal([]byte(m["push_app_key"]), &pushAppKey)
	if err != nil {
		log15.Error("convertAppInfo u-push_app_key can not unmarshal", "err", err)
	}
	err = json.Unmarshal([]byte(m["push_app_master_secret"]), &pushAppMs)
	if err != nil {
		log15.Error("convertAppInfo push_app_master_secret-push can not unmarshal", "err", err)
	}
	err = json.Unmarshal([]byte(m["push_mi_active"]), &pushMiActive)
	if err != nil {
		log15.Error("convertAppInfo u-push_mi_active can not unmarshal", "err", err)
	}

	return &types.App{
		AppId:               m["app_id"],
		AccountServer:       m["user_info_url"],
		RPPid:               m["redpacket_pid"],
		RPServer:            m["redpacket_server"],
		RPUrl:               m["redpacket_url"],
		IsInner:             utility.ToInt(m["is_inner"]),
		IsOtc:               utility.ToInt(m["is_otc"]),
		OtcServer:           m["otc_server"],
		MainCoin:            m["main_coin"],
		AppKey:              m["backend_app_key"],
		AppSecret:           m["backend_app_secret"],
		PushAppKey:          pushAppKey,
		PushAppMasterSecret: pushAppMs,
		PushMiActive:        pushMiActive,
		Modules:             module,
		OssConfig:           oss,
	}
}

//根据appid查询模块启用信息
func findAppModuleByAppId(appId string) ([]*types.Module, error) {
	sql := "select * from app_module where app_id = ?"
	maps, err := conn.Query(sql, appId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}

	list := make([]*types.Module, 0)
	for _, info := range maps {
		item := &types.Module{
			Type:   utility.ToInt(info["type"]),
			Name:   info["name"],
			Enable: utility.ToBool(info["enable"]),
		}
		list = append(list, item)
	}
	return list, nil
}

func FindAllAppInfo() ([]*types.App, error) {
	sql := "select * from app"
	maps, err := conn.Query(sql)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	list := make([]*types.App, 0)
	for _, m := range maps {
		//查询模块启用状态
		item := convertAppInfo(m)
		list = append(list, item)
	}
	return list, nil
}

func FindAppInfoByAppId(appId string) (*types.App, error) {
	sql := "select * from app where app_id = ?"
	maps, err := conn.Query(sql, appId)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	info := maps[0]
	return convertAppInfo(info), nil
}

func FindAllVersionConfig() ([]*types.Version, error) {
	sql := "select * from app_update"
	maps, err := conn.Query(sql)
	if err != nil {
		return nil, err
	}
	if len(maps) < 1 {
		return nil, nil
	}
	list := make([]*types.Version, 0)
	for _, m := range maps {
		item := &types.Version{
			AppId:          m["app_id"],
			Name:           m["name"],
			Compatible:     utility.ToBool(m["compatible"]),
			MinVersionCode: utility.ToInt(m["min_version_code"]),
			VersionCode:    utility.ToInt(m["version_code"]),
			VersionName:    m["version_name"],
			Description:    m["description"],
			Url:            m["url"],
			Size:           utility.ToInt64(m["size"]),
			Md5:            m["md5"],
			ForceList:      m["force"],
		}
		list = append(list, item)
	}
	return list, nil
}

func UpdateOssConfig(appId string, oss *types.OssAddress) error {
	b, err := json.Marshal(oss)
	if err != nil {
		log15.Error("UpdateOssConfig Marshal", "err", err.Error())
		return err
	}
	sql := "update app set oss_config = ? where app_id = ?"
	_, _, err = conn.Exec(sql, string(b), appId)
	if err != nil {
		log15.Error("UpdateOssConfig", "err", err.Error())
	}
	return err
}
