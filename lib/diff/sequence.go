package diff

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"strings"
)

// Comparable interface
type Comparable interface {
	// Hash for comparison
	Hash() []byte
}

var _ Comparable = Char(0)

// Char is a rune
type Char rune

// Hash interface
func (c Char) Hash() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[:], uint32(c))
	return b
}

// String of Char list
type String []Comparable

// NewString from string
func NewString(s string) String {
	cs := make([]Comparable, len(s))
	for i, c := range s {
		cs[i] = Char(c)
	}
	return cs
}

// String interface
func (s String) String() string {
	out := ""
	for _, c := range s {
		out += string(c.(Char))
	}
	return out
}

var _ Comparable = &Line{}

// Line of a string for fast comparison.
type Line struct {
	hash []byte
	str  []byte
}

// NewLine from bytes
func NewLine(b []byte) *Line {
	// For testing, md5 should be sufficient
	sum := md5.Sum(b)
	return &Line{
		hash: sum[:],
		str:  b,
	}
}

// Hash interface
func (c *Line) Hash() []byte {
	return c.hash
}

// Text of Char list
type Text []Comparable

// NewText from string. It will split the s via newlines.
func NewText(s string) Text {
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

// String interface
func (s Text) String() string {
	out := []string{}
	for _, c := range s {
		out = append(out, string(c.(*Line).str))
	}
	return strings.Join(out, "\n")
}
