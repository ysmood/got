package diff_test

import (
	"testing"

	"github.com/ysmood/got/lib/diff"
)

func TestNewString(t *testing.T) {
	g := setup(t)
	testCases := []struct {
		in   string
		want []string
	}{{
		"find a place to eat 热干面",
		[]string{"find", " ", "a", " ", "place", " ", "to", " ", "eat", " ", "热", "干", "面"},
	}, {
		"	as.Equal(arr, arr)",
		[]string{"	", "as", ".", "Equal", "(", "arr", ",", " ", "arr", ")"},
	}, {
		"English-Words紧接着中文",
		[]string{"English", "-", "Words", "紧", "接", "着", "中", "文"},
	}, {
		"WordVeryVeryVeryVeryVeryVeryVerylong",
		[]string{"WordVeryVeryVeryVery", "VeryVeryVerylong"},
	}}

	for _, c := range testCases {
		output := diff.NewString(c.in)
		g.Eq(len(output), len(c.want))
		for i, got := range output {
			g.Eq(got.String(), c.want[i])
		}
	}
}
