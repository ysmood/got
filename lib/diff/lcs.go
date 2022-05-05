package diff

import (
	"context"
)

// LCS between x and y.
// This implementation converts the LCS problem into LIS sub problems.
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func (x Sequence) LCS(ctx context.Context, y Sequence) Sequence {
	left, right := x.Common(y)
	l, r, x, y := x[:left], x[len(x)-right:], x[left:len(x)-right], y[left:len(y)-right]

	return append(append(Sequence{}, l...), append(x.findLCS(ctx, y), r...)...)
}

func (x Sequence) findLCS(ctx context.Context, y Sequence) Sequence {
	x = x.Reduce(y)
	y = y.Reduce(x)
	m := x.Occurrence(y)

	l := len(m)
	s := make([]int, l)
	p := make([]int, l)
	var longest []int
	for l > 0 && ctx.Err() == nil {
		for i := 0; i < l; i++ {
			s[i] = m[i][p[i]]
		}

		lis := LIS(s)
		if len(lis) > len(longest) {
			longest = lis
		}

		p[0]++
		for i := 0; i < l; i++ {
			if p[i] < len(m[i]) {
				break
			} else {
				p[i] = 0
				if i+1 == l {
					goto end
				}
				p[i+1]++
			}
		}
	}

end:

	lcs := make(Sequence, len(longest))
	for i, index := range longest {
		lcs[i] = x[index]
	}
	return lcs
}

// Common returns the common prefix and suffix between x and y.
// This function scales down the problem via the first property:
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#First_property
func (x Sequence) Common(y Sequence) (left, right int) {
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
func (x Sequence) Reduce(y Sequence) Sequence {
	rest := Sequence{}
	h := y.Histogram()
	for _, c := range x {
		if _, has := h[c.Hash()]; has {
			rest = append(rest, c)
		}
	}
	return rest
}

// Histogram of a Comparables
type Histogram map[string][]int

// Histogram of each Comparable
func (x Sequence) Histogram() Histogram {
	h := Histogram{}
	for i, c := range x {
		if _, has := h[c.Hash()]; !has {
			h[c.Hash()] = []int{}
		}
		h[c.Hash()] = append(h[c.Hash()], i)
	}
	return h
}

// IsSubsequenceOf returns true if x is a subsequence of y
func (x Sequence) IsSubsequenceOf(y Sequence) bool {
	for i, j := 0, 0; i < len(x); i++ {
		for {
			if j >= len(y) {
				return false
			}
			if eq(x[i], y[j]) {
				j++
				break
			}
			j++
		}
	}

	return true
}

// Occurrence returns the position of each element of x in y.
func (x Sequence) Occurrence(y Sequence) [][]int {
	m := make([][]int, len(y))
	h := x.Histogram()

	for i, c := range y {
		if indexes, has := h[c.Hash()]; has {
			m[i] = indexes
		}
	}

	return m
}

// LIS returns the longest increasing subsequence of s.
// https://en.wikipedia.org/wiki/Longest_increasing_subsequence
func LIS(x []int) []int {
	p := make([]int, len(x))
	m := make([]int, len(x)+1)

	l := 0
	for i := range x {
		// Binary search for the largest positive j ≤ L
		// such that X[M[j]] < X[i]
		lo := 1
		hi := l + 1
		for lo < hi {
			mid := lo + (hi-lo)/2
			if x[m[mid]] < x[i] {
				lo = mid + 1
			} else {
				hi = mid
			}
		}

		// After searching, lo is 1 greater than the
		// length of the longest prefix of X[i]
		newL := lo

		// The predecessor of X[i] is the last index of
		// the subsequence of length newL-1
		p[i] = m[newL-1]
		m[newL] = i

		if newL > l {
			// If we found a subsequence longer than any we've
			// found yet, update L
			l = newL
		}
	}

	// Reconstruct the longest increasing subsequence
	s := make([]int, l)
	k := m[l]
	for i := l - 1; i >= 0; i-- {
		s[i] = x[k]
		k = p[k]
	}

	return s
}
