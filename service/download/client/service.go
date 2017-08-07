package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/httpclient"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/download"
	svccmd "github.com/dpb587/ssoca/service/download/client/cmd"
	svchttpclient "github.com/dpb587/ssoca/service/download/httpclient"
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
	return "Download environment artifacts"
}

func (s Service) GetCommand() interface{} {
	cmd := clientcmd.ServiceCommand{
		Runtime:     s.runtime,
		ServiceName: s.Type(),
	}

	return &struct {
		Get  svccmd.Get  `command:"get" description:"Get an artifact"`
		List svccmd.List `command:"list" description:"List available artifacts"`
	}{
		Get: svccmd.Get{
			ServiceCommand: cmd,
			FS:             s.fs,
			GetClient:      s.GetClient,
		},
		List: svccmd.List{
			ServiceCommand: cmd,
			GetClient:      s.GetClient,
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
