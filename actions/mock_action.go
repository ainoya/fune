package actions

// MockAction is mock implementation for tests.
type MockAction struct {
	ch      chan interface{}
	name    string
	Memory  []interface{}
	Stopped chan struct{}
}

// On returns functions which produces dummy messages.
// After producing them, close input channel.
// these dummy messages is saved to `MockAction` instance for testing later.
func (a *MockAction) On() func(event interface{}) {

	f := func(e interface{}) {
		a.Memory = append(a.Memory, e)
		if len(a.Memory) >= 10 {
			close(a.Stopped)
		}
	}

	return f
}

// Ch returns value `ch` of struct `MockAction`.
func (a *MockAction) Ch() chan interface{} {
	return a.ch
}

//NewMockAction returns instantiated `MockAction`.
func NewMockAction(name string) *MockAction {

	a := &MockAction{
		name:    name,
		ch:      make(chan interface{}),
		Stopped: make(chan struct{}),
	}
	return a
}
