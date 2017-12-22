package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"golang.org/x/crypto/ssh"
)

type ExecuteOptions struct {
	Exec      string
	ExtraArgs []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Host string

	SkipAuthRetry bool
}

func (s Service) Execute(opts ExecuteOptions) (int, error) {
	if err := s.executeOpts(&opts); err != nil {
		return -1, err
	}

	tmpdir, err := s.fs.TempDir("ssh")
	if err != nil {
		return -1, bosherr.WrapError(err, "Creating certificate tmpdir")
	}

	defer s.fs.RemoveAll(tmpdir)

	privateKeyBytes, publicKeyBytes, err := makeSSHKeyPair()
	if err != nil {
		return -1, bosherr.WrapError(err, "Creating ephemeral ssh key")
	}

	tmpPrivateKey := fmt.Sprintf("%s/id_rsa", tmpdir)

	err = s.fs.WriteFile(tmpPrivateKey, nil)
	if err != nil {
		return -1, bosherr.WrapError(err, "Touching private key")
	}

	err = s.fs.Chmod(tmpPrivateKey, 0600)
	if err != nil {
		return -1, bosherr.WrapError(err, "Setting permissions of private key")
	}

	err = s.fs.WriteFile(tmpPrivateKey, privateKeyBytes)
	if err != nil {
		return -1, bosherr.WrapError(err, "Writing private key")
	}

	err = s.fs.WriteFile(fmt.Sprintf("%s/id_rsa.pub", tmpdir), publicKeyBytes)
	if err != nil {
		return -1, bosherr.WrapError(err, "Writing public key")
	}

	certificate, target, err := s.SignPublicKey(SignPublicKeyOptions{
		PublicKey: publicKeyBytes,
	})
	if err != nil {
		return -1, bosherr.WrapError(err, "Requesting signed public keys")
	}

	sshargs := []string{
		"-o", "ForwardAgent=no",
		"-o", "ServerAliveInterval=30",
		"-o", "IdentitiesOnly=yes",
	}

	tmpCertificate := fmt.Sprintf("%s/id_rsa-cert.pub", tmpdir)

	err = s.fs.WriteFile(tmpCertificate, certificate)
	if err != nil {
		return -1, bosherr.WrapError(err, "Writing certificate")
	}

	sshargs = append(sshargs, "-o", fmt.Sprintf("IdentityFile=%s", tmpPrivateKey))
	sshargs = append(sshargs, "-o", fmt.Sprintf("CertificateFile=%s", tmpCertificate))
	sshargs = append(sshargs, opts.ExtraArgs...)

	if target != nil {
		if target.Port != 0 {
			sshargs = append(sshargs, "-p", string(target.Port))
		}

		if target.User != "" {
			sshargs = append(sshargs, "-l", target.User)
		}

		if target.PublicKey != "" {
			sshargs = append(sshargs, "-o", "StrictHostKeyChecking=yes")

			tmpKnownHosts := fmt.Sprintf("%s/known_hosts", tmpdir)

			err = s.fs.WriteFileString(tmpKnownHosts, fmt.Sprintf("%s %s\n", target.Host, target.PublicKey))
			if err != nil {
				return -1, bosherr.WrapError(err, "Writing certificate")
			}

			sshargs = append(sshargs, "-o", fmt.Sprintf("UserKnownHostsFile=%s", tmpKnownHosts))
		}

		if target.Host != "" {
			if opts.Host != "" {
				return -1, errors.New("Cannot specify user or host (already configured by remote)")
			}

			sshargs = append(sshargs, target.Host)
		}
	}

	if opts.Host != "" {
		sshargs = append(sshargs, opts.Host)
	}

	_, _, exitStatus, err := s.cmdRunner.RunComplexCommand(boshsys.Command{
		Name: opts.Exec,
		Args: sshargs,

		Stdin:  opts.Stdin,
		Stdout: opts.Stdout,
		Stderr: opts.Stderr,

		KeepAttached: true,
	})

	return exitStatus, err
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

func (s Service) executeOpts(opts *ExecuteOptions) error {
	if opts.Exec == "" {
		opts.Exec = "ssh"
	}

	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}

	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}

	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}

	return nil
}
