package diff

import (
	"bufio"
	"bytes"
	"hash/fnv"
	"regexp"
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
		hash := word
		if len(word) > 8 {
			h := fnv.New128()
			_, _ = h.Write([]byte(word))
			hash = string(h.Sum(nil))
		}
		cs[i] = &Element{hash, word}
	}
	return cs
}

// NewLines from string. It will split the s via newlines.
func NewLines(s string) Sequence {
	sc := bufio.NewScanner(bytes.NewBufferString(s))
	cs := []Comparable{}
	for i := 0; sc.Scan(); i++ {
		h := fnv.New128()
		_, _ = h.Write(sc.Bytes())
		cs = append(cs, &Element{string(h.Sum(nil)), sc.Text()})
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
	`|.` + // match others as single-char words
	``)

// RegRune to match a rune
var regRune = regexp.MustCompile(`(?s).`)

// SplitKey for context
var SplitKey = struct{}{}

// Split a line into words
func Split(s string) []string {
	var reg *regexp.Regexp
	if len(s) <= 100 {
		reg = regRune
	} else {
		reg = regWord
	}

	return reg.FindAllString(s, -1)
}
