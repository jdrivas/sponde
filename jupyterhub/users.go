package jupyterhub

import (
	"encoding/json"
	"io/ioutil"
)

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

// GetUsers returns a list of logged in JupyterHub users.
func GetUsers() (users []User, err error) {
	resp, err := callJHGet("/users")
	body := []byte{}
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	json.Unmarshal(body, &users)
	return users, err
}
