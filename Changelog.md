
# v1.0.0

MVP
- minimal html UI
- JSON API
- Sqlite backend

```
goos: linux
goarch: amd64
pkg: github.com/bchabanne/urlshortener
cpu: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
BenchmarkInsertHdd-8                   7         155672037 ns/op
BenchmarkReadHdd-8                 10424            109758 ns/op
BenchmarkReadWriteHdd-8              169          16628722 ns/op
BenchmarkInsertMemory-8            21298             56252 ns/op
BenchmarkReadMemory-8              13527             85455 ns/op
BenchmarkReadWriteMemory-8         13507             85256 ns/op
PASS
```
