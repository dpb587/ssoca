package client

import (
	"fmt"
	"net/http"

	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/dpb587/ssoca/auth/authn/uaa/api"
	"github.com/dpb587/ssoca/config"
	env_api "github.com/dpb587/ssoca/service/env/api"
)

func (s Service) AuthLogin(remoteService env_api.InfoServiceResponse) (interface{}, error) {
	metadata := api.Metadata{}

	err := config.RemarshalJSON(remoteService.Metadata, &metadata)
	if err != nil {
		return nil, bosherr.WrapError(err, "Parsing metadata")
	}

	client, err := s.uaaClientFactory.CreateClient(metadata.URL, "bosh_cli", "", metadata.CACertificate)
	if err != nil {
		return nil, bosherr.WrapError(err, "Creating UAA client")
	}

	prompts, err := client.Prompts()
	if err != nil {
		return nil, bosherr.WrapError(err, "Discovering UAA prompts")
	}

	ui := s.runtime.GetUI()
	var answers []boshuaa.PromptAnswer

	for _, prompt := range prompts {
		var askFunc func(string) (string, error)

		if prompt.IsPassword() {
			askFunc = ui.AskForPassword
		} else {
			askFunc = ui.AskForText
		}

		value, err1 := askFunc(prompt.Label)
		if err1 != nil {
			return nil, err1
		}

		if value != "" {
			answer := boshuaa.PromptAnswer{Key: prompt.Key, Value: value}
			answers = append(answers, answer)
		}
	}

	accessToken, err := client.OwnerPasswordCredentialsGrant(answers)
	if err != nil {
		return nil, bosherr.WrapError(err, "Fetching credentials grant")
	}

	auth := AuthConfig{
		URL:           metadata.URL,
		CACertificate: metadata.CACertificate,
		RefreshToken:  accessToken.RefreshToken().Value(),
	}

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

	client, err := s.uaaClientFactory.CreateClient(authConfig.URL, "bosh_cli", "", authConfig.CACertificate)
	if err != nil {
		return bosherr.WrapError(err, "Creating UAA client")
	}

	staleToken := client.NewStaleAccessToken(authConfig.RefreshToken)
	accessToken, err := staleToken.Refresh()
	if err != nil {
		return bosherr.WrapError(err, "Refreshing token")
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", accessToken.Type(), accessToken.Value()))

	return nil
}
