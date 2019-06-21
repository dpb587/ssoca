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

func (c client) GetAuth() (api.AuthResponse, error) {
	out := api.AuthResponse{}

	err := c.client.APIGet("/env/auth", &out)
	if err != nil {
		// naive fallback to older server endpoints
		if c.client.APIGet("/auth/info", &out) != nil {
			return out, errors.Wrap(err, "getting /env/auth")
		}
	}

	return out, nil
}

func (c client) GetInfo() (api.InfoResponse, error) {
	out := api.InfoResponse{}

	err := c.client.APIGet("/env/info", &out)
	if err != nil {
		return out, errors.Wrap(err, "getting /env/info")
	}

	return out, nil
}
