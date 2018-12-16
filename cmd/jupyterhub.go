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
		Use:     "connection <connection-name>",
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
		Use:   "users [<user-id> ...]",
		Short: "Users accessing the hub.",
		Long: `Returns a list of users from the connected Hub, 
or if users are specified, data on those users`,
		Run: doUsers(listUsers),
	}
	listCmd.AddCommand(listUsersCmd)

	var describeUsersCmd = &cobra.Command{
		Use:                   "users [flags] [<user-id> ...]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"user"},
		Short:                 "Hub users.",
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
		Use:   "admin <true|false>",
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
		Use:   "name <user-id>",
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

	// User Severs
	startCmd.AddCommand(&cobra.Command{
		Use:   "server <user-id>",
		Short: "Starts a users notebook server.",
		Long:  "Starts a users notebook server and will tell you if the server has started or pending starting on return.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			started, resp, err := getCurrentConnection().StartServer(args[0])
			DisplayF(displayServerStartedF(started, resp, err), resp, err)
		},
	})

	stopCmd.AddCommand(&cobra.Command{
		Use:   "server <user-id>",
		Short: "Stops a users notebook server.",
		Long:  "Stops a users notebook server and will tell you if the server has stopped or is pending stop on return.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			stopped, resp, err := getCurrentConnection().StopServer(args[0])
			DisplayF(displayServerStartedF(stopped, resp, err), resp, err)
		},
	})

	startCmd.AddCommand(&cobra.Command{
		Use:   "named-server <user-id> <server-name>",
		Short: "Start a named server",
		Long:  "Start a named server for a user.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			started, resp, err := getCurrentConnection().StartNamedServer(args[0], args[1])
			DisplayF(displayServerStartedF(started, resp, err), resp, err)
		},
	})

	stopCmd.AddCommand(&cobra.Command{
		Use:   "named-server <server-id> <server-name>",
		Short: "Command for managing named servers.",
		Long:  "Command for starting/stopping named servers.",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			stopped, resp, err := getCurrentConnection().StopNamedServer(args[0], args[1])
			DisplayF(displpayServerStopedF(stopped, resp, err), resp, err)
		},
	})

	// User Tokens
	listTokensCmd := &cobra.Command{
		Use:     "tokens <user-id>",
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

	createTokenCmd := &cobra.Command{
		Use:   "token [flags] <user-id> <note> ....",
		Short: "Create an API token for a user.",
		Long: `Creates a new API token for <user-id> with identifying text <note> 
(all text typed after the <user-id> is taken as a the test of the note.).

NOTE: This will display a token independently of the show-tokens command or any settings. 
This is the only place where this token will be displayed and you cannot get it back 
any other way. So, write it down if you intend to use it.`,
		DisableFlagsInUseLine: true,
		Args:                  cobra.MinimumNArgs(2),
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
		Use:   "token <user-id> <token-id>",
		Short: "Delete a users secrutity token",
		Long:  `Deletes the token specified by <user-id> and <token-id>`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().DeleteToken(args[0], args[1])
			Display(resp, err)
		},
	}
	deleteCmd.AddCommand(deleteTokenCmd)

	describeTokenCmd := &cobra.Command{
		Use:   "token <user-id> <token-id>",
		Short: "Detail for a security tokens",
		Long:  `Provides detail about a specific token for <user-id>> and <token-id>.`,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			token, resp, err := getCurrentConnection().GetToken(args[0], args[1])
			Describe(APIToken(token), resp, err)
		},
	}
	describeCmd.AddCommand(describeTokenCmd)

	// Hub Tokens

	// TODO: This currently only gets users back
	// and needs to get Services as well and then display correctly.
	getCmd.AddCommand(&cobra.Command{
		Use:   "owner <hub-token>",
		Short: "Identify a user or service from a Hub API token.",
		Long:  "Returns and displays a user or service from a Hub API token.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			user, resp, err := getCurrentConnection().GetTokenOwner(args[0])
			Describe(UserList(UserList{user}), resp, err)
		},
	})

	/* DEPRECATED so we'll leave it out.
	createCmd.AddCommand(&cobra.Command{
		Use:   "hub-token",
		Short: "Create a new token for communicating with this Hub.",
		Long:  "Retruns a new API token for communicating with this Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			token, resp, err := getCurrentConnection().CreateAPIToken()
			df := func() {
				fmt.Printf("The new token is: %v\n", token)
			}
			DisplayF(df, resp, err)
		},
	})
	*/
	// Groups
	listCmd.AddCommand(&cobra.Command{
		Use:   "groups",
		Short: "Groups registered with the Hub.",
		Long:  "Returns details of the groups that are defined with this Hub.",
		Run: func(cmd *cobra.Command, args []string) {
			groups, resp, err := getCurrentConnection().GetGroups()
			List(Groups(groups), resp, err)
		},
	})

	describeCmd.AddCommand(&cobra.Command{
		Use:   "group <group-name>",
		Short: "Details of a group",
		Long:  "Returns the details of the Hub user group <group-name>.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			group, resp, err := getCurrentConnection().GetGroup(args[0])
			Describe(Group(group), resp, err)
		},
	})

	createCmd.AddCommand(&cobra.Command{
		Use:   "group <group-name>",
		Short: "Create a group on the JupyterHub hub.",
		Long:  "Create a a new group named <group-name> on the hub.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().CreateGroup(args[0])
			Display(resp, err)
		},
	})

	deleteCmd.AddCommand(&cobra.Command{
		Use:   "group <group-name>",
		Short: "Delete a group from the Hub.",
		Long:  "Delete a group <group-nmae> from the JupyterHub hub.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := getCurrentConnection().DeleteGroup(args[0])
			Display(resp, err)
		},
	})

	// Users in groups
	addCmd.AddCommand(&cobra.Command{
		Use:                   "user [flags] <user-id> <group-name> ...",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"users"},
		Short:                 "Add user to a group.",
		Long:                  "Add a user <user-id> to a group <group-name> or list of groups <group1> <group2> ....",
		Example:               "  sponde add user david admin ee201-spring2019",
		Args:                  cobra.MinimumNArgs(2),
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
		Use:                   "user [flags] <group-name> <user-id> ...",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"users"},
		Short:                 "Remove a user or users from a group",
		Long:                  "Remove a user or list of users from the Hub user group <group-name>.",
		Args:                  cobra.MinimumNArgs(2),
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
	// TODO: Consider validating the HTTP verbs.
	httpCmd.AddCommand(&cobra.Command{
		Use:                   "send [flags] <method> <command> [<json-string> ....]",
		DisableFlagsInUseLine: true,
		Aliases:               []string{"SEND"},
		Short:                 "HTTP <method> <command> to hub.",
		Long: `Sends an HTTP <method> <command> to the Jupyterhub hub.
<method> is an HTTP verb (e.g. "GET")

All of the args following <command> are caputred as a single json 
string and placed in the body of the request, 
with the ContentType header set to application/json.`,
		Example: `  sponde http send post /groups/test/users {"name": "admin", "users": ["david"]}`,
		Args:    cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 2 {
				httpDisplay(getCurrentConnection().Send(strings.ToUpper(args[0]), args[1], nil))
			} else {
				httpDisplay(getCurrentConnection().SendJSONString(strings.ToUpper(args[0]), args[1], strings.Join(args[2:], " "), nil))
			}
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:                   "get [flags]  <command>",
		Aliases:               []string{"GET"},
		DisableFlagsInUseLine: true,
		Short:                 "HTTP GET <arg> to hub.",
		Args:                  cobra.MinimumNArgs(1),
		Long:                  " Sends an HTTP GET <arg> to the Jupyterhub hub.",
		Example:               `  sponde http get /users`,
		Run: func(cmd *cobra.Command, args []string) {
			httpDisplay(getCurrentConnection().Get(args[0], nil))
		},
	})

	httpCmd.AddCommand(&cobra.Command{
		Use:                   "post [flags] <command> [<json-string> ....]",
		Aliases:               []string{"POST"},
		DisableFlagsInUseLine: true,
		Short:                 "HTTP POST <command> to hub.",
		Long: `Sends an HTTP POST <command> to the Hub.  

All of the args follwing <command> are caputred as a single json 
string and placed in the body of the request, 
with the ContentType header set to application/json.`,
		Example: `  sponde http post /groups/test/users {"name": "admin", "users": ["david"]}`,
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
		Use:     "delete [flags] <command> [<json-string> ....]",
		Aliases: []string{"DELETE"},
		Short:   "HTTP DELETE <arg> to hub.",
		Long: `Sends an HTTP DELETE <command> to the Hub.  

All of the args following <command> are caputred as a single json 
string and placed in the body of the request, 
with the ContentType header set to application/json.`,
		Example: `  sponde http delete /groups/test/users {"name": "admin", "users": ["david"]}`,
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
