package client

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/ssh/api"
)

type SignPublicKeyOptions struct {
	SkipAuthRetry bool
	PublicKey     []byte
}

func (s Service) SignPublicKey(opts SignPublicKeyOptions) ([]byte, *api.SignPublicKeyTargetResponse, error) {
	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Getting client")
	}

	requestPayload := api.SignPublicKeyRequest{
		PublicKey: string(opts.PublicKey),
	}

	response, err := client.PostSignPublicKey(requestPayload)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Requesting signed public keys")
	}

	return []byte(response.Certificate), response.Target, err
}
