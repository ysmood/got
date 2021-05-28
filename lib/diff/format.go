package diff

import (
	"github.com/ysmood/got/lib/gop"
)

// DefaultTheme colors for Sprint
var DefaultTheme = func(t Type) gop.Color {
	switch t {
	case AddSymbol, AddLine:
		return gop.Green
	case DelSymbol, DelLine:
		return gop.Red
	}
	return gop.None
}

// NoTheme colors for Sprint
var NoTheme = func(t Type) gop.Color {
	return gop.None
}

// Diff x and y into a human readable string.
func Diff(x, y string) string {
	return Format(TokenizeText(x, y), DefaultTheme)
}

// Format tokens into a human readable string
func Format(ts []*Token, theme func(Type) gop.Color) string {
	out := ""

	for _, t := range ts {
		out += gop.ColorStr(theme(t.Type), t.Literal)
	}

	return out
}
