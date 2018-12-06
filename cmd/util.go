package cmd

import (
	"fmt"

	t "github.com/jdrivas/sponde/term"
)

func cmdError(e error) {
	fmt.Printf("Error: %s\n", t.Fail(e.Error()))
}

func checkForEmptyString(s string) (r string) {
	r = s
	if s == "" {
		r = "<empty>"
	}
	return r
}
