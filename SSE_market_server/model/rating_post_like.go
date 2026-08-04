package model

import "time"

type RatingPostLike struct {
	LikeID       int       `gorm:"primary_key;column:like_id"`
	RatingPostID int       `gorm:"unique_index:rating_post_like_user;column:rating_post_id;type:int;not null"`
	UserID       int       `gorm:"unique_index:rating_post_like_user;column:user_id;type:int;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime"`
}
