# go-ulid

Universally Unique Lexicographically Sortable Identifier (ULID) generator in Go.

## Why ULID?

Check out [the spec of ULID](https://github.com/ulid/spec).

## About this implementation

- It uses cryptographically secure random numbers by default.
- high performance.

vs [oklog/ulid](https://github.com/oklog/ulid) v2.1.1:

```plain
goos: darwin
goarch: arm64
pkg: github.com/shogo82148/go-ulid/_benchmark
cpu: Apple M1 Pro
BenchmarkOklogMake
BenchmarkOklogMake-10      	 6851478	       157.5 ns/op	      16 B/op	       1 allocs/op
BenchmarkShogoMake
BenchmarkShogoMake-10      	 8778632	       136.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkOklogString
BenchmarkOklogString-10    	49443757	        23.36 ns/op	      32 B/op	       1 allocs/op
BenchmarkShogoString
BenchmarkShogoString-10    	53509219	        21.78 ns/op	      32 B/op	       1 allocs/op
BenchmarkOklogParse
BenchmarkOklogParse-10     	75034783	        15.56 ns/op	       0 B/op	       0 allocs/op
BenchmarkShogoParse
BenchmarkShogoParse-10     	128026156	         9.382 ns/op	       0 B/op	       0 allocs/op
PASS
coverage: [no statements]
ok  	github.com/shogo82148/go-ulid/_benchmark	7.173s
```

## Prior Art

- [oklog/ulid](https://github.com/oklog/ulid)
- [darccio/go-ulid](https://github.com/darccio/go-ulid)
