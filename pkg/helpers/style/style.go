package style

import (
	"fmt"
	"github.com/muesli/termenv"
)

const (
	Red     = "#E88388"
	Green   = "#A8CC8C"
	Yellow  = "#DBAB79"
	Blue    = "#71BEF2"
	Magenta = "#D290E4"
	Cyan    = "#66C2CD"
	Gray    = "#B9BFCA"
	Info    = Blue
	Success = Green
	Warning = Yellow
	Danger  = Red
)

var termProfile = termenv.ColorProfile()

// Fg build foreground color style
func Fg(color string, f string, args ...interface{}) termenv.Style {
	return termenv.String(fmt.Sprintf(f, args...)).Foreground(termProfile.Color(color))
}

// Faint build faint style
func Faint(f string, args ...interface{}) termenv.Style {
	return termenv.String(fmt.Sprintf(f, args...)).Faint()
}
