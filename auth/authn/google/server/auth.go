package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"

	"github.com/dpb587/ssoca/auth"
)

type userinfoPayload struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (s Service) ParseRequestAuth(req http.Request) (*auth.Token, error) {
	return s.oauth.ParseRequestAuth(req)
}

func (s Service) OAuthUserProfileLoader(client *http.Client) (token auth.Token, _ error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return token, errors.Wrap(err, "Fetching user info")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return token, errors.New("Failed to request user info")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, errors.Wrap(err, "Reading user info")
	}

	userinfo := userinfoPayload{}
	err = json.Unmarshal(data, &userinfo)
	if err != nil {
		return token, errors.Wrap(err, "Unmarshaling user info")
	}

	if !userinfo.EmailVerified {
		return token, errors.New("Refusing to authenticate account with unverified email")
	}

	token.ID = userinfo.Email

	emailSplit := strings.Split(userinfo.Email, "@")

	token.Attributes = map[auth.TokenAttribute]*string{}
	token.Attributes[auth.TokenUsernameAttribute] = &userinfo.Email
	token.Attributes[auth.TokenEmailAttribute] = &userinfo.Email
	token.Attributes[auth.TokenNameAttribute] = &userinfo.Name

	token.Groups = []string{
		userinfo.Email,
		fmt.Sprintf("email/mailbox/%s", emailSplit[0]),
		fmt.Sprintf("email/domain/%s", emailSplit[1]),
	}

	if s.config.Scopes.CloudProject != nil {
		err = s.oauthUserProfileCloudProjectLoader(client, &token)
		if err != nil {
			return token, errors.Wrap(err, "Loading Cloud project scopes")
		}
	}

	return token, nil
}

func (s Service) oauthUserProfileCloudProjectLoader(client *http.Client, token *auth.Token) error {
	cloudresourcemanagerService, err := cloudresourcemanager.New(client)
	if err != nil {
		return errors.Wrap(err, "Creating API client")
	}

	res, err := cloudresourcemanagerService.Projects.List().PageSize(1024).Do()
	if err != nil {
		return errors.Wrap(err, "Listing projects")
	}

	allProjects := len(s.config.Scopes.CloudProject.Projects) == 0
	refProjects := map[string]bool{}
	for _, project := range s.config.Scopes.CloudProject.Projects {
		refProjects[project] = true
	}

	allRoles := len(s.config.Scopes.CloudProject.Roles) == 0
	refRoles := map[string]bool{}
	for _, role := range s.config.Scopes.CloudProject.Roles {
		refRoles[role] = true
	}

	refMember := fmt.Sprintf("user:%s", token.Email())

	for _, project := range res.Projects {
		if !allProjects && !refProjects[project.ProjectId] {
			continue
		}

		projectIam, err := cloudresourcemanagerService.Projects.GetIamPolicy(project.ProjectId, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
		if err != nil {
			return errors.Wrap(err, "Getting IAM policy")
		}

		for _, binding := range projectIam.Bindings {
			if !allRoles && !refRoles[binding.Role] {
				continue
			}

			for _, member := range binding.Members {
				if member == refMember {
					token.Groups = append(token.Groups, fmt.Sprintf("cloud/project/%s/%s", project.ProjectId, binding.Role))
				}
			}
		}
	}

	return nil
}
