package cli

import (
	"os"

	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/openvpn/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Commands struct {
	BaseProfile              BaseProfile              `command:"base-profile" description:"Show the base connection profile of the OpenVPN server"`
	Exec                     Exec                     `command:"exec" description:"Execute openvpn to connect to the remote server" alias:"connect"`
	CreateProfile            CreateProfile            `command:"create-profile" description:"Create and sign an OpenVPN configuration profile"`
	CreateTunnelblickProfile CreateTunnelblickProfile `command:"create-tunnelblick-profile" description:"Create a Tunnelblick profile"`
	CreateLaunchdService     CreateLaunchdService     `command:"create-launchd-service" description:"Create a launchd service"`

	sf svc.ServiceFactory
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *Commands {
	cmd := clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: svc.Service{}.Type(),
	}

	return &Commands{
		BaseProfile: BaseProfile{
			serviceFactory: sf,
			ServiceCommand: cmd,
		},
		Exec: Exec{
			serviceFactory: sf,
			ServiceCommand: cmd,
		},
		CreateProfile: CreateProfile{
			serviceFactory: sf,
			ServiceCommand: cmd,
		},
		CreateTunnelblickProfile: CreateTunnelblickProfile{
			serviceFactory: sf,
			fs:             fs,
			cmdRunner:      cmdRunner,
			ServiceCommand: cmd,
			SsocaExec:      os.Args[0],
		},
		CreateLaunchdService: CreateLaunchdService{
			serviceFactory: sf,
			fs:             fs,
			cmdRunner:      cmdRunner,
			ServiceCommand: cmd,
			SsocaExec:      os.Args[0],
		},
	}
}
