# Overview

```txt
$ go test -benchmem -bench .
goos: darwin
goarch: amd64
pkg: github.com/ysmood/got/lib/benchmark
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkRandomYad-12       	     459	   2443507 ns/op	  746207 B/op	   44452 allocs/op
BenchmarkRandomGoogle-12    	     171	   6946139 ns/op	  680123 B/op	   11572 allocs/op
BenchmarkLinesYad-12        	   18100	     60055 ns/op	   61062 B/op	     564 allocs/op
BenchmarkLinesGoogle-12     	      24	  47627760 ns/op	 4895833 B/op	   37645 allocs/op
PASS
ok  	github.com/ysmood/got/lib/benchmark	6.781s
```

YadLCS is very good at line based diff, it's about 800x faster than Google's diff lib.
YadLCS 3x faster on random long string.

It will use more memory than Google's when the comparable candidates are similar, such as X is `ababab`, Y is `bababa`,
they both only have `a` and `b`. But it's not common in real-world text based comparison.
