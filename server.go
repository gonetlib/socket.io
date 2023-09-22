package socketio

import "net/http"

type serverOpts struct{}

type Server interface {
	On(event string, processor func(so Socket, err error))
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewServer(opts *serverOpts) (server Server, err error) {
	return
}
