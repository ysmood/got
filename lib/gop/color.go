package gop

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Color type
type Color int

const (
	// Black type
	Black Color = iota + 30
	// Red type
	Red
	// Green type
	Green
	// Yellow type
	Yellow
	// Blue type
	Blue
	// Magenta type
	Magenta
	// Cyan type
	Cyan
	// White type
	White

	// None type
	None Color = -1
)

// ColorStr string
func ColorStr(c Color, s string) string {
	if c == None || !SupportsColor {
		return s
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}

// SupportsColor returns true if current shell supports ANSI color
var SupportsColor = func() bool {
	b, _ := exec.Command("tput", "colors").CombinedOutput()
	n, _ := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 32)
	return n > 0
}()

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var regStripColor = regexp.MustCompile(ansi)

// StripColor is copied from https://github.com/acarl005/stripansi
func StripColor(str string) string {
	return regStripColor.ReplaceAllString(str, "")
}
