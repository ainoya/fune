package integration

import (
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
	as := actions.NewActions()

	// Setup actions
	a := actions.NewMockAction("system_test")
	actions.RegisterAction(a)
	actions.ActivateActions()

	// Load actions to Emitter
	e.LoadActions(as)

	// BroadCast message
	e.BroadCast()

	<-e.Stopped()
	<-a.Stopped

	assert.Len(t, a.Memory, 10, "produced message size")
}
