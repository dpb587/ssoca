package client

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/openvpn/client/profile"
)

type CreateProfileOptions struct {
	SkipAuthRetry bool
}

func (s Service) CreateProfile(opts CreateProfileOptions) (profile.Profile, error) {
	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return profile.Profile{}, errors.Wrap(err, "Getting client")
	}

	profileManager, err := profile.CreateManagerAndPrivateKey(client, s.name)
	if err != nil {
		return profile.Profile{}, errors.Wrap(err, "Getting profile manager")
	}

	return profileManager.GetProfile()
}
