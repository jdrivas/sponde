package cmd

import (
	"fmt"
	"strings"
)

func checkForEmptyString(s string) (r string) {
	r = s
	if s == "" {
		r = "<empty>"
	}
	return r
}

var trueValues = []string{"true", "True", "yes", "Yes"}
var falseValues = []string{"false", "False", "no", "No"}
var allValues = strings.Join(append(trueValues, falseValues...), ", ")

func truthyString(v string) (bool, error) {

	var err error
	for _, t := range trueValues {
		if v == t {
			return true, err
		}
	}

	for _, t := range falseValues {
		if v == t {
			return false, err
		}
	}
	return false, fmt.Errorf("couldn't make a truth value from \"%s\"; try one of: %s.", v, allValues)
}
