package flags

import (
	"flag"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestEnabledActionsFromFlags(t *testing.T) {
	args := []string{"-actions=mock"}
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	fs.Var(NewActionNames(), "actions", "")

	fs.Parse(args)

	enabledActions, _ := EnabledActionsFromFlags(fs, "actions")

	assert.Equal(t, enabledActions, []string{"mock"})
}
