package diff_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ysmood/got/lib/diff"
)

// buildLines returns n newline-separated lines of the form "prefix N".
func buildLines(prefix string, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%s %d\n", prefix, i)
	}
	return sb.String()
}

// mutate returns x with every stepth line replaced by a differing line.
func mutate(x string, step int) string {
	lines := strings.Split(x, "\n")
	for i := 0; i < len(lines); i += step {
		lines[i] = lines[i] + " changed"
	}
	return strings.Join(lines, "\n")
}

func BenchmarkDiff(b *testing.B) {
	cases := []struct {
		name string
		x, y string
	}{
		{"Identical/1000", buildLines("line", 1000), buildLines("line", 1000)},
		{"FewChanges/1000", buildLines("line", 1000), mutate(buildLines("line", 1000), 50)},
		{"ManyChanges/1000", buildLines("line", 1000), mutate(buildLines("line", 1000), 3)},
		{"Disjoint/1000", buildLines("old", 1000), buildLines("new", 1000)},
		{"Small", "the quick brown fox\njumps over\nthe lazy dog\n", "the quick red fox\njumps over\nthe happy dog\n"},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = diff.Diff(c.x, c.y)
			}
		})
	}
}
