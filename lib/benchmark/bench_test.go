package benchmark

import (
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/ysmood/got"
	"github.com/ysmood/got/lib/lcs"
)

func BenchmarkLinesYad(b *testing.B) {
	// We use the same test file as https://github.com/sergi/go-diff/tree/master/testdata

	g := got.T(b)
	xs := lcs.NewLines(g.ReadFile("fixtures/speedtest1.txt").String())
	ys := lcs.NewLines(g.ReadFile("fixtures/speedtest2.txt").String())

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		xs.YadLCS(g.Context(), ys)
	}
}

func BenchmarkLinesGoogle(b *testing.B) {
	g := got.T(b)
	xs := g.ReadFile("fixtures/speedtest1.txt").String()
	ys := g.ReadFile("fixtures/speedtest2.txt").String()
	dmp := diffmatchpatch.New()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		dmp.DiffMain(xs, ys, true)
	}
}

func BenchmarkRandomYad(b *testing.B) {
	g := got.T(b)
	xs := lcs.NewChars(g.ReadFile("fixtures/rand_x.txt").String())
	ys := lcs.NewChars(g.ReadFile("fixtures/rand_y.txt").String())

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		xs.YadLCS(g.Context(), ys)
	}
}

func BenchmarkRandomGoogle(b *testing.B) {
	g := got.T(b)
	xs := g.ReadFile("fixtures/rand_x.txt").String()
	ys := g.ReadFile("fixtures/rand_y.txt").String()
	dmp := diffmatchpatch.New()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		dmp.DiffMain(xs, ys, false)
	}
}
