package common

import (
	"fmt"
	//_ "github.com/alexbrainman/odbc"
	"loginTest/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// 使用viper从配置文件中读取数据库配置
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	// 若没有相应数据库，运行时将根据对应结构体自动创建数据库
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.Plike{})
	db.AutoMigrate(&model.Psave{})
	db.AutoMigrate(&model.Cclike{})
	db.AutoMigrate(&model.Ccdeny{})
	db.AutoMigrate(&model.Pclike{})
	db.AutoMigrate(&model.Pcdeny{})
	db.AutoMigrate(&model.Pcomment{})
	db.AutoMigrate(&model.Ccomment{})
	db.AutoMigrate(&model.Pbrowse{})
	db.AutoMigrate(&model.Admin{})
	db.AutoMigrate(&model.Feedback{})
	db.AutoMigrate(&model.Notice{})
	db.AutoMigrate(&model.Sue{})
	db.AutoMigrate(&model.CDKey{})
	db.AutoMigrate(&model.Tag{})
	db.AutoMigrate(&model.ChatMsg{})
	db.AutoMigrate(&model.OAuth2App{})
	db.AutoMigrate(&model.OAuth2Code{})
	db.AutoMigrate(&model.OAuth2Token{})
	db.AutoMigrate(&model.OAuth2AuthRecord{})
	db.AutoMigrate(&model.Pvote{})
	db.AutoMigrate(&model.Resume{})
	RunRatingMigration(db)
	// 历史简历：浏览次数初始化为下载量 × 3（仅 view_num 仍为 0 的行，避免重复覆盖）
	db.Exec("UPDATE resumes SET view_num = download_num * 3 WHERE view_num = 0 AND download_num > 0")
	db.Model(&model.Pcomment{}).AddForeignKey("ptargetID", "posts(postID)", "CASCADE", "CASCADE")
	db.Model(&model.Ccomment{}).AddForeignKey("ctargetID", "pcomments(pcommentID)", "CASCADE", "CASCADE")
	db.Model(&model.Plike{}).AddForeignKey("ptargetID", "posts(postID)", "CASCADE", "CASCADE")
	db.Model(&model.Cclike{}).AddForeignKey("cctargetID", "ccomments(ccommentID)", "CASCADE", "CASCADE")
	db.Model(&model.Pclike{}).AddForeignKey("pctargetID", "pcomments(pcommentID)", "CASCADE", "CASCADE")
	db.Model(&model.Ccdeny{}).AddForeignKey("cctargetID", "ccomments(ccommentID)", "CASCADE", "CASCADE")
	db.Model(&model.Pcdeny{}).AddForeignKey("pctargetID", "pcomments(pcommentID)", "CASCADE", "CASCADE")
	db.Model(&model.Pvote{}).AddForeignKey("postID", "posts(postID)", "CASCADE", "CASCADE")

	//db.Model(&model.Post{}).AddForeignKey("userID", "users(userID)", "CASCADE", "CASCADE")
	//db.Model(&model.Pcomment{}).AddForeignKey("userID", "users(userID)", "CASCADE", "CASCADE")
	//db.Model(&model.Ccomment{}).AddForeignKey("userID", "users(userID)", "CASCADE", "CASCADE")
	//db.Model(&model.Plike{}).AddForeignKey("userID", "users(userID)", "CASCADE", "CASCADE")
	//db.Model(model.Cclike{}).AddForeignKey("userID", "users(userID)", "CASCADE", "CASCADE")

	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
