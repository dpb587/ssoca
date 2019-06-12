package authorized

import (
	"net/http"

	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/auth/authz/filter"
	serverservice "github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/dpb587/ssoca/service"
)

type Service struct {
	service     serverservice.Service
	requirement filter.Requirement
}

func NewService(service serverservice.Service, requirement filter.Requirement) Service {
	return Service{
		service:     service,
		requirement: requirement,
	}
}

func (s Service) Type() service.Type {
	return s.service.Type()
}

func (s Service) Version() string {
	return s.service.Version()
}

func (s Service) Name() string {
	return s.service.Name()
}

func (s Service) Metadata() interface{} {
	return s.service.Metadata()
}

func (s Service) GetRoutes() []req.RouteHandler {
	return s.service.GetRoutes()
}

func (s Service) VerifyAuthorization(req http.Request, token *auth.Token) error {
	err := s.requirement.VerifyAuthorization(&req, token)
	if err != nil {
		return err
	}

	return s.service.VerifyAuthorization(req, token)
}
