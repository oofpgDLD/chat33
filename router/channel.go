package router

import (
	"sync"
)

var channels *chMap

func init() {
	channels = &chMap{
		data: make(map[string]*Channel),
	}
	channels.add("default", NewChannel("default"))
}

func GetChannel(key string) (*Channel, bool) {
	return channels.get(key)
}

func AppendChannel(key string, value *Channel) {
	channels.add(key, value)
}

func DeleteChannel(key string) {
	channels.delete(key)
}

// chMap maintains all chat and private rooms
type chMap struct {
	sync.RWMutex

	data map[string]*Channel // channelId --> Ch
}

func (cls *chMap) get(key string) (*Channel, bool) {
	cls.RLock()
	defer cls.RUnlock()

	ch, ok := cls.data[key]
	return ch, ok
}

func (cls *chMap) add(key string, value *Channel) {
	cls.Lock()
	defer cls.Unlock()

	cls.data[key] = value
}

func (cls *chMap) delete(key string) {
	cls.Lock()
	defer cls.Unlock()

	delete(cls.data, key)
}

// Channel struct of room or personal chan dec
type Channel struct {
	sync.RWMutex

	route string
	hub   *Hub
}

func NewChannel(route string) *Channel {
	return &Channel{
		route: route,
		hub:   NewHub(),
	}
}

func (cl *Channel) Register(userId string, client *Client) {
	cl.hub.register <- client
}

func (cl *Channel) UnRegister(userId string, client *Client) {
	cl.hub.unregister <- client
}

func (cl *Channel) Broadcast(msg interface{}) {
	cl.hub.broadcast <- msg
}

func (cl *Channel) GetRegisterNumber() int {
	users := make(map[string]bool)
	clients := cl.hub.GetAllClients()
	for _, v := range clients {
		users[v.user.Id] = true
	}
	return len(users)
}
