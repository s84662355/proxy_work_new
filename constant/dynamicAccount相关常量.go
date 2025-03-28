package constant

import (
	"time"
)

const (
	DynamicAccountDataCacheRedisKeyPrefix     = "ProxyDynamicAccountFlowHashPrefix_" ///缓存DynamicAccount数据的key前缀
	DynamicAccountDataCacheByIdRedisKeyPrefix = DynamicAccountDataCacheRedisKeyPrefix + "id_"

	DynamicAccountRedisFlowPrefix = "ProxyIncrementFlow_" // DynamicAccount redis 流量前缀

	DynamicAccountIDFlowRedisQueueSet = "FlowAccountIDQueueSet" ///
)

const (
	BatcheUpdateDynamicAccountDataCacheSize = 100
)

const (
	DynamicAccountDataCacheRedisTtl       = 600 * time.Second
	DynamicAccountBatchUpdateDataLoopTime = 300 * time.Second
)
