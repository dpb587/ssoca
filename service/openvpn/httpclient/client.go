package httpclient

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/dpb587/ssoca/httpclient"
	"github.com/dpb587/ssoca/service/openvpn/api"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func New(baseclient httpclient.Client, service string) (Client, error) {
	if baseclient == nil {
		return nil, errors.New("client is nil")
	}

	return &client{
		client:  baseclient,
		service: service,
	}, nil
}

type client struct {
	client  httpclient.Client
	service string
}

func (c client) BaseProfile() (string, error) {
	path := fmt.Sprintf("/%s/base-profile", c.service)
	res, err := c.client.Get(path)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Getting %s", path)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Reading response body")
	}

	return string(body), nil
}

func (c client) SignUserCSR(in api.SignUserCSRRequest) (api.SignUserCSRResponse, error) {
	out := api.SignUserCSRResponse{}

	path := fmt.Sprintf("/%s/sign-user-csr", c.service)
	err := c.client.APIPost(path, &out, in)
	if err != nil {
		return out, bosherr.WrapErrorf(err, "Posting %s", path)
	}

	return out, nil
}
