package jupyterhub

import (
	"fmt"
)

type Info struct {
	Authenticator Auth    `json:"authenticator"`
	Version       string  `json:"version"`
	Python        string  `json:"python"`
	SysExecutable string  `json:"sys_executable"`
	Spawner       Spawner `json:"spawner"`
}

type Auth struct {
	Class   string `json:"class"`
	Version string `json:"version"`
}

type Spawner struct {
	Class   string `json:"class"`
	Version string `json:"version"`
}

// GetInfo returns the Hub's system information.
func GetInfo() (info Info, err error) {
	_, err = Get(fmt.Sprintf("/info"), &info)
	return info, err
}
