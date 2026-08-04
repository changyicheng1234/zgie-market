package core

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
)

// 全局单例：L1本地缓存（全项目共用）
var LocalCache *bigcache.BigCache

// InitLocalCache 初始化本地缓存
func InitLocalCache() {

	// 生产级配置
	config := bigcache.Config{
		Shards: 2,
		// 缓存过期时间：5分钟（热点数据）
		LifeWindow: 500 * time.Minute,
		// 后台清理过期数据间隔
		CleanWindow: 1 * time.Minute,
		// 最大缓存数量，避免OOM
		MaxEntriesInWindow: 1000 * 1000,
		// 关闭日志
		Verbose: false,
	}

	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		panic(fmt.Sprintf("本地缓存初始化失败: %v", err))
	}

	LocalCache = cache
	fmt.Println("✅ L1本地缓存(BigCache)初始化成功")
}
