package diff

import (
	"bytes"
)

// LCS extends the standard lcs algorithm with dynamic programming and recursive common-lines reduction.
// The base algorithm we use is here: https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#LCS_function_defined.
// TODO: implement Patience Diff http://alfedenzo.livejournal.com/170301.html
func LCS(x, y []Comparable) []Comparable {
	var search func(xi, xj, yi, yj int) []Comparable
	mem := map[[4]int][]Comparable{}

	search = func(xi, xj, yi, yj int) []Comparable {
		k := [4]int{xi, xj, yi, yj}
		var lcs []Comparable
		var has bool
		if lcs, has = mem[k]; has {
			return lcs
		}

		if (xj-xi)*(yj-yi) == 0 {
			lcs = []Comparable{}
		} else if l, r := Common(x[xi:xj], y[yi:yj]); l+r > 0 {
			lcs = append([]Comparable{}, x[xi:xi+l]...)
			lcs = append(lcs, search(xi+l, xj-r, yi+l, yj-r)...)
			lcs = append(lcs, x[xj-r:xj]...)
		} else {
			left, right := search(xi, xj, yi, yj-1), search(xi, xj-1, yi, yj)
			if len(left) > len(right) {
				lcs = left
			} else {
				lcs = right
			}
		}

		mem[k] = lcs
		return lcs
	}

	return search(0, len(x), 0, len(y))
}

// Common returns the common prefix and suffix between x and y.
// This function scales down the problem via the first property:
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#First_property
func Common(x, y []Comparable) (left, right int) {
	l := min(len(x), len(y))
	for ; left < l; left++ {
		if !equal(x[left], y[left]) {
			break
		}
	}

	lx, ly := len(x), len(y)
	l = min(lx-left, ly-left)
	for ; right < l; right++ {
		if !equal(x[lx-right-1], y[ly-right-1]) {
			break
		}
	}

	return
}

func equal(x, y Comparable) bool {
	return bytes.Equal(x.Hash(), y.Hash())
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
