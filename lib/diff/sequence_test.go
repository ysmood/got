package diff_test

import (
	"strings"
	"testing"

	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/diff"
	"github.com/ysmood/got/lib/gop"
)

func TestSplit(t *testing.T) {
	g := got.T(t)

	check := func(in string, expect ...string) {
		g.Helper()
		out := diff.Split(strings.Repeat("\t", 100) + in)[100:]
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

	check(gop.S("test", gop.Red),
		gop.Red.Set, "test", gop.Red.Unset)
}
