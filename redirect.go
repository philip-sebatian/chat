package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

func (c *Client) broadcast_messsages(h *Hub) {
	defer func() {
		h.unregister_hub <- c
		c.con.Close()
	}()
	for {
		_, message, err := c.con.ReadMessage()
		if err != nil {
			fmt.Println("error in reading message from connection")
			return
		}
		splited := strings.SplitN(string(message), ":", 2)
		fmt.Println(splited[1])
		h.broadcaster <- Message{
			content:    splited[1],
			to_address: splited[0],
		}
	}

}

func (c *Client) write_to_a_connection() {
	for messages := range c.send {
		if err := c.con.WriteMessage(websocket.TextMessage, []byte(messages.content)); err != nil {
			fmt.Println("Error in writing message to the websocket connection")
			break
		}
	}
	c.con.Close()
}
func (h *Hub) Handlews(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error has occurred while upgrading http to websocket")
		return
	}
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Username not provided")
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}
	client := &Client{
		username: username,
		send:     make(chan Message, 256),
		con:      conn,
	}

	h.register_hub <- client
	//client give message func implement ****

	//client read message implement *******
	go client.write_to_a_connection()
	go client.broadcast_messsages(h)

}
