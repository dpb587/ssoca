package authn

import (
	"errors"
	"net/http"
	"strings"

	apierr "github.com/dpb587/ssoca/server/api/errors"
)

func ExtractBearerTokenValue(r http.Request) (string, error) {
	authValue := r.Header.Get("Authorization")
	if authValue == "" {
		return "", nil
	}

	authValuePieces := strings.SplitN(authValue, " ", 2)
	if len(authValuePieces) != 2 {
		return "", apierr.NewError(errors.New("invalid Authorization format"), http.StatusForbidden, "")
	} else if strings.ToLower(authValuePieces[0]) != "bearer" {
		return "", apierr.NewError(errors.New("invalid Authorization method"), http.StatusForbidden, "")
	}

	return authValuePieces[1], nil
}
