package jupyterhub

// GetInfo returns the Hub's system information.
func GetServices() (services []map[string]interface{}, err error) {
	_, err = Get("/services", services)
	return services, err
}
