package actions

import (
	"reflect"
)

// Action is interface of actions triggered by listener events.
type Action interface {
	Name() string
	On() func(interface{})
	Ch() chan interface{}
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

var (
	actions          map[string]Action
	installedActions = make(map[string]*ActionAddress)
)

// NewActions returns defined action list as singleton.
func NewActions() map[string]Action {
	if actions == nil {
		actions = make(map[string]Action)
	}

	return actions
}

// InstallAction is called from every actions on `init()` functions,
// and register information as action's list.
func InstallAction(name string, action Action, newActionFunc func() Action) {
	installedActions[name] = &ActionAddress{
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
		addr := installedActions[actionName]
		// TODO : error handling when target action name was not found in installedActions
		if addr.NewFunc != nil {
			RegisterAction(addr)
		}
	}
}

// ClearActions removes all registered `actions`.
func ClearActions() {
	actions = nil
}

// ApplyConfig applies map values parsed from commandline flags to actions.
func ApplyConfig(actionsConfig map[string]string) {
	for _, action := range actions {
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

// ReadAllConfigKeys reads all `ConfigUnits` from `installedActions` and
// return them as flatten array
func ReadAllConfigKeys(label string) [](*ConfigUnit) {
	var keys [](*ConfigUnit)
	for _, action := range installedActions {
		keys = append(keys, action.ConfigUnits...)
	}
	return keys
}

// RegisterAction registers defined action to `actions` list.
func RegisterAction(addr *ActionAddress) {
	a := addr.NewFunc()
	actions[a.Name()] = a
}

// Actions return action singleton
// TODO : add error handling
func Actions() map[string]Action {
	return actions
}

// ActivateActions runs all registered actions in `actions` as goroutine.
func ActivateActions() {
	for _, action := range actions {
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
	if actions != nil {
		for _, action := range actions {
			ch := action.Ch()
			if ch != nil {
				close(ch)
			}
		}
	}
}
