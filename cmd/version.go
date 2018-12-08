package cmd

import (
	"fmt"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
)

type Version jh.Version

func (v Version) List() {
	version := jh.Version(v)
	fmt.Printf("%s %s\n", t.Title("JupyterHub Version:"), t.Text(version.Version))
}
