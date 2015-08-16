package actions

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

type StdOutAction struct {
	ch   chan interface{}
	out  chan *docker.APIEvents
	name string
}

// StdOutActionName is used for identify name of itself.
var StdOutActionName = "stdout"

// init function called once at launching program.
// In init() function, InstallAction is called and registers itself
// to `installedAction`
func init() {
	InstallAction(StdOutActionName, &StdOutAction{}, NewStdOutAction)
}

// Name returns value `name` of struct `StdOutAction`.
func (a *StdOutAction) Name() string {
	return a.name
}

// Ch returns value `ch` of struct `StdOutAction`.
func (a *StdOutAction) Ch() chan interface{} {
	return a.ch
}

// On returns functions that outputs to STDOUT.
func (a *StdOutAction) On() func(event interface{}) {

	f := func(e interface{}) {
		a.out <- e.(*docker.APIEvents)
	}

	return f
}

//NewStdOutAction returns instantiated `StdOutAction`.
func NewStdOutAction() Action {
	a := &StdOutAction{
		name: StdOutActionName,
		ch:   make(chan interface{}),
		out:  make(chan *docker.APIEvents),
	}

	return a
}

// Prepare is called when `actions` loads action instance.
// printMsg is called with goroutine to receive docker events
// from channel `StdOutAction.ch` and print out them.
func (a *StdOutAction) Prepare() {
	go a.printMsg()
}

func (a *StdOutAction) printMsg() Action {
	for {
		e := <-a.out
		fmt.Printf("Status: %s ID: %s From: %s Time: %d\n",
			e.Status,
			e.ID,
			e.From,
			e.Time,
		)
	}
}
