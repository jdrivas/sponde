package jupyterhub

import "fmt"

// Version is the version nof the running JupyterHub
type Groups []Group
type Group struct {
	Name  string                   `json:"name"`
	Kind  string                   `json:"kind"`
	Users []map[string]interface{} `json:"users"`
}

// GetVersion returns the version of the JupyterHub from querying JupyterHub API.
func GetGroups() (groups Groups, err error) {
	_, err = getResult("/groups", &groups)
	return groups, err
}

func CreateGroup(name string) (err error) {
	_, err = Post(fmt.Sprintf("/groups/%s", name))
	return err
}

func DeleteGroup(name string) (err error) {
	_, err = Delete(fmt.Sprintf("/groups/%s", name))
	return err
}
