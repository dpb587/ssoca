package cli

import (
	"os"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/env/client"
)

type Commands struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	Set          Set          `command:"set" description:"Configure the connection to an environment" alias:"add"`
	Info         Info         `command:"info" description:"Show environment information"`
	Services     Services     `command:"services" description:"Show current services available from the environment"`
	List         List         `command:"list" description:"List all locally-configured environments"`
	Rename       Rename       `command:"rename" description:"Set a new name for the environment"`
	SetOption    SetOption    `command:"set-option" description:"Set a local client option in the environment"`
	UpdateClient UpdateClient `command:"update-client" description:"Download the latest client from the environment"`
	Unset        Unset        `command:"unset" description:"Remove all configuration for an environment" alias:"remove"`
}

// TODO convert to service factory rather than single env instance
func CreateCommands(runtime client.Runtime, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem, s *svc.Service) *Commands {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: s.Type(),
	}

	return &Commands{
		ServiceCommand: cmd,

		Set: Set{
			ServiceCommand: cmd,
			FS:             fs,
			GetClient:      s.GetClient,
		},
		Info: Info{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		Services: Services{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		List: List{
			ServiceCommand: cmd,
		},
		Rename: Rename{
			ServiceCommand: cmd,
		},
		SetOption: SetOption{
			ServiceCommand: cmd,
		},
		UpdateClient: UpdateClient{
			ServiceCommand:    cmd,
			SsocaExec:         os.Args[0],
			FS:                fs,
			CmdRunner:         cmdRunner,
			GetClient:         s.GetClient,
			GetDownloadClient: s.GetDownloadClient,
		},
		Unset: Unset{
			ServiceCommand: cmd,
		},
	}
}
