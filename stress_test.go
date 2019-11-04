package stress

import (
	"math/rand"
	"testing"

	"github.com/karlmcguire/stress/basic"
	"github.com/karlmcguire/stress/syncMap"
)

const (
	numKeys       = 1e7
	keyMask       = numKeys - 1
	numGoroutines = 8
	segSize       = numKeys / numGoroutines
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

func genBenchmarks() []*Benchmark {
	return []*Benchmark{
		{"basic", basic.New(numKeys)},
		{"sync.Map", syncMap.New(numKeys)},
	}
}

func BenchmarkMaps(b *testing.B) {
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
	}
}
