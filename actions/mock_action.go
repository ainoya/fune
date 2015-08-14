package actions

type MockAction struct {
	ch      chan interface{}
	name    string
	Memory  []string
	Stopped chan struct{}
}

func (a *MockAction) On() {
	for {
		select {
		case e := <-a.ch:
			a.Memory = append(a.Memory, e.(string))
			if len(a.Memory) >= 10 {
				close(a.Stopped)
			}
		}
	}
}

func (a *MockAction) Ch() chan interface{} {
	return a.ch
}

func NewMockAction(name string) *MockAction {

	a := &MockAction{
		name:    name,
		ch:      make(chan interface{}),
		Stopped: make(chan struct{}),
	}
	return a
}
