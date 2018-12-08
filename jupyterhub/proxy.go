package jupyterhub

import (
	"net/http"
)

// Routes are map of the routespec string to the Route information.
type Routes map[string]Route

type Route struct {
	RouteSpec string    `json:"routespec"`
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
func GetProxy() (routes Routes, resp *http.Response, err error) {
	resp, err = Get("/proxy", &routes)
	return routes, resp, err
}
