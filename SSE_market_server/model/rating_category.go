package model

import "time"

type RatingCategory struct {
	CategoryID  int       `gorm:"primary_key;column:category_id"`
	Name        string    `gorm:"column:name;type:varchar(50);unique_index"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	CreatorID   int       `gorm:"column:creator_id;type:int;default:0"`
	PostCount   int       `gorm:"column:post_count;type:int;default:0"`
	IsOfficial  bool      `gorm:"column:is_official;type:bool;default:0"`
	SortOrder   int       `gorm:"column:sort_order;type:int;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
}
