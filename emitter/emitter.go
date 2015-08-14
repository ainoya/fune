package emitter

import (
	"container/list"
	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/listener"
)

type Emitter struct {
	l       listener.Listener
	actions *list.List
}

func NewEmitter(l listener.Listener) *Emitter {
	emitter := &Emitter{
		l: l,
	}

	return emitter
}

func (e *Emitter) LoadActions(as *list.List) {
	e.actions = as
}

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

func (e *Emitter) Stopped() chan struct{} {
	return e.l.Stopped()
}

func (e *Emitter) publishEvent(event interface{}) {
	for elem := e.actions.Front(); elem != nil; elem = elem.Next() {
		action := elem.Value.(actions.Action)
		go func(a actions.Action) { a.Ch() <- event }(action)
	}
}
