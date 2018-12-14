package jupyterhub

// GetInfo returns the Hub's system information.

// Services is a list of Service
type Services []Service

// Service is the State the hub keeps on a hub managed process.
type Service struct {
	Name    string                 `json:"name"`
	Admin   bool                   `json:"admin"`
	URL     string                 `json:"url"`
	Prefix  string                 `json:"prefix"`
	PID     int                    `json:"pid"`
	Command []string               `json:"command"`
	Info    map[string]interface{} `json:"info"`
}

// GetServices lists the services on the Hub.
func (conn Connection) GetServices() (services Services, err error) {
	_, err = conn.Get("/services", services)
	return services, err
}
