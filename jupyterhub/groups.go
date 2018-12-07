package jupyterhub

import "fmt"

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

// GetVersion returns the version of the JupyterHub from querying JupyterHub API.
func GetGroups() (groups Groups, err error) {
	_, err = getResult("/groups", &groups)
	return groups, err
}

func CreateGroup(name string) (err error) {
	_, err = Post(fmt.Sprintf("/groups/%s", name), nil)
	return err
}

func DeleteGroup(name string) (err error) {
	_, err = Delete(fmt.Sprintf("/groups/%s", name), nil)
	return err
}

func RemoveUserFromGroup(user UserGroup) (err error) {
	_, err = Post(fmt.Sprintf("/groups/%s/users", user.Name), user)
	return err
}
