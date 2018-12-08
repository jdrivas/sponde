package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jdrivas/sponde/config"
	t "github.com/jdrivas/sponde/term"

	// "github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

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
		token := config.GetSafeToken(true, false)
		spacer := ""
		if token != "" {
			spacer = " "
		}
		status := statusDisplay()
		prompt := fmt.Sprintf("%s [%s%s %s]: ", t.Title("sponde"), t.Info(status), t.Highlight(connName), t.SubTitle("%s%s%s", hubURL, spacer, token))
		line, err := readline.Line(prompt)
		if err == io.EOF {
			moreCommands = false
		} else if err != nil {
			fmt.Printf("Readline Error: %s\n", t.Fail(err.Error()))
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

// Yes, I'm sure there's some kind of []rune
// thing to do here instead.
func statusDisplay() (s string) {
	if config.Verbose() {
		s = fmt.Sprintf("%s%s", s, "v")
	}
	if config.Debug() {
		s = fmt.Sprintf("%s%s", s, "d")
	}
	if len(s) > 0 {
		s = fmt.Sprintf("%s%s", s, " ")
	}
	return s
}

// DoInteractive sets up a readline loop that reads and executes comands.
func DoInteractive() {
	readline.SetHistoryPath("./.sponde_history")
	xICommand := func(line string) (err error) { return doICommand(line) }
	err := promptLoop(xICommand)
	if err != nil {
		fmt.Printf("Error exiting prompter: %s\n", t.Fail(err.Error()))
	}
}
