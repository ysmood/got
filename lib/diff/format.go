package diff

import (
	"github.com/ysmood/got/lib/gop"
)

// Theme for diff
type Theme func(t Type) gop.Style

// ThemeDefault colors for Sprint
var ThemeDefault = func(t Type) gop.Style {
	switch t {
	case AddSymbol, AddWords:
		return gop.BgGreen
	case DelSymbol, DelWords:
		return gop.BgRed
	case ChunkStart:
		return gop.BgMagenta
	}
	return gop.None
}

// ThemeNone colors for Sprint
var ThemeNone = func(t Type) gop.Style {
	return gop.None
}

// Diff x and y into a human readable string.
func Diff(x, y string) string {
	return Format(Tokenize(x, y), ThemeDefault)
}

// Tokenize x and y into diff tokens with diff words and narrow chunks.
func Tokenize(x, y string) []*Token {
	ts := TokenizeText(x, y)
	lines := ParseTokenLines(ts)
	lines = Narrow(3, lines)
	ChunkWords(lines)
	return SpreadTokenLines(lines)
}

// Format tokens into a human readable string
func Format(ts []*Token, theme func(Type) gop.Style) string {
	out := ""

	for _, t := range ts {
		s := t.Literal
		out += gop.Stylize(theme(t.Type), s)
	}

	return gop.FixNestedStyle(out)
}

// Narrow the context around each diff section to n lines.
func Narrow(n int, lines []*TokenLine) []*TokenLine {
	if n < 0 {
		n = 0
	}

	keep := map[int]bool{}
	for i, l := range lines {
		switch l.Type {
		case AddSymbol, DelSymbol:
			for j := max(i-n, 0); j <= i+n && j < len(lines); j++ {
				keep[j] = true
			}
		}
	}

	out := []*TokenLine{}
	for i, l := range lines {
		if !keep[i] {
			continue
		}

		if _, has := keep[i-1]; !has {
			ts := []*Token{{ChunkStart, "@@ diff chunk @@"}, {Newline, "\n"}}
			out = append(out, &TokenLine{ChunkStart, ts})
		}

		out = append(out, l)

		if _, has := keep[i+1]; !has {
			ts := []*Token{{ChunkEnd, ""}, {Newline, "\n"}}
			out = append(out, &TokenLine{ChunkEnd, ts})
		}
	}

	return out
}

// ChunkWords with words
func ChunkWords(lines []*TokenLine) {
	delLines := []*TokenLine{}
	addLines := []*TokenLine{}

	df := func() {
		if len(delLines) == 0 || len(delLines) != len(addLines) {
			return
		}

		for i := 0; i < len(delLines); i++ {
			d := delLines[i]
			a := addLines[i]

			dts, ats := TokenizeLine(d.Tokens[2].Literal, a.Tokens[2].Literal)
			d.Tokens = append(d.Tokens[0:2], append(dts, d.Tokens[3:]...)...)
			a.Tokens = append(a.Tokens[0:2], append(ats, a.Tokens[3:]...)...)
		}

		delLines = []*TokenLine{}
		addLines = []*TokenLine{}
	}

	for _, l := range lines {
		switch l.Type {
		case DelSymbol:
			delLines = append(delLines, l)
		case AddSymbol:
			addLines = append(addLines, l)
		default:
			df()
		}
	}

	df()
}
