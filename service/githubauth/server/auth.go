package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
)

func (s Service) ParseRequestAuth(req http.Request) (*auth.Token, error) {
	return s.oauth.ParseRequestAuth(req)
}

func (s Service) OAuthUserProfileLoader(client *http.Client) (token auth.Token, _ error) {
	ctx := context.Background()

	ghclient := github.NewClient(client)

	user, _, err := ghclient.Users.Get(ctx, "")
	if err != nil {
		return token, errors.Wrap(err, "fetching user info")
	}

	token.ID = *user.Login
	token.Attributes = map[auth.TokenAttribute]*string{}
	token.Attributes[auth.TokenUsernameAttribute] = user.Login

	if user.Name != nil {
		token.Attributes[auth.TokenNameAttribute] = user.Name
	}

	token.Groups = []string{}

	for nextPage := 1; nextPage != 0; {
		teams, resp, err := ghclient.Organizations.ListUserTeams(ctx, &github.ListOptions{Page: nextPage})
		if err != nil {
			return token, errors.Wrap(err, "listing user teams")
		}

		for _, team := range teams {
			token.Groups = append(token.Groups, fmt.Sprintf("%s/%s", *team.Organization.Login, *team.Slug))
		}

		nextPage = resp.NextPage
	}

	for nextPage := 1; nextPage != 0; {
		orgs, resp, err := ghclient.Organizations.List(ctx, "", &github.ListOptions{Page: nextPage})
		if err != nil {
			return token, errors.Wrap(err, "listing user organizations")
		}

		for _, org := range orgs {
			token.Groups = append(token.Groups, *org.Login)
		}

		nextPage = resp.NextPage
	}

	return token, nil
}
