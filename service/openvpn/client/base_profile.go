package client

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type BaseProfileOptions struct {
	SkipAuthRetry bool
}

func (s Service) BaseProfile(opts BaseProfileOptions) (string, error) {
	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return "", bosherr.WrapError(err, "Getting client")
	}

	return client.BaseProfile()
}
