package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/gop"
)

func TestNewString(t *testing.T) {
	g := got.T(t)

	check := func(in string, expect ...string) {
		g.Helper()
		out := []string{}
		in = strings.Repeat("\t", 100) + in
		for _, w := range diff.NewWords(diff.Split, in) {
			out = append(out, w.String())
		}
		g.Eq(out[100:], expect)
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

	check(gop.S("test", gop.Red),
		gop.Red.Set, "test", gop.Red.Unset)
}
