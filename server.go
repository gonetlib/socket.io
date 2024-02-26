package socketio

import (
	"errors"
	"net/http"
	"sync"

	"github.com/liuliqiang/log4go"
)

type Server interface {
	Namespace

	WithNamespace(namespace string) Namespace
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type server struct {
	opts *serverOpts

	connections sync.Map // session id -> connection
}

func NewServer(opts *serverOpts) (s Server, err error) {
	if opts == nil {
		opts = NewServerOpts()
	}

	return &server{
		opts:        opts,
		connections: sync.Map{},
	}, nil
}

func (s *server) On(event string, handler EventHandler) {
	log4go.Debug("on event: %s", event)
}

func (s *server) WithNamespace(nsp string) Namespace {
	return &namespace{
		nsp:         nsp,
		nspHandlers: map[string]EventHandler{},
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess, err := s.parseRequest(r)
	if err != nil {
		log4go.Error("parse request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	switch sess.Transport {
	case "polling":
		sess.Transport = ""
		s.handlePolling(w, r, sess)
	case "websocket":
		sess.Transport = ""
		s.handleWebsocket(w, r, sess)
	}
}

func (s *server) handleWebsocket(w http.ResponseWriter, r *http.Request, sess *Session) {
	log4go.Debug("handle websocket")
	conn, err := NewGorillaWebsocketConnection(s, w, r)
	if err != nil {
		log4go.Error("new websocket connection:", err)
		return
	}

	if err = conn.Open(sess); err != nil { // engine.io open protocol
		log4go.Error("open websocket:", err)
		return
	}
	s.connections.Store(sess.Sid, conn)

	log4go.Info("Websocket connected :%s", conn.RemoteAddr())
	if err = conn.Serve(); err != nil {
		log4go.Error("serve websocket: %v", err)
	}
}

func (s *server) handlePolling(w http.ResponseWriter, r *http.Request, sess *Session) {
	log4go.Debug("handle polling")
}

func (s *server) parseRequest(r *http.Request) (sess *Session, err error) {
	log4go.Debug("parse request: %+v", r.URL.Query())
	sess = &Session{
		TransportID: r.URL.Query().Get("t"),
		Sid:         r.URL.Query().Get("sid"),
		Upgrades:    []string{"polling"},
	}
	if sess.Sid == "" {
		sess.Sid = "lv_VI97HAXpY6yYWAAAC"
	}

	// if sess.EIOVersion = r.URL.Query().Get("EIO"); sess.EIOVersion == "" {
	// 	return nil, errors.New("missing EIO version")
	// }
	if sess.Transport = r.URL.Query().Get("transport"); sess.Transport == "" {
		return nil, errors.New("missing transport")
	}

	return sess, nil
}

type serverOpts struct {
	// TODO: add tls support
	listenAddr     string
	transports     []string
	pingTimeoutMs  int
	pingIntervalMs int
}

func NewServerOpts() *serverOpts {
	return &serverOpts{
		listenAddr:     "localhost:8000",
		transports:     []string{"polling", "websocket"},
		pingTimeoutMs:  5000,
		pingIntervalMs: 25000,
	}
}

func (o *serverOpts) WithListenAddr(addr string) *serverOpts {
	o.listenAddr = addr
	return o
}
