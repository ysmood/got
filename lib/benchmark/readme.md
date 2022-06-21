# Overview

```txt
$ go test -benchmem -bench .
goos: darwin
goarch: amd64
pkg: github.com/ysmood/got/lib/benchmark
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkLinesYad-12        	   20043	     56490 ns/op	   62905 B/op	     333 allocs/op
BenchmarkLinesGoogle-12     	      24	  47257804 ns/op	 4895688 B/op	   37644 allocs/op
BenchmarkRandomYad-12       	     506	   2360644 ns/op	  468221 B/op	   16931 allocs/op
BenchmarkRandomGoogle-12    	     174	   6910260 ns/op	  680122 B/op	   11572 allocs/op
PASS
ok  	github.com/ysmood/got/lib/benchmark	6.396s
```

YadLCS is very good at line based diff (BenchmarkLines), it's about 800x faster than [Google Myer's algorithm](https://github.com/sergi/go-diff/blob/849d7ebc9716f43ec1295e9bc00e5c8cffef3d9f/diffmatchpatch/diff.go#L5-L7), and it uses 100x less memory.
YadLCS 3x faster on random long string.

It will use more memory when the comparable candidates are similar, such as when we compare `ababab` and `bababa`,
they both only have 2 types of comparable which is `a` and `b`. But it's not common in real-world text based comparison.
