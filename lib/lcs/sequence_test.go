package lcs_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/lcs"
)

func TestSplit(t *testing.T) {
	g := got.T(t)

	g.Len(lcs.Split(""), 0)

	check := func(in string, expect ...string) {
		g.Helper()
		out := lcs.Split(strings.Repeat("\t", 100) + in)[100:]
		g.Eq(out, expect)
	}

	check("find a place to eat 热干面",
		"find", " ", "a", " ", "place", " ", "to", " ", "eat", " ", "热", "干", "面")

	check("	as.Equal(arr, arr) test",
		"	", "as", ".", "Equal", "(", "arr", ",", " ", "arr", ")", " ", "test")

	check("English-Words紧接着中文",
		"English", "-", "Words", "紧", "接", "着", "中", "文")

	check("123456test12345",
		"123", "456", "test", "123", "45")

	check("WordVeryVeryVeryVeryVeryVeryVerylong",
		"WordVeryVery", "VeryVeryVery", "VeryVerylong")
}

func TestIsSubsequenceOf(t *testing.T) {
	g := got.T(t)

	y := lcs.NewChars("abc")

	g.True(lcs.NewChars("ab").IsSubsequenceOf(y))
	g.True(lcs.NewChars("ac").IsSubsequenceOf(y))
	g.True(lcs.NewChars("bc").IsSubsequenceOf(y))
	g.False(lcs.NewChars("cb").IsSubsequenceOf(y))
	g.False(lcs.NewChars("ba").IsSubsequenceOf(y))
	g.False(lcs.NewChars("ca").IsSubsequenceOf(y))
}

func TestNew(t *testing.T) {
	g := setup(t)
	g.Len(lcs.NewLines("a"), 1)
	g.Len(lcs.NewLines("a\n"), 2)
	g.Len(lcs.NewLines("a\n\n"), 3)
	g.Len(lcs.NewLines("\na"), 2)
	g.Eq(lcs.NewLines("\nabc\nabc").String(), "\nabc\nabc")

	g.Len(lcs.NewWords([]string{"a", "b"}), 2)
}
