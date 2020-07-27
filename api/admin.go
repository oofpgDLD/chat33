package api

import (
	"fmt"
	"path"

	"github.com/33cn/chat33/utility"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/model"
	"github.com/33cn/chat33/pkg/account"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
)

var logAdmin = log15.New("model", "api/admin")

//应用启动
func AppOpen(c *gin.Context) {
	err := model.AppOpen(c.MustGet(UserId).(string), c.MustGet(AppId).(string), c.MustGet(DeviceType).(string), c.MustGet(Version).(string))
	c.Set(ReqError, err)
}

//管理员登录
func AdminLogin(c *gin.Context) {
	type requestParams struct {
		Account  string `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	appId := c.GetHeader("FZM-APP-ID")
	if appId == "" {
		c.Set(ReqError, result.NewError(result.LackParam).SetChildErr(result.ServiceChat, nil, "lack http header: FZM-APP-ID"))
		return
	}

	adminId, err := model.AdminLogin(appId, params.Account, params.Password)
	if err != nil {
		c.Set(ReqError, err)
		return
	}

	session, err := store.Get(c.Request, sessionAdmin)
	if session == nil {
		logUser.Error("can not get session store", "err", err.Error())
		c.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}
	session.Values["admin_id"] = adminId
	session.Values["appId"] = appId
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		logUser.Error("can not save session", "err", err.Error())
		c.Set(ReqError, result.NewError(result.ServerInterError))
		return
	}
	c.Set(ReqError, err)
}

//最近统计消息
func AdminAccount(c *gin.Context) {
	account, err := model.AdminAccount(c.MustGet(AdminId).(string))
	if err != nil {
		c.Set(ReqError, err)
		return
	}
	ret := make(map[string]string)
	ret["account"] = account
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//------------------------------统计数据------------------------------//
//最近统计消息
func LatestData(c *gin.Context) {
	ret, err := model.LatestData(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//近期统计数据 根据平台分
func LatestDataPlatfrom(c *gin.Context) {
	ret, err := model.LatestDataPlatform(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//明细数据
func SumDetails(c *gin.Context) {
	type requestParams struct {
		StartTime int64 `json:"startTime" binding:"required"`
		EndTime   int64 `json:"endTime" binding:"required"`
		Count     int   `json:"count" binding:"required"`
		Page      int   `json:"page" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.SumDetails(c.MustGet(AppId).(string), params.StartTime, params.EndTime, params.Count, params.Page)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//导出明细数据
func ExportSumDetails(c *gin.Context) {
	type requestParams struct {
		StartTime int64 `json:"startTime" binding:"required"`
		EndTime   int64 `json:"endTime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	url, err := model.ExportSumDetails(c.MustGet(AppId).(string), params.StartTime, params.EndTime)
	ret := make(map[string]string)
	ret["filename"] = url
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//折线图
func SumGraph(c *gin.Context) {
	type requestParams struct {
		Types     []int `json:"types" binding:"required"`
		StartTime int64 `json:"startTime" binding:"required"`
		EndTime   int64 `json:"endTime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.SumGraph(c.MustGet(AppId).(string), params.StartTime, params.EndTime, params.Types)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//-----------------------------版本统计--------------------//
//获取所有版本
func GetAppVersions(c *gin.Context) {
	ret, err := model.GetAppVersions(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//版本统计折线图
func VersionGraph(c *gin.Context) {
	type requestParams struct {
		Types     []int    `json:"types" binding:"required"`
		Versions  []string `json:"versions" binding:"required"`
		StartTime int64    `json:"startTime" binding:"required"`
		EndTime   int64    `json:"endTime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.VersionGraph(c.MustGet(AppId).(string), params.StartTime, params.EndTime, params.Types, params.Versions)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//全部版本明细
func AppVersionDetails(c *gin.Context) {
	type requestParams struct {
		Count      int     `json:"count" binding:"required"`
		Page       int     `json:"page" binding:"required"`
		StartTime  int64   `json:"startTime" binding:"required"`
		EndTime    int64   `json:"endTime" binding:"required"`
		SortTarget *string `json:"sortTarget" binding:"required"`
		SortRule   *int    `json:"sortRule" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AppVersionDetails(c.MustGet(AppId).(string), params.StartTime, params.EndTime, params.Count, params.Page, *params.SortTarget, *params.SortRule)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//导出全部版本明细
func ExportAppVersionDetails(c *gin.Context) {
	type requestParams struct {
		StartTime  int64   `json:"startTime" binding:"required"`
		EndTime    int64   `json:"endTime" binding:"required"`
		SortTarget *string `json:"sortTarget" binding:"required"`
		SortRule   *int    `json:"sortRule" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	url, err := model.ExportAppVersionDetails(c.MustGet(AppId).(string), params.StartTime, params.EndTime, *params.SortTarget, *params.SortRule)
	ret := make(map[string]string)
	ret["filename"] = url
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//单个版本明细
func AppSpecificVersionDetails(c *gin.Context) {
	type requestParams struct {
		Version   string `json:"version" binding:"required"`
		Count     int    `json:"count" binding:"required"`
		Page      int    `json:"page" binding:"required"`
		StartTime int64  `json:"startTime" binding:"required"`
		EndTime   int64  `json:"endTime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AppSpecificVersionDetails(c.MustGet(AppId).(string), params.Version, params.StartTime, params.EndTime, params.Count, params.Page)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//导出单个版本明细
func ExportAppSpecificVersionDetails(c *gin.Context) {
	type requestParams struct {
		Version   string `json:"version" binding:"required"`
		StartTime int64  `json:"startTime" binding:"required"`
		EndTime   int64  `json:"endTime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	url, err := model.ExportAppSpecificVersionDetails(c.MustGet(AppId).(string), params.Version, params.StartTime, params.EndTime)
	ret := make(map[string]string)
	ret["filename"] = url
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//版本平台信息
func AppSpecificVersionAsPlatform(c *gin.Context) {
	type requestParams struct {
		Version   string `json:"version" binding:"required"`
		StartTime int64  `json:"startTime" binding:"required"`
		EndTime   int64  `json:"endTime" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AppSpecificVersionAsPlatform(c.MustGet(AppId).(string), params.Version, params.StartTime, params.EndTime)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//-----------------------------用户管理--------------------//

//用户数量统计
func AdminUsersCount(c *gin.Context) {
	ret, err := model.AdminUsersCount(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//用户明细查询
func AdminUsersList(c *gin.Context) {
	type requestParams struct {
		Types []int   `json:"types" binding:"required"`
		Query *string `json:"Query" binding:"required"`
		Page  int64   `json:"page" binding:"required"`
		Count int64   `json:"count" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AdminUsersList(c.MustGet(AppId).(string), params.Types, *params.Query, params.Page, params.Count)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//封禁账号
func BanUser(c *gin.Context) {
	type requestParams struct {
		Id        string `json:"id" binding:"required"`
		StartTime int64  `json:"startTime" binding:"required"`
		EndTime   int64  `json:"endTime" binding:"required"`
		Reason    string `json:"reason"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.BanUser(c.MustGet(AppId).(string), c.MustGet(AdminId).(string), params.Id, params.Reason, params.EndTime)
	c.Set(ReqError, err)
}

//解封账号
func BanUserCancel(c *gin.Context) {
	type requestParams struct {
		Id string `json:"id" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.BanUserCancel(c.MustGet(AppId).(string), c.MustGet(AdminId).(string), params.Id)
	c.Set(ReqError, err)
}

/*func AdminGetUsers(c *gin.Context) {
	type requestParams struct {
		Types []string   `json:"types" binding:"required"`
		Query string   `json:"query"`
		Page int    `json:"page" binding:"required"`
		Count int   `json:"count" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.BanUserCancel(c.MustGet(AdminId).(string), params.Id)
	c.Set(ReqError, err)
}*/
//-----------------------群管理----------------------//
//群统计
func SetLimit(c *gin.Context) {
	type requestParams struct {
		MemberLimit     *int `json:"memberLimit" binding:"required"`
		RoomCreateLimit *int `json:"roomCreateLimit" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}
	err := model.SetLimit(c.MustGet(AppId).(string), *params.MemberLimit, *params.RoomCreateLimit)
	c.Set(ReqError, err)
}

//群统计
func GetLimit(c *gin.Context) {
	ret, err := model.GetLimit(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//群统计
func AdminRoomsCount(c *gin.Context) {
	ret, err := model.AdminRoomsCount(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//群明细查询
func AdminRoomsList(c *gin.Context) {
	type requestParams struct {
		Types []int   `json:"types" binding:"required"`
		Query *string `json:"Query" binding:"required"`
		Page  int64   `json:"page" binding:"required"`
		Count int64   `json:"count" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AdminRoomsList(c.MustGet(AppId).(string), params.Types, *params.Query, params.Page, params.Count)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//封禁群
func BanRoom(c *gin.Context) {
	type requestParams struct {
		Id        string `json:"id" binding:"required"`
		StartTime int64  `json:"startTime" binding:"required"`
		EndTime   int64  `json:"endTime" binding:"required"`
		Reason    string `json:"reason"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.BanRoom(c.MustGet(AppId).(string), c.MustGet(AdminId).(string), params.Id, params.Reason, params.EndTime)
	c.Set(ReqError, err)
}

//解封群
func BanRoomCancel(c *gin.Context) {
	type requestParams struct {
		Id string `json:"id" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.BanRoomCancel(c.MustGet(AppId).(string), c.MustGet(AdminId).(string), params.Id)
	c.Set(ReqError, err)
}

//设置推荐群
func SetRecommend(c *gin.Context) {
	type requestParams struct {
		Id        string `json:"id" binding:"required"`
		Recommend *int   `json:"recommend" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetRecommend(c.MustGet(AppId).(string), params.Id, *params.Recommend)
	c.Set(ReqError, err)
}

//查询管理员操作记录
func AdminOperateLog(c *gin.Context) {
	type requestParams struct {
		Types []int  `json:"types" binding:"required"`
		Query string `json:"query"`
		Page  int    `json:"page" binding:"required"`
		Count int    `json:"count" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.AdminOperateLog(c.MustGet(AppId).(string), params.Query, params.Types, params.Page, params.Count)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//获取文件
func GetExcel(c *gin.Context) {
	type requestParams struct {
		Filename string `form:"filename"`
	}
	var params requestParams
	if err := c.BindQuery(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filePath))
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	if ok, _ := model.FindExcel(types.ExcelAddr + params.Filename); !ok {
		c.Set(ReqError, result.NewError(result.AdminExportFailed))
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", params.Filename))
	c.File(types.ExcelAddr + params.Filename)
	c.Set(RespMiddleWareDisabled, true)
}

//获取模块启用状态
func ModuleEnable(c *gin.Context) {
	ret, err := model.ModuleEnable(c.MustGet(AppId).(string), c.MustGet(UserId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//获取oss配置信息
func GetOssConfig(c *gin.Context) {
	appId := c.MustGet(AppId)
	app := app.GetApp(appId.(string))
	if app == nil {
		c.Set(ReqError, result.NewError(result.ServerInterError).SetExtMessage(types.ERR_APPNOTFIND.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, app.OssConfig)
}

//--------------------------开屏广告-------------------------//
func CreateAd(c *gin.Context) {
	type requestParams struct {
		Name     string `json:"name"`
		Url      string `json:"url" binding:"required"`
		Duration int    `json:"duration" binding:"required"`
		Link     string `json:"link"`
		IsActive *int   `json:"isActive" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.CreateAd(c.MustGet(AppId).(string), params.Name, params.Url, params.Duration, params.Link, *params.IsActive)
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

func GetAllAd(c *gin.Context) {
	ret, err := model.GetAllAd(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

func SetAdName(c *gin.Context) {
	type requestParams struct {
		Id   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetAdName(params.Id, params.Name)
	c.Set(ReqError, err)
}

//激活广告
func ActiveAd(c *gin.Context) {
	type requestParams struct {
		Id       string `json:"id" binding:"required"`
		IsActive *int   `json:"isActive" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.ActiveAd(params.Id, *params.IsActive)
	c.Set(ReqError, err)
}

//删除广告
func DeleteAd(c *gin.Context) {
	type requestParams struct {
		Id string `json:"id" binding:"required"`
	}
	var params requestParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.DeleteAd(params.Id)
	c.Set(ReqError, err)
}

func UploadAd(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logAdmin.Error("UploadAd", "err", err.Error())
		c.Set(ReqError, result.NewError(result.AdUploadFailed))
		return
	}

	ext := path.Ext(header.Filename)
	if ext != ".jpg" {
		c.Set(ReqError, result.NewError(result.AdUploadFailed).SetExtMessage("仅支持jpg格式"))
		return
	}

	appId := c.MustGet(AppId).(string)
	//获取文件名
	//filename := header.Filename
	filename := "ad-" + appId + "-" + utility.ToString(utility.NowMillionSecond()) + ext

	//创建广告
	url, err := model.CreateAdFile(filename, &file)
	ret := make(map[string]string)
	ret["url"] = url

	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//-------------------------资金管理--------------------------//
func EditReward(c *gin.Context) {
	var params account.EditRewardParam
	params.BaseOpen = 1
	params.AdvanceOpen = 1
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.EditReward(c.MustGet(AppId).(string), &params)
	if err != nil {
		c.Set(ReqError, result.NewError(result.EditRewardFailed).SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
}

func ShowReward(c *gin.Context) {
	ret, err := model.ShowReward(c.MustGet(AppId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.ShowRewardFailed).SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

func RewardRule(c *gin.Context) {
	ret, err := model.RewardRule(c.MustGet(AppId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.ShowRewardFailed).SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//获取支持的币种
func CoinSupport(c *gin.Context) {
	coins, err := model.CoinSupport(c.MustGet(AppId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	ret := make(map[string]interface{})
	ret["coins"] = coins
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//获取奖励统计信息系
func RewardStatistics(c *gin.Context) {
	ret, err := model.RewardStatistics(c.MustGet(AppId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//获取奖励列表
func RewardList(c *gin.Context) {
	var params account.RewardListParam
	params.Page = 1
	params.Size = 15
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	ret, err := model.RewardList(c.MustGet(AppId).(string), &params)
	if err != nil {
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//获取红包手续费列表
func RPFeeConfig(c *gin.Context) {
	ret, err := model.RPFeeConfig(c.MustGet(AppId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//设置红包手续费
func SetRPFeeConfig(c *gin.Context) {
	var params account.SetRPFeeParam
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	err := model.SetRPFeeConfig(c.MustGet(AppId).(string), &params)
	if err != nil {
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
}

//红包手续费统计信息
func RPFeeStatistics(c *gin.Context) {
	ret, err := model.RPFeeStatistics(c.MustGet(AppId).(string))
	if err != nil {
		c.Set(ReqError, result.NewError(result.ServiceReqFailed).JustShowExtMsg().SetExtMessage(err.Error()))
		return
	}
	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

//-------------------------------认证------------------//
func VerifyApprove(c *gin.Context) {
	var params struct {
		Id     string `json:"id" binding:"required"`
		Accept *int   `json:"accept" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	switch *params.Accept {
	case types.VerifyApproveReject:
	case types.VerifyApproveAccept:
	case types.VerifyApproveFeeBack:
	default:
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, "accept type not exists"))
		return
	}

	err := model.VerifyApprove(c.MustGet(AppId).(string), params.Id, *params.Accept)
	c.Set(ReqError, err)
}

func PersonalVerifyList(c *gin.Context) {
	var params struct {
		Search *string `json:"search"`
		Page   int     `json:"page"`
		Size   int     `json:"size"`
		State  *int    `json:"state"`
	}
	params.Page = 1
	params.Size = 16

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	count, apply, err := model.PersonalVerifyList(c.MustGet(AppId).(string), params.Search, params.Page, params.Size, params.State)
	if err != nil {
		c.Set(ReqError, err)
		return
	}

	list := make([]interface{}, 0)
	for _, v := range apply {
		item := map[string]interface{}{
			"id":          v.Id,
			"uid":         v.Uid,
			"account":     v.Phone,
			"name":        v.Username,
			"avatar":      v.Avatar,
			"currency":    v.Currency,
			"amount":      utility.ToString(v.Amount),
			"state":       v.State,
			"description": v.Description,
			"coinState":   v.FeeState,
			"time":        v.VerifyApply.UpdateTime,
		}
		list = append(list, item)
	}

	ret := make(map[string]interface{})
	ret["count"] = utility.ToString(count)
	ret["list"] = list

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

func RoomVerifyList(c *gin.Context) {
	var params struct {
		Search *string `json:"search"`
		Page   int     `json:"page"`
		Size   int     `json:"size"`
		State  *int    `json:"state"`
	}
	params.Page = 1
	params.Size = 16

	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	count, apply, err := model.RoomVerifyList(c.MustGet(AppId).(string), params.Search, params.Page, params.Size, params.State)
	if err != nil {
		c.Set(ReqError, err)
		return
	}

	list := make([]interface{}, 0)
	for _, v := range apply {
		item := map[string]interface{}{
			"id":          v.VerifyApply.Id,
			"master":      v.User.Phone,
			"account":     v.Room.MarkId,
			"name":        v.Room.Name,
			"avatar":      v.Room.Avatar,
			"currency":    v.Currency,
			"amount":      utility.ToString(v.Amount),
			"state":       v.State,
			"description": v.Description,
			"coinState":   v.FeeState,
			"time":        v.VerifyApply.UpdateTime,
		}
		list = append(list, item)
	}

	ret := make(map[string]interface{})
	ret["count"] = utility.ToString(count)
	ret["list"] = list

	c.Set(ReqError, nil)
	c.Set(ReqResult, ret)
}

func VerifyGetConfigByType(c *gin.Context) {
	var params struct {
		Type *int `json:"type"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	appId := c.MustGet(AppId).(string)
	//调用打币接口
	app := app.GetApp(appId)
	if app == nil {
		c.Set(ReqError, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error()))
		return
	}
	defCoin := app.MainCoin

	fee, err := model.VerifyGetConfig(c.MustGet(AppId).(string))

	type queryFee struct {
		Type     int    `json:"type"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	}

	feeRlt := map[int]*queryFee{
		types.VerifyForUser: {
			Type:     types.VerifyForUser,
			Currency: defCoin,
			Amount:   "0",
		},
		types.VerifyForRoom: {
			Type:     types.VerifyForRoom,
			Currency: defCoin,
			Amount:   "0",
		},
	}

	for _, v := range fee {
		if f, ok := feeRlt[v.Type]; ok {
			f.Currency = v.Currency
			f.Amount = utility.ToString(v.Amount)
		}
	}

	list := make([]interface{}, 0)
	//查询全部
	if params.Type == nil || *params.Type == 0 {
		list = append(list, feeRlt[types.VerifyForUser])
		list = append(list, feeRlt[types.VerifyForRoom])
	} else {
		list = append(list, feeRlt[*params.Type])
	}

	ret := map[string]interface{}{
		"list": list,
	}

	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

func VerifyGetConfig(c *gin.Context) {
	fee, err := model.VerifyGetConfig(c.MustGet(AppId).(string))
	if len(fee) < 2 {
		appId := c.MustGet(AppId).(string)
		//调用打币接口
		app := app.GetApp(appId)
		if app == nil {
			c.Set(ReqError, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error()))
			return
		}
		currency := app.MainCoin

		fee = []*types.VerifyFee{
			{
				AppId:    appId,
				Type:     1,
				Currency: currency,
				Amount:   0,
			},
			{
				AppId:    appId,
				Type:     2,
				Currency: currency,
				Amount:   0,
			},
		}
	}

	ret := map[string]interface{}{
		"config": fee,
	}
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

func VerifySetFee(c *gin.Context) {
	var params struct {
		Config []*types.VerifyFee `json:"config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.Set(ReqError, result.NewError(result.ParamsError).SetChildErr(result.ServiceChat, nil, err.Error()))
		return
	}

	appId := c.MustGet(AppId).(string)
	//调用打币接口
	app := app.GetApp(appId)
	if app == nil {
		c.Set(ReqError, result.NewError(result.PermissionDeny).SetExtMessage(types.ERR_APPNOTFIND.Error()))
		return
	}
	currency := app.MainCoin
	if currency == "" {
		c.Set(ReqError, result.NewError(result.ServerInterError).SetExtMessage("未设置主币种"))
		return
	}

	for _, v := range params.Config {
		v.AppId = appId
		v.Currency = currency
	}

	err := model.VerifySetFee(params.Config)
	c.Set(ReqError, err)
}

func VerifyFeeStatistics(c *gin.Context) {
	ret, err := model.VerifyFeeStatistics(c.MustGet(AppId).(string))
	c.Set(ReqError, err)
	c.Set(ReqResult, ret)
}

//获取用户在线数量

func GetUserNums(c *gin.Context) {
	ret := model.UserOnlineNums()
	c.Set(ReqResult, ret)
	c.Set(ReqError, nil)
}

//获取用户链接数量

func GetUserClientNums(c *gin.Context) {
	ret := model.UserClientNums()
	c.Set(ReqResult, ret)
	c.Set(ReqError, nil)
}
