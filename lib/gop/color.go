package gop

import (
	"fmt"
	"os"
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
	// Forground type
	Forground
	// Default type
	Default

	// BgBlack type
	BgBlack
	// BgRed type
	BgRed
	// BgGreen type
	BgGreen
	// BgYellow type
	BgYellow
	// BgBlue type
	BgBlue
	// BgMagenta type
	BgMagenta
	// BgCyan type
	BgCyan
	// BgWhite type
	BgWhite
	// Background type
	Background
	// BgDefault type
	BgDefault

	// None type
	None Color = -1
)

var regNewline = regexp.MustCompile(`\r?\n`)

// ColorStr string
func ColorStr(c Color, s string) string {
	if NoColor || c == None || !SupportsColor {
		return s
	}

	newline := regNewline.FindString(s)

	lines := regNewline.Split(s, -1)
	out := []string{}

	for _, l := range lines {
		out = append(out, fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, l))
	}

	return strings.Join(out, newline)
}

// SupportsColor returns true if current shell supports ANSI color
var SupportsColor = func() bool {
	b, _ := exec.Command("tput", "colors").CombinedOutput()
	n, _ := strconv.ParseInt(strings.TrimSpace(string(b)), 10, 32)
	return n > 0
}()

// NoColor respects https://no-color.org/
var NoColor = func() bool {
	_, has := os.LookupEnv("NO_COLOR")
	return has
}()

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var regStripColor = regexp.MustCompile(ansi)

// StripColor is copied from https://github.com/acarl005/stripansi
func StripColor(str string) string {
	return regStripColor.ReplaceAllString(str, "")
}
