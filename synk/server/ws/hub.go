package ws

import (
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte  //广播事件
	register   chan *Client //监听事件
	unregister chan *Client //取消监听
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

var once sync.Once
var singleton *Hub

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: //当有人注册后
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
