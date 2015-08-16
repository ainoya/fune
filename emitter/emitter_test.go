package emitter

import (
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/listener"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewEmitter(t *testing.T) {
	c := listener.GetDefaultClient()
	l := listener.NewDockerListener(c)

	NewEmitter(l)
}

func TestBroadCast(t *testing.T) {
	l := listener.NewMockListener()

	e := NewEmitter(l)

	e.actions = make(map[string]actions.Action)

	msg := "message"
	actionNum := 10

	var count uint64

	var wg sync.WaitGroup

	// Prepare mock actions
	for i := 1; i <= actionNum; i++ {
		actionName := fmt.Sprintf("receiver_%d", i)
		a := actions.NewMockAction().(*actions.MockAction)
		a.SetName(actionName)

		e.actions[fmt.Sprintf("mock_%d", i)] = a
	}

	for i := 1; i <= actionNum; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			select {
			case ev, ok := <-e.actions[fmt.Sprintf("mock_%d", j)].Ch():
				if ok {
					assert.Equal(t, ev, msg, "received broadcasted message")
					atomic.AddUint64(&count, 1)
				}
			}
		}(i)
	}

	e.BroadCast()
	e.l.(*listener.MockListener).Ch <- msg
	close(e.l.Stopped())

	wg.Wait()

	expected := strconv.Itoa(actionNum)
	actual := strconv.FormatUint(atomic.LoadUint64(&count), 10)

	assert.Equal(t, expected, actual, "All Actions received message")
}
