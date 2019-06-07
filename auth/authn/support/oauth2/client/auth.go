package client

import (
	"fmt"
	"net/http"

	"github.com/dpb587/ssoca/client/auth"
	"github.com/dpb587/ssoca/client/config"
	"github.com/pkg/errors"
)

func (s Service) AuthLogin() error {
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	authBind := config.EnvironmentOptionAuthBind
	err = env.GetOption(&authBind)
	if err != nil {
		return errors.Wrap(err, "loading bind option")
	}

	openCommand := config.EnvironmentOptionAuthOpenCommand
	err = env.GetOption(&openCommand)
	if err != nil {
		return errors.Wrap(err, "loading open option")
	}

	str := auth.NewServerTokenRetrieval(env.URL, s.runtime.GetVersion(), s.cmdRunner, authBind.GetValue(), openCommand.GetValue(), s.runtime.GetStderr(), s.runtime.GetStdin())

	token, err := str.Retrieve(fmt.Sprintf("/%s/initiate", s.name))
	if err != nil {
		return errors.Wrap(err, "waiting for user token")
	}

	env, err = s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	env.Auth = &config.EnvironmentAuthState{
		Type: s.name,
		Options: AuthConfig{
			Token: token,
		},
	}

	configManager, err := s.runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting config manager")
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

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authConfig.Token))

	return nil
}
