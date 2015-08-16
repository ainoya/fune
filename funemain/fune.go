package funemain

import (
	//"fmt"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/coreos/etcd/pkg/osutil"
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/coreos/pkg/capnslog"
	"github.com/ainoya/fune/funeagent"
	"os"
)

var plog = capnslog.NewPackageLogger("github.com/ainoya/fune", "funemain")

// Main sets up fune agent
func Main() {
	cfg := NewConfig()
	err := cfg.Parse(os.Args[1:])
	if err != nil {
		plog.Errorf("error verifying flags, %v. See 'fune --help'.", err)
		os.Exit(1)
	}
	setupLogging(cfg)

	var stopped <-chan struct{}

	stopped, err = startFune(cfg)

	osutil.HandleInterrupts()

	//TODO : systemd integration

	<-stopped
	osutil.Exit(0)
}

// startFune launches the fune agent.
func startFune(cfg *Config) (<-chan struct{}, error) {
	agentCfg := &funeagent.AgentConfig{
		EnabledActions: cfg.enabledActions,
		ActionsConfig:  cfg.actionsConfig,
		ListenerType:   cfg.listenerType,
	}

	var a *funeagent.FuneAgent
	a, err := funeagent.NewAgent(agentCfg)

	if err != nil {
		return nil, err
	}
	a.Start()
	osutil.RegisterInterruptHandler(a.Stop)

	return a.StopNotify(), nil
}

func setupLogging(cfg *Config) {
	capnslog.SetGlobalLogLevel(capnslog.INFO)
	if cfg.debug {
		capnslog.SetGlobalLogLevel(capnslog.DEBUG)
	}
	if cfg.logPkgLevels != "" {
		repoLog := capnslog.MustRepoLogger("github.com/ainoya/fune")
		settings, err := repoLog.ParseLogLevelConfig(cfg.logPkgLevels)
		if err != nil {
			plog.Warningf("couldn't parse log level string: %s, continuing with default levels", err.Error())
			return
		}
		repoLog.SetLogLevel(settings)
	}
}
