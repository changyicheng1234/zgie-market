package service

import (
	"loginTest/model"

	"github.com/jinzhu/gorm"
)

// AdjustUserReceivedLike 调整用户收到的点赞累计（帖子被赞总和）
func AdjustUserReceivedLike(db *gorm.DB, authorUserID int, delta int) {
	if authorUserID == 0 || delta == 0 {
		return
	}
	db.Model(&model.User{}).Where("userID = ?", authorUserID).
		UpdateColumn("received_like_num", gorm.Expr("GREATEST(received_like_num + ?, 0)", delta))
}

// AdjustUserReceivedVote 调整用户收到的投票累计（其帖子被投票总和）
func AdjustUserReceivedVote(db *gorm.DB, authorUserID int, delta int) {
	if authorUserID == 0 || delta == 0 {
		return
	}
	db.Model(&model.User{}).Where("userID = ?", authorUserID).
		UpdateColumn("received_vote_num", gorm.Expr("GREATEST(received_vote_num + ?, 0)", delta))
}

// AdjustPostVoteNum 调整帖子票数字段
func AdjustPostVoteNum(db *gorm.DB, postID int, delta int) {
	if postID == 0 || delta == 0 {
		return
	}
	db.Model(&model.Post{}).Where("postID = ?", postID).
		UpdateColumn("vote_num", gorm.Expr("GREATEST(vote_num + ?, 0)", delta))
}

// OnPostDeleted 删帖前扣减作者排行统计
func OnPostDeleted(db *gorm.DB, post *model.Post) {
	if post.PostID == 0 || post.UserID == 0 {
		return
	}
	var voteCount int
	db.Model(&model.Pvote{}).Where("postID = ?", post.PostID).Count(&voteCount)
	AdjustUserReceivedLike(db, post.UserID, -post.LikeNum)
	AdjustUserReceivedVote(db, post.UserID, -voteCount)
}

// RebuildRankStats 从 pvotes / posts.like_num 全量重建排行冗余字段（部署或修复时调用）
func RebuildRankStats(db *gorm.DB) {
	db.Exec(`
		UPDATE posts p
		LEFT JOIN (
			SELECT postID, COUNT(*) AS cnt FROM pvotes GROUP BY postID
		) v ON p.postID = v.postID
		SET p.vote_num = COALESCE(v.cnt, 0)
	`)

	var users []model.User
	db.Select("userID").Find(&users)
	for _, u := range users {
		var likeSum int
		db.Model(&model.Post{}).Where("userID = ?", u.UserID).
			Select("COALESCE(SUM(like_num), 0)").Row().Scan(&likeSum)

		var voteSum int
		db.Table("pvotes").
			Joins("INNER JOIN posts ON posts.postID = pvotes.postID").
			Where("posts.userID = ?", u.UserID).
			Count(&voteSum)

		db.Model(&model.User{}).Where("userID = ?", u.UserID).Updates(map[string]interface{}{
			"received_like_num": likeSum,
			"received_vote_num": voteSum,
		})
	}
}
