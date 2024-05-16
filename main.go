package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//creates new hub
	hub := NewHub()
	//start the hub
	go hub.start()
	//start the http router
	r := mux.NewRouter()
	//any http request as at /ws will be passed to hub.Handlews for an upgrader to websockets
	r.HandleFunc("/ws", hub.Handlews)

	http.Handle("/", r)
	log.Println("Server started at :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
