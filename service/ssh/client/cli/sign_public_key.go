package cli

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	svc "github.com/dpb587/ssoca/service/ssh/client"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type SignPublicKey struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	serviceFactory svc.ServiceFactory
	fs             boshsys.FileSystem

	Args SignPublicKeyArgs `positional-args:"true" required:"true"`
}

var _ flags.Commander = SignPublicKey{}

type SignPublicKeyArgs struct {
	Path string `positional-arg-name:"PATH"`
}

func (c SignPublicKey) Execute(_ []string) error {
	service := c.serviceFactory.New(c.ServiceName)

	expandedPath, err := c.fs.ExpandPath(c.Args.Path)
	if err != nil {
		return bosherr.WrapError(err, "Expanding path")
	}

	publicKey, err := c.fs.ReadFile(expandedPath)
	if err != nil {
		return bosherr.WrapError(err, "Reading public key")
	}

	certificate, _, err := service.SignPublicKey(svc.SignPublicKeyOptions{
		PublicKey: publicKey,
	})
	if err != nil {
		return bosherr.WrapError(err, "Getting profile")
	}

	ui := c.Runtime.GetUI()
	ui.PrintLinef("%s", certificate)

	return nil
}
