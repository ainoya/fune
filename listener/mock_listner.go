package listener

import (
	"fmt"
)

type MockListener struct {
	Ch      chan interface{}
	stopped chan struct{}
}

func (l *MockListener) StartListen() {
	go l.produceMockMessages()
}

func (l *MockListener) Events() <-chan interface{} {
	return l.Ch
}

func NewMockListener() *MockListener {
	return &MockListener{
		Ch:      make(chan interface{}),
		stopped: make(chan struct{}),
	}
}

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
