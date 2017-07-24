package client

import (
	"fmt"
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/dpb587/ssoca/client/auth"
	"github.com/dpb587/ssoca/client/config"
	env_api "github.com/dpb587/ssoca/service/env/api"
)

func (s Service) AuthLogin(_ env_api.InfoServiceResponse) (interface{}, error) {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return nil, bosherr.WrapError(err, "Getting environment")
	}

	openCommand := config.NewStringSliceEnvironmentOption(config.EnvironmentOptionAuthOpenCommand)
	err = env.GetOption(&openCommand, []string{"open"})
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading option")
	}

	str := auth.NewServerTokenRetrieval(env.URL, s.cmdRunner, openCommand.GetValue(), s.runtime.GetStdout(), s.runtime.GetStdin())

	token, err := str.Retrieve("/auth/initiate")
	if err != nil {
		return nil, bosherr.WrapError(err, "Waiting for user token")
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
		return bosherr.WrapError(err, "Getting environment")
	}

	authConfig := AuthConfig{}
	err = env.Auth.UnmarshalOptions(&authConfig)
	if err != nil {
		return bosherr.WrapError(err, "Parsing authentication options")
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authConfig.Token))

	return nil
}
