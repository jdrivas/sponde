package jupyterhub

import (
	"fmt"
	"net/http"
)

type UserList []User

type User struct {
	Kind         string            `json:"kind"`
	Name         string            `json:"name"`
	Admin        bool              `json:"admin"`
	Groups       []string          `json:"groups"`
	ServerURL    string            `json:"server"`
	Pending      string            `json:"pending"`
	Created      string            `json:"created"`
	LastActivity string            `json:"last_activity"`
	Servers      map[string]Server `json:"servers"`
}

type Server struct {
	Name         string      `json:"name"`
	Ready        bool        `json:"ready"`
	Pending      string      `json:"pending"`
	URL          string      `json:"url"`
	ProgressURL  string      `json:"progress_url"`
	Started      string      `json:"started"`
	LastActivity string      `json:"last_activity"`
	State        StateValues `json:"state"`
}

type StateValues struct {
	PodName string `json:"progress_url"`
}

// GetUser retruns a users information
func GetUser(username string) (user User, resp *http.Response, err error) {
	resp, err = Get(fmt.Sprintf("%s%s", "/users", username), &user)
	return user, resp, err
}

// GetUsers gets users details from the hub.
// It returns a list of users for those that are found
// on the hub,  list of usernamess that were not found,
// and an errorif there we any problems.
// TODO: This make one call for each user. This is inefficient for
// hubs with a small number of users, but probably more efficent for hubs
// with a large numbr of users. Decide if this should change.
// TODO: Depdneing on the previous TODO, note that only the last calls
// http.Response is returned, this is indicative of needing a better solution (perhaps
// move this logic into a CMD function that is used rather than here in JH.)
func GetUsers(usernames []string) (users UserList, badUsers []string, resp *http.Response, err error) {
	for _, un := range usernames {
		user := new(User)
		resp, err = Get(fmt.Sprintf("%s%s", "/users/", un), user)
		if err == nil {
			users = append(users, *user)
		} else {
			if resp.StatusCode == http.StatusNotFound {
				badUsers = append(badUsers, un)
				err = nil
			} else {
				break
			}
		}
	}
	return users, badUsers, resp, err
}

// GetAllUsers returns a list of logged in JupyterHub users.
func GetAllUsers() (users UserList, resp *http.Response, err error) {
	resp, err = Get("/users", &users)
	return users, resp, err
}

type Tokens struct {
	APITokens   []APIToken   `json:"api_tokens"`
	OAuthTokens []OAuthToken `json:"oauth_tokens"`
}

type APIToken struct {
	Kind         string `json:"kind"`
	ID           string `json:"id"`
	User         string `json:"user"`
	Service      string `json:"service"`
	Note         string `json:"note"`
	Created      string `json:"created"`
	Expires      string `json:"expires"`
	LastActivity string `json:"last_activity"`
}

type OAuthToken struct {
	Kind         string `json:"kind"`
	ID           string `json:"id"`
	User         string `json:"user"`
	Service      string `json:"service"`
	Note         string `json:"note"`
	Created      string `json:"created"`
	Expires      string `json:"expires"`
	LastActivity string `json:"last_activity"`
	OAuthClient  string `json:"oauth_client"`
}

func GetTokens(username string) (token Tokens, resp *http.Response, err error) {
	resp, err = Get(fmt.Sprintf("/users/%s/tokens", username), &token)
	return token, resp, err
}
