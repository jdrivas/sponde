package jupyterhub

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/viper"
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
	_, err = get(fmt.Sprintf("/info"), &info)
	return info, err
}

// Print prints to stdout a list view of hub's information.
func (info *Info) Print() {
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "JupyterHub API URL:\t%s\n", viper.GetString("hubURL"))
	fmt.Fprintf(w, "JupyterHub Version:\t%s\n", info.Version)
	fmt.Fprintf(w, "JupyterHub System Executable:\t%s\n", info.SysExecutable)
	fmt.Fprintf(w, "Authenticator Class:\t%s\n", info.Authenticator.Class)
	fmt.Fprintf(w, "Authenticator Version:\t%s\n", info.Authenticator.Version)
	fmt.Fprintf(w, "Spawner Class:\t%s\n", info.Spawner.Class)
	fmt.Fprintf(w, "Spawner Version:\t%s\n", info.Spawner.Version)

	python := strings.Split(info.Python, "\n")
	if len(python) > 0 {
		fmt.Fprintf(w, "Python:\t%s\n", python[0])
		for _, l := range python[1:] {
			fmt.Fprintf(w, "\t%s\n", l)
		}
	}
	w.Flush()
}
