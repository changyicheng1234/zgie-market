package middleware

import (
	"loginTest/common"
	"loginTest/controller"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// BloomFilterMiddleware 布隆过滤器中间件
func BloomFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// 1. 帖子详情接口
		if path == "/api/auth/showDetails" && method == http.MethodPost {
			var req controller.PostDetailsMsg
			// 官方原生缓存Body，最简单！
			if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil || req.PostID == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败"})
				c.Abort()
				return
			}
			if !common.CheckPostExist(int(req.PostID)) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "帖子不存在"})
				c.Abort()
				return
			}
		}

		// 2. 评论接口
		if path == "/api/auth/showPcomments" && method == http.MethodPost {
			var req controller.Commentsmsg
			// 官方原生缓存Body
			if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil ||
				req.UserTelephone == "" || req.PostID == 0 || req.TypePost == "" {
				c.JSON(http.StatusBadRequest, gin.H{"msg": "服务器无法成功解析请求"})
				c.Abort()
				return
			}
			if !common.CheckPostExist(req.PostID) {
				c.JSON(http.StatusNotFound, gin.H{"msg": "帖子不存在，无评论可查"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
