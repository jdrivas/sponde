package cmd

import (
	"fmt"
	"os"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
)

// Groups are lists of groups of users kept on the hub.
type Groups jh.Groups

// Group is the detail of the hub group.
type Group jh.Group

// List displays a list of groups
func (groups Groups) List() {
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

// Describe is a more detailed description of a group
func (group Group) Describe() {
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

// UserGroup captures the changes for an update.
type UserGroup jh.UserGroup

// List displays the details of a UserGroup
func (userGroup UserGroup) List() {
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", t.Title("Name\tUsers"))
	fmt.Fprintf(w, "%s\n", t.SubTitle("%s\t%v", userGroup.Name, userGroup.UserNames))
	w.Flush()
}
