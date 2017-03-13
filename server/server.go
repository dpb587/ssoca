package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/dpb587/ssoca/authz/filter"
	"github.com/dpb587/ssoca/authz/service"
	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/server/api"
	"github.com/dpb587/ssoca/server/config"
	"github.com/dpb587/ssoca/server/service"

	srv_auth "github.com/dpb587/ssoca/service/auth/server"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Server struct {
	config   config.ServerConfig
	services service.Manager
	fs       boshsys.FileSystem
	logger   logrus.FieldLogger

	server *http.Server
}

func CreateFromConfig(
	cfg config.Config,
	fs boshsys.FileSystem,
	certauthFactory certauth.Factory,
	serviceFactory service.Factory,
	certauthManager certauth.Manager,
	filterManager filter.Manager,
	serviceManager service.Manager,
	logger logrus.FieldLogger,
) (Server, error) {
	if cfg.Auth.Type == "" {
		return Server{}, errors.New("Configuration missing: auth.type")
	}

	if cfg.Env.URL == "" {
		return Server{}, errors.New("Configuration missing: env.url")
	}

	if cfg.Server.CertificatePath == "" {
		return Server{}, errors.New("Configuration missing: server.certificate_path")
	} else if !fs.FileExists(cfg.Server.CertificatePath) {
		return Server{}, fmt.Errorf("Configuration key invalid: server.certificate_path: file does not exist: %s", cfg.Server.CertificatePath)
	}

	if cfg.Server.PrivateKeyPath == "" {
		return Server{}, errors.New("Configuration key missing: server.private_key_path:")
	} else if !fs.FileExists(cfg.Server.PrivateKeyPath) {
		return Server{}, fmt.Errorf("Configuration key invalid: server.private_key_path: file does not exist: %s", cfg.Server.PrivateKeyPath)
	}

	knownCertAuths := map[string]bool{}

	for certauthIdx, certauth := range cfg.CertAuths {
		if certauth.Name == "" {
			cfg.CertAuths[certauthIdx].Name = "default"
			certauth = cfg.CertAuths[certauthIdx]
		}

		_, found := knownCertAuths[certauth.Name]
		if found {
			return Server{}, fmt.Errorf("Configuration value duplicated: certauths[%d].name: %s", certauthIdx, certauth.Name)
		}

		knownCertAuths[certauth.Name] = true
	}

	knownServices := map[string]bool{}

	for serviceIdx, service := range cfg.Services {
		if service.Name == "" {
			cfg.Services[serviceIdx].Name = service.Type
			service = cfg.Services[serviceIdx]
		}

		_, found := knownServices[service.Name]
		if found {
			return Server{}, fmt.Errorf("Configuration value duplicated: services[%d].name: %s", serviceIdx, service.Name)
		}

		knownServices[service.Name] = true
	}

	// end validation

	for _, caConfig := range cfg.CertAuths {
		ca, err := certauthFactory.Create(caConfig.Name, caConfig.Type, caConfig.Options)

		if err != nil {
			return Server{}, bosherr.WrapErrorf(err, "Creating certauth (%s)", caConfig.Name)
		}

		certauthManager.Add(ca)
	}

	for _, svcConfig := range cfg.Services {
		svc, err := serviceFactory.Create(svcConfig.Type, svcConfig.Name, svcConfig.Options)
		if err != nil {
			return Server{}, bosherr.WrapError(err, "Creating service")
		}

		filteredService, err := filterService(svc, svcConfig, cfg.Auth.Require, filterManager)
		if err != nil {
			return Server{}, bosherr.WrapErrorf(err, "Applying authorization filters to %s", svc.Name())
		}

		serviceManager.Add(filteredService)
	}

	svc, err := serviceFactory.Create(fmt.Sprintf("%s_authn", cfg.Auth.Type), "auth", cfg.Auth.Options)
	if err != nil {
		return Server{}, bosherr.WrapError(err, "Creating auth service")
	}

	serviceManager.Add(srv_auth.NewService(svc.(service.AuthService)))

	return NewServer(cfg.Server, serviceManager, fs, logger), nil
}

func filterService(service service.Service, config config.ServiceConfig, authFilters []filter.RequireConfig, filterManager filter.Manager) (service.Service, error) {
	var merged []filter.RequireConfig

	for _, a := range authFilters {
		merged = append(merged, a)
	}

	for _, a := range config.Require {
		merged = append(merged, a)
	}

	if len(merged) == 0 {
		return service, nil
	}

	and, err := filterManager.Get("and")
	if err != nil {
		panic(err)
	}

	requirement, err := and.Create(merged)
	if err != nil {
		panic(err)
	}

	return authorized.NewService(service, requirement), nil
}

func NewServer(cfg config.ServerConfig, services service.Manager, fs boshsys.FileSystem, logger logrus.FieldLogger) Server {
	res := Server{
		config:   cfg,
		services: services,
		fs:       fs,
		logger:   logger,
	}

	return res
}

func (s Server) Run() error {
	authSvc, err := s.services.GetAuth()
	if err != nil {
		return bosherr.WrapError(err, "Loading authentication service")
	}

	mux := http.NewServeMux()

	for _, svcName := range s.services.Services() {
		svc, _ := s.services.Get(svcName)

		for _, handler := range svc.GetRoutes() {
			apiPath := fmt.Sprintf("/%s/%s", svc.Name(), handler.Route())
			apiHandler, err := api.CreateHandler(authSvc, svc, handler, s.logger)
			if err != nil {
				return bosherr.WrapErrorf(err, "Creating handler for %s", apiPath)
			}

			mux.Handle(apiPath, apiHandler)
		}
	}

	if s.config.Redirect.Root != "" {
		mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/" {
				http.NotFound(res, req)

				return
			}

			http.Redirect(res, req, s.config.Redirect.Root, http.StatusFound)
		})
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler: mux,
	}

	s.logger.WithFields(logrus.Fields{
		"server.local_addr": s.server.Addr,
	}).Info("Server is ready")

	// @todo gofunc
	return s.server.ListenAndServeTLS(s.config.CertificatePath, s.config.PrivateKeyPath)
}
