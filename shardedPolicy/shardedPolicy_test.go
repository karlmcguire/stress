package shardedPolicy

import "testing"

func BenchmarkShardedPolicy(b *testing.B) {
	p := New(1e6)

}
