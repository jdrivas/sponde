package jupyterhub

import (
	"fmt"
	"net/http"
)

// Groups a list of groups.
type Groups []Group

// Group is the hub respresentation of a group of users.
type Group struct {
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	UserNames []string `json:"users"`
}

// UserGroup is state requried for Adding/Removing a user to a group
type UserGroup struct {
	Name      string   `json:"name"`
	UserNames []string `json:"users"`
}

// GetGroup return a named group.
func (conn Connection) GetGroup(name string) (group Group, resp *http.Response, err error) {
	resp, err = conn.Get(fmt.Sprintf("/groups/%s", name), &group)
	return group, resp, err
}

// GetGroups returns all of the groups on the hub.
func (conn Connection) GetGroups() (groups Groups, resp *http.Response, err error) {
	resp, err = conn.Get("/groups", &groups)
	return groups, resp, err
}

// CreateGroup creates a group with name on the hub.
func (conn Connection) CreateGroup(name string) (resp *http.Response, err error) {
	resp, err = conn.Post(fmt.Sprintf("/groups/%s", name), nil, nil)
	return resp, err
}

// DeleteGroup deletes the group named name from the hub.
func (conn Connection) DeleteGroup(name string) (resp *http.Response, err error) {
	resp, err = conn.Delete(fmt.Sprintf("/groups/%s", name), nil, nil)
	return resp, err
}

// AddUserToGroup adds the UserGroup.UserNames to the group UserGroup.Name on the hub.
func (conn Connection) AddUserToGroup(user UserGroup) (returnUsers UserGroup, resp *http.Response, err error) {
	resp, err = conn.Post(fmt.Sprintf("/groups/%s/users", user.Name), user, &returnUsers)
	return returnUsers, resp, err
}

// RemoveUserFromGroup  removes the UserGroup.UserNames from the group UserGroup.Name from the hub.
func (conn Connection) RemoveUserFromGroup(user UserGroup) (returnUsers UserGroup, resp *http.Response, err error) {
	resp, err = conn.Delete(fmt.Sprintf("/groups/%s/users", user.Name), user, &returnUsers)
	return returnUsers, resp, err
}
