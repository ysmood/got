package lcs

import (
	"context"
)

// YadLCS returns the x index of each Comparable that are in the YadLCS between x and y.
// The complexity is O(M * log(L)), M is the number of char matches between x and y, L is the length of YadLCS.
// The worst memory complexity is O(M), usually it's much less.
func (xs Sequence) YadLCS(ctx context.Context, ys Sequence) []int {
	h := xs.Histogram()
	r := result{}
	rest := len(ys)

	for _, c := range ys {
		if ctx.Err() != nil {
			break
		}

		if xi, has := h[c.String()]; has {
			from := len(r)
			for i := len(xi) - 1; i >= 0; i-- {
				from = r.add(from, xi[i], rest)
			}
		}

		rest--
	}

	return r.lcs()
}

type node struct {
	x int
	p *node
}

type result []*node

func (rp *result) add(from, x, rest int) int {
	r := *rp
	l := len(r)

	next, n := r.find(from, x)
	if n != nil && l-next < rest { // only when we have enough rest xs
		if next == l {
			*rp = append(r, &node{x, n})
		} else if x < r[next].x {
			r[next] = &node{x, n}
		}
		return next
	}

	if l == 0 {
		*rp = append(r, &node{x, nil})
		return 1
	}

	if l-1 < rest && x < r[0].x {
		r[0] = &node{x, nil}
	}

	return 0
}

// binary search to find the largest r[i].x that is smaller than x
func (rp result) find(from, x int) (int, *node) {
	var found *node
	for i, j := 0, from; i < j; {
		h := (i + j) >> 1
		n := rp[h]
		if n.x < x {
			from = h
			found = n
			i = h + 1
		} else {
			j = h
		}
	}
	return from + 1, found
}

func (rp *result) lcs() []int {
	r := *rp
	l := len(r)
	lcs := make([]int, l)

	if l == 0 {
		return lcs
	}

	for n, i := r[l-1], l-1; i >= 0; i-- {
		lcs[i] = n.x
		n = n.p
	}

	return lcs
}
