package actions

import (
	"github.com/ainoya/fune/listener"
	"reflect"
)

// Action is interface of actions triggered by listener events.
type Action interface {
	Name() string
	On() func(interface{})
	Ch() chan interface{}
	Prepare()
}

// ActionAddress is information of actions
type ActionAddress struct {
	NewFunc     func() Action
	ConfigUnits []*ConfigUnit
}

// ConfigUnit stores config information of actions.
type ConfigUnit struct {
	Name        string
	Description string
	Value       string
}

// SetConfig save config key value pair to `ConfigUnit`.
func (addr *ActionAddress) SetConfig(name, value string) {
	for i := 0; i < len(addr.ConfigUnits); i++ {
		c := addr.ConfigUnits[i]
		if c.Name == name {
			c.Value = value
		}
	}
}

// repository is singleton to manage actions
var repository = &ActionRepository{
	EnabledActions:   make(map[string]Action),
	InstalledActions: make(map[string]*ActionAddress),
}

// ActionRepository binds requirement variables for which manage actions
type ActionRepository struct {
	EnabledActions   map[string]Action
	InstalledActions map[string]*ActionAddress
	Listener         listener.Listener
}

// NewActions returns defined action list as singleton.
func NewActions(listener listener.Listener) *ActionRepository {
	if repository.EnabledActions == nil {
		repository.EnabledActions = make(map[string]Action)
	}

	repository.Listener = listener

	return repository
}

// InstallAction is called from every actions on `init()` functions,
// and register information as action's list.
func InstallAction(name string, action Action, newActionFunc func() Action) {
	repository.InstalledActions[name] = &ActionAddress{
		NewFunc:     newActionFunc,
		ConfigUnits: readConfigKeysFromAction(action),
	}
}

func readConfigKeysFromAction(action Action) []*ConfigUnit {
	var configUnits [](*ConfigUnit)
	s := reflect.ValueOf(action).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		configKey := typeOfT.Field(i).Tag.Get("config")
		description := typeOfT.Field(i).Tag.Get("description")
		if configKey != "" {
			configUnits = append(configUnits, &ConfigUnit{Name: configKey, Description: description})
		}
	}
	return configUnits
}

// EnableActions instantiates enabled actions specified by args.
func EnableActions(targetActionNames []string) {
	for _, actionName := range targetActionNames {
		addr := repository.InstalledActions[actionName]
		// TODO : error handling when target action name was not found in InstalledActions
		if addr.NewFunc != nil {
			RegisterAction(addr)
		}
	}
}

// ClearActions removes all registered `actions`.
func ClearActions() {
	repository.EnabledActions = nil
}

// ApplyConfig applies map values parsed from commandline flags to actions.
func ApplyConfig(actionsConfig map[string]string) {
	for _, action := range repository.EnabledActions {
		s := reflect.ValueOf(action).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			configKey := typeOfT.Field(i).Tag.Get("config")
			configVal := actionsConfig[configKey]
			if configVal != "" {
				f.SetString(configVal)
			}
		}
	}
}

// ReadAllConfigKeys reads all `ConfigUnits` from `InstalledActions` and
// return them as flatten array
func ReadAllConfigKeys(label string) [](*ConfigUnit) {
	var keys [](*ConfigUnit)
	for _, action := range repository.InstalledActions {
		keys = append(keys, action.ConfigUnits...)
	}
	return keys
}

// RegisterAction registers defined action to `actions` list.
func RegisterAction(addr *ActionAddress) {
	a := addr.NewFunc()
	repository.EnabledActions[a.Name()] = a
}

// Actions return action singleton
// TODO : add error handling
func Actions() map[string]Action {
	return repository.EnabledActions
}

// ActivateActions runs all registered actions in `actions` as goroutine.
func ActivateActions() {
	for _, action := range repository.EnabledActions {
		action.Prepare()
		go processOn(action)
	}
}

func processOn(a Action) {
	for {
		select {
		case e := <-a.Ch():
			a.On()(e)
		}
	}
}

// DeactivateActions closes all input channel of registered `actions` list.
func DeactivateActions() {
	actions := repository.EnabledActions
	if actions != nil {
		for _, action := range actions {
			ch := action.Ch()
			if ch != nil {
				close(ch)
			}
		}
	}
}
