package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ysmood/got/lib/utils"
)

func TestCompare(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		x interface{}
		y interface{}
		s interface{}
	}{
		{1, 1, 0.0},
		{1, 3.0, -2.0},
		{"b", "a", 1.0},
		{1, nil, -1.0},
		{now.Add(time.Second), now, float64(time.Second)},
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
