package common

import (
	"log"
	"loginTest/model"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// 仅用于从旧表迁移的一次性结构
type legacyRstar struct {
	RstarID int `gorm:"column:rstarID;primary_key"`
	PostID  int `gorm:"column:postID"`
	UserID  int `gorm:"column:userID"`
	StarNum int `gorm:"column:star_num"`
}

func (legacyRstar) TableName() string { return "rstars" }

var officialRatingCategories = []struct {
	Name        string
	Description string
	SortOrder   int
}{
	{"其他", "课程、食堂、宿舍等校园生活评分", 1},
}

func RunRatingMigration(db *gorm.DB) {
	db.AutoMigrate(
		&model.RatingCategory{},
		&model.RatingPost{},
		&model.RatingStar{},
		&model.RatingComment{},
		&model.RatingPostLike{},
		&model.RatingCommentLike{},
	)

	seedOfficialRatingCategories(db)

	var migrated int
	db.Model(&model.RatingPost{}).Count(&migrated)
	if migrated > 0 {
		recountRatingCategoryPosts(db)
		return
	}

	var legacyCount int
	db.Model(&model.Post{}).Where("`partition` = ?", "打分").Count(&legacyCount)
	if legacyCount == 0 {
		return
	}

	log.Printf("[rating] migrating %d legacy rating posts from posts table", legacyCount)
	migrateLegacyRatingPosts(db)
	recountRatingCategoryPosts(db)
	log.Printf("[rating] migration finished")
}

func seedOfficialRatingCategories(db *gorm.DB) {
	var n int
	db.Model(&model.RatingCategory{}).Count(&n)
	if n > 0 {
		return
	}
	now := time.Now()
	for _, item := range officialRatingCategories {
		db.Create(&model.RatingCategory{
			Name:        item.Name,
			Description: item.Description,
			CreatorID:   0,
			IsOfficial:  true,
			SortOrder:   item.SortOrder,
			CreatedAt:   now,
		})
	}
}

func ensureOtherCategory(db *gorm.DB) int {
	var cat model.RatingCategory
	db.Where("name = ?", "其他").First(&cat)
	if cat.CategoryID != 0 {
		return cat.CategoryID
	}
	now := time.Now()
	cat = model.RatingCategory{
		Name:        "其他",
		Description: "课程、食堂、宿舍等校园生活评分",
		CreatorID:   0,
		IsOfficial:  true,
		SortOrder:   1,
		CreatedAt:   now,
	}
	db.Create(&cat)
	return cat.CategoryID
}

func defaultOtherCategoryID(db *gorm.DB) int {
	return ensureOtherCategory(db)
}

// pruneRatingCategoriesToOtherOnly 仅保留「其他」，其余分类的帖子归入「其他」后删除。
func pruneRatingCategoriesToOtherOnly(db *gorm.DB) {
	otherID := ensureOtherCategory(db)
	if otherID == 0 {
		return
	}
	moved := db.Model(&model.RatingPost{}).
		Where("category_id != ?", otherID).
		Update("category_id", otherID).RowsAffected
	deleted := db.Where("name != ?", "其他").Delete(&model.RatingCategory{}).RowsAffected
	if moved > 0 || deleted > 0 {
		log.Printf("[rating] pruned categories: moved_posts=%d deleted_categories=%d", moved, deleted)
	}
	recountRatingCategoryPosts(db)
}

func resolveCategoryID(db *gorm.DB, tag string) int {
	_ = strings.TrimSpace(tag)
	return defaultOtherCategoryID(db)
}

func migrateLegacyRatingPosts(db *gorm.DB) {
	var posts []model.Post
	db.Where("`partition` = ?", "打分").Order("postID ASC").Find(&posts)
	for _, p := range posts {
		categoryID := resolveCategoryID(db, p.Tag)
		rp := model.RatingPost{
			RatingPostID: p.PostID,
			CategoryID:   categoryID,
			UserID:       p.UserID,
			Title:        p.Title,
			Content:      p.Ptext,
			CommentNum:   p.CommentNum,
			LikeNum:      p.LikeNum,
			BrowseNum:    p.BrowseNum,
			Photos:       p.Photos,
			CreatedAt:    p.PostTime,
		}
		db.Create(&rp)

		var stars []legacyRstar
		db.Where("postID = ?", p.PostID).Find(&stars)
		for _, s := range stars {
			db.Create(&model.RatingStar{
				RatingPostID: p.PostID,
				UserID:       s.UserID,
				StarNum:      s.StarNum,
			})
		}

		var comments []model.Pcomment
		db.Where("ptargetID = ?", p.PostID).Find(&comments)
		for _, c := range comments {
			starNum := 0
			var rs legacyRstar
			if db.Where("postID = ? AND userID = ?", p.PostID, c.UserID).First(&rs).Error == nil {
				starNum = rs.StarNum
			}
			rc := model.RatingComment{
				CommentID:    c.PcommentID,
				RatingPostID: p.PostID,
				UserID:       c.UserID,
				Content:      c.Pctext,
				StarNum:      starNum,
				LikeNum:      c.LikeNum,
				DenyNum:      c.DenyNum,
				CreatedAt:    c.Time,
			}
			db.Create(&rc)

			var pclikes []model.Pclike
			db.Where("pctargetID = ?", c.PcommentID).Find(&pclikes)
			for _, lk := range pclikes {
				db.Create(&model.RatingCommentLike{
					CommentID: c.PcommentID,
					UserID:    lk.UserID,
					CreatedAt: time.Now(),
				})
			}
		}

		var plikes []model.Plike
		db.Where("ptargetID = ?", p.PostID).Find(&plikes)
		for _, lk := range plikes {
			t := lk.Time
			if t.IsZero() {
				t = time.Now()
			}
			db.Create(&model.RatingPostLike{
				RatingPostID: p.PostID,
				UserID:       lk.UserID,
				CreatedAt:    t,
			})
		}

		db.Where("postID = ?", p.PostID).Delete(&legacyRstar{})
		db.Where("ptargetID = ?", p.PostID).Delete(&model.Pcomment{})
		db.Where("ptargetID = ?", p.PostID).Delete(&model.Plike{})
		db.Where("ptargetID = ?", p.PostID).Delete(&model.Psave{})
		db.Delete(&p)
	}
}

func recountRatingCategoryPosts(db *gorm.DB) {
	db.Model(&model.RatingCategory{}).Update("post_count", 0)
	var rows []struct {
		CategoryID int
		Cnt        int
	}
	db.Model(&model.RatingPost{}).
		Select("category_id, count(*) as cnt").
		Group("category_id").
		Scan(&rows)
	for _, r := range rows {
		db.Model(&model.RatingCategory{}).Where("category_id = ?", r.CategoryID).
			Update("post_count", r.Cnt)
	}
}
