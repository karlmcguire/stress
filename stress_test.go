package stress

import (
	"math/rand"
	"runtime"
	"testing"

	"github.com/karlmcguire/stress/syncChanDropSharded"
)

const (
	numKeys       = 1e7
	keyMask       = numKeys - 1
	numGoroutines = 8
	segSize       = numKeys / numGoroutines
	batchSize     = 1e3
)

type (
	Map interface {
		Get(uint64) uint64
		Set(uint64, uint64)
		SetAll([][2]uint64)
		Del(uint64)
	}
	Benchmark struct {
		Name string
		Map  Map
	}
)

func genKeys() [numKeys]uint64 {
	var keys [numKeys]uint64
	for i := range keys {
		keys[i] = rand.Uint64() % numKeys
	}
	return keys
}

func genPairs() [][2]uint64 {
	keys := genKeys()
	pairs := make([][2]uint64, batchSize)
	for i := range pairs {
		pairs[i] = [2]uint64{keys[i], keys[i]}
	}
	return pairs
}

func genBenchmarks() []*Benchmark {
	return []*Benchmark{
		/*
			{"lock", lock.New(numKeys)},
			{"lockSharded", lockSharded.New(numKeys)},
			{"chanDrop", chanDrop.New(numKeys)},
			{"chanDropSharded", chanDropSharded.New(numKeys)},
			{"sync", sync.New(numKeys)},
		*/
		{"syncChanDropSharded", syncChanDropSharded.New(numKeys)},
		//{"ring", ring.New(numKeys)},
	}
}

func BenchmarkGet(b *testing.B) {
	keys, benchmarks := genKeys(), genBenchmarks()
	for _, benchmark := range benchmarks {
		for _, key := range keys {
			benchmark.Map.Set(key, key)
		}
		b.Run(benchmark.Name, func(b *testing.B) {
			b.SetBytes(1)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for i := rand.Int() & keyMask; pb.Next(); i++ {
					benchmark.Map.Get(keys[i&keyMask])
				}
			})
		})
	}
}

func BenchmarkSet(b *testing.B) {
	keys, benchmarks := genKeys(), genBenchmarks()
	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			b.SetBytes(1)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for i := rand.Int() & keyMask; pb.Next(); i++ {
					benchmark.Map.Set(keys[i&keyMask], keys[i&keyMask])
				}
			})
		})
		runtime.GC()
	}
}

func BenchmarkSetAll(b *testing.B) {
	pairs, benchmarks := genPairs(), genBenchmarks()
	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			b.SetBytes(batchSize)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					benchmark.Map.SetAll(pairs)
				}
			})
		})
		runtime.GC()
	}
}
