package cmd

import "github.com/dpb587/ssoca/client"

type ServiceCommand struct {
	Runtime client.Runtime

	ServiceName string `long:"service" short:"s" description:"Service name" env:"SSOCA_SERVICE"`
}

type InteractiveAuthCommand struct {
	SkipAuthRetry bool `long:"skip-auth-retry" description:"Skip interactive authentication retries when logged out"`
}
