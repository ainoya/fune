package emitter

import (
	"container/list"
	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/listener"
)

// Emitter listens events from `Listener`, and emits registered `Action`.
type Emitter struct {
	l       listener.Listener
	actions *list.List
}

// NewEmitter returns instantiated `NewEmitter`.
func NewEmitter(l listener.Listener) *Emitter {
	emitter := &Emitter{
		l: l,
	}

	return emitter
}

// LoadActions loads registered `Actions` into `Emitter.actions`
func (e *Emitter) LoadActions(as *list.List) {
	e.actions = as
}

// BroadCast publish received events from `Listener` to all registered actions.
func (e *Emitter) BroadCast() {
	events := e.l.Events()
	go func() {
		for {
			select {
			case event := <-events:
				e.publishEvent(event)
			case <-e.l.Stopped():
				return
			}
		}
	}()
}

// Stopped returns value `Emitter.Listner.Stopped()` for checking whether listner is finished
// to produce events.
func (e *Emitter) Stopped() chan struct{} {
	return e.l.Stopped()
}

func (e *Emitter) publishEvent(event interface{}) {
	for elem := e.actions.Front(); elem != nil; elem = elem.Next() {
		action := elem.Value.(actions.Action)
		go func(a actions.Action) { a.Ch() <- event }(action)
	}
}
