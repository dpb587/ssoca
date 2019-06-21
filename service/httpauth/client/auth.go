package client

import (
	"net/http"

	"github.com/dpb587/ssoca/client/config"
	"github.com/pkg/errors"
)

func (s Service) AuthLogin() error {
	configManager, err := s.runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting config manager")
	}

	ui := s.runtime.GetUI()
	auth := AuthConfig{}

	username, err := ui.AskForText("username")
	if err != nil {
		return errors.Wrap(err, "requesting username")
	}

	auth.Username = username

	password, err := ui.AskForPassword("password")
	if err != nil {
		return errors.Wrap(err, "requesting password")
	}

	auth.Password = password

	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	env.Auth = &config.EnvironmentAuthState{
		Name:    s.name,
		Type:    string(s.Type()),
		Options: auth,
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "updating environment")
	}

	return nil
}

func (s Service) AuthLogout() error {
	return nil
}

func (s Service) AuthRequest(req *http.Request) error {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	authConfig := AuthConfig{}
	err = env.Auth.UnmarshalOptions(&authConfig)
	if err != nil {
		return errors.Wrap(err, "parsing authentication options")
	}

	req.SetBasicAuth(authConfig.Username, authConfig.Password)

	return nil
}
