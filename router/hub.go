package router

import (
	"sync"

	"github.com/inconshreveable/log15"
)

var rlog = log15.New("module", "chat33/router")

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	sync.RWMutex

	clients    map[*Client]bool // Registered clients
	broadcast  chan interface{} // Inbound messages from the clients
	register   chan *Client     // Register requests from the clients
	unregister chan *Client     // Unregister requests from clients
}

func NewHub() *Hub {
	hub := &Hub{
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go hub.Run()
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.deleteClient(client)
		case message, ok := <-h.broadcast:
			if !ok {
				rlog.Error("hub closed")
				break
			}
			clients := h.GetAllClients()
			for _, client := range clients {
				client.Send(message)
			}
		}
	}
}

func (h *Hub) IsExist(key *Client) bool {
	h.RLock()
	defer h.RUnlock()

	_, ok := h.clients[key]
	return ok
}

func (h *Hub) addClient(c *Client) {
	h.Lock()
	h.clients[c] = true
	h.Unlock()
}

func (h *Hub) deleteClient(c *Client) {
	h.Lock()
	delete(h.clients, c)
	h.Unlock()
}

func (h *Hub) GetAllClients() []*Client {
	h.RLock()
	defer h.RUnlock()

	ret := []*Client{}
	for c := range h.clients {
		ret = append(ret, c)
	}
	return ret
}
