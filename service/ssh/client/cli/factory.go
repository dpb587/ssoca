package cli

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/ssh"
	svc "github.com/dpb587/ssoca/service/ssh/client"
)

type Commands struct {
	*clientcmd.ServiceCommand

	Agent         Agent         `command:"agent" description:"Start an SSH agent"`
	Exec          Exec          `command:"exec" description:"Connect to a remote SSH server"`
	SignPublicKey SignPublicKey `command:"sign-public-key" description:"Create a certificate for a specific public key"`

	sf svc.ServiceFactory
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *Commands {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: string(ssh.Type),
	}

	return &Commands{
		ServiceCommand: cmd,

		Agent: Agent{
			ServiceCommand: cmd,
			serviceFactory: sf,
			cmdRunner:      cmdRunner,
			fs:             fs,
		},
		Exec: Exec{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
		SignPublicKey: SignPublicKey{
			ServiceCommand: cmd,
			serviceFactory: sf,
			fs:             fs,
		},
	}
}
