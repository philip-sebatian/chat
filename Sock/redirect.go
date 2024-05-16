package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// this is an upgrader that upgrades the http connection to websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// the message written to the connection is made into a Message struct which has the content and user name
// of the person who client 'c' want to send the message to
// Message struct is then pushed to the broadcaster channel from where the hub send the message
// send channel of the recepient
// the func write to connections continuesly check for any message in send channel
// if there is any message it writes the message in the connection
// from where the recepient can get the message client c has send him (through javascript or with another server)
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
		splited := strings.SplitN(string(message), ":", 2) //format of message "to_username":"messaage
		fmt.Println(splited[1])                            //split the message to who to send the message to
		h.broadcaster <- Message{                          //and the contents of the message
			content:    splited[1],
			to_address: splited[0],
		}
	}

}

// clients lopps throught the messages in the send channel
// and write all the message received by the client in send channel to its connection
// client code can read the text received in this connection
func (c *Client) write_to_a_connection() {
	for messages := range c.send {
		if err := c.con.WriteMessage(websocket.TextMessage, []byte(messages.content)); err != nil {
			fmt.Println("Error in writing message to the websocket connection")
			break
		}
	}
	c.con.Close()
}

// gets http request upgrades it and gets username of the person who sent this requets
// creates a CLient struct and registers with the hub
// starts listening
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
	//register with the hub

	h.register_hub <- client
	//start go routines of the fuction
	go client.write_to_a_connection()
	go client.broadcast_messsages(h)

}
