package diff

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"regexp"

	"github.com/ysmood/got/lib/gop"
)

// Sequence list
type Sequence []Comparable

// Comparable interface
type Comparable interface {
	// Hash for comparison
	Hash() string
	// String returns the full content
	String() string
}

// Element of a line or a word
type Element struct {
	hash string
	str  string
}

// Hash for comparison
func (e *Element) Hash() string {
	return e.hash
}

// String returns the full content
func (e *Element) String() string {
	return e.str
}

// NewWords from string
func NewWords(words []string) Sequence {
	cs := make([]Comparable, len(words))
	for i, word := range words {
		if gop.RegANSI.MatchString(word) {
			cs[i] = &Element{"", word}
		} else {
			cs[i] = &Element{word, word}
		}
	}
	return cs
}

// NewText from string. It will split the s via newlines.
func NewText(s string) Sequence {
	sc := bufio.NewScanner(bytes.NewBufferString(s))
	cs := []Comparable{}
	for i := 0; sc.Scan(); i++ {
		sum := md5.Sum(sc.Bytes())
		cs = append(cs, &Element{string(sum[:]), sc.Text()})
	}

	if len(s) > 0 && s[len(s)-1] == '\n' {
		cs = append(cs, &Element{"", ""})
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
