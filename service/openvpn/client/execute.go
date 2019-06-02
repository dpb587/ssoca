package client

import (
	"fmt"
	"os"
	"path"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/openvpn/client/internal"
	"github.com/dpb587/ssoca/service/openvpn/client/management"
	"github.com/dpb587/ssoca/service/openvpn/client/profile"
)

type ExecuteOptions struct {
	Exec      string
	ExtraArgs []string

	SkipAuthRetry     bool
	SkipInstall       bool
	StaticCertificate bool
	Sudo              bool
}

func (s Service) requireExecutable(skipInstall bool) (string, error) {
	executable, guessed, err := s.executableFinder.Find()
	if err != nil {
		if skipInstall {
			return "", errors.Wrap(err, "finding executable")
		}
	}

	if guessed {
		s.logger.Warnf("openvpn executable found outside of $PATH (using %s)", executable)
	}

	if executable != "" {
		return executable, nil
	}

	s.logger.Warnf("openvpn executable not found (attempting automatic installation)")

	err = s.executableInstaller.Install(s.logger)
	if err != nil {
		return "", errors.Wrap(err, "installing executable")
	}

	return s.requireExecutable(true)
}

func (s Service) Execute(opts ExecuteOptions) error {
	var executable string

	if opts.Exec != "" {
		executable = opts.Exec
	} else {
		var err error

		executable, err = s.requireExecutable(opts.SkipInstall)
		if err != nil {
			return errors.Wrap(err, "requiring executable")
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

	if !opts.StaticCertificate {
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

	if opts.StaticCertificate {
		err = s.fs.WriteFileString(configPath, profile.StaticConfig())
	} else {
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
	}
	if err != nil {
		return errors.Wrap(err, "writing certificate")
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
