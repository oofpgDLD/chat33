package model

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
	"github.com/inconshreveable/log15"
)

/**
	这个文件是关于消息的消息通知的
**/

var noticeLog = log15.New("logic", "logic/event_notice")

//返回，content 和 用户名称的列表
func ConvertInviteAlertContent(roomId, operater string, users []string) ([]string, string, string) {
	var menbers string
	var isMoreThanSix bool
	names := make([]string, 0)
	operName := orm.GetMemberName(roomId, operater)
	for i, v := range users {
		if i > 5 {
			isMoreThanSix = true
		}
		//more than seven
		if i > 6 {
			break
		}
		name := orm.GetMemberName(roomId, v)
		if !isMoreThanSix {
			if i != 0 {
				menbers += "、"
			}
			menbers += name
		}
		names = append(names, name)
	}
	if isMoreThanSix {
		return names, operName, fmt.Sprintf(" %s 邀请 %s 等%d人 进入群聊", operName, menbers, len(users))
	} else {
		return names, operName, fmt.Sprintf(" %s 邀请 %s 进入群聊", operName, menbers)
	}
}

//组成消息通知格式:
//caller:调用的即协议中的发送者
//target:发送对象类型，3 to user; 2 to room
//targetId:发送对象id
//members: 如果target为群,则改字段为nil；如果target为用户，则members不为空
//messageType:消息通知类型，详见
func SendAlertCreate(caller, targetId string, target int, members []string, messageType int, content map[string]interface{}, msgTime int64) (*proto.Proto, error) {
	//生成msgId：规则 uuid 大写
	msgId := uuid.NewV4().String()
	p, err := proto.NewCmnProto("", msgId, targetId, target, messageType, types.IsNotSnap, content, msgTime, nil)
	if err != nil {
		noticeLog.Error("create cmn proto", "err", err.Error())
		return nil, err
	}
	err = p.AppendChatLog(caller, types.NotRead, msgTime)
	if err != nil {
		noticeLog.Error("append proto log", "err", err.Error())
		return nil, err
	}

	isSpecial := p.IsSpecial()
	//append special
	if target == types.ToRoom {
		if members == nil && isSpecial {
			go func() {
				//AppendReceiveLog(p.GetTargetId(), p.GetLogId(), types.NotRead)
				members, err := orm.FindNotDelMembers(p.GetTargetId(), types.SearchAll)
				if err != nil {
					return
				}
				for _, u := range members {
					_, err := orm.AppendMemberRevLog(u.RoomMember.UserId, p.GetLogId(), types.NotRead)
					if err != nil {
						noticeLog.Warn("appendLog", "warn", "db err")
					}
				}
			}()
		} else if isSpecial {
			go func() {
				//AppendReceiveLogByList(members, p.GetLogId(), types.NotRead)
				for _, u := range members {
					_, err := orm.AppendMemberRevLog(u, p.GetLogId(), types.NotRead)
					if err != nil {
						noticeLog.Warn("appendLog", "warn", "db err")
					}
				}
			}()
		}
	}
	return p, nil
}

func SendAlert(caller, targetId string, target int, members []string, messageType int, content map[string]interface{}) {
	msgTime := utility.NowMillionSecond()
	p, err := SendAlertCreate(caller, targetId, target, members, messageType, content, msgTime)
	if err != nil {
		return
	}
	u, _ := orm.GetUserInfoById(caller)
	if u == nil {
		noticeLog.Error("wrap resp 查询发送者信息失败")
		return
	}
	resp, err2 := p.WrapResp(u, msgTime)
	if err2 != nil {
		noticeLog.Error("wrap resp", "err", err2.Error())
		return
	}

	if target == types.ToRoom {
		if members == nil {
			/*go func() {
				err := PushToRoom(caller, p.GetTargetId(), p)
				if err != nil{
					noticeLog.Warn("push to room failed", "error", err)
				}
			}()*/
			clId := types.GetRoomRouteById(targetId)
			cl, _ := router.GetChannel(clId)
			if cl != nil {
				cl.Broadcast(resp)
			}
		} else {
			var unConnects = make(map[string]string)
			for _, memId := range members {
				if client, ok := router.GetUser(memId); ok && client != nil {
					client.SendToAllClients(resp)
				} else {
					//TODO 判断用户是否存在
					unConnects[memId] = ""
				}
			}
			/*go func() {
				PushToMember(caller, p.GetTargetId(), unConnects, p)
			}()*/
		}
	}

	if target == types.ToUser {
		//text := "收到一条好友消息"
		for _, memId := range members {
			if u, ok := router.GetUser(memId); ok && u != nil {
				u.SendToAllClients(resp)
			}
		}
	}
}
