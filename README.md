# stress

Stress testing hashmap implemenations.


## Results
```
goos: linux
goarch: amd64

BenchmarkSet/lock-12                    20000000        86.7 ns/op    11.53 MB/s
BenchmarkSet/lockSharded-12            100000000        13.5 ns/op    74.22 MB/s
BenchmarkSet/chanDrop-12               500000000        3.26 ns/op   306.48 MB/s
BenchmarkSet/chanDropSharded-12        500000000        2.52 ns/op   396.79 MB/s
BenchmarkSet/syncMap-12                 10000000         170 ns/op     5.85 MB/s
BenchmarkSet/ring-12                    30000000        34.1 ns/op    29.29 MB/s

BenchmarkSetAll/lock-12                   100000       11802 ns/op    84.73 MB/s
BenchmarkSetAll/lockSharded-12            100000       15131 ns/op    66.09 MB/s
BenchmarkSetAll/chanDrop-12               200000       11799 ns/op    84.75 MB/s
BenchmarkSetAll/chanDropSharded-12        100000       14825 ns/op    67.45 MB/s
BenchmarkSetAll/syncMap-12                 10000      154386 ns/op     6.48 MB/s
BenchmarkSetAll/ring-12                   100000       11670 ns/op    85.69 MB/s

BenchmarkGet/lock-12                    30000000        41.0 ns/op    24.37 MB/s
BenchmarkGet/lockSharded-12            200000000        6.74 ns/op   148.31 MB/s
BenchmarkGet/chanDrop-12                10000000         119 ns/op     8.40 MB/s
BenchmarkGet/chanDropSharded-12        200000000        6.80 ns/op   147.17 MB/s
BenchmarkGet/syncMap-12                200000000        6.08 ns/op   164.45 MB/s
BenchmarkGet/ring-12                    30000000        43.5 ns/op    22.99 MB/s
```

## Ranked

### 1. chanDropSharded

```
BenchmarkSet/chanDropSharded-12        500000000        2.52 ns/op   396.79 MB/s
BenchmarkSetAll/chanDropSharded-12        100000       14825 ns/op    67.45 MB/s
BenchmarkGet/chanDropSharded-12        200000000        6.80 ns/op   147.17 MB/s
```

### 2. chanDrop

```
BenchmarkSet/chanDrop-12               500000000        3.26 ns/op   306.48 MB/s
BenchmarkSetAll/chanDrop-12               200000       11799 ns/op    84.75 MB/s
BenchmarkGet/chanDrop-12                10000000         119 ns/op     8.40 MB/s
```

---

**chanDropSharded** outperforms **chanDrop** on every benchmark except `SetAll`,
this is because of the overhead required to select a shard for every `Set` call:

```go
func (m *Map) SetAll(pairs [][2]uint64) {
	for _, pair := range pairs {
		m.shards[pair[0]&shardMask].Set(pair[0], pair[1])
	}
}
```

However, this overhead is acceptable because **chanDropSharded** has over 10x
better performance for `Get` operations, due to the lowered contention.
