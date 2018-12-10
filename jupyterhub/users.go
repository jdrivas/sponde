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

// Description
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

// StateValues are returned from the sever.
type StateValues struct {
	PodName string `json:"pod_name"`
}

// UpdatedUser is the object to send to the server with udser updates.
type UpdatedUser struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
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

// UpdateUser changes a users name or admin status. Use the UpdatedUser object to specify and you only need
// to fill in the values that are changing, though it all works with a full object.
func UpdateUser(name string, user UpdatedUser) (returnUser UpdatedUser, resp *http.Response, err error) {
	resp, err = Patch(fmt.Sprintf("/users/%s", name), user, &returnUser)
	return returnUser, resp, err
}

//
// User Tokens
//
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

// GetTokens returns all of the tokens for the specified users
func GetTokens(username string) (token Tokens, resp *http.Response, err error) {
	resp, err = Get(fmt.Sprintf("/users/%s/tokens", username), &token)
	return token, resp, err
}

// GetToken returns the details of a usres particular token given by username and tokenID.
func GetToken(username, tokenID string) (token APIToken, resp *http.Response, err error) {
	resp, err = Get(fmt.Sprintf("/users/%s/tokens/%s", username, tokenID), &token)
	return token, resp, err
}

// CreateToken will create a single APIToken from the Template provided,
// for the user and return the newly created token.
// Only Note will be saved in the new token.
func CreateToken(username string, newToken APIToken) (createdToken APIToken, resp *http.Response, err error) {
	resp, err = Post(fmt.Sprintf("/users/%s/tokens", username), newToken, &createdToken)
	return createdToken, resp, err
}

// DeleteToken deletes the token identified by usernamd and tokenID.
func DeleteToken(username, tokenID string) (resp *http.Response, err error) {
	resp, err = Delete(fmt.Sprintf("/users/%s/tokens/%s", username, tokenID), nil, nil)
	return resp, err
}
