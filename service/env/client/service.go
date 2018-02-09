package client

import (
	"os"

	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/httpclient"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	svcdownloadhttpclient "github.com/dpb587/ssoca/service/download/httpclient"
	svc "github.com/dpb587/ssoca/service/env"
	svccmd "github.com/dpb587/ssoca/service/env/client/cmd"
	svchttpclient "github.com/dpb587/ssoca/service/env/httpclient"
)

type Service struct {
	svc.Service

	runtime client.Runtime
	fs      boshsys.FileSystem
}

var _ service.Service = Service{}

func NewService(runtime client.Runtime, fs boshsys.FileSystem) Service {
	return Service{
		runtime: runtime,
		fs:      fs,
	}
}

func (s Service) Description() string {
	return "Manage environment references"
}

func (s Service) GetCommand() interface{} {
	cmd := clientcmd.ServiceCommand{
		Runtime:     s.runtime,
		ServiceName: s.Type(),
	}

	return &struct {
		Set          svccmd.Set          `command:"set" description:"Configure the connection to an environment" alias:"add"`
		Info         svccmd.Info         `command:"info" description:"Show environment information"`
		Services     svccmd.Services     `command:"services" description:"Show current services available from the environment"`
		List         svccmd.List         `command:"list" description:"List all locally-configured environments"`
		SetOption    svccmd.SetOption    `command:"set-option" description:"Set a local client option in the environment"`
		UpdateClient svccmd.UpdateClient `command:"update-client" description:"Download the latest client from the environment"`
		Unset        svccmd.Unset        `command:"unset" description:"Remove all configuration for an environment" alias:"remove"`
	}{
		Set: svccmd.Set{
			ServiceCommand: cmd,
			FS:             s.fs,
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
		SetOption: svccmd.SetOption{
			ServiceCommand: cmd,
		},
		UpdateClient: svccmd.UpdateClient{
			ServiceCommand:    cmd,
			SsocaExec:         os.Args[0],
			FS:                s.fs,
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
