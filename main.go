package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	hub := NewHub()
	go hub.start()

	r := mux.NewRouter()
	r.HandleFunc("/ws", hub.Handlews)

	http.Handle("/", r)
	log.Println("Server started at :5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
