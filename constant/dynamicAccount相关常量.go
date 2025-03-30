package constant

import (
	"time"
)

const (
	DynamicAccountDataCacheRedisKeyPrefix     = "ProxyDynamicAccountFlowHashPrefix_" ///缓存DynamicAccount数据的key前缀
	DynamicAccountDataCacheByIdRedisKeyPrefix = DynamicAccountDataCacheRedisKeyPrefix + "id_"
	DynamicAccountRedisFlowPrefix             = "ProxyIncrementFlow_"         // DynamicAccount redis 流量前缀
	DynamicAccountIDFlowRedisQueueSet         = "FlowAccountIDQueueSet"       ///
	DynamicAccountIDFlowRedisQueueSortedSet   = "FlowAccountIDQueueSortedSet" ///子账号流量队列 redis的有序集合
	DynamicAccountDbLock                      = "DynamicAccountDbLock_"       // 更新子账号流量的分布式锁前缀
)

const (
	BatcheUpdateDynamicAccountDataCacheSize = 100
)

const (
	DynamicAccountDataCacheRedisTtl       = 60 * 60 * time.Second
	DynamicAccountBatchUpdateDataLoopTime = 15 * 60 * time.Second
	DynamicAccountRedisFlowTtl            = 60 * 60 * time.Second
)

const ExistsFlowDynamicAccountIDbyRedisLua = `
			-- 批量获取传入的键的值
			local keys = redis.call('MGET', unpack(KEYS))
			if not keys then
			    -- 如果 MGET 执行失败，返回错误信息
			    return redis.error_reply('MGET command failed')
			end

			-- 初始化一个空数组，用于存储值大于 0 的键
			local resultKeys = {}
			-- 初始化一个索引，用于记录 resultKeys 数组的位置
			local index = 1

			-- 遍历获取到的键的值
			for i, value in ipairs(keys) do
			    -- 检查值是否存在且能成功转换为数字，并且转换后的数字大于 0
			    local num = tonumber(value)
			    if num and num > 0 then
			        -- 将符合条件的键添加到 resultKeys 数组
			        resultKeys[index] = KEYS[i]
			        -- 索引加 1
			        index = index + 1
			    end
			end

			-- 从 ARGV 中获取集合的键名
			local setKey = ARGV[1]
			if not setKey then
			    -- 如果没有传入集合键名，返回错误信息
			    return redis.error_reply('Set key is not provided')
			end

			-- 如果 resultKeys 数组不为空，则批量添加到集合中
			if #resultKeys > 0 then
			    -- 调用 SADD 命令将 resultKeys 中的键添加到集合中
			    local addedCount = redis.call('SADD', setKey, unpack(resultKeys))
			    if not addedCount then
			        -- 如果 SADD 执行失败，返回错误信息
			        return redis.error_reply('SADD command failed')
			    end
			    -- 返回成功添加到集合中的新元素数量
			    return addedCount
			end

			-- 如果 resultKeys 数组为空，返回 0
			return 0
    `
