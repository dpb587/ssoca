package api

import (
	"net/http"

	"github.com/dpb587/ssoca/auth/authn"
	"github.com/dpb587/ssoca/auth/authz"
	apierr "github.com/dpb587/ssoca/server/api/errors"
	"github.com/dpb587/ssoca/server/requtil"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	"github.com/sirupsen/logrus"

	uuid "github.com/nu7hatch/gouuid"
)

type apiHandler struct {
	authService    service.AuthService
	apiService     service.Service
	handler        req.RouteHandler
	clientIPGetter requtil.ClientIPGetter
	logger         logrus.FieldLogger
}

func CreateHandler(
	authService service.AuthService,
	apiService service.Service,
	handler req.RouteHandler,
	clientIPGetter requtil.ClientIPGetter,
	logger logrus.FieldLogger,
) (http.Handler, error) {
	return apiHandler{
		authService:    authService,
		apiService:     apiService,
		handler:        handler,
		clientIPGetter: clientIPGetter,
		logger:         logger,
	}, nil
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := req.Request{
		RawRequest:  r,
		RawResponse: w,
	}

	requestUUID, err := uuid.NewV4()
	if err != nil {
		h.sendGenericErrorResponse(request, apierr.WrapError(err, "generating request ID"))

		return
	}

	clientIP, err := h.clientIPGetter(r)
	if err != nil {
		h.sendGenericErrorResponse(request, apierr.WrapError(err, "determining Client IP"))
	}

	request.ID = requestUUID.String()
	request.LoggerContext = logrus.Fields{
		"server.request.id":              request.ID,
		"server.request.client_ip":       clientIP.String(),
		"server.request.remote_addr":     r.RemoteAddr,
		"server.request.x_forwarded_for": r.Header.Get("x-forwarded-for"),
		"server.request.user_agent":      r.Header.Get("user-agent"),
		"service.name":                   h.apiService.Name(),
		"service.type":                   h.apiService.Type(),
	}

	token, err := h.authService.ParseRequestAuth(*r)
	if err != nil {
		// never allow a token if there was an error
		token = nil

		// differentiate unauthorized (essentially unauthorized, aka expired) vs forbidden (apparent auth, but invalid)
		if matchederr, matched := err.(apierr.Error); matched {
			if matchederr.Status == http.StatusUnauthorized {
				h.getRequestLogger(request).Debug(err)

				err = nil
			}
		}

		if err != nil {
			h.sendGenericErrorResponse(request, apierr.WrapError(err, "parsing authentication token"))

			return
		}
	}

	request.AuthToken = token

	err = h.apiService.VerifyAuthorization(*r, request.AuthToken)
	if err != nil {
		statusCode := http.StatusInternalServerError

		if _, ok := err.(authn.Error); ok {
			statusCode = http.StatusUnauthorized
		} else if _, ok := err.(authz.Error); ok {
			statusCode = http.StatusForbidden
		}

		h.sendGenericErrorResponse(request, apierr.WrapError(apierr.NewError(err, statusCode, ""), "checking service authorization"))

		return
	}

	if token != nil {
		request.LoggerContext["auth.user_id"] = token.ID
	}

	err = h.handler.Execute(request)
	if err != nil {
		h.sendGenericErrorResponse(request, apierr.WrapError(err, "executing handler"))
	}

	h.getRequestLogger(request).Info("finished request")
}

func (h apiHandler) sendGenericErrorResponse(request req.Request, err error) {
	h.sendErrorResponse(request, apierr.NewError(err, http.StatusInternalServerError, ""))
}

func (h apiHandler) sendErrorResponse(request req.Request, err apierr.Error) {
	request.RawResponse.WriteHeader(err.Status)

	var loggerFunc func(args ...interface{})
	logger := h.getRequestLogger(request)

	if err.Status >= 500 {
		loggerFunc = logger.Error
	} else {
		loggerFunc = logger.Warn
	}

	loggerFunc(err.Error())

	request.WritePayload(map[string]interface{}{
		"error": map[string]interface{}{
			"status":  err.Status,
			"message": err.PublicError,
		},
	})
}

func (h apiHandler) getRequestLogger(request req.Request) *logrus.Entry {
	return h.logger.WithFields(request.LoggerContext).WithFields(logrus.Fields{
		"server.request.method": request.RawRequest.Method,
		"server.request.path":   request.RawRequest.URL.Path,
	})
}
