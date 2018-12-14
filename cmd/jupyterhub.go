package cmd

import (
	"fmt"
	"net/http"
	"strings"

	jh "github.com/jdrivas/sponde/jupyterhub"
	t "github.com/jdrivas/sponde/term"
	"github.com/spf13/cobra"
)

var listConnsCmd *cobra.Command

func buildJupyterHub(mode runMode) {

	//
	// Application Util
	//
	rootCmd.AddCommand(&cobra.Command{
		Use:   "show-tokens",
		Short: "Toggle display of tokens",
		Long:  "Toggles displapy of tokens in application. This will have no effect if the configuration variable neverShowTokens has been set.",
		Run: func(cmd *cobra.Command, args []string) {
			toggleShowTokens()
			if getShowTokens() {
				fmt.Printf("Showing tokens on.\n")
			} else {
				fmt.Printf("Showing tokens off.\n")
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
			err := setConnection(args[0])
			List(getAllConnections(), nil, err)
		},
	})

	listCmd.AddCommand(listConnsCmd)

	//
	// Hub Commands
	//

	// General Hub State
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "The version of JupyterHub.",
		Long:  "Returns the version number of the running JupyterHub.",
		Run: func(cmd *cobra.Command, args []string) {
			version, resp, err := getCurrentConnection().GetVersion()
			List(Version(version), resp, err)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "info",
		Short: "Hub operational details.",
		Long:  "Returns detailed information about the running Hub",
		Run: func(cmd *cobra.Command, args []string) {
			info, resp, err := getCurrentConnection().GetInfo()
			List(Info(info), resp, err)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "shutdown",
		Short: "Shutdown the hub",
		Long:  "Shutdown the connected JupyterHub hub",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().Shutdown()
			display := func() {
				var result string
				switch resp.StatusCode {
				case http.StatusAccepted:
					result = t.Success("shutting down.")
				default:
					result = t.Fail("Probably not shutting down.")
				}
				fmt.Printf("%s %s\n", t.Title("Hub is"), result)
			}
			DisplayF(display, resp, err)
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
			routes, resp, err := getCurrentConnection().GetProxy()
			List(Routes(routes), resp, err)
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
		Run: doUsers(listUsers),
	}
	listCmd.AddCommand(listUsersCmd)
	listUsersCmd.SetUsageTemplate(userArgsTemplate)

	var describeUsersCmd = &cobra.Command{
		Use:     "users",
		Aliases: []string{"user"},
		Short:   "Hub users.",
		Long: `Returns a longer description of hub users.
If no user-id is provided then all Hub users are described.`,
		Run: doUsers(describeUsers),
	}
	describeCmd.AddCommand(describeUsersCmd)

	var updateUsersCmd = &cobra.Command{
		Use:   "user",
		Short: "Change the name or admin status of an existing user.",
		Long:  "Change the name or admin status of an exiting hub user.",
	}
	updateCmd.AddCommand(updateUsersCmd)

	updateUsersCmd.AddCommand(&cobra.Command{
		Use:   "admin",
		Short: "Change the admin status of a user.",
		Long:  "Set the admin status of an existing hub user to \"true\" or \"false\".",
		Args: func(cmd *cobra.Command, args []string) (err error) {
			_, err = truthyString(args[1])
			return err
		},
		Run: func(cmd *cobra.Command, args []string) {
			v, _ := truthyString(args[1])
			u := jh.UpdatedUser{
				Name:  args[0],
				Admin: v,
			}
			updatedUser, resp, err := getCurrentConnection().UpdateUser(args[0], u)
			List(UpdatedUser(updatedUser), resp, err)
		},
	})

	updateUsersCmd.AddCommand(&cobra.Command{
		Use:   "name",
		Short: "Change the name status of a user.",
		Long:  "Set the name status of an existing hub user from <old-name> to <new-name>.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			u := jh.UpdatedUser{
				Name: args[1],
			}
			updatedUser, resp, err := getCurrentConnection().UpdateUser(args[0], u)
			List(UpdatedUser(updatedUser), resp, err)
		},
	})

	// User Tokens
	listTokensCmd := &cobra.Command{
		// listCmd.AddCommand(&cobra.Command{
		Use:     "tokens",
		Aliases: []string{"token"},
		Short:   "Users security tokens",
		Long: `Returns a list of a Hub user's seurity tokens.
This must be called with at least one user-id, but  may be called with a list.
The API, and so this command does not actually obtain the token itself.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tokens, resp, err := getCurrentConnection().GetTokens(args[0])
			List(Tokens(tokens), resp, err)
		},
	}
	listCmd.AddCommand(listTokensCmd)
	listTokensCmd.SetUsageTemplate(userOneArgsTemplate)

	describeTokenCmd := &cobra.Command{
		Use:   "token",
		Short: "Detail for a security tokens",
		Long:  `Provides detail about a specific token for <username> and <token-id>.`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			token, resp, err := getCurrentConnection().GetToken(args[0], args[1])
			Describe(APIToken(token), resp, err)
		},
	}
	describeCmd.AddCommand(describeTokenCmd)

	createTokenCmd := &cobra.Command{
		Use:   "token",
		Short: "Create an API token for a user.",
		Long: `Creates a new API token for <username> with identifying text <note>.
NOTE: This will display a token independently of the show-tokens command or any settings. 
This is the only place where this token will be displayed and you cannot get it back 
any other way. So, write it down if you intend to use it.`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			notes := strings.Join(args[1:], " ")
			tokenTemplate := jh.APIToken{
				User: name,
				Note: notes,
			}
			token, resp, err := getCurrentConnection().CreateToken(name, tokenTemplate)
			display := func() {
				if err == nil && token.Token != "" {
					fmt.Printf("\n%s %s\n\n", t.Success("New token:"), t.Title(token.Token))
					APIToken(token).Describe()
				}
			}
			DisplayF(display, resp, err)
		},
	}
	createCmd.AddCommand(createTokenCmd)

	deleteTokenCmd := &cobra.Command{
		Use:   "token",
		Short: "Delete a users secrutity token",
		Long:  `Deletes the token specified by <username> and <token-id>`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().DeleteToken(args[0], args[1])
			Display(resp, err)
		},
	}
	deleteCmd.AddCommand(deleteTokenCmd)

	// User Severs
	startCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Starts a users notebook server.",
		Long:  "Starts a users notebook server and will tell you if the server has started yet.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			started, resp, err := getCurrentConnection().StartServer(args[0])
			DisplayF(displayServerStartedF(started, resp, err), resp, err)
		},
	})

	stopCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Stops a users notebook server.",
		Long:  "Stops a users notebook server and will tell you if the server has started yet.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stopped, resp, err := getCurrentConnection().StopServer(args[0])
			DisplayF(displayServerStartedF(stopped, resp, err), resp, err)
		},
	})

	startCmd.AddCommand(&cobra.Command{
		Use:   "named-server",
		Short: "Start a named server",
		Long:  "Start a named server for a user.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			started, resp, err := getCurrentConnection().StartNamedServer(args[0], args[1])
			DisplayF(displayServerStartedF(started, resp, err), resp, err)
		},
	})

	stopCmd.AddCommand(&cobra.Command{
		Use:   "named-server",
		Short: "Command for managing named servers.",
		Long:  "Command for starting/stopping named servers.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			stopped, resp, err := getCurrentConnection().StopNamedServer(args[0], args[1])
			DisplayF(displpayServerStopedF(stopped, resp, err), resp, err)
		},
	})

	// Groups
	listCmd.AddCommand(&cobra.Command{
		Use:   "groups",
		Short: "Groups registered with the Hub.",
		Long:  "Returns details the groups that are defined with this Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			groups, resp, err := getCurrentConnection().GetGroups()
			List(Groups(groups), resp, err)
		},
	})

	describeCmd.AddCommand(&cobra.Command{
		Use:   "group",
		Short: "Groups registered with the Hub.",
		Long:  "Returns details the groups that are defined with this Hub.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			group, resp, err := getCurrentConnection().GetGroup(args[0])
			Describe(Group(group), resp, err)
		},
	})

	createCmd.AddCommand(&cobra.Command{
		Use:   "group",
		Short: "Create a group on the JupyterHub hub.",
		Long:  "Create a a new group on the JupyterHub hub. Requires a name as the first and only argument.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().CreateGroup(args[0])
			Display(resp, err)
		},
	})

	deleteCmd.AddCommand(&cobra.Command{
		Use:   "group",
		Short: "Delete a group on the JupyterHub hub.",
		Long:  "Delete a group on the JupyterHub hub. Requires a name as the first and only argument.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().DeleteGroup(args[0])
			Display(resp, err)
		},
	})

	// Users in groups
	addCmd.AddCommand(&cobra.Command{
		Use:   "user",
		Short: "Add user to group",
		Long:  "Add a named user <user> to the gruoup <group>.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ug := jh.UserGroup{
				Name:      args[len(args)-1],
				UserNames: args[:len(args)-1],
			}
			userGroup, resp, err := getCurrentConnection().AddUserToGroup(ug)
			List(UserGroup(userGroup), resp, err)
		},
	})

	removeCmd.AddCommand(&cobra.Command{
		Use:   "user",
		Short: "Remove a user from a group",
		Long:  "Remove a named user <user> from the gruoup <group>.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ug := jh.UserGroup{
				Name:      args[len(args)-1],
				UserNames: args[:len(args)-1],
			}
			userGroup, resp, err := getCurrentConnection().RemoveUserFromGroup(ug)
			List(UserGroup(userGroup), resp, err)
		},
	})

	// Services
	listCmd.AddCommand(&cobra.Command{
		Use:   "services",
		Short: "Services registered with the Hub.",
		Long:  "Returns details of the services that the Hub supports.",
		Run: func(cmd *cobra.Command, args []string) {
			services, err := getCurrentConnection().GetServices()
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

	// HTTP Util
	httpCmd.AddCommand(&cobra.Command{
		Use:     "send",
		Aliases: []string{"SEND"},
		Short:   "HTTP <method> <arg> to hub.",
		Long: `Sends an HTTP <method> <arg> to the Jupyterhub hub.
			<method> is an HTTP verb (e.g. "GET")`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 2 {
				httpDisplay(getCurrentConnection().Send(args[0], args[1], nil))
			} else {
				httpDisplay(getCurrentConnection().SendJSONString(args[0], args[1], strings.Join(args[2:], " "), nil))
			}
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:     "get",
		Aliases: []string{"GET"},
		Short:   "HTTP GET <arg> to hub.",
		Long:    "Sends an HTTP GET <arg> to the Jupyterhub hub.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			httpDisplay(getCurrentConnection().Get(args[0], nil))
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:     "post",
		Aliases: []string{"POST"},
		Short:   "HTTP POST <arg> to hub.",
		Long:    "Sends an HTTP POST <arg> to the Jupyterhub hub.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				httpDisplay(getCurrentConnection().Post(args[0], strings.Join(args[1:], " "), nil))
			} else {
				httpDisplay(getCurrentConnection().Post(args[0], nil, nil))
			}
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:     "delete",
		Aliases: []string{"DELETE"},
		Short:   "HTTP DELETE <arg> to hub.",
		Long:    "Sends an HTTP DELETE <arg> to the Jupyterhub hub.",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			httpDisplay(getCurrentConnection().Delete(args[0], nil, nil))
		},
	})

}

func init() {
	// This is a bit like rootCmd, we want to initialize this and hang flags off of it before it gets
	// executed. I'm not sure whty ....
	listConnsCmd = &cobra.Command{
		Use:     "connections",
		Aliases: []string{"conn", "con", "conns", "cons", "connection"},
		Short:   "Available connections to a JupyterHub hub.",
		Long:    "List all o fthe aviallable JupyterHub hub connections.",
		Run: func(cmd *cobra.Command, args []string) {
			conns := getAllConnections()
			// TODO: DisplayFunction when there is no resp or error
			List(conns, nil, nil)
		},
	}

}

const showTokensOnceFlagKey = "show-tokens"

var showTokensOnceFlagV bool

func initJupyterHubFlags() {

	// This flag should only work on the single command ie. it's not durable across
	// invocations in interactive mode.
	// This ensures that the value of showTokensOnce is only inited at the before
	// anything has executed and will make sure that
	// this flag only works when it is set each time through the interactive loop.
	// cobra.OnInitialize(func() {
	// fmt.Printf("** OnInitialize(cmd/jupytethub.go init()).\n")
	listConnsCmd.ResetFlags()
	listConnsCmd.PersistentFlags().BoolVarP(&showTokensOnceFlagV, showTokensOnceFlagKey, "s", false, "Show tokens when listing connecitions.")
	// })
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
