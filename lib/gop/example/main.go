package main

import (
	"time"

	"github.com/ysmood/got/lib/gop"
)

func main() {
	val := map[string]interface{}{
		"bool":   true,
		"number": 1 + 1i,
		"bytes":  []byte{97, 98, 99},
		"lines":  "multiline string\nline two",
		"slice":  []interface{}{1, 2},
		"time":   time.Now(),
		"chan":   make(chan int, 1),
		"struct": struct{ test int32 }{
			test: 13,
		},
		"json": `{"a"   : 1}`,
		"func": func(int) int { return 0 },
	}
	val["slice"].([]interface{})[1] = val["slice"]

	_ = gop.P(val)
}
