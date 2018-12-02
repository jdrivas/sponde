package jupyterhub

import (
	"fmt"
)

// Version is the version nof the running JupyterHub
type Version struct {
	Version string `json:"version"`
}

// GetVersion returns the version of the JupyterHub from querying JupyterHub API.
func GetVersion() (version Version, err error) {
	resp, err := callJHGet("/")
	if err == nil {
		unmarshal(resp, &version)
	}
	return version, err
}

func (version *Version) Print() {
	fmt.Printf("JupyterHub Version: %s\n", version.Version)
}
