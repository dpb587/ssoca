package cli

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/file/client"
)

type Commands struct {
	*clientcmd.ServiceCommand

	Get  Get  `command:"get" description:"Get a file"`
	List List `command:"list" description:"List files"`
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory) *Commands {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: svc.Service{}.Type(),
	}

	return &Commands{
		ServiceCommand: cmd,
		Get: Get{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
		List: List{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
	}
}
