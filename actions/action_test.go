package actions

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestNewActions(t *testing.T) {
	a := NewActions()

	assert.IsType(t, a, make(map[string]Action))
	assert.Equal(t, a, Actions())
	assert.Equal(t, a, actions)

	ClearActions()

	assert.Nil(t, actions)
	assert.Nil(t, Actions())
}

func TestRegisterAction(t *testing.T) {
	NewActions()

	addr := &ActionAddress{
		NewFunc: NewMockAction,
	}

	RegisterAction(addr)

	registered := actions["mock"]
	assert.IsType(t, registered, NewMockAction())

	ClearActions()
}

func TestInstallAction(t *testing.T) {
	NewActions()
	actionName := "mock"

	InstallAction(actionName, &MockAction{}, NewMockAction)

	a := installedActions[actionName].NewFunc()

	assert.IsType(t, a, NewMockAction())
	ClearActions()
}

func TestApplyConfig(t *testing.T) {
	NewActions()
	actionName := "mock"

	configKey := "mock-config-value"
	configValue := "value"

	addr := &ActionAddress{
		NewFunc: NewMockAction,
		ConfigUnits: []*ConfigUnit{
			&ConfigUnit{
				Name: configKey,
			},
		},
	}

	addr.SetConfig(configKey, configValue)

	RegisterAction(addr)

	cfg := make(map[string]string)
	cfg[configKey] = configValue

	ApplyConfig(cfg)

	registered := actions[actionName]

	configured := registered.(*MockAction).ConfigValue

	assert.Equal(t, configured, configValue)

	ClearActions()
}

func TestEnableActions(t *testing.T) {
	NewActions()
	actionName := "mock"

	InstallAction(actionName, &MockAction{}, NewMockAction)
	EnableActions([]string{actionName})

	enabled := actions[actionName]

	assert.IsType(t, enabled, NewMockAction())
	ClearActions()
}
