package jupyterhub

import (
	"net/http"
)

// TODO: This needs to work for users and services!

// GetTokenOwner will return the user assocaited wit the token.
func (conn Connection) GetTokenOwner(token string) (u User, resp *http.Response, err error) {
	resp, err = conn.Get("/authorizations/token/"+token, &u)
	return u, resp, err
}

/* DPRECATED API, so we'll leave it out.

// CreateAPIToken returns a new token for communication with the connected Hub.
func (conn Connection) CreateAPIToken() (t map[string]interface{}, resp *http.Response, err error) {
	resp, err = conn.Post("/authorizations/token", nil, &t)
	return t, resp, err
}

*/
