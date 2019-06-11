package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/httpclient"

	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/file"
	svchttpclient "github.com/dpb587/ssoca/service/file/httpclient"
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

func (s Service) Name() string {
	return s.name
}

func (s Service) Description() string {
	return "Download environment artifacts"
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
