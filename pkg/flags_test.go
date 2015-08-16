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

	if err := fs.Parse(args); err != nil {
		t.Errorf("#%d: failed to parse flags: %v", err)
	}

	enabledActions, _ := EnabledActionsFromFlags(fs, "actions")

	assert.Equal(t, enabledActions, []string{"mock"})
}
