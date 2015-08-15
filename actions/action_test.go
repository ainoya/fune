package actions

import (
	"container/list"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestNewActions(t *testing.T) {
	a := NewActions()

	assert.IsType(t, a, list.New())
	assert.Equal(t, a, Actions())
	assert.Equal(t, a, actions)

	ClearActions()

	assert.Nil(t, actions)
	assert.Nil(t, Actions())
}

func TestRegisterAction(t *testing.T) {
	NewActions()
	a := NewMockAction("test")
	RegisterAction(a)

	registered := Actions().Front().Value.(Action)
	assert.IsType(t, a, registered)
}
