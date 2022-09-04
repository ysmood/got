# Overview

```txt
$ go test -benchmem -bench .
goos: darwin
goarch: amd64
pkg: github.com/ysmood/got/lib/benchmark
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkRandomYad-8               44662             24634 ns/op           15147 B/op        234 allocs/op
BenchmarkRandomGoogle-8            19030             57792 ns/op           42208 B/op        776 allocs/op
BenchmarkRandomMyers-8             17997             66763 ns/op          379044 B/op        335 allocs/op
PASS
ok  	github.com/ysmood/got/lib/benchmark	3.886s
```

YadLCS is faster and uses less memory than [Google Myer's algorithm](https://github.com/sergi/go-diff/blob/849d7ebc9716f43ec1295e9bc00e5c8cffef3d9f/diffmatchpatch/diff.go#L5-L7) when the item histogram is large, it's common in line based diff in large text context.
