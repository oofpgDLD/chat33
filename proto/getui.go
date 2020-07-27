package proto

import (
	"fmt"

	"github.com/33cn/chat33/orm"

	"github.com/33cn/chat33/types"
)

// TODO
func (p *Proto) GetGTMsg(fromId string) string {
	if p.msg.isEncryptedMsg() {
		return "[你收到一条消息]"
	}
	if p.chType == types.ToUser {
		if p.isSnap == types.IsSnap {
			return "【阅后即焚】"
		}
		switch p.msgType {
		case types.System:
			textHead := "【公告】"
			textBody := p.msg["content"]
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.Text:
			return p.msg["content"].(string)
		case types.Audio:
			return "【语音】"
		case types.Photo:
			return "【图片】"
		case types.RedPack:
			textHead := "【红包】"
			textBody := p.msg["remark"]
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.Video:
			return "【视频】"
		case types.Alert:
			textHead := "【通知】"
			textBody := p.msg["content"]
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		default:
			return ""
		}
	}
	if p.chType == types.ToRoom {
		name := orm.GetMemberName(p.targetId, fromId)
		if p.isSnap == types.IsSnap {
			textHead := name
			textBody := "【阅后即焚】"
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		}
		switch p.msgType {
		case types.System:
			textHead := "【公告】"
			textBody := p.msg["content"]
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.Text:
			textHead := name
			textBody := p.msg["content"].(string)
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.Audio:
			textHead := name
			textBody := "【语音】"
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.Photo:
			textHead := name
			textBody := "【图片】"
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.RedPack:
			textHead := "【红包】"
			textBody := p.msg["remark"]
			text := fmt.Sprintf("%v: %v%v", name, textHead, textBody)
			return text
		case types.Video:
			textHead := name
			textBody := "【视频】"
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		case types.Alert:
			textHead := "【通知】"
			textBody := p.msg["content"]
			text := fmt.Sprintf("%v: %v", textHead, textBody)
			return text
		default:
			return ""
		}
	}
	return ""
}
