package shardMap

import (
	"fmt"
	"sync"
)

const SHARD_COUNT = 32

type ShardMap[K comparable, V any] struct {
	shards []*ShardMapShared[K, V]
}
