package jupyterhub

import (
	"net/http"
)

// Shutdown tells the running hub to shutdown. In a Kubernetes environment
// this usually means restart.
func Shutdown() (resp *http.Response, err error) {
	return Post("/shutdown", nil, nil)
}
