package actions

import (
	"container/list"
)

// Action is interface of actions triggered by listner events.
type Action interface {
	On() func(interface{})
	Ch() chan interface{}
}

var actions *list.List

// NewActions returns defined action list as singleton.
func NewActions() *list.List {
	if actions == nil {
		actions = list.New()
	}

	return actions
}

// ClearActions removes all registered `actions`.
func ClearActions() {
	actions = nil
}

// RegisterAction registers defined action to `actions` list.
func RegisterAction(a Action) {
	actions.PushBack(a)
}

// Actions return action singleton
// TODO : add error handling
func Actions() *list.List {
	return actions
}

// ActivateActions runs all registered actions in `actions` as goroutine.
func ActivateActions() {
	for elem := actions.Front(); elem != nil; elem = elem.Next() {
		action := elem.Value.(Action)
		go processOn(action)
	}
}

func processOn(a Action) {
	for {
		select {
		case e := <-a.Ch():
			a.On()(e)
		}
	}
}

// DeactivateActions closes all input channel of registered `actions` list.
func DeactivateActions() {
	if actions != nil {
		for elem := actions.Front(); elem != nil; elem = elem.Next() {
			ch := elem.Value.(Action).Ch()
			if ch != nil {
				close(ch)
			}
		}
	}
}
