package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/dpb587/ssoca/auth"
	"github.com/dpb587/ssoca/server/service"
	"github.com/dpb587/ssoca/server/service/req"
	uuid "github.com/nu7hatch/gouuid"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type apiHandler struct {
	authService service.AuthService
	apiService  service.Service
	handler     reflect.Value
	handlerIn   []func(http.ResponseWriter, *http.Request) (reflect.Value, error)
	handlerOut  func([]reflect.Value) (interface{}, error)
	logger      logrus.FieldLogger
}

func CreateAPIHandler(authService service.AuthService, apiService service.Service, handler req.RouteHandler, logger logrus.FieldLogger) (http.Handler, error) {
	handlerIn := []func(http.ResponseWriter, *http.Request) (reflect.Value, error){}

	handlerMethod := reflect.ValueOf(handler).MethodByName("Execute")
	if handlerMethod.IsNil() {
		return nil, errors.New("Route handler is missing Execute method")
	}

	handlerType := handlerMethod.Type()

	// inputs

	for numIdx := 0; numIdx < handlerType.NumIn(); numIdx++ {
		var handlerArgTransform func(http.ResponseWriter, *http.Request) (reflect.Value, error)

		handlerArg := handlerType.In(numIdx)

		if handlerArg.String() == "http.ResponseWriter" {
			handlerArgTransform = func(w http.ResponseWriter, _ *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(w), nil
			}
		} else if handlerArg.String() == "*http.Request" {
			handlerArgTransform = func(_ http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(r), nil
			}
		} else if handlerArg.String() == "logrus.Fields" {
			handlerArgTransform = func(_ http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				return reflect.ValueOf(r.Context().Value("loggerContext")), nil
			}
		} else if handlerArg.String() == "auth.Token" {
			handlerArgTransform = func(_ http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				token := r.Context().Value(auth.RequestToken)
				if token == nil {
					return reflect.Value{}, NewAPIError(errors.New("Token missing from request context"), http.StatusUnauthorized, "")
				}

				return reflect.ValueOf(token), nil
			}
		} else {
			// better be an api type
			handlerArgTransform = func(_ http.ResponseWriter, r *http.Request) (reflect.Value, error) {
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return reflect.Value{}, bosherr.WrapError(err, "Reading request body")
				}

				apiRequestData := reflect.New(handlerArg)

				err = json.Unmarshal(bytes, apiRequestData.Interface())
				if err != nil {
					return reflect.Value{}, NewAPIError(WrapError(err, "Unmarshaling request payload"), http.StatusBadRequest, "Invalid body")
				}

				return reflect.ValueOf(apiRequestData.Elem().Interface()), nil
			}
		}

		handlerIn = append(handlerIn, handlerArgTransform)
	}

	// outputs

	var handlerOut func([]reflect.Value) (interface{}, error)

	switch handlerType.NumOut() {
	case 0:
		handlerOut = func(_ []reflect.Value) (interface{}, error) {
			return nil, nil
		}
	case 1:
		handlerOut = func(v []reflect.Value) (interface{}, error) {
			if !v[0].IsNil() {
				return nil, v[0].Elem().Interface().(error)
			}

			return nil, nil
		}
	case 2:
		handlerOut = func(v []reflect.Value) (interface{}, error) {
			if !v[1].IsNil() {
				return v[0].Interface(), v[1].Elem().Interface().(error)
			}

			return v[0].Interface(), nil
		}
	default:
		return apiHandler{}, errors.New("Invalid handler function return values")
	}

	// container

	return apiHandler{
		authService: authService,
		apiService:  apiService,
		handler:     handlerMethod,
		handlerIn:   handlerIn,
		handlerOut:  handlerOut,
		logger:      logger,
	}, nil
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUUID, err := uuid.NewV4()
	if err != nil {
		h.sendGenericErrorResponse(w, r, WrapError(err, "Generating request ID"))

		return
	}

	loggerContext := logrus.Fields{
		"server.request.id":          requestUUID.String(),
		"server.request.remote_addr": r.RemoteAddr,
		"service.name":               h.apiService.Name(),
		"service.type":               h.apiService.Type(),
	}

	r = r.WithContext(context.WithValue(r.Context(), "loggerContext", loggerContext))

	token, err := h.authService.ParseRequestAuth(*r)
	if err != nil {
		h.sendGenericErrorResponse(w, r, WrapError(err, "Parsing authentication token"))

		return
	}

	authz, err := h.apiService.IsAuthorized(*r, token)
	if err != nil {
		h.sendGenericErrorResponse(w, r, WrapError(err, "Checking service authorization"))

		return
	} else if !authz {
		h.sendErrorResponse(w, r, NewAPIError(errors.New("Not authorized"), http.StatusForbidden, ""))

		return
	}

	if token != nil {
		loggerContext["auth.username"] = token.Username()

		r = r.WithContext(context.WithValue(r.Context(), auth.RequestToken, token))
		r = r.WithContext(context.WithValue(r.Context(), "loggerContext", loggerContext))
	}

	handlerIn := []reflect.Value{}

	for argTransformIdx, argTransform := range h.handlerIn {
		argValue, err1 := argTransform(w, r)
		if err1 != nil {
			h.sendGenericErrorResponse(w, r, WrapErrorf(err1, "Converting argument %d", argTransformIdx))

			return
		}

		handlerIn = append(handlerIn, argValue)
	}

	handlerReturn := h.handler.Call(handlerIn)

	res, err := h.handlerOut(handlerReturn)
	if err != nil {
		h.sendGenericErrorResponse(w, r, WrapError(err, "Executing handler"))
	} else if res != nil {
		h.sendResponse(w, r, res)
	} else {
		// they hadnled it directly; log what happened
		h.logResponse(w, r)
	}
}

func (h apiHandler) getRequestContext(w http.ResponseWriter, r *http.Request) logrus.FieldLogger {
	l := h.logger.WithFields(logrus.Fields{
		"server.request.method": r.Method,
		"server.request.path":   r.URL.Path,
	})

	loggerContext := r.Context().Value("loggerContext")
	if loggerContext != nil {
		loggerContext, _ := loggerContext.(logrus.Fields)
		l = l.WithFields(loggerContext)
	}

	return l
}

func (h apiHandler) logResponse(w http.ResponseWriter, r *http.Request) {
	h.getRequestContext(w, r).Debug("Handled server request")
}

func (h apiHandler) sendGenericErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.sendErrorResponse(w, r, NewAPIError(err, http.StatusInternalServerError, ""))
}

func (h apiHandler) sendErrorResponse(w http.ResponseWriter, r *http.Request, err APIError) {
	w.WriteHeader(err.Status)

	if err.Status > 500 {
		h.getRequestContext(w, r).Error(err.Error())
	} else {
		h.getRequestContext(w, r).Warn(err.Error())
	}

	h.sendResponse(w, r, map[string]interface{}{
		"error": map[string]interface{}{
			"status":  err.Status,
			"message": err.PublicError,
		},
	})
}

func (h apiHandler) sendResponse(w http.ResponseWriter, r *http.Request, payload interface{}) {
	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		// @todo hope we don't have recursive errors...
		h.sendErrorResponse(w, r, NewAPIError(bosherr.WrapError(err, "Marshaling response payload"), http.StatusInternalServerError, ""))

		return
	}

	w.Header().Add("Content-Type", "application/javascript")
	io.WriteString(w, string(bytes))
	io.WriteString(w, "\n")

	h.logResponse(w, r)
}
