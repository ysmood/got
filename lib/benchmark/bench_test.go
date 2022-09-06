package benchmark

import (
	"context"
	"crypto/rand"
	"strings"
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

	xs, ys := []rune(x), []rune(y)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		dmp.DiffMainRunes(xs, ys, false)
	}
}

func BenchmarkRandomMyers(b *testing.B) {
	xs, ys := split(x), split(y)

	for i := 0; i < b.N; i++ {
		_ = myers.Diff(xs, ys)
	}
}

func randStr(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return string(b)
}

func split(text string) []string {
	return strings.Split(text, "")
}
