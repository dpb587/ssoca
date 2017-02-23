package client

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	env_api "github.com/dpb587/ssoca/service/env/api"
)

func (s Service) AuthLogin(_ env_api.InfoServiceResponse) (interface{}, error) {
	ui := s.runtime.GetUI()
	auth := AuthConfig{}

	username, err := ui.AskForText("username")
	if err != nil {
		return auth, bosherr.WrapError(err, "Requesting username")
	}

	auth.Username = username

	password, err := ui.AskForPassword("password")
	if err != nil {
		return auth, bosherr.WrapError(err, "Requesting password")
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
		return bosherr.WrapError(err, "Getting environment")
	}

	authConfig := AuthConfig{}
	err = env.Auth.UnmarshalOptions(&authConfig)
	if err != nil {
		return bosherr.WrapError(err, "Parsing authentication options")
	}

	req.SetBasicAuth(authConfig.Username, authConfig.Password)

	return nil
}
