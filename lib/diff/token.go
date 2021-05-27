package diff

import (
	"fmt"
	"strings"
)

// Type of token
type Type int

const (
	// LineNum
	LineNum Type = iota

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
)

// Token presents a symbol in diff layout
type Token struct {
	Type    Type
	Literal string
}

// Tokenize text block a and b into diff tokens.
func Tokenize(x, y string) []*Token {
	xls := NewText(x) // x lines
	yls := NewText(y) // y lines
	s := LCS(xls, yls)
	var i, j, k int

	ts := []*Token{}

	xNum, yNum, sNum := numFormat(xls, yls)

	for i < len(xls) && j < len(yls) && k < len(s) {
		xl := xls[i].(*Line)
		yl := yls[j].(*Line)
		l := s[k].(*Line)

		if !equal(xl, l) {
			ts = append(ts,
				&Token{LineNum, fmt.Sprintf(xNum, i+1)},
				&Token{DelSymbol, "- "},
				&Token{DelLine, string(xl.str) + "\n"})
			i++
		} else if !equal(yl, l) {
			ts = append(ts,
				&Token{LineNum, fmt.Sprintf(yNum, j+1)},
				&Token{AddSymbol, "+ "},
				&Token{AddLine, string(yl.str) + "\n"})
			j++
		} else {
			ts = append(ts,
				&Token{LineNum, fmt.Sprintf(sNum, i+1, j+1)},
				&Token{SameSymbol, "  "},
				&Token{SameLine, string(l.str) + "\n"})
			i, j, k = i+1, j+1, k+1
		}
	}

	return ts
}

func numFormat(x, y []Comparable) (string, string, string) {
	xl := len(fmt.Sprintf("%d", len(x)))
	yl := len(fmt.Sprintf("%d", len(y)))

	return fmt.Sprintf("%%0%dd "+strings.Repeat(" ", yl+1), xl),
		fmt.Sprintf(strings.Repeat(" ", xl)+" %%0%dd ", yl),
		fmt.Sprintf("%%0%dd %%0%dd ", xl, yl)
}
