package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jdrivas/sponde/config"
	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
)

var listConnsCmd *cobra.Command
var showTokens, showTokensOnce bool

const showTokensOnceFlagKey = "show-tokens"

func buildJupyterHub(mode runMode) {

	// HTTP Util
	httpCmd.AddCommand(&cobra.Command{
		Use:   "send",
		Short: "HTTP <method> <arg> to hub.",
		Long: `Sends an HTTP <method> <arg> to the Jupyterhub hub.
		<method> is an HTTP verb (e.g. "GET")`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			doHTTPResponse(jh.Send(args[0], args[1], nil))
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "HTTP GET <arg> to hub.",
		Long:  "Sends an HTTP GET <arg> to the Jupyterhub hub.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doHTTPResponse(jh.Get(args[0], nil))
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:   "post",
		Short: "HTTP POST <arg> to hub.",
		Long:  "Sends an HTTP POST <arg> to the Jupyterhub hub.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doHTTPResponse(jh.Post(args[0], nil, nil))
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:   "post-content",
		Short: "HTTP POST <arg> to hub.",
		Long:  "Sends an HTTP POST <arg> to the Jupyterhub hub.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			doHTTPResponse(jh.SendJSONString("POST", args[0], strings.Join(args[1:], " "), nil))
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "HTTP DELETE <arg> to hub.",
		Long:  "Sends an HTTP DELETE <arg> to the Jupyterhub hub.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			doHTTPResponse(jh.Delete(args[0], nil, nil))
		},
	})

	// Util
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "The version of JupyterHub.",
		Long:  "Returns the version number of the running JupyterHub.",
		Run: func(cmd *cobra.Command, args []string) {
			version, err := jh.GetVersion()
			if err == nil {
				PrintVersion(version)
			} else {
				cmdError(err)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "Hub operational details.",
		Long:  "Returns detailed information about the running Hub",
		Run: func(cmd *cobra.Command, args []string) {
			info, err := jh.GetInfo()
			if err == nil {
				if err == nil {
					PrintInfo(info)
				}
			} else {
				cmdError(err)
			}
		},
	})

	// Connections
	setCmd.AddCommand(&cobra.Command{
		Use:     "connection",
		Aliases: []string{"conn", "con"},
		Short:   "Use the named connection to the Hub.",
		Long:    "Sets the connection to the JupyterHub Hub to the named connection. ",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.SetConnection(args[0])
			if err != nil {
				cmdError(err)
			}
		},
	})

	listConnsCmd = &cobra.Command{
		Use:     "connections",
		Aliases: []string{"conn", "con", "conns", "cons", "connection"},
		Short:   "Available connections to a JupyterHub hub.",
		Long:    "List all o fthe aviallable JupyterHub hub connections.",
		Run: func(cmd *cobra.Command, args []string) {
			conns := config.GetConnections()
			currentName := config.GetConnectionName()
			w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", t.Title("\tName\tURL\tToken"))
			for _, c := range conns {
				name := t.Text(c.Name)
				current := ""
				if c.Name == currentName {
					name = t.Highlight("%s", c.Name)
					current = t.Highlight("%s", "*")
				}
				token := config.MakeSafeTokenString(c, false, true)
				fmt.Fprintf(w, "%s\t%s\t%s\n", current, name, t.Text("%s\t%s", c.HubURL, token))
			}
			w.Flush()
		},
	}
	listCmd.AddCommand(listConnsCmd)
	// This flag should only work on the single command
	listConnsCmd.PersistentFlags().BoolVarP(&showTokensOnce, showTokensOnceFlagKey, "s", false, "Show tokens when listing connecitions.")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "show-tokens",
		Short: "Toggle display of tokens",
		Long:  "Toggles displapy of tokens in application. This will have no effect if the configuration variable neverShowTokens has been set.",
		Run: func(cmd *cobra.Command, args []string) {
			config.SetShowTokens(!config.GetShowTokens())
			if config.GetShowTokens() {
				fmt.Printf("Showing tokens on.\n")
			} else {
				fmt.Printf("Showing tokens off.\n")
			}
		},
	})
	// To make show-tokens command, and the underlying config file variables
	// durable on commands.
	cobra.OnInitialize(initShowTokensOnce)

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
				ListRoutes(routes)
			} else {
				cmdError(err)
			}
		},
	}
	rootCmd.AddCommand(proxyCmd)
	listCmd.AddCommand(proxyCmd)

	// Users
	var listUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "Users accessing the hub.",
		Long: `Returns a list of users from the connected Hub, 
or if users are specified, data on those users`,
		Run: doUsers(listUsers, cmdError),
	}
	listCmd.AddCommand(listUsersCmd)
	listUsersCmd.SetUsageTemplate(userArgsTemplate)

	var describeUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "Hub users.",
		Long: `Returns a longer description of hub users.
If no user-id is provided then all Hub users are described.`,
		Run: doUsers(describeUsers, cmdError),
	}
	describeCmd.AddCommand(describeUsersCmd)
	describeUsersCmd.SetUsageTemplate(userArgsTemplate)

	// Tokens
	tokenCmd := &cobra.Command{
		// listCmd.AddCommand(&cobra.Command{
		Use:     "tokens",
		Aliases: []string{"token"},
		Short:   "Users security tokens",
		Long: `Returns a list of a Hub user's seurity tokens.
This must be called with at least one user-id, but  may be called with a list.
The API, and so this command does not actually obtain the token itself.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := jh.GetTokens(args[0])
			if err == nil {
				listTokens(tokens)
			}
		},
	}
	listCmd.AddCommand(tokenCmd)
	tokenCmd.SetUsageTemplate(userOneArgsTemplate)

	// Groups
	listCmd.AddCommand(&cobra.Command{
		Use:   "groups",
		Short: "Groups registered with the Hub.",
		Long:  "Returns details the groups that are defined with this Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			groups, err := jh.GetGroups()
			if err == nil && len(groups) > 0 {
				listGroups(groups)
			} else {
				if err != nil {
					cmdError(err)
				} else {
					fmt.Println("There were no groups.")
				}
			}
		},
	})

	createCmd.AddCommand(&cobra.Command{
		Use:   "group",
		Short: "Create a group on the JupyterHub hub.",
		Long:  "Create a a new group on the JupyterHub hub. Requires a name as the first and only argument.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := jh.CreateGroup(args[0])
			if err != nil {
				cmdError(err)
			}
		},
	})

	deleteCmd.AddCommand(&cobra.Command{
		Use:   "group",
		Short: "Delete a group on the JupyterHub hub.",
		Long:  "Delete a group on the JupyterHub hub. Requires a name as the first and only argument.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := jh.DeleteGroup(args[0])
			if err != nil {
				cmdError(err)
			}
		},
	})

	// Services
	listCmd.AddCommand(&cobra.Command{
		Use:   "services",
		Short: "Services registered with the Hub.",
		Long:  "Returns details of the services that the Hub supports.",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := jh.GetServices()
			if err == nil && len(services) > 0 {
				listServices(services)
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

// var prevShowFlagSet = false

// // Yes this is goofy.
// // We want the flag to only effect once, not permentaly
// // set the state.
// // But we want the state to be durable if was set by command (show-tokens)
// // it was set by commands.
// // TODO: There are some patterns here to extract for next time.
func initShowTokensOnce() {
	if showTokensOnce {
		config.ShowTokensOnce()
	} else {
		config.ResetShowTokensOnce()
	}
}

// For use when the command can take, but doesn't have to, an arbitrary number of
// <user-id> arguments.
var userArgsTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}} [<user-id> ...]{{if gt (len .Aliases) 0}}

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// For use when the command must have a least one, but can take oer user arguemnts.
var userOneArgsTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}} <user-id> [<user-id> ...]{{if gt (len .Aliases) 0}}

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
