package route

import (
	"loginTest/api"
	"loginTest/controller"
	"loginTest/middleware"

	"github.com/gin-gonic/gin"
)

// 创建路由

func CollectRoute(r *gin.Engine) *gin.Engine {

	r.Use(middleware.CORSMiddleware())
	// 这里把路由分成了三组，其中auth是需要token验证的，也就是需要用户登录
	// noauth是不需要token的，也就是不需要用户登录
	// oauth是OAuth带限制权限认证的，需要用户先授权
	auth := r.Group("")
	noauth := r.Group("")
	oauth := r.Group("")

	auth.Use(middleware.AuthMiddleware())
	auth.Use(middleware.BloomFilterMiddleware()) // 使用布隆过滤器防止恶意攻击
	oauth.Use(middleware.OAuthMiddleware())

	//OAuth2认证
	auth.POST("/api/auth/connect/check", controller.CheckOAuth2App)
	auth.POST("/api/auth/connect/authorize", controller.AuthorizeOAuth2App)
	noauth.POST("/api/auth/connect/token", controller.GetTokenbyOAuth2)
	oauth.POST("/api/auth/connect/userinfo", controller.GetInfobyOAuth2)
	//OAuth2授权记录管理
	auth.GET("/api/auth/connect/records", controller.GetOAuth2AuthRecords)

	// 打分专区（独立数据与接口）
	auth.POST("/api/auth/rating/category/list", controller.RatingCategoryList)
	auth.POST("/api/auth/rating/category/create", controller.RatingCategoryCreate)
	auth.POST("/api/auth/rating/post/list", controller.RatingPostList)
	auth.POST("/api/auth/rating/post/count", controller.RatingPostCount)
	auth.POST("/api/auth/rating/post/create", controller.RatingPostCreate)
	auth.POST("/api/auth/rating/post/detail", controller.RatingPostDetail)
	auth.POST("/api/auth/rating/post/delete", controller.RatingPostDelete)
	auth.POST("/api/auth/rating/post/like", controller.RatingPostLike)
	auth.POST("/api/auth/rating/star/submit", controller.RatingStarSubmit)
	auth.POST("/api/auth/rating/comment/list", controller.RatingCommentList)
	auth.POST("/api/auth/rating/comment/create", controller.RatingCommentCreate)
	auth.POST("/api/auth/rating/comment/like", controller.RatingCommentLike)

	noauth.POST("/api/auth/register", controller.Register)
	noauth.POST("/api/auth/login", controller.Login)
	noauth.POST("/api/auth/refresh", controller.RefreshToken)
	auth.POST("/api/auth/apiTest", api.ApiTest)
	auth.POST("/api/auth/post", controller.Post)
	auth.POST("/api/auth/browse", controller.Browse)
	auth.POST("/api/auth/userProfilePosts", controller.UserProfilePosts)
	auth.POST("/api/auth/userProfilePostNum", controller.UserProfilePostNum)
	auth.POST("/api/auth/getPostNum", controller.GetPostNum)
	auth.POST("/api/auth/deleteMe", controller.DeleteMe)
	auth.POST("/api/auth/updateLike", controller.UpdateLike)
	auth.POST("/api/auth/showDetails", controller.ShowDetails)
	auth.POST("api/auth/showPcomments", controller.GetComments)
	auth.POST("api/auth/postPcomment", controller.PostPcomment)
	auth.POST("api/auth/postCcomment", controller.PostCcomment)
	auth.POST("api/auth/updateCcommentLike", controller.UpdateCcommentLike)
	auth.POST("api/auth/updatePcommentLike", controller.UpdatePcommentLike)
	auth.POST("api/auth/updateCcommentDeny", controller.UpdateCcommentDeny)
	auth.POST("api/auth/updatePcommentDeny", controller.UpdatePcommentDeny)
	auth.POST("/api/auth/updateSave", controller.UpdateSave)
	auth.POST("/api/auth/updateBrowseNum", controller.UpdateBrowseNum)
	auth.POST("/api/auth/updatePostPrivacy", controller.UpdatePostPrivacy)
	auth.POST("/api/auth/deletePost", controller.DeletePost)
	auth.POST("/api/auth/deletePcomment", controller.DeletePcomment)
	auth.POST("/api/auth/deleteCcomment", controller.DeleteCcomment)
	auth.POST("/api/auth/submitReport", controller.SubmitReport)
	auth.GET("/api/auth/info", controller.Info)
	noauth.POST("/api/auth/uploadPhotos", controller.UploadPhotos)
	noauth.POST("/api/auth/uploadAvatar", controller.UploadAvatar)
	noauth.POST("/api/auth/updateAvatar", controller.UpdateAvatar)
	auth.POST("/api/auth/getAvatar", controller.GetAvatar)
	auth.POST("/api/auth/getInfo", controller.GetInfo)
	auth.POST("/api/auth/updateUserInfo", controller.UpdateUserInfo)
	noauth.POST("/api/auth/uploadZip", controller.UploadZip)
	auth.POST("/api/auth/submitFeedback", controller.SubmitFeedback)
	auth.GET("/api/auth/getAllFeedback", controller.GetAllFeedback)
	auth.GET("/api/auth/calculateHeat", controller.CalculateHeat)
	auth.POST("/api/auth/searchPostsByTitle", controller.SearchPostsByTitle)
	auth.POST("/api/auth/submitVote", controller.SubmitVote)
	auth.GET("/api/auth/voteRanking", controller.GetVoteRanking)
	auth.GET("/api/auth/postRanking", controller.GetPostRanking)
	auth.GET("/api/auth/userRanking", controller.GetUserRanking)
	auth.POST("/api/auth/checkUserVoted", controller.CheckUserVoted)
	auth.GET("/api/auth/getNotice", controller.GetNotice)
	auth.GET("/api/auth/getNoticeNum", controller.GetNoticeNum)
	auth.PATCH("api/auth/readNotice/:noticeID", controller.ReadNotice)
	noauth.POST("/api/auth/modifyPassword", controller.ModifyPassword)
	noauth.POST("/api/auth/validateEmail", controller.ValidateEmail)
	auth.POST("/api/auth/identityValidate", controller.IdentityValidate)
	auth.POST("/api/auth/getFileList", controller.GetFileList)
	auth.POST("/api/auth/fileDelete", controller.DeleteFile)
	auth.POST("/api/auth/folderDelete", controller.DeleteFolder)
	auth.POST("/api/auth/uploadFile", controller.TeachUpload)
	auth.POST("/api/auth/addFolder", controller.AddFolder)
	auth.POST("/api/auth/getObjectUrl", controller.GetObjectUrl)
	auth.POST("/api/auth/changeEmailPush", controller.ChangeEmailPushStatus)
	auth.GET("/api/auth/getTags", controller.GetTags)
	auth.POST("/api/auth/gettitle", controller.GetTitle)
	auth.POST("api/auth/statistics", controller.GetStatistics)

	//获得所有聊天用户列表，并且建立ws链接
	auth.GET("/websocket/auth/chat", controller.ChatHandler)
	//点击与某人聊天主页的接口，需要监听用户是否在聊天页面
	//auth.GET("/api/auth/intoChatView", controller.IntoChatView)
	//退出与某人聊天主页的接口，需要监听用户是否在聊天页面
	//auth.POST("/api/auth/leaveChatView", controller.LeaveChatView)
	auth.GET("/api/auth/getChatHistory", controller.GetChatHistory)
	auth.GET("/api/auth/getChatNotice", controller.GetChatNotice)
	//商城板块
	auth.POST("/api/auth/postProduct", controller.PostProduct)
	auth.POST("/api/auth/getProducts", controller.GetProducts)
	auth.POST("/api/auth/getProductDetail", controller.GetProductDetail)
	auth.POST("/api/auth/getProductNum", controller.GetProductNum)
	auth.POST("/api/auth/deleteProduct", controller.DeleteProduct)
	auth.POST("/api/auth/saleProduct", controller.SaleProduct)
	auth.POST("/api/auth/getCarouselImg", controller.GetCarouselImg)

	// 简历专区
	auth.POST("/api/auth/resume/upload", controller.ResumeUpload)
	auth.POST("/api/auth/resume/list", controller.ResumeList)
	auth.POST("/api/auth/resume/detail", controller.ResumeDetail)
	auth.POST("/api/auth/resume/download", controller.ResumeDownload)
	auth.POST("/api/auth/resume/preview", controller.ResumePreview)
	auth.GET("/api/auth/resume/preview/stream/:resumeID", controller.ResumePreviewStream)
	auth.POST("/api/auth/resume/update", controller.ResumeUpdate)
	auth.POST("/api/auth/resume/delete", controller.ResumeDelete)

	// 给管理员设置一个新的路由分组
	adminAuth := r.Group("")
	adminAuth.Use(middleware.AuthMiddleware_admin())
	//adminAuth.POST("/api/auth/passUsers", controller.PassUsers)
	adminAuth.POST("/api/auth/addAdmin", controller.AddAdmin)
	adminAuth.POST("/api/auth/changePassword", controller.ChangeAdminPassword)
	adminAuth.POST("/api/auth/deleteUser", controller.DeleteUser)
	adminAuth.POST("/api/auth/deleteAdmin", controller.DeleteAdmin)
	adminAuth.POST("/api/auth/showUsers", controller.ShowFilterUsers)
	r.POST("/api/auth/adminLogin", controller.AdminLogin)
	adminAuth.GET("/api/auth/admininfo", controller.AdminInfo)
	adminAuth.GET("/api/auth/getSues", controller.GetSues)
	adminAuth.POST("/api/auth/noViolation", controller.NoViolation)
	adminAuth.POST("/api/auth/violation", controller.Violation)
	adminAuth.POST("/api/auth/adminPost", controller.AdminPost)

	return r
}
