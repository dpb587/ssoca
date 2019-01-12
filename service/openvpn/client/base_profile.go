package client

import (
	"github.com/pkg/errors"
)

type BaseProfileOptions struct {
	SkipAuthRetry bool
}

func (s Service) BaseProfile(opts BaseProfileOptions) (string, error) {
	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return "", errors.Wrap(err, "Getting client")
	}

	return client.BaseProfile()
}
