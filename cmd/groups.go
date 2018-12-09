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
	if len(groups) < 3 {
		fmt.Fprintf(w, "%s\n", t.Title("Name\tKind\t# Users"))
		for _, g := range groups {
			fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%v", g.Name, g.Kind, g.UserNames))
		}
	} else {
		fmt.Fprintf(w, "%s\n", t.Title("Name\tKind\t# Users"))
		for _, g := range groups {
			fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%d", g.Name, g.Kind, len(g.UserNames)))
		}
	}
	w.Flush()
}

type Group jh.Group

func (g Group) Describe() {
	group := jh.Group(g)
	userNames := group.UserNames
	firstUserName := "<no-users>"
	if len(userNames) > 0 {
		firstUserName = userNames[0]
		userNames = userNames[1:]
	}
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("Name\tKind\tUsers"))
	fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%s\t%s", group.Name, group.Kind, firstUserName))
	for _, name := range userNames {
		fmt.Fprintf(w, "\t\t%s\n", t.SubTitle(name))
	}
	w.Flush()
}

type UserGroup jh.UserGroup

func (userGroup UserGroup) List() {
	ug := jh.UserGroup(userGroup)
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("Name\tUsers"))
	fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%v", ug.Name, ug.UserNames))
	w.Flush()
}
