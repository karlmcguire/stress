package syncPoolChanDropSharded

const (
	numShards = 256
	shardMask = numShards - 1
)

type Map struct {
	shards [numShards]lockedMap
	items  chan [][3]uint64
}
