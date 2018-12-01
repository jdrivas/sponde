package jupyterhub

import (
	"encoding/json"
	"io/ioutil"
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

// GetUsers returns a list of logged in JupyterHub users.
func GetInfo() (info Info, err error) {
	resp, err := callJHGet("/info")
	body := []byte{}
	if err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}
	json.Unmarshal(body, &info)
	return info, err
}
