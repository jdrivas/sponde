package jupyterhub

import "fmt"

// GetInfo returns the Hub's system information.
func GetServices() (services []map[string]interface{}, err error) {
	_, err = get("/services", services)
	return services, err
}

// ListServices print services to stdout.
func ListServices(services interface{}) {
	fmt.Printf("Service: %#v\n", services)
	// w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
	// fmt.Fprintf(w, "JupyterHub Version:\t%s\n", info.Version)
	// fmt.Fprintf(w, "JupyterHub System Executable:\t%s\n", info.SysExecutable)
	// fmt.Fprintf(w, "Authenticator Class:\t%s\n", info.Authenticator.Class)
	// fmt.Fprintf(w, "Authenticator Version:\t%s\n", info.Authenticator.Version)
	// fmt.Fprintf(w, "Spawner Class:\t%s\n", info.Spawner.Class)
	// fmt.Fprintf(w, "Spawner Version:\t%s\n", info.Spawner.Version)

	// python := strings.Split(info.Python, "\n")
	// if len(python) > 0 {
	// 	fmt.Fprintf(w, "Python:\t%s\n", python[0])
	// 	for _, l := range python[1:] {
	// 		fmt.Fprintf(w, "\t%s\n", l)
	// 	}
	// }
	// w.Flush()
}
