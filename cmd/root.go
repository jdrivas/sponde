// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	"github.com/jdrivas/sponde/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd, setCmd, httpCmd, interactiveCmd   *cobra.Command
	listCmd, describeCmd, createCmd, deleteCmd *cobra.Command
	addCmd, removeCmd                          *cobra.Command
	cfgFile, tokenFlagVar, hubURLFlagVar       string
	verbose, debug                             bool
)

const defaultHubURL = "http://127.0.0.1:8081"

// This is pulled out specially, because for interactive
// it gets run before each line is parsed.
// runMode allows us to add or remove commands
// as necessary for interadtive use
type runMode int

const (
	interactive runMode = 0
	commandline runMode = 1
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

	setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set a value or values.",
		Long:  "Sets the value or a list of values on an object or  applicaiton state.",
	}
	rootCmd.AddCommand(setCmd)

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

	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Use HTTP verbs.",
		Long:  "Send requests to the HUB using with HTTP verbs and arguments.",
	}
	rootCmd.AddCommand(httpCmd)

	// Add the commands to the rootCmd node (e.g. http get /users).
	if mode != interactive {
		rootCmd.AddCommand(interactiveCmd)
	}

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

func init() {

	// Root is created here, rather than in build root, because for interative
	// any root command flags set on the original command line should persist
	// to each interactive command. They can  be explicitly overridden if needed.
	rootCmd = &cobra.Command{
		Use:   "sponde",
		Short: "Connect and report on a JupyterHub Hub.",
		Long:  "A tool for managing a JuyterHub Hub through the JupyterHub API",
	}

	// Flags available to everyone.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location. (default is .sponde.{yaml,json,toml}")
	rootCmd.PersistentFlags().StringVarP(&tokenFlagVar, "token", "t", "", "connect to the JupyterhHub with this authorization token.")
	rootCmd.PersistentFlags().StringVarP(&hubURLFlagVar, "hub-url", "u", "",
		fmt.Sprintf("connect to the JupyterhHub at this URL. (default is %s)", defaultHubURL))

	// viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	// viper.BindPFlag("hubURL", rootCmd.PersistentFlags().Lookup("hub-url"))

	// To suport configuration files to populate, as well as flags, use these viper variables
	// to access this global state
	// e.g. token := viper.GetBool("debug")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Describe what is happening as its happening.")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Describe details about what's happening.")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Called before any command, and so in interactive mode, each time a command is executed.
	cobra.OnInitialize(initConfig)

	// // Let's just read this in once at
	// initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if viper.GetBool("debug") {
		fmt.Printf("Reading config file.\n")
	}
	if cfgFile != "" {
		// Use config file from the flag.
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

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if viper.GetBool("debug") {
			fmt.Printf("Error loading config file: %s - %v\n", viper.ConfigFileUsed(), err)
		}
	}

	// Set up config managed durable variables
	initConnections()
}

func initConnections() {

	config.InitConnections(defaultHubURL)
	if rootCmd.PersistentFlags().Lookup("hub-url").Changed {
		config.UpdateDefaultHubURL(hubURLFlagVar)
	}
	if rootCmd.PersistentFlags().Lookup("token").Changed {
		config.UpdateDefaultToken(tokenFlagVar)
	}
}
