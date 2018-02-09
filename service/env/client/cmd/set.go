package cmd

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/jessevdk/go-flags"

	clientcmd "github.com/dpb587/ssoca/client/cmd"

	"github.com/dpb587/ssoca/client/config"
)

type Set struct {
	clientcmd.ServiceCommand

	GetClient GetClient

	Args              SetArgs `positional-args:"true"`
	CACertificatePath string  `long:"ca-cert" description:"Environment CA certificate path"`
	SkipVerify        bool    `long:"skip-verify" description:"Skip verification of environment availability"`

	FS boshsys.FileSystem
}

var _ flags.Commander = Set{}

type SetArgs struct {
	URL string `positional-arg-name:"URL" description:"Environment URL"`
}

func (c Set) Execute(_ []string) error {
	envURL := c.Args.URL

	if !strings.Contains(envURL, "://") {
		envURL = fmt.Sprintf("https://%s", envURL)
	} else if !strings.HasPrefix(envURL, "https://") {
		return fmt.Errorf("Environment URL must use https scheme: %s", envURL)
	}

	env := config.EnvironmentState{
		Alias: c.Runtime.GetEnvironmentName(),
		URL:   envURL,
	}

	if c.CACertificatePath != "" {
		absPath, err := c.FS.ExpandPath(c.CACertificatePath)
		if err != nil {
			return bosherr.WrapError(err, "Expanding path")
		}

		cacert, err := c.FS.ReadFileString(absPath)
		if err != nil {
			return bosherr.WrapError(err, "Reading file")
		}

		env.CACertificate = cacert
	}

	configManager, err := c.Runtime.GetConfigManager()
	if err != nil {
		return bosherr.WrapError(err, "Getting state manager")
	}

	err = configManager.SetEnvironment(env)
	if err != nil {
		return bosherr.WrapError(err, "Setting environment")
	}

	if c.SkipVerify {
		return nil
	}

	err = c.verify()
	if err != nil {
		return bosherr.WrapError(err, "Verifying environment")
	}

	return nil
}

func (c Set) verify() error {
	ui := c.Runtime.GetUI()

	client, err := c.GetClient()
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	info, err := client.GetInfo()
	if err != nil {
		return bosherr.WrapError(err, "Getting remote environment info")
	}

	ui.PrintBlock([]byte(fmt.Sprintf("Successfully connected to %s\n", info.Env.Title)))

	return nil
}
