package server

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/github"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth"
	oauth2supportconfig "github.com/dpb587/ssoca/authn/support/oauth2/config"
)

func (s Service) ParseRequestAuth(req http.Request) (auth.Token, error) {
	return s.oauth.ParseRequestAuth(req)
}

func (s Service) OAuthUserProfileLoader(client *http.Client) (oauth2supportconfig.UserProfile, error) {
	userProfile := oauth2supportconfig.UserProfile{}
	ghclient := github.NewClient(client)

	user, _, err := ghclient.Users.Get("")
	if err != nil {
		return userProfile, bosherr.WrapError(err, "Fetching username")
	}

	userProfile.Username = *user.Login
	userProfile.Scopes = []string{userProfile.Username}

	for nextPage := 1; nextPage != 0; {
		teams, resp, err := ghclient.Organizations.ListUserTeams(&github.ListOptions{Page: nextPage})
		if err != nil {
			return userProfile, bosherr.WrapError(err, "Listing user teams")
		}

		for _, team := range teams {
			userProfile.Scopes = append(userProfile.Scopes, fmt.Sprintf("%s/%s", *team.Organization.Login, *team.Slug))
		}

		nextPage = resp.NextPage
	}

	for nextPage := 1; nextPage != 0; {
		orgs, resp, err := ghclient.Organizations.List("", &github.ListOptions{Page: nextPage})
		if err != nil {
			return userProfile, bosherr.WrapError(err, "Listing user organizations")
		}

		for _, org := range orgs {
			userProfile.Scopes = append(userProfile.Scopes, *org.Login)
		}

		nextPage = resp.NextPage
	}

	return userProfile, nil
}
