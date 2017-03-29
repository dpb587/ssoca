package client

import (
	"github.com/dpb587/ssoca/client"
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"

	svc "github.com/dpb587/ssoca/service/ssh"
	svccmd "github.com/dpb587/ssoca/service/ssh/client/cmd"
	svchttpclient "github.com/dpb587/ssoca/service/ssh/httpclient"

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
	return "Interact with remote SSH servers"
}

func (s Service) GetCommand() interface{} {
	cmd := clientcmd.ServiceCommand{
		Runtime:     s.runtime,
		ServiceName: s.Type(),
	}

	return &struct {
		Agent         svccmd.Agent         `command:"agent" description:"Start an SSH agent"`
		Exec          svccmd.Exec          `command:"exec" description:"Connect to a remote SSH server"`
		SignPublicKey svccmd.SignPublicKey `command:"sign-public-key" description:"Create a certificate for a specific public key"`
	}{
		Agent: svccmd.Agent{
			ServiceCommand: cmd,
			CmdRunner:      s.cmdRunner,
			FS:             s.fs,
			GetClient:      s.GetClient,
		},
		Exec: svccmd.Exec{
			ServiceCommand: cmd,
			CmdRunner:      s.cmdRunner,
			FS:             s.fs,
			GetClient:      s.GetClient,
		},
		SignPublicKey: svccmd.SignPublicKey{
			ServiceCommand: cmd,
			FS:             s.fs,
			GetClient:      s.GetClient,
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
