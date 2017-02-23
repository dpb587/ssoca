package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dpb587/ssoca/auth"
	oauth2supportconfig "github.com/dpb587/ssoca/authn/support/oauth2/config"
	cloudresourcemanager "github.com/google/google-api-go-client/cloudresourcemanager/v1"
)

type userinfoPayload struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (s Service) ParseRequestAuth(req http.Request) (auth.Token, error) {
	return s.oauth.ParseRequestAuth(req)
}

func (s Service) OAuthUserProfileLoader(client *http.Client) (oauth2supportconfig.UserProfile, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return oauth2supportconfig.UserProfile{}, bosherr.WrapError(err, "Fetching user info")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return oauth2supportconfig.UserProfile{}, errors.New("Failed to request user info")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return oauth2supportconfig.UserProfile{}, bosherr.WrapError(err, "Reading user info")
	}

	userinfo := userinfoPayload{}
	err = json.Unmarshal(data, &userinfo)
	if err != nil {
		return oauth2supportconfig.UserProfile{}, bosherr.WrapError(err, "Unmarshaling user info")
	}

	if !userinfo.EmailVerified {
		return oauth2supportconfig.UserProfile{}, errors.New("Refusing to authenticate account with unverified email")
	}

	emailSplit := strings.Split(userinfo.Email, "@")

	userProfile := oauth2supportconfig.UserProfile{
		Username: userinfo.Email,
		Scopes: []string{
			userinfo.Email,
			fmt.Sprintf("email/mailbox/%s", emailSplit[0]),
			fmt.Sprintf("email/domain/%s", emailSplit[1]),
		},
		Attributes: map[string]string{
			"name": userinfo.Name,
		},
	}

	if s.config.Scopes.CloudProject != nil {
		err = s.oauthUserProfileCloudProjectLoader(client, &userProfile)
		if err != nil {
			return oauth2supportconfig.UserProfile{}, bosherr.WrapError(err, "Loading Cloud project scopes")
		}
	}

	return userProfile, nil
}

func (s Service) oauthUserProfileCloudProjectLoader(client *http.Client, userProfile *oauth2supportconfig.UserProfile) error {
	cloudresourcemanagerService, err := cloudresourcemanager.New(client)
	if err != nil {
		return bosherr.WrapError(err, "Creating API client")
	}

	res, err := cloudresourcemanagerService.Projects.List().PageSize(1024).Do()
	if err != nil {
		return bosherr.WrapError(err, "Listing projects")
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

	refMember := fmt.Sprintf("user:%s", userProfile.Username)

	for _, project := range res.Projects {
		if !allProjects && !refProjects[project.ProjectId] {
			continue
		}

		projectIam, err := cloudresourcemanagerService.Projects.GetIamPolicy(project.ProjectId, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
		if err != nil {
			return bosherr.WrapError(err, "Getting IAM policy")
		}

		for _, binding := range projectIam.Bindings {
			if !allRoles && !refRoles[binding.Role] {
				continue
			}

			for _, member := range binding.Members {
				if member == refMember {
					userProfile.Scopes = append(userProfile.Scopes, fmt.Sprintf("cloud/project/%s/%s", project.ProjectId, binding.Role))
				}
			}
		}
	}

	return nil
}
