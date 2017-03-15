package client

import (
	"github.com/dpb587/ssoca/client"

	svc "github.com/dpb587/ssoca/auth/authn/http"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Service struct {
	svc.Service

	runtime client.Runtime
	fs      boshsys.FileSystem
}

func NewService(runtime client.Runtime) Service {
	return Service{
		runtime: runtime,
	}
}

func (s Service) Description() string {
	return "Authenticate with HTTP username/password"
}

func (s Service) GetCommand() interface{} {
	return nil
}
