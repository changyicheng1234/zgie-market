package controller

import (
	"fmt"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

type RatingBrowseMsg struct {
	CategoryID int    `json:"categoryID"`
	Searchinfo string `json:"searchinfo"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type RatingCategoryCreateMsg struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SyncFeedPost *bool  `json:"syncFeedPost"`
}

type RatingPostCreateMsg struct {
	CategoryID int    `json:"categoryID"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Photos     string `json:"photos"`
}

type RatingPostIDMsg struct {
	RatingPostID int `json:"ratingPostID"`
}

type RatingStarSubmitMsg struct {
	RatingPostID int `json:"ratingPostID"`
	StarNum      int `json:"starNum"`
}

type RatingCommentCreateMsg struct {
	RatingPostID int    `json:"ratingPostID"`
	Content      string `json:"content"`
	StarNum      int    `json:"starNum"`
}

type RatingCommentListMsg struct {
	RatingPostID int `json:"ratingPostID"`
}

type RatingLikeMsg struct {
	RatingPostID int `json:"ratingPostID"`
}

type RatingCommentLikeMsg struct {
	CommentID int `json:"commentID"`
}

type RatingCategoryResponse struct {
	CategoryID  int    `json:"categoryID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PostCount   int    `json:"postCount"`
	IsOfficial  bool   `json:"isOfficial"`
	SortOrder   int    `json:"sortOrder"`
}

type RatingPostListItem struct {
	RatingPostID   int       `json:"ratingPostID"`
	CategoryID     int       `json:"categoryID"`
	CategoryName   string    `json:"categoryName"`
	UserID         int       `json:"userID"`
	UserName       string    `json:"userName"`
	UserScore      int       `json:"userScore"`
	UserTelephone  string    `json:"userTelephone"`
	UserAvatar     string    `json:"userAvatar"`
	UserIdentity   string    `json:"userIdentity"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Like           int       `json:"like"`
	Comment        int       `json:"comment"`
	Browse         int       `json:"browse"`
	CreatedAt      time.Time `json:"createdAt"`
	IsLiked        bool      `json:"isLiked"`
	Photos         string    `json:"photos"`
	AverageRating  float64   `json:"averageRating"`
	Stars          [5]int    `json:"stars"`
	RaterCount     int       `json:"raterCount"`
}

type RatingPostDetailResponse struct {
	RatingPostListItem
	UserStar int `json:"userStar"`
}

type RatingCommentResponse struct {
	CommentID       int       `json:"commentID"`
	AuthorID        int       `json:"authorID"`
	Author          string    `json:"author"`
	AuthorTelephone string    `json:"authorTelephone"`
	AuthorScore     int       `json:"authorScore"`
	AuthorAvatar    string    `json:"authorAvatar"`
	AuthorIdentity  string    `json:"authorIdentity"`
	CommentTime     time.Time `json:"commentTime"`
	Content         string    `json:"content"`
	LikeNum         int       `json:"likeNum"`
	IsLiked         bool      `json:"isLiked"`
	AuthorRating    int       `json:"authorRating"`
}

func RatingCategoryList(c *gin.Context) {
	db := common.GetDB()
	var list []model.RatingCategory
	db.Order("created_at DESC, category_id DESC").Find(&list)
	out := make([]RatingCategoryResponse, 0, len(list))
	for _, cat := range list {
		out = append(out, RatingCategoryResponse{
			CategoryID:  cat.CategoryID,
			Name:        cat.Name,
			Description: cat.Description,
			PostCount:   cat.PostCount,
			IsOfficial:  cat.IsOfficial,
			SortOrder:   cat.SortOrder,
		})
	}
	c.JSON(http.StatusOK, out)
}

func RatingCategoryCreate(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingCategoryCreateMsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	name := strings.TrimSpace(msg.Name)
	if utf8.RuneCountInString(name) < 2 || utf8.RuneCountInString(name) > 30 {
		response.Response(c, http.StatusBadRequest, 400, nil, "分类名称需 2～30 个字")
		return
	}
	desc := strings.TrimSpace(msg.Description)
	if utf8.RuneCountInString(desc) > 200 {
		response.Response(c, http.StatusBadRequest, 400, nil, "简介最多 200 字")
		return
	}
	db := common.GetDB()
	var exists model.RatingCategory
	db.Where("name = ?", name).First(&exists)
	if exists.CategoryID != 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "该分类名称已存在")
		return
	}
	var maxOrder int
	db.Model(&model.RatingCategory{}).Select("COALESCE(MAX(sort_order), 0)").Row().Scan(&maxOrder)
	syncFeed := true
	if msg.SyncFeedPost != nil {
		syncFeed = *msg.SyncFeedPost
	}

	cat := model.RatingCategory{
		Name:        name,
		Description: desc,
		CreatorID:   user.UserID,
		IsOfficial:  false,
		SortOrder:   maxOrder + 1,
		CreatedAt:   time.Now(),
	}

	tx := db.Begin()
	if err := tx.Create(&cat).Error; err != nil {
		tx.Rollback()
		response.Response(c, http.StatusInternalServerError, 500, nil, "创建失败")
		return
	}
	if syncFeed {
		if err := createRatingCategoryFeedPost(tx, user, cat); err != nil {
			tx.Rollback()
			response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "创建失败")
		return
	}

	c.JSON(http.StatusOK, RatingCategoryResponse{
		CategoryID:  cat.CategoryID,
		Name:        cat.Name,
		Description: cat.Description,
		PostCount:   0,
		IsOfficial:  false,
		SortOrder:   cat.SortOrder,
	})
}

func RatingPostList(c *gin.Context) {
	db := common.GetDB()
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingBrowseMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	limit, offset := msg.Limit, msg.Offset
	if limit <= 0 || limit > 50 {
		limit = 15
	}
	if offset < 0 {
		offset = 0
	}
	var posts []model.RatingPost
	q := ratingApplyBrowseFilter(db, db.Model(&model.RatingPost{}), &msg)
	q.Order("rating_post_id DESC").Offset(offset).Limit(limit).Find(&posts)
	if len(posts) == 0 {
		c.JSON(http.StatusOK, []RatingPostListItem{})
		return
	}
	c.JSON(http.StatusOK, buildRatingPostListItems(db, posts, user.UserID))
}

func RatingPostCount(c *gin.Context) {
	db := common.GetDB()
	var msg RatingBrowseMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	q := ratingApplyBrowseFilter(db, db.Model(&model.RatingPost{}), &msg)
	var count int
	q.Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func RatingPostCreate(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingPostCreateMsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	title := strings.TrimSpace(msg.Title)
	if title == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题不能为空")
		return
	}
	if utf8.RuneCountInString(title) > 50 {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题最多 50 个字")
		return
	}
	content := strings.TrimSpace(msg.Content)
	if utf8.RuneCountInString(content) > 5000 {
		response.Response(c, http.StatusBadRequest, 400, nil, "内容过长")
		return
	}
	if msg.CategoryID <= 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "请选择分类")
		return
	}
	db := common.GetDB()
	var cat model.RatingCategory
	if db.Where("category_id = ?", msg.CategoryID).First(&cat).Error != nil || cat.CategoryID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "分类不存在")
		return
	}
	now := time.Now()
	rp := model.RatingPost{
		CategoryID: msg.CategoryID,
		UserID:     user.UserID,
		Title:      title,
		Content:    content,
		Photos:     msg.Photos,
		CreatedAt:  now,
	}
	db.Create(&rp)
	db.Model(&cat).Update("post_count", gorm.Expr("post_count + ?", 1))
	c.JSON(http.StatusOK, gin.H{"ratingPostID": rp.RatingPostID})
}

func RatingPostDetail(c *gin.Context) {
	db := common.GetDB()
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingPostIDMsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil || msg.RatingPostID <= 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "无效的评分帖 ID")
		return
	}
	var rp model.RatingPost
	if db.Where("rating_post_id = ?", msg.RatingPostID).First(&rp).Error != nil || rp.RatingPostID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "评分帖不存在")
		return
	}
	db.Model(&rp).Update("browse_num", gorm.Expr("browse_num + ?", 1))
	rp.BrowseNum++
	items := buildRatingPostListItems(db, []model.RatingPost{rp}, user.UserID)
	if len(items) == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "评分帖不存在")
		return
	}
	detail := RatingPostDetailResponse{RatingPostListItem: items[0]}
	var userStar model.RatingStar
	db.Where("rating_post_id = ? AND user_id = ?", rp.RatingPostID, user.UserID).First(&userStar)
	detail.UserStar = userStar.StarNum
	c.JSON(http.StatusOK, detail)
}

func RatingPostDelete(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingPostIDMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	if msg.RatingPostID <= 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "无效的评分帖 ID")
		return
	}
	db := common.GetDB()
	var rp model.RatingPost
	db.Where("rating_post_id = ?", msg.RatingPostID).First(&rp)
	if rp.RatingPostID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "评分帖不存在")
		return
	}
	if rp.UserID != user.UserID {
		response.Response(c, http.StatusForbidden, 403, nil, "无权删除")
		return
	}
	categoryID := rp.CategoryID
	db.Where("rating_post_id = ?", rp.RatingPostID).Delete(&model.RatingStar{})
	db.Where("rating_post_id = ?", rp.RatingPostID).Delete(&model.RatingComment{})
	db.Where("rating_post_id = ?", rp.RatingPostID).Delete(&model.RatingPostLike{})
	db.Delete(&rp)
	db.Model(&model.RatingCategory{}).Where("category_id = ?", categoryID).
		Update("post_count", gorm.Expr("GREATEST(post_count - 1, 0)"))
	response.Response(c, http.StatusOK, 200, nil, "删除成功")
}

func RatingPostLike(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingLikeMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	if msg.RatingPostID <= 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "无效的评分帖 ID")
		return
	}
	db := common.GetDB()
	var rp model.RatingPost
	db.Where("rating_post_id = ?", msg.RatingPostID).First(&rp)
	if rp.RatingPostID == 0 {
		return
	}
	var like model.RatingPostLike
	db.Where("rating_post_id = ? AND user_id = ?", msg.RatingPostID, user.UserID).First(&like)
	if like.LikeID != 0 {
		db.Delete(&like)
		db.Model(&rp).Update("like_num", gorm.Expr("GREATEST(like_num - 1, 0)"))
	} else {
		db.Create(&model.RatingPostLike{
			RatingPostID: msg.RatingPostID,
			UserID:       user.UserID,
			CreatedAt:    time.Now(),
		})
		db.Model(&rp).Update("like_num", gorm.Expr("like_num + 1"))
	}
	response.Response(c, http.StatusOK, 200, nil, "操作成功")
}

func RatingStarSubmit(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingStarSubmitMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	if msg.RatingPostID <= 0 || msg.StarNum < 1 || msg.StarNum > 5 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评分须为 1～5 星")
		return
	}
	db := common.GetDB()
	var rp model.RatingPost
	db.Where("rating_post_id = ?", msg.RatingPostID).First(&rp)
	if rp.RatingPostID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "评分帖不存在")
		return
	}
	var star model.RatingStar
	err := db.Where("rating_post_id = ? AND user_id = ?", msg.RatingPostID, user.UserID).First(&star).Error
	if err != nil {
		db.Create(&model.RatingStar{
			RatingPostID: msg.RatingPostID,
			UserID:       user.UserID,
			StarNum:      msg.StarNum,
		})
	} else {
		db.Model(&star).Update("star_num", msg.StarNum)
	}
	response.Response(c, http.StatusOK, 200, nil, "评分成功")
}

func RatingCommentList(c *gin.Context) {
	db := common.GetDB()
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingCommentListMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	if msg.RatingPostID <= 0 {
		c.JSON(http.StatusOK, []RatingCommentResponse{})
		return
	}
	var comments []model.RatingComment
	db.Where("rating_post_id = ?", msg.RatingPostID).Order("comment_id DESC").Find(&comments)
	if len(comments) == 0 {
		c.JSON(http.StatusOK, []RatingCommentResponse{})
		return
	}
	commentIDs := make([]int, 0, len(comments))
	userIDs := make(map[int]struct{})
	for _, cm := range comments {
		commentIDs = append(commentIDs, cm.CommentID)
		userIDs[cm.UserID] = struct{}{}
	}
	userMap := ratingBatchUserMap(db, userIDs)
	likeMap := ratingBatchCommentLiked(db, user.UserID, commentIDs)
	out := make([]RatingCommentResponse, 0, len(comments))
	for _, cm := range comments {
		u := ratingGetUser(userMap, cm.UserID)
		out = append(out, RatingCommentResponse{
			CommentID:       cm.CommentID,
			AuthorID:        u.UserID,
			Author:          u.Name,
			AuthorTelephone: u.Phone,
			AuthorScore:     u.Score,
			AuthorAvatar:    u.AvatarURL,
			AuthorIdentity:  u.Identity,
			CommentTime:     cm.CreatedAt,
			Content:         cm.Content,
			LikeNum:         cm.LikeNum,
			IsLiked:         likeMap[cm.CommentID],
			AuthorRating:    cm.StarNum,
		})
	}
	c.JSON(http.StatusOK, out)
}

func RatingCommentCreate(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingCommentCreateMsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	content := strings.TrimSpace(msg.Content)
	if content == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论内容不能为空")
		return
	}
	if utf8.RuneCountInString(content) > 1000 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评论内容过长")
		return
	}
	if msg.StarNum < 0 || msg.StarNum > 5 {
		response.Response(c, http.StatusBadRequest, 400, nil, "评分须在 0～5 之间")
		return
	}
	db := common.GetDB()
	var rp model.RatingPost
	db.Where("rating_post_id = ?", msg.RatingPostID).First(&rp)
	if rp.RatingPostID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "评分帖不存在")
		return
	}
	currentDateTime := time.Now()
	if user.Banend.After(currentDateTime) {
		response.Response(c, http.StatusBadRequest, 400, nil, "你尚处于禁言状态中，不得评论")
		return
	}
	cm := model.RatingComment{
		RatingPostID: msg.RatingPostID,
		UserID:       user.UserID,
		Content:      content,
		StarNum:      msg.StarNum,
		CreatedAt:    time.Now(),
	}
	db.Create(&cm)
	db.Model(&rp).Update("comment_num", gorm.Expr("comment_num + 1"))

	if msg.StarNum >= 1 && msg.StarNum <= 5 {
		var star model.RatingStar
		if db.Where("rating_post_id = ? AND user_id = ?", msg.RatingPostID, user.UserID).First(&star).Error != nil {
			db.Create(&model.RatingStar{
				RatingPostID: msg.RatingPostID,
				UserID:       user.UserID,
				StarNum:      msg.StarNum,
			})
		} else {
			db.Model(&star).Update("star_num", msg.StarNum)
		}
		cm.StarNum = msg.StarNum
	}

	if rp.UserID != user.UserID {
		notice := model.Notice{
			Receiver: rp.UserID,
			Sender:   user.UserID,
			Type:     "rating_comment",
			Ntext:    content,
			Time:     time.Now(),
			Read:     false,
			Target:   cm.CommentID,
		}
		db.Create(&notice)
		var receiver model.User
		db.Where("userID = ?", rp.UserID).First(&receiver)
		if receiver.UserID != 0 && !receiver.ISEmailNotificationBlocked {
			EmailAlerter(receiver.Email, 1, content)
		}
	}

	c.JSON(http.StatusOK, RatingCommentResponse{
		CommentID:       cm.CommentID,
		AuthorID:        user.UserID,
		Author:          user.Name,
		AuthorTelephone: user.Phone,
		AuthorScore:     user.Score,
		AuthorAvatar:    user.AvatarURL,
		AuthorIdentity:  user.Identity,
		CommentTime:     cm.CreatedAt,
		Content:         cm.Content,
		LikeNum:         0,
		IsLiked:         false,
		AuthorRating:    cm.StarNum,
	})
}

func RatingCommentLike(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}
	var msg RatingCommentLikeMsg
	_ = c.ShouldBindBodyWith(&msg, binding.JSON)
	if msg.CommentID <= 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "无效的评论 ID")
		return
	}
	db := common.GetDB()
	var cm model.RatingComment
	db.Where("comment_id = ?", msg.CommentID).First(&cm)
	if cm.CommentID == 0 {
		return
	}
	var like model.RatingCommentLike
	db.Where("comment_id = ? AND user_id = ?", msg.CommentID, user.UserID).First(&like)
	if like.LikeID != 0 {
		db.Delete(&like)
		db.Model(&cm).Update("like_num", gorm.Expr("GREATEST(like_num - 1, 0)"))
	} else {
		db.Create(&model.RatingCommentLike{
			CommentID: msg.CommentID,
			UserID:    user.UserID,
			CreatedAt: time.Now(),
		})
		db.Model(&cm).Update("like_num", gorm.Expr("like_num + 1"))
	}
	response.Response(c, http.StatusOK, 200, nil, "操作成功")
}

func buildRatingPostListItems(db *gorm.DB, posts []model.RatingPost, viewerUserID int) []RatingPostListItem {
	postIDs := make([]int, 0, len(posts))
	categoryIDs := make(map[int]struct{})
	userIDs := make(map[int]struct{})
	for _, p := range posts {
		postIDs = append(postIDs, p.RatingPostID)
		categoryIDs[p.CategoryID] = struct{}{}
		if p.UserID != 0 {
			userIDs[p.UserID] = struct{}{}
		}
	}
	catMap := ratingBatchCategoryMap(db, categoryIDs)
	userMap := ratingBatchUserMap(db, userIDs)
	likedSet := ratingBatchPostLiked(db, viewerUserID, postIDs)
	starsMap := ratingBatchStars(db, postIDs)

	out := make([]RatingPostListItem, 0, len(posts))
	for _, p := range posts {
		u := ratingGetUser(userMap, p.UserID)
		cat := catMap[p.CategoryID]
		stars, avg, cnt := ratingAggregateStars(starsMap[p.RatingPostID])
		_, isLiked := likedSet[p.RatingPostID]
		out = append(out, RatingPostListItem{
			RatingPostID:  p.RatingPostID,
			CategoryID:    p.CategoryID,
			CategoryName:  cat.Name,
			UserID:        p.UserID,
			UserName:      u.Name,
			UserScore:     u.Score,
			UserTelephone: u.Phone,
			UserAvatar:    u.AvatarURL,
			UserIdentity:  u.Identity,
			Title:         p.Title,
			Content:       p.Content,
			Like:          p.LikeNum,
			Comment:       p.CommentNum,
			Browse:        p.BrowseNum,
			CreatedAt:     p.CreatedAt,
			IsLiked:       isLiked,
			Photos:        p.Photos,
			AverageRating: avg,
			Stars:         stars,
			RaterCount:    cnt,
		})
	}
	return out
}

func ratingBatchCategoryMap(db *gorm.DB, ids map[int]struct{}) map[int]model.RatingCategory {
	slice := make([]int, 0, len(ids))
	for id := range ids {
		slice = append(slice, id)
	}
	m := make(map[int]model.RatingCategory)
	if len(slice) == 0 {
		return m
	}
	var cats []model.RatingCategory
	db.Where("category_id IN (?)", slice).Find(&cats)
	for _, c := range cats {
		m[c.CategoryID] = c
	}
	return m
}

func ratingBatchUserMap(db *gorm.DB, ids map[int]struct{}) map[int]model.User {
	slice := make([]int, 0, len(ids))
	for id := range ids {
		slice = append(slice, id)
	}
	m := make(map[int]model.User)
	if len(slice) == 0 {
		return m
	}
	var users []model.User
	db.Where("userID IN (?)", slice).Find(&users)
	for _, u := range users {
		m[u.UserID] = u
	}
	return m
}

func ratingGetUser(m map[int]model.User, id int) model.User {
	if u, ok := m[id]; ok {
		return u
	}
	if id == 0 {
		return model.User{Name: "管理员", Phone: "11111111111"}
	}
	return model.User{}
}

func ratingBatchPostLiked(db *gorm.DB, userID int, postIDs []int) map[int]struct{} {
	set := make(map[int]struct{})
	if userID == 0 || len(postIDs) == 0 {
		return set
	}
	var likes []model.RatingPostLike
	db.Where("user_id = ? AND rating_post_id IN (?)", userID, postIDs).Find(&likes)
	for _, l := range likes {
		set[l.RatingPostID] = struct{}{}
	}
	return set
}

func ratingBatchStars(db *gorm.DB, postIDs []int) map[int][]model.RatingStar {
	m := make(map[int][]model.RatingStar)
	if len(postIDs) == 0 {
		return m
	}
	var list []model.RatingStar
	db.Where("rating_post_id IN (?)", postIDs).Find(&list)
	for _, s := range list {
		m[s.RatingPostID] = append(m[s.RatingPostID], s)
	}
	return m
}

func ratingAggregateStars(stars []model.RatingStar) ([5]int, float64, int) {
	var cnt [5]int
	var total float64
	var n int
	for _, s := range stars {
		if s.StarNum < 1 || s.StarNum > 5 {
			continue
		}
		cnt[s.StarNum-1]++
		total += float64(s.StarNum)
		n++
	}
	avg := 0.0
	if n > 0 {
		avg = total / float64(n)
	}
	return cnt, avg, n
}

func ratingBatchCommentLiked(db *gorm.DB, userID int, commentIDs []int) map[int]bool {
	m := make(map[int]bool)
	if userID == 0 || len(commentIDs) == 0 {
		return m
	}
	var likes []model.RatingCommentLike
	db.Where("user_id = ? AND comment_id IN (?)", userID, commentIDs).Find(&likes)
	for _, l := range likes {
		m[l.CommentID] = true
	}
	return m
}

// ratingApplyBrowseFilter 列表/计数共用：支持按分区筛选；搜索词匹配标题/正文，且匹配分区名时包含该分区下全部帖子。
func ratingApplyBrowseFilter(db *gorm.DB, q *gorm.DB, msg *RatingBrowseMsg) *gorm.DB {
	if msg.CategoryID > 0 {
		q = q.Where("category_id = ?", msg.CategoryID)
	}
	s := strings.TrimSpace(msg.Searchinfo)
	if s == "" {
		return q
	}
	like := "%" + s + "%"
	if msg.CategoryID > 0 {
		return q.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	var catIDs []int
	db.Model(&model.RatingCategory{}).Where("name LIKE ?", like).Pluck("category_id", &catIDs)
	if len(catIDs) > 0 {
		return q.Where("(title LIKE ? OR content LIKE ?) OR category_id IN (?)", like, like, catIDs)
	}
	return q.Where("title LIKE ? OR content LIKE ?", like, like)
}

func createRatingCategoryFeedPost(db *gorm.DB, user model.User, cat model.RatingCategory) error {
	title := common.FeedTitleWithPrefix("新增分区-", cat.Name)
	content := ratingCategoryFeedPostContent(cat.Name, cat.Description, cat.CategoryID)
	return common.CreateHomeFeedPost(db, user, title, content)
}

func ratingCategoryFeedPostContent(name, desc string, categoryID int) string {
	link := fmt.Sprintf("%s/rating?category=%d", common.NewFrontendBaseURL, categoryID)
	var b strings.Builder
	b.WriteString("打分专区新增了分区「")
	b.WriteString(name)
	b.WriteString("」，欢迎大家来评分、留言。\n\n")
	b.WriteString("直达链接：\n")
	b.WriteString(link)
	b.WriteString("\n")
	if strings.TrimSpace(desc) != "" {
		b.WriteString("\n分区简介：")
		b.WriteString(strings.TrimSpace(desc))
		b.WriteString("\n")
	}
	return b.String()
}
