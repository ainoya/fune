package actions

import (
	"github.com/ainoya/fune/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

func TestNewActions(t *testing.T) {
	a := NewActions(nil)

	assert.IsType(t, a, &ActionRepository{})
	assert.Equal(t, a.EnabledActions, Actions())
	assert.Equal(t, a, repository)

	ClearActions()

	assert.Nil(t, repository.EnabledActions)
	assert.Nil(t, Actions())
}

func TestRegisterAction(t *testing.T) {
	NewActions(nil)

	addr := &ActionAddress{
		NewFunc: NewMockAction,
	}

	RegisterAction(addr)

	registered := repository.EnabledActions["mock"]
	assert.IsType(t, registered, NewMockAction())

	ClearActions()
}

func TestInstallAction(t *testing.T) {
	NewActions(nil)
	actionName := "mock"

	InstallAction(actionName, &MockAction{}, NewMockAction)

	a := repository.InstalledActions[actionName].NewFunc()

	assert.IsType(t, a, NewMockAction())
	ClearActions()
}

func TestApplyConfig(t *testing.T) {
	NewActions(nil)
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

	registered := repository.EnabledActions[actionName]

	configured := registered.(*MockAction).ConfigValue

	assert.Equal(t, configured, configValue)

	ClearActions()
}

func TestEnableActions(t *testing.T) {
	NewActions(nil)
	actionName := "mock"

	InstallAction(actionName, &MockAction{}, NewMockAction)
	EnableActions([]string{actionName})

	enabled := repository.EnabledActions[actionName]

	assert.IsType(t, enabled, NewMockAction())
	ClearActions()
}
