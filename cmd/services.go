package cmd

import (
	"fmt"
)

// ListServices print services to stdout.
func listServices(services interface{}) {
	fmt.Printf("Service: %#v\n", services)
}
