package goflags

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"time"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/config"
	authintercept "github.com/dpb587/ssoca/client/httpclient"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/config/storage"
	"github.com/dpb587/ssoca/httpclient"
)

type Runtime struct {
	ConfigPath  string `long:"config" env:"SSOCA_CONFIG" description:"Configuration file path" default:"~/.config/ssoca/config"`
	Environment string `short:"e" long:"environment" env:"SSOCA_ENVIRONMENT" description:"Environment name"`

	serviceManager service.Manager
	fs             boshsys.FileSystem
	ui             boshui.UI
	logger         boshlog.Logger

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	config config.Manager
}

var _ client.Runtime = Runtime{}

func NewRuntime(serviceManager service.Manager, ui boshui.UI, stdin io.Reader, stdout io.Writer, stderr io.Writer, fs boshsys.FileSystem, logger boshlog.Logger) Runtime {
	return Runtime{
		serviceManager: serviceManager,
		fs:             fs,
		ui:             ui,
		logger:         logger,
		stdin:          stdin,
		stdout:         stdout,
		stderr:         stderr,
	}
}

func (r Runtime) GetEnvironment() (config.EnvironmentState, error) {
	configManager, err := r.GetConfigManager()
	if err != nil {
		return config.EnvironmentState{}, bosherr.WrapError(err, "Getting config manager")
	}

	return configManager.GetEnvironment(r.GetEnvironmentName())
}

func (r Runtime) GetEnvironmentName() string {
	return r.Environment
}

func (r Runtime) GetUI() boshui.UI {
	return r.ui
}

func (r Runtime) GetStdin() io.Reader {
	return r.stdin
}

func (r Runtime) GetStdout() io.Writer {
	return r.stdout
}

func (r Runtime) GetStderr() io.Writer {
	return r.stderr
}

func (r Runtime) GetClient() (httpclient.Client, error) {
	env, err := r.GetEnvironment()
	if err != nil {
		return nil, bosherr.WrapError(err, "Getting environment")
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading trusted system CA certificates")
	}

	if env.CACertificate != "" {
		cert, err := env.GetCACertificate()
		if err != nil {
			return nil, bosherr.WrapError(err, "Getting CA certificate")
		}

		certPool.AddCert(cert)
	}

	baseTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 30 * time.Second,
		DisableKeepAlives:   true,
	}

	goclient := &http.Client{
		Transport: client.NewAuthTransport(r, r.serviceManager, baseTransport),
	}

	client := httpclient.NewClient(env.URL, goclient)

	return client, nil
}

func (r Runtime) GetAuthInterceptClient() (httpclient.Client, error) {
	client, err := r.GetClient()
	if err != nil {
		return nil, err
	}

	configManager, err := r.GetConfigManager()
	if err != nil {
		return nil, err
	}

	return authintercept.NewClient(client, r.serviceManager, configManager, r.GetEnvironmentName()), nil
}

func (r Runtime) GetConfigManager() (config.Manager, error) {
	if r.config == nil {
		r.config = config.NewDefaultManager(storage.NewFormattedFS(r.fs, storage.YAMLFormat{}), r.ConfigPath)
	}

	return r.config, nil
}
