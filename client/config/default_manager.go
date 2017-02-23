package config

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/config/storage"
)

type DefaultManager struct {
	storage storage.Storage
	path    string
}

func NewDefaultManager(storage storage.Storage, path string) DefaultManager {
	return DefaultManager{
		storage: storage,
		path:    path,
	}
}

func (m DefaultManager) GetSource() string {
	return m.path
}

func (m DefaultManager) GetEnvironments() (EnvironmentsState, error) {
	state := State{}

	err := m.storage.Get(m.path, &state)
	if err != nil {
		return EnvironmentsState{}, bosherr.WrapError(err, "Getting environments")
	}

	return state.Environments, nil
}

func (m DefaultManager) GetEnvironment(name string) (EnvironmentState, error) {
	envs, err := m.GetEnvironments()
	if err != nil {
		return EnvironmentState{}, bosherr.WrapError(err, "Getting environment")
	}

	for _, env := range envs {
		if env.Alias == name || env.URL == name {
			return env, nil
		}
	}

	return EnvironmentState{}, fmt.Errorf("Environment not found: %s", name)
}

func (m DefaultManager) SetEnvironment(env EnvironmentState) error {
	envs, err := m.GetEnvironments()
	if err != nil {
		return bosherr.WrapError(err, "Getting environments")
	}

	newState := State{}
	isSet := false

	for _, existingEnv := range envs {
		if existingEnv.Alias == env.Alias || existingEnv.URL == env.URL {
			if isSet == false {
				newState.Environments = append(newState.Environments, env)
			}

			isSet = true
		} else {
			newState.Environments = append(newState.Environments, existingEnv)
		}
	}

	if !isSet {
		newState.Environments = append(newState.Environments, env)
	}

	// @todo sort

	_, err = m.storage.Put(m.path, newState)
	if err != nil {
		return bosherr.WrapError(err, "Putting environment")
	}

	return nil
}
