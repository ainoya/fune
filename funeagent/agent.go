package funeagent

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/coreos/pkg/capnslog"
)

var plog = capnslog.NewPackageLogger("github.com/ainoya/fune", "funeagent")

// FuneAgent is the agent implementation
type FuneAgent struct {
	cfg *AgentConfig

	stop   chan struct{}
	done   chan struct{}
	errorc chan error
}

// NewAgent creates a new FuneAgent from the supplied configuration. The
// configuration is considered static for the lifetime of the FuneAgent.
func NewAgent(cfg *AgentConfig) (*FuneAgent, error) {

	agent := &FuneAgent{
		cfg:    cfg,
		errorc: make(chan error, 1),
	}

	return agent, nil
}

// Start prepares and starts agent in a new goroutine. It is no longer safe to
// modify a agent's fields after it has been sent to Start.
// It also starts a goroutine to publish its agent information.
func (a *FuneAgent) Start() {
	a.start()
	//TODO :  goroutine to publish its agent infromation
}

// Stop stops the agent gracefully, and shuts down the running goroutine.
// Stop should be called after a Start(s), otherwise it will block forever.
func (a *FuneAgent) Stop() {
	select {
	case a.stop <- struct{}{}:
	case <-a.done:
		return
	}
	<-a.done
}

// start is called by Start() function internally.
func (a *FuneAgent) start() {
	a.done = make(chan struct{})
	a.stop = make(chan struct{})

	go a.run()
}

func (a *FuneAgent) run() {

	defer func() {
		close(a.done)
	}()

	for {
		select {
		case err := <-a.errorc:
			plog.Errorf("%s", err)
			return
		case <-a.stop:
			return
		}
	}
}

// StopNotify returns a channel that receives a empty struct
// when the server is stopped.
func (a *FuneAgent) StopNotify() <-chan struct{} { return a.done }
