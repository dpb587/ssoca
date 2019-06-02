package client

import (
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/service/file/api"
)

type ListOptions struct {
	SkipAuthRetry bool
}

func (s Service) List(opts ListOptions) ([]api.ListFileResponse, error) {
	client, err := s.GetClient(opts.SkipAuthRetry)
	if err != nil {
		return nil, errors.Wrap(err, "getting client")
	}

	list, err := client.GetList()
	if err != nil {
		return nil, errors.Wrap(err, "getting remote environment info")
	}

	return list.Files, nil
}
