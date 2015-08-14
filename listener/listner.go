package listener

type Listener interface {
	StartListen()
	Events() <-chan interface{}
	Stopped() chan struct{}
}
