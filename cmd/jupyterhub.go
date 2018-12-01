package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	jh "github.com/jdrivas/jhmon/jupyterhub"
	"github.com/spf13/cobra"
)

// var listUsersCmd, describeUsersCmd *cobra.Command

func buildJupyterHub(mode runMode) {

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get detailed about hub",
		Long:  "Returns detailed information about the running Hub",
		Run: func(cmd *cobra.Command, args []string) {
			info, err := jh.GetInfo()
			if err == nil {
				python := strings.Split(info.Python, "\n")
				w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
				fmt.Fprintf(w, "JupyterHub Version:\t%s\n", info.Version)
				fmt.Fprintf(w, "JupyterHub System Executable:\t%s\n", info.SysExecutable)
				fmt.Fprintf(w, "Authenticator Class:\t%s\n", info.Authenticator.Class)
				fmt.Fprintf(w, "Authenticator Version:\t%s\n", info.Authenticator.Version)
				fmt.Fprintf(w, "Spawner Class:\t%s\n", info.Spawner.Class)
				fmt.Fprintf(w, "Spawner Version:\t%s\n", info.Spawner.Version)
				// fmt.Fprintf(w, "Python\n")
				for _, l := range python {
					fmt.Fprintf(w, "Python:\t%s\n", l)
				}
				w.Flush()

			} else {
				fmt.Printf("ERROR: %#v", err)
			}
		},
	}
	rootCmd.AddCommand(infoCmd)

	listUsersCmd := &cobra.Command{
		Use:   "users",
		Short: "Get a list of current hub users",
		Long:  "Returns a list of users from the connected Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := jh.GetUsers()
			if err == nil {
				w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
				fmt.Fprintf(w, "Name\tAdmin\tCreated\tServer\tLast\n")
				for _, u := range users {
					fmt.Fprintf(w, "%s\t%t\t%s\t%s\t%s\n", u.Name, u.Admin, u.Created, u.ServerURL, u.LastActivity)
				}
				w.Flush()
			} else {
				fmt.Printf("ERROR: %#v", err)
			}
		},
	}
	listCmd.AddCommand(listUsersCmd)

	describeUsersCmd := &cobra.Command{
		Use:   "users",
		Short: "Long description of a hub user",
		Long:  "Returns a description of the user on the Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			users, err := jh.GetUsers()
			if err == nil {
				for _, u := range users {
					w := tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
					fmt.Fprintf(w, "Name\tKind\tAdmin\tServer\n")
					fmt.Fprintf(w, "%s\t%s\t%t\t%s\n", u.Name, u.Kind, u.Admin, u.ServerURL)
					w.Flush()
					fmt.Println()
					w = tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
					fmt.Fprintf(w, "Created\tLast Activity\tPending\n")
					pending := "<empty>"
					if u.Pending != "" {
						pending = u.Pending
					}
					fmt.Fprintf(w, "%s\t%s\t%s\n", u.Created, u.LastActivity, pending)
					w.Flush()
					fmt.Printf("\nServers\n")
					for _, s := range u.Servers {
						w = tabwriter.NewWriter(os.Stdout, 4, 4, 3, ' ', 0)
						fmt.Fprintf(w, "Name\tReady\tPending\tStarted\tLast Activity\n")
						name := "<empty>"
						if s.Name != "" {
							name = s.Name
						}
						pending := "<empty>"
						if s.Pending != "" {
							pending = u.Pending
						}
						fmt.Fprintf(w, "%s\t%t\t%s\t%s\t%s\n", name, s.Ready, pending, s.Started, s.LastActivity)
						w.Flush()
					}
				}
			} else {
				fmt.Printf("ERROR: %#v", err)
			}
		},
	}
	describeCmd.AddCommand(describeUsersCmd)

}
