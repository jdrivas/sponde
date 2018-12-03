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
		Short: "The version of JupyterHub.",
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
		Short: "Hub operational details.",
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
	var listUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "Users accessing the hub.",
		Long: `Returns a list of users from the connected Hub, 
or if users are specified, data on those users`,
		Run: doUsers(jh.ListUsers, cmdError),
	}
	listCmd.AddCommand(listUsersCmd)
	listUsersCmd.SetUsageTemplate(userArgsTemplate)

	var describeUsersCmd = &cobra.Command{
		Use:   "users",
		Short: "Hub users.",
		Long: `Returns a longer description of hub users.
If no user-id is provided then all Hub users are described.`,
		Run: doUsers(jh.DescribeUsers, cmdError),
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
This must be called with at least one user-id, but  may be called with a list.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, err := jh.GetTokens(args[0])
			if err == nil {
				tokens.Print()
			}
		},
	}
	listCmd.AddCommand(tokenCmd)
	tokenCmd.SetUsageTemplate(userOneArgsTemplate)

	// Services
	listCmd.AddCommand(&cobra.Command{
		Use:   "services",
		Short: "Services registered with the Hub.",
		Long:  "Returns details of the services that the Hub supports.",
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
