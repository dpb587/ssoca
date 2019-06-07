package client

import (
	"github.com/dpb587/ssoca/client"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	name      string
	runtime   client.Runtime
	cmdRunner boshsys.CmdRunner
}

func NewService(name string, runtime client.Runtime, cmdRunner boshsys.CmdRunner) Service {
	return Service{
		name:      name,
		runtime:   runtime,
		cmdRunner: cmdRunner,
	}
}

func (s Service) Name() string {
	return s.name
}
