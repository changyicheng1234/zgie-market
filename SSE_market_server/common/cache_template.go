package common

import (
	"context"
	"encoding/json"
	"log"
	"loginTest/core"
	"loginTest/util"
	"time"
)

var ctx = context.Background()

// MultiLevelGet 多级缓存统一读取模板（全项目读接口调用）
// 流程：本地缓存 → Redis → 数据库 → 回写两级缓存
func MultiLevelGet[T any](key string, dbFunc func() (T, error), expire time.Duration) (T, error) {
	var zero T

	// 1. 查 L1 本地缓存（最快）
	data, err := core.LocalCache.Get(key)
	if err == nil {
		var res T
		if errUnmarshal := json.Unmarshal(data, &res); errUnmarshal == nil {
			return res, nil
		}
	}

	// 2. 查 L2 Redis 缓存
	redisData, err := core.MyRedis.Get(ctx, key).Result()
	if err == nil {
		// 回写到本地缓存
		_ = core.LocalCache.Set(key, []byte(redisData))
		var res T
		if errUnmarshal := json.Unmarshal([]byte(redisData), &res); errUnmarshal == nil {
			return res, nil
		}
	}

	// 3. 缓存未命中：加分布式锁，防止缓存击穿
	lock := util.NewLock("lock:" + key)
	if !lock.TryLock() {
		return zero, nil
	}
	defer lock.UnLock()

	// 双重检查：防止并发重复查库
	redisData, _ = core.MyRedis.Get(ctx, key).Result()
	if redisData != "" {
		var res T
		if errUnmarshal := json.Unmarshal([]byte(redisData), &res); errUnmarshal == nil {
			return res, nil
		}
	}

	// 4. 查询数据库
	dbData, err := dbFunc()
	if err != nil {
		return zero, err
	}

	// 5. 回写两级缓存（带随机过期，防雪崩）
	jsonData, _ := json.Marshal(dbData)
	_ = core.MyRedis.Set(ctx, key, jsonData, util.GetRandomExpire(expire)).Err()
	_ = core.LocalCache.Set(key, jsonData)

	return dbData, nil
}

// DelMultiLevelCache 删除两级缓存（写接口调用）
func DelMultiLevelCache(keys ...string) {
	for _, key := range keys {
		core.LocalCache.Delete(key)
		core.MyRedis.Del(ctx, key)
		log.Printf("[DelMultiLevelCache] 删除缓存 | key: %s", key)
	}
}

// InvalidatePostDetailCache 帖子详情变更时失效两级缓存（评论数、点赞数等）
func InvalidatePostDetailCache(postID int) {
	if postID <= 0 {
		return
	}
	DelMultiLevelCache(util.GetPostDetailKey(postID))
}
