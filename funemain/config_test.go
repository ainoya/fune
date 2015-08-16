package funemain

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/ainoya/fune/pkg"
	"testing"
)

func TestParseActionsConfig(t *testing.T) {
	config := NewConfig()

	args := []string{"-mock-config-value=testvalue"}

	config.Parse(args)

	assert.Equal(t, config.actionsConfig["mock-config-value"], "testvalue")
}

func TestParseEnabledActions(t *testing.T) {
	config := NewConfig()

	args := []string{"--actions=a,b,c"}

	config.Parse(args)

	assert.Equal(t, config.enabledActions[0], "a")
	assert.Equal(t, config.enabledActions, flags.ActionNames{"a", "b", "c"})
}
