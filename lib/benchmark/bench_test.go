package benchmark

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/ysmood/got/lib/benchmark/myers"
	"github.com/ysmood/got/lib/lcs"
)

var x = randStr(100)
var y = randStr(100)

func BenchmarkRandomYad(b *testing.B) {
	c := context.Background()

	xs := lcs.NewChars(x)
	ys := lcs.NewChars(y)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		xs.YadLCS(c, ys)
	}
}

func BenchmarkRandomGoogle(b *testing.B) {
	dmp := diffmatchpatch.New()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		dmp.DiffMain(x, y, false)
	}
}

func BenchmarkRandomMyers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = myers.Diff(x, y)
	}
}

func randStr(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return string(b)
}
