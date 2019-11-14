package stress

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync/atomic"
	"testing"
	"unsafe"

	sim "github.com/dgraph-io/ristretto/sim"
	"github.com/karlmcguire/stress/chanDropSharded"
	"github.com/karlmcguire/stress/lock"
	"github.com/karlmcguire/stress/lockSharded"
	"github.com/karlmcguire/stress/sync"
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

func genKeysZipf() [numKeys]uint64 {
	var keys [numKeys]uint64
	zipf := sim.NewZipfian(1.5, 2, numKeys)
	for i := range keys {
		keys[i], _ = zipf()
	}
	return keys
}

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
		{"sync", sync.New(numKeys)},
		{"lock", lock.New(numKeys)},
		{"lockSharded", lockSharded.New(numKeys)},
		//{"chanDrop", chanDrop.New(numKeys)},
		{"chanDropSharded", chanDropSharded.New(numKeys)},
		//{"syncChanDropSharded", syncChanDropSharded.New(numKeys)},
		//{"ring", ring.New(numKeys)},
		//{"ringChanDropSharded", ringChanDropSharded.New(numKeys)},
	}
}

//go:linkname FastRand runtime.fastrand
func FastRand() uint32

func BenchmarkTesting(b *testing.B) {
	rc := uint64(0)
	_ = unsafe.Pointer(&rc)
	sets, gets := uint64(0), uint64(0)
	b.RunParallel(func(pb *testing.PB) {
		mc := atomic.AddUint64(&rc, 1)
		for pb.Next() {
			if 25*mc/100 != 25*(mc-1)/100 {
				atomic.AddUint64(&sets, 1)
			} else {
				atomic.AddUint64(&gets, 1)
			}
		}
	})
	fmt.Printf("%0.2f%% writes\n", float64(sets)/(float64(sets)+float64(gets)))
}

func BenchmarkMixed(b *testing.B) {
	keys, benchmarks := genKeysZipf(), genBenchmarks()
	for _, benchmark := range benchmarks {
		rc := uint64(0)
		b.Run(benchmark.Name, func(b *testing.B) {
			b.SetBytes(1)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				mc := atomic.AddUint64(&rc, 1)
				if 50*mc/100 != 50*(mc-1)/100 {
					for i := rand.Int(); pb.Next(); i++ {
						benchmark.Map.Set(keys[i&keyMask], 0)
					}
				} else {
					for i := rand.Int(); pb.Next(); i++ {
						benchmark.Map.Get(keys[i&keyMask])
					}
				}
			})
		})
	}
}

func BenchmarkGet(b *testing.B) {
	keys, benchmarks := genKeysZipf(), genBenchmarks()
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
	keys, benchmarks := genKeysZipf(), genBenchmarks()
	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			b.SetBytes(1)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for i := rand.Int(); pb.Next(); i++ {
					benchmark.Map.Set(keys[i&keyMask], 0)
				}
			})
		})
		runtime.GC()
	}
}

/*
func BenchmarkSetAll(b *testing.B) {
	pairs, benchmarks := genPairs(), genBenchmarks()
	for _, benchmark := range benchmarks {
		b.Run(benchmark.Name, func(b *testing.B) {
			b.SetBytes(batchSize)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					benchmark.Map.SetAll(0)
				}
			})
		})
		runtime.GC()
	}
}
*/
