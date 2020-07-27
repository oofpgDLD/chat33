package api

import (
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/33cn/chat33/comet/ws"
	"github.com/33cn/chat33/orm"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/sessions"
	"github.com/inconshreveable/log15"
	"gopkg.in/go-playground/validator.v8"
)

var logBase = log15.New("module", "api/base")
var (
	sessionName  = "session-login"
	sessionAdmin = "background"
	sessionKey   = "veryHardToGuess"
	store        *sessions.CookieStore
)

const apiAuthTypeToken = false

var cfg *types.Config

// WriteMsg 给接口请求返回消息
func WriteMsg(w io.Writer, msg []byte) {
	var err error
	_, err = w.Write(msg)
	if err != nil {
		//logAPI.Warn("Write err", "err ", err)
	}
}

const (
	RespMiddleWareDisabled = "RespMiddleWareDisabled"

	UserId     = "userId"
	Token      = "token"
	AppId      = "appId"
	DeviceName = "deviceName"
	DeviceType = "deviceType"
	Uuid       = "uuid"
	Version    = "version"
	Time       = "time"

	ReqError  = "error"
	ReqResult = "result"

	AdminId = "adminId"
)

func RespMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if v, ok := c.Get(RespMiddleWareDisabled); ok && v == true {
			return
		}
		err := c.MustGet(ReqError)
		var info interface{}
		var errMsg string
		var errCode int
		var errData map[string]interface{}
		if err == nil {
			errCode = result.CodeOK
		} else {
			errCode = err.(*result.Error).ErrorCode
			errData = err.(*result.Error).Data
			errMsg = err.(*result.Error).Error()
		}
		if result, ok := c.Get(ReqResult); !ok || result == nil {
			info = gin.H{}
		} else {
			info = result
		}
		if errData != nil {
			info = errData
		}
		ret := result.ComposeHttpAck(errCode, errMsg, info)
		c.PureJSON(http.StatusOK, ret)
		//c.PureJSON()
	}
}

func authToken(c *gin.Context) {
	appId := c.GetHeader("FZM-APP-ID")
	token := c.GetHeader("FZM-AUTH-TOKEN")
	device := c.GetHeader("FZM-DEVICE")

	//鉴权  通过token查询user
	user, err := orm.GetUserInfoByToken(appId, token)
	if err != nil {
		if err == types.ERR_LOGIN_EXPIRED {
			c.Set(ReqError, result.NewError(result.LoginExpired))
			c.Abort()
			return
		}
		c.Set(ReqError, result.NewError(result.TokenLoginFailed).SetExtMessage(err.Error()))
		c.Abort()
		return
	}

	//user, err := orm.GetUserInfoByUid(appId, tokenInfo.Uid)
	//if err != nil {
	//	c.Set(ReqError, result.NewError(result.DbConnectFail))
	//	c.Abort()
	//	return
	//}
	if user == nil {
		c.Set(ReqError, result.NewError(result.UserNotReg))
		c.Abort()
		return
	}
	userId := user.UserId

	time := utility.NowMillionSecond()
	//TODO 之前的接口为空，以后可扩展
	uuid := ""

	c.Set(UserId, userId)
	c.Set(AppId, appId)
	c.Set(Token, token)
	c.Set(DeviceType, device)
	c.Set(Time, time)
	c.Set(Uuid, uuid)
	logBase.Info("auth info", "appId", appId, "userId", userId, "device type", device, "token", "time", time)
}

func authSession(c *gin.Context) {
	session, err := store.Get(c.Request, sessionName)
	if session == nil {
		logBase.Error("AuthMiddleWare get session", "err", err.Error())
		c.Set(ReqError, result.NewError(result.ServerInterError))
		c.Abort()
		return
	}

	userId := session.Values["user_id"]
	appId := session.Values["appId"]
	token := session.Values["token"]
	device := c.GetHeader("FZM-DEVICE")
	uuid := c.GetHeader("FZM-UUID")
	logUser.Info("authSession success", "appId", appId, "userId", userId, "token", token)
	if userId == nil || appId == nil || token == nil {
		c.Set(ReqError, result.NewError(result.LoginExpired))
		c.Abort()
		return
	}
	c.Set(UserId, userId)
	c.Set(AppId, appId)
	c.Set(Token, token)
	c.Set(Uuid, uuid)
	c.Set(DeviceType, device)
}

func authTest(c *gin.Context) {
	userId := c.GetHeader("FZM-USER-ID")
	appId := c.GetHeader("FZM-APP-ID")
	token := c.GetHeader("FZM-AUTH-TOKEN")
	device := c.GetHeader("FZM-DEVICE")

	c.Set(UserId, userId)
	c.Set(AppId, appId)
	c.Set(Token, token)
	c.Set(DeviceType, device)
}

func AuthMiddleWare() gin.HandlerFunc {
	/*//用于测试 不经过托管账户
	if true {
		return authTest
	}*/
	if apiAuthTypeToken {
		return authToken
	} else {
		return authSession
	}
}

func AdminAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, sessionAdmin)
		if session == nil {
			logBase.Error("AdminAuthMiddleWare get session", "err", err.Error())
			c.Set(ReqError, result.NewError(result.ServerInterError))
			c.Abort()
			return
		}

		adminId := session.Values["admin_id"]
		appId := session.Values["appId"]
		headerAppId := c.GetHeader("FZM-APP-ID")
		if adminId == nil || appId == nil || headerAppId != appId {
			c.Set(ReqError, result.NewError(result.LoginExpired))
			c.Abort()
			return
		}
		c.Set(AdminId, adminId)
		c.Set(AppId, appId)
	}
}

func ParseHeaderMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, sessionName)
		if session == nil {
			logBase.Error("ParseHeaderMiddleWare get session", "err", err.Error())
			c.Set(ReqError, result.NewError(result.ServerInterError))
			c.Abort()
			return
		}

		userId := session.Values["user_id"]
		appId := c.GetHeader("FZM-APP-ID")
		token := c.GetHeader("FZM-AUTH-TOKEN")
		device := c.GetHeader("FZM-DEVICE")
		uuid := c.GetHeader("FZM-UUID")
		devName := c.GetHeader("FZM-DEVICE-NAME")
		version := c.GetHeader("FZM-VERSION")

		c.Set(UserId, userId)
		c.Set(AppId, appId)
		c.Set(Token, token)
		c.Set(DeviceType, device)
		c.Set(Uuid, uuid)
		c.Set(DeviceName, devName)
		c.Set(Version, version)
	}
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*") //Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,FZM-APP-ID
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

func Init(c *types.Config) *gin.Engine {
	cfg = c
	store = sessions.NewCookieStore([]byte(sessionKey))
	r := gin.Default()
	//Custom Validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("LegalNumber", CheckNumber)
		if err != nil {
			logBase.Error("gin custom validator init error, field:LegalNumber")
		}
	}

	// websocket
	r.GET("/ws", func(context *gin.Context) {
		if apiAuthTypeToken {
			ws.ServeWs(store, context, getUserInfoV2)
		} else {
			ws.ServeWs(store, context, ws.GetUserInfo)
		}
	})

	//inner api
	inner := r.Group("/inner")
	inner.GET("/ping", Ping)
	inner.GET("/allTS", AllTickerSymbol, RespMiddleWare())
	//inner.GET("/find-user-session", FindUserSession)
	appConfig := inner.Group("/app", RespMiddleWare())
	appConfig.GET("/update", AppUpdate)
	appConfig.GET("/watch", AppConfig)
	appConfig.POST("/oss/update", OssUpdate)

	root := r.Group("", RespMiddleWare())
	// 获取币种信息
	//root.POST("/coin/coinInfo", api.CoinList)
	// （弃用）获取app信息
	root.POST("/app/appInfo", AppInfo)
	//版本更新
	root.POST("/version", ParseHeaderMiddleWare(), VersionController)
	//（弃用） 传入个推cid
	root.POST("/GTCid", AuthMiddleWare(), SaveGTCid)

	//应用启动
	root.POST("/open", ParseHeaderMiddleWare(), AppOpen)
	//发送消息
	root.POST("/push", ParseHeaderMiddleWare(), Push)

	public := root.Group("/public")
	//public.POST("/reward-rule", ParseHeaderMiddleWare(), RewardRule)
	public.Any("/reward-rule", ParseHeaderMiddleWare(), Cors(), RewardRule)
	public.POST("/module-state", AuthMiddleWare(), ModuleEnable)
	public.Any("/verification-fee", ParseHeaderMiddleWare(), Cors(), VerifyGetConfigByType)
	public.POST("/oss-config", AuthMiddleWare(), GetOssConfig)

	chat33 := r.Group("/chat33", RespMiddleWare())
	//精确搜索用户或群
	chat33.POST("/search", ParseHeaderMiddleWare(), ClearlySearch)

	//获取群会话秘钥
	chat33.POST("/roomSessionKey", AuthMiddleWare(), RoomSessionKey)

	//获取开屏广告
	chat33.POST("/getAdvertisement", ParseHeaderMiddleWare(), Advertisement)

	chat33.Use(AuthMiddleWare())
	{
		//上传秘钥接口
		chat33.POST("/uploadSecretKey", UploadSecretKey)
		//获取入群/好友申请列表
		chat33.POST("/applyList", GetApplyList)
		//获取未处理申请数量
		chat33.POST("/unreadApplyNumber", GetUnreadApplyNumber)
		//撤回指定的一条消息
		chat33.POST("/RevokeMessage", RevokeMsg)
		//批量撤回 文件、图片、视频消息
		chat33.POST("/RevokeFiles", RevokeFiles)
		//阅读指定一条阅后即焚消息
		chat33.POST("/readSnapMsg", ReadSnapMsg)
		//转发消息
		chat33.POST("/forward", ForwardMsg)
		//转发加密消息
		chat33.POST("/encryptForward", ForwardEncryptMsg)
	}

	room := r.Group("/room", RespMiddleWare())
	//获取群在线人数 TODO 是否需要验证登录
	room.POST("/getOnlineNumber", GetOnlineNumber)
	//获取推荐群列表
	room.POST("/recommend", ParseHeaderMiddleWare(), RecommendRooms)
	room.Use(AuthMiddleWare())
	{
		//判断用户是否在群里
		room.POST("/userIsInRoom", UserIsInRoom)
		//创建群
		room.POST("/create", CreateRoom)
		//删除群
		room.POST("/delete", RemoveRoom)
		//退出群
		room.POST("/loginOut", LoginOutRoom)
		//踢出群
		room.POST("/kickOut", KickOutRoom)
		//获取群列表
		room.POST("/list", GetRoomList)
		//获取群信息
		room.POST("/info", GetRoomInfo)
		//获取群成员列表
		room.POST("/userList", GetRoomUserList)
		//获取群成员信息
		room.POST("/userInfo", GetRoomUserInfo)
		//管理员设置群
		room.POST("/setPermission", AdminSetPermission)
		//修改群名称
		room.POST("/setName", SetRoomName)
		//修改群头像
		room.POST("/setAvatar", SetRoomAvatar)
		//群内用户身份设置
		room.POST("/setLevel", SetLevel)
		//群成员设置免打扰
		room.POST("/setNoDisturbing", SetNoDisturbing)
		//群成员设置群内昵称
		room.POST("/setMemberNickname", SetMemberNickname)
		//邀请入群
		room.POST("/joinRoomInvite", JoinRoomInvite)
		//申请入群
		room.POST("/joinRoomApply", JoinRoomApply)
		//批量申请入群
		room.POST("/batchJoinRoomApply", BatchJoinRoomApply)
		//入群申请处理
		room.POST("/joinRoomApprove", JoinRoomApprove)
		//群消息记录
		room.POST("/chatLog", GetRoomChatLog)
		//群文件记录
		room.POST("/historyFiles", GetRoomFileList)
		//群图片记录
		room.POST("/historyPhotos", GetRoomPhotoList)
		//群成员设置群置顶
		room.POST("/stickyOnTop", SetStickyOnTop)
		//搜索群成员信息
		room.POST("/searchMember", GetRoomSearchMember)
		//设置群禁言列表
		room.POST("/setMutedList", SetRoomMuted)
		//设置群禁言
		room.POST("/setMutedSingle", SetRoomMutedSingle)
		//获取群成公告
		room.POST("/systemMsgs", GetSystemMsg)
		//发布群成公告
		room.POST("/sendSystemMsgs", SetSystemMsg)
		//群认证申请
		room.POST("/verification-apply", VerifyApplyForRoom)
		//用户加v审核状态
		room.POST("/verification-state", RoomVerificationState)
	}

	friend := r.Group("/friend", RespMiddleWare())
	//校验答案
	friend.POST("/checkAnswer", CheckAnswer)
	friend.Use(AuthMiddleWare())
	{
		// 判断是否是好友
		friend.POST("/isFriend", IsFriend)
		// 获取好友列表
		friend.POST("/list", FriendList)
		// 添加好友
		friend.POST("/add", AddFriend)
		// 添加好友  不需要验证
		friend.POST("/addForGM", AddFriendNotConfirm)
		// 删除好友
		friend.POST("/delete", DeleteFriend)
		// 处理好友请求
		friend.POST("/response", HandleFriendRequest)
		// 修改好友备注
		friend.POST("/setRemark", FriendSetRemark)
		// 修改好友扩展备注
		friend.POST("/setExtRemark", FriendSetExtRemark)
		// 设置好友免打扰
		friend.POST("/setNoDisturbing", FriendSetDND)
		// 设置好友消息置顶
		friend.POST("/stickyOnTop", FriendSetTop)
		//获取好友消息记录
		friend.POST("/chatLog", FriendChatLog)
		//获取文件列表
		friend.POST("/historyFiles", FriendFliesLog)
		//获取图片列表
		friend.POST("/historyPhotos", FriendPhotosLog)
		//获取所有好友未读消息统计
		friend.POST("/unread", GetAllFriendUnreadMsg)
		//发送将阅后即焚消息截图
		friend.POST("/printScreen", PrintScreen)
		//设置答案验证
		friend.POST("/question", Question)
		//设置好友验证
		friend.POST("/confirm", Confirm)

		//加入黑名单
		friend.POST("/block", BlockFriend)
		//移出黑名单
		friend.POST("/unblock", UnblockFriend)
		//黑名单列表
		friend.POST("/blocked-list", FriendBlockList)
	}

	user := r.Group("/user", RespMiddleWare())
	//是否设置支付密码
	user.POST("/isSetPayPwd", ParseHeaderMiddleWare(), UserIsSetPayPwd)
	//修改支付密码
	user.POST("/setPayPwd", ParseHeaderMiddleWare(), UserSetPayPwd)
	//验证支付密码
	user.POST("/checkPayPwd", ParseHeaderMiddleWare(), UserCheckPayPwd)
	// 用户token登录
	user.POST("/tokenLogin", ParseHeaderMiddleWare(), UserTokenLogin)
	// 验证验证码(返回Token)
	user.POST("/phoneLogin", ParseHeaderMiddleWare(), PhoneLogin)
	// 发送验证码
	user.POST("/sendCode", ParseHeaderMiddleWare(), SendCode)
	// 单次邀请奖励列表
	user.POST("/single-invite-info", ParseHeaderMiddleWare(), SingleInviteInfos)
	// 叠加邀请奖励列表
	user.POST("/accumulate-invite-info", ParseHeaderMiddleWare(), AccumulateInviteInfos)
	// 邀请奖励统计信息
	user.POST("/invite-statistics", ParseHeaderMiddleWare(), InviteStatistics)
	// 用户退出登陆
	user.POST("/logout", UserLoginOut)
	user.Use(AuthMiddleWare())
	{
		//更新推送deviceToken (友盟)
		user.POST("/set-device-token", UpdatePushDeviceToken)
		// 用户修改昵称
		user.POST("/editNickname", UserEditNickname)
		// 用户自定义头像
		user.POST("/editAvatar", UserEditAvatar)
		//查看用户详情(弃用)
		user.POST("/userInfo", FriendInfo)
		//查看用户详情(根据Uid批量获取)
		user.POST("/usersInfo", UserInfo)
		//通过uid(地址)获取用户信息
		user.POST("/userInfoByUid", UserListByUid)
		//查看用户设置
		user.POST("/userConf", UserConf)
		//设置邀请入群确认
		user.POST("/set-invite-confirm", SetInviteConfirm)
		//加v认证请求
		user.POST("/verification-apply", VerifyApplyForUser)
		//用户加v审核状态
		user.POST("/verification-state", UserVerificationState)
		//上链
		user.POST("/isChain", IsChain)
	}

	redPacket := r.Group("/red-packet", RespMiddleWare())
	// 查询红包余额
	redPacket.POST("/balance", ParseHeaderMiddleWare(), Balance)
	// 获取代币信息
	redPacket.POST("/coin", ParseHeaderMiddleWare(), GetCoinInfo)

	redPacket.Use(AuthMiddleWare())
	{
		// 发红包
		redPacket.POST("/send", Send)
		// 领红包
		redPacket.POST("/receive-entry", ReceiveEntry)
		// 红包详情
		redPacket.POST("/detail", RedEnvelopeDetail)
		// 红包领取详情
		redPacket.POST("/receiveDetail", RedEnvelopeReceiveDetail)
		// 红包统计信息
		redPacket.POST("/statistic", RPStatistics)
	}

	pay := r.Group("/pay", RespMiddleWare())
	pay.Use(AuthMiddleWare())
	{
		// 付款
		pay.POST("/payment", Payment)
	}

	r.GET("/getExcel2", GetExcel)
	manage := r.Group("/manage", RespMiddleWare())
	manage.POST("/login", AdminLogin)
	manage.GET("/downloadExcel", GetExcel)
	//统计在线用户数
	manage.GET("/userNums", GetUserNums)
	//统计在线用户连接数
	manage.GET("/userClientNums", GetUserClientNums)
	manage.Use(AdminAuthMiddleWare())
	{
		ad := manage.Group("/startPage")
		ad.POST("create", CreateAd)
		ad.POST("getAll", GetAllAd)
		ad.POST("upload", UploadAd)
		ad.POST("setName", SetAdName)
		ad.POST("active", ActiveAd)
		ad.POST("delete", DeleteAd)

		funds := manage.Group("/reward")
		funds.POST("/set-config", EditReward)
		funds.POST("/get-config", ShowReward)
		funds.POST("/coin-support", CoinSupport)
		funds.POST("/statistics", RewardStatistics)
		funds.POST("/list", RewardList)

		verification := manage.Group("/verification")
		//认证审批
		verification.POST("/approve", VerifyApprove)
		//个人认证申请列表
		verification.POST("/personal-list", PersonalVerifyList)
		//群认证申请列表
		verification.POST("/group-list", RoomVerifyList)
		//获取认证手续费配置
		verification.POST("/fee-config", VerifyGetConfig)
		//设置认证收费
		verification.POST("/set-fee", VerifySetFee)
		//认证收费统计
		verification.POST("/fee-statistics", VerifyFeeStatistics)

		rp := manage.Group("/red-packet")
		//查询红包手续费
		rp.POST("/fee-config", RPFeeConfig)
		//设置红包手续费
		rp.POST("/set-fee", SetRPFeeConfig)
		//手续费统计信息
		rp.POST("/fee-statistics", RPFeeStatistics)

		//近期数据
		manage.POST("/account", AdminAccount)
		//汇总数据
		//近期数据
		manage.POST("/latestData", LatestData)
		//平台统计
		manage.POST("/platform", LatestDataPlatfrom)
		//明细数据
		manage.POST("/sumDetails", SumDetails)
		//导出明细数据
		manage.POST("/sum/export", ExportSumDetails)
		//折线统计图
		manage.POST("/sumChart", SumGraph)

		//版本
		//获取所有版本
		manage.POST("/version/versions", GetAppVersions)
		//版本统计折线图
		manage.POST("/version/chart", VersionGraph)
		//版本统计明细
		manage.POST("/version/detailAll", AppVersionDetails)
		//导出版本统计明细
		manage.POST("/version/exportAll", ExportAppVersionDetails)
		//具体版本明细
		manage.POST("/version/detailVersion", AppSpecificVersionDetails)
		//导出具体版本明细
		manage.POST("/version/export", ExportAppSpecificVersionDetails)

		//版本平台信息
		manage.POST("/version/platform", AppSpecificVersionAsPlatform)

		//用户统计
		manage.POST("/user/count", AdminUsersCount)
		//用户列表
		manage.POST("/user/list", AdminUsersList)
		//封禁用户
		manage.POST("/user/ban", BanUser)
		//解封用户
		manage.POST("/user/cancelBan", BanUserCancel)

		//数量上限设置
		manage.POST("/room/setLimit", SetLimit)
		//数量上限获取
		manage.POST("/room/getLimit", GetLimit)

		//群统计
		manage.POST("/room/count", AdminRoomsCount)
		//群列表
		manage.POST("/room/list", AdminRoomsList)
		//封禁群
		manage.POST("/room/ban", BanRoom)
		//解封群
		manage.POST("/room/cancelBan", BanRoomCancel)
		//设置推荐群
		manage.POST("/room/setRecommend", SetRecommend)

		//操作日志
		manage.POST("/log", AdminOperateLog)
	}

	//赞赏模块
	praise := r.Group("/praise", RespMiddleWare())
	praise.Use(AuthMiddleWare())
	{
		// 赞赏列表
		praise.POST("/list", PraiseList)
		// 赏赐详情
		praise.POST("/details", PraiseDetails)
		// 赏赐详情列表
		praise.POST("/detailList", PraiseDetailList)
		// 打赏
		praise.POST("/reward", PraiseReward)
		// 点赞
		praise.POST("/like", PraiseLike)
		//打赏用户
		praise.POST("/rewardUser", PraiseRewardUser)
		//赞赏榜单
		praise.POST("/leaderboard", LeaderBoard)
		//历史榜单
		praise.POST("/leaderboardHistory", LeaderBoardHistory)
	}

	return r
}

func CheckNumber(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	val := field.Interface()

	switch val.(type) {
	case int:
		return true
	case string:
		if val.(string) == "" {
			return true
		}
		_, err := strconv.ParseInt(val.(string), 10, 64)
		if err != nil {
			return false
		}
		return true
	case int64:
		return true
	default:
		//utility_log.Error("func ToInt error unknow type")
		return false
	}
}

//支持 token 获取用户信息
func getUserInfoV2(store *sessions.CookieStore, c *gin.Context) (userId, device, uuid, appId string, level int, loginTime int64) {
	authToken(c)
	userId = utility.ToString(c.MustGet(UserId))
	device = utility.ToString(c.MustGet(DeviceType))
	uuid = utility.ToString(c.MustGet(Uuid))
	appId = utility.ToString(c.MustGet(AppId))
	loginTime = utility.ToInt64(c.MustGet(Time))

	switch device {
	case types.DevicePC:
	case types.DeviceAndroid:
	case types.DeviceIOS:
	default:
		device = ""
	}

	if userId == "" || device == "" {
		level = router.VISITOR
	} else {
		level = router.NOMALUSER
	}
	return
}
