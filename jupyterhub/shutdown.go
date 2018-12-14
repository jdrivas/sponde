package jupyterhub

import (
	"net/http"
)

// Shutdown tells the running hub to shutdown. In a Kubernetes environment
// this usually means restart.
func (conn Connection) Shutdown() (resp *http.Response, err error) {
	return conn.Post("/shutdown", nil, nil)
}
