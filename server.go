package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	send     chan Message //message received to be written to this connection
	con      *websocket.Conn
}

type Message struct {
	content    string
	to_address string
}

type Hub struct {
	clients        map[string]*Client
	broadcaster    chan Message
	register_hub   chan *Client
	unregister_hub chan *Client
	mu             sync.RWMutex
}

func (h *Hub) start() {

	for {
		select {
		case client := <-h.register_hub:
			h.mu.Lock()
			h.clients[client.username] = client
			h.mu.Unlock()

		case client := <-h.unregister_hub:
			h.mu.Lock()
			delete(h.clients, client.username)
			h.mu.Unlock()
		case message := <-h.broadcaster:
			h.mu.RLock()
			h.clients[message.to_address].send <- message
			h.mu.RUnlock()
		}

	}

}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		broadcaster:    make(chan Message),
		register_hub:   make(chan *Client),
		unregister_hub: make(chan *Client),
	}
}
