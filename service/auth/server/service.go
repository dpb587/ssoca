package server

import (
	"net/http"

	svc "github.com/dpb587/ssoca/service/auth"
	svcreq "github.com/dpb587/ssoca/service/auth/server/req"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
)

type Service struct {
	svc.Service

	auth service.AuthService
}

func NewService(authService service.AuthService) Service {
	return Service{
		auth: authService,
	}
}

func (s Service) Name() string {
	return "auth"
}

func (s Service) Type() string {
	return s.auth.Type()
}

func (s Service) Metadata() interface{} {
	return s.auth.Metadata()
}

func (s Service) GetRoutes() []req.RouteHandler {
	return append(s.auth.GetRoutes(), svcreq.Info{})
}

func (s Service) ParseRequestAuth(r http.Request) (auth.Token, error) {
	return s.auth.ParseRequestAuth(r)
}

func (s Service) IsAuthorized(req http.Request, token auth.Token) (bool, error) {
	return s.auth.IsAuthorized(req, token)
}
