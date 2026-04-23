package diff

import (
	"bufio"
	"bytes"
	"regexp"
)

type contextSplitKey struct{}

// SplitKey is the context key used to override the default word-splitting
// function consumed by TokenizeLine.
var SplitKey = contextSplitKey{}

// splitLines splits s on newlines, preserving a trailing empty line when s
// ends with '\n'. Each returned line has no trailing newline.
func splitLines(s string) []string {
	sc := bufio.NewScanner(bytes.NewBufferString(s))
	lines := []string{}
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if len(s) > 0 && s[len(s)-1] == '\n' {
		lines = append(lines, "")
	}
	return lines
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
