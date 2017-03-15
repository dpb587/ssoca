package req

import (
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/env/api"
	"github.com/dpb587/ssoca/service/env/config"
)

type Info struct {
	Config   config.Config
	Services service.Manager

	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Info{}

func (h Info) Route() string {
	return "info"
}

func (h Info) Execute(request req.Request) error {
	response := api.InfoResponse{
		Env: api.InfoEnvResponse{
			Banner:   h.Config.Banner,
			Metadata: h.Config.Metadata,
			Name:     h.Config.Name,
			Title:    h.Config.Title,
			URL:      h.Config.URL,
		},
	}

	for _, svcName := range h.Services.Services() {
		svc, _ := h.Services.Get(svcName)

		authz, _ := svc.IsAuthorized(*request.RawRequest, request.AuthToken)
		if !authz {
			continue
		}

		svcInfo := api.InfoServiceResponse{
			Type:     svc.Type(),
			Version:  svc.Version(),
			Metadata: svc.Metadata(),
		}

		if svc.Name() == "auth" {
			response.Auth = svcInfo
		} else if svc.Name() == "env" {
			response.Version = svc.Version()
		} else {
			svcInfo.Name = svc.Name()
			response.Services = append(response.Services, svcInfo)
		}
	}

	return request.WritePayload(response)
}
