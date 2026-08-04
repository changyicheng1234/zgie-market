package model

import "time"

// OAuth2应用信息表
type OAuth2App struct {
	AppID       string    `gorm:"primary_key;column:app_id;type:varchar(100)"`
	AppName     string    `gorm:"column:app_name;type:varchar(100)"`
	AppSecret   string    `gorm:"column:app_secret;type:varchar(255)"`
	AppIcon     string    `gorm:"column:app_icon;type:varchar(255)"`
	RedirectURI string    `gorm:"column:redirect_uri;type:varchar(255)"`
	Description string    `gorm:"column:description;type:varchar(500)"`
	IsActive    bool      `gorm:"column:is_active;type:tinyint(1);default:1"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
}

// 指定表名
func (OAuth2App) TableName() string {
	return "oauth2_apps"
}

// OAuth2临时授权码表
type OAuth2Code struct {
	CodeID      int       `gorm:"primary_key;column:code_id"`
	Code        string    `gorm:"column:code;type:varchar(100);unique"`
	AppID       string    `gorm:"column:app_id;type:varchar(100)"`
	UserID      int       `gorm:"column:user_id;type:int"`
	Scope       string    `gorm:"column:scope;type:varchar(100)"`
	RedirectURI string    `gorm:"column:redirect_uri;type:varchar(255)"`
	State       string    `gorm:"column:state;type:varchar(100)"`
	ExpiresAt   time.Time `gorm:"column:expires_at;type:datetime"`
	IsUsed      bool      `gorm:"column:is_used;type:tinyint(1);default:0"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
}

// 指定表名
func (OAuth2Code) TableName() string {
	return "oauth2_codes"
}

// OAuth2临时访问令牌表
type OAuth2Token struct {
	TokenID      int       `gorm:"primary_key;column:token_id"`
	AccessToken  string    `gorm:"column:access_token;type:varchar(255);unique"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(255)"`
	AppID        string    `gorm:"column:app_id;type:varchar(100)"`
	UserID       int       `gorm:"column:user_id;type:int"`
	Scope        string    `gorm:"column:scope;type:varchar(100)"`
	ExpiresAt    time.Time `gorm:"column:expires_at;type:datetime"`
	IsActive     bool      `gorm:"column:is_active;type:tinyint(1);default:1"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime"`
}

// 指定表名
func (OAuth2Token) TableName() string {
	return "oauth2_tokens"
}

// OAuth2授权记录表
type OAuth2AuthRecord struct {
	RecordID     int       `gorm:"primary_key;column:record_id"`
	UserID       int       `gorm:"column:user_id;type:int"`
	AppID        string    `gorm:"column:app_id;type:varchar(100)"`
	AppName      string    `gorm:"column:app_name;type:varchar(100)"`
	Scope        string    `gorm:"column:scope;type:varchar(100)"`
	IPAddress    string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent    string    `gorm:"column:user_agent;type:varchar(500)"`
	Status       string    `gorm:"column:status;type:enum('authorized','expired');default:'authorized'"`
	AuthorizedAt time.Time `gorm:"column:authorized_at;type:datetime"`
	ExpiresAt    time.Time `gorm:"column:expires_at;type:datetime"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime"`
}

// 指定表名
func (OAuth2AuthRecord) TableName() string {
	return "oauth2_auth_records"
}
