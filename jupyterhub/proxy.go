package jupyterhub

import (
	"net/http"
)

// Routes are map of the routespec string to the Route information.
type Routes map[string]Route

// Route is the data the Hub provides on proxy routing.
type Route struct {
	RouteSpec string    `json:"routespec"`
	Target    string    `json:"target"`
	Data      RouteData `json:"data"`
}

// RouteData is the for whom detail.
type RouteData struct {
	User         string `json:"user"`
	ServerName   string `json:"server_name"`
	Hub          bool   `json:"hub"`
	LastActivity string `json:"last_activity"`
}

// GetProxy returns alist of routes maintained on the hub.
func (conn Connection) GetProxy() (routes Routes, resp *http.Response, err error) {
	resp, err = conn.Get("/proxy", &routes)
	return routes, resp, err
}
