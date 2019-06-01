package cli

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/file/client"
)

type Commands struct {
	*clientcmd.ServiceCommand

	Exec Exec `command:"exec" description:"Temporarily get and then execute a file"`
	Get  Get  `command:"get" description:"Download a file and verify its checksum" alias:"download"`
	List List `command:"list" description:"List available files"`
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory) *Commands {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: svc.Service{}.Type(),
	}

	return &Commands{
		ServiceCommand: cmd,
		Exec: Exec{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
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
