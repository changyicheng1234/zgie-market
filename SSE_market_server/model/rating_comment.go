package model

import "time"

type RatingComment struct {
	CommentID    int       `gorm:"primary_key;column:comment_id"`
	RatingPostID int       `gorm:"index:rating_comment_post;column:rating_post_id;type:int;not null"`
	UserID       int       `gorm:"index:rating_comment_user;column:user_id;type:int;not null"`
	Content      string    `gorm:"column:content;type:varchar(1000)"`
	StarNum      int       `gorm:"column:star_num;type:int;default:0"`
	LikeNum      int       `gorm:"column:like_num;type:int;default:0"`
	DenyNum      int       `gorm:"column:deny_num;type:int;default:0"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime"`
}
