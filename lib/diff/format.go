package diff

import (
	"github.com/ysmood/got/lib/gop"
)

// DefaultTheme colors for Sprint
var DefaultTheme = func(t Type) gop.Color {
	return gop.None
}

// NoTheme colors for Sprint
var NoTheme = func(t Type) gop.Color {
	return gop.None
}

// Diff x and y into a human readable string.
func Diff(x, y string) string {
	return Format(Tokenize(x, y), DefaultTheme)
}

// Narrow the context around each diff section to n lines.
func Narrow(n int, ts []*Token) []*Token {
	if n < 0 {
		n = 0
	}

	out := []*Token{}
	pivot := 0

	lines := ParseTokenLines(ts)
	hunks := ParseTokenHunks(lines)

	for _, h := range hunks {
		from, to := h.From(), h.To()
		if from-n >= 0 {
			from = from - n
		}
		if to+n < len(lines) {
			to = to + n
		}

		if pivot < from {
			out = append(out, lines[from:to]...)
		}

	}

	return out
}

// Format tokens into a human readable string
func Format(ts []*Token, theme func(Type) gop.Color) string {
	return ""
}
