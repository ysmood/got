# Go Pretty Print Value

Make a random Go value human readable. The output format uses valid golang syntax, so you don't have to learn any new knowledge to understand the output.

## Features

- Uses valid golang syntax to print the data
- Correctly prints cyclic data
- Make rune, []byte, time, etc. data human readable
- No invisible char in string
- Color output with customizable theme
- Stable map output with sorted by keys
- Auto split multiline large string block
- Low-level API to extend the lib

## Usage

Usually, you only need to use `gop.P` function:

```go
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
        "struct": struct{ test int32 }{
            test: 13,
        },
    }
    val["slice"].([]interface{})[1] = val["slice"]

    gop.P(val)
}
```

The output will be:

```go
map[string]interface {}/* len=7 */{
    "bool": true,
    "bytes": []byte("abc")/* len=3 */,
    "lines": "" +
        "multiline string\n" +
        "line two"/* len=25 */,
    "number": 1+1i,
    "slice": []interface {}/* len=2 */{
        1,
        gop.Circular("slice").([]interface {}),
    },
    "struct": struct { test int32 }{
        test: int32(13),
    },
    "time": gop.Time("2022-03-16T21:38:22.135011+08:00"),
}
```
