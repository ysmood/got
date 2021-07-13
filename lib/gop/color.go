package gop

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Color type
type Color int

const (

	// Red type
	Red Color = iota + 31
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

	// None type
	None Color = -1
)

// ColorStr string
func ColorStr(c Color, s string) string {
	if c == -1 || !SupportsColor {
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
