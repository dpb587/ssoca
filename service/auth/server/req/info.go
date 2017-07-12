package req

import (
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/auth/api"
)

type Info struct {
	req.WithoutAdditionalAuthorization
}

var _ req.RouteHandler = Info{}

func (h Info) Route() string {
	return "info"
}

func (h Info) Execute(request req.Request) error {
	response := api.InfoResponse{}

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
