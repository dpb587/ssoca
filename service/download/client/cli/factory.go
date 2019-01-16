package cli

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/download/client"
)

type Commands struct {
	Get  Get  `command:"get" description:"Get an artifact"`
	List List `command:"list" description:"List available artifacts"`
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory) *Commands {
	cmd := clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: svc.Service{}.Type(),
	}

	return &Commands{
		Get: Get{
			serviceFactory: sf,
			ServiceCommand: cmd,
		},
		List: List{
			serviceFactory: sf,
			ServiceCommand: cmd,
		},
	}
}
