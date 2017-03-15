package client

import (
	"github.com/dpb587/ssoca/auth/authn/uaa/helper"
	"github.com/dpb587/ssoca/client"

	svc "github.com/dpb587/ssoca/auth/authn/uaa"
)

type Service struct {
	svc.Service

	runtime          client.Runtime
	uaaClientFactory helper.ClientFactory
}

func NewService(runtime client.Runtime, uaaClientFactory helper.ClientFactory) Service {
	return Service{
		runtime:          runtime,
		uaaClientFactory: uaaClientFactory,
	}
}

func (s Service) Description() string {
	return "Authenticate with a Cloud Foundry UAA server"
}

func (s Service) GetCommand() interface{} {
	return nil
}
