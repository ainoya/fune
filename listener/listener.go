package listener

// Listener listens event stream, like docker events API.
type Listener interface {
	StartListen()
	Events() <-chan interface{}
	Stopped() chan struct{}
}
