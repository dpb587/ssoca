package cli

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/service/file"
)

type Commands struct {
	*ServiceCommand

	Exec Exec `command:"exec" description:"Temporarily get and then execute a file"`
	Get  Get  `command:"get" description:"Download a file and verify its checksum" alias:"download"`
	List List `command:"list" description:"List available files"`
}

func CreateCommands(runtime client.Runtime, manager service.Manager) *Commands {
	cmd := &ServiceCommand{
		clientcmd.ServiceCommand{
			Runtime:        runtime,
			ServiceManager: manager,
			ServiceType:    file.Type,
			ServiceName:    string(file.Type),
		},
	}

	return &Commands{
		ServiceCommand: cmd,
		Exec: Exec{
			ServiceCommand: cmd,
		},
		Get: Get{
			ServiceCommand: cmd,
		},
		List: List{
			ServiceCommand: cmd,
		},
	}
}
