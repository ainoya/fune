package action

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

type actionStdOut struct {
	ch chan interface{}
}

func (a *actionStdOut) Ch() chan interface{} {
	return a.ch
}

func (a *actionStdOut) On() func(event interface{}) {

	f := func(e interface{}) {
		fmt.Printf("new event received: %s %s %s %s",
			e.(docker.APIEvents).Status,
			e.(docker.APIEvents).ID,
			e.(docker.APIEvents).From,
			e.(docker.APIEvents).Time,
		)
	}

	return f
}
