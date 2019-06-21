package server

import (
	"errors"
	"net/http"

	"github.com/dpb587/ssoca/auth"
)

func (s Service) SupportsRequestAuth(r http.Request) (bool, error) {
	username, password, _ := r.BasicAuth()

	return username != "" || password != "", nil
}

func (s Service) ParseRequestAuth(r http.Request) (*auth.Token, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, nil
	}

	for _, user := range s.config.Users {
		if user.Username != username {
			continue
		} else if user.Password != password {
			continue
		}

		token := auth.Token{}
		token.ID = username
		token.Groups = user.Groups

		token.Attributes = user.Attributes
		if token.Attributes == nil {
			token.Attributes = map[auth.TokenAttribute]*string{}
		}

		token.Attributes[auth.TokenUsernameAttribute] = &username

		return &token, nil
	}

	return nil, errors.New("invalid authentication")
}
