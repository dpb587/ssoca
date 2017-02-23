package server

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/service/auth/api"
)

type Info struct{}

func (h Info) Route() string {
	return "info"
}

func (h Info) Execute(req *http.Request) (api.InfoResponse, error) {
	res := api.InfoResponse{}

	rawToken := req.Context().Value(auth.RequestToken)

	if rawToken == nil {
		return res, nil
	}

	token, ok := rawToken.(auth.Token)
	if !ok {
		panic("Invalid request authentication token")
	}

	res.Username = token.Username()
	res.Attributes = token.Attributes()

	return res, nil
}
