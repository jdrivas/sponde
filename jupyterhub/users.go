package jupyterhub

import (
	"fmt"
	"net/http"
)

type UserList []User

type User struct {
	Name         string            `json:"name"`
	Kind         string            `json:"kind"`
	Admin        bool              `json:"admin"`
	Created      string            `json:"created"`
	LastActivity string            `json:"last_activity"`
	ServerURL    string            `json:"server"`
	Pending      string            `json:"pending"`
	Servers      map[string]Server `json:"servers"`
}

type Server struct {
	Name         string      `json:"name"`
	LastActivity string      `json:"last_activity"`
	Started      string      `json:"started"`
	Pending      string      `json:"pending"`
	Ready        bool        `json:"ready"`
	State        StateValues `json:"state"`
	URL          string      `json:"url"`
	ProgressURL  string      `json:"progress_url"`
}

type StateValues struct {
	PodName string `json:"progress_url"`
}

// GetUser retruns a users information
func GetUser(username string) (user User, err error) {
	_, err = getResult(fmt.Sprintf("%s%s", "/users", username), &user)
	return user, err
}

// GetUsers gets users details from the hub.
// It returns a list of users for those that are found
// on the hub,  list of usernamess that were not found,
// and an errorif there we any problems.
func GetUsers(usernames []string) (users UserList, badUsers []string, err error) {
	for _, un := range usernames {
		user := new(User)
		var resp *http.Response
		resp, err = getResult(fmt.Sprintf("%s%s", "/users/", un), user)
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
	return users, badUsers, err
}

// GetAllUsers returns a list of logged in JupyterHub users.
func GetAllUsers() (users UserList, err error) {
	_, err = getResult("/users", &users)
	return users, err
}

type Tokens struct {
	APITokens   []APIToken   `json:"api_tokens"`
	OAuthTokens []OAuthToken `json:"oauth_tokens"`
}

type APIToken struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	User         string `json:"user"`
	Created      string `json:"created"`
	LastActivity string `json:"last_activity"`
	Note         string `json:"note"`
}

type OAuthToken struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	User         string `json:"user"`
	Created      string `json:"created"`
	LastActivity string `json:"last_activity"`
	OAuthClient  string `json:"oauth_client"`
}

func GetTokens(username string) (token Tokens, err error) {
	_, err = getResult(fmt.Sprintf("/users/%s/tokens", username), &token)
	return token, err
}
