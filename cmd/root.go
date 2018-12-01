// Copyright © 2018 David Rivas

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd, listCmd, describeCmd, interactiveCmd *cobra.Command
	cfgFile, token                                string
	verbose                                       bool
)

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

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List elements of a collection.",
		Long:  "Provides a short description of each element of a collection.",
	}

	describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "Provides a description of something.",
		Long:  "Provides a long description of an object.",
	}

	// Other Commands.
	if mode != interactive {
		rootCmd.AddCommand(interactiveCmd)
	}
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
		Use:   "cobra_test",
		Short: "Test application to take cobra for a spin",
		Long: `
		A simple application that will serve as a place to try thing out in 
		Cobra and probably Viper too.`,
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location. (defaul is ./.jhmon")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Describe what is happening as its happening.")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "secert token for connecting to server.")
	// This means that the values are obtained from viper when they are used.
	// This is in contradistinction say, from Traitlets which bind values
	// to an object and manage flags, config etc.
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

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
			fmt.Printf("Error loading config: %v\n", err)
		}
	}

}
