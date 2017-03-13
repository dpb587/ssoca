package req

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service/auth/api"
)

type Info struct{}

var _ req.RouteHandler = Info{}

func (h Info) Route() string {
	return "info"
}

func (h Info) Execute(req *http.Request) (api.InfoResponse, error) {
	res := api.InfoResponse{}

	rawToken := req.Context().Value(auth.RequestToken)

	if rawToken == nil {
		return res, nil
	}

	token, ok := rawToken.(*auth.Token)
	if !ok {
		panic("invalid token in request context")
	}

	res.ID = token.ID
	res.Groups = token.Groups
	res.Attributes = map[string]string{}

	for k, v := range token.Attributes {
		if v == nil {
			continue
		}

		res.Attributes[string(k)] = *v
	}

	return res, nil
}
