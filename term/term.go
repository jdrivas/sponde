package term

import (
	"fmt"

	"github.com/mgutz/ansi"
)

// Holds semantic color definitions.
type colorField struct {
	null, reset, title, header, text, info, fail string
}

var (
	// Basic Colors
	defaultColor = fmt.Sprintf(ansi.ColorCode("default"))
	defaultEmph  = fmt.Sprintf(ansi.ColorCode("default+b"))
	emphBlue     = fmt.Sprintf(ansi.ColorCode("blue+b"))
	emphRed      = fmt.Sprintf(ansi.ColorCode("red+b"))

	// Default prompt colors
	pColor = colorField{
		null:  fmt.Sprintf("%s", "\x00\x00\x00\x00\x00\x00\x00"),
		reset: fmt.Sprintf(ansi.ColorCode("reset")),
		title: emphBlue,
		info:  emphBlue,
		fail:  emphRed,
	}

	// Default Display Colors
	displayColor = colorField{
		null:   fmt.Sprintf("%s", "\x00\x00\x00\x00\x00\x00\x00"),
		reset:  fmt.Sprintf(ansi.ColorCode("reset")),
		title:  defaultEmph,
		header: defaultEmph,
		text:   defaultColor,
		info:   defaultColor,
		fail:   emphRed,
	}
)

// Header returns a string wrapped in ANSI colors for header.
func Header(s string) string {
	return fmt.Sprintf("%s%s%s", displayColor.header, s, displayColor.reset)
	// return s
}

func Text(s string) string {
	return fmt.Sprintf("%s%s%s", displayColor.text, s, displayColor.reset)
}
