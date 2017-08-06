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

type agent_ struct {
	parent sshagent.Agent
	client httpclient.Client
}

var _ sshagent.Agent = agent_{}

func NewAgent(parent sshagent.Agent, client httpclient.Client) sshagent.Agent {
	return agent_{
		parent: parent,
		client: client,
	}
}

func (a agent_) List() ([]*sshagent.Key, error) {
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

		allKeys = append(allKeys, &sshagent.Key{
			Blob:    decoded,
			Comment: fmt.Sprintf("%s (ssoca agent)", key.Comment),
			Format:  split[0],
		})
	}

	return allKeys, nil
}

func (a agent_) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return a.parent.Sign(key, data)
}

func (a agent_) Add(key sshagent.AddedKey) error {
	return a.parent.Add(key)
}

func (a agent_) Remove(key ssh.PublicKey) error {
	return a.parent.Remove(key)
}

func (a agent_) RemoveAll() error {
	return a.parent.RemoveAll()
}

func (a agent_) Lock(passphrase []byte) error {
	return a.parent.Lock(passphrase)
}

func (a agent_) Unlock(passphrase []byte) error {
	return a.parent.Unlock(passphrase)
}

func (a agent_) Signers() ([]ssh.Signer, error) {
	return a.parent.Signers()
}
