package actions

// MockAction is mock implementation for tests.
type MockAction struct {
	ch          chan interface{}
	name        string
	ConfigValue string `config:"mock-config-value" description:"config key for test"`
	Memory      []interface{}
	Stopped     chan struct{}
}

// MockActionName is used for identify name of itself.
var MockActionName = "mock"

func init() {
	InstallAction("mock", &MockAction{}, NewMockAction)
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

// Prepare is dummy implementation of Action.
func (a *MockAction) Prepare() {
}

// Ch returns value `ch` of struct `MockAction`.
func (a *MockAction) Ch() chan interface{} {
	return a.ch
}

// Name returns value `name` of struct `MockAction`.
func (a *MockAction) Name() string {
	return a.name
}

// SetName setter of `MockAction.name` for test.
func (a *MockAction) SetName(name string) {
	a.name = name
}

//NewMockAction returns instantiated `MockAction`.
func NewMockAction() Action {
	a := &MockAction{
		name:    MockActionName,
		ch:      make(chan interface{}),
		Stopped: make(chan struct{}),
	}
	return a
}
