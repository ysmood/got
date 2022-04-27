package diff

import (
	"strings"
)

func neq(x, y Comparable) bool {
	return x.Hash() != y.Hash()
}

func eq(x, y Comparable) bool {
	return x.Hash() == y.Hash()
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// String interface
func (x Comparables) String() string {
	if len(x) == 0 {
		return ""
	}

	switch x[0].(type) {
	case Line:
		out := []string{}
		for _, c := range x {
			out = append(out, c.String())
		}
		return strings.Join(out, "\n")

	default:
		out := ""
		for _, c := range x {
			out += c.String()
		}
		return out
	}
}
