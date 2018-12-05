// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	"github.com/jdrivas/jhmon/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd, setCmd, listCmd, describeCmd, interactiveCmd *cobra.Command
	cfgFile, token, hubURL                                string
	verbose, debug                                        bool
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

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"get"},
		Short:   "Short description of a collection of objects.",
		Long:    "Provides a short description of each element of a collection.",
	}

	describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "Longer  description of a a collection of objects.",
		Long:  "Provides a longer, more complete  description of a collection object.",
	}

	// Other Commands.
	if mode != interactive {
		rootCmd.AddCommand(interactiveCmd)
	}
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(describeCmd)
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
		Use:   "jhmon",
		Short: "Connect and report on a JupyterHub Hub.",
		Long:  "A tool for managing a JuyterHub Hub through the JupyterHub API",
	}

	// Flags available to everyone.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location. (default is .jhmon.{yaml,json,toml}")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "connect to the JupyterhHub with this authorization token.")
	rootCmd.PersistentFlags().StringVarP(&hubURL, "hub-url", "u", "",
		fmt.Sprintf("connect to the JupyterhHub at this URL. (default is%s)", defaultHubURL))

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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("jhmon")

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
		if viper.GetBool("verbose") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if viper.GetBool("verbose") {
			fmt.Printf("Error loading config file: %s - %v\n", viper.ConfigFileUsed(), err)
		}
	}

	// Set up a default connection
	initConnections()

}

func initConnections() {
	config.InitConnections(defaultHubURL)
	if rootCmd.PersistentFlags().Lookup("hub-url").Changed {
		config.UpdateDefaultHubURL(hubURL)
	}
	if rootCmd.PersistentFlags().Lookup("token").Changed {
		config.UpdateDefaultToken(token)
	}
}
