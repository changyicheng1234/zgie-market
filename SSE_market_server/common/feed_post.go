package common

import (
	"fmt"
	"loginTest/model"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/jinzhu/gorm"
)

const (
	HomeFeedPartition  = "主页"
	NewFrontendBaseURL = "https://ssemarket.cn/new"
	MaxFeedTitleRunes  = 30
)

// FeedTitleWithPrefix 生成不超过 maxRunes 字的「前缀+名称」标题。
func FeedTitleWithPrefix(prefix, name string) string {
	full := prefix + name
	if utf8.RuneCountInString(full) <= MaxFeedTitleRunes {
		return full
	}
	allowed := MaxFeedTitleRunes - utf8.RuneCountInString(prefix)
	runes := []rune(name)
	if len(runes) > allowed {
		runes = runes[:allowed]
	}
	return prefix + string(runes)
}

// CreateHomeFeedPost 在主页发布同步帖。
func CreateHomeFeedPost(db *gorm.DB, user model.User, title, content string) error {
	if user.Banend.After(time.Now()) {
		return fmt.Errorf("你尚处于禁言状态中，不得同步发帖")
	}
	if utf8.RuneCountInString(content) > 10000 {
		return fmt.Errorf("同步发帖内容过长")
	}
	newPost := model.Post{
		UserID:     user.UserID,
		Partition:  HomeFeedPartition,
		Title:      title,
		Ptext:      content,
		LikeNum:    0,
		CommentNum: 0,
		BrowseNum:  0,
		Heat:       0,
		PostTime:   time.Now(),
		Photos:     "",
		Tag:        "",
	}
	if err := db.Create(&newPost).Error; err != nil {
		return fmt.Errorf("同步发帖失败")
	}
	BloomFilter.Add([]byte(strconv.FormatUint(uint64(newPost.PostID), 10)))
	return nil
}
