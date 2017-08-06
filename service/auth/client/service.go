package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/auth"
	svccmd "github.com/dpb587/ssoca/service/auth/client/cmd"
	svchttpclient "github.com/dpb587/ssoca/service/auth/httpclient"
)

type Service struct {
	svc.Service

	runtime        client.Runtime
	serviceManager service.Manager
}

var _ service.Service = Service{}

func NewService(runtime client.Runtime, serviceManager service.Manager) Service {
	return Service{
		runtime:        runtime,
		serviceManager: serviceManager,
	}
}

func (s Service) Description() string {
	return "Manage authentication"
}

func (s Service) GetCommand() interface{} {
	cmd := clientcmd.ServiceCommand{
		Runtime:     s.runtime,
		ServiceName: s.Type(),
	}

	return &struct {
		Info   svccmd.Info   `command:"info" description:"Show current authentication information"`
		Login  svccmd.Login  `command:"login" description:"Authenticate for a new token"`
		Logout svccmd.Logout `command:"logout" description:"Revoke an authentication token"`
	}{
		Info: svccmd.Info{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		Login: svccmd.Login{
			ServiceCommand: cmd,
			ServiceManager: s.serviceManager,
			GetClient:      s.GetClient,
		},
		Logout: svccmd.Logout{
			ServiceCommand: cmd,
			ServiceManager: s.serviceManager,
		},
	}
}

func (s Service) GetClient() (svchttpclient.Client, error) {
	client, err := s.runtime.GetClient()
	if err != nil {
		return nil, err
	}

	return svchttpclient.New(client)
}
