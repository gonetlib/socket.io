package socketio

import (
	"errors"

	"github.com/liuliqiang/log4go"
)

type EngineIOHandler struct {
	conn  Connection
	siHdr *SocketIOHandler
}

func NewEngineIOHandler(conn Connection, siHdr *SocketIOHandler) *EngineIOHandler {
	return &EngineIOHandler{
		conn:  conn,
		siHdr: siHdr,
	}
}

func (h *EngineIOHandler) processText(msg []byte) (err error) {
	// engine.io protocol
	if len(msg) == 0 {
		log4go.Error("empty message")
		return errors.New("empty message")
	}
	log4go.Debug("engine.io message: '%s'", msg)
	switch msg[0] {
	case '0': // open
	case '1': // close
	case '2': // ping
	case '3': // pong
	case '4': // message
		h.siHdr.ProcessMessage(msg[1:])
	case '5': // upgrade
	case '6': // noop
	}
	return nil
}
