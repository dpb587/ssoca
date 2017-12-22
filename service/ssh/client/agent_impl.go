package client

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"
)

type Agent struct {
	parent  sshagent.Agent
	service Service
}

var _ sshagent.Agent = Agent{}

func (a Agent) List() ([]*sshagent.Key, error) {
	keys, err := a.parent.List()
	if err != nil {
		return nil, err
	}

	allKeys := []*sshagent.Key{}

	// @todo probably don't want to sign everything, everytime
	for keyIdx, key := range keys {
		parsedKey, err := ssh.ParsePublicKey(key.Blob)
		if err != nil {
			return nil, bosherr.WrapError(err, "parsing public key")
		}

		allKeys = append(allKeys, key)

		_, signed := parsedKey.(*ssh.Certificate)
		if signed {
			continue
		}

		certificate, _, err := a.service.SignPublicKey(SignPublicKeyOptions{
			PublicKey: []byte(fmt.Sprintf("%s %s", parsedKey.Type(), base64.StdEncoding.EncodeToString(parsedKey.Marshal()))),
		})
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "signing public key %d", keyIdx+1)
		}

		split := strings.SplitN(string(certificate), " ", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("signing public key: unexpected server response: %s", certificate)
		}

		decoded, err := base64.StdEncoding.DecodeString(split[1])
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "decoding certificate")
		}

		allKeys = append(allKeys, &sshagent.Key{
			Blob:    decoded,
			Comment: fmt.Sprintf("%s (ssoca agent)", key.Comment),
			Format:  split[0],
		})
	}

	return allKeys, nil
}

func (a Agent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return a.parent.Sign(key, data)
}

func (a Agent) Add(key sshagent.AddedKey) error {
	return a.parent.Add(key)
}

func (a Agent) Remove(key ssh.PublicKey) error {
	return a.parent.Remove(key)
}

func (a Agent) RemoveAll() error {
	return a.parent.RemoveAll()
}

func (a Agent) Lock(passphrase []byte) error {
	return a.parent.Lock(passphrase)
}

func (a Agent) Unlock(passphrase []byte) error {
	return a.parent.Unlock(passphrase)
}

func (a Agent) Signers() ([]ssh.Signer, error) {
	return a.parent.Signers()
}

func (a Agent) Listen(socket net.Listener) error {
	for {
		handle, err := socket.Accept()
		if err != nil {
			return err
		}

		go func(handle net.Conn) {
			defer handle.Close()

			_ = sshagent.ServeAgent(a, handle)
		}(handle)
	}
}
