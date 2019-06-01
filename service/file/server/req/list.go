package server

import (
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/file/api"
	svcconfig "github.com/dpb587/ssoca/service/file/server/config"
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
				Name: path.Name,
				Size: path.Size,
				Digest: api.ListFileDigestResponse{
					SHA1:   path.Digest.SHA1,
					SHA256: path.Digest.SHA256,
					SHA512: path.Digest.SHA512,
				},
			},
		)
	}

	return request.WritePayload(response)
}
