package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	svc "github.com/dpb587/ssoca/service/file"
)

type ServiceFactory struct {
	svc.ServiceType

	runtime   client.Runtime
	fs        boshsys.FileSystem
	cmdRunner boshsys.CmdRunner
}

func NewServiceFactory(runtime client.Runtime, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) ServiceFactory {
	return ServiceFactory{
		runtime:   runtime,
		fs:        fs,
		cmdRunner: cmdRunner,
	}
}

func (sf ServiceFactory) New(name string) service.Service {
	return NewService(name, sf.runtime, sf.fs, sf.cmdRunner)
}
