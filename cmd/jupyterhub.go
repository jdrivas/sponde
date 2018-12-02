package cmd

import (
	"fmt"

	jh "github.com/jdrivas/jhmon/jupyterhub"
	"github.com/spf13/cobra"
)

func cmdError(e error) {
	fmt.Printf("Error: %s\n", e)
}

func buildJupyterHub(mode runMode) {

	// Util
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Get the hub version",
		Long:  "Returns the version number of the running JupyterHub.",
		Run: func(cmd *cobra.Command, args []string) {
			version, err := jh.GetVersion()
			if err == nil {
				version.Print()
			} else {
				cmdError(err)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "Get detailed about hub",
		Long:  "Returns detailed information about the running Hub",
		Run: func(cmd *cobra.Command, args []string) {
			info, err := jh.GetInfo()
			if err == nil {
				info.Print()
			} else {
				cmdError(err)
			}
		},
	})

	// Proxy Routes
	var proxyCmd = &cobra.Command{
		Use:     "proxy",
		Aliases: []string{"routes"},
		Short:   "The proxy's routing table",
		Long:    "Returns the routing table from the JupyterHub proxy",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			routes, err := jh.GetProxy()
			if err == nil {
				routes.Print()
			} else {
				cmdError(err)
			}
		},
	}
	rootCmd.AddCommand(proxyCmd)
	listCmd.AddCommand(proxyCmd)

	// Users
	listCmd.AddCommand(&cobra.Command{
		Use:   "users",
		Short: "Get a data on a user or all users",
		Long:  "Returns a list of users from the connected Hub, or if users are specified, data on those usrs.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				users, err := jh.GetUsers()
				if err == nil {
					jh.ListUsers(users)
				} else {
					cmdError(err)
				}
			} else {
				var users jh.Users
				var badUsers []string
				for _, username := range args {
					user, err := jh.GetUser(username)
					if err == nil {
						users = append(users, user)
					} else {
						badUsers = append(badUsers, username)
					}
				}
				if len(users) > 0 {
					jh.ListUsers(users)
				}
				if len(badUsers) > 0 {
					fmt.Printf("\nThese users were not found:\n")
					for _, u := range badUsers {
						fmt.Printf("%s\n", u)
					}
				}
			}
		},
	})

	describeCmd.AddCommand(&cobra.Command{
		Use:   "users",
		Short: "Long description of a hub user",
		Long:  "Returns a description of the user on the Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := jh.GetUsers()
			if err == nil {
				jh.DescribeUsers(users)
			} else {
				cmdError(err)
			}
		},
	})

	// Tokens
	listCmd.AddCommand(&cobra.Command{
		Use:     "tokens",
		Aliases: []string{"token"},
		Short:   "get a users tokens",
		Long:    "Returns a list of tokens associated with the user.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := jh.GetTokens(args[0])
			if err == nil {
				tokens.Print()
			}
		},
	})

	// Services
	listCmd.AddCommand(&cobra.Command{
		Use:   "services",
		Short: "List of services",
		Long:  "Returns infomration of the services that the Hub supports.",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := jh.GetServices()
			if err == nil && len(services) > 0 {
				jh.ListServices(services)
			} else {
				if err != nil {
					cmdError(err)
				} else {
					fmt.Println("There were no services.")
				}
			}
		},
	})

}
