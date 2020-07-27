package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/33cn/chat33/app"
	"github.com/33cn/chat33/orm"
	push "github.com/33cn/chat33/pkg/u-push"
	android_push "github.com/33cn/chat33/pkg/u-push/android"
	ios_push "github.com/33cn/chat33/pkg/u-push/ios"
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

var pushLog = log15.New("module", "push")

//
type IPusher interface {
	singlePush(title, text, targetId string, channelType int) error
}

type androidPusher struct {
	AppKey          string
	AppMasterSecret string
	MiActivity      string
	DeviceToken     string
}

type iOSPusher struct {
	AppKey          string
	AppMasterSecret string
	MiActivity      string
	DeviceToken     string
}

func checkDeviceToken(app *types.App, userId, deviceToken string) (IPusher, error) {
	deviceType, nowUser, err := orm.FindUserIdByDeviceToken(deviceToken)
	if err != nil {
		return nil, errors.New("check device token failed :" + err.Error())
	}

	if nowUser != userId {
		return nil, nil
	}

	var pusher IPusher
	switch deviceType {
	case types.DeviceAndroid:
		pusher = &androidPusher{
			AppKey:          app.PushAppKey.Android,
			AppMasterSecret: app.PushAppMasterSecret.Android,
			MiActivity:      app.PushMiActive.Android,
			DeviceToken:     deviceToken,
		}
	case types.DeviceIOS:
		pusher = &iOSPusher{
			AppKey:          app.PushAppKey.IOS,
			AppMasterSecret: app.PushAppMasterSecret.IOS,
			MiActivity:      app.PushMiActive.IOS,
			DeviceToken:     deviceToken,
		}
	default:
		return nil, errors.New(fmt.Sprintf("checkDeviceToken device : %s type illegal", deviceType))
	}
	return pusher, nil
}

func (t *androidPusher) singlePush(title, text, targetId string, channelType int) error {
	var client push.PushClient
	unicast := android_push.NewAndroidUnicast(t.AppKey, t.AppMasterSecret)

	fmt.Println(t.AppKey, t.AppMasterSecret, t.DeviceToken, title, text)
	unicast.SetDeviceToken(t.DeviceToken)
	unicast.SetTitle(title)
	unicast.SetText(text)
	unicast.GoCustomAfterOpen("")
	unicast.SetDisplayType(push.NOTIFICATION)
	unicast.SetMipush(true, t.MiActivity)
	// 线上模式
	unicast.SetReleaseMode()
	// Set customized fields
	unicast.SetExtraField("targetId", targetId)
	unicast.SetExtraField("channelType", utility.ToString(channelType))
	return client.Send(unicast)
}

func (t *iOSPusher) singlePush(title, text, targetId string, channelType int) error {
	var client push.PushClient
	unicast := ios_push.NewIOSUnicast(t.AppKey, t.AppMasterSecret)

	fmt.Println(t.AppKey, t.AppMasterSecret, t.DeviceToken, title, text)
	unicast.SetDeviceToken(t.DeviceToken)

	//unicast.SetAlert("IOS 单播测试")
	unicast.SetAlertJson(push.IOSAlert{
		Title: title,
		Body:  text,
	})
	//unicast.SetBadge(0)
	unicast.SetSound("default")
	//// 测试模式
	//unicast.SetTestMode()
	// 线上模式
	unicast.SetReleaseMode()

	// Set customized fields
	unicast.SetCustomizedField("targetId", targetId)
	unicast.SetCustomizedField("channelType", utility.ToString(channelType))

	return client.Send(unicast)
}

func PushToFriend(appId, userId, friendId, text string, msg *proto.Proto) error {
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	user, err := orm.FindFriendById(friendId, userId)
	if err != nil {
		return err
	}
	//如果不是红包消息 需要判断消息免打扰
	if msg.GetMsgType() != types.RedPack {
		//查询对方是否对你设置了消息免打扰
		if user == nil || user.DND == types.NoDisturbingOn {
			//消息免打扰  直接推送websocket
			return nil
		}
	}

	//如果是焚毁通知 则不推送
	if msg.GetMsgType() == types.Alert {
		content := msg.GetMsg()
		alertType := utility.ToInt(content["type"])
		if alertType == types.AlertHadBurntMsg {
			return nil
		}
	}

	friend, err := orm.GetUserInfoById(friendId)
	if err != nil {
		return err
	}

	//非透传消息 title:自己名称
	var userName string
	if user != nil {
		if user.Remark != "" {
			userName = user.Remark
		} else {
			userName = user.Username
		}
	}
	title := userName

	var pusher IPusher
	if pusher, err = checkDeviceToken(app, friendId, friend.DeviceToken); pusher == nil {
		return err
	}
	if cfg.Log.Level == "debug" {
		b, err := json.Marshal(pusher)
		if err != nil {
		}
		pushLog.Debug("push to room", "friendId", friendId, "userId", userId, "pusher", string(b))
	}
	return pusher.singlePush(title, text, userId, types.ToUser)
}

func PushToMember(appId, userId, roomId string, members map[string]string, msg *proto.Proto) error {
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	//如果是焚毁通知
	if msg.GetMsgType() == types.Alert {
		content := msg.GetMsg()
		alertType := utility.ToInt(content["type"])
		if alertType == types.AlertHadBurntMsg {
			return nil
		}
	}

	//筛选群中未设置消息免打扰的 红包消息除外
	if msg.GetChannelType() == types.ToRoom {
		if msg.GetMsgType() != types.RedPack {
			//查找设置消息免打扰的
			dndUsers, err := orm.FindSetNoDisturbingMembers(msg.GetTargetId())
			if err != nil {
				return err
			}

			//移除设置了消息免打扰的
			for _, v := range dndUsers {
				key := v.UserId
				if _, ok := members[key]; ok {
					delete(members, key)
				}
			}
		}
	}

	var roomName string
	info, err := orm.FindRoomById(msg.GetTargetId(), types.RoomNotDeleted)
	if err == nil && info != nil {
		roomName = info.Name
	}

	title := roomName
	text := msg.GetGTMsg(userId)

	for k, v := range members {
		var pusher IPusher
		if pusher, err = checkDeviceToken(app, k, v); pusher == nil {
			if err != nil {
				//TODO 输出错误日志
			}
			continue
		}
		err = pusher.singlePush(title, text, roomId, types.ToRoom)
	}
	return nil
}

func PushToRoom(appId, userId, roomId string, msg *proto.Proto) error {
	app := app.GetApp(appId)
	if app == nil {
		return types.ERR_APPNOTFIND
	}

	//如果是焚毁通知
	if msg.GetMsgType() == types.Alert {
		content := msg.GetMsg()
		alertType := utility.ToInt(content["type"])
		if alertType == types.AlertHadBurntMsg {
			return nil
		}
	}

	//查询离线成员
	var members = make(map[string]string)
	users, _ := orm.FindNotDelMembers(roomId, types.SearchAll)
	for _, v := range users {
		if user, ok := router.GetUser(v.RoomMember.UserId); !ok || len(user.GetClients()) == 0 {
			members[v.RoomMember.UserId] = v.DeviceToken
		}
	}

	//筛选群中未设置消息免打扰的 红包消息除外
	if msg.GetChannelType() == types.ToRoom {
		if msg.GetMsgType() != types.RedPack {
			//查找设置消息免打扰的
			dndUsers, err := orm.FindSetNoDisturbingMembers(msg.GetTargetId())
			if err != nil {
				return err
			}

			//移除设置了消息免打扰的
			for _, v := range dndUsers {
				key := v.UserId
				delete(members, key)
			}
		}
	}

	var roomName string
	info, err := orm.FindRoomById(msg.GetTargetId(), types.RoomNotDeleted)
	if err == nil && info != nil {
		roomName = info.Name
	}

	title := roomName
	text := msg.GetGTMsg(userId)

	for k, v := range members {
		var pusher IPusher
		if pusher, err = checkDeviceToken(app, k, v); pusher == nil {
			if err != nil {
				//TODO 输出错误日志
			}
			continue
		}
		if cfg.Log.Level == "debug" {
			b, err := json.Marshal(pusher)
			if err != nil {
			}
			pushLog.Debug("push to room", "roomId", roomId, "userId", userId, "pusher", string(b))
		}
		err = pusher.singlePush(title, text, roomId, types.ToRoom)
		if err != nil {
			pushLog.Error(err.Error())
		}
	}
	return nil
}
