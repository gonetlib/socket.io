package socketio

type Namespace interface {
	On(event string, handler EventHandler)
}

type EventHandler func(so Connection, err error)

type namespace struct {
	nsp         string
	nspHandlers map[string]EventHandler
}

func (n *namespace) On(event string, handler EventHandler) {
	n.nspHandlers[event] = handler
}
