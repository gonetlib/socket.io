package socketio

type Socket interface {
	Join(room string) error
	Emit(event string, msg []byte) error
}

type Broadcast interface {
	Emit(event string, msg []byte) error
}

type Namespace interface {
	Emit(event string, msg []byte) error
	To(room string) Namespace
	In(room string) Namespace
}
