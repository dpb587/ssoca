package cmd

import (
	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/service"
	"github.com/sirupsen/logrus"
)

type ServiceCommand struct {
	Runtime        client.Runtime
	ServiceManager service.Manager `no-flag:"true"`

	ServiceType string `no-flag:"true"`
	ServiceName string `long:"service" short:"s" description:"Service name" env:"SSOCA_SERVICE"`
}

func (sc ServiceCommand) GetLogger() logrus.FieldLogger {
	return sc.Runtime.GetLogger().WithFields(logrus.Fields{
		"service.type": sc.ServiceType,
		"service.name": sc.ServiceName,
	})
}

func (sc ServiceCommand) GetService() (service.Service, error) {
	return sc.ServiceManager.Get(sc.ServiceType, sc.ServiceName)
}

type InteractiveAuthCommand struct {
	SkipAuthRetry bool `long:"skip-auth-retry" description:"Skip interactive authentication retries when logged out"`
}
