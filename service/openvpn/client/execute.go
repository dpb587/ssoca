package client

import (
	"fmt"
	"os"

	"github.com/dpb587/ssoca/service/openvpn/client/management"
	"github.com/dpb587/ssoca/service/openvpn/client/profile"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type ExecuteOptions struct {
	Exec      string
	ExtraArgs []string

	SkipAuthRetry     bool
	StaticCertificate bool
	Sudo              bool
}

func (s Service) Execute(opts ExecuteOptions) error {
	var executable string

	if opts.Exec != "" {
		executable = opts.Exec
	} else {
		var err error

		executable, err = s.executableFinder.Find()
		if err != nil {
			return bosherr.WrapError(err, "Finding executable")
		}
	}

	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	profileManager, err := profile.CreateManagerAndPrivateKey(client, s.name)
	if err != nil {
		return bosherr.WrapError(err, "Getting profile manager")
	}

	tmpdir, err := s.fs.TempDir("openvpn")
	if err != nil {
		return bosherr.WrapError(err, "Creating tmpdir")
	}

	defer s.fs.RemoveAll(tmpdir)

	err = s.fs.Chmod(tmpdir, 0700)
	if err != nil {
		return bosherr.WrapError(err, "Chmod'ing tmpdir")
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
			s.runtime.GetLogger(),
		)

		mgmt.Start()
		defer mgmt.Stop()
	}

	profile, err := profileManager.GetProfile()
	if err != nil {
		return bosherr.WrapError(err, "Getting profile")
	}

	if opts.StaticCertificate {
		err = s.fs.WriteFileString(configPath, profile.StaticConfig())
	} else {
		err = s.fs.WriteFileString(configPath, profile.ManagementConfig(mgmt.ManagementConfigValue()))
	}
	if err != nil {
		return bosherr.WrapError(err, "Writing certificate")
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
