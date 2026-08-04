package controller

import (
	"net/http"
	"sort"
	"time"

	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"loginTest/service"

	"github.com/gin-gonic/gin"
)

const userRankPenaltyID = 278
const userRankPenalty = 100

func effectiveUserRankScore(u model.User) int {
	score := u.ReceivedLikeNum + u.ReceivedVoteNum
	if u.UserID == userRankPenaltyID {
		score -= userRankPenalty
	}
	return score
}

type VoteSearchMsg struct {
	Keyword       string `json:"keyword"`
	UserTelephone string `json:"userTelephone"`
}

type VoteSearchResult struct {
	PostID    int    `json:"postID"`
	Title     string `json:"title"`
	Partition string `json:"partition"`
	PostTime  string `json:"postTime"`
	VoteCount int    `json:"voteCount"`
}

func SearchPostsByTitle(c *gin.Context) {
	db := common.GetDB()
	var req VoteSearchMsg
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	user, ok := viewerFromAuth(c)
	if !ok {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}

	keyword := req.Keyword
	if len(keyword) == 0 {
		response.Success(c, gin.H{"posts": []VoteSearchResult{}}, "ok")
		return
	}

	var posts []model.Post
	db.Where("title LIKE ? AND (is_private = ? OR userID = ?)", "%"+keyword+"%", false, user.UserID).
		Order("postID DESC").
		Limit(20).
		Find(&posts)

	results := make([]VoteSearchResult, 0, len(posts))
	for _, post := range posts {
		results = append(results, VoteSearchResult{
			PostID:    post.PostID,
			Title:     post.Title,
			Partition: post.Partition,
			PostTime:  post.PostTime.Format("2006-01-02 15:04:05"),
			VoteCount: post.VoteNum,
		})
	}

	response.Success(c, gin.H{"posts": results}, "ok")
}

type SubmitVoteMsg struct {
	PostID        int    `json:"postID"`
	UserTelephone string `json:"userTelephone"`
}

func SubmitVote(c *gin.Context) {
	db := common.GetDB()
	var req SubmitVoteMsg
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	user, ok := viewerFromAuth(c)
	if !ok {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}

	var post model.Post
	db.Where("postID = ?", req.PostID).First(&post)
	if post.PostID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "帖子不存在")
		return
	}
	if !canViewPost(post, user.UserID) {
		response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
		return
	}

	var existing model.Pvote
	db.Where("userID = ? AND postID = ?", user.UserID, post.PostID).First(&existing)
	if existing.PvoteID != 0 {
		response.Response(c, http.StatusOK, 400, nil, "您已为该帖投过票")
		return
	}

	newVote := model.Pvote{
		PostID: post.PostID,
		UserID: user.UserID,
		Time:   time.Now(),
	}
	if err := db.Create(&newVote).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "投票失败")
		return
	}

	service.AdjustPostVoteNum(db, post.PostID, 1)
	service.AdjustUserReceivedVote(db, post.UserID, 1)

	response.Success(c, nil, "投票成功")
}

type VoteRankItem struct {
	PostID    int    `json:"postID"`
	Title     string `json:"title"`
	Partition string `json:"partition"`
	VoteCount int    `json:"voteCount"`
	PostTime  string `json:"postTime"`
}

// GetPostRanking 优质榜（帖子榜）：按帖子得票数（不含任何私密帖）
func GetPostRanking(c *gin.Context) {
	db := freshDB()

	var posts []model.Post
	err := db.Where("vote_num > 0 AND (is_private = ? OR is_private IS NULL)", false).
		Order("vote_num DESC, postID DESC").
		Limit(120).
		Find(&posts).Error
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取优质榜失败")
		return
	}

	items := make([]VoteRankItem, 0, len(posts))
	for _, post := range posts {
		if post.IsPrivate {
			continue
		}
		if len(items) >= 100 {
			break
		}
		items = append(items, VoteRankItem{
			PostID:    post.PostID,
			Title:     post.Title,
			Partition: post.Partition,
			VoteCount: post.VoteNum,
			PostTime:  post.PostTime.Format("2006-01-02 15:04:05"),
		})
	}

	response.Success(c, gin.H{"ranking": items}, "ok")
}

type UserRankItem struct {
	UserID          int    `json:"userID"`
	UserName        string `json:"userName"`
	AvatarURL       string `json:"avatarURL"`
	ReceivedLikeNum int    `json:"receivedLikeNum"`
	ReceivedVoteNum int    `json:"receivedVoteNum"`
	RankScore       int    `json:"rankScore"`
}

// GetUserRanking 用户榜：其全部帖子获赞数 + 获投票数（userID=278 得分减 100）
func GetUserRanking(c *gin.Context) {
	db := common.GetDB()

	var users []model.User
	err := db.Where("received_like_num > 0 OR received_vote_num > 0 OR userID = ?", userRankPenaltyID).
		Find(&users).Error
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取用户榜失败")
		return
	}

	type rankedUser struct {
		user  model.User
		score int
	}
	ranked := make([]rankedUser, 0, len(users))
	for _, u := range users {
		score := effectiveUserRankScore(u)
		if score <= 0 && u.UserID != userRankPenaltyID {
			continue
		}
		ranked = append(ranked, rankedUser{user: u, score: score})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].score != ranked[j].score {
			return ranked[i].score > ranked[j].score
		}
		return ranked[i].user.UserID < ranked[j].user.UserID
	})
	if len(ranked) > 100 {
		ranked = ranked[:100]
	}

	items := make([]UserRankItem, 0, len(ranked))
	for _, row := range ranked {
		u := row.user
		items = append(items, UserRankItem{
			UserID:          u.UserID,
			UserName:        u.Name,
			AvatarURL:       u.AvatarURL,
			ReceivedLikeNum: u.ReceivedLikeNum,
			ReceivedVoteNum: u.ReceivedVoteNum,
			RankScore:       row.score,
		})
	}

	response.Success(c, gin.H{"ranking": items}, "ok")
}

// GetVoteRanking 兼容旧接口，等同优质榜
func GetVoteRanking(c *gin.Context) {
	GetPostRanking(c)
}

type CheckVotedMsg struct {
	PostID        int    `json:"postID"`
	UserTelephone string `json:"userTelephone"`
}

func CheckUserVoted(c *gin.Context) {
	db := common.GetDB()
	var req CheckVotedMsg
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	user, ok := viewerFromAuth(c)
	if !ok {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}

	var post model.Post
	db.Where("postID = ?", req.PostID).First(&post)
	if post.PostID == 0 || !canViewPost(post, user.UserID) {
		response.Response(c, http.StatusNotFound, 404, nil, "帖子不存在")
		return
	}

	var existing model.Pvote
	db.Where("userID = ? AND postID = ?", user.UserID, req.PostID).First(&existing)

	response.Success(c, gin.H{"voted": existing.PvoteID != 0}, "ok")
}
