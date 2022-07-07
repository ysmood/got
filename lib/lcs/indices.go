package lcs

// Indices is the index list of items in xs that forms the LCS between xs and ys.
type Indices []int

// Compress indices into negative indices
// The negative integer whose absolute value represents the length of the continuous subsequence.
// Such as:
//     Indices{1,2,3,4,5,10,20,21,22}
// will be compressed into:
//     Indices{1,-4,10,20,-2}
func (idx Indices) Compress() Indices {
	l := len(idx)

	if l == 0 {
		return idx
	}

	c := Indices{idx[0]}
	for i, j := 1, 1; i < l; i++ {
		p := c[j-1]
		if p < 0 {
			pp := c[j-2]
			if idx[i] == pp-p+1 {
				c[j-1]--
				continue
			}
		} else if idx[i] == p+1 {
			c = append(c, -1)
			j++
			continue
		}
		c = append(c, idx[i])
		j++
	}

	return c
}

// Decompress negative indices
func (idx Indices) Decompress() Indices {
	s := Indices{}
	for i, ix := range idx {
		if ix >= 0 {
			s = append(s, ix)
		} else {
			p := idx[i-1]
			for j := 0; j < -ix; j++ {
				p++
				s = append(s, p)
			}
		}
	}
	return s
}
