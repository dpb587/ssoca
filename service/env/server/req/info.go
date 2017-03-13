package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/env/api"
	"github.com/dpb587/ssoca/service/env/config"
)

type Info struct {
	Config   config.Config
	Services service.Manager
}

var _ req.RouteHandler = Info{}

func (h Info) Route() string {
	return "info"
}

func (h Info) Execute(req *http.Request) (api.InfoResponse, error) {
	res := api.InfoResponse{
		Env: api.InfoEnvResponse{
			Banner:   h.Config.Banner,
			Metadata: h.Config.Metadata,
			Name:     h.Config.Name,
			Title:    h.Config.Title,
			URL:      h.Config.URL,
		},
	}

	var token *auth.Token
	rawToken := req.Context().Value(auth.RequestToken)

	if rawToken != nil {
		var ok bool
		token, ok = rawToken.(*auth.Token)
		if !ok {
			panic("invalid token in request context")
		}
	}

	for _, svcName := range h.Services.Services() {
		svc, _ := h.Services.Get(svcName)

		authz, _ := svc.IsAuthorized(*req, token)
		if !authz {
			continue
		}

		svcInfo := api.InfoServiceResponse{
			Type:     svc.Type(),
			Version:  svc.Version(),
			Metadata: svc.Metadata(),
		}

		if svc.Name() == "auth" {
			res.Auth = svcInfo
		} else if svc.Name() == "env" {
			res.Version = svc.Version()
		} else {
			svcInfo.Name = svc.Name()
			res.Services = append(res.Services, svcInfo)
		}
	}

	return res, nil
}
