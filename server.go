package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	username string
	send     chan Message //message received to be written to this connection ...better rename it to received
	con      *websocket.Conn
}

type Message struct {
	content    string
	to_address string
}

// structure of a hub
type Hub struct {
	clients        map[string]*Client //maps username to Client pointers
	broadcaster    chan Message       // channel that stores the Mesages that needs to be redirected
	register_hub   chan *Client       // channedl for storing the users that wants to register to the hub
	unregister_hub chan *Client       // channel that is used to unregistrer users from hub
	mu             sync.RWMutex       //mutex locks for accessing the clients map
}

func (h *Hub) start() {

	for {
		select {
		// if there is client is the register channel
		//add them to the client set
		case client := <-h.register_hub:
			h.mu.Lock()
			h.clients[client.username] = client
			h.mu.Unlock()
		// if there is client is the unregister channel
		//delete them from the client set
		case client := <-h.unregister_hub:
			h.mu.Lock()
			delete(h.clients, client.username)
			h.mu.Unlock()
		//if there is any messages to be broadcaster to any other person
		//ie if any Message in broadcaster must be send be put to
		//the clients send channnel so that they can read it
		case message := <-h.broadcaster:
			h.mu.RLock()
			h.clients[message.to_address].send <- message
			h.mu.RUnlock()
		}

	}

}

// returns a new hub
func NewHub() *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		broadcaster:    make(chan Message),
		register_hub:   make(chan *Client),
		unregister_hub: make(chan *Client),
	}
}
