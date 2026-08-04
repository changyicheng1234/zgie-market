package model

import "time"

// Pvote 帖子投票记录，每用户每帖仅可投一票
type Pvote struct {
	PvoteID int       `gorm:"primary_key;column:pvoteID"`
	PostID  int       `gorm:"index:pvotepost;column:postID;unique_index:pvote_user_post"`
	Post    Post      `gorm:"association_foreignkey:postID;foreignkey:postID"`
	UserID  int       `gorm:"index:pvoteuser;column:userID;unique_index:pvote_user_post"`
	User    User      `gorm:"association_foreignkey:userID;foreignkey:userID"`
	Time    time.Time `gorm:"column:time;type:datetime"`
}
