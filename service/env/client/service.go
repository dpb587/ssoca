package client

import (
	"os"

	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/httpclient"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/env"
	svccmd "github.com/dpb587/ssoca/service/env/client/cmd"
	svchttpclient "github.com/dpb587/ssoca/service/env/httpclient"
	svcdownloadhttpclient "github.com/dpb587/ssoca/service/file/httpclient"
)

type Service struct {
	svc.ServiceType

	runtime   client.Runtime
	fs        boshsys.FileSystem
	cmdRunner boshsys.CmdRunner
}

var _ service.Service = Service{}

func NewService(runtime client.Runtime, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *Service {
	return &Service{
		runtime:   runtime,
		fs:        fs,
		cmdRunner: cmdRunner,
	}
}

func (s Service) Description() string {
	return "Manage environment references"
}

func (s Service) GetCommand() interface{} {
	cmd := &clientcmd.ServiceCommand{
		Runtime:     s.runtime,
		ServiceName: s.Type(),
	}

	return &struct {
		*clientcmd.ServiceCommand `no-flag:"true"`

		Set          svccmd.Set          `command:"set" description:"Configure the connection to an environment" alias:"add"`
		Info         svccmd.Info         `command:"info" description:"Show environment information"`
		Services     svccmd.Services     `command:"services" description:"Show current services available from the environment"`
		List         svccmd.List         `command:"list" description:"List all locally-configured environments"`
		Rename       svccmd.Rename       `command:"rename" description:"Set a new name for the environment"`
		SetOption    svccmd.SetOption    `command:"set-option" description:"Set a local client option in the environment"`
		UpdateClient svccmd.UpdateClient `command:"update-client" description:"Download the latest client from the environment"`
		Unset        svccmd.Unset        `command:"unset" description:"Remove all configuration for an environment" alias:"remove"`
	}{
		ServiceCommand: cmd,

		Set: svccmd.Set{
			ServiceCommand: cmd,
			FS:             s.fs,
			GetClient:      s.GetClient,
		},
		Info: svccmd.Info{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		Services: svccmd.Services{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
		},
		List: svccmd.List{
			ServiceCommand: cmd,
		},
		Rename: svccmd.Rename{
			ServiceCommand: cmd,
		},
		SetOption: svccmd.SetOption{
			ServiceCommand: cmd,
		},
		UpdateClient: svccmd.UpdateClient{
			ServiceCommand:    cmd,
			SsocaExec:         os.Args[0],
			FS:                s.fs,
			CmdRunner:         s.cmdRunner,
			GetClient:         s.GetClient,
			GetDownloadClient: s.getDownloadClient,
		},
		Unset: svccmd.Unset{
			ServiceCommand: cmd,
		},
	}
}

func (s Service) GetClient() (svchttpclient.Client, error) {
	client, err := s.runtime.GetClient()
	if err != nil {
		return nil, err
	}

	return svchttpclient.New(client)
}

func (s Service) getDownloadClient(service string, skipAuthRetry bool) (svcdownloadhttpclient.Client, error) {
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

	return svcdownloadhttpclient.New(client, service)
}
