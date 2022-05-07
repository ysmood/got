package diff

import (
	"context"
)

// LCS between x and y.
// This implementation converts the LCS problem into LIS sub problems without recursion.
// The memory complexity is O(x.Occurrence(y)).
// The time complexicy is O(x.Occurrence(y).Complexity()).
// The time complexicy is similar with Myer's diff algorithm, but with more modulized steps, which allows further optimization easier.
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
	p := make([]int, len(o))
	n := o.Complexity()

	lis := NewLIS(len(o), func(i int) int {
		return o[i][p[i]]
	})

	var longest int
	var longestI int
	for i := 0; i < n && ctx.Err() == nil; i++ {
		o.Permutate(p, i)

		l := lis.Length()
		if l > longest {
			longestI = i
			longest = l
		}
	}

	p = make([]int, len(o))
	o.Permutate(p, longestI)
	s := lis.Get()

	lcs := make(Sequence, longest)
	for i := 0; i < longest; i++ {
		lcs[i] = x[s[i]]
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

// Complexity to find the LCS in m
func (o Occurrence) Complexity() int {
	if len(o) == 0 {
		return 0
	}

	n := 1
	for _, i := range o {
		n *= len(i)
	}
	return n
}

// Permutate p with i
func (o Occurrence) Permutate(p []int, i int) {
	for j := 0; i > 0; j++ {
		p[j] = i % len(o[j])
		i = i / len(o[j])
	}
}

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

// LIS https://en.wikipedia.org/wiki/Longest_increasing_subsequence
type LIS struct {
	p   []int
	m   []int
	len int
	x   func(int) int
}

// NewLIS helper
func NewLIS(len int, x func(int) int) *LIS {
	return &LIS{
		p:   make([]int, len),
		m:   make([]int, len+1),
		len: len,
		x:   x,
	}
}

// Length of LIS
func (lis *LIS) Length() int {
	l := 0
	for i := 0; i < lis.len; i++ {
		lo := 1
		hi := l + 1
		for lo < hi {
			mid := lo + (hi-lo)/2
			if lis.x(lis.m[mid]) < lis.x(i) {
				lo = mid + 1
			} else {
				hi = mid
			}
		}

		newL := lo

		lis.p[i] = lis.m[newL-1]
		lis.m[newL] = i

		if newL > l {
			l = newL
		}
	}
	return l
}

// Get the LIS
func (lis *LIS) Get() []int {
	l := lis.Length()
	s := make([]int, l)
	k := lis.m[l]
	for i := l - 1; i >= 0; i-- {
		s[i] = lis.x(k)
		k = lis.p[k]
	}
	return s
}
