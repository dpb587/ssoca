package client

import (
	"crypto/tls"
	"crypto/x509"
	"io"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/config/storage"
	"github.com/dpb587/ssoca/httpclient"
)

type FlagsRuntime struct {
	ConfigPath  string `long:"config" env:"SSOCA_CONFIG" description:"Configuration file path" default:"~/.ssoca/config"`
	Environment string `short:"e" long:"environment" env:"SSOCA_ENVIRONMENT" description:"Environment name"`

	serviceManager service.Manager
	fs             boshsys.FileSystem
	ui             boshui.UI
	logger         boshlog.Logger
	stdout         io.Writer
	stdin          io.Reader

	config config.Manager
}

func NewFlagsRuntime(serviceManager service.Manager, ui boshui.UI, stdout io.Writer, stdin io.Reader, fs boshsys.FileSystem, logger boshlog.Logger) FlagsRuntime {
	return FlagsRuntime{
		serviceManager: serviceManager,
		fs:             fs,
		ui:             ui,
		logger:         logger,
		stdout:         stdout,
		stdin:          stdin,
	}
}

func (r FlagsRuntime) GetEnvironment() (config.EnvironmentState, error) {
	configManager, err := r.GetConfigManager()
	if err != nil {
		return config.EnvironmentState{}, bosherr.WrapError(err, "Getting config manager")
	}

	return configManager.GetEnvironment(r.GetEnvironmentName())
}

func (r FlagsRuntime) GetEnvironmentName() string {
	return r.Environment
}

func (r FlagsRuntime) GetUI() boshui.UI {
	return r.ui
}

func (r FlagsRuntime) GetStdout() io.Writer {
	return r.stdout
}

func (r FlagsRuntime) GetStdin() io.Reader {
	return r.stdin
}

func (r FlagsRuntime) GetClient() (*httpclient.Client, error) {
	env, err := r.GetEnvironment()
	if err != nil {
		return nil, bosherr.WrapError(err, "Getting environment")
	}

	certPool := x509.NewCertPool()

	if env.CACertificate != "" {
		cert, err := env.GetCACertificate()
		if err != nil {
			return nil, bosherr.WrapError(err, "Getting CA certificate")
		}

		certPool.AddCert(cert)
	}

	client := httpclient.NewClient(env.URL, &tls.Config{
		RootCAs: certPool,
	})

	client.Transport = NewAuthTransport(&r, r.serviceManager, client.Transport)

	return client, nil
}

func (r FlagsRuntime) GetConfigManager() (config.Manager, error) {
	if r.config == nil {
		r.config = config.NewDefaultManager(storage.NewFormattedFS(r.fs, storage.YAMLFormat{}), r.ConfigPath)
	}

	return r.config, nil
}
