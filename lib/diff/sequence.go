package diff

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"regexp"

	"github.com/ysmood/got/lib/gop"
)

// Comparables list
type Comparables []Comparable

// Comparable interface
type Comparable interface {
	// Hash for fast comparison
	Hash() string
	// String returns the full content
	String() string
}

// Word block
type Word struct {
	str  string
	hash string
}

// Hash interface
func (w Word) Hash() string {
	return w.hash
}

// String interface
func (w Word) String() string {
	return w.str
}

// NewWords from string
func NewWords(split func(string) []string, s string) Comparables {
	words := split(s)
	cs := make([]Comparable, len(words))
	for i, word := range words {
		if gop.RegANSI.MatchString(word) {
			cs[i] = Word{word, ""}
		} else {
			cs[i] = Word{word, word}
		}
	}
	return cs
}

// Line of a string for fast comparison.
type Line struct {
	str  string
	hash string
}

// NewLine from bytes
func NewLine(b []byte) Line {
	// For testing, md5 should be sufficient
	sum := md5.Sum(b)
	return Line{string(b), string(sum[:])}
}

// Hash interface
func (c Line) Hash() string {
	return c.hash
}

// String interface
func (c Line) String() string {
	return c.str
}

// NewText from string. It will split the s via newlines.
func NewText(s string) Comparables {
	sc := bufio.NewScanner(bytes.NewBufferString(s))
	cs := []Comparable{}
	for sc.Scan() {
		cs = append(cs, NewLine(sc.Bytes()))
	}

	if len(s) > 0 && s[len(s)-1] == '\n' {
		cs = append(cs, NewLine([]byte{}))
	}

	return cs
}

// RegWord to match a word
var regWord = regexp.MustCompile(`(?s)` + // enable . to match newline
	`[[:alpha:]]{1,12}` + // match alphabets, length limit is 12
	`|[[:digit:]]{1,3}` + // match digits, length limit is 3
	`|` + gop.RegANSI.String() + // match terminal color escape sequences
	`|.` + // match others as single-char words
	``)

// RegRune to match a rune
var regRune = regexp.MustCompile(`(?s)` + // enable . to match newline
	gop.RegANSI.String() + // match terminal color escape sequences
	`|.` + // match others as single-char words
	``)

// SplitKey for context
var SplitKey = struct{}{}

// Split a line into words
func Split(s string) []string {
	var reg *regexp.Regexp
	if len(gop.StripANSI(s)) <= 100 {
		reg = regRune
	} else {
		reg = regWord
	}

	return reg.FindAllString(s, -1)
}
