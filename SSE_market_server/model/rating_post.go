package model

import "time"

type RatingPost struct {
	RatingPostID int       `gorm:"primary_key;column:rating_post_id"`
	CategoryID   int       `gorm:"index:rating_post_category;column:category_id;type:int;not null"`
	UserID       int       `gorm:"index:rating_post_user;column:user_id;type:int;not null"`
	Title        string    `gorm:"column:title;type:varchar(100)"`
	Content      string    `gorm:"column:content;type:varchar(5000)"`
	CommentNum   int       `gorm:"column:comment_num;type:int;default:0"`
	LikeNum      int       `gorm:"column:like_num;type:int;default:0"`
	BrowseNum    int       `gorm:"column:browse_num;type:int;default:0"`
	Photos       string    `gorm:"column:photos;type:varchar(1000)"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime"`
}
