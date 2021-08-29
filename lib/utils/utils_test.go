package utils_test

import (
	"fmt"
	"testing"

	"github.com/ysmood/got/lib/utils"
)

func TestCompare(t *testing.T) {
	testCases := []struct {
		x interface{}
		y interface{}
		s interface{}
	}{
		{1, 1, 0.0},
		{1, 3.0, -2.0},
		{"b", "a", 1.0},
		{1, nil, -1.0},
	}
	for _, c := range testCases {
		t.Run(fmt.Sprintf("%v", c), func(t *testing.T) {
			s := utils.Compare(c.x, c.y)
			if s != c.s {
				t.Fail()
				t.Log("expect s to be", c.s, "but got", s)
			}
		})
	}
}
