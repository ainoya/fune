package funeagent

// AgentConfig holds the configuration of fune as taken from the command line.
type AgentConfig struct {
	Name           string
	EnabledActions []string
	ActionsConfig  map[string]string
	ListenerType   string
}
