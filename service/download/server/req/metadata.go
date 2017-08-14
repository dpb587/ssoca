package server

import (
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/download/api"
)

type Metadata struct {
	Metadata map[string]string

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Metadata{}

func (h Metadata) Route() string {
	return "metadata"
}

func (h Metadata) Execute(request req.Request) error {
	return request.WritePayload(api.MetadataResponse{
		Metadata: h.Metadata,
	})
}
