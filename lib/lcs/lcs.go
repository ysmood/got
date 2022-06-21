package lcs

import (
	"context"
)

// YadLCS returns the x index of each Comparable that are in the YadLCS between x and y.
// The complexity is O(M * log(L)), M is the number of char matches between x and y, L is the length of LCS.
// The worst memory complexity is O(M), but usually it's much less.
//
// The advantage of this algorithm is it's easy to understand and implement. It converts the LCS
// problem into problems that are familiar to us, such as LIS, binary-search, object-pool, etc., which give us
// more room to do the optimization for each streamline.
func (xs Sequence) YadLCS(ctx context.Context, ys Sequence) []int {
	o := xs.Occurrence(ys)
	r := result{}
	rest := len(ys)

	for _, xi := range o {
		if ctx.Err() != nil {
			break
		}

		from := len(r.list)
		for i := len(xi) - 1; i >= 0; i-- {
			from = r.add(from, xi[i], rest)
		}

		rest--
	}

	return r.lcs()
}

type node struct {
	x int
	p *node
	c int // pointer count
}

type result struct {
	list    []*node
	listLen int

	// reuse mem allocation of node
	pool    []*node
	poolLen int
}

func (r *result) new(x int, n *node) *node {
	if n != nil {
		n.c++
	}

	// reuse node if possible
	if r.poolLen > 0 {
		nw := r.pool[r.poolLen-1]
		r.pool = r.pool[:r.poolLen-1]
		r.poolLen--
		nw.x = x
		nw.p = n
		return nw
	}

	return &node{x, n, 0}
}

func (r *result) replace(i, x int, n *node) {
	// recycle nodes
	if m := r.list[i]; m.c == 0 {
		r.pool = append(r.pool, m)
		r.poolLen++

		for m = m.p; m != nil && m != n; m = m.p {
			m.c--
			if m.c == 0 {
				r.pool = append(r.pool, m)
				r.poolLen++
			} else {
				break
			}
		}
	}

	r.list[i] = r.new(x, n)
}

func (r *result) append(x int, n *node) {
	r.list = append(r.list, r.new(x, n))
	r.listLen++
}

func (r *result) add(from, x, rest int) int {
	l := r.listLen

	next, n := r.find(from, x)
	if n != nil && l-next < rest { // only when we have enough rest xs
		if next == l {
			r.append(x, n)
		} else if x < r.list[next].x {
			r.replace(next, x, n)
		}
		return next
	}

	if l == 0 {
		r.append(x, nil)
		return 1
	}

	if l-1 < rest && x < r.list[0].x {
		r.replace(0, x, nil)
	}

	return 0
}

// binary search to find the largest r[i].x that is smaller than x
func (r *result) find(from, x int) (int, *node) {
	var found *node
	for i, j := 0, from; i < j; {
		h := (i + j) >> 1
		n := r.list[h]
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

func (r *result) lcs() []int {
	l := r.listLen
	lcs := make([]int, l)

	if l == 0 {
		return lcs
	}

	for n, i := r.list[l-1], l-1; i >= 0; i-- {
		lcs[i] = n.x
		n = n.p
	}

	return lcs
}
