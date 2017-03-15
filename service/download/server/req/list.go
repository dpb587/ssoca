package server

import (
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/download/api"
	svcconfig "github.com/dpb587/ssoca/service/download/config"
)

type List struct {
	Paths []svcconfig.PathConfig

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = List{}

func (h List) Route() string {
	return "list"
}

func (h List) Execute(request req.Request) error {
	response := api.ListResponse{
		Files: []api.ListFileResponse{},
	}

	for _, path := range h.Paths {
		response.Files = append(
			response.Files,
			api.ListFileResponse{
				Name:   path.Name,
				Size:   path.Size,
				Digest: path.Digest,
			},
		)
	}

	return request.WritePayload(response)
}
