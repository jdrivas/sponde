package jupyterhub

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// Routes are map of the routespec string to the Route information.
type Routes map[string]Route

type Route struct {
	RouteSpec string    `json:"routespec`
	Target    string    `json:"target"`
	Data      RouteData `json:"data"`
}

type RouteData struct {
	User         string `json:"user"`
	ServerName   string `json:"server_name"`
	Hub          bool   `json:"hub"`
	LastActivity string `json:"last_activity"`
}

// GetUsers returns a list of logged in JupyterHub users.
func GetProxy() (routes Routes, err error) {
	_, err = get("/proxy", &routes)
	return routes, err
}

func (routes *Routes) Print() {
	// fmt.Printf("Proxy: %#v\n", routes)
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, "Routespec\tTarget\tUser\tLast Activity\n")
	for _, ri := range *routes {
		user := "<empty>"
		if ri.Data.Hub && ri.Data.User != "" {
			user = fmt.Sprintf("Hub : $s", ri.Data.User)
		} else if ri.Data.Hub {
			user = "Hub"
		} else if ri.Data.User != "" {
			user = ri.Data.User
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ri.RouteSpec, ri.Target, user, ri.Data.LastActivity)
	}

	w.Flush()

}
