// +build !windows

package gop

import (
	"fmt"
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
	if c == -1 {
		return s
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}
