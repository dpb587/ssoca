package cli

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/auth/client"
)

type Commands struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	Info   Info   `command:"info" description:"Show current authentication information"`
	Login  Login  `command:"login" description:"Authenticate for a new token"`
	Logout Logout `command:"logout" description:"Revoke an authentication token"`
}

func CreateCommands(runtime client.Runtime, s svc.Service) *Commands {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: svc.Service{}.Type(),
	}

	return &Commands{
		ServiceCommand: cmd,

		Info: Info{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		Login: Login{
			ServiceCommand: cmd,
			ServiceManager: s.GetServiceManager(),
			GetClient:      s.GetClient,
		},
		Logout: Logout{
			ServiceCommand: cmd,
			ServiceManager: s.GetServiceManager(),
		},
	}
}
