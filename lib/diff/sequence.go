package diff

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"unicode"
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
type Word string

// Hash interface
func (c Word) Hash() string {
	return string(c)
}

// String interface
func (c Word) String() string {
	return string(c)
}

// NewString from string
func NewString(s string) Comparables {
	runes := []rune(s)
	words := SplitSentence(runes)
	cs := make([]Comparable, len(words))
	for i, word := range words {
		cs[i] = Word(word)
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
	return Line{
		hash: string(sum[:]),
		str:  string(b),
	}
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

// Split sentence into words, for word-wise comparing.
// Any kind of splitter or non-English character is kept alone.
func SplitSentence(sentence []rune) []string {
	words := make([]string, 0)
	var prev []rune
	for _, r := range sentence {
		if notPartOfWord(r) {
			if len(prev) > 0 {
				words = safelyAppendWord(words, string(prev))
				prev = make([]rune, 0)
			}
			words = append(words, string(r))
		} else {
			prev = append(prev, r)
		}
	}
	if len(prev) > 0 {
		words = safelyAppendWord(words, string(prev))
	}
	return words
}

const MaxWordLen = 20

// try not append very long word
func safelyAppendWord(words []string, w string) []string {
	if len(w) <= MaxWordLen {
		return append(words, w)
	}
	parts := (len(w) + MaxWordLen - 1) / MaxWordLen
	for i := 0; i < (parts - 1); i++ {
		words = append(words, w[:MaxWordLen])
		w = w[MaxWordLen:]
	}
	if len(w) > 0 {
		words = append(words, w)
	}
	return words
}

// We assume that only alphabetic letters and digits can form word.
func notPartOfWord(r rune) bool {
	return !(inAlphabet(r) || unicode.IsDigit(r))
}

func inAlphabet(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}
