package diff

import (
	"context"
	"fmt"
	"strings"
)

// Type of token
type Type int

const (
	// Newline type
	Newline Type = iota
	// Space type
	Space

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
func TokenizeText(ctx context.Context, x, y string) []*Token {
	xls := splitLines(x)
	yls := splitLines(y)

	xi, yi := intern(xls, yls)
	matches := histogramDiff(ctx, xi, yi)

	ts := []*Token{}
	xNum, yNum, sNum := numFormat(len(xls), len(yls))

	i, j, mi := 0, 0, 0
	for i < len(xls) || j < len(yls) {
		xStop, yStop := len(xls), len(yls)
		hasMatch := false
		if mi < len(matches) {
			xStop = matches[mi].xStart
			yStop = matches[mi].yStart
			hasMatch = true
		}

		if i < xStop {
			ts = append(ts,
				&Token{DelSymbol, fmt.Sprintf(xNum, i+1) + "-"},
				&Token{Space, " "},
				&Token{DelLine, xls[i]},
				&Token{Newline, "\n"})
			i++
		} else if j < yStop {
			ts = append(ts,
				&Token{AddSymbol, fmt.Sprintf(yNum, j+1) + "+"},
				&Token{Space, " "},
				&Token{AddLine, yls[j]},
				&Token{Newline, "\n"})
			j++
		} else if hasMatch {
			m := matches[mi]
			for k := 0; k < m.length; k++ {
				ts = append(ts,
					&Token{SameSymbol, fmt.Sprintf(sNum, m.xStart+k+1, m.yStart+k+1) + " "},
					&Token{Space, " "},
					&Token{SameLine, xls[m.xStart+k] + "\n"})
			}
			i = m.xStart + m.length
			j = m.yStart + m.length
			mi++
		} else {
			break
		}
	}

	return ts
}

// TokenizeLine two different lines
func TokenizeLine(ctx context.Context, x, y string) ([]*Token, []*Token) {
	split := Split
	val := ctx.Value(SplitKey)
	if val != nil {
		split = val.(func(string) []string)
	}

	xs := split(x)
	ys := split(y)

	xi, yi := intern(xs, ys)
	matches := myersDiff(ctx, xi, yi)

	xTokens := []*Token{}
	yTokens := []*Token{}

	merge := func(ts []*Token) []*Token {
		last := len(ts) - 1
		if last > 0 && ts[last].Type == ts[last-1].Type {
			ts[last-1].Literal += ts[last].Literal
			ts = ts[:last]
		}
		return ts
	}

	i, j, mi := 0, 0, 0
	for i < len(xs) || j < len(ys) {
		xStop, yStop := len(xs), len(ys)
		hasMatch := false
		if mi < len(matches) {
			xStop = matches[mi].xStart
			yStop = matches[mi].yStart
			hasMatch = true
		}

		if i < xStop {
			xTokens = append(xTokens, &Token{DelWords, xs[i]})
			xTokens = merge(xTokens)
			i++
		} else if j < yStop {
			yTokens = append(yTokens, &Token{AddWords, ys[j]})
			yTokens = merge(yTokens)
			j++
		} else if hasMatch {
			m := matches[mi]
			for k := 0; k < m.length; k++ {
				xTokens = append(xTokens, &Token{SameWords, xs[m.xStart+k]})
				yTokens = append(yTokens, &Token{SameWords, ys[m.yStart+k]})
				xTokens = merge(xTokens)
				yTokens = merge(yTokens)
			}
			i = m.xStart + m.length
			j = m.yStart + m.length
			mi++
		} else {
			break
		}
	}

	return xTokens, yTokens
}

func numFormat(xLen, yLen int) (string, string, string) {
	xl := len(fmt.Sprintf("%d", xLen))
	yl := len(fmt.Sprintf("%d", yLen))

	return fmt.Sprintf("%%0%dd "+strings.Repeat(" ", yl+1), xl),
		fmt.Sprintf(strings.Repeat(" ", xl)+" %%0%dd ", yl),
		fmt.Sprintf("%%0%dd %%0%dd ", xl, yl)
}
