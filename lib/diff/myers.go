package diff

import "context"

// myersDiff returns contiguous matching regions between xs and ys in order,
// computed by Myers' O(ND) shortest-edit-script algorithm. Well suited for
// short sequences (e.g. words/chars in a line) where there is little or no
// "rarity" signal for histogram diff to exploit.
// If ctx is cancelled before the shortest path is found, nil is returned.
func myersDiff(ctx context.Context, xs, ys []int) []histMatch {
	n, m := len(xs), len(ys)
	if n == 0 || m == 0 {
		return nil
	}

	maxD := n + m
	offset := maxD
	v := make([]int, 2*maxD+1)

	trace := make([][]int, 0, maxD+1)

	dEnd := -1
	for d := 0; d <= maxD && dEnd < 0; d++ {
		if ctx.Err() != nil {
			return nil
		}

		snap := make([]int, 2*maxD+1)
		copy(snap, v)
		trace = append(trace, snap)

		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[offset+k-1] < v[offset+k+1]) {
				x = v[offset+k+1]
			} else {
				x = v[offset+k-1] + 1
			}
			y := x - k

			for x < n && y < m && xs[x] == ys[y] {
				x++
				y++
			}
			v[offset+k] = x

			if x >= n && y >= m {
				dEnd = d
				break
			}
		}
	}

	// Backtrace, collecting matches in reverse.
	pairsX := make([]int, 0, n)
	pairsY := make([]int, 0, n)
	x, y := n, m
	for d := dEnd; d > 0; d-- {
		prev := trace[d]
		k := x - y
		var prevK int
		if k == -d || (k != d && prev[offset+k-1] < prev[offset+k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}
		prevX := prev[offset+prevK]
		prevY := prevX - prevK

		for x > prevX && y > prevY {
			pairsX = append(pairsX, x-1)
			pairsY = append(pairsY, y-1)
			x--
			y--
		}
		x = prevX
		y = prevY
	}
	for x > 0 && y > 0 {
		pairsX = append(pairsX, x-1)
		pairsY = append(pairsY, y-1)
		x--
		y--
	}

	// Pairs are reversed; walk from end to group consecutive pairs into regions.
	var out []histMatch
	for i := len(pairsX) - 1; i >= 0; {
		sx, sy := pairsX[i], pairsY[i]
		length := 1
		for i-1 >= 0 && pairsX[i-1] == pairsX[i]+1 && pairsY[i-1] == pairsY[i]+1 {
			i--
			length++
		}
		out = append(out, histMatch{xStart: sx, yStart: sy, length: length})
		i--
	}

	return out
}
