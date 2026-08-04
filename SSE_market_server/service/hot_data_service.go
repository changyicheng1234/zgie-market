package service

import (
	"context"
	"errors"
	"log"
	"loginTest/common"
	"loginTest/core"
	"loginTest/model"
	"loginTest/util"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// 全局单例
var HotData *HotDataService
var ctx = context.Background()

// HotDataService 浏览数计数服务：Redis INCR 为跨实例唯一真相，周期同步到 MySQL。
type HotDataService struct {
	dbSyncInterval  time.Duration
	browseDirty     sync.Map // postID -> struct{}，该帖 Redis 计数有变更、待落库
	redisHealthy    bool
	stopChan        chan struct{}
	lastSyncedPosts map[int]int64 // 上次已成功写入 DB 的浏览数，避免重复 UPDATE
	lastSyncedMutex sync.Mutex
}

// InitHotDataService 初始化
func InitHotDataService() {
	dbSyncInterval := viper.GetDuration("hot_data.db_sync_interval")
	if dbSyncInterval == 0 {
		dbSyncInterval = 5 * time.Second
	}

	HotData = &HotDataService{
		dbSyncInterval:  dbSyncInterval,
		redisHealthy:    true,
		stopChan:        make(chan struct{}),
		lastSyncedPosts: make(map[int]int64),
	}

	go HotData.startHealthCheck()
	go HotData.startDBSyncTask()
}

func (s *HotDataService) browseNumFromDB(postID int) int64 {
	db := common.GetDB()
	if db == nil {
		return 0
	}
	var p model.Post
	if err := db.Select("browse_num").Where("postID = ?", postID).First(&p).Error; err != nil {
		return 0
	}
	return int64(p.BrowseNum)
}

// ensureBrowseCountKeySeeded 在 Redis 中不存在计数 key 时用 DB 的 browse_num 做 SetNX，避免 Incr 从 0 起算导致与 MySQL 分叉。
func (s *HotDataService) ensureBrowseCountKeySeeded(postID int) {
	if !s.redisHealthy {
		return
	}
	key := util.GetPostBrowseCountKey(postID)
	_, err := core.MyRedis.Get(ctx, key).Result()
	if err == nil {
		return
	}
	if !errors.Is(err, redis.Nil) {
		return
	}
	n := s.browseNumFromDB(postID)
	ok, errNX := core.MyRedis.SetNX(ctx, key, n, 0).Result()
	if errNX == nil && ok {
		log.Printf("[hot_data] browse count key missing, seeded from DB postID=%d val=%d", postID, n)
	}
}

// IncrBrowseCount 浏览数 +1：直接 Redis INCR，多副本下与 MySQL 最终由 syncToDB 对齐。
func (s *HotDataService) IncrBrowseCount(postID int) int64 {
	if !s.redisHealthy {
		// Redis 不可用时不上报虚假增长，避免与恢复后的 Redis 剧烈跳变
		return s.browseNumFromDB(postID)
	}
	s.ensureBrowseCountKeySeeded(postID)
	key := util.GetPostBrowseCountKey(postID)
	n, err := core.MyRedis.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("[hot_data] IncrBrowseCount redis Incr failed postID=%d: %v", postID, err)
		return s.browseNumFromDB(postID)
	}
	s.browseDirty.Store(postID, struct{}{})
	return n
}

// SyncBrowseRedisFloor 在仅更新 MySQL 后把 Redis 计数至少抬到 minVal，避免与 ShowDetails 分叉。
func (s *HotDataService) SyncBrowseRedisFloor(postID int, minVal int64) {
	if minVal <= 0 || !s.redisHealthy {
		return
	}
	key := util.GetPostBrowseCountKey(postID)
	cur, err := core.MyRedis.Get(ctx, key).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return
	}
	if errors.Is(err, redis.Nil) || minVal > cur {
		if err := core.MyRedis.Set(ctx, key, minVal, 0).Err(); err == nil {
			s.browseDirty.Store(postID, struct{}{})
		}
	}
}

// Redis健康检查
func (s *HotDataService) startHealthCheck() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := core.MyRedis.Ping(ctx).Result()
			if err != nil {
				if s.redisHealthy {
					s.redisHealthy = false
				}
			} else {
				if !s.redisHealthy {
					s.redisHealthy = true
				}
			}
		case <-s.stopChan:
			return
		}
	}
}

// startDBSyncTask 启动数据库同步任务
func (s *HotDataService) startDBSyncTask() {
	ticker := time.NewTicker(s.dbSyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.syncToDB()
		case <-s.stopChan:
			s.syncToDB()
			return
		}
	}
}

// syncToDB 将 Redis 中的浏览数同步到数据库（处理所有 browseDirty 中的帖子，与本地 pending 无关）
func (s *HotDataService) syncToDB() {
	if !s.redisHealthy {
		return
	}

	db := common.GetDB()
	if db == nil {
		return
	}

	postIDs := make([]int, 0, 64)
	s.browseDirty.Range(func(k, v interface{}) bool {
		postIDs = append(postIDs, k.(int))
		return true
	})
	if len(postIDs) == 0 {
		return
	}

	s.lastSyncedMutex.Lock()
	defer s.lastSyncedMutex.Unlock()

	for _, postID := range postIDs {
		key := util.GetPostBrowseCountKey(postID)
		s.ensureBrowseCountKeySeeded(postID)
		redisBrowseNum, err := core.MyRedis.Get(ctx, key).Int64()
		if err != nil {
			continue
		}

		lastSynced, exists := s.lastSyncedPosts[postID]
		if exists && lastSynced == redisBrowseNum {
			s.browseDirty.Delete(postID)
			continue
		}

		err = db.Model(&model.Post{}).
			Where("postID = ?", postID).
			Update("browse_num", redisBrowseNum).Error
		if err != nil {
			continue
		}
		s.lastSyncedPosts[postID] = redisBrowseNum
		s.browseDirty.Delete(postID)
	}
}

// Shutdown 关闭
func (s *HotDataService) Shutdown() {
	close(s.stopChan)
}
