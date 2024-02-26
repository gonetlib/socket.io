package socketio

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liuliqiang/log4go"
)

type Connection interface {
	net.Conn

	On(event string, handler EventHandler)
	Emit(event string, msg []byte) error

	Close() (err error)
	Serve() (err error)

	Handle(event string, msg []byte) (err error)

	Open(sess *Session) (err error)
}

type gorillaWebsocketConnection struct {
	*websocket.Conn
	server *server

	eiHdr         *EngineIOHandler
	eventHandlers map[string]EventHandler
}

func NewGorillaWebsocketConnection(s *server, w http.ResponseWriter, r *http.Request) (Connection, error) {
	upgrader := &websocket.Upgrader{} // TODO: add more options

	log4go.Debug("headers: %+v", r.Header)
	log4go.Debug("url: %+v", r.URL.Query())
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("upgrade websocket: %w", err)
	}

	log4go.Info("Websocket connecting")

	conn := &gorillaWebsocketConnection{
		server:        s,
		Conn:          c,
		eventHandlers: make(map[string]EventHandler),
		// siHandler: &SocketIOHandler{},
	}
	conn.eiHdr = NewEngineIOHandler(conn, NewSocketIOHandler(conn))

	return conn, nil
}

// func (c *gorillaWebsocketConnection) accept() error {
// 	ft, reader, err := c.NextReader()
// 	if err != nil {
// 		return fmt.Errorf("get websocket reader: %w", err)
// 	}

// 	reader.Read(make([]byte, 0))

// 	return nil
// }

func (c *gorillaWebsocketConnection) On(event string, handler EventHandler) {
	log4go.Debug("on event: %s", event)
	c.eventHandlers[event] = handler
}
func (c *gorillaWebsocketConnection) Handle(event string, msg []byte) (err error) {
	log4go.Debug("handle event: %s", event)
	h, ok := c.eventHandlers[event]
	if !ok {
		log4go.Debug("event not found")
		return
	}

	if err = h(c, msg); err != nil {
		return err
	}
	return nil
}

func (c *gorillaWebsocketConnection) Emit(event string, msg []byte) error {
	return c.WriteMessage(websocket.TextMessage, msg)
}

func (c *gorillaWebsocketConnection) Serve() (err error) {
	for {
		// c.SetReadDeadline(time.Now().Add(3 * time.Second))
		msgType, bytes, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("get websocket reader: %w", err)
		}

		if msgType != websocket.TextMessage {
			log4go.Warn("unsupported message type")
			continue
		}

		c.eiHdr.processText(bytes)
	}
}

func (c *gorillaWebsocketConnection) Read(p []byte) (n int, err error) {
	_, reader, err := c.NextReader()
	if err != nil {
		return 0, fmt.Errorf("get websocket reader: %w", err)
	}

	return reader.Read(p)
}

func (c *gorillaWebsocketConnection) Write(p []byte) (n int, err error) {
	return len(p), c.WriteMessage(websocket.TextMessage, p)
}

func (c *gorillaWebsocketConnection) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

//https://github.com/socketio/engine.io-protocol/tree/v3#websocket

type Session struct {
	// {"sid":"lv_VI97HAXpY6yYWAAAC","upgrades":["websocket"],"pingInterval":25000,"pingTimeout":5000}
	Transport    string   `json:"transport,omitempty"`
	TransportID  string   `json:"t,omitempty"`
	EIOVersion   string   `json:"EIO,omitempty"`
	Sid          string   `json:"sid,omitempty"`
	Upgrades     []string `json:"upgrades,omitempty"`
	PingInterval int      `json:"pingInterval,omitempty"`
	PingTimeout  int      `json:"pingTimeout,omitempty"`
}

func (c *gorillaWebsocketConnection) Open(sess *Session) (err error) {
	sess.PingInterval = 25000 // TODO: come from configure
	sess.PingTimeout = 5000   // TODO: come from configure
	bytes, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	length := len(bytes) + 1
	lengthStr := strconv.Itoa(length)

	resp := make([]byte, 0, len(lengthStr)+2)
	// resp = append(resp, lengthStr...)
	resp = append(resp, '0')
	resp = append(resp, bytes...)
	resp = append(resp, '\n')
	log4go.Debug(string(resp))
	w, err := c.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return fmt.Errorf("get websocket writer: %w", err)
	}

	if _, err = w.Write(resp); err != nil {
		return fmt.Errorf("write session: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("close session: %w", err)
	}
	return nil
}
