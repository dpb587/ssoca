package server

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

func (s Service) ParseRequestAuth(req http.Request) (auth.Token, error) {
	username, password, ok := req.BasicAuth()
	if !ok {
		return nil, nil
	}

	for _, user := range s.config.Users {
		if user.Username != username {
			continue
		} else if user.Password != password {
			continue
		}

		return auth.NewSimpleToken(username, user.Attributes), nil
	}

	return nil, errors.New("Invalid authentication")
}
