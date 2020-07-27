package proto

import (
	"encoding/json"
	"unsafe"

	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

type msg map[string]interface{}

func NewMsg(mp map[string]interface{}) msg {
	return mp
}

// 系统消息
func (m msg) checkSysMsg() bool {
	if m.isEncryptedMsg() {
		return true
	}

	if ss, ok := m["content"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	return true
}

// 文本消息
func (m msg) checkTxtMsg() bool {
	if m.isEncryptedMsg() {
		return true
	}

	if ss, ok := m["content"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	return true
}

// 语音消息
func (m msg) checkAudioMsg() bool {
	if m.isEncryptedMsg() {
		return true
	}

	if ss, ok := m["mediaUrl"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	if ss, ok := m["time"]; !ok {
		return false
	} else {
		if utility.ToFloat32(ss) < 0 {
			return false
		}
	}

	return true
}

// 图片消息
func (m msg) checkImgMsg() bool {
	if m.isEncryptedMsg() {
		return true
	}

	if ss, ok := m["imageUrl"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	if ss, ok := m["height"]; !ok {
		return false
	} else {
		if utility.ToInt32(ss) <= 0 {
			return false
		}
	}

	if ss, ok := m["width"]; !ok {
		return false
	} else {
		if utility.ToInt32(ss) <= 0 {
			return false
		}
	}

	return true
}

// 红包消息
func (m msg) checkRedPacketMsg() bool {
	if _, ok := m["coin"]; !ok {
		return false
	}

	if ss, ok := m["packetId"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	if _, ok := m["packetType"]; !ok {
		return false
	}

	return true
}

// 视频消息
func (m msg) checkVideoMsg() bool {
	if m.isEncryptedMsg() {
		return true
	}

	if ss, ok := m["mediaUrl"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	if ss, ok := m["time"]; !ok {
		return false
	} else {
		if utility.ToFloat32(ss) < 0 {
			return false
		}
	}

	return true
}

// 文件消息
func (m msg) checkFileMsg() bool {
	if m.isEncryptedMsg() {
		return true
	}

	if ss, ok := m["fileUrl"]; !ok {
		return false
	} else {
		if utility.ToString(ss) == "" {
			return false
		}
	}

	if ss, ok := m["size"]; !ok {
		return false
	} else {
		if utility.ToFloat32(ss) < 0 {
			return false
		}
	}

	return true
}

// 检查长度
func (m msg) checkMsgLength() bool {
	// TODO
	bytes, err := json.Marshal(m)
	if err != nil {
		return false
	}
	array := []rune(*(*string)(unsafe.Pointer(&bytes)))
	return len(array) <= types.DbMsgLength
}

func (m msg) CheckSpecialMsg() bool {
	var typeCode int
	if v, ok := m["type"]; !ok {
		return false
	} else {
		typeCode = utility.ToInt(v)
	}
	return isSpecialType(typeCode)
}

func (m msg) GetTypeCode() int {
	var typeCode int
	if v, ok := m["type"]; !ok {
		return 0
	} else {
		typeCode = utility.ToInt(v)
	}
	return typeCode
}

//是否是加密消息
func (m msg) isEncryptedMsg() bool {
	if _, ok := m["encryptedMsg"]; ok {
		return true
	}
	return false
}
