package listener

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"
)

// MockListener is mock implementation for tests.
type MockListener struct {
	Ch      chan interface{}
	stopped chan struct{}
}

// StartListen produce dummy messages.
func (l *MockListener) StartListen() {
	go l.produceMockMessages()
}

// Events returns value `l.Ch`
func (l *MockListener) Events() <-chan interface{} {
	return l.Ch
}

// NewMockListener returns instantiated `MockListner`.
func NewMockListener() *MockListener {
	return &MockListener{
		Ch:      make(chan interface{}),
		stopped: make(chan struct{}),
	}
}

// Stopped returns value `l.stopped`
func (l *MockListener) Stopped() chan struct{} {
	return l.stopped
}

func (l *MockListener) produceMockMessages() {
	msgNum := 10

	for i := 1; i <= msgNum; i++ {
		l.Ch <- fmt.Sprintf("mock_message_%d", i)
	}

	close(l.stopped)
}

func (l *MockListener) produceDockerEvents() {
	msgNum := 10

	for i := 1; i <= msgNum; i++ {
		l.Ch <- &docker.APIEvents{
			Status: "create",
			ID:     "dfdf82bd3881",
			From:   "base:latest",
			Time:   1374067970,
		}
	}

	close(l.stopped)
}

// StartProduceDockerEvents produces dummy docker events.
func (l *MockListener) StartProduceDockerEvents() {
	go l.produceDockerEvents()
}
