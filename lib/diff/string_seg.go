package diff

import (
	"regexp"
	"strings"
)

// strSeg locates a substring inside a source string by byte offset and length.
// No per-segment copy: src is the single owner of the bytes.
type strSeg struct {
	src            string
	offset, length int
}

// text returns the segment content as a slice of src (no copy).
func (s strSeg) text() string {
	return s.src[s.offset : s.offset+s.length]
}

// Build writes the segment into sb.
func (s strSeg) Build(sb *strings.Builder) {
	sb.WriteString(s.src[s.offset : s.offset+s.length])
}

// seg wraps a whole string into a strSeg spanning all of it.
func seg(s string) strSeg {
	return strSeg{src: s, offset: 0, length: len(s)}
}

// Static literal segments shared across every token that uses them, so the
// backing string and segment are allocated once per package init rather than
// per token.
var (
	segNewline    = seg("\n")
	segSpace      = seg(" ")
	segEmpty      = seg("")
	segChunkStart = seg("@@ diff chunk @@")
)

type contextSplitKey struct{}

// SplitKey is the context key used to override the default word-splitting
// function consumed by TokenizeLine.
var SplitKey = contextSplitKey{}

// splitLines splits s on '\n' into line segments. A trailing '\n' yields an
// extra empty segment (to match line-oriented scanner semantics).
func splitLines(s string) []strSeg {
	if len(s) == 0 {
		return nil
	}

	n := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n++
		}
	}

	segments := make([]strSeg, 0, n)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			segments = append(segments, strSeg{src: s, offset: start, length: i - start})
			start = i + 1
		}
	}
	segments = append(segments, strSeg{src: s, offset: start, length: len(s) - start})
	return segments
}

var regWord = regexp.MustCompile(`(?s)` + // enable . to match newline
	`[[:alpha:]]{1,12}` + // match alphabets, length limit is 12
	`|[[:digit:]]{1,3}` + // match digits, length limit is 3
	`|.`) // match others as single-char words

var regRune = regexp.MustCompile(`(?s).`)

// Split splits a line into words. Short lines are split per-rune; longer
// lines are split into alphabetic, digit, or single-character tokens.
func Split(s string) []string {
	reg := regWord
	if len(s) <= 100 {
		reg = regRune
	}
	return reg.FindAllString(s, -1)
}
