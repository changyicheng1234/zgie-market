package model

import "time"

type Product struct {
	ProductID   int       `gorm:"primary_key;column:productID"`
	UserID      int       `gorm:"index:idx_product_user;column:userID;type:int"`
	User        User      `gorm:"association_foreignkey:userID;foreignkey:userID"`
	Name        string    `gorm:"column:product_name;type:varchar(90)"`
	Description string    `gorm:"column:description;type:varchar(2000)"`
	PostTime    time.Time `gorm:"column:post_time;type:datetime"`
	Photos      string    `gorm:"column:photos;type:varchar(2000)"`
	ISSold      bool      `gorm:"column:is_sold;type:bool;default:false"`
	Price       int       `gorm:"column:price;type:int"`
	ISAnonymous bool      `gorm:"column:is_anonymous;type:bool;default:false"`
	ISDelete    bool      `gorm:"column:is_delete;type:bool;default:false"`
}
