package cli

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	globalservice "github.com/dpb587/ssoca/service"
	envapi "github.com/dpb587/ssoca/service/env/api"
	envclient "github.com/dpb587/ssoca/service/env/client"
)

type Login struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	SkipVerify bool `long:"skip-verify" description:"Skip verification of authentication, once complete"`

	ServiceManager service.Manager
	GetClient      GetClient
}

var _ flags.Commander = Login{}

func (c Login) Execute(_ []string) error {
	rawEnvService, err := c.ServiceManager.Get("env", "env")
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

	// find service named auth
	var authServiceListing envapi.InfoServiceResponse

	for _, serviceListing := range envInfo.Services {
		if serviceListing.Name != "auth" {
			continue
		}

		authServiceListing = serviceListing

		break
	}

	// deprecated; fallback to older style of reporting
	if authServiceListing.Name == "" && envInfo.Auth != nil {
		authServiceListing = *envInfo.Auth
		authServiceListing.Name = "auth"
	}

	if authServiceListing.Name == "" {
		return errors.New("failed to find auth service")
	}

	authServiceType := globalservice.Type(authServiceListing.Type)

	svc, err := c.ServiceManager.Get(authServiceType, "auth")
	if err != nil {
		return errors.Wrap(err, "loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("cannot authenticate with service: %s", authServiceType)
	}

	err = authService.AuthLogin()
	if err != nil {
		return errors.Wrapf(err, "executing login with %s", authService.Type())
	}

	if c.SkipVerify {
		return nil
	}

	err = c.verify()
	if err != nil {
		return errors.Wrap(err, "verifying authentication")
	}

	return nil
}

func (c Login) verify() error {
	ui := c.Runtime.GetUI()

	client, err := c.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting client")
	}

	authInfo, err := client.GetAuth()
	if err != nil {
		return errors.Wrap(err, "getting remote authentication info")
	}

	if authInfo.ID == "" {
		return errors.New("failed to use authentication credentials")
	}

	ui.PrintBlock([]byte(fmt.Sprintf("Successfully logged in as %s\n", authInfo.ID)))

	return nil
}
