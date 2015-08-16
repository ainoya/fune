package flags

import (
	"flag"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/coreos/pkg/capnslog"
	"strings"
)

var (
	plog = capnslog.NewPackageLogger("github.com/ainoya/fune/pkg", "flags")
)

// ActionNames implements the flag.Value interface.
type ActionNames []string

// Set implements the function of flag.Value interface.
// It splits by a comma separeted value to array of string.
func (actionNames *ActionNames) Set(s string) error {
	strs := strings.Split(s, ",")
	*actionNames = strs
	return nil
}

// Set implements the function of flag.Value interface.
// It joins an array of strings to a comma separeted value(string).
func (actionNames *ActionNames) String() string {
	all := make([]string, len(*actionNames))
	for i, u := range *actionNames {
		all[i] = u
	}
	return strings.Join(all, ",")
}

// NewActionNames returns instantiated `ActionNames`
func NewActionNames() *ActionNames {
	v := &ActionNames{}
	return v
}

// EnabledActionsFromFlags returns value as an array of string
// that is assigned to a command line flag name `enabledActionsFlagName`.
func EnabledActionsFromFlags(fs *flag.FlagSet, enabledActionsFlagName string) ([]string, error) {
	visited := make(map[string]struct{})
	fs.Visit(func(f *flag.Flag) {
		visited[f.Name] = struct{}{}
	})

	_, enabledActionsFlagIsSet := visited[enabledActionsFlagName]

	if enabledActionsFlagIsSet {
		enabledActions := *fs.Lookup(enabledActionsFlagName).Value.(*ActionNames)
		return enabledActions, nil
	}

	return []string{}, nil
}
