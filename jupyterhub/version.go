package jupyterhub

// Version is the version nof the running JupyterHub
type Version struct {
	Version string `json:"version"`
}

// GetVersion returns the version of the JupyterHub from querying JupyterHub API.
func GetVersion() (version Version, err error) {
	_, err = Get("/", &version)
	return version, err
}
