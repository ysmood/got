package diff

import "context"

// histMatch is a contiguous matching region between xs and ys.
type histMatch struct {
	xStart, yStart, length int
}

// intern assigns an int ID to each distinct string across xs and ys so a diff
// algorithm can compare and hash ints instead of full strings.
// Strings missing from the other side get a unique negative ID per position so
// they never compare equal to anything.
func intern(xs, ys []string) ([]int, []int) {
	ids := make(map[string]int, len(xs))
	for _, s := range xs {
		if _, ok := ids[s]; !ok {
			ids[s] = len(ids)
		}
	}

	xi := make([]int, len(xs))
	for i, s := range xs {
		xi[i] = ids[s]
	}

	yi := make([]int, len(ys))
	unique := -1
	for i, s := range ys {
		if id, ok := ids[s]; ok {
			yi[i] = id
		} else {
			yi[i] = unique
			unique--
		}
	}

	return xi, yi
}

// internLines is intern for line segments, resolved against their source texts.
// Map keys are string headers that slice into each segment's src — no line data
// is copied.
func internLines(xs, ys []strSeg) ([]int, []int) {
	ids := make(map[string]int, len(xs))
	for _, s := range xs {
		k := s.text()
		if _, ok := ids[k]; !ok {
			ids[k] = len(ids)
		}
	}

	xi := make([]int, len(xs))
	for i, s := range xs {
		xi[i] = ids[s.text()]
	}

	yi := make([]int, len(ys))
	unique := -1
	for i, s := range ys {
		k := s.text()
		if id, ok := ids[k]; ok {
			yi[i] = id
		} else {
			yi[i] = unique
			unique--
		}
	}

	return xi, yi
}

// histogramDiff returns contiguous matching regions between xs and ys in order.
// It uses a histogram-based algorithm: recursively pick the rarest element as
// an anchor, extend it to the longest common run, then diff the halves.
// If ctx is cancelled the partial result collected so far is returned.
func histogramDiff(ctx context.Context, xs, ys []int) []histMatch {
	var out []histMatch
	histRec(ctx, xs, ys, 0, len(xs), 0, len(ys), &out)
	return out
}

func histRec(ctx context.Context, xs, ys []int, ax, bx, ay, by int, out *[]histMatch) {
	if ax >= bx || ay >= by || ctx.Err() != nil {
		return
	}

	m := findHistAnchor(ctx, xs, ys, ax, bx, ay, by)
	if m.length == 0 {
		return
	}

	histRec(ctx, xs, ys, ax, m.xStart, ay, m.yStart, out)
	*out = append(*out, m)
	histRec(ctx, xs, ys, m.xStart+m.length, bx, m.yStart+m.length, by, out)
}

// findHistAnchor returns the best anchor region within xs[ax..bx] and ys[ay..by].
// Preference: lower occurrence count of the anchor element, then longer match,
// then lower xStart.
func findHistAnchor(ctx context.Context, xs, ys []int, ax, bx, ay, by int) histMatch {
	h := make(map[int][]int, bx-ax)
	for i := ax; i < bx; i++ {
		h[xs[i]] = append(h[xs[i]], i)
	}

	best := histMatch{}
	bestCount := 0

	for j := ay; j < by; j++ {
		if ctx.Err() != nil {
			return histMatch{}
		}

		positions, ok := h[ys[j]]
		if !ok {
			continue
		}
		count := len(positions)

		for _, xi := range positions {
			sx, sy := xi, j
			for sx > ax && sy > ay && xs[sx-1] == ys[sy-1] {
				sx--
				sy--
			}
			ex, ey := xi, j
			for ex+1 < bx && ey+1 < by && xs[ex+1] == ys[ey+1] {
				ex++
				ey++
			}
			length := ex - sx + 1

			if best.length == 0 ||
				count < bestCount ||
				(count == bestCount && length > best.length) ||
				(count == bestCount && length == best.length && sx < best.xStart) {
				best = histMatch{sx, sy, length}
				bestCount = count
			}
		}
	}

	return best
}
