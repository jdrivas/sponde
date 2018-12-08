package jupyterhub

// GetInfo returns the Hub's system information.

type Services []Service
type Service struct {
	Name    string                 `json:"name"`
	Admin   bool                   `json:"admin"`
	URL     string                 `json:"url"`
	Prefix  string                 `json:"prefix"`
	PID     int                    `json:"pid"`
	Command []string               `json:"command"`
	Info    map[string]interface{} `json:"info"`
}

func GetServices() (services Services, err error) {
	_, err = Get("/services", services)
	return services, err
}
