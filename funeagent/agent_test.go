package funeagent

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"

	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/emitter"
	"github.com/ainoya/fune/listener"
)

func TestNewAgent(t *testing.T) {
	cfg := &AgentConfig{}
	agent, _ := NewAgent(cfg)

	assert.NotNil(t, agent)
}

func TestStart(t *testing.T) {
	cfg := &AgentConfig{}
	agent, _ := NewAgent(cfg)
	cfg.EnabledActions = []string{"mock"}

	agent.listener = listener.NewMockListener()
	agent.emitter = emitter.NewEmitter(agent.listener)

	agent.Start()

	<-agent.emitter.Stopped()
	<-actions.Actions()["mock"].(*actions.MockAction).Stopped

	assert.Len(t, actions.Actions()["mock"].(*actions.MockAction).Memory, 10, "produced message size")

	actions.ClearActions()
}

func TestStopNotify(t *testing.T) {
	s := &FuneAgent{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		<-s.stop
		close(s.done)
	}()

	notifier := s.StopNotify()
	select {
	case <-notifier:
		t.Fatalf("received unexpected stop notification")
	default:
	}
	s.Stop()
	select {
	case <-notifier:
	default:
		t.Fatalf("cannot receive stop notification")
	}
}
