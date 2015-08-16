package funemain

import (
	"flag"
	"fmt"
	"github.com/ainoya/fune/actions"
	"github.com/ainoya/fune/pkg"
	"os"
)

// Config stores configurations parsed from command line flags.
type Config struct {
	*flag.FlagSet

	debug        bool
	logPkgLevels string

	enabledActions flags.ActionNames
	actionsConfig  map[string]string

	listenerType string
}

// NewConfig initializes config struct
func NewConfig() *Config {
	cfg := &Config{
		actionsConfig: make(map[string]string),
	}

	cfg.FlagSet = flag.NewFlagSet("fune", flag.ContinueOnError)
	fs := cfg.FlagSet
	fs.Usage = func() {
		fmt.Println(usageline)
	}

	// logging
	fs.BoolVar(&cfg.debug, "debug", false, "Enable debug output to the logs.")

	fs.Var(flags.NewActionNames(), "actions", "List of Action names you want to use")

	for _, actionConfigUnit := range actions.ReadAllConfigKeys("config") {
		var dummy string
		fs.StringVar(
			&dummy,
			actionConfigUnit.Name, "",
			actionConfigUnit.Description,
		)
	}

	return cfg
}

// Parse parses command line flags to FuneConfig.
func (cfg *Config) Parse(arguments []string) error {
	perr := cfg.FlagSet.Parse(arguments)
	switch perr {
	case nil:
	case flag.ErrHelp:
		fmt.Println(flagsline)
		os.Exit(0)
	default:
		os.Exit(2)
	}
	if len(cfg.FlagSet.Args()) != 0 {
		return fmt.Errorf("'%s' is not a valid flag", cfg.FlagSet.Arg(0))
	}

	cfg.enabledActions, _ = flags.EnabledActionsFromFlags(cfg.FlagSet, "actions")

	for _, actionConfigUnit := range actions.ReadAllConfigKeys("config") {
		actionFlagName := actionConfigUnit.Name
		cfg.actionsConfig[actionFlagName] = cfg.FlagSet.Lookup(actionFlagName).Value.String()
	}

	//TODO : show version

	return nil
}
