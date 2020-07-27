package api

import (
	"net/http"
	"strings"

	"github.com/33cn/chat33/pkg/excRate"

	"github.com/33cn/chat33/types"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/utility"
	"github.com/gin-gonic/gin"
)

//用户session
func FindUserSession(c *gin.Context) {
	/*var params struct {
		Session string    `form:"session"`
		Id string    `form:"id"`
	}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.PureJSON(http.StatusOK, err.Error())
		return
	}

	c.Header("Cookie", params.Session)
	session := sessions.Default(c)
	userId := session.Get("user_id")
	appId := session.Get("appId")
	token := session.Get("token")

	c.PureJSON(http.StatusOK, fmt.Sprintf("%s%s%s", userId, appId, token))*/
}

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "success")
}

//更新App配置
func AppUpdate(c *gin.Context) {
	app.UpdateAppsConfig()
	app.UpdateVersionConfig()
	c.Set(ReqError, nil)
}

//查看App配置
func AppConfig(c *gin.Context) {
	apps := app.GetApps()
	versions := app.GetVersions()

	ret := make(map[string]interface{})
	ret["apps"] = apps
	ret["versions"] = versions
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

// 获取app列表api
func AppInfo(c *gin.Context) {
	ret := make(map[string]interface{})
	ret["app_list"] = app.GetApps()

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//检查版本更新
func VersionController(c *gin.Context) {
	type VersionCodeParas struct {
		NowVersionCode int    `json:"nowVersionCode"`
		NowVersionName string `json:"nowVersionName"`
	}
	var params VersionCodeParas
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	if params.NowVersionCode == 0 && params.NowVersionName == "" {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "code=0 name is empty"))
		return
	}

	deviceType := c.MustGet(DeviceType)
	appId := c.MustGet(AppId)
	if deviceType == "" && appId == "" {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "lack FZM-DEVICE and FZM-APP-ID"))
		return
	}

	var ret = make(map[string]interface{})
	vConfig := app.GetVersion(utility.ToString(appId), utility.ToString(deviceType))
	if vConfig == nil {
		c.Set(ReqError, result.NewError(result.AppUpdateFailed).SetExtMessage("找不到对应配置"))
		return
	}

	if deviceType == "Android" {
		//如果兼容旧版本
		if vConfig.Compatible {
			//最小版本号
			if vConfig.MinVersionCode <= params.NowVersionCode {
				ret["forceUpdate"] = false
			} else {
				ret["forceUpdate"] = true
			}
		} else if params.NowVersionCode < vConfig.VersionCode {
			ret["forceUpdate"] = true
		} else {
			ret["forceUpdate"] = false
		}
		//判断是否是特殊的几个版本
		if !ret["forceUpdate"].(bool) {
			lists := strings.Split(vConfig.ForceList, "#")
			for _, list := range lists {
				if utility.ToInt(list) == params.NowVersionCode {
					ret["forceUpdate"] = true
					break
				}
			}
		}
		ret["versionCode"] = vConfig.VersionCode
		ret["versionName"] = vConfig.VersionName
		ret["description"] = vConfig.Description
		ret["url"] = vConfig.Url
		ret["size"] = vConfig.Size
		ret["md5"] = vConfig.Md5

		c.Set(ReqError, nil)
		c.Set(ReqResult, ret)
		return
	}
	if deviceType == "iOS" {
		//根据版本号计算出code
		codeParas := strings.Split(params.NowVersionName, ".")
		if len(codeParas) != 3 {
			c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "version name illegal"))
			return
		}
		code := utility.ToInt(codeParas[0])*10000 + utility.ToInt(codeParas[1])*100 + utility.ToInt(codeParas[2])
		if vConfig.Compatible {
			if vConfig.MinVersionCode <= code {
				ret["forceUpdate"] = false
			} else {
				ret["forceUpdate"] = true
			}
		} else {
			nowCodeParas := strings.Split(vConfig.VersionName, ".")
			if len(nowCodeParas) != 3 {
				c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "version name illegal"))
				return
			}
			nowCode := utility.ToInt(nowCodeParas[0])*10000 + utility.ToInt(nowCodeParas[1])*100 + utility.ToInt(nowCodeParas[2])
			if nowCode <= code {
				ret["forceUpdate"] = false
			} else {
				ret["forceUpdate"] = true
			}
		}
		//判断是否是特殊的几个版本
		if !ret["forceUpdate"].(bool) {
			lists := strings.Split(vConfig.ForceList, "#")
			for _, list := range lists {
				if list == params.NowVersionName {
					ret["forceUpdate"] = true
					break
				}
			}
		}
		ret["versionName"] = vConfig.VersionName
		ret["description"] = vConfig.Description
		ret["url"] = vConfig.Url
		ret["size"] = vConfig.Size

		c.Set(ReqError, nil)
		c.Set(ReqResult, ret)
		return
	}
}

//oss相关接口
func OssUpdate(c *gin.Context) {
	var params struct {
		AppId           string `form:"AppId" json:"AppId" binding:"required"`
		AccessKeyID     string `form:"AccessKeyID" json:"AccessKeyID" binding:"required"`
		AccessKeySecret string `form:"AccessKeySecret" json:"AccessKeySecret" binding:"required"`
	}
	if err := c.ShouldBind(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	oss := &types.OssAddress{
		AccessKeyID:     params.AccessKeyID,
		AccessKeySecret: params.AccessKeySecret,
	}

	err := app.UpdateOss(params.AppId, oss)
	if err != nil {
		c.Set(ReqError, result.NewError(result.DbConnectFail))
		return
	}
	app := app.GetApp(params.AppId)
	if app != nil {
		app.OssConfig = oss
	}
	c.Set(ReqError, nil)
}

//oss相关接口
func AllTickerSymbol(c *gin.Context) {
	ret := excRate.ReadAllTickerSymbol()
	c.Set(ReqResult, ret)
	c.Set(ReqError, nil)
}
