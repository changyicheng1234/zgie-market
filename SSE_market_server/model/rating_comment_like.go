package model

import "time"

type RatingCommentLike struct {
	LikeID    int       `gorm:"primary_key;column:like_id"`
	CommentID int       `gorm:"unique_index:rating_comment_like_user;column:comment_id;type:int;not null"`
	UserID    int       `gorm:"unique_index:rating_comment_like_user;column:user_id;type:int;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime"`
}
