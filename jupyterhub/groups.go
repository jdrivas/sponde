package jupyterhub

import (
	"fmt"
	"net/http"
)

// Version is the version nof the running JupyterHub
type Groups []Group
type Group struct {
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	UserNames []string `json:"users"`
}

// This is used to post deltes of users from a group.
type UserGroup struct {
	Name      string   `json:"name"`
	UserNames []string `json:"users"`
}

// GetGroup return a named group.
func GetGroup(name string) (group Group, resp *http.Response, err error) {
	resp, err = Get(fmt.Sprintf("/groups/%s", name), &group)
	return group, resp, err
}

// GetGroups returns all of the groups on the hub.
func GetGroups() (groups Groups, resp *http.Response, err error) {
	resp, err = Get("/groups", &groups)
	return groups, resp, err
}

// CreateGroup creates a group with name on the hub.
func CreateGroup(name string) (resp *http.Response, err error) {
	resp, err = Post(fmt.Sprintf("/groups/%s", name), nil, nil)
	return resp, err
}

// DeleteGroup deletes the group named name from the hub.
func DeleteGroup(name string) (resp *http.Response, err error) {
	resp, err = Delete(fmt.Sprintf("/groups/%s", name), nil, nil)
	return resp, err
}

// AddUserToGroup adds the UserGroup.UserNames to the group UserGroup.Name on the hub.
func AddUserToGroup(user UserGroup) (returnUsers UserGroup, resp *http.Response, err error) {
	resp, err = Post(fmt.Sprintf("/groups/%s/users", user.Name), user, &returnUsers)
	return returnUsers, resp, err
}

// RemoveUserFromGroup  removes the UserGroup.UserNames from the group UserGroup.Name from the hub.
func RemoveUserFromGroup(user UserGroup) (returnUsers UserGroup, resp *http.Response, err error) {
	resp, err = Delete(fmt.Sprintf("/groups/%s/users", user.Name), user, &returnUsers)
	return returnUsers, resp, err
}
