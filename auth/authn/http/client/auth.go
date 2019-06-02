package client

import (
	"net/http"

	"github.com/pkg/errors"

	env_api "github.com/dpb587/ssoca/service/env/api"
)

func (s Service) AuthLogin(_ env_api.InfoServiceResponse) (interface{}, error) {
	ui := s.runtime.GetUI()
	auth := AuthConfig{}

	username, err := ui.AskForText("username")
	if err != nil {
		return auth, errors.Wrap(err, "requesting username")
	}

	auth.Username = username

	password, err := ui.AskForPassword("password")
	if err != nil {
		return auth, errors.Wrap(err, "requesting password")
	}

	auth.Password = password

	return auth, nil
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
