package main

import (
	"log"
	"net/http"

	socketio "github.com/gonetlib/socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On(socketio.EventConnection, func(so socketio.Socket, _ error) {
		log.Println("on connection")
		so.Join("chat")

	})
	server.On(socketio.EventDisConnection, func(so socketio.Socket, err error) {
		log.Println("on disconnect")
	})
	server.On(socketio.EventError, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/", server)

	addr := "localhost:8000"
	log.Printf("Serving at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
