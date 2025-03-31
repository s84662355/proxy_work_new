package constant

import (
	"time"
)

const (
	VsIPTransitDynamicCacheRedisKeyPrefix     = "VsIPTransitDynamicPrefix_" // 主账号redis缓存的key前缀
	FlowUserIdQueueSortedSet                  = "FlowUserIdQueueSortedSet"  ///主账号流量redis有序集合队列
	VsIPTransitDynamicUpdateFlowRedisLock     = "updatedbflow_lock"         // 主账号更新流量的分布式锁前缀
	VsIPTransitDynamicUpdateFlowRedisLockTtl  = 60 * time.Second            // 主账号更新流量的分布式锁的超时时间
	VsIPTransitDynamicAutoBuyFlowRedisLock    = "BuyFlow__lock"             // 主账号自动购买流量的分布式锁前缀
	VsIPTransitDynamicAutoBuyFlowRedisLockTtl = 60 * time.Second            ///主账号自动购买流量的分布式锁过期时间
)

const (
	BatchUpdateDynamicDataCacheSize = 100
)

const (
	VsIPTransitDynamicCacheRedisTtl           = 60 * 60 * time.Second ///主账号redis缓存的过期时间
	VsIPTransitDynamicBatchUpdateDataLoopTime = 15 * 60 * time.Second
)
