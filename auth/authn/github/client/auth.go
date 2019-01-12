package client

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/client/auth"
	"github.com/dpb587/ssoca/client/config"
	env_api "github.com/dpb587/ssoca/service/env/api"
	"github.com/pkg/errors"
)

func (s Service) AuthLogin(_ env_api.InfoServiceResponse) (interface{}, error) {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "Getting environment")
	}

	authBind := config.NewStringEnvironmentOption(config.EnvironmentOptionAuthBind)
	err = env.GetOption(&authBind, "0.0.0.0:0")
	if err != nil {
		return nil, errors.Wrap(err, "Loading bind option")
	}

	openCommand := config.NewStringSliceEnvironmentOption(config.EnvironmentOptionAuthOpenCommand)
	err = env.GetOption(&openCommand, []string{"open"})
	if err != nil {
		return nil, errors.Wrap(err, "Loading open option")
	}

	str := auth.NewServerTokenRetrieval(env.URL, s.runtime.GetVersion(), s.cmdRunner, authBind.GetValue(), openCommand.GetValue(), s.runtime.GetStderr(), s.runtime.GetStdin())

	token, err := str.Retrieve("/auth/initiate")
	if err != nil {
		return nil, errors.Wrap(err, "Waiting for user token")
	}

	authConfig := AuthConfig{
		Token: token,
	}

	return authConfig, nil
}

func (s Service) AuthLogout() error {
	return nil
}

func (s Service) AuthRequest(req *http.Request) error {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "Getting environment")
	}

	authConfig := AuthConfig{}
	err = env.Auth.UnmarshalOptions(&authConfig)
	if err != nil {
		return errors.Wrap(err, "Parsing authentication options")
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authConfig.Token))

	return nil
}
