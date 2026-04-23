package diff

import (
	"context"
	"strconv"
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

// Token is a lazily-rendered piece of diff output.
type Token interface {
	Type() Type
	Build(sb *strings.Builder)
}

// Text renders a token into a string via its Build method.
func Text(t Token) string {
	var sb strings.Builder
	t.Build(&sb)
	return sb.String()
}

// SegToken carries a single strSeg tagged with a Type.
type SegToken struct {
	T Type
	s strSeg
}

// Type of the token.
func (x SegToken) Type() Type { return x.T }

// Build writes the segment into sb.
func (x SegToken) Build(sb *strings.Builder) { x.s.Build(sb) }

// ConcatToken holds a run of strSegs sharing one Type. Produced by mergeRuns
// so that merging is zero-copy: parts are concatenated only at Build time.
type ConcatToken struct {
	T     Type
	parts []strSeg
}

// Type of the token.
func (x ConcatToken) Type() Type { return x.T }

// Build writes each part into sb in order.
func (x ConcatToken) Build(sb *strings.Builder) {
	for _, p := range x.parts {
		p.Build(sb)
	}
}

// DelGutter renders the deleted-line gutter: "<N pad Xl> <Yl+1 spaces>-".
type DelGutter struct{ N, Xl, Yl int }

// Type of the token.
func (DelGutter) Type() Type { return DelSymbol }

// Build writes the gutter into sb.
func (g DelGutter) Build(sb *strings.Builder) {
	appendPad(sb, g.N, g.Xl)
	sb.WriteByte(' ')
	for i := 0; i <= g.Yl; i++ {
		sb.WriteByte(' ')
	}
	sb.WriteByte('-')
}

// AddGutter renders the added-line gutter: "<Xl+1 spaces><N pad Yl> +".
type AddGutter struct{ N, Xl, Yl int }

// Type of the token.
func (AddGutter) Type() Type { return AddSymbol }

// Build writes the gutter into sb.
func (g AddGutter) Build(sb *strings.Builder) {
	for i := 0; i <= g.Xl; i++ {
		sb.WriteByte(' ')
	}
	appendPad(sb, g.N, g.Yl)
	sb.WriteByte(' ')
	sb.WriteByte('+')
}

// SameGutter renders the unchanged-line gutter: "<A pad Xl> <B pad Yl>  ".
type SameGutter struct{ A, B, Xl, Yl int }

// Type of the token.
func (SameGutter) Type() Type { return SameSymbol }

// Build writes the gutter into sb.
func (g SameGutter) Build(sb *strings.Builder) {
	appendPad(sb, g.A, g.Xl)
	sb.WriteByte(' ')
	appendPad(sb, g.B, g.Yl)
	sb.WriteByte(' ')
	sb.WriteByte(' ')
}

// digits returns the number of base-10 digits of n for n >= 0.
func digits(n int) int {
	d := 1
	for n >= 10 {
		d++
		n /= 10
	}
	return d
}

// appendPad writes n in base 10, zero-padded to at least width digits.
func appendPad(sb *strings.Builder, n, width int) {
	for i := digits(n); i < width; i++ {
		sb.WriteByte('0')
	}
	sb.WriteString(strconv.Itoa(n))
}

// TokenizeText text block a and b into diff tokens.
func TokenizeText(ctx context.Context, x, y string) []Token {
	xls := splitLines(x)
	yls := splitLines(y)

	xi, yi := internLines(xls, yls)
	matches := histogramDiff(ctx, xi, yi)

	xl := digits(len(xls))
	yl := digits(len(yls))

	ts := []Token{}

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
				DelGutter{N: i + 1, Xl: xl, Yl: yl},
				SegToken{T: Space, s: segSpace},
				SegToken{T: DelLine, s: xls[i]},
				SegToken{T: Newline, s: segNewline})
			i++
		} else if j < yStop {
			ts = append(ts,
				AddGutter{N: j + 1, Xl: xl, Yl: yl},
				SegToken{T: Space, s: segSpace},
				SegToken{T: AddLine, s: yls[j]},
				SegToken{T: Newline, s: segNewline})
			j++
		} else if hasMatch {
			m := matches[mi]
			for k := 0; k < m.length; k++ {
				ts = append(ts,
					SameGutter{A: m.xStart + k + 1, B: m.yStart + k + 1, Xl: xl, Yl: yl},
					SegToken{T: Space, s: segSpace},
					SegToken{T: SameLine, s: xls[m.xStart+k]},
					SegToken{T: Newline, s: segNewline})
			}
			i = m.xStart + m.length
			j = m.yStart + m.length
			mi++
		}
	}

	return ts
}

// TokenizeLine two different lines
func TokenizeLine(ctx context.Context, x, y string) ([]Token, []Token) {
	split := Split
	val := ctx.Value(SplitKey)
	if val != nil {
		split = val.(func(string) []string)
	}

	xs := split(x)
	ys := split(y)

	xi, yi := intern(xs, ys)
	matches := myersDiff(ctx, xi, yi)

	xTokens := []Token{}
	yTokens := []Token{}

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
			xTokens = append(xTokens, SegToken{T: DelWords, s: seg(xs[i])})
			i++
		} else if j < yStop {
			yTokens = append(yTokens, SegToken{T: AddWords, s: seg(ys[j])})
			j++
		} else if hasMatch {
			m := matches[mi]
			for k := 0; k < m.length; k++ {
				xTokens = append(xTokens, SegToken{T: SameWords, s: seg(xs[m.xStart+k])})
				yTokens = append(yTokens, SegToken{T: SameWords, s: seg(ys[m.yStart+k])})
			}
			i = m.xStart + m.length
			j = m.yStart + m.length
			mi++
		}
	}

	return mergeRuns(xTokens), mergeRuns(yTokens)
}

// mergeRuns collapses each maximal run of same-Type tokens into a ConcatToken
// that holds the run's strSegs. Concatenation is deferred until Build time.
func mergeRuns(ts []Token) []Token {
	if len(ts) < 2 {
		return ts
	}

	out := make([]Token, 0, len(ts))
	for i := 0; i < len(ts); {
		start := i
		for i+1 < len(ts) && ts[i+1].Type() == ts[start].Type() {
			i++
		}
		if i == start {
			out = append(out, ts[start])
		} else {
			parts := make([]strSeg, 0, i-start+1)
			for k := start; k <= i; k++ {
				parts = append(parts, ts[k].(SegToken).s)
			}
			out = append(out, ConcatToken{T: ts[start].Type(), parts: parts})
		}
		i++
	}
	return out
}
