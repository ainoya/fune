package emitter

import (
	"container/list"
	"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/listener"
	"sync"
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

	e.actions = list.New()

	msg := "message"
	actionNum := 10

	count := 0

	var wg sync.WaitGroup

	for i := 1; i <= actionNum; i++ {
		go func(j int) {
			a := actions.NewMockAction(fmt.Sprintf("receiver_%d", j))
			e.actions.PushBack(
				a,
			)
			wg.Add(1)
			defer wg.Done()
			select {
			case e, ok := <-a.Ch():
				if ok {
					assert.Equal(t, e, msg, "received broadcasted message")
					count += 1
				}
			}
		}(i)
	}

	e.BroadCast()
	e.l.(*listener.MockListener).Ch <- msg
	close(e.l.Stopped())

	wg.Wait()
	assert.Equal(t, actionNum, count, "All Actions received message")

}
