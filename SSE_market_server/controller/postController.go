package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"loginTest/api"
	"loginTest/common"
	"loginTest/config"
	"loginTest/core"
	"loginTest/dto"
	"loginTest/model"
	"loginTest/response"
	"loginTest/service"
	"loginTest/util"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

const (
	// 课程交流分区名称
	CoursePartition = "课程交流"
)

type PostMsg struct {
	UserTelephone string
	Title         string
	Content       string
	Partition     string
	Photos        string
	TagList       string
	IsPrivate     bool `json:"isPrivate"`
}

func Post(c *gin.Context) {
	db := common.GetDB()
	var requestPostMsg PostMsg
	c.Bind(&requestPostMsg)
	// 获取参数
	userTelephone := requestPostMsg.UserTelephone
	title := requestPostMsg.Title
	content := requestPostMsg.Content
	partition := requestPostMsg.Partition
	photos := requestPostMsg.Photos
	tagList := requestPostMsg.TagList
	tags := strings.Split(tagList, "|")
	tagString := strings.Join(tags, ",")
	// 验证数据
	if len(userTelephone) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "返回的手机号为空")
		return
	}
	if !isTelephoneExist(db, userTelephone) {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}
	if len(title) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题不能为空")
		return
	}
	// 这里不能直接用len()否则会出现中文字符的计算错误
	if utf8.RuneCountInString(title) > 30 {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题最多为30个字")
		return
	}

	if len(content) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "内容不能为空")
		return
	}

	if utf8.RuneCountInString(content) > config.PostContentMaxRunes {
		response.Response(c, http.StatusBadRequest, 400, nil, fmt.Sprintf("内容最多为%d个字", config.PostContentMaxRunes))
		return
	}

	if len(partition) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "分区不能为空")
		return
	}

	// if api.GetSuggestion(title) == "Block" || api.GetSuggestion(content) == "Block" {
	// 	// response.Response(c, http.StatusBadRequest, 400, nil, "标题或内容可能含有不良信息")
	// }

	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	if user.UserID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}

	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}

	currentDateTime := time.Now()
	if user.Banend.After(currentDateTime) {
		response.Response(c, http.StatusBadRequest, 400, nil, "你尚处于禁言状态中，不得发帖")
		return
	}

	if partition == CoursePartition {
		response.Response(c, http.StatusBadRequest, 400, nil, "权限不足")
		return
	}

	newPost := model.Post{
		UserID:     int(user.UserID),
		Partition:  partition,
		Title:      title,
		Ptext:      content,
		LikeNum:    0,
		CommentNum: 0,
		BrowseNum:  0,
		Heat:       0,
		PostTime:   time.Now(),
		Photos:     photos,
		Tag:        tagString,
		IsPrivate:  requestPostMsg.IsPrivate,
	}
	db.Create(&newPost)
	common.BloomFilter.Add([]byte(strconv.FormatUint(uint64(newPost.PostID), 10)))

	// 判断tag，看是否需要发通知
	tempTag := model.Tag{}
	db.Where("name = ?", tags[0]).First(&tempTag)
	if tempTag.TagID == 0 || tempTag.Type != "course" {
	} else {
		receiver := model.User{}
		db.Where("email = ?", tempTag.Value).First(&receiver)
		fmt.Println("receiver:" + receiver.Name + " email:" + receiver.Email)
		if receiver.UserID != 0 && receiver.Phone != userTelephone {
			newNotice := model.Notice{
				Receiver: receiver.UserID,
				Sender:   user.UserID,
				Type:     "post",
				Ntext:    title,
				Time:     time.Now(),
				Read:     false,
				Target:   newPost.PostID,
			}
			db.Create(&newNotice)
		}
	}
	response.Response(c, http.StatusOK, 200, nil, "发帖成功")
}

type PostResponse struct {
	PostID        uint
	Partition     string
	UserID        int
	UserName      string
	UserScore     int
	UserTelephone string
	UserAvatar    string
	UserIdentity  string
	Title         string
	Content       string
	Like          int
	Comment       int
	Browse        int
	Heat          float64
	PostTime      time.Time
	IsSaved       bool
	IsLiked       bool
	Photos        string
	Tag           string
	IsPrivate     bool
}

type BrowseMeg struct {
	UserTelephone string `json:"userTelephone"`
	Partition     string `json:"partition"`
	Searchinfo    string `json:"searchinfo"`
	Tag           string `json:"tag"`
	Searchsort    string `json:"searchsort"` // home | history | save
	OrderBy       string `json:"orderBy"`    // time_desc|time_asc|like_desc|like_asc|comment_desc|comment_asc
	TimeFrom      string `json:"timeFrom"`
	TimeTo        string `json:"timeTo"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
	PostType      string `json:"postType"`
}

func freshDB() *gorm.DB {
	return common.GetDB().New()
}

// viewerFromAuth 从 JWT 中间件注入的上下文取当前用户（以 userID 为准，不用手机号）
func viewerFromAuth(c *gin.Context) (model.User, bool) {
	value, exists := c.Get("user")
	if !exists {
		return model.User{}, false
	}
	user, ok := value.(model.User)
	if !ok || user.UserID == 0 {
		return model.User{}, false
	}
	return user, true
}

func canViewPost(post model.Post, viewerUserID int) bool {
	return !post.IsPrivate || post.UserID == viewerUserID
}

func applyPostVisibilityScope(query *gorm.DB, viewerUserID int) *gorm.DB {
	return query.Where("(is_private = ? OR userID = ?)", false, viewerUserID)
}

func loadViewablePost(db *gorm.DB, postID int, viewerUserID int) (model.Post, bool) {
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 || !canViewPost(post, viewerUserID) {
		return post, false
	}
	return post, true
}

func canViewPostByPcommentID(db *gorm.DB, pcommentID int, viewerUserID int) bool {
	var pcomment model.Pcomment
	db.Where("pcommentID = ?", pcommentID).First(&pcomment)
	if pcomment.PcommentID == 0 {
		return false
	}
	var post model.Post
	db.Where("postID = ?", pcomment.PtargetID).First(&post)
	return canViewPost(post, viewerUserID)
}

func canViewPostByCcommentID(db *gorm.DB, ccommentID int, viewerUserID int) bool {
	var ccomment model.Ccomment
	db.Where("ccommentID = ?", ccommentID).First(&ccomment)
	if ccomment.CcommentID == 0 {
		return false
	}
	var pcomment model.Pcomment
	db.Where("pcommentID = ?", ccomment.CtargetID).First(&pcomment)
	if pcomment.PcommentID == 0 {
		return false
	}
	var post model.Post
	db.Where("postID = ?", pcomment.PtargetID).First(&post)
	return canViewPost(post, viewerUserID)
}

func buildPostResponse(post model.Post, user model.User, isSaved bool, isLiked bool) PostResponse {
	return PostResponse{
		PostID:        uint(post.PostID),
		Partition:     post.Partition,
		UserID:        post.UserID,
		UserName:      user.Name,
		UserScore:     user.Score,
		UserTelephone: user.Phone,
		UserAvatar:    user.AvatarURL,
		UserIdentity:  user.Identity,
		Title:         post.Title,
		Content:       post.Ptext,
		Like:          post.LikeNum,
		Comment:       post.CommentNum,
		Browse:        post.BrowseNum,
		Heat:          post.Heat,
		PostTime:      post.PostTime,
		IsSaved:       isSaved,
		IsLiked:       isLiked,
		Photos:        post.Photos,
		Tag:           post.Tag,
		IsPrivate:     post.IsPrivate,
	}
}

func parseBrowseTime(s string, endOfDay bool) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, s, time.Local)
		if err != nil {
			continue
		}
		if layout == "2006-01-02" && endOfDay {
			t = t.Add(24*time.Hour - time.Nanosecond)
		}
		return t, true
	}
	return time.Time{}, false
}

func browseApplyTimeRange(query *gorm.DB, timeFrom, timeTo string) *gorm.DB {
	if t, ok := parseBrowseTime(timeFrom, false); ok {
		query = query.Where("post_time >= ?", t)
	}
	if t, ok := parseBrowseTime(timeTo, true); ok {
		query = query.Where("post_time <= ?", t)
	}
	return query
}

func browseApplyOrder(query *gorm.DB, orderBy string) *gorm.DB {
	switch orderBy {
	case "time_asc":
		return query.Order("post_time ASC")
	case "like_desc":
		return query.Order("like_num DESC, postID DESC")
	case "like_asc":
		return query.Order("like_num ASC, postID DESC")
	case "comment_desc":
		return query.Order("comment_num DESC, postID DESC")
	case "comment_asc":
		return query.Order("comment_num ASC, postID DESC")
	default:
		return query.Order("post_time DESC")
	}
}

func Browse(c *gin.Context) {
	viewer, ok := viewerFromAuth(c)
	if !ok {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	db := freshDB()
	var requestBrowseMeg BrowseMeg
	if err := c.ShouldBindJSON(&requestBrowseMeg); err != nil {
		_ = c.ShouldBind(&requestBrowseMeg)
	}
	searchsort := requestBrowseMeg.Searchsort
	if searchsort == "" {
		searchsort = "home"
	}
	requestBrowseMeg.Searchsort = searchsort

	posts := browseFetchPosts(db, &requestBrowseMeg, &viewer)
	if len(posts) == 0 {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	postIDs := make([]int, 0, len(posts))
	userIDs := make(map[int]struct{})
	for _, p := range posts {
		if !canViewPost(p, viewer.UserID) {
			continue
		}
		postIDs = append(postIDs, p.PostID)
		if p.UserID != 0 {
			userIDs[p.UserID] = struct{}{}
		}
	}
	savedSet := browseBatchSavedSet(db, viewer.UserID, postIDs)
	likedSet := browseBatchLikedSet(db, viewer.UserID, postIDs)
	userMap := browseBatchUserMap(db, userIDs)
	postResponses := make([]PostResponse, 0, len(posts))
	isSaveMode := (searchsort == "save")
	for _, post := range posts {
		if !canViewPost(post, viewer.UserID) {
			continue
		}
		user := browseGetUserForPost(userMap, post.UserID)
		_, isSaved := savedSet[post.PostID]
		if isSaveMode {
			isSaved = true
		}
		_, isLiked := likedSet[post.PostID]
		postResponses = append(postResponses, buildPostResponse(post, user, isSaved, isLiked))
	}
	c.JSON(http.StatusOK, postResponses)
}

type UserProfilePostsMsg struct {
	UserTelephone string `json:"userTelephone"`
	TargetUserID  int    `json:"targetUserID"`
	Limit         int    `json:"limit"`
	Offset        int    `json:"offset"`
}

// UserProfilePosts 个人主页发帖列表；对方开启「隐藏发帖」且浏览者非本人时返回空。
func UserProfilePosts(c *gin.Context) {
	db := freshDB()
	var req UserProfilePostsMsg
	_ = c.Bind(&req)
	if req.TargetUserID <= 0 {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	viewer, ok := viewerFromAuth(c)
	if !ok {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	var target model.User
	if err := db.Where("userID = ?", req.TargetUserID).First(&target).Error; err != nil {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	if target.HideProfilePosts && viewer.UserID != target.UserID {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	limit, offset := req.Limit, req.Offset
	if limit <= 0 || limit > 50 {
		limit = 15
	}
	if offset < 0 {
		offset = 0
	}
	var posts []model.Post
	query := db.Where("userID = ?", target.UserID)
	if viewer.UserID != target.UserID {
		query = query.Where("is_private = ?", false)
	}
	query.
		Order("postID DESC").Offset(offset).Limit(limit).Find(&posts)
	if len(posts) == 0 {
		c.JSON(http.StatusOK, []PostResponse{})
		return
	}
	postIDs := make([]int, 0, len(posts))
	userIDs := make(map[int]struct{})
	for _, p := range posts {
		postIDs = append(postIDs, p.PostID)
		if p.UserID != 0 {
			userIDs[p.UserID] = struct{}{}
		}
	}
	savedSet := browseBatchSavedSet(db, viewer.UserID, postIDs)
	likedSet := browseBatchLikedSet(db, viewer.UserID, postIDs)
	userMap := browseBatchUserMap(db, userIDs)
	postResponses := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		u := browseGetUserForPost(userMap, post.UserID)
		_, isSaved := savedSet[post.PostID]
		_, isLiked := likedSet[post.PostID]
		postResponses = append(postResponses, buildPostResponse(post, u, isSaved, isLiked))
	}
	c.JSON(http.StatusOK, postResponses)
}

type UserProfilePostNumMsg struct {
	UserTelephone string `json:"userTelephone"`
	TargetUserID  int    `json:"targetUserID"`
}

// UserProfilePostNum 个人主页可展示的发帖总数（不含「打分」）；隐藏时对非本人返回 0。
func UserProfilePostNum(c *gin.Context) {
	db := freshDB()
	var req UserProfilePostNumMsg
	_ = c.Bind(&req)
	count := 0
	if req.TargetUserID <= 0 {
		c.JSON(http.StatusOK, gin.H{"Postcount": count})
		return
	}
	viewer, ok := viewerFromAuth(c)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"Postcount": count})
		return
	}
	var target model.User
	if err := db.Where("userID = ?", req.TargetUserID).First(&target).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"Postcount": count})
		return
	}
	if target.HideProfilePosts && viewer.UserID != target.UserID {
		c.JSON(http.StatusOK, gin.H{"Postcount": 0})
		return
	}
	query := db.Model(&model.Post{}).Where("userID = ?", target.UserID)
	if viewer.UserID != target.UserID {
		query = query.Where("is_private = ?", false)
	}
	query.Count(&count)
	c.JSON(http.StatusOK, gin.H{"Postcount": count})
}

func browseFetchPosts(db *gorm.DB, req *BrowseMeg, temUser *model.User) []model.Post {
	limit, offset := req.Limit, req.Offset
	if limit <= 0 {
		limit = 15
	}
	if offset < 0 {
		offset = 0
	}
	partition := req.Partition
	searchinfo := req.Searchinfo
	tag := req.Tag
	searchsort := req.Searchsort
	if searchsort == "" {
		searchsort = "home"
	}

	if searchsort == "save" {
		var saves []model.Psave
		db.Model(&model.Psave{}).
			Where("userID = ?", temUser.UserID).
			Order("psaveID DESC").
			Offset(offset).Limit(limit).
			Find(&saves)
		if len(saves) == 0 {
			return nil
		}
		postIDs := make([]int, 0, len(saves))
		for _, s := range saves {
			postIDs = append(postIDs, s.PtargetID)
		}
		var posts []model.Post
		applyPostVisibilityScope(db.Model(&model.Post{}).Where("postID IN (?)", postIDs), temUser.UserID).Find(&posts)
		postMap := make(map[int]model.Post, len(posts))
		for _, p := range posts {
			postMap[p.PostID] = p
		}
		ordered := make([]model.Post, 0, len(saves))
		for _, s := range saves {
			if p, ok := postMap[s.PtargetID]; ok && p.PostID != 0 {
				ordered = append(ordered, p)
			}
		}
		return ordered
	}

	var posts []model.Post
	base := db.Model(&model.Post{})
	base = browseApplyTimeRange(base, req.TimeFrom, req.TimeTo)
	base = applyPostVisibilityScope(base, temUser.UserID)
	if searchsort == "home" {
		base = browseApplyOrder(base, req.OrderBy)
	} else {
		base = base.Order("postID DESC")
	}
	if searchsort == "history" {
		base = base.Where("userID = ?", temUser.UserID)
		base.Offset(offset).Limit(limit).Find(&posts)
		return posts
	}
	if partition == "主页" || len(partition) == 0 {
		if len(searchinfo) > 0 {
			if tag == "" {
				base = base.Where("(title LIKE ? OR ptext LIKE ?)", "%"+searchinfo+"%", "%"+searchinfo+"%")
			} else {
				base = base.Where("(title LIKE ? OR ptext LIKE ?) AND tag = ?", "%"+searchinfo+"%", "%"+searchinfo+"%", tag)
			}
		}
	} else if partition == "优质贴" {
		base = base.Where("is_high_quality = ?", 1)
	} else {
		base = base.Where("`partition` = ?", partition)
		if len(searchinfo) > 0 {
			if tag == "" {
				base = base.Where("(title LIKE ? OR ptext LIKE ?)", "%"+searchinfo+"%", "%"+searchinfo+"%")
			} else {
				base = base.Where("(title LIKE ? OR ptext LIKE ?) AND tag = ?", "%"+searchinfo+"%", "%"+searchinfo+"%", tag)
			}
		}
	}
	base.Offset(offset).Limit(limit).Find(&posts)
	return posts
}

func browseBatchSavedSet(db *gorm.DB, userID int, postIDs []int) map[int]struct{} {
	if len(postIDs) == 0 {
		return nil
	}
	var saves []model.Psave
	db.Where("userID = ? AND ptargetID IN (?)", userID, postIDs).Find(&saves)
	set := make(map[int]struct{}, len(saves))
	for _, s := range saves {
		set[s.PtargetID] = struct{}{}
	}
	return set
}

func browseBatchLikedSet(db *gorm.DB, userID int, postIDs []int) map[int]struct{} {
	if len(postIDs) == 0 {
		return nil
	}
	var likes []model.Plike
	db.Where("userID = ? AND ptargetID IN (?)", userID, postIDs).Find(&likes)
	set := make(map[int]struct{}, len(likes))
	for _, l := range likes {
		set[l.PtargetID] = struct{}{}
	}
	return set
}

func browseBatchUserMap(db *gorm.DB, userIDs map[int]struct{}) map[int]model.User {
	if len(userIDs) == 0 {
		return nil
	}
	ids := make([]int, 0, len(userIDs))
	for id := range userIDs {
		ids = append(ids, id)
	}
	var users []model.User
	db.Where("userID IN (?)", ids).Find(&users)
	m := make(map[int]model.User, len(users))
	for _, u := range users {
		m[u.UserID] = u
	}
	return m
}

// browseGetUserForPost 从 userMap 取发帖用户信息；userID==0 时返回“管理员”占位。
func browseGetUserForPost(userMap map[int]model.User, userID int) model.User {
	if userID == 0 {
		return model.User{Name: "管理员", Phone: "11111111111"}
	}
	if u, ok := userMap[userID]; ok {
		return u
	}
	return model.User{}
}

func GetPostNum(c *gin.Context) {
	viewer, ok := viewerFromAuth(c)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"Postcount": 0})
		return
	}
	db := freshDB()
	var requestBrowseMeg BrowseMeg
	if err := c.ShouldBindJSON(&requestBrowseMeg); err != nil {
		_ = c.ShouldBind(&requestBrowseMeg)
	}
	partition := requestBrowseMeg.Partition
	searchinfo := requestBrowseMeg.Searchinfo
	searchsort := requestBrowseMeg.Searchsort
	if searchsort == "" {
		searchsort = "home"
	}
	timeFrom := requestBrowseMeg.TimeFrom
	timeTo := requestBrowseMeg.TimeTo
	var count int
	viewerUserID := viewer.UserID
	if searchsort == "home" {
		if partition == "主页" || len(partition) == 0 {
			if len(searchinfo) == 0 && timeFrom == "" && timeTo == "" {
				applyPostVisibilityScope(db.Model(&model.Post{}), viewerUserID).Count(&count)
			} else {
				q := db.Model(&model.Post{}).Where("postId != ?", 0)
				q = applyPostVisibilityScope(q, viewerUserID)
				q = browseApplyTimeRange(q, timeFrom, timeTo)
				if len(searchinfo) > 0 {
					q = q.Where("title LIKE ? OR ptext LIKE ? OR tag LIKE ?", "%"+searchinfo+"%", "%"+searchinfo+"%", "%"+searchinfo+"%")
				}
				q.Count(&count)
			}
		} else if partition == "优质贴" {
			q := applyPostVisibilityScope(db.Model(&model.Post{}).Where("is_high_quality =?", 1), viewerUserID)
			q = browseApplyTimeRange(q, timeFrom, timeTo)
			q.Count(&count)
		} else {
			q := applyPostVisibilityScope(db.Model(&model.Post{}).Where("`partition` = ?", partition), viewerUserID)
			q = browseApplyTimeRange(q, timeFrom, timeTo)
			if len(searchinfo) > 0 {
				q = q.Where("title LIKE ? OR ptext LIKE ? OR tag LIKE ?", "%"+searchinfo+"%", "%"+searchinfo+"%", "%"+searchinfo+"%")
			}
			q.Count(&count)
		}
	} else {
		if searchsort == "save" {
			db.Model(&model.Psave{}).
				Joins("INNER JOIN posts ON psaves.ptargetID = posts.postID").
				Where("psaves.userID = ? AND posts.postID != ? AND (posts.is_private = ? OR posts.userID = ?)", viewer.UserID, 0, false, viewer.UserID).
				Count(&count)
		} else if searchsort == "history" {
			db.Model(&model.Post{}).
				Where("userID = ?", viewer.UserID).
				Where("postId != ?", 0).
				Count(&count)
		}
	}
	// 将结果返回给客户端
	c.JSON(http.StatusOK, gin.H{
		"Postcount": count,
	})
}

type PostPrivacyMsg struct {
	UserTelephone string `json:"userTelephone"`
	PostID        uint   `json:"postID"`
	IsPrivate     bool   `json:"isPrivate"`
}

func UpdatePostPrivacy(c *gin.Context) {
	viewer, ok := viewerFromAuth(c)
	if !ok {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db := freshDB()
	var req PostPrivacyMsg
	if err := c.Bind(&req); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数解析错误")
		return
	}
	var post model.Post
	db.Where("postID = ?", req.PostID).First(&post)
	if post.PostID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "帖子不存在")
		return
	}
	if post.UserID != viewer.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db.Model(&post).Update("is_private", req.IsPrivate)
	common.DelMultiLevelCache(util.GetPostDetailKey(int(req.PostID)))
	response.Success(c, gin.H{"isPrivate": req.IsPrivate}, "更新成功")
}

type SaveMsg struct {
	UserTelephone string
	PostID        uint
	//IsSaved       bool
}

func UpdateSave(c *gin.Context) {
	db := common.GetDB()
	var requestSaveMsg SaveMsg
	c.Bind(&requestSaveMsg)
	userTelephone := requestSaveMsg.UserTelephone
	postID := requestSaveMsg.PostID
	//isSaved := requestSaveMsg.IsSaved

	// Find the user by telephone
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	if user.UserID == 0 {
		return
	}
	var post model.Post
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 {
		return
	}
	if !canViewPost(post, user.UserID) {
		response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
		return
	}

	isSaved := false
	var save model.Psave
	db.Where("userID = ? AND ptargetID = ?", user.UserID, post.PostID).First(&save)
	if save.PsaveID != 0 {
		isSaved = true
	}

	if isSaved {
		db.Delete(&save)
	} else {
		newSave := model.Psave{
			UserID:    user.UserID,
			PtargetID: post.PostID,
		}
		if newSave.UserID != 0 && newSave.PtargetID != 0 {
			db.Create(&newSave)
		}
	}
}

type LikeMsg struct {
	UserTelephone string
	PostID        uint
	//IsLiked       bool
}

func UpdateLike(c *gin.Context) {
	db := common.GetDB()
	var requestLikeMsg LikeMsg
	c.Bind(&requestLikeMsg)
	userTelephone := requestLikeMsg.UserTelephone
	postID := requestLikeMsg.PostID

	// Find the user by telephone
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	if user.UserID == 0 {
		return
	}
	var post model.Post
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 {
		return
	}
	if !canViewPost(post, user.UserID) {
		response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
		return
	}

	isLiked := false
	var like model.Plike
	db.Where("userID = ? AND ptargetID = ?", user.UserID, post.PostID).First(&like)
	if like.PlikeID != 0 {
		isLiked = true
	}

	if isLiked { //已点赞
		db.Model(&post).Update("like_num", post.LikeNum-1)
		db.Delete(&like)
		service.AdjustUserReceivedLike(db, post.UserID, -1)
	} else { //未点赞
		newLike := model.Plike{
			UserID:    user.UserID,
			PtargetID: post.PostID,
			Time:      time.Now(),
		}
		if newLike.UserID != 0 && newLike.PtargetID != 0 {
			db.Model(&post).Update("like_num", post.LikeNum+1)
			// 在这里设置 点赞 的权重
			weightLike := float64(3)
			db.Model(&post).Update("heat", post.Heat+weightLike)
			db.Create(&newLike)
			service.AdjustUserReceivedLike(db, post.UserID, 1)
		}
	}

	cacheKey := util.GetPostDetailKey(int(postID))
	common.DelMultiLevelCache(cacheKey)
}

type BrowseMsg struct {
	UserTelephone string
	PostID        uint
	PostType      string
	// BrowseNum     int
}

// UpdateBrowseNum 仅记录「该用户是否已浏览过该帖」（Pbrowse）并更新 heat。
// 帖子浏览人次 browse_num 仅由 ShowDetails 中的 HotData.IncrBrowseCount（Redis）及 syncToDB 写入 MySQL，避免双写与随机注水。
func UpdateBrowseNum(c *gin.Context) {
	db := common.GetDB()
	var requestBrowseMsg BrowseMsg
	c.Bind(&requestBrowseMsg)
	userTelephone := requestBrowseMsg.UserTelephone
	postID := requestBrowseMsg.PostID
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	if user.UserID == 0 {
		return
	}
	var post model.Post
	db.Where("postID = ?", postID).First(&post)
	if post.PostID == 0 {
		return
	}
	if !canViewPost(post, user.UserID) {
		return
	}
	var browsed model.Pbrowse
	db.Where("userID = ? AND ptargetID = ?", user.UserID, postID).First(&browsed)
	if browsed.PbrowseID == 0 {
		weightBrowse := float64(1)
		db.Model(&post).Update("heat", post.Heat+weightBrowse)
		newBrowse := model.Pbrowse{
			UserID:    user.UserID,
			PtargetID: post.PostID,
			Time:      time.Now(),
		}
		db.Create(&newBrowse)
	}
}

type IDmsg struct {
	PostID   uint
	PostType string
}

func DeletePost(c *gin.Context) {
	db := common.GetDB()
	var ID IDmsg
	c.Bind(&ID)
	PostID := ID.PostID
	var post model.Post
	db.Where("postID = ?", PostID).First(&post)
	if post.PostID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "需要删除的帖子不存在")
		return
	}
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != post.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	service.OnPostDeleted(db, &post)
	db.Delete(&post)
	common.DelMultiLevelCache(util.GetPostDetailKey(int(PostID)))
	c.JSON(http.StatusOK, gin.H{"message": "帖子删除成功"})
}

type Reportmsg struct {
	TargetID      uint
	Targettype    string
	UserTelephone string
	Reason        string
}

func SubmitReport(c *gin.Context) {
	db := common.GetDB()
	var reportmsg Reportmsg
	c.Bind(&reportmsg)
	TargetID := reportmsg.TargetID
	Targettype := reportmsg.Targettype
	userTelephone := reportmsg.UserTelephone
	Reason := reportmsg.Reason
	if len(Reason) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "举报内容不能为空")
		return
	}
	var user model.User
	db.Where("phone = ?", userTelephone).First(&user)
	if user.UserID == 0 {
		return
	}
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}
	if Targettype == "post" {
		if _, ok := loadViewablePost(db, int(TargetID), user.UserID); !ok {
			response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
			return
		}
	}
	newSue := model.Sue{
		Targettype: Targettype,
		TargetID:   int(TargetID),
		UserID:     int(user.UserID),
		User:       user,
		Reason:     Reason,
		Time:       time.Now(),
		Status:     "wait",
		Finish:     false,
	}
	db.Create(&newSue)
	response.Response(c, http.StatusOK, 200, nil, "举报发送成功")
}

type PostDetailsResponse struct {
	PostID        uint
	UserID        uint
	UserName      string
	UserScore     int
	UserTelephone string
	UserAvatar    string
	UserIdentity  string
	Title         string
	Content       string
	Like          int
	Comment       int
	Browse        int
	Heat          float64
	PostTime      time.Time
	IsSaved       bool
	IsLiked       bool
	Photos        string
	Tag           string
	IsPrivate     bool
}

type PostDetailsMsg struct {
	UserTelephone string `json:"userTelephone"`
	PostID        uint   `json:"postID"`
	PostType      string `json:"postType"`
}

// showDetailStepTimer 单次 ShowDetails 请求内各阶段耗时，结束时由 defer 统一打日志
type showDetailStepTimer struct {
	start  time.Time
	last   time.Time
	steps  []showDetailTimedStep
	postID uint
}

type showDetailTimedStep struct {
	name string
	d    time.Duration
}

func newShowDetailStepTimer() *showDetailStepTimer {
	now := time.Now()
	return &showDetailStepTimer{start: now, last: now}
}

func (s *showDetailStepTimer) setPostID(id uint) {
	s.postID = id
}

func (s *showDetailStepTimer) step(name string) {
	now := time.Now()
	s.steps = append(s.steps, showDetailTimedStep{name: name, d: now.Sub(s.last)})
	s.last = now
}

func (s *showDetailStepTimer) finish() {
	total := time.Since(s.start)
	var b strings.Builder
	b.WriteString("[ShowDetails] postID=")
	b.WriteString(strconv.FormatUint(uint64(s.postID), 10))
	b.WriteString(" total=")
	b.WriteString(total.String())
	for _, st := range s.steps {
		b.WriteString(" ")
		b.WriteString(st.name)
		b.WriteString("=")
		b.WriteString(st.d.String())
	}
	// 与 Gin 默认访问日志一致写到 stdout；标准库 log 默认 stderr，在部分终端/采集里不可见
	fmt.Println(b.String())
}

// ShowDetails 重构最终版：两级缓存 + 浏览数高并发优化 + 无HEAT维护
func ShowDetails(c *gin.Context) {
	timing := newShowDetailStepTimer()
	defer timing.finish()

	db := common.GetDB()
	var requestPostDetailsMsg PostDetailsMsg

	// 绑定参数
	if err := c.ShouldBindBodyWith(&requestPostDetailsMsg, binding.JSON); err != nil {
		timing.step("bind_body")
		response.Response(c, http.StatusBadRequest, 400, nil, "参数解析失败")
		return
	}
	timing.step("bind_body")

	postID := requestPostDetailsMsg.PostID

	viewer, ok := viewerFromAuth(c)
	if !ok {
		timing.step("load_viewer_from_auth")
		response.Response(c, http.StatusNotFound, 404, nil, "无法解析当前用户")
		return
	}
	timing.step("load_viewer_from_auth")

	// 校验帖子ID
	if postID == 0 {
		timing.setPostID(0)
		timing.step("validate_post_id")
		response.Response(c, http.StatusBadRequest, 404, nil, "接收到的postID为空")
		return
	}
	timing.setPostID(postID)
	timing.step("validate_post_id")

	// ===================== 【核心：两级缓存读取冷数据】 =====================
	var post model.Post
	var user model.User
	postIDInt := int(postID)
	cacheKey := util.GetPostDetailKey(postIDInt)
	baseExpire := 30 * time.Minute

	// 调用你原有的两级缓存模板 L1 → L2 → DB
	dtoObj, err := common.MultiLevelGet(cacheKey, func() (*dto.PostDetailCacheDTO, error) {

		var dbPost model.Post
		// 查询帖子+用户
		err := db.Preload("User").Where("postID = ?", postID).First(&dbPost).Error
		if err != nil {
			return nil, err
		}

		key := util.GetPostBrowseCountKey(int(postID))
		core.MyRedis.SetNX(context.Background(), key, dbPost.BrowseNum, 0).Err()

		// 构建 DTO（解决time.Time反序列化崩溃）
		var cacheUser model.User
		if dbPost.UserID == 0 {
			cacheUser = model.User{Name: "管理员", Phone: "11111111111"}
		} else {
			cacheUser = dbPost.User
		}

		var dto dto.PostDetailCacheDTO
		dto.ConvertFromModel(dbPost, cacheUser)
		return &dto, nil
	}, baseExpire)
	timing.step("multilevel_get")

	if err != nil {
		response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
		return
	}

	post, user = dtoObj.ConvertToModel()
	timing.step("dto_to_model")

	if !canViewPost(post, viewer.UserID) {
		response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
		return
	}
	timing.step("check_visibility")

	// ===================== 【浏览数高并发优化：本地聚合 + 批量同步】 =====================
	post.BrowseNum = int(service.HotData.IncrBrowseCount(postIDInt))
	timing.step("browse_incr")

	// ===================== 以下所有逻辑 100% 完全不变 =====================
	queryIsLiked := func() bool {
		isLiked := false
		var like model.Plike
		db.Where("userID = ? AND ptargetID = ?", viewer.UserID, postID).First(&like)
		if like.PlikeID != 0 {
			isLiked = true
		}
		return isLiked
	}

	isLiked := queryIsLiked()
	timing.step("query_plike")
	isSaved := false
	var save model.Psave
	db.Where("userID = ? AND ptargetID = ?", viewer.UserID, postID).First(&save)
	if save.PsaveID != 0 {
		isSaved = true
	}
	timing.step("query_psave")

	postDetailsResponse := PostDetailsResponse{
		PostID:        uint(post.PostID),
		UserID:        uint(user.UserID),
		UserName:      user.Name,
		UserScore:     user.Score,
		UserTelephone: user.Phone,
		UserAvatar:    user.AvatarURL,
		UserIdentity:  user.Identity,
		Title:         post.Title,
		Content:       post.Ptext,
		Like:          post.LikeNum,
		Comment:       post.CommentNum,
		PostTime:      post.PostTime,
		IsSaved:       isSaved,
		IsLiked:       isLiked,
		Browse:        post.BrowseNum,
		Heat:          post.Heat,
		Photos:        post.Photos,
		Tag:           post.Tag,
		IsPrivate:     post.IsPrivate,
	}
	c.JSON(http.StatusOK, postDetailsResponse)
	timing.step("write_json")
}

func UploadPhotos(c *gin.Context) {
	// UserID := c.PostForm("UserID")
	// 获取前端传过来的图片
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}
	// 文件保存路径和文件名可以根据实际情况修改
	// 文件名我们采用了当前时间戳和原始文件名的组合，以避免文件名冲突
	// 时间戳采用 nanoseconds 级别，可以几乎确保每个文件名都是唯一的
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := "public/uploads/" + filename
	// 保存文件到本地
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	_, err = api.UploadImage("/src/images/uploads/"+filename, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "原图保存失败" + err.Error()})
		return
	}

	// 生成略缩图
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件打开失败"})
		return
	}
	defer src.Close()
	// 读取上传文件的内容，并解码为 image.Image 对象
	img, _, err := image.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件解码失败"})
		return
	}

	// 对图像进行缩略图处理
	resizedImage := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)
	// 这里似乎只支持绝对路径，这里要根据服务器修改
	resizedPath := "public/resized/" + filename
	err = imaging.Save(resizedImage, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 更新 Post 的 Photos 字段
	// 我们存储的是可以通过 HTTP 访问的 URL，而不是服务器本地的文件路径
	fileURL := ""
	fileURL, err = api.UploadImage("/src/images/resized/"+filename, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"fileURL": fileURL, "message": "上传成功"})
}

func UploadZip(c *gin.Context) {
	const maxUploadSize = 10 << 20 // 10 MB

	// 获取前端传过来的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件上传失败"})
		return
	}

	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件太大，不能超过10MB"})
		return
	}

	fileBytes, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}
	defer fileBytes.Close()

	buffer := make([]byte, 512)
	_, err = fileBytes.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取文件"})
		return
	}

	if http.DetectContentType(buffer) != "application/zip" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件必须是zip格式"})
		return
	}

	// 文件保存路径和文件名可以根据实际情况修改
	// 文件名我们采用了当前时间戳和原始文件名的组合，以避免文件名冲突
	// 时间戳采用 nanoseconds 级别，可以几乎确保每个文件名都是唯一的
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := "public/uploads/" + filename

	// 保存文件到本地
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件保存失败"})
		return
	}

	// 更新 Post 的 Photos 字段
	// 我们存储的是可以通过 HTTP 访问的 URL，而不是服务器本地的文件路径
	fileURL := "https://ssemarket.cn/api/uploads/" + filename

	// 返回成功
	c.JSON(http.StatusOK, gin.H{"zipURL": fileURL, "message": "上传成功"})
}

func SubmitFeedback(c *gin.Context) {
	db := common.GetDB()

	// Create a struct to hold the incoming JSON body
	var feedbackInput struct {
		Ftext      string `json:"ftext"`
		Attachment string `json:"attachment"`
	}

	// Bind the incoming JSON to the struct
	if err := c.Bind(&feedbackInput); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Invalid request body")
		return
	}

	// Create a new feedback entry
	feedback := model.Feedback{
		Ftext:      feedbackInput.Ftext,
		Attachment: feedbackInput.Attachment,
		Time:       time.Now(),
		Status:     "wait",
	}

	db.Create(&feedback)

	if db.NewRecord(feedback) {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Failed to submit feedback")
		return
	}

	// Convert to JSON and respond
	response.Success(c, gin.H{
		"feedbackID": feedback.FeedbackID,
		"ftext":      feedback.Ftext,
		"attachment": feedback.Attachment,
	}, "Feedback submitted successfully")
}

func GetAllFeedback(c *gin.Context) {
	db := common.GetDB()

	var feedbacks []model.Feedback
	db.Find(&feedbacks)

	if len(feedbacks) == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "No feedback found")
		return
	}
	// 对feedbacks按照时间进行排序
	sort.Slice(feedbacks, func(i, j int) bool {
		return feedbacks[j].Time.Before(feedbacks[i].Time)
	})
	response.Success(c, gin.H{"feedbacks": feedbacks}, "Feedback retrieved successfully")
}

//func GetFeedback(c *gin.Context) {
//	db := common.GetDB()
//
//	feedbackID, err := strconv.Atoi(c.PostForm("feedbackID"))
//	if err != nil {
//		response.Response(c, http.StatusBadRequest, 400, nil, "Invalid feedback ID")
//		return
//	}
//
//	var feedback model.Feedback
//	if err := db.First(&feedback, feedbackID).Error; err != nil {
//		if gorm.IsRecordNotFoundError(err) {
//			response.Response(c, http.StatusNotFound, 404, nil, "Feedback not found")
//		} else {
//			response.Response(c, http.StatusInternalServerError, 500, nil, "Database error")
//		}
//		return
//	}
//
//	response.Success(c, gin.H{"feedback": feedback}, "Feedback retrieved successfully")
//}

type HeatResponse struct {
	PostID uint
	Title  string
	Heat   float64
}

// 表示存在zset的val
type PostInfo struct {
	PostID    uint    `json:"postID"`
	Title     string  `json:"title"`
	BrowseNum float64 `json:"browseNum"`
}

func CalculateHeat(c *gin.Context) {
	db := common.GetDB()
	// Redis ZSet的key
	redisKey := "hot_posts_zset"
	ctx := context.Background()

	// 从Redis中获取热度前10的帖子
	topPosts, err := core.MyRedis.ZRevRangeWithScores(ctx, redisKey, 0, 19).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热度前10帖子失败"})
		return
	}

	var heatResponsesTop10 []HeatResponse

	// 处理 Redis 返回的数据（跳过私密帖，最多返回 10 条公开帖）
	for _, post := range topPosts {
		if len(heatResponsesTop10) >= 10 {
			break
		}
		var postInfo PostInfo
		err := json.Unmarshal([]byte(post.Member.(string)), &postInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析帖子信息失败"})
			return
		}

		var dbPost model.Post
		if err := db.Where("postID = ? AND is_private = ?", postInfo.PostID, false).First(&dbPost).Error; err != nil {
			continue
		}

		heatResponse := HeatResponse{
			PostID: postInfo.PostID,
			Title:  dbPost.Title,
			Heat:   postInfo.BrowseNum,
		}
		heatResponsesTop10 = append(heatResponsesTop10, heatResponse)
	}

	c.JSON(http.StatusOK, heatResponsesTop10)
}

func GetTags(c *gin.Context) {
	//从中间件存入的user中获取userID
	value, exists := c.Get("user")
	if !exists {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	var user model.User
	user = value.(model.User)
	if user.UserID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户错误")
		return
	}
	db := common.GetDB()
	// 读取tag表
	var tags []model.Tag
	db.Find(&tags)
	response.Success(c, gin.H{"tags": tags}, "获取tag列表成功！")
}
func GetTitle(c *gin.Context) {
	// 从中间件存入的user中获取userID
	_, exists := c.Get("user")
	if !exists {
		response.Response(c, http.StatusBadRequest, 400, nil, "检测不到用户，无法访问")
		return
	}
	text := c.PostForm("text")
	if text == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "请输入正文内容")
		return
	}
	// 构造POST请求的body
	// 使用bytes.Buffer创建一个新的buffer
	var requestBody bytes.Buffer

	// 创建multipart.Writer实例，将文件流写入buffer
	writer := multipart.NewWriter(&requestBody)

	// 写入text字段
	writer.WriteField("text", text)

	// 关闭multipart.Writer以设置边界
	writer.Close()

	// 创建一个通道用于接收异步操作结果
	resultChan := make(chan []byte, 1)

	// 启动goroutine执行异步HTTP请求
	go func() {
		// 向ssemarket.cn:5000/title转发POST请求
		resp, err := http.Post("http://ssemarket.cn:5000/title", writer.FormDataContentType(), &requestBody)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "请求失败")
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			response.Response(c, http.StatusInternalServerError, 500, nil, "请求titlegpt失败")
			return
		}
		// 读取response
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "读取响应失败")
			return
		}

		// 将结果发送到通道
		resultChan <- respBody
	}()

	// 设置超时时间，等待结果
	select {
	case result := <-resultChan:
		var r map[string]map[string]string
		err := json.Unmarshal(result, &r)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "解析响应失败")
			return
		}
		c.JSON(http.StatusOK, gin.H{"title": r["data"]["title"]})
	case <-time.After(60 * time.Second):
		// 如果超时，返回超时信息
		response.Response(c, http.StatusGatewayTimeout, 504, nil, "请求超时")
	}

}
