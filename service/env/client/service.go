package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/httpclient"

	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/env"
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

func (s Service) Name() string {
	return "env"
}

func (s Service) GetClient() (svchttpclient.Client, error) {
	client, err := s.runtime.GetClient()
	if err != nil {
		return nil, err
	}

	return svchttpclient.New(client)
}

func (s Service) GetDownloadClient(service string, skipAuthRetry bool) (svcdownloadhttpclient.Client, error) {
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
