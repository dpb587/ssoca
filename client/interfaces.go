package client

import (
	"io"

	"github.com/cloudfoundry/bosh-cli/ui"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/httpclient"
)

//go:generate counterfeiter . Runtime
type Runtime interface {
	GetEnvironment() (config.EnvironmentState, error)
	GetEnvironmentName() string
	GetConfigManager() (config.Manager, error)

	GetClient() (httpclient.Client, error)
	GetAuthInterceptClient() (httpclient.Client, error)

	GetUI() ui.UI

	GetStderr() io.Writer
	GetStdout() io.Writer
	GetStdin() io.Reader
}

//go:generate counterfeiter . ExecutableFinder
type ExecutableFinder interface {
	Find() (string, error)
}
