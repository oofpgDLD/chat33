package router

import (
	"sync"
)

const (
	VISITOR   = 0
	NOMALUSER = 1
	MANAGER   = 2
)

var users *userMap

func init() {
	users = &userMap{
		data: make(map[string]*User),
	}
}

func GetUser(key string) (*User, bool) {
	return users.get(key)
}

func AddUser(u *User) {
	if u == nil {
		return
	}
	users.add(u)
}

func GetUserNums() int {
	return users.len()
}

func GetUserClientNums() int {
	return users.userClientNums()
}

// userMap maintains all online users
type userMap struct {
	sync.RWMutex

	//key: uuid
	data map[string]*User
}

func (um *userMap) get(key string) (*User, bool) {
	um.RLock()
	defer um.RUnlock()

	u, ok := um.data[key]
	return u, ok
}

func (um *userMap) add(u *User) {
	um.Lock()
	defer um.Unlock()

	um.data[u.Id] = u
}

func (um *userMap) delete(u *User) {
	um.Lock()
	defer um.Unlock()

	delete(um.data, u.Id)
}

func (um *userMap) len() int {
	um.RLock()
	defer um.RUnlock()

	return len(um.data)
}

func (um *userMap) userClientNums() int {
	um.RLock()
	defer um.RUnlock()
	var n int
	for _, v := range users.data {
		n += len(v.clients)
	}
	return n
}

// User struct
type User struct {
	sync.RWMutex

	Id    string
	Level int

	//key: uuid, value: client
	clients map[string]*Client
}

func NewUser(userId string, level int) *User {
	return &User{
		Id:      userId,
		Level:   level,
		clients: make(map[string]*Client),
	}
}

func (u *User) GetClients() map[string]*Client {
	u.RLock()
	defer u.RUnlock()

	ret := make(map[string]*Client)
	for k, v := range u.clients {
		ret[k] = v
	}
	return ret
}

func (u *User) GetClient(clientId string) (*Client, bool) {
	u.RLock()
	defer u.RUnlock()

	cli, ok := u.clients[clientId]
	return cli, ok
}

func (u *User) AppendClient(c *Client) {
	u.Lock()
	defer u.Unlock()

	u.clients[c.GetId()] = c
}

func (u *User) DeleteClient(clientId string) {
	u.Lock()
	delete(u.clients, clientId)
	/*if len(u.clients) == 0 {
		users.delete(u)
	}*/
	u.Unlock()
}

// all client regist current channel
func (u *User) Subscribe(channel *Channel) {
	clients := u.GetClients()
	for _, cli := range clients {
		channel.Register(u.Id, cli)
	}
}

func (u *User) UnSubscribe(channel *Channel) {
	clients := u.GetClients()
	for _, cli := range clients {
		channel.UnRegister(u.Id, cli)
	}
}

func (u *User) ClientUnRegister(clientId string, channel *Channel) {
	if cli, ok := u.GetClient(clientId); ok {
		channel.UnRegister(u.Id, cli)
	}
}

func (u *User) SendToAllClients(msg interface{}) {
	clients := u.GetClients()
	for _, c := range clients {
		c.Send(msg)
	}
}

func IsUserOnline(userId string) bool {
	_, ok := users.get(userId)
	return ok
}
