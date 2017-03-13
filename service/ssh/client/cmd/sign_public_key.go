package cmd

import (
	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/ssh/api"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type SignPublicKey struct {
	clientcmd.ServiceCommand

	GetClient GetClient
	FS        boshsys.FileSystem

	Args SignPublicKeyArgs `positional-args:"true" required:"true"`
}

var _ flags.Commander = SignPublicKey{}

type SignPublicKeyArgs struct {
	Path string `positional-arg-name:"PATH"`
}

func (c SignPublicKey) Execute(args []string) error {
	client, err := c.GetClient(c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	expandedPath, err := c.FS.ExpandPath(c.Args.Path)
	if err != nil {
		return bosherr.WrapError(err, "Expanding path")
	}

	publicKey, err := c.FS.ReadFileString(expandedPath)
	if err != nil {
		return bosherr.WrapError(err, "Reading public key")
	}

	requestPayload := api.SignPublicKeyRequest{
		PublicKey: publicKey,
	}

	response, err := client.PostSignPublicKey(requestPayload)
	if err != nil {
		return bosherr.WrapError(err, "Requesting signed public keys")
	}

	ui := c.Runtime.GetUI()
	ui.PrintLinef("%s", response.Certificate)

	return nil
}
