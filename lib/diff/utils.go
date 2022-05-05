package diff

import (
	"strings"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func neq(x, y Comparable) bool {
	return x.Hash() != y.Hash()
}

func eq(x, y Comparable) bool {
	return x.Hash() == y.Hash()
}

// String interface
func (x Sequence) String() string {
	if len(x) == 0 {
		return ""
	}

	if x[0].Hash() == x[0].String() {
		out := ""
		for _, c := range x {
			out += c.String()
		}
		return out
	}

	out := []string{}
	for _, c := range x {
		out = append(out, c.String())
	}
	return strings.Join(out, "\n")
}
