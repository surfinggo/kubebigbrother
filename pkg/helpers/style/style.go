package style

import (
	"fmt"
	"github.com/muesli/termenv"
)

var termProfile = termenv.ColorProfile()

func Danger(f string, args ...interface{}) termenv.Style  { return Red(f, args...) }
func Info(f string, args ...interface{}) termenv.Style    { return Blue(f, args...) }
func Success(f string, args ...interface{}) termenv.Style { return Green(f, args...) }
func Warning(f string, args ...interface{}) termenv.Style { return Yellow(f, args...) }

func Red(f string, args ...interface{}) termenv.Style     { return Fg("#E88388", f, args...) }
func Green(f string, args ...interface{}) termenv.Style   { return Fg("#A8CC8C", f, args...) }
func Yellow(f string, args ...interface{}) termenv.Style  { return Fg("#DBAB79", f, args...) }
func Blue(f string, args ...interface{}) termenv.Style    { return Fg("#71BEF2", f, args...) }
func Magenta(f string, args ...interface{}) termenv.Style { return Fg("#D290E4", f, args...) }
func Cyan(f string, args ...interface{}) termenv.Style    { return Fg("#66C2CD", f, args...) }
func Gray(f string, args ...interface{}) termenv.Style    { return Fg("#B9BFCA", f, args...) }

func Fg(color string, f string, args ...interface{}) termenv.Style {
	return termenv.String(fmt.Sprintf(f, args...)).Foreground(termProfile.Color(color))
}

func Faint(f string, args ...interface{}) termenv.Style {
	return termenv.String(fmt.Sprintf(f, args...)).Faint()
}
