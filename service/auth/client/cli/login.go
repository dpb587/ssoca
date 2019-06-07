package cli

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/client/service"
)

type Login struct {
	*clientcmd.ServiceCommand `no-flag:"true"`

	SkipVerify bool `long:"skip-verify" description:"Skip verification of authentication, once complete"`

	ServiceManager service.Manager
	GetClient      GetClient
}

var _ flags.Commander = Login{}

func (c Login) Execute(_ []string) error {
	// TODO refactor cli + login

	// if c.SkipVerify {
	// 	return nil
	// }
	//
	// err = c.verify()
	// if err != nil {
	// 	return errors.Wrap(err, "verifying authentication")
	// }
	//
	return errors.New("TODO")
}

func (c Login) verify() error {
	ui := c.Runtime.GetUI()

	client, err := c.GetClient()
	if err != nil {
		return errors.Wrap(err, "getting client")
	}

	authInfo, err := client.GetInfo()
	if err != nil {
		return errors.Wrap(err, "getting remote authentication info")
	}

	if authInfo.ID == "" {
		return errors.New("failed to use authentication credentials")
	}

	ui.PrintBlock([]byte(fmt.Sprintf("Successfully logged in as %s\n", authInfo.ID)))

	return nil
}
