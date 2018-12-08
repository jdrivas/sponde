package jupyterhub

import "net/http"

// Version is the version nof the running JupyterHub
type Version struct {
	Version string `json:"version"`
}

// GetVersion returns the version of the JupyterHub from querying JupyterHub API.
func GetVersion() (version Version, resp *http.Response, err error) {
	resp, err = Get("/", &version)
	return version, resp, err
}
