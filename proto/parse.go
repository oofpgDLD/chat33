package proto

import (
	"github.com/33cn/chat33/utility"
)

func (p *Proto) parse(m map[string]interface{}) bool {
	if val, ok := m["msgId"]; !ok {
		return false
	} else {
		p.msgId = utility.ToString(val)
	}

	if val, ok := m["eventType"]; !ok {
		return false
	} else {
		evType := utility.ToInt(val)
		if !isCmnEvTypeOk(evType) {
			protoLog.Error("invalid event type")
			return false
		}
		p.evType = evType
	}

	if val, ok := m["channelType"]; !ok {
		return false
	} else {
		chType := utility.ToInt(val)
		if !isChTypeOk(chType) {
			protoLog.Error("invalid channel type")
			return false
		}
		p.chType = chType
	}

	if val, ok := m["targetId"]; !ok {
		return false
	} else {
		p.targetId = utility.ToString(val)
	}

	if val, ok := m["msgType"]; !ok {
		return false
	} else {
		msgType := utility.ToInt(val)
		if !isMsgTypeOk(msgType) {
			return false
		}
		p.msgType = msgType
	}

	if val, ok := m["isSnap"]; !ok {
		return false
	} else {
		isSnap := utility.ToInt(val)
		if !isSnapTypeOk(isSnap) {
			return false
		}
		p.isSnap = isSnap
	}

	if val, ok := m["msg"]; !ok {
		return false
	} else {
		if val2, ok := val.(map[string]interface{}); !ok {
			return false
		} else {
			if !isTypedMsgOk(p.msgType, val2) {
				protoLog.Error("invalid msg")
				return false
			}
			p.msg = val2
		}
	}

	if val, ok := m["ext"]; ok {
		p.ext = val
	}
	return true
}

func (p *EventSyncMsg) parse(m map[string]interface{}) bool {
	if val, ok := m["time"]; !ok {
		return false
	} else {
		p.time = utility.ToInt64(val)
	}
	return true
}

func (p *EventUpdateSKey) parse(m map[string]interface{}) bool {
	if val, ok := m["roomId"]; !ok {
		return false
	} else {
		p.roomId = utility.ToString(val)
	}

	if val, ok := m["fromKey"]; !ok {
		//return false
	} else {
		p.fromKey = utility.ToString(val)
	}

	if val, ok := m["secret"]; !ok {
		return false
	} else {
		if _, ok := val.([]interface{}); !ok {
			return false
		}
		var secrets []secretInfo
		for _, secret := range val.([]interface{}) {
			if _, ok := secret.(map[string]interface{}); !ok {
				return false
			}
			v := secret.(map[string]interface{})
			item := secretInfo{
				userId: utility.ToString(v["userId"]),
				key:    utility.ToString(v["key"]),
			}
			secrets = append(secrets, item)
		}
		p.secret = secrets
	}
	return true
}

func (p *EventGetSKey) parse(m map[string]interface{}) bool {
	if val, ok := m["datetime"]; !ok {
		return false
	} else {
		p.datetime = utility.ToInt64(val)
	}
	return true
}

func (p *EventUpdatePubKey) parse(m map[string]interface{}) bool {
	if val, ok := m["publicKey"]; !ok {
		return false
	} else {
		p.PublicKey = utility.ToString(val)
	}

	if val, ok := m["privateKey"]; !ok {
		return false
	} else {
		p.PrivateKey = utility.ToString(val)
	}

	return true
}

func (p *EventAck) parse(m map[string]interface{}) bool {
	if val, ok := m["begin"]; !ok {
		return false
	} else {
		p.begin = utility.ToInt64(val)
	}

	if val, ok := m["end"]; !ok {
		return false
	} else {
		p.end = utility.ToInt64(val)
	}

	if val, ok := m["total"]; !ok {
		return false
	} else {
		p.total = utility.ToInt64(val)
	}
	return true
}
