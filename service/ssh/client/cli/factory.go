package cli

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/ssh/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Commands struct {
	Agent         Agent         `command:"agent" description:"Start an SSH agent"`
	Exec          Exec          `command:"exec" description:"Connect to a remote SSH server"`
	SignPublicKey SignPublicKey `command:"sign-public-key" description:"Create a certificate for a specific public key"`

	sf svc.ServiceFactory
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *Commands {
	cmd := clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: svc.Service{}.Type(),
	}

	return &Commands{
		Agent: Agent{
			serviceFactory: sf,
			ServiceCommand: cmd,
			cmdRunner:      cmdRunner,
			fs:             fs,
		},
		Exec: Exec{
			serviceFactory: sf,
			ServiceCommand: cmd,
		},
		SignPublicKey: SignPublicKey{
			serviceFactory: sf,
			ServiceCommand: cmd,
			fs:             fs,
		},
	}
}
