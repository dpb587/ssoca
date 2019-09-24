package client

import (
	"fmt"
	"os"
	"path"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/openvpn/client/internal"
	"github.com/dpb587/ssoca/service/openvpn/client/management"
	"github.com/dpb587/ssoca/service/openvpn/client/profile"
)

type ExecuteOptions struct {
	Exec      string
	ExtraArgs []string

	ManagementMode executeManagementMode
	SkipAuthRetry  bool
	Sudo           bool
}

type executeManagementMode string

const ExecuteManagementModeAuto executeManagementMode = "auto"
const ExecuteManagementModeDisabled executeManagementMode = "disabled"
const ExecuteManagementModeEnabled executeManagementMode = "enabled"

func ExecuteManagementModeFromString(in string) (executeManagementMode, error) {
	switch executeManagementMode(in) {
	case ExecuteManagementModeAuto:
		return ExecuteManagementModeAuto, nil
	case ExecuteManagementModeDisabled:
		return ExecuteManagementModeDisabled, nil
	case ExecuteManagementModeEnabled:
		return ExecuteManagementModeEnabled, nil
	}

	return "", fmt.Errorf("invalid value: %v", in)
}

func (s Service) Execute(opts ExecuteOptions) error {
	var executable string

	if opts.Exec != "" {
		executable = opts.Exec
	} else {
		var err error
		var guessed bool

		executable, guessed, err = s.executableFinder.Find()
		if err != nil {
			return errors.Wrap(err, "finding executable")
		} else if guessed {
			s.logger.Warnf("openvpn executable found outside of $PATH (using %s)", executable)
		}
	}

	if opts.ManagementMode == "" {
		opts.ManagementMode = ExecuteManagementModeAuto
	}

	if opts.ManagementMode == ExecuteManagementModeAuto {
		// this is best effort; ignore errors in favor of surfacing them in non-magic paths
		stdout, _, _, _ := s.cmdRunner.RunCommand(executable, "--version")
		if strings.Contains(stdout, "library versions: OpenSSL 1.1.1") {
			s.logger.Warnf("disabling management interface (https://github.com/dpb587/ssoca/issues/13): detected openssl 1.1.1")

			opts.ManagementMode = ExecuteManagementModeDisabled
		}
	}

	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return errors.Wrap(err, "getting client")
	}

	profileManager, err := profile.CreateManagerAndPrivateKey(client, s.name)
	if err != nil {
		return errors.Wrap(err, "getting profile manager")
	}

	tmpdir, err := s.fs.TempDir("openvpn")
	if err != nil {
		return errors.Wrap(err, "creating tmpdir")
	}

	defer s.fs.RemoveAll(tmpdir)

	err = s.fs.Chmod(tmpdir, 0700)
	if err != nil {
		return errors.Wrap(err, "chmod'ing tmpdir")
	}

	configPath := fmt.Sprintf("%s/openvpn.ovpn", tmpdir)

	openvpnargs := []string{}

	if opts.Sudo {
		openvpnargs = append(openvpnargs, executable)
		executable = "sudo"
	}

	openvpnargs = append(openvpnargs, "--config", configPath)
	openvpnargs = append(openvpnargs, opts.ExtraArgs...)

	var mgmt management.Server

	if opts.ManagementMode == ExecuteManagementModeEnabled {
		mgmt = management.NewServer(
			management.NewDefaultHandler(&profileManager),
			"tcp",
			"127.0.0.1:0",
			internal.GeneratePassword(32),
			s.runtime.GetLogger(),
		)

		mgmt.Start()
		defer mgmt.Stop()
	}

	profile, err := profileManager.GetProfile()
	if err != nil {
		return errors.Wrap(err, "getting profile")
	}

	if opts.ManagementMode == ExecuteManagementModeEnabled {
		managementPasswordPath := path.Join(tmpdir, "management.pw")

		err = s.fs.WriteFileString(managementPasswordPath, mgmt.ManagementPassword()+"\n")
		if err != nil {
			return errors.Wrap(err, "writing management password file")
		}

		// the containing directory already has restricted permissions;
		// this avoids a warning message from openvpn
		err = s.fs.Chmod(managementPasswordPath, 0700)
		if err != nil {
			return errors.Wrap(err, "chmod'ing management password file")
		}

		err = s.fs.WriteFileString(configPath, profile.ManagementConfig(mgmt.ManagementConfigValue(), managementPasswordPath))
	} else {
		err = s.fs.WriteFileString(configPath, profile.StaticConfig())
	}
	if err != nil {
		return errors.Wrap(err, "writing profile")
	}

	_, _, _, err = s.cmdRunner.RunComplexCommand(executeRewriteCommand(boshsys.Command{
		Name: executable,
		Args: openvpnargs,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		KeepAttached: true,
	}))

	return err
}
