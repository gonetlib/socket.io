package socketio

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Namespace

	WithNamespace(namespace string) Namespace
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type server struct {
	opts *serverOpts
}

func NewServer(opts *serverOpts) (s Server, err error) {
	if opts == nil {
		opts = NewServerOpts()
	}

	return &server{
		opts: opts,
	}, nil
}

func (s *server) On(event string, handler EventHandler) {
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
		log.Error("parse request:", err)
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
	log.Debug("handle websocket")
	conn, err := NewGorillaWebsocketConnection(w, r)
	if err != nil {
		log.Error("new websocket connection:", err)
		return
	}

	if err = conn.Open(sess); err != nil { // engine.io open protocol
		log.Error("open websocket:", err)
		return
	}
	// if err := s.acceptConnection(conn); err != nil {
	// TODO: handle error event
	// log.Error("accept connection:", err)
	// return
	// }
	log.Info("Websocket connected", log.WithField("conn", conn.RemoteAddr()))
	if err = conn.Serve(); err != nil {
		log.Error("serve websocket:", err)
	}
}

func (s *server) handlePolling(w http.ResponseWriter, r *http.Request, sess *Session) {
	log.Debug("handle polling")
}

func (s *server) parseRequest(r *http.Request) (sess *Session, err error) {
	log.Debug("parse request", r.URL.Query())
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
