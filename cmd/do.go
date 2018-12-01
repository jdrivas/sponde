package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var special bool
var doCmd, somethingCmd *cobra.Command

func buildDo(mode runMode) {

	doCmd = &cobra.Command{
		Use:   "do",
		Short: "Do something",
		Long:  "Do whatever comes after this.",
	}

	somethingCmd = &cobra.Command{
		Use:   "something",
		Short: "do something",
		Long:  "The first something that we can do",
		Run: func(cmd *cobra.Command, args []string) {
			s := ""
			if special {
				s = " special"
			}
			fmt.Printf("Doing something%s.\n", s)
		},
	}

	somethingCmd.Flags().BoolVarP(&special, "special", "s", false, "do something special")
	doCmd.AddCommand(somethingCmd)
	rootCmd.AddCommand(doCmd)
}
