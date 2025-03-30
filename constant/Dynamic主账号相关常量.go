package constant

import (
	"time"
)

const (
	VsIPTransitDynamicCacheRedisKeyPrefix = "VsIPTransitDynamicPrefix_" // 主账号redis缓存的key前缀

	FlowUserIdQueueSortedSet = "FlowUserIdQueueSortedSet" ///主账号流量redis有序集合队列
)

const (
	BatchUpdateDynamicDataCacheSize = 100
)

const (
	VsIPTransitDynamicCacheRedisTtl           = 60 * 60 * time.Second ///主账号redis缓存的过期时间
	VsIPTransitDynamicBatchUpdateDataLoopTime = 15 * 60 * time.Second
)
