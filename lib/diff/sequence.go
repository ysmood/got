package diff

import (
	"bufio"
	"bytes"
	"crypto/md5"
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

// Char is a rune
type Char string

// Hash interface
func (c Char) Hash() string {
	return string(c)
}

// String interface
func (c Char) String() string {
	return string(c)
}

// NewString from string
func NewString(s string) Comparables {
	cs := make([]Comparable, len(s))
	for i, c := range s {
		cs[i] = Char(c)
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
