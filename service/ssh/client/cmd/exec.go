package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh/agent"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/ssh/api"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Exec struct {
	clientcmd.ServiceCommand

	GetClient GetClient
	CmdRunner boshsys.CmdRunner
	FS        boshsys.FileSystem

	Args ExecArgs `positional-args:"true"`
}

type ExecArgs struct {
	Host string `positional-arg-name:"HOST"`
}

func (c *Exec) Execute(args []string) error {
	client, err := c.GetClient(c.ServiceName)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	requestPayload, err := c.generatePayload()
	if err != nil {
		return bosherr.WrapError(err, "Generating request payload")
	}

	response, err := client.PostSignPublicKey(requestPayload)
	if err != nil {
		return bosherr.WrapError(err, "Requesting signed public keys")
	}

	tmpdir, err := c.FS.TempDir("ssh")
	if err != nil {
		return bosherr.WrapError(err, "Creating certificate tmpdir")
	}

	defer c.FS.RemoveAll(tmpdir)

	sshargs := []string{}

	targetPath := fmt.Sprintf("%s/id_rsa-cert.pub", tmpdir)

	err = c.FS.WriteFileString(targetPath, response.Certificate)
	if err != nil {
		return bosherr.WrapError(err, "Writing certificate")
	}

	sshargs = append(sshargs, "-o", fmt.Sprintf("CertificateFile=%s", targetPath))

	if response.Target != nil {
		if response.Target.Port != 0 {
			sshargs = append(sshargs, "-p", string(response.Target.Port))
		}

		if response.Target.User != "" {
			sshargs = append(sshargs, "-l", response.Target.User)
		}

		if response.Target.Host != "" {
			if c.Args.Host != "" {
				return errors.New("Cannot specify user or host (already configured by remote)")
			}

			sshargs = append(sshargs, response.Target.Host)
		}
	}

	if c.Args.Host != "" {
		sshargs = append(sshargs, c.Args.Host)
	}

	_, _, _, err = c.CmdRunner.RunComplexCommand(boshsys.Command{
		Name: "ssh",
		Args: sshargs,

		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		KeepAttached: true,
	})

	return err
}

func (c *Exec) generatePayload() (api.SignPublicKeyRequest, error) {
	payload := api.SignPublicKeyRequest{}

	socket, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return payload, bosherr.WrapError(err, "Dialing $SSH_AUTH_SOCK")
	}

	a := agent.NewClient(socket)
	keys, err := a.List()
	if err != nil {
		return payload, bosherr.WrapError(err, "Listing SSH agent keys")
	}

	for _, key := range keys {
		payload.PublicKey = fmt.Sprintf("%s %s %s", key.Format, base64.StdEncoding.EncodeToString(key.Marshal()), key.Comment)

		break
	}

	if payload.PublicKey == "" {
		return payload, errors.New("Cannot find public key from agent")
	}

	return payload, nil
}
