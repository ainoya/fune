package actions

import (
	"container/list"
)

type Action interface {
	On()
	Ch() chan interface{}
}

var actions *list.List

func NewActions() *list.List {
	if actions == nil {
		actions = list.New()
	}

	return actions
}

func RegisterAction(a Action) {
	actions.PushBack(a)
}

func Actions() *list.List {
	return actions
}

func ActivateActions() {
	for elem := actions.Front(); elem != nil; elem = elem.Next() {
		action := elem.Value.(Action)
		go action.On()
	}
}

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
