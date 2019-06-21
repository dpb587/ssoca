package cli

import (
	"os"

	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/openvpn"
	svc "github.com/dpb587/ssoca/service/openvpn/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Commands struct {
	*clientcmd.ServiceCommand

	BaseProfile              BaseProfile              `command:"base-profile" description:"Show the base connection profile of the OpenVPN server"`
	Exec                     Exec                     `command:"exec" description:"Execute openvpn to connect to the remote server" alias:"connect"`
	CreateONCProfile         CreateONCProfile         `command:"create-onc-profile" description:"Create an ONC profile"`
	CreateProfile            CreateProfile            `command:"create-profile" description:"Create and sign an OpenVPN configuration profile"`
	CreateTunnelblickProfile CreateTunnelblickProfile `command:"create-tunnelblick-profile" description:"Create a Tunnelblick profile"`
	CreateLaunchdService     CreateLaunchdService     `command:"create-launchd-service" description:"Create a launchd service"`

	sf svc.ServiceFactory
}

func CreateCommands(runtime client.Runtime, sf svc.ServiceFactory, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *Commands {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     runtime,
		ServiceName: string(openvpn.Type),
	}

	return &Commands{
		ServiceCommand: cmd,

		BaseProfile: BaseProfile{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
		Exec: Exec{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
		CreateONCProfile: CreateONCProfile{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
		CreateProfile: CreateProfile{
			ServiceCommand: cmd,
			serviceFactory: sf,
		},
		CreateTunnelblickProfile: CreateTunnelblickProfile{
			ServiceCommand: cmd,
			serviceFactory: sf,
			fs:             fs,
			cmdRunner:      cmdRunner,
			SsocaExec:      os.Args[0],
		},
		CreateLaunchdService: CreateLaunchdService{
			ServiceCommand: cmd,
			serviceFactory: sf,
			fs:             fs,
			cmdRunner:      cmdRunner,
			SsocaExec:      os.Args[0],
		},
	}
}
