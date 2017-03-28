package agent

import (
	"encoding/base64"
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/service/ssh/api"
	"github.com/dpb587/ssoca/service/ssh/httpclient"
	"golang.org/x/crypto/ssh"
	sshagent "golang.org/x/crypto/ssh/agent"
)

type Agent struct {
	parent sshagent.Agent
	client *httpclient.Client
}

var _ sshagent.Agent = Agent{}

func NewAgent(parent sshagent.Agent, client *httpclient.Client) sshagent.Agent {
	return Agent{
		parent: parent,
		client: client,
	}
}

func (a Agent) List() ([]*sshagent.Key, error) {
	keys, err := a.parent.List()
	if err != nil {
		return nil, err
	}

	for keyIdx, key := range keys {
		cert, err := a.client.PostSignPublicKey(api.SignPublicKeyRequest{
			PublicKey: fmt.Sprintf("%s %s", key.Format, base64.StdEncoding.EncodeToString(key.Blob)),
		})
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "signing public key %d", keyIdx+1)
		}

		split := strings.SplitN(cert.Certificate, " ", 2)
		if len(split) != 2 {
			return nil, fmt.Errorf("signing public key: unexpected server response: %s", cert.Certificate)
		}

		decoded, err := base64.StdEncoding.DecodeString(split[1])
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "decoding certificate")
		}

		keys[keyIdx] = &sshagent.Key{
			Blob:    decoded,
			Comment: key.Comment,
			Format:  split[0],
		}

		break
	}

	return keys, nil
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
