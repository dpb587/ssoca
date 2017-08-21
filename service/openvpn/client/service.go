package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/httpclient"

	svc "github.com/dpb587/ssoca/service/openvpn"
	svchttpclient "github.com/dpb587/ssoca/service/openvpn/httpclient"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.Service

	runtime          client.Runtime
	fs               boshsys.FileSystem
	cmdRunner        boshsys.CmdRunner
	executableFinder client.ExecutableFinder
}

var _ service.Service = Service{}

func NewService(runtime client.Runtime, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner, executableFinder client.ExecutableFinder) Service {
	return Service{
		runtime:          runtime,
		fs:               fs,
		cmdRunner:        cmdRunner,
		executableFinder: executableFinder,
	}
}

func (s Service) Description() string {
	return "Establish an OpenVPN to a remote server"
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
