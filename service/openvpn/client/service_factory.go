package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	svc "github.com/dpb587/ssoca/service/openvpn"
	"github.com/sirupsen/logrus"
)

type ServiceFactory struct {
	svc.ServiceType

	runtime          client.Runtime
	fs               boshsys.FileSystem
	cmdRunner        boshsys.CmdRunner
	executableFinder client.ExecutableFinder
}

func NewServiceFactory(runtime client.Runtime, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner, executableFinder client.ExecutableFinder) ServiceFactory {
	return ServiceFactory{
		runtime:          runtime,
		fs:               fs,
		cmdRunner:        cmdRunner,
		executableFinder: executableFinder,
	}
}

func (sf ServiceFactory) New(name string) *Service {
	return NewService(
		name,
		sf.runtime,
		sf.runtime.GetLogger().WithFields(logrus.Fields{
			"service.type": sf.Type(),
			"service.name": name,
		}),
		sf.fs,
		sf.cmdRunner,
		sf.executableFinder,
	)
}
