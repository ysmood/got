package benchmark

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/lcs"
)

func BenchmarkYad(b *testing.B) {
	g := got.T(b)
	xs := lcs.NewChars(g.ReadFile("fixtures/rand_x.txt").String())
	ys := lcs.NewChars(g.ReadFile("fixtures/rand_y.txt").String())

	g.Log(len(xs.Sub(xs.YadLCS(g.Context(), ys))))

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		xs.YadLCS(g.Context(), ys)
	}
}

func BenchmarkGoogle(b *testing.B) {
	g := got.T(b)
	xs := g.ReadFile("fixtures/rand_x.txt").String()
	ys := g.ReadFile("fixtures/rand_y.txt").String()
	dmp := diffmatchpatch.New()

	df := dmp.DiffMain(xs, ys, false)
	l := ""
	for _, d := range df {
		if d.Type == diffmatchpatch.DiffEqual {
			l += d.Text
		}
	}
	c := lcs.NewChars(l)
	g.Log(len(c), c.IsSubsequenceOf(lcs.NewChars(xs)))

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		dmp.DiffMain(xs, ys, false)
	}
}
