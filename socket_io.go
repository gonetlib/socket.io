package socketio

const (
	SIEventOpen        = "open"
	SIEventClose       = "close"
	SIEventEvent       = "event"
	SIEventAck         = "ack"
	SIEventError       = "error"
	SIEventBinaryEvent = "binary event"
	SIEventBinaryAck   = "binary ack"
)

type SocketIOHandler struct {
	conn Connection
}

func NewSocketIOHandler(conn Connection) *SocketIOHandler {
	return &SocketIOHandler{
		conn: conn,
	}
}

func (h *SocketIOHandler) ProcessMessage(payload []byte) {
	switch payload[0] {
	case '0':
		// open connection
		h.conn.Handle(SIEventAck, payload[1:])
	case '1':
		// close connection
		h.conn.Handle(SIEventClose, payload[1:])
	case '2':
		// event
		h.conn.Handle(SIEventEvent, payload[1:])
	case '3':
		// ack
		h.conn.Handle(SIEventAck, payload[1:])
	case '4':
		// connection error
		h.conn.Handle(SIEventError, payload[1:])
	case '5':
		// binary event
		h.conn.Handle(SIEventBinaryEvent, payload[1:])
	case '6':
		// binary ack
		h.conn.Handle(SIEventBinaryAck, payload[1:])
	}
}
