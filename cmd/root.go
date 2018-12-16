// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	// "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	rootCmd, setCmd, getCmd, httpCmd, interactiveCmd *cobra.Command
	listCmd, describeCmd, createCmd, deleteCmd       *cobra.Command
	addCmd, updateCmd, removeCmd                     *cobra.Command
	startCmd, stopCmd                                *cobra.Command
)

// This is pulled out specially, because for interactive
// it gets run before each line is parsed.
// runMode allows us to add or remove commands
// as necessary for interadtive use
type runMode int

const (
	interactive runMode = iota + 1
	commandline
)

func buildRoot(mode runMode) {

	interactiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Interactive mode",
		Long:  "Runs a command line interpreter with sematnics to make session use easy.",
		Run: func(cmd *cobra.Command, args []string) {
			DoInteractive()
		},
	}
	// Add the commands to the rootCmd node (e.g. http get /users).
	if mode != interactive {
		rootCmd.AddCommand(interactiveCmd)
	}

	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set a value or values.",
		Long:  "Sets the value or a list of values on an object or  applicaiton state.",
	}
	rootCmd.AddCommand(setCmd)

	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a value",
		Long:  "Get a value or list of values from an object or applicaiton state.",
	}
	rootCmd.AddCommand(getCmd)

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a resource on the hub.",
		Long:  "Create an object on the JupyterHub hub.",
	}
	rootCmd.AddCommand(createCmd)

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource on the hub.",
		Long:  "Delete an object on the JupyterHub hub.",
	}
	rootCmd.AddCommand(deleteCmd)

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add an element to a resource.",
		Long:  "Add an element to a resource on the JupyterHub hub.",
	}
	rootCmd.AddCommand(addCmd)

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an element from a resource.",
		Long:  "Update an element from a resource  on the JupyterHub hub.",
	}
	rootCmd.AddCommand(updateCmd)

	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove an element from a resource.",
		Long:  "Remove an element from a resource  on the JupyterHub hub.",
	}
	rootCmd.AddCommand(removeCmd)

	listCmd = &cobra.Command{
		Use: "list",
		// Aliases: []string{""},
		Short: "Short description of a collection of objects.",
		Long:  "Provides a short description of each element of a collection.",
	}
	rootCmd.AddCommand(listCmd)

	describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "Longer  description of a a collection of objects.",
		Long:  "Provides a longer, more complete  description of a collection object.",
	}
	rootCmd.AddCommand(describeCmd)

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start some process on the hub.",
		Long:  "Start some process or action on the hub, usually reporting back status.",
	}
	rootCmd.AddCommand(startCmd)

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stops some process on the hub.",
		Long:  "Stops some process or action on the hub, usually reporting back status.",
	}
	rootCmd.AddCommand(stopCmd)

	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Use HTTP verbs.",
		Long:  "Send requests to the HUB using with HTTP verbs and arguments.",
	}
	rootCmd.AddCommand(httpCmd)

	buildJupyterHub(mode)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	buildRoot(commandline)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//
// Flag and config file init.
//

const (
	configFlagKey       = "config"
	hubURLFlagKey       = "hub-url"
	tokenFlagKey        = "token"
	authRedirectFlagKey = "auth-redirect-url"
	clientIDFlagKey     = "client-id"
	clientSecretFlagKey = "client-secret"
	verboseFlagKey      = "verbose"
	debugFlagKey        = "debug"
)

var (
	cfgFile, tokenFV, hubURLFV                         string
	authClientIDFV, authClientSecretFV, authRedirectFV string

	verbose, debug bool
)

func init() {

	// Root is created here, rather than in build root, because for interative
	// any root command flags set on the original command line should persist
	// to _each_ interactive command. They can  be explicitly overridden if needed.
	rootCmd = &cobra.Command{
		Use:   "sponde <command> [<args>]",
		Short: "Connect and report on a JupyterHub Hub.",
		Long:  "A tool for managing a JuyterHub Hub through the JupyterHub API",
	}

	initFlags()
	// initConfig()
	// initConnectionWithFlags()

	cobra.OnInitialize(cobraInit)
}

// Intended to be executed once before each commend.
// This happens after the commands line has been parsed
// but before any CMDs have been executed.
func cobraInit() {
	initConfig()
	initConnectionWithFlags()
}

func initFlags() {

	// Rest flags to start
	rootCmd.ResetFlags()

	// Flags available to everyone.
	rootCmd.PersistentFlags().StringVar(&cfgFile, configFlagKey, "", "config file location. (default is .sponde.{yaml,json,toml}")

	// Connection paramaters
	rootCmd.PersistentFlags().StringVarP(&tokenFV, tokenFlagKey, "t", "", "connect to the JupyterhHub with this authorization token.")
	rootCmd.PersistentFlags().StringVarP(&hubURLFV, hubURLFlagKey, "u", "",
		fmt.Sprintf("connect to the JupyterhHub at this URL. (default is %s)", defaultHubURL))

	//  Auth paramaters
	rootCmd.PersistentFlags().StringVarP(&authRedirectFV, authRedirectFlagKey, "", "", "OAuth redirect url - only need for auth commands.")
	rootCmd.PersistentFlags().StringVarP(&authClientIDFV, clientIDFlagKey, "", "", "OAuth client id - only need for auth commands.")
	rootCmd.PersistentFlags().StringVarP(&authClientSecretFV, clientSecretFlagKey, "", "", "OAuth client secret - only need for auth commands.")

	// To suport configuration files populating values, as well as flags, bind the variables to
	// the viper instance.

	rootCmd.PersistentFlags().BoolVarP(&verbose, verboseFlagKey, "v", false, "Describe what is happening as its happening.")
	viper.BindPFlag(verboseFlagKey, rootCmd.PersistentFlags().Lookup(verboseFlagKey))

	rootCmd.PersistentFlags().BoolVarP(&debug, debugFlagKey, "d", false, "Describe details about what's happening.")
	viper.BindPFlag(debugFlagKey, rootCmd.PersistentFlags().Lookup(debugFlagKey))

	// Now init the Juphterhub specific flags.
	initJupyterHubFlags()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Fin a config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("sponde")

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra_test" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Read in the config file.
	if err := viper.ReadInConfig(); err == nil {
		if Debug() {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if Debug() {
			fmt.Printf("Error loading config file: %s - %v\n", viper.ConfigFileUsed(), err)
		}
	}

}

// This should be called AFTER the config file has been read.
func initConnectionWithFlags() {
	// Do the normal config file default
	initConnections()
	conn := getCurrentConnection()

	update := false
	// Then overide with flags as apprpriate
	if rootCmd.PersistentFlags().Lookup(hubURLFlagKey).Changed {
		conn.HubURL = hubURLFV
		update = true
	}
	if rootCmd.PersistentFlags().Lookup(tokenFlagKey).Changed {
		conn.Token = tokenFV
		update = true
	}
	if update {
		updateCurrentConnection(conn)
	}

	if listConnsCmd.PersistentFlags().Lookup(showTokensOnceFlagKey).Changed {
		setShowTokensOnce()
	} else {
		resetShowTokensOnce()
	}

}
