package integration

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/emitter"
	"github.com/ainoya/fune/listener"
	"testing"
)

func TestSystemWithAllMock(t *testing.T) {
	l := listener.NewMockListener()
	l.StartListen()
	e := emitter.NewEmitter(l)
	as := actions.NewActions(nil)

	// Setup actions
	addr := &actions.ActionAddress{
		NewFunc: actions.NewMockAction,
	}
	actions.RegisterAction(addr)
	actions.ActivateActions()

	// Load actions to Emitter
	e.LoadActions(as.EnabledActions)

	// BroadCast message
	e.BroadCast()

	<-e.Stopped()
	<-actions.Actions()["mock"].(*actions.MockAction).Stopped

	assert.Len(t, actions.Actions()["mock"].(*actions.MockAction).Memory, 10, "produced message size")

	actions.ClearActions()
}

func TestSystemWithDockerEvents(t *testing.T) {
	l := listener.NewMockListener()
	l.StartProduceDockerEvents()
	e := emitter.NewEmitter(l)
	as := actions.NewActions(nil).EnabledActions

	// Setup actions
	addr := &actions.ActionAddress{
		NewFunc: actions.NewMockAction,
	}
	actions.RegisterAction(addr)
	actions.ActivateActions()

	// Load actions to Emitter
	e.LoadActions(as)

	// BroadCast message
	e.BroadCast()

	<-e.Stopped()
	<-actions.Actions()["mock"].(*actions.MockAction).Stopped

	assert.Len(t, actions.Actions()["mock"].(*actions.MockAction).Memory, 10, "produced message size")

	event := actions.Actions()["mock"].(*actions.MockAction).Memory[0].(*docker.APIEvents)

	assert.Equal(t, event.From, "base:latest")

	actions.ClearActions()
}
