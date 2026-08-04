package model

import "time"

type Resume struct {
	ResumeID    int       `gorm:"primary_key;column:resumeID"`
	UserID      int       `gorm:"index;column:userID;type:int"`
	Title       string    `gorm:"column:title;type:varchar(100)"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	FileName    string    `gorm:"column:file_name;type:varchar(255)"`
	CosKey      string    `gorm:"column:cos_key;type:varchar(512)"`
	FileSize    int64     `gorm:"column:file_size;type:bigint"`
	MimeType    string    `gorm:"column:mime_type;type:varchar(64)"`
	DownloadNum int       `gorm:"column:download_num;type:int;default:0"`
	ViewNum     int       `gorm:"column:view_num;type:int;default:0"`
	IsAnonymous bool      `gorm:"column:is_anonymous;type:bool;default:0"`
	IsDelete    bool      `gorm:"column:is_delete;type:bool;default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
}
