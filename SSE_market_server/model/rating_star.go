package model

type RatingStar struct {
	StarID       int `gorm:"primary_key;column:star_id"`
	RatingPostID int `gorm:"unique_index:rating_star_post_user;column:rating_post_id;type:int;not null"`
	UserID       int `gorm:"unique_index:rating_star_post_user;column:user_id;type:int;not null"`
	StarNum      int `gorm:"column:star_num;type:int;not null"`
}
