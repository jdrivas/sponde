package cmd

import (
	"fmt"

	jh "github.com/jdrivas/sponde/jupyterhub"
	"github.com/spf13/cobra"
)

// doUsers is a command handler that will print a list of all users on the hub
// if no arguments are provided, or treat arguments as user names and print a list of users
// found on the Hub with details, and the names of users not found on the hub.
func doUsers(listFunc func(jh.UserList), logError func(error)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		var users jh.UserList
		var badNames []string
		var err error
		if len(args) > 0 {
			users, badNames, err = jh.GetUsers(args)
		} else {
			users, err = jh.GetAllUsers()
		}
		if err == nil {
			if len(users) > 0 {
				listFunc(users)
			}
			if len(users) > 0 && len(badNames) > 0 {
				fmt.Println("")
			}
			if len(badNames) > 0 {
				// TODO: Pluralize
				fmt.Printf("There were %d user names not found on the Hub:\n", len(badNames))
				for _, n := range badNames {
					fmt.Printf("%s\n", n)
				}
			}
		} else {
			logError(err)
		}
	}
}
