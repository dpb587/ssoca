package client

import (
	"os"

	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/httpclient"

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

var _ service.Service = Service{}

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
		Connect: svccmd.Connect{
			ServiceCommand: cmd,
			FS:             s.fs,
			CmdRunner:      s.cmdRunner,
			GetClient:      s.GetClient,
			ExecutableFinder: svcclienthelper.ExecutableFinder{
				FS: s.fs,
			},
		},
		CreateProfile: svccmd.CreateProfile{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
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

func (s Service) GetClient(service string, skipAuthRetry bool) (svchttpclient.Client, error) {
	var client httpclient.Client
	var err error

	if skipAuthRetry {
		client, err = s.runtime.GetClient()
	} else {
		client, err = s.runtime.GetAuthInterceptClient()
	}

	if err != nil {
		return nil, err
	}

	return svchttpclient.New(client, service)
}
