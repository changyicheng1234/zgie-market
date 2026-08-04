package heat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"loginTest/config"
	"loginTest/controller"
	"math"
	"time"

	"github.com/redis/go-redis/v9"

	"loginTest/common"
	"loginTest/core"
	"loginTest/model"

	"github.com/robfig/cron/v3"
)

const (
	// 课程交流分区名称
	CoursePartition = "课程交流"
)

// 启动热度计算定时任务
func StartHeat() {
	RefreshHeat()
	// 创建 cron 实例并添加任务
	c := cron.New()

	// 每h刷新
	c.AddFunc("@hourly", RefreshHeat)
	// 每 48 小时清理一次已使用的 refresh token
	c.AddFunc("@every 48h", CleanUsedRefreshTokens)
	c.Start()

	// 防止程序退出
	select {}
}

// 刷新热度并存入 Redis ZSet
func RefreshHeat() {
	ctx := context.Background()
	db := common.GetDB()

	// 获取当前时间并计算10天前的时间
	now := time.Now()
	tenDaysAgo := now.AddDate(0, 0, -config.HEAT_DAY)

	// 获取10天内的所有帖子，排除课程交流分区
	var posts []model.Post
	db.Where("post_time >= ? AND `partition` != ? AND is_private = ?", tenDaysAgo, CoursePartition, false).Find(&posts)

	// Redis ZSet的key，用来存储帖子热度
	redisKey := "hot_posts_zset"
	// 清空之前的 ZSet，确保只有当天最新的热度数据
	core.MyRedis.Del(ctx, redisKey)

	// 计算每个帖子的热度并存入 Redis ZSet
	for _, post := range posts {
		heatScore := calculateHeat(post.BrowseNum, post.CommentNum, post.LikeNum, math.Floor(now.Sub(post.PostTime).Hours()/24))

		// 创建包含 PostID 和 Title 的 JSON 字符串
		postInfo := controller.PostInfo{
			PostID:    uint(post.PostID),
			Title:     post.Title,
			BrowseNum: float64(post.BrowseNum),
		}
		postInfoJSON, err := json.Marshal(postInfo)
		if err != nil {
			log.Println("Error marshalling postInfo:", err)
			continue
		}

		// 将热度和帖子信息存入 Redis ZSet
		core.MyRedis.ZAdd(ctx, redisKey, redis.Z{
			Score:  heatScore,
			Member: postInfoJSON, // 存储 JSON 字符串
		})
	}

	// 获取Redis中热度前10的帖子
	topPosts, err := core.MyRedis.ZRevRangeWithScores(ctx, redisKey, 0, 9).Result()
	if err != nil {
		log.Println("获取Redis热度前10帖子失败:", err)
	} else {
		fmt.Println("Redis热度前10帖子已刷新:", topPosts)
	}
}

// CleanUsedRefreshTokens 定期清理已使用的 refresh token 记录
func CleanUsedRefreshTokens() {
	db := common.GetDB()

	// 清理3天前的已使用token
	retentionDays := 3
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	result := db.Where("is_used = ? AND created_at < ?", true, cutoffTime).
		Delete(&model.RefreshToken{})

	if result.Error != nil {
		log.Printf("清理已使用 refresh token 失败: %v", result.Error)
		return
	}

	log.Printf("定期清理已使用 refresh token 完成，清理了 %d 条记录", result.RowsAffected)
}

// 计算帖子热度值的函数
func calculateHeat(browseNum, commentNum, likeNum int, daysSincePost float64) float64 {
	return math.Log10(float64(browseNum)+float64(commentNum)/10+float64(likeNum)/2) - daysSincePost
}
