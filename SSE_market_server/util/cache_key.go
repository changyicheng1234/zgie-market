package util

import (
	"fmt"
	"math/rand"
	"time"
)

// 缓存全局前缀
const (
	// 原有基础缓存前缀
	PrefixPostDetail      = "post:detail:"       // 帖子详情（ShowDetails）
	PrefixHotPost         = "post:hot:list"      // 热点帖子列表
	PrefixUserInfo        = "user:info:"         // 用户信息
	PrefixCommentList     = "comment:list:"      // 帖子评论列表（主评论+子评论）
	PrefixPostBrowse      = "post:browse:"       // 帖子浏览列表（home/history/save/rating）
	PrefixProductList     = "product:list:"      // 商品列表（home/个人商品）
	PrefixPostBrowseCount = "post:browse:count:" // 帖子浏览数独立Key
)

// GetPostBrowseCountKey 生成帖子浏览数缓存Key（适配ShowDetails函数，postID为int）
func GetPostBrowseCountKey(postID int) string {
	return PrefixPostBrowseCount + fmt.Sprintf("%d", postID)
}

// GetPostDetailKey 生成帖子详情缓存Key（适配ShowDetails函数，postID为int）
func GetPostDetailKey(postID int) string {
	return PrefixPostDetail + fmt.Sprintf("%d", postID)
}

// GetHotPostKey 生成热点帖子列表缓存Key
func GetHotPostKey() string {
	return PrefixHotPost
}

// GetUserInfoKey 生成用户信息缓存Key（userID为int）
func GetUserInfoKey(userID int) string {
	return PrefixUserInfo + fmt.Sprintf("%d", userID)
}

// GetCommentListKey 生成帖子评论列表缓存Key
// 参数：postID 帖子ID，postType 帖子类型（post/评分贴）- 两者共同决定评论列表唯一性
func GetCommentListKey(postID int, postType string) string {
	return PrefixCommentList + fmt.Sprintf("%d:%s", postID, postType)
}

// GetPostBrowseKey 生成帖子浏览列表缓存Key
func GetPostBrowseKey(userID int, searchsort, partition, searchinfo, tag string) string {
	// 对空值做默认处理，避免Key中出现空串导致的隐性问题
	if partition == "" {
		partition = "default"
	}
	if searchinfo == "" {
		searchinfo = "empty"
	}
	if tag == "" {
		tag = "empty"
	}
	return PrefixPostBrowse + fmt.Sprintf("%d:%s:%s:%s:%s", userID, searchsort, partition, searchinfo, tag)
}

// GetProductListKey 生成商品列表缓存Key
// 参数：userID 登录用户ID，searchsort 筛选类型(home/个人商品)，searchinfo 搜索关键词
// limit/offset不加入Key：分页参数单独处理，避免缓存Key爆炸
func GetProductListKey(userID int, searchsort, searchinfo string) string {
	if searchinfo == "" {
		searchinfo = "empty"
	}
	return PrefixProductList + fmt.Sprintf("%d:%s:%s", userID, searchsort, searchinfo)
}

// GetRandomExpire 随机过期时间（防缓存雪崩）
func GetRandomExpire(base time.Duration) time.Duration {
	rand.Seed(time.Now().UnixNano())
	// 随机偏移1~300秒，避免大量Key同时过期
	offset := time.Duration(rand.Intn(300)) * time.Second
	return base + offset
}
