package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/httpclient"

	svc "github.com/dpb587/ssoca/service/ssh"
	svchttpclient "github.com/dpb587/ssoca/service/ssh/httpclient"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.ServiceType

	name      string
	runtime   client.Runtime
	fs        boshsys.FileSystem
	cmdRunner boshsys.CmdRunner
}

var _ service.Service = Service{}

func NewService(name string, runtime client.Runtime, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *Service {
	return &Service{
		name:      name,
		runtime:   runtime,
		fs:        fs,
		cmdRunner: cmdRunner,
	}
}

func (s Service) GetClient(skipAuthRetry bool) (svchttpclient.Client, error) {
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

	return svchttpclient.New(client, s.name)
}
