package main

import (
	"os"

	"github.com/dpb587/ssoca/certauth"
	"github.com/dpb587/ssoca/config/storage"
	"github.com/dpb587/ssoca/server"
	"github.com/dpb587/ssoca/server/config"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/dynamicvalue"

	certauth_fs "github.com/dpb587/ssoca/certauth/fs"
	certauth_memory "github.com/dpb587/ssoca/certauth/memory"

	srv_docroot "github.com/dpb587/ssoca/service/docroot/server"
	srv_env "github.com/dpb587/ssoca/service/env/server"
	srv_file "github.com/dpb587/ssoca/service/file/server"
	srv_openvpn "github.com/dpb587/ssoca/service/openvpn/server"
	srv_ssh "github.com/dpb587/ssoca/service/ssh/server"

	srv_githubauth "github.com/dpb587/ssoca/service/githubauth/server"
	srv_googleauth "github.com/dpb587/ssoca/service/googleauth/server"
	srv_httpauth "github.com/dpb587/ssoca/service/httpauth/server"
	srv_uaaauth "github.com/dpb587/ssoca/service/uaaauth/server"

	"github.com/dpb587/ssoca/auth/authz/filter"
	filter_and "github.com/dpb587/ssoca/auth/authz/filter/and"
	filter_authenticated "github.com/dpb587/ssoca/auth/authz/filter/authenticated"
	filter_or "github.com/dpb587/ssoca/auth/authz/filter/or"
	filter_public "github.com/dpb587/ssoca/auth/authz/filter/public"
	filter_remote_ip "github.com/dpb587/ssoca/auth/authz/filter/remote_ip"
	filter_scope "github.com/dpb587/ssoca/auth/authz/filter/scope"
	filter_username "github.com/dpb587/ssoca/auth/authz/filter/username"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/sirupsen/logrus"
)

func main() {
	fs := boshsys.NewOsFileSystem(boshlog.NewLogger(boshlog.LevelError))

	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout

	configStorage := storage.NewFormattedFS(fs, storage.YAMLFormat{})

	absPath, err := fs.ExpandPath(os.Args[1])
	if err != nil {
		panic(err)
	}

	cfg := config.Config{}

	err = configStorage.Get(absPath, &cfg)
	if err != nil {
		panic(err)
	}

	cfg.ApplyMigrations()
	cfg.ApplyDefaults()

	cfgval := dynamicvalue.DefaultFactory{}

	certauthManager := certauth.NewDefaultManager()

	certauthFactory := certauth.NewDefaultFactory()
	certauthFactory.Register("fs", certauth_fs.NewFactory(fs, logger))
	certauthFactory.Register("memory", certauth_memory.NewFactory(logger))

	filterManager := filter.NewDefaultManager()
	filterManager.Add("and", filter_and.NewFilter(&filterManager))
	filterManager.Add("authenticated", filter_authenticated.Filter{})
	filterManager.Add("or", filter_or.NewFilter(&filterManager))
	filterManager.Add("public", filter_public.Filter{})
	filterManager.Add("remote_ip", filter_remote_ip.Filter{})
	filterManager.Add("scope", filter_scope.Filter{})
	filterManager.Add("username", filter_username.Filter{})

	serviceManager := service.NewDefaultManager()
	serviceManager.Add(srv_env.NewService(cfg.Env, &serviceManager))

	serviceFactory := service.NewDefaultFactory()
	serviceFactory.Register(srv_githubauth.NewServiceFactory(cfg.Env.URL, cfg.Server.Redirect))
	serviceFactory.Register(srv_googleauth.NewServiceFactory(cfg.Env.URL, cfg.Server.Redirect))
	serviceFactory.Register(srv_httpauth.NewServiceFactory())
	serviceFactory.Register(srv_uaaauth.NewServiceFactory())

	serviceFactory.Register(srv_file.NewServiceFactory(fs))
	serviceFactory.Register(srv_ssh.NewServiceFactory(cfgval, certauthManager))
	serviceFactory.Register(srv_docroot.NewServiceFactory(fs))
	serviceFactory.Register(srv_openvpn.NewServiceFactory(certauthManager))

	srv, err := server.CreateFromConfig(
		cfg,
		fs,
		certauthFactory,
		serviceFactory,
		certauthManager,
		&filterManager,
		&serviceManager,
		logger,
	)
	if err != nil {
		panic(err)
	}

	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
