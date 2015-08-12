package funemain

import (
	"flag"
	"fmt"
	"os"
)

type config struct {
	*flag.FlagSet

	debug        bool
	logPkgLevels string
}

//newConfig initializes config struct
func newConfig() *config {
	cfg := &config{}

	cfg.FlagSet = flag.NewFlagSet("etcd", flag.ContinueOnError)
	fs := cfg.FlagSet
	fs.Usage = func() {
		fmt.Println(usageline)
	}

	// logging
	fs.BoolVar(&cfg.debug, "debug", false, "Enable debug output to the logs.")

	return cfg
}

func (cfg *config) Parse(arguments []string) error {
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

	//TODO : show version

	return nil
}
