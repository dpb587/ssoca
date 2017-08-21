package cmd

import (
	"github.com/dpb587/ssoca/client"
	"github.com/sirupsen/logrus"
)

type ServiceCommand struct {
	Runtime client.Runtime

	ServiceName string `long:"service" short:"s" description:"Service name" env:"SSOCA_SERVICE"`
}

func (sc ServiceCommand) GetLogger() logrus.FieldLogger {
	return sc.Runtime.GetLogger().WithFields(logrus.Fields{
		"service.name": sc.ServiceName,
	})
}

type InteractiveAuthCommand struct {
	SkipAuthRetry bool `long:"skip-auth-retry" description:"Skip interactive authentication retries when logged out"`
}
