package cmd

import (
	"fmt"
	"os"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
)

type Routes jh.Routes

func (r Routes) List() {
	routes := jh.Routes(r)
	if len(routes) > 0 {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, " %s\n", t.Title("Routespec\tTarget\tUser\tLast Activity"))
		for _, ri := range routes {
			user := "<empty>"
			if ri.Data.Hub && ri.Data.User != "" {
				user = fmt.Sprintf("Hub : $s", ri.Data.User)
			} else if ri.Data.Hub {
				user = "Hub"
			} else if ri.Data.User != "" {
				user = ri.Data.User
			}
			fmt.Fprintf(w, "%s\n", t.Text("%s\t%s\t%s\t%s", ri.RouteSpec, ri.Target, user, ri.Data.LastActivity))
		}

		w.Flush()
	} else {
		fmt.Printf("There were no proxy routes.\n")
	}

}
