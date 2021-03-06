package config

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/config/storage"
)

type DefaultManager struct {
	storage storage.Storage
	path    string
}

var _ Manager = DefaultManager{}

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
		return EnvironmentsState{}, errors.Wrap(err, "getting environments")
	}

	return state.Environments, nil
}

func (m DefaultManager) GetEnvironment(name string) (EnvironmentState, error) {
	envs, err := m.GetEnvironments()
	if err != nil {
		return EnvironmentState{}, errors.Wrap(err, "getting environment")
	}

	for _, env := range envs {
		if env.Alias == name || env.URL == name {
			return env, nil
		}
	}

	return EnvironmentState{}, fmt.Errorf("environment not found: %s", name)
}

func (m DefaultManager) SetEnvironment(env EnvironmentState) error {
	envs, err := m.GetEnvironments()
	if err != nil {
		return errors.Wrap(err, "getting environments")
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
		return errors.Wrap(err, "putting environment")
	}

	return nil
}

func (m DefaultManager) UnsetEnvironment(name string) error {
	envs, err := m.GetEnvironments()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	newState := State{}

	for _, env := range envs {
		if env.Alias == name || env.URL == name {
			continue
		}

		newState.Environments = append(newState.Environments, env)
	}

	_, err = m.storage.Put(m.path, newState)
	if err != nil {
		return errors.Wrap(err, "putting environment")
	}

	return nil
}
