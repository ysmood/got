// TODO: we should optimize the algorithm based on other modern papers, but for easy to understand
// we use the most basic one for now: https://en.wikipedia.org/wiki/Longest_common_subsequence_problem
// A good candidate is http://www.xmailserver.org/diff2.pdf .

package diff

import (
	"bytes"
)

// LCS computes the lcs of x and y
func LCS(x, y []Comparable) []Comparable {
	x, y = sort(x, y)

	x, y, prefix, suffix := scaleDown(x, y)

	c := lcsTable(x, y)
	return append(append(
		prefix,
		lcsFromTable(c, x, y, len(x)-1, len(y)-1)...),
		suffix...,
	)
}

// If the beginning and ending are equal they must be in the LCS.
// This function scales down the problem via the first property:
// https://en.wikipedia.org/wiki/Longest_common_subsequence_problem#First_property
func scaleDown(x, y []Comparable) (xs, ys, p, s []Comparable) {
	p, s = []Comparable{}, []Comparable{}
	for i := 0; i < len(x); i++ {
		if equal(x[i], y[i]) {
			p = append(p, x[i])
		} else {
			break
		}
	}

	x, y = x[len(p):], y[len(p):]

	for i, j := len(x)-1, len(y)-1; i >= 0 && j >= 0; {
		if equal(x[i], y[j]) {
			s = append([]Comparable{x[i]}, s...)
		} else {
			break
		}

		i--
		j--
	}

	xs, ys = x[:len(x)-len(s)], y[:len(y)-len(s)]

	return
}

func sort(x, y []Comparable) (small, large []Comparable) {
	for i, j := 0, 0; i < len(x) && j < len(y); {
		xh, yh := x[i].Hash(), y[j].Hash()
		if bytes.Compare(xh[:], yh[:]) > 0 {
			return y, x
		}
		i++
		j++
	}
	return x, y
}

// Computes the lcs table for x and y. Time complexity is O(w*h)
func lcsTable(x, y []Comparable) [][]int {
	w := len(x)
	h := len(y)
	c := make([][]int, w)

	for i := 0; i < w; i++ {
		row := make([]int, h)
		for j := 0; j < h; j++ {
			var top, left int
			if i != 0 {
				top = c[i-1][j]
			}
			if j != 0 {
				left = row[j-1]
			}

			maxAdjoin := max(top, left)

			if equal(x[i], y[j]) {
				row[j] = maxAdjoin + 1
			} else {
				row[j] = maxAdjoin
			}
		}
		c[i] = row
	}

	return c
}

// Backtrack the longest common subsequence for x and y via table c.
func lcsFromTable(c [][]int, x, y []Comparable, i, j int) []Comparable {
	if i < 0 || j < 0 {
		return []Comparable{}
	}

	if equal(x[i], y[j]) {
		return append(lcsFromTable(c, x, y, i-1, j-1), x[i])
	}

	var top, left int
	if i != 0 {
		top = c[i-1][j]
	}
	if j != 0 {
		left = c[i][j-1]
	}

	if top > left {
		return lcsFromTable(c, x, y, i-1, j)
	}
	return lcsFromTable(c, x, y, i, j-1)
}

func equal(x, y Comparable) bool {
	return bytes.Equal(x.Hash(), y.Hash())
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
