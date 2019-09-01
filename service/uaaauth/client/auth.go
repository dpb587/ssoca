package client

import (
	"context"
	"fmt"
	"net/http"

	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	"github.com/pkg/errors"

	clientconfig "github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/config"
	"github.com/dpb587/ssoca/service"
	"github.com/dpb587/ssoca/service/env"
	env_api "github.com/dpb587/ssoca/service/env/api"
	envclient "github.com/dpb587/ssoca/service/env/client"
	"github.com/dpb587/ssoca/service/uaaauth/api"
)

func (s Service) AuthLogin(_ context.Context) error {
	configManager, err := s.runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting config manager")
	}

	rawEnvService, err := s.serviceManager.Get(env.Type, "env")
	if err != nil {
		return errors.Wrap(err, "getting env service")
	}

	envService, ok := rawEnvService.(*envclient.Service)
	if !ok {
		return errors.Wrap(err, "expecting env service")
	}

	envClient, err := envService.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting env HTTP client")
	}

	envInfo, err := envClient.GetInfo()
	if err != nil {
		return errors.Wrap(err, "getting environment info")
	}

	var authServiceListing env_api.InfoServiceResponse

	for _, serviceListing := range envInfo.Services {
		if serviceListing.Name != s.name || service.Type(serviceListing.Type) != s.Type() {
			continue
		}

		authServiceListing = serviceListing

		break
	}

	if authServiceListing.Type == "" {
		return fmt.Errorf("expected to find remote service: type=%s, name=%s", s.Type(), s.name)
	}

	metadata := api.Metadata{}

	err = config.RemarshalJSON(authServiceListing.Metadata, &metadata)
	if err != nil {
		return errors.Wrap(err, "parsing metadata")
	}

	client, err := s.uaaClientFactory.CreateClient(metadata.URL, metadata.ClientID, metadata.ClientSecret, metadata.CACertificate)
	if err != nil {
		return errors.Wrap(err, "creating UAA client")
	}

	prompts, err := client.Prompts()
	if err != nil {
		return errors.Wrap(err, "discovering UAA prompts")
	}

	ui := s.runtime.GetUI()
	var answers []boshuaa.PromptAnswer

	for _, prompt := range prompts {
		if len(metadata.Prompts) > 0 {
			var matchedPrompt bool

			for _, expectedKey := range metadata.Prompts {
				if prompt.Key == expectedKey {
					matchedPrompt = true

					break
				}
			}

			if !matchedPrompt {
				continue
			}
		}

		var askFunc func(string) (string, error)

		if prompt.IsPassword() {
			askFunc = ui.AskForPassword
		} else {
			askFunc = ui.AskForText
		}

		value, err1 := askFunc(prompt.Label)
		if err1 != nil {
			return err1
		}

		if value != "" {
			answer := boshuaa.PromptAnswer{Key: prompt.Key, Value: value}
			answers = append(answers, answer)
		}
	}

	accessToken, err := client.OwnerPasswordCredentialsGrant(answers)
	if err != nil {
		return errors.Wrap(err, "fetching credentials grant")
	}

	// get the very latest version
	env, err := s.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	env.Auth = &clientconfig.EnvironmentAuthState{
		Name: s.name,
		Type: string(s.Type()),
		Options: AuthConfig{
			URL:           metadata.URL,
			CACertificate: metadata.CACertificate,
			ClientID:      metadata.ClientID,
			ClientSecret:  metadata.ClientSecret,
			RefreshToken:  accessToken.RefreshToken().Value(),
		},
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "updating environment")
	}

	return nil
}

func (s Service) AuthLogout(_ context.Context) error {
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

	client, err := s.uaaClientFactory.CreateClient(authConfig.URL, authConfig.ClientID, authConfig.ClientSecret, authConfig.CACertificate)
	if err != nil {
		return errors.Wrap(err, "creating UAA client")
	}

	staleToken := client.NewStaleAccessToken(authConfig.RefreshToken)
	accessToken, err := staleToken.Refresh()
	if err != nil {
		return errors.Wrap(err, "refreshing token")
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", accessToken.Type(), accessToken.Value()))

	return nil
}
