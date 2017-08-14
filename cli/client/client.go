package main

import (
	"os"

	"github.com/jessevdk/go-flags"

	boshui "github.com/cloudfoundry/bosh-cli/ui"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/dpb587/ssoca/client/goflags"
	"github.com/dpb587/ssoca/client/service"
	"github.com/dpb587/ssoca/version"

	clierrors "github.com/dpb587/ssoca/cli/errors"

	srv_github_auth "github.com/dpb587/ssoca/auth/authn/github/client"
	srv_google_auth "github.com/dpb587/ssoca/auth/authn/google/client"
	srv_http_auth "github.com/dpb587/ssoca/auth/authn/http/client"
	srv_auth "github.com/dpb587/ssoca/service/auth/client"
	srv_download "github.com/dpb587/ssoca/service/download/client"
	srv_env "github.com/dpb587/ssoca/service/env/client"
	srv_openvpn "github.com/dpb587/ssoca/service/openvpn/client"
	srv_ssh "github.com/dpb587/ssoca/service/ssh/client"
	// srv_uaa_auth "github.com/dpb587/ssoca/auth/authn/uaa/client"
	// srv_uaa_auth_helper "github.com/dpb587/ssoca/auth/authn/uaa/helper"
)

var appName = "ssoca-client"
var appSemver, appCommit, appBuilt string

func main() {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)
	ui := boshui.NewConfUI(logger)
	serviceManager := service.NewDefaultManager()

	runtime := goflags.NewRuntime(version.MustVersion(appName, appSemver, appCommit, appBuilt), serviceManager, ui, os.Stdin, os.Stdout, os.Stderr, fs, logger)
	var parser = flags.NewParser(&runtime, flags.Default)

	serviceManager.Add(srv_auth.NewService(&runtime, serviceManager))
	serviceManager.Add(srv_download.NewService(&runtime, fs))
	serviceManager.Add(srv_env.NewService(&runtime, fs))
	serviceManager.Add(srv_github_auth.NewService(&runtime, cmdRunner))
	serviceManager.Add(srv_google_auth.NewService(&runtime, cmdRunner))
	serviceManager.Add(srv_http_auth.NewService(&runtime))
	serviceManager.Add(srv_openvpn.NewService(&runtime, fs, cmdRunner))
	serviceManager.Add(srv_ssh.NewService(&runtime, fs, cmdRunner))
	// serviceManager.Add(srv_uaa_auth.NewService(&runtime, srv_uaa_auth_helper.DefaultClientFactory{}))

	for _, name := range serviceManager.Services() {
		service, err := serviceManager.Get(name)
		if err != nil {
			panic(err)
		}

		command := service.GetCommand()
		if command != nil {
			parser.AddCommand(name, service.Description(), service.Description(), command)
		}
	}

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else if exitErr, ok := err.(clierrors.Exit); ok {
			os.Exit(exitErr.Code)
		} else {
			os.Exit(1)
		}
	}
}
