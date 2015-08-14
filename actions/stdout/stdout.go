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

func (a *actionStdOut) on() {
	for {
		select {
		case event := <-a.Ch():
			fmt.Printf("new event received: %s %s %s %s",
				event.(docker.APIEvents).Status,
				event.(docker.APIEvents).ID,
				event.(docker.APIEvents).From,
				event.(docker.APIEvents).Time,
			)
		}
	}
}
