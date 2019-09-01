package cli

import (
	"context"
	"fmt"
	"time"

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

	SkipVerify  bool          `long:"skip-verify" description:"Skip verification of authentication, once complete"`
	WaitTimeout time.Duration `long:"wait-timeout" description:"Timeout to wait for authentication before erroring" default:"15m"`

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

	// intentionally not defaulting to any existing auth token's service to force
	// explicit invocations when users run the `login` command; the exception
	// being when internal reauthentication attempts occur and it automatically
	// passes whatever service the previous token was using.
	authServiceName := c.ServiceCommand.ServiceName

	if authServiceName == "" {
		authServiceName = envInfo.Env.DefaultAuthService

		if authServiceName == "" {
			// deprecated; fallback to older style of reporting
			authServiceName = "auth"
		}
	}

	// find service named auth
	var authServiceListing envapi.InfoServiceResponse

	for _, serviceListing := range envInfo.Services {
		if serviceListing.Name != authServiceName {
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

	svc, err := c.ServiceManager.Get(authServiceType, authServiceListing.Name)
	if err != nil {
		return errors.Wrap(err, "loading auth service")
	}

	authService, ok := svc.(service.AuthService)
	if !ok {
		return fmt.Errorf("cannot authenticate with service: %s", authServiceType)
	}

	ctx := context.Background()

	if c.WaitTimeout > 0 {
		var ctxCancel context.CancelFunc

		ctx, ctxCancel = context.WithTimeout(ctx, c.WaitTimeout)
		defer ctxCancel()
	}

	err = authService.AuthLogin(ctx)
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
