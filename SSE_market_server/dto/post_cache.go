package dto

import (
	"loginTest/model"
	"time"
)

// PostDetailCacheDTO 帖子详情缓存专用 DTO
type PostDetailCacheDTO struct {
	Post PostCacheItem `json:"post"`
	User UserCacheItem `json:"user"`
}

// PostCacheItem 仅缓存帖子【静态冷数据】
// 时间使用 string 存储，彻底避免 JSON 解析崩溃
type PostCacheItem struct {
	PostID     uint   `json:"post_id"`
	UserID     uint   `json:"user_id"`
	Title      string `json:"title"`
	Ptext      string `json:"ptext"`
	LikeNum    int    `json:"like_num"`
	CommentNum int    `json:"comment_num"`
	BrowseNum  int    `json:"browse_num"`
	Photos     string `json:"photos"`
	Tag        string `json:"tag"`
	IsPrivate  bool   `json:"is_private"`
	PostTime   string `json:"post_time"`
}

// UserCacheItem 用户信息缓存结构
type UserCacheItem struct {
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Score     int    `json:"score"`
	AvatarURL string `json:"avatar_url"`
	Identity  string `json:"identity"`
}

// ---------------------------------------------------------------------------------
// ConvertFromModel
// 将 model.Post + model.User 转换成 DTO（写入缓存前调用）
// 类型 100% 安全，time.Time → RFC3339 字符串
// ---------------------------------------------------------------------------------
func (d *PostDetailCacheDTO) ConvertFromModel(post model.Post, user model.User) {
	// 转换帖子信息
	d.Post = PostCacheItem{
		PostID:     uint(post.PostID),
		UserID:     uint(post.UserID),
		Title:      post.Title,
		Ptext:      post.Ptext,
		LikeNum:    post.LikeNum,
		CommentNum: post.CommentNum,
		BrowseNum:  post.BrowseNum,
		Photos:     post.Photos,
		Tag:        post.Tag,
		IsPrivate:  post.IsPrivate,
		PostTime:   post.PostTime.Format(time.RFC3339), // 时间格式化，无崩溃
	}

	// 转换用户信息
	d.User = UserCacheItem{
		UserID:    uint(user.UserID),
		Name:      user.Name,
		Phone:     user.Phone,
		Score:     user.Score,
		AvatarURL: user.AvatarURL,
		Identity:  user.Identity,
	}
}

// ---------------------------------------------------------------------------------
// ConvertToModel
// 将 DTO 转回 model.Post + model.User（读取缓存后调用）
// 字符串时间 → time.Time，类型完全匹配
// 不会修改 heat 字段，heat 由外部其他函数维护
// ---------------------------------------------------------------------------------
func (d *PostDetailCacheDTO) ConvertToModel() (model.Post, model.User) {
	// 时间字符串转回 time.Time（一定成功，因为是我们自己存的 RFC3339）
	postTime, _ := time.Parse(time.RFC3339, d.Post.PostTime)

	// 组装 Post model
	post := model.Post{
		PostID:     int(d.Post.PostID),
		UserID:     int(d.Post.UserID),
		Title:      d.Post.Title,
		Ptext:      d.Post.Ptext,
		LikeNum:    d.Post.LikeNum,
		CommentNum: d.Post.CommentNum,
		BrowseNum:  d.Post.BrowseNum,
		Photos:     d.Post.Photos,
		Tag:        d.Post.Tag,
		IsPrivate:  d.Post.IsPrivate,
		PostTime:   postTime,
	}

	// 组装 User model
	user := model.User{
		UserID:    int(d.User.UserID),
		Name:      d.User.Name,
		Phone:     d.User.Phone,
		Score:     d.User.Score,
		AvatarURL: d.User.AvatarURL,
		Identity:  d.User.Identity,
	}

	return post, user
}
