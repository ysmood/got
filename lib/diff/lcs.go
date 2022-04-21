package diff

import (
	"bytes"
)

// LCS extends the standard lcs algorithm with dynamic programming and recursive common-lines reduction.
// The base algorithm we use is here: https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#LCS_function_defined.
// TODO: implement Patience Diff http://alfedenzo.livejournal.com/170301.html
func LCS(x, y []Comparable) []Comparable {
	var search func(i, j int) []Comparable
	mem := map[[2]int][]Comparable{}

	search = func(i, j int) []Comparable {
		k := [2]int{i, j}
		var lcs []Comparable
		var has bool
		if lcs, has = mem[k]; has {
			return lcs
		}

		var p, s []Comparable

		if i == 0 || j == 0 {
			lcs = []Comparable{}
		} else if x, y, p, s = reduce(x[:i], y[:j]); len(p) > 0 || len(s) > 0 {
			lcs = append(append(p, search(len(x), len(y))...), s...)
		} else {
			left, right := search(i, j-1), search(i-1, j)
			if len(left) > len(right) {
				lcs = left
			} else {
				lcs = right
			}
		}

		mem[k] = lcs
		return lcs
	}

	return search(len(x), len(y))
}

// If the beginning and ending are equal they must be in the LCS.
// This function scales down the problem via the first property:
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#First_property
func reduce(x, y []Comparable) (xs, ys, prefix, suffix []Comparable) {
	prefix, suffix = []Comparable{}, []Comparable{}
	for i := 0; i < len(x); i++ {
		if equal(x[i], y[i]) {
			prefix = append(prefix, x[i])
		} else {
			break
		}
	}

	x, y = x[len(prefix):], y[len(prefix):]

	for i, j := len(x)-1, len(y)-1; i >= 0 && j >= 0; {
		if equal(x[i], y[j]) {
			suffix = append([]Comparable{x[i]}, suffix...)
		} else {
			break
		}

		i--
		j--
	}

	xs, ys = x[:len(x)-len(suffix)], y[:len(y)-len(suffix)]

	return
}

func equal(x, y Comparable) bool {
	return bytes.Equal(x.Hash(), y.Hash())
}
