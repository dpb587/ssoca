package client

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/dpb587/ssoca/client"
)

type ServiceFactory struct {
	runtime client.Runtime
	fs      boshsys.FileSystem
}

func NewServiceFactory(runtime client.Runtime, fs boshsys.FileSystem) ServiceFactory {
	return ServiceFactory{
		runtime: runtime,
		fs:      fs,
	}
}

func (sf ServiceFactory) New(name string) Service {
	return NewService(name, sf.runtime, sf.fs)
}
