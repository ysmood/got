package diff

import "context"

// LCS extends the standard lcs algorithm with dynamic programming and recursive line-reduction.
// The base algorithm we use is here: https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#LCS_function_defined.
// TODO: implement Patience Diff http://alfedenzo.livejournal.com/170301.html
func (x Comparables) LCS(ctx context.Context, y Comparables) Comparables {
	var search func(xi, xj, yi, yj int) Comparables
	mem := map[[4]int]Comparables{}

	search = func(xi, xj, yi, yj int) Comparables {
		if ctx.Err() != nil {
			return Comparables{}
		}

		k := [4]int{xi, xj, yi, yj}
		var lcs Comparables
		var has bool
		if lcs, has = mem[k]; has {
			return lcs
		}

		if (xj-xi)*(yj-yi) == 0 {
			lcs = Comparables{}
		} else if l, r := x[xi:xj].Common(y[yi:yj]); l+r > 0 {
			lcs = append(Comparables{}, x[xi:xi+l]...)
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

	x = x.Reduce(y)
	y = y.Reduce(x)

	return search(0, len(x), 0, len(y))
}

// Common returns the common prefix and suffix between x and y.
// This function scales down the problem via the first property:
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#First_property
func (x Comparables) Common(y Comparables) (left, right int) {
	l := min(len(x), len(y))
	for ; left < l && eq(x[left], y[left]); left++ {
	}

	lx, ly := len(x), len(y)
	l = min(lx-left, ly-left)
	for ; right < l && eq(x[lx-right-1], y[ly-right-1]); right++ {
	}

	return
}

// Reduce Comparables from x that doesn't exist in y
func (x Comparables) Reduce(y Comparables) Comparables {
	rest := Comparables{}
	h := y.Histogram()
	for _, c := range x {
		if _, has := h[c.Hash()]; has {
			rest = append(rest, c)
		}
	}
	return rest
}

// Histogram of each Comparable
func (x Comparables) Histogram() map[string]int {
	his := map[string]int{}
	for _, c := range x {
		his[c.Hash()]++
	}
	return his
}
