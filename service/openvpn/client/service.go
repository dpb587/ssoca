package client

import (
	"os"

	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"

	svc "github.com/dpb587/ssoca/service/openvpn"
	svccmd "github.com/dpb587/ssoca/service/openvpn/client/cmd"
	svcclienthelper "github.com/dpb587/ssoca/service/openvpn/client/helper"
	svchttpclient "github.com/dpb587/ssoca/service/openvpn/httpclient"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.Service

	runtime   client.Runtime
	fs        boshsys.FileSystem
	cmdRunner boshsys.CmdRunner
}

func NewService(runtime client.Runtime, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) Service {
	return Service{
		runtime:   runtime,
		fs:        fs,
		cmdRunner: cmdRunner,
	}
}

func (s Service) Description() string {
	return "Establish an OpenVPN to a remote server"
}

func (s Service) GetCommand() interface{} {
	cmd := clientcmd.ServiceCommand{
		Runtime:     s.runtime,
		ServiceName: s.Type(),
	}

	connect := svccmd.Connect{
		ServiceCommand: cmd,
		FS:             s.fs,
		CmdRunner:      s.cmdRunner,
		CreateProfile:  s.CreateProfile,
		ExecutableFinder: svcclienthelper.ExecutableFinder{
			FS: s.fs,
		},
	}

	return &struct {
		BaseProfile              svccmd.BaseProfile              `command:"base-profile" description:"Show the base connection profile of the OpenVPN server"`
		Connect                  svccmd.Connect                  `command:"connect" description:"Connect to a remote OpenVPN server"`
		CreateProfile            svccmd.CreateProfile            `command:"create-profile" description:"Create and sign an OpenVPN configuration profile"`
		CreateTunnelblickProfile svccmd.CreateTunnelblickProfile `command:"create-tunnelblick-profile" description:"Create a Tunnelblick profile"`
	}{
		BaseProfile: svccmd.BaseProfile{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		Connect: connect,
		CreateProfile: svccmd.CreateProfile{
			ServiceCommand:    cmd,
			CreateUserProfile: s.CreateProfile,
		},
		CreateTunnelblickProfile: svccmd.CreateTunnelblickProfile{
			ServiceCommand: cmd,
			Name:           s.Type(),
			SssocaExec:     os.Args[0],
			GetClient:      s.GetClient,
			FS:             s.fs,
		},
	}
}

func (s Service) GetClient(service string) (*svchttpclient.Client, error) {
	client, err := s.runtime.GetClient()
	if err != nil {
		return &svchttpclient.Client{}, err
	}

	return svchttpclient.New(client, service)
}
