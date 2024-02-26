package main

import (
	"log"
	"net/http"

	socketio "github.com/gonetlib/socket.io"
	"github.com/liuliqiang/log4go"
)

func main() {
	log4go.SetLevel(log4go.LogLevelDebug)
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On(socketio.EventConnection, func(conn socketio.Connection, msg []byte) (err error) {
		conn.Emit("message", []byte("hello"))
		log.Println("on connection")
		return nil
	})
	server.On(socketio.EventDisConnection, func(conn socketio.Connection, msg []byte) (err error) {
		log.Println("on disconnect")
		return nil
	})
	server.On(socketio.EventError, func(conn socketio.Connection, msg []byte) (err error) {
		log.Println("error:", err)
		return nil
	})

	server.On("chat", func(ns socketio.Connection, msg []byte) (err error) {
		return nil
	})

	server.WithNamespace("/chat").On("", func(ns socketio.Connection, msg []byte) (err error) {
		return nil
	})

	http.Handle("/", server)

	addr := "localhost:5123"
	log.Printf("Serving at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
