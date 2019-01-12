package httpclient

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/httpclient"
	"github.com/dpb587/ssoca/service/env/api"
)

func New(baseclient httpclient.Client) (Client, error) {
	if baseclient == nil {
		return nil, errors.New("client is nil")
	}

	return &client{client: baseclient}, nil
}

type client struct {
	client httpclient.Client
}

func (c client) GetInfo() (api.InfoResponse, error) {
	out := api.InfoResponse{}

	err := c.client.APIGet("/env/info", &out)
	if err != nil {
		return out, errors.Wrap(err, "Getting /env/info")
	}

	return out, nil
}
