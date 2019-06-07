package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/httpclient"
	"github.com/sirupsen/logrus"

	svc "github.com/dpb587/ssoca/service/openvpn"
	svchttpclient "github.com/dpb587/ssoca/service/openvpn/httpclient"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.ServiceType

	name                string
	runtime             client.Runtime
	logger              logrus.FieldLogger
	fs                  boshsys.FileSystem
	cmdRunner           boshsys.CmdRunner
	executableFinder    client.ExecutableFinder
	executableInstaller client.ExecutableInstaller
}

var _ service.Service = Service{}

func NewService(name string, runtime client.Runtime, logger logrus.FieldLogger, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner, executableFinder client.ExecutableFinder, executableInstaller client.ExecutableInstaller) Service {
	return Service{
		name:                name,
		runtime:             runtime,
		logger:              logger,
		fs:                  fs,
		cmdRunner:           cmdRunner,
		executableFinder:    executableFinder,
		executableInstaller: executableInstaller,
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
