package controller

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"loginTest/api"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

const (
	resumeMaxUploadSize = 10 << 20 // 10 MB
	resumeMaxPerUser    = 5
	resumeCOSPrefix     = "resume/"
	resumeURLExpire     = 30 * time.Minute
)

var resumeAllowedExt = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

func resumeCanPreview(ext string) bool {
	switch strings.ToLower(ext) {
	case ".pdf", ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

func resumeContentType(ext, mimeType string) string {
	if mimeType != "" && mimeType != "application/octet-stream" {
		return mimeType
	}
	switch strings.ToLower(ext) {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

type ResumeBrowseMsg struct {
	Searchinfo string `json:"searchinfo"`
	Searchsort string `json:"searchsort"` // home | my
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type ResumeIDMsg struct {
	ResumeID int `json:"resumeID"`
}

type ResumeUpdateMsg struct {
	ResumeID    int    `json:"resumeID"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func validateResumeTitleDescription(title, description string) (string, string, string) {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if title == "" {
		return "", "", "标题不能为空"
	}
	if utf8.RuneCountInString(title) > 50 {
		return "", "", "标题最多50个字"
	}
	if utf8.RuneCountInString(description) > 500 {
		return "", "", "简介最多500个字"
	}
	return title, description, ""
}

type ResumeResponse struct {
	ResumeID    int    `json:"resumeID"`
	UserID      int    `json:"userID"`
	Uploader    string `json:"uploader"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	MimeType    string `json:"mimeType"`
	FileExt     string `json:"fileExt"`
	DownloadNum int    `json:"downloadNum"`
	ViewNum     int    `json:"viewNum"`
	IsAnonymous bool   `json:"isAnonymous"`
	IsOwner     bool   `json:"isOwner"`
	CanPreview  bool   `json:"canPreview"`
	CreatedAt   string `json:"createdAt"`
}

func requireLoginUser(c *gin.Context) (model.User, bool) {
	userID := GetTokenUserID(c)
	if userID == 0 {
		response.Response(c, http.StatusUnauthorized, 401, nil, "请先登录")
		return model.User{}, false
	}
	db := common.GetDB()
	var user model.User
	db.Where("userID = ?", userID).First(&user)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnauthorized, 401, nil, "用户不存在")
		return model.User{}, false
	}
	return user, true
}

func viewerIsOwner(resume model.Resume, viewer model.User) bool {
	return resume.UserID == viewer.UserID
}

func shouldMaskAnonymous(resume model.Resume, viewer model.User) bool {
	return resume.IsAnonymous && !viewerIsOwner(resume, viewer)
}

// anonymousDisplayFileName 对外展示的文件名，避免真实姓名出现在 URL/界面。
func anonymousDisplayFileName(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	if ext == "" {
		return "简历文件"
	}
	return "简历" + ext
}

// downloadDisplayFileName 下载保存名：{用户名|匿名用户}-简历.ext
func downloadDisplayFileName(resume model.Resume, uploader model.User) string {
	displayName := strings.TrimSpace(uploader.Name)
	if resume.IsAnonymous {
		displayName = "匿名用户"
	}
	if displayName == "" {
		displayName = "用户"
	}
	displayName = sanitizeDownloadFileSegment(displayName)
	ext := strings.ToLower(filepath.Ext(resume.FileName))
	if ext == "" {
		ext = ".pdf"
	}
	return displayName + "-简历" + ext
}

func sanitizeDownloadFileSegment(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
			continue
		}
		b.WriteRune(r)
	}
	out := strings.TrimSpace(b.String())
	if out == "" {
		return "用户"
	}
	return out
}

func resumeToResponse(r model.Resume, uploader model.User, viewer model.User) ResumeResponse {
	mask := shouldMaskAnonymous(r, viewer)
	owner := viewerIsOwner(r, viewer)

	uploaderName := uploader.Name
	uploaderID := uploader.UserID
	// 对外一律不返回真实文件名，避免原始上传名泄露身份或敏感信息
	fileName := anonymousDisplayFileName(r.FileName)
	if mask {
		uploaderName = "匿名用户"
		uploaderID = 0
	}

	ext := strings.ToLower(filepath.Ext(r.FileName))
	fileExt := strings.TrimPrefix(ext, ".")
	canPreview := resumeCanPreview(ext)

	return ResumeResponse{
		ResumeID:    r.ResumeID,
		UserID:      uploaderID,
		Uploader:    uploaderName,
		Title:       r.Title,
		Description: r.Description,
		FileName:    fileName,
		FileSize:    r.FileSize,
		MimeType:    r.MimeType,
		FileExt:     fileExt,
		DownloadNum: r.DownloadNum,
		ViewNum:     r.ViewNum,
		IsAnonymous: r.IsAnonymous,
		IsOwner:     owner,
		CanPreview:  canPreview,
		CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04"),
	}
}

func newResumeCosKey(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// 路径不含 userID，避免 COS URL 泄露上传者身份
	return fmt.Sprintf("%s%s/%d%s", resumeCOSPrefix, hex.EncodeToString(b), time.Now().UnixNano(), ext), nil
}

func resumeAccessURL(cosKey string, inline bool, attachmentFilename string) (string, error) {
	url, err := api.GetPresignedURL(cosKey, resumeURLExpire, inline, attachmentFilename)
	if err != nil {
		return api.GetUrl(cosKey), nil
	}
	return url, nil
}

func ResumeUpload(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	db := common.GetDB()
	var count int
	db.Model(&model.Resume{}).Where("userID = ? AND is_delete = ?", user.UserID, false).Count(&count)
	if count >= resumeMaxPerUser {
		response.Response(c, http.StatusBadRequest, 400, nil, fmt.Sprintf("每人最多上传%d份简历", resumeMaxPerUser))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "请选择要上传的文件")
		return
	}
	if file.Size > resumeMaxUploadSize {
		response.Response(c, http.StatusBadRequest, 400, nil, "文件不能超过10MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !resumeAllowedExt[ext] {
		response.Response(c, http.StatusBadRequest, 400, nil, "仅支持 PDF、Word、JPG、PNG、WebP、GIF 格式")
		return
	}

	title, description, errMsg := validateResumeTitleDescription(c.PostForm("title"), c.PostForm("description"))
	if errMsg != "" {
		response.Response(c, http.StatusBadRequest, 400, nil, errMsg)
		return
	}
	isAnonymous := c.PostForm("isAnonymous") == "true" || c.PostForm("isAnonymous") == "1"

	f, err := file.Open()
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "无法读取文件")
		return
	}
	defer f.Close()

	buffer := make([]byte, 512)
	n, _ := f.Read(buffer)
	mimeType := http.DetectContentType(buffer[:n])
	if !isAllowedResumeMime(mimeType, ext) {
		response.Response(c, http.StatusBadRequest, 400, nil, "文件格式不正确")
		return
	}

	safeName := filepath.Base(file.Filename)
	safeName = strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' {
			return '_'
		}
		return r
	}, safeName)

	cosKey, err := newResumeCosKey(ext)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "生成存储路径失败")
		return
	}
	if _, err = f.Seek(0, 0); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "无法读取文件")
		return
	}
	if err = api.UploadFile(cosKey, f); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件上传失败")
		return
	}

	now := time.Now()
	resume := model.Resume{
		UserID:      user.UserID,
		Title:       title,
		Description: description,
		FileName:    safeName,
		CosKey:      cosKey,
		FileSize:    file.Size,
		MimeType:    mimeType,
		IsAnonymous: isAnonymous,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := db.Create(&resume).Error; err != nil {
		_ = api.FileDelete(cosKey)
		response.Response(c, http.StatusInternalServerError, 500, nil, "保存失败")
		return
	}

	response.Success(c, gin.H{"resumeID": resume.ResumeID}, "简历上传成功")
}

func isAllowedResumeMime(mimeType, ext string) bool {
	switch ext {
	case ".pdf":
		return strings.Contains(mimeType, "pdf") || mimeType == "application/octet-stream"
	case ".doc":
		return strings.Contains(mimeType, "msword") || mimeType == "application/octet-stream"
	case ".docx":
		return strings.Contains(mimeType, "wordprocessingml") || strings.Contains(mimeType, "zip") || mimeType == "application/octet-stream"
	case ".jpg", ".jpeg":
		return strings.Contains(mimeType, "jpeg") || mimeType == "application/octet-stream"
	case ".png":
		return strings.Contains(mimeType, "png") || mimeType == "application/octet-stream"
	case ".webp":
		return strings.Contains(mimeType, "webp") || mimeType == "application/octet-stream"
	case ".gif":
		return strings.Contains(mimeType, "gif") || mimeType == "application/octet-stream"
	default:
		return false
	}
}

func ResumeList(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	var req ResumeBrowseMsg
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	db := common.GetDB()
	query := db.Model(&model.Resume{}).Where("is_delete = ?", false)
	if req.Searchsort == "my" {
		query = query.Where("userID = ?", user.UserID)
	}
	if req.Searchinfo != "" {
		like := "%" + req.Searchinfo + "%"
		query = query.Where("(title LIKE ? OR description LIKE ? OR file_name LIKE ?)", like, like, like)
	}

	var resumes []model.Resume
	query.Order("resumeID DESC").Offset(req.Offset).Limit(req.Limit).Find(&resumes)

	userIDs := make([]int, 0, len(resumes))
	for _, r := range resumes {
		userIDs = append(userIDs, r.UserID)
	}
	userMap := make(map[int]model.User)
	if len(userIDs) > 0 {
		var users []model.User
		db.Where("userID IN (?)", userIDs).Find(&users)
		for _, u := range users {
			userMap[u.UserID] = u
		}
	}

	responses := make([]ResumeResponse, 0, len(resumes))
	for _, r := range resumes {
		u := userMap[r.UserID]
		responses = append(responses, resumeToResponse(r, u, user))
	}
	c.JSON(http.StatusOK, responses)
}

func incrementResumeViewNum(db *gorm.DB, resume *model.Resume) {
	db.Model(resume).UpdateColumn("view_num", gorm.Expr("view_num + ?", 1))
	resume.ViewNum++
}

func ResumeDetail(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	var req ResumeIDMsg
	if err := c.ShouldBindJSON(&req); err != nil || req.ResumeID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	db := common.GetDB()
	var resume model.Resume
	db.Where("resumeID = ? AND is_delete = ?", req.ResumeID, false).First(&resume)
	if resume.ResumeID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "简历不存在")
		return
	}

	incrementResumeViewNum(db, &resume)

	var uploader model.User
	db.Where("userID = ?", resume.UserID).First(&uploader)
	c.JSON(http.StatusOK, resumeToResponse(resume, uploader, user))
}

func resumeFilePayload(resume model.Resume, uploader model.User, viewer model.User, inline bool) (gin.H, error) {
	var fileName string
	var attachmentName string
	if inline {
		fileName = anonymousDisplayFileName(resume.FileName)
		attachmentName = ""
	} else {
		fileName = downloadDisplayFileName(resume, uploader)
		attachmentName = fileName
	}
	fileURL, err := resumeAccessURL(resume.CosKey, inline, attachmentName)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"fileURL":    fileURL,
		"fileName":   fileName,
		"canPreview": resumeCanPreview(filepath.Ext(resume.FileName)),
	}, nil
}

func ResumePreview(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	var req ResumeIDMsg
	if err := c.ShouldBindJSON(&req); err != nil || req.ResumeID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	db := common.GetDB()
	var resume model.Resume
	db.Where("resumeID = ? AND is_delete = ?", req.ResumeID, false).First(&resume)
	if resume.ResumeID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "简历不存在")
		return
	}

	if !resumeCanPreview(filepath.Ext(resume.FileName)) {
		response.Response(c, http.StatusBadRequest, 400, nil, "当前文件类型不支持在线预览，请下载查看")
		return
	}

	var uploader model.User
	db.Where("userID = ?", resume.UserID).First(&uploader)
	payload, err := resumeFilePayload(resume, uploader, user, true)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取预览链接失败")
		return
	}
	response.Success(c, payload, "获取预览链接成功")
}

// ResumePreviewStream 经服务端代理 PDF 流，响应 Content-Disposition: inline，供前端 blob 内嵌预览。
func ResumePreviewStream(c *gin.Context) {
	if _, ok := requireLoginUser(c); !ok {
		return
	}

	resumeID, err := strconv.Atoi(c.Param("resumeID"))
	if err != nil || resumeID <= 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	db := common.GetDB()
	var resume model.Resume
	db.Where("resumeID = ? AND is_delete = ?", resumeID, false).First(&resume)
	if resume.ResumeID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "简历不存在")
		return
	}

	ext := strings.ToLower(filepath.Ext(resume.FileName))
	if !resumeCanPreview(ext) {
		response.Response(c, http.StatusBadRequest, 400, nil, "当前文件类型不支持在线预览")
		return
	}

	body, contentType, err := api.GetObject(resume.CosKey)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "读取文件失败")
		return
	}
	defer body.Close()

	// 预览流不携带 filename，避免微信等内置浏览器误判为下载
	c.Header("Content-Type", resumeContentType(ext, contentType))
	c.Header("Content-Disposition", "inline")
	c.Header("Cache-Control", "private, no-store")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "SAMEORIGIN")
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, body)
}

func ResumeDownload(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	var req ResumeIDMsg
	if err := c.ShouldBindJSON(&req); err != nil || req.ResumeID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	db := common.GetDB()
	var resume model.Resume
	db.Where("resumeID = ? AND is_delete = ?", req.ResumeID, false).First(&resume)
	if resume.ResumeID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "简历不存在")
		return
	}

	var uploader model.User
	db.Where("userID = ?", resume.UserID).First(&uploader)

	db.Model(&resume).UpdateColumn("download_num", gorm.Expr("download_num + ?", 1))
	payload, err := resumeFilePayload(resume, uploader, user, false)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取下载链接失败")
		return
	}
	response.Success(c, gin.H{
		"downloadURL": payload["fileURL"],
		"fileName":    payload["fileName"],
	}, "获取下载链接成功")
}

func ResumeUpdate(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	var req ResumeUpdateMsg
	if err := c.ShouldBindJSON(&req); err != nil || req.ResumeID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	title, description, errMsg := validateResumeTitleDescription(req.Title, req.Description)
	if errMsg != "" {
		response.Response(c, http.StatusBadRequest, 400, nil, errMsg)
		return
	}

	db := common.GetDB()
	var resume model.Resume
	db.Where("resumeID = ? AND is_delete = ?", req.ResumeID, false).First(&resume)
	if resume.ResumeID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "简历不存在")
		return
	}
	if resume.UserID != user.UserID {
		response.Response(c, http.StatusForbidden, 403, nil, "无权编辑该简历")
		return
	}

	now := time.Now()
	if err := db.Model(&resume).Updates(map[string]interface{}{
		"title":       title,
		"description": description,
		"updated_at":  now,
	}).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "保存失败")
		return
	}

	resume.Title = title
	resume.Description = description
	resume.UpdatedAt = now

	var uploader model.User
	db.Where("userID = ?", resume.UserID).First(&uploader)
	c.JSON(http.StatusOK, resumeToResponse(resume, uploader, user))
}

func ResumeDelete(c *gin.Context) {
	user, ok := requireLoginUser(c)
	if !ok {
		return
	}

	var req ResumeIDMsg
	if err := c.ShouldBindJSON(&req); err != nil || req.ResumeID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}

	db := common.GetDB()
	var resume model.Resume
	db.Where("resumeID = ? AND is_delete = ?", req.ResumeID, false).First(&resume)
	if resume.ResumeID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "简历不存在")
		return
	}
	if resume.UserID != user.UserID {
		response.Response(c, http.StatusForbidden, 403, nil, "无权删除该简历")
		return
	}

	db.Model(&resume).Update("is_delete", true)
	_ = api.FileDelete(resume.CosKey)
	response.Success(c, nil, "简历已删除")
}
