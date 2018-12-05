package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jdrivas/sponde/config"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Holds semantic color definitions.
type promptColor struct {
	null, reset, title, info, fail string
}

var (
	// Handy designations
	emphBlue = fmt.Sprintf(ansi.ColorCode("blue+b"))
	emphRed  = fmt.Sprintf(ansi.ColorCode("red+b"))

	// Default prompt colors
	pColor = promptColor{
		null:  fmt.Sprintf("%s", "\x00\x00\x00\x00\x00\x00\x00"),
		reset: fmt.Sprintf(ansi.ColorCode("reset")),
		title: emphBlue,
		info:  emphBlue,
		fail:  emphRed,
	}

	// Type exit instead of just control-d, Note: We actually os.exit() here.
	// Which eans no post-processing of any kind including simply falling through
	// to the orignial Execute command.
	// if this is a problem, move the definition of the promptLoop moreCommands up
	// to module scope and set it to false in the Run function directly below.
	exitCmd = &cobra.Command{
		Use:   "exit",
		Short: "Exit from the application",
		Long:  "Stop reading input lines and terminate the application.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nGoodbye and thank you.\n")
			os.Exit(0)
		},
	}

	verboseCmd = &cobra.Command{
		Use:   "verbose",
		Short: "Toggle verbose mode and print status.",
		Long:  "Toggle verbose, verbose will print out detailed status as its happening.",
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set("verbose", !viper.GetBool("verbose"))
			vs := "Off"
			if viper.GetBool("verbose") {
				vs = "On"
			}
			fmt.Printf("Verbose is %s\n", vs)
		},
	}

	debugCmd = &cobra.Command{
		Use:   "debug",
		Short: "Toggle debug mode and print status.",
		Long:  "Toggle debug, verbose will print out detailed status as its happening.",
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set("debug", !viper.GetBool("debug"))
			vs := "Off"
			if viper.GetBool("debug") {
				vs = "On"
			}
			fmt.Printf("Debug is %s\n", vs)
		},
	}
)

// Parse the line and execute the command
func doICommand(line string) (err error) {
	// Start from clean slate each time.
	rootCmd.ResetCommands()
	buildRoot(interactive)

	rootCmd.AddCommand(exitCmd)
	rootCmd.AddCommand(verboseCmd)
	rootCmd.AddCommand(debugCmd)

	rootCmd.SetArgs(strings.Split(line, " "))
	err = rootCmd.Execute()
	return err
}

func promptLoop(process func(string) error) (err error) {

	for moreCommands := true; moreCommands; {
		hubURL := config.GetHubURL()
		connName := config.GetConnectionName()
		token = config.GetSafeToken()
		prompt := fmt.Sprintf("%ssponde [%s - %s %s]:%s ", pColor.title, connName, hubURL, token, pColor.reset)
		line, err := readline.Line(prompt)
		if err == io.EOF {
			moreCommands = false
		} else if err != nil {
			fmt.Printf("%sReadline Error: %s%s\n", pColor.fail, err, pColor.reset)
		} else {
			readline.AddHistory(line)
			err = process(line)
			if err == io.EOF {
				moreCommands = false
			}
		}
	}
	return nil
}

// DoInteractive sets up a readline loop that reads and executes comands.
func DoInteractive() {
	readline.SetHistoryPath("./.sponde_history")
	xICommand := func(line string) (err error) { return doICommand(line) }
	err := promptLoop(xICommand)
	if err != nil {
		fmt.Printf("%sError exiting prompter: %s%s\n", pColor.fail, err, pColor.reset)
	}
}
