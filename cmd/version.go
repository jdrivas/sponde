package cmd

import (
	"fmt"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
)

func PrintVersion(version jh.Version) {
	fmt.Printf("%s %s\n", t.Title("JupyterHub Version:"), t.Text(version.Version))
}
