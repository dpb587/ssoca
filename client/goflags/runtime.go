package goflags

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"time"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/dpb587/ssoca/client"
	"github.com/dpb587/ssoca/client/config"
	"github.com/dpb587/ssoca/client/goflags/cmd"
	authintercept "github.com/dpb587/ssoca/client/httpclient"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/config/storage"
	"github.com/dpb587/ssoca/httpclient"
	"github.com/dpb587/ssoca/version"
)

type Runtime struct {
	ConfigPath  string `long:"config" env:"SSOCA_CONFIG" description:"Configuration file path" default:"~/.config/ssoca/config"`
	Environment string `short:"e" long:"environment" env:"SSOCA_ENVIRONMENT" description:"Environment name"`

	LogLevel string `long:"log-level" env:"SSOCA_LOG_LEVEL" description:"Log level" default:"WARN"`

	Version cmd.Version `command:"version" description:"Show the current version"`

	exec           string
	version        version.Version
	serviceManager service.Manager
	fs             boshsys.FileSystem
	ui             boshui.UI
	logger         logrus.FieldLogger

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	config config.Manager
}

var _ client.Runtime = Runtime{}

func NewRuntime(exec string, version_ version.Version, serviceManager service.Manager, ui boshui.UI, stdin io.Reader, stdout io.Writer, stderr io.Writer, fs boshsys.FileSystem) *Runtime {
	runtime := Runtime{
		exec:           exec,
		version:        version_,
		serviceManager: serviceManager,
		fs:             fs,
		ui:             ui,
		stdin:          stdin,
		stdout:         stdout,
		stderr:         stderr,
	}

	runtime.Version = cmd.Version{Runtime: runtime, Version: runtime.version}

	return &runtime
}

func (r Runtime) GetLogger() logrus.FieldLogger {
	if r.logger == nil {
		logger := logrus.New()

		level, err := logrus.ParseLevel(r.LogLevel)
		if err != nil {
			panic(err)
		}

		logger.Level = level
		// logger.Formatter = &logrus.JSONFormatter{}

		r.logger = logger.WithFields(logrus.Fields{
			"cli.name":   r.Version.Version.Name,
			"cli.semver": r.Version.Version.Semver,
		})
	}

	return r.logger
}

func (r Runtime) GetExec() string {
	return r.exec
}

func (r Runtime) GetVersion() version.Version {
	return r.version
}

func (r Runtime) GetEnvironment() (config.EnvironmentState, error) {
	configManager, err := r.GetConfigManager()
	if err != nil {
		return config.EnvironmentState{}, errors.Wrap(err, "getting config manager")
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
		return nil, errors.Wrap(err, "getting environment")
	}

	var certPool *x509.CertPool

	if env.CACertificate != "" {
		certPool = x509.NewCertPool()

		cert, err := env.GetCACertificate()
		if err != nil {
			return nil, errors.Wrap(err, "getting CA certificate")
		}

		certPool.AddCert(cert)
	} else {
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, errors.Wrap(err, "loading trusted system CA certificates")
		}
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
		Timeout:   60 * time.Second,
		Transport: client.NewAuthTransport(r, r.serviceManager, baseTransport),
	}

	client := httpclient.NewClient(goclient, r.version, env.URL)

	return client, nil
}

func (r Runtime) GetAuthInterceptClient() (httpclient.Client, error) {
	client, err := r.GetClient()
	if err != nil {
		return nil, err
	}

	return authintercept.NewClient(client, r), nil
}

func (r Runtime) GetConfigManager() (config.Manager, error) {
	if r.config == nil {
		r.config = config.NewDefaultManager(storage.NewFormattedFS(r.fs, storage.YAMLFormat{}), r.ConfigPath)
	}

	return r.config, nil
}
