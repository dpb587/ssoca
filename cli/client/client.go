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

	srv_auth "github.com/dpb587/ssoca/service/auth/client"
	srv_auth_cli "github.com/dpb587/ssoca/service/auth/client/cli"
	srv_env "github.com/dpb587/ssoca/service/env/client"
	srv_env_cli "github.com/dpb587/ssoca/service/env/client/cli"
	srv_file "github.com/dpb587/ssoca/service/file/client"
	srv_file_cli "github.com/dpb587/ssoca/service/file/client/cli"
	srv_githubauth "github.com/dpb587/ssoca/service/githubauth/client"
	srv_googleauth "github.com/dpb587/ssoca/service/googleauth/client"
	srv_httpauth "github.com/dpb587/ssoca/service/httpauth/client"
	srv_openvpn "github.com/dpb587/ssoca/service/openvpn/client"
	srv_openvpn_cli "github.com/dpb587/ssoca/service/openvpn/client/cli"
	srv_openvpn_executable_finder "github.com/dpb587/ssoca/service/openvpn/client/executable/finder"
	srv_ssh "github.com/dpb587/ssoca/service/ssh/client"
	srv_ssh_cli "github.com/dpb587/ssoca/service/ssh/client/cli"
	srv_uaaauth "github.com/dpb587/ssoca/service/uaaauth/client"
	srv_uaaauth_helper "github.com/dpb587/ssoca/service/uaaauth/helper"
)

var appName = "ssoca-client"
var appSemver, appCommit, appBuilt string

func main() {
	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)
	ui := boshui.NewConfUI(logger)
	serviceManager := service.NewDefaultManager()

	runtime := goflags.NewRuntime(os.Args[0], version.MustVersion(appName, appSemver, appCommit, appBuilt), serviceManager, ui, os.Stdin, os.Stdout, os.Stderr, fs)
	var parser = flags.NewParser(runtime, flags.Default)

	envService := srv_env.NewService(runtime, fs, cmdRunner)

	serviceManager.Add(envService)

	serviceManager.AddFactory(srv_file.NewServiceFactory(runtime, fs, cmdRunner))
	serviceManager.AddFactory(srv_githubauth.NewServiceFactory(runtime, cmdRunner))
	serviceManager.AddFactory(srv_googleauth.NewServiceFactory(runtime, cmdRunner))
	serviceManager.AddFactory(srv_httpauth.NewServiceFactory(runtime))
	serviceManager.AddFactory(srv_uaaauth.NewServiceFactory(runtime, serviceManager, srv_uaaauth_helper.DefaultClientFactory{}))

	parser.AddCommand(
		"env",
		"Manage environment references",
		"Manage environment references",
		srv_env_cli.CreateCommands(runtime, cmdRunner, fs, envService),
	)

	parser.AddCommand(
		"auth",
		"Manage authentication",
		"Manage authentication",
		srv_auth_cli.CreateCommands(runtime, srv_auth.NewService(runtime, serviceManager)),
	)

	parser.AddCommand(
		"file",
		"Access files from the environment",
		"Access files from the environment",
		srv_file_cli.CreateCommands(runtime, serviceManager),
	)

	parser.Find("file").Aliases = []string{"download"}

	openvpnFinder := srv_openvpn_executable_finder.Finder{FS: fs}
	parser.AddCommand(
		"openvpn",
		"Establish OpenVPN connections to remote servers",
		"Establish OpenVPN connections to remote servers",
		srv_openvpn_cli.CreateCommands(runtime, srv_openvpn.NewServiceFactory(runtime, fs, cmdRunner, openvpnFinder), fs, cmdRunner),
	)

	parser.AddCommand(
		"ssh",
		"Establish SSH connections to remote servers",
		"Establish SSH connections to remote servers",
		srv_ssh_cli.CreateCommands(runtime, srv_ssh.NewServiceFactory(runtime, fs, cmdRunner), fs, cmdRunner),
	)

	// execute

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		runtime.GetLogger().Fatal(err)

		if exitErr, ok := err.(clierrors.Exit); ok {
			os.Exit(exitErr.Code)
		} else {
			os.Exit(1)
		}
	}
}
