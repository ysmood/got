package diff

import (
	"context"
)

// LCS between x and y.
// This implementation converts the LCS problem into LIS sub problems without recursion.
// The memory complexity is O(x + y).
// The time complexity is similar with Myer's diff algorithm, but with more modularized steps, which allows further optimization easier,
// it doesn't use recursion, it's much easer to understand and implement.
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem
func (x Sequence) LCS(ctx context.Context, y Sequence) Sequence {
	left, right := x.Common(y)
	l, r, x, y := x[:left], x[len(x)-right:], x[left:len(x)-right], y[left:len(y)-right]

	return append(append(Sequence{}, l...), append(x.lcs(ctx, y), r...)...)
}

func (x Sequence) lcs(ctx context.Context, y Sequence) Sequence {
	x = x.Reduce(y)
	y = y.Reduce(x)

	o := x.Occurrence(y)
	lo := len(o)

	var lis []int          // longest increasing subsequence
	can := make([]int, lo) // candidate lis

	for i := 0; lo-i > len(lis) && ctx.Err() == nil; i++ { // only when the rest are more than lis
		l := 0
		can[0] = o[i][0]
		for j := i + 1; j < lo; j++ {
			oj := o[j]
			if gt, found := BTreeFindGreater(oj, can[l]); found {
				l++
				can[l] = gt
				goto next
			}

			break
		next:
		}

		if l >= len(lis) {
			lis = make([]int, l+1)
			copy(lis, can[:l+1])
		}
	}

	lcs := make(Sequence, len(lis))
	for i, j := range lis {
		lcs[i] = x[j]
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

// Occurrence histogram
type Occurrence [][]int

// Occurrence returns the position of each element of y in x.
func (x Sequence) Occurrence(y Sequence) Occurrence {
	m := make(Occurrence, len(y))
	h := x.Histogram()

	for i, c := range y {
		if indexes, has := h[c.Hash()]; has {
			m[i] = indexes
		}
	}

	return m
}
