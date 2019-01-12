package httpclient

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/httpclient"
	"github.com/dpb587/ssoca/service/ssh/api"
)

func New(baseclient httpclient.Client, service string) (Client, error) {
	if baseclient == nil {
		return nil, errors.New("client is nil")
	}

	return client{
		client:  baseclient,
		service: service,
	}, nil
}

type client struct {
	client  httpclient.Client
	service string
}

func (c client) GetCAPublicKey() (api.CAPublicKeyResponse, error) {
	out := api.CAPublicKeyResponse{}

	path := fmt.Sprintf("/%s/ca-public-key", c.service)
	err := c.client.APIGet(path, &out)
	if err != nil {
		return out, errors.Wrapf(err, "Getting %s", path)
	}

	return out, nil
}

func (c client) PostSignPublicKey(in api.SignPublicKeyRequest) (api.SignPublicKeyResponse, error) {
	out := api.SignPublicKeyResponse{}

	path := fmt.Sprintf("/%s/sign-public-key", c.service)
	err := c.client.APIPost(path, &out, in)
	if err != nil {
		return out, errors.Wrapf(err, "Posting %s", path)
	}

	return out, nil
}
