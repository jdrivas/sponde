package term

import (
	"fmt"
	"github.com/fatih/color"
)

// Use this with github.com/juju/ansi term to get a TabWriter that works with color.

type ColorSprintfFunc func(string, ...interface{}) string

var (
	// Text Formatting
	Title    = color.New(color.FgBlack).SprintfFunc()
	SubTitle = color.New(color.FgHiBlack).SprintfFunc()
	Text     = color.New(color.FgHiBlack).SprintfFunc()

	// Semantic Formatting
	Info      = color.New(color.FgBlack).SprintfFunc()
	Highlight = color.New(color.FgGreen).SprintfFunc()
	Success   = color.New(color.FgGreen).SprintfFunc()
	Warn      = color.New(color.FgYellow).SprintfFunc()
	Fail      = color.New(color.FgRed).SprintfFunc()
	Alert     = color.New(color.FgRed).SprintfFunc()
)

func Error(err error) string {
	return (fmt.Sprintf("%s %s", Title("Error: "), Fail("%v", err)))
}
