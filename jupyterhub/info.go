package jupyterhub

import (
	"fmt"
	"net/http"
)

// Info has basic Jupyterhub state information.
type Info struct {
	Authenticator Authenticator `json:"authenticator"`
	Version       string        `json:"version"`
	Python        string        `json:"python"`
	SysExecutable string        `json:"sys_executable"`
	Spawner       Spawner       `json:"spawner"`
}

// Authenticator is the Python class handling authentication for the hub.
type Authenticator struct {
	Class   string `json:"class"`
	Version string `json:"version"`
}

// Spawner is the class spawning notebooks servers.
type Spawner struct {
	Class   string `json:"class"`
	Version string `json:"version"`
}

// GetInfo returns the Hub's system information.
func (conn Connection) GetInfo() (info Info, resp *http.Response, err error) {
	resp, err = conn.Get(fmt.Sprintf("/info"), &info)
	return info, resp, err
}
