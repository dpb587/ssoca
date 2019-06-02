package client

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/service/uaaauth/helper"

	svc "github.com/dpb587/ssoca/service/uaaauth"
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
