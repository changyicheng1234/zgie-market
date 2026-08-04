package common

import (
	"log"
	"loginTest/model"
	"strconv"

	"github.com/bits-and-blooms/bloom/v3"
)

// 全局布隆过滤器（永远不为nil）
var BloomFilter *bloom.BloomFilter

// InitBloomFilter 企业级初始化：失败不panic，服务降级运行
func InitBloomFilter() {
	// 1. 兜底初始化（杜绝空指针）
	BloomFilter = bloom.NewWithEstimates(100000, 0.0001)
	log.Println("布隆过滤器基础实例初始化完成")

	// 2. 依赖检查：无DB直接跳过
	if DB == nil {
		log.Println("DB未连接，布隆过滤器降级运行")
		return
	}

	// 3. 安全加载数据：【修复字段名】id → post_id
	var postIDs []string
	err := DB.Model(&model.Post{}).Pluck("postID", &postIDs).Error
	if err != nil {
		log.Printf("⚠️ 布隆过滤器加载帖子ID失败：%v", err)
		log.Println("⚠️ 布隆过滤器已降级，不影响服务运行")
		return
	}

	// 4. 数据加载成功
	for _, pid := range postIDs {
		BloomFilter.Add([]byte(pid))
	}
	log.Printf("✅ 布隆过滤器初始化完成，加载 %d 条帖子数据", len(postIDs))
}

// CheckPostExist 安全校验：永远不panic
func CheckPostExist(postID int) bool {
	// 降级规则：过滤器未就绪 / ID非法 → 直接放行
	if BloomFilter == nil || postID <= 0 {
		return true
	}
	return BloomFilter.Test([]byte(strconv.Itoa(postID)))
}

func AddPostID(postID int) {
	if BloomFilter == nil || postID <= 0 {
		return
	}
	BloomFilter.Add([]byte(strconv.Itoa(postID)))
}
