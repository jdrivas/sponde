package cmd

import (
	"fmt"
	"os"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
)

type Groups jh.Groups

func (gs Groups) List() {
	groups := jh.Groups(gs)
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("Name\tKind\tUsers"))
	for _, g := range groups {
		fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%v", g.Name, g.Kind, g.UserNames))
	}
	w.Flush()
}
