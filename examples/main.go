package main

import (
	"log"
	"net/http"

	socketio "github.com/gonetlib/socket.io"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On(socketio.EventConnection, func(conn socketio.Connection, _ error) {
		conn.Emit("message", []byte("hello"))
		log.Println("on connection")
	})
	server.On(socketio.EventDisConnection, func(conn socketio.Connection, err error) {
		log.Println("on disconnect")
	})
	server.On(socketio.EventError, func(conn socketio.Connection, err error) {
		log.Println("error:", err)
	})

	server.On("chat", func(ns socketio.Connection, err error) {
	})

	server.WithNamespace("/chat").On("", func(ns socketio.Connection, err error) {

	})

	http.Handle("/", server)

	addr := "localhost:5123"
	log.Printf("Serving at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
