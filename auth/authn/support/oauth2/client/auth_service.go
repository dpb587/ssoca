package client

import (
	"fmt"
	"net/http"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/auth"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/service"
	"github.com/pkg/errors"
)

type AuthService struct {
	serviceName string
	serviceType service.Type
	runtime     client.Runtime
	cmdRunner   boshsys.CmdRunner
}

func NewAuthService(serviceName string, serviceType service.Type, runtime client.Runtime, cmdRunner boshsys.CmdRunner) *AuthService {
	return &AuthService{
		serviceName: serviceName,
		serviceType: serviceType,
		runtime:     runtime,
		cmdRunner:   cmdRunner,
	}
}

func (as AuthService) AuthLogin() error {
	configManager, err := as.runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting config manager")
	}

	env, err := as.runtime.GetEnvironment()
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

	str := auth.NewServerTokenRetrieval(
		env.URL,
		as.runtime.GetVersion(),
		as.cmdRunner,
		authBind.GetValue(),
		openCommand.GetValue(),
		as.runtime.GetStderr(),
		as.runtime.GetStdin(),
	)

	token, err := str.Retrieve(fmt.Sprintf("/%s/initiate", as.serviceName))
	if err != nil {
		return errors.Wrap(err, "waiting for user token")
	}

	// get the very latest version
	env, err = as.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	env.Auth = &config.EnvironmentAuthState{
		Name: as.serviceName,
		Type: string(as.serviceType),
		Options: authn.AuthorizationToken{
			Type:  "Bearer",
			Value: token,
		},
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "updating environment")
	}

	return nil
}

func (as AuthService) AuthLogout() error {
	configManager, err := as.runtime.GetConfigManager()
	if err != nil {
		return errors.Wrap(err, "getting config manager")
	}

	env, err := as.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	env.Auth = nil

	err = configManager.SetEnvironment(env)
	if err != nil {
		return errors.Wrap(err, "updating environment")
	}

	return nil
}

func (as AuthService) AuthRequest(req *http.Request) error {
	env, err := as.runtime.GetEnvironment()
	if err != nil {
		return errors.Wrap(err, "getting environment")
	}

	if env.Auth.Options == nil {
		// should never happen
		return nil
	}

	authConfig := authn.AuthorizationToken{}
	err = env.Auth.UnmarshalOptions(&authConfig)
	if err != nil {
		return errors.Wrap(err, "parsing authentication options")
	}

	authType := authConfig.Type
	authValue := authConfig.Value

	// deprecated fallback
	if authValue == "" && authConfig.Token != "" {
		authType = "bearer"
		authValue = authConfig.Token
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", authType, authValue))

	return nil
}
