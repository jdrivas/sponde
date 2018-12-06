package term

import (
	"github.com/fatih/color"
)

// Use this with github.com/juju/ansi term to get a TabWriter that works with color.
var (
	// Text Formatting
	Title    = color.New(color.FgBlack).SprintfFunc()
	SubTitle = color.New(color.FgHiBlack).SprintfFunc()
	Text     = color.New(color.FgHiBlack).SprintfFunc()

	// Semantic Formatting
	Info      = color.New(color.FgBlack).SprintfFunc()
	Highlight = color.New(color.FgGreen).SprintfFunc()
	Success   = color.New(color.FgGreen).SprintfFunc()
	Fail      = color.New(color.FgRed).SprintfFunc()
)
