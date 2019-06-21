package req

import (
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/env/api"
)

type Auth struct {
	req.WithoutAdditionalAuthorization

	RouteName string
}

var _ req.RouteHandler = Auth{}

func (h Auth) Route() string {
	if h.RouteName != "" {
		// TODO deprecate after legacy auth service no longer uses this
		return h.RouteName
	}

	return "auth"
}

func (h Auth) Execute(request req.Request) error {
	response := api.AuthResponse{}

	token := request.AuthToken

	if token != nil {
		response.ID = token.ID
		response.Groups = token.Groups
		response.Attributes = map[string]string{}

		for k, v := range token.Attributes {
			if v == nil {
				continue
			}

			response.Attributes[string(k)] = *v
		}
	}

	return request.WritePayload(response)
}
