package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"

	clientcmd "github.com/dpb587/ssoca/client/cmd"
	"github.com/dpb587/ssoca/service/ssh/api"
	"github.com/jessevdk/go-flags"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Exec struct {
	clientcmd.ServiceCommand
	clientcmd.InteractiveAuthCommand

	GetClient GetClient
	CmdRunner boshsys.CmdRunner
	FS        boshsys.FileSystem

	Opts []string `long:"opts" description:"Options to pass through to SSH"`

	Args ExecArgs `positional-args:"true" optional:"true"`
}

var _ flags.Commander = Exec{}

type ExecArgs struct {
	Host string `positional-arg-name:"HOST"`
}

func (c Exec) Execute(_ []string) error {
	client, err := c.GetClient(c.ServiceName, c.SkipAuthRetry)
	if err != nil {
		return bosherr.WrapError(err, "Getting client")
	}

	tmpdir, err := c.FS.TempDir("ssh")
	if err != nil {
		return bosherr.WrapError(err, "Creating certificate tmpdir")
	}

	defer c.FS.RemoveAll(tmpdir)

	privateKeyBytes, publicKeyBytes, err := makeSSHKeyPair()
	if err != nil {
		return bosherr.WrapError(err, "Creating ephemeral ssh key")
	}

	tmpPrivateKey := fmt.Sprintf("%s/id_rsa", tmpdir)

	err = c.FS.WriteFile(tmpPrivateKey, privateKeyBytes)
	if err != nil {
		return bosherr.WrapError(err, "Writing private key")
	}

	err = c.FS.Chmod(tmpPrivateKey, 0600)
	if err != nil {
		return bosherr.WrapError(err, "Setting permissions of private key")
	}

	err = c.FS.WriteFile(fmt.Sprintf("%s/id_rsa.pub", tmpdir), publicKeyBytes)
	if err != nil {
		return bosherr.WrapError(err, "Writing public key")
	}

	response, err := client.PostSignPublicKey(api.SignPublicKeyRequest{
		PublicKey: string(publicKeyBytes),
	})
	if err != nil {
		return bosherr.WrapError(err, "Requesting signed public keys")
	}

	sshargs := []string{
		"-o", "ForwardAgent=no",
		"-o", "ServerAliveInterval=30",
		"-o", "IdentitiesOnly=yes",
	}

	tmpCertificate := fmt.Sprintf("%s/id_rsa-cert.pub", tmpdir)

	err = c.FS.WriteFileString(tmpCertificate, response.Certificate)
	if err != nil {
		return bosherr.WrapError(err, "Writing certificate")
	}

	sshargs = append(sshargs, "-o", fmt.Sprintf("IdentityFile=%s", tmpPrivateKey))
	sshargs = append(sshargs, "-o", fmt.Sprintf("CertificateFile=%s", tmpCertificate))

	for _, arg := range c.Opts {
		sshargs = append(sshargs, arg)
	}

	if response.Target != nil {
		if response.Target.Port != 0 {
			sshargs = append(sshargs, "-p", string(response.Target.Port))
		}

		if response.Target.User != "" {
			sshargs = append(sshargs, "-l", response.Target.User)
		}

		if response.Target.PublicKey != "" {
			sshargs = append(sshargs, "-o", "StrictHostKeyChecking=yes")

			tmpKnownHosts := fmt.Sprintf("%s/known_hosts", tmpdir)

			err = c.FS.WriteFileString(tmpKnownHosts, fmt.Sprintf("%s %s\n", response.Target.Host, response.Target.PublicKey))
			if err != nil {
				return bosherr.WrapError(err, "Writing certificate")
			}

			sshargs = append(sshargs, "-o", fmt.Sprintf("UserKnownHostsFile=%s", tmpKnownHosts))
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

// https://github.com/cloudfoundry/bosh-cli/blob/a0c78a59b5eeac11a32e953451a497eb1cb9ba7d/director/ssh_opts.go#L43
func makeSSHKeyPair() ([]byte, []byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}

	privKeyBuf := bytes.NewBufferString("")

	err = pem.Encode(privKeyBuf, privKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	pub, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privKeyBuf.Bytes(), ssh.MarshalAuthorizedKey(pub), nil
}
