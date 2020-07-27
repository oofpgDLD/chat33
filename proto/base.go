package proto

import (
	"encoding/json"
	"fmt"

	"github.com/inconshreveable/log15"

	"github.com/33cn/chat33/router"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

type ProtocolCreator func() IBaseProto

var proto_creators map[int]ProtocolCreator = make(map[int]ProtocolCreator)

type ProtocolReceiver func(IBaseProto, router.Connection, []byte)

var proto_receiver map[int]ProtocolReceiver = make(map[int]ProtocolReceiver)

type IBaseProto interface {
	Receive(router.Connection, []byte)
}

func init() {
	proto_creators[types.EventCommonMsg] = func() IBaseProto { return new(Proto) }
	proto_creators[types.EventSyncMsg] = func() IBaseProto { return new(EventSyncMsg) }
	proto_creators[types.EventUpdateSKey] = func() IBaseProto { return new(EventUpdateSKey) }
	proto_creators[types.EventUpdatePublicKey] = func() IBaseProto { return new(EventUpdatePubKey) }
	proto_creators[types.EventStartAck] = func() IBaseProto { return new(EventAck) }
	proto_creators[types.EventStartGetAllSKey] = func() IBaseProto { return new(EventGetSKey) }
}

func GetProtocolReceiver() map[int]ProtocolReceiver {
	return proto_receiver
}

func CreateProto(b []byte) (IBaseProto, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, fmt.Errorf("prase err")
	}
	//调试收到的消息
	log15.Debug("receive msg:", "content", m)
	var eventType int
	if val, ok := m["eventType"]; !ok {
		return nil, fmt.Errorf("eventType err")
	} else {
		eventType = utility.ToInt(val)
	}

	creator := proto_creators[eventType]
	if creator == nil {
		return nil, fmt.Errorf("proto not registe")
	}
	return creator(), nil
}
