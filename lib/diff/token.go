package diff

import (
	"fmt"
	"strings"
)

// Type of token
type Type int

const (
	// LineNum type
	LineNum Type = iota

	// ChunkStart type
	ChunkStart
	// ChunkEnd type
	ChunkEnd

	// SameSymbol type
	SameSymbol
	// SameLine type
	SameLine

	// AddSymbol type
	AddSymbol
	// AddLine type
	AddLine

	// DelSymbol typ
	DelSymbol
	// DelLine type
	DelLine

	// SameWords type
	SameWords
	// AddWords type
	AddWords
	// DelWords type
	DelWords

	// EmptyLine type
	EmptyLine
)

// Token presents a symbol in diff layout
type Token struct {
	Type    Type
	Literal string
}

// TokenizeText text block a and b into diff tokens.
func TokenizeText(x, y string) []*Token {
	xls := NewText(x) // x lines
	yls := NewText(y) // y lines
	s := LCS(xls, yls)

	ts := []*Token{}

	xNum, yNum, sNum := numFormat(xls, yls)
	chunkStarted := false

	for i, j, k := 0, 0, 0; i < len(xls) || j < len(yls); {
		if i < len(xls) && (k == len(s) || !equal(xls[i], s[k])) {
			ts, chunkStarted = tokenizeChunk(true, chunkStarted, ts)

			ts = append(ts,
				&Token{LineNum, fmt.Sprintf(xNum, i+1)},
				&Token{DelSymbol, "- "},
				&Token{DelLine, string(xls[i].(*Line).str) + "\n"})
			i++
		} else if j < len(yls) && (k == len(s) || !equal(yls[j], s[k])) {
			ts, chunkStarted = tokenizeChunk(true, chunkStarted, ts)

			ts = append(ts,
				&Token{LineNum, fmt.Sprintf(yNum, j+1)},
				&Token{AddSymbol, "+ "},
				&Token{AddLine, string(yls[j].(*Line).str) + "\n"})
			j++
		} else {
			ts, chunkStarted = tokenizeChunk(false, chunkStarted, ts)

			ts = append(ts,
				&Token{LineNum, fmt.Sprintf(sNum, i+1, j+1)},
				&Token{SameSymbol, "  "},
				&Token{SameLine, string(s[k].(*Line).str) + "\n"})
			i, j, k = i+1, j+1, k+1
		}
	}

	ts, _ = tokenizeChunk(false, chunkStarted, ts)

	return ts
}

// TokenizeLine two different lines
func TokenizeLine(x, y string) ([]*Token, []*Token) {
	xs := NewString(x)
	ys := NewString(y)

	s := LCS(xs, ys)

	xTokens := []*Token{}
	for i, j := 0, 0; i < len(xs); i++ {
		if j < len(s) && equal(xs[i], s[j]) {
			xTokens = append(xTokens, &Token{SameWords, string(s[j].(Char))})
			j++
		} else {
			xTokens = append(xTokens, &Token{DelWords, string(xs[i].(Char))})
		}
	}

	yTokens := []*Token{}
	for i, j := 0, 0; i < len(ys); i++ {
		if j < len(s) && equal(ys[i], s[j]) {
			yTokens = append(yTokens, &Token{SameWords, string(s[j].(Char))})
			j++
		} else {
			yTokens = append(yTokens, &Token{AddWords, string(ys[i].(Char))})
		}
	}

	return xTokens, yTokens
}

func tokenizeChunk(start bool, started bool, ts []*Token) ([]*Token, bool) {
	if start && !started {
		return append(ts, &Token{ChunkStart, ""}), true
	}
	if !start && started {
		return append(ts, &Token{ChunkEnd, ""}), false
	}
	return ts, started
}

func numFormat(x, y []Comparable) (string, string, string) {
	xl := len(fmt.Sprintf("%d", len(x)))
	yl := len(fmt.Sprintf("%d", len(y)))

	return fmt.Sprintf("%%0%dd "+strings.Repeat(" ", yl+1), xl),
		fmt.Sprintf(strings.Repeat(" ", xl)+" %%0%dd ", yl),
		fmt.Sprintf("%%0%dd %%0%dd ", xl, yl)
}
