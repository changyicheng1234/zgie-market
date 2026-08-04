package controller

import (
	"context"
	"fmt"
	"image"
	"log"
	"loginTest/api"
	"loginTest/common"
	"loginTest/core"
	"loginTest/dto"
	"loginTest/model"
	"loginTest/response"
	"loginTest/util"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const (
	VALIDATESUCCESS = 0
	EMAILERROR      = 1
	EMAILNOTEXIST   = 2
	VALIDATEFAILED  = 3
)

const MAX_VALIDATE_FAIL = 5

func GetTokenUserID(c *gin.Context) int {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || len(tokenString) <= 7 || !strings.HasPrefix(tokenString, "Bearer ") {
		return 0
	}
	tokenString = tokenString[7:]
	token, claims, err := common.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return 0
	}

	// 获取token中的用户标识符
	tokenUserID := claims.UserID
	return tokenUserID
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("phone = ?", telephone).First(&user)
	if user.UserID != 0 {
		return true
	}
	return false
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.UserID != 0 {
		return true
	}
	return false
}

//	func isNumExist(db *gorm.DB, num int) bool {
//		var user model.User
//		db.Where("num = ?", num).First(&user)
//		if user.UserID != 0 {
//			return true
//		}
//		return false
//	}
type UserMsg struct {
	UserID int `json:"user_id"`
}
type registerUser struct {
	Name      string `gorm:"type:varchar(20);not null"`
	Phone     string `gorm:"type:varchar(11);not null"`
	Password  string `gorm:"size:255;not null"`
	Password2 string `gorm:"size:255;not null"`
	Email     string `gorm:"type:varchar(11);not null"`
	//Num       string `gorm:"type:varchar(8);not null"`
	ValiCode string `gorm:"type:varchar(10);not null"`
	CDKey    string `json:"CDKey"`
}

type identity struct {
	Email    string `gorm:"type:varchar(50);not null"`
	ValiCode string `gorm:"type:varchar(10);not null"`
}

type ModifyUser struct {
	Email     string `gorm:"type:varchar(50);not null"`
	Phone     string `gorm:"type:varchar(11);not null"`
	Password  string `gorm:"size:255;not null"`
	Password2 string `gorm:"size:255;not null"`
	ValiCode  string
}

type requestEmail struct {
	Email string `gorm:"type:varchar(50);not null"`
	Mode  int
}

func DeleteMe(c *gin.Context) {
	db := common.GetDB()

	var requestUser = registerUser{}
	c.Bind(&requestUser)
	phone := requestUser.Phone
	email := requestUser.Email
	valiCode := requestUser.ValiCode

	var checkUser model.User
	db.Where("phone = ?", phone).First(&checkUser)
	userID := checkUser.UserID
	if checkUser.Email != email {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "请用注册时使用的email完成注销")
		return
	}
	if userID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "未找到该用户")
		return
	}

	resCode := IdentityValidateService(email, valiCode)
	if resCode == VALIDATEFAILED {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	} else if resCode == EMAILERROR {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	} else if resCode == EMAILNOTEXIST {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
		return
	}

	checkUser.Name = "用户已注销"
	checkUser.Phone = "0"
	checkUser.Email = ""
	checkUser.AvatarURL = ""

	db.Save(&checkUser)

	response.Response(c, http.StatusOK, 200, nil, "成功删除该用户")
}

func Register(c *gin.Context) {
	// 连接数据库
	db := common.GetDB()
	// 从前端获取信息
	var requestUser = registerUser{}
	c.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Phone
	password1 := requestUser.Password
	password2 := requestUser.Password2
	password1 = util.Decrypt(password1)
	password2 = util.Decrypt(password2)
	email := strings.TrimSpace(requestUser.Email)
	print(email)
	//num := requestUser.Num
	valiCode := requestUser.ValiCode
	cdkey := requestUser.CDKey
	//若使用postman等工具，写法如下：
	// name := c.PostForm("name")
	// telephone := c.PostForm("telephone")
	// password := c.PostForm("password")
	// PostForm中的参数与在postman中发送信息的名字要一致
	if len(telephone) == 0 {
		for {
			telephone = util.GenerateRandomDigits(11)
			// 如果存在，继续生成新的随机字符串
			// 如果不存在，退出循环
			if !isTelephoneExist(db, telephone) {
				break
			}
		}
	}
	//验证数据
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "手机号必须为11位!!!")
		println(telephone)
		return
	}
	if len(password1) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	if len(password2) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}
	if password1 != password2 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "两次密码不相同，请重新输入")
		return
	}
	//如果名称没有传，给一个10位的随机字符串
	if len(name) == 0 {
		name = gofakeit.Name()
	}
	// 判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该手机号已存在")
		return
	}
	check := false
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			check = true
		}
	}
	if !check {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	}
	if isEmailExist(db, email) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱已存在")
		return
	}
	rds := core.MyRedis
	ctx := context.Background()
	codeKey := "vcode:" + email
	failKey := "vcode_fail:" + email
	correctValiCode, err := rds.Get(ctx, codeKey).Result()
	if err != nil || correctValiCode != valiCode {
		failCountStr, _ := rds.Get(ctx, failKey).Result()
		failCount := 0
		if failCountStr != "" {
			if v, convErr := strconv.Atoi(failCountStr); convErr == nil {
				failCount = v
			}
		}
		failCount++
		rds.Set(ctx, failKey, strconv.Itoa(failCount), 5*time.Minute)
		if failCount >= MAX_VALIDATE_FAIL {
			rds.Del(ctx, codeKey)
		}
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}
	var key model.CDKey
	db.Where("content =? AND used = 0", cdkey).First(&key)
	if key.CDKeyID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "激活码错误")
		return
	}
	// 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
		return
	}
	//numInt, err := strconv.Atoi(num)
	//fmt.Println(num, numInt)
	//if err != nil {
	//	response.Response(c, http.StatusInternalServerError, 400, nil, "学号异常")
	//	return
	//}
	//if isNumExist(db, numInt) {
	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该学号已存在")
	//	return
	//}
	// 创建新用户结构体
	newUser := model.User{
		Name:     name,
		Phone:    telephone,
		Password: string(hasedPassword),
		Banend:   time.Now(),
		Email:    email,
		//Num:      numInt,
	}
	// 将结构体传进Create函数即可在数据库中添加一条记录
	// 其他的增删查改功能参见postController里的updateLike函数
	db.Create(&newUser)
	if newUser.UserID == 0 {
		response.Response(c, http.StatusInternalServerError, 400, nil, "userID为0")
		return
	}
	// 修改激活码
	// 警告：当使用 struct 更新时，GORM只会更新那些非零值的字段
	db.Model(&key).Update(model.CDKey{
		Used:     true,
		UsedTime: time.Now(),
	})
	//发放 access token 和 refresh token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "账号已注册，但无法返回token")
		log.Printf("token generate error: %v", err)
		return
	}
	refreshToken, err := common.ReleaseRefreshToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "账号已注册，但无法返回refresh token")
		log.Printf("refresh token generate error: %v", err)
		return
	}
	// 将 refresh token 记录到数据库，使其一次性使用
	refreshRecord := model.RefreshToken{
		UserID:    newUser.UserID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}
	db.Create(&refreshRecord)
	//返回结果
	response.Success(c, gin.H{"token": token, "refresh_token": refreshToken}, "注册成功")
	rds.Del(ctx, codeKey)
	rds.Del(ctx, failKey)
}

func IdentityValidate(c *gin.Context) {
	//db := common.GetDB()
	// 从前端获取信息
	var requestUser = identity{}
	c.Bind(&requestUser)
	//获取参数

	email := requestUser.Email
	valiCode := requestUser.ValiCode
	resCode := IdentityValidateService(email, valiCode)
	if resCode == VALIDATEFAILED {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
	} else if resCode == EMAILERROR {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
	} else if resCode == EMAILNOTEXIST {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
	} else {
		response.Response(c, http.StatusOK, 200, nil, "身份验证成功")
	}
}

// 验证码校验服务
func IdentityValidateService(email string, valiCode string) int {
	db := common.GetDB()
	check := false
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			check = true
		}
	}

	if !check {
		return EMAILERROR
	}
	if !isEmailExist(db, email) {
		return EMAILNOTEXIST
	}
	rds := core.MyRedis
	ctx := context.Background()
	codeKey := "vcode:" + email
	failKey := "vcode_fail:" + email
	correctValiCode, err := rds.Get(ctx, codeKey).Result()
	if err != nil || correctValiCode != valiCode {
		failCountStr, _ := rds.Get(ctx, failKey).Result()
		failCount := 0
		if failCountStr != "" {
			if v, convErr := strconv.Atoi(failCountStr); convErr == nil {
				failCount = v
			}
		}
		failCount++
		rds.Set(ctx, failKey, strconv.Itoa(failCount), 5*time.Minute)
		if failCount >= MAX_VALIDATE_FAIL {
			rds.Del(ctx, codeKey)
		}
		return VALIDATEFAILED
	}
	//返回结果
	rds.Del(ctx, codeKey)
	rds.Del(ctx, failKey)
	return VALIDATESUCCESS
}

type loginuser struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	//Phone    string `gorm:"type:varchar(11);not null"`
	Email    string `json:"email"`
	Password string `gorm:"size:255;not null"`
}

func Login(c *gin.Context) {
	// 邮箱登录
	db := common.GetDB()
	var requestUser = loginuser{}
	c.Bind(&requestUser)
	//获取参数
	//telephone := requestUser.Phone
	email := requestUser.Email
	password := requestUser.Password
	password = util.Decrypt(password)
	//数据验证
	//if len(telephone) == 0 {
	//	msg := "手机号为空！" + telephone
	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, msg)
	//	return
	//}
	//if len(telephone) != 11 {
	//	msg := "手机号必须为11位" + telephone
	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, msg)
	//	return
	//}
	if len(email) == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱号为空！"+email)
		return
	}
	//if !(strings.HasSuffix(email, "@mail.sysu.edu.cn") || strings.HasSuffix(email, "@mail2.sysu.edu.cn")) {
	//	response.Response(c, http.StatusUnprocessableEntity, 400, nil, "请输入正确的中大邮箱")
	//	return
	//}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码不能少于6位")
		return
	}

	//判断手机号是否存在
	var user model.User
	//db.Where("phone = ?", telephone).First(&user)
	db.Where("email = ?", email).First(&user)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	// if user.IDpass == false {
	// 	response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户注册尚未通过审核，请耐心等待")
	// 	return
	// }
	//发放 access token 和 refresh token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	refreshToken, err := common.ReleaseRefreshToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 400, nil, "系统异常")
		log.Printf("refresh token generate error: %v", err)
		return
	}
	// 将 refresh token 记录到数据库，使其一次性使用
	refreshRecord := model.RefreshToken{
		UserID:    user.UserID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}
	db.Create(&refreshRecord)
	//返回结果
	response.Success(c, gin.H{"token": token, "refresh_token": refreshToken}, "登录成功")
}

// 通过 refresh token 刷新 access token
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&req); err != nil || req.RefreshToken == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "请求参数错误")
		return
	}

	token, claims, err := common.ParseRefreshToken(req.RefreshToken)
	if err != nil || !token.Valid {
		response.Response(c, http.StatusUnauthorized, 401, nil, "refresh token 无效")
		return
	}

	userID := claims.UserID
	if userID == 0 {
		response.Response(c, http.StatusUnauthorized, 401, nil, "refresh token 无效")
		return
	}

	db := common.GetDB()
	var user model.User
	if err := db.Where("userID = ?", userID).First(&user).Error; err != nil || user.UserID == 0 {
		response.Response(c, http.StatusUnauthorized, 401, nil, "用户不存在")
		return
	}

	// 在数据库中检查 refresh token 是否存在且未使用，并未过期
	var refreshRecord model.RefreshToken
	if err := db.Where("token = ? AND user_id = ? AND is_used = ?", req.RefreshToken, userID, false).First(&refreshRecord).Error; err != nil || refreshRecord.ID == 0 {
		response.Response(c, http.StatusUnauthorized, 401, nil, "refresh token 无效或已使用")
		return
	}
	if time.Now().After(refreshRecord.ExpiresAt) {
		// 标记为已使用，避免重复提示
		refreshRecord.IsUsed = true
		db.Save(&refreshRecord)
		response.Response(c, http.StatusUnauthorized, 401, nil, "refresh token 已过期")
		return
	}

	// 将当前 refresh token 标记为已使用
	refreshRecord.IsUsed = true
	db.Save(&refreshRecord)

	// 重新发放新的 access token 和 refresh token
	newToken, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "生成新的token失败")
		return
	}
	newRefreshToken, err := common.ReleaseRefreshToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "生成新的refresh token失败")
		return
	}
	// 将新的 refresh token 记录到数据库
	newRefreshRecord := model.RefreshToken{
		UserID:    user.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}
	db.Create(&newRefreshRecord)

	response.Success(c, gin.H{
		"token":         newToken,
		"refresh_token": newRefreshToken,
	}, "刷新token成功")
}

func ModifyPassword(c *gin.Context) {
	db := common.GetDB()
	var user model.User
	var inputUser ModifyUser
	c.Bind(&inputUser)
	email := inputUser.Email
	password := inputUser.Password
	password = util.Decrypt(password)
	password2 := inputUser.Password2
	password2 = util.Decrypt(password2)
	valiCode := inputUser.ValiCode
	if password != password2 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "密码输入不一致")
		return
	}
	if !isEmailExist(db, email) {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "未找到邮箱")
		return
	}

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
		return
	}
	db.Where("email = ?", email).First(&user)
	//fmt.Println(phone)
	//fmt.Println(password)
	if user.UserID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "用户不存在")
		return
	}

	//校验验证码是否正确
	resCode := IdentityValidateService(email, valiCode)
	if resCode == VALIDATEFAILED {
		response.Response(c, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	} else if resCode == EMAILERROR {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "邮箱不符合要求")
		return
	} else if resCode == EMAILNOTEXIST {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "该邮箱不存在")
		return
	} else {
		user.Password = string(hasedPassword)
		db.Save(&user)
		response.Success(c, gin.H{"data": user}, "修改密码成功")
		return
	}
}

type OAuth2Info struct {
	AppId        string `json:"app_id"`
	RedirectURL  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
}

// 检查OAuth2应用信息，前端传入appid，返回对应的app_name、app_icon
func CheckOAuth2App(c *gin.Context) {
	db := common.GetDB()
	var requestInfo OAuth2Info
	if err := c.Bind(&requestInfo); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "请求参数格式错误")
		return
	}

	appId := requestInfo.AppId

	// 验证应用是否存在
	var app model.OAuth2App
	if err := db.Where("app_id = ? AND is_active = ?", appId, true).First(&app).Error; err != nil {
		response.Response(c, http.StatusNotFound, 404, nil, "OAuth2认证应用不存在或未开通")
		return
	}

	// 返回应用信息
	appInfo := gin.H{
		"app_id":      app.AppID,
		"app_name":    app.AppName,
		"app_icon":    app.AppIcon,
		"description": app.Description,
	}

	response.Response(c, http.StatusOK, 200, appInfo, "应用信息获取成功")
}

// OAuth2授权接口，传入app_id、redirect_uri、response_type、scope、state
// 如果app_id跟数据库里的redirect_url对应，则返回带着临时有效的code的redirect_uri
func AuthorizeOAuth2App(c *gin.Context) {
	db := common.GetDB()
	var requestInfo OAuth2Info
	if err := c.Bind(&requestInfo); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "请求参数格式错误")
		return
	}

	appId := requestInfo.AppId
	redirectURI := requestInfo.RedirectURL
	responseType := requestInfo.ResponseType
	scope := requestInfo.Scope
	state := requestInfo.State

	// 验证response_type，response_type目前只允许code(authorization_code)
	if responseType != "code" {
		response.Response(c, http.StatusBadRequest, 400, nil, "不支持的response_type")
		return
	}

	// 验证scope，scope目前只允许read_basic，用于授权当前用户的头像、用户名、邮箱等基本信息
	if scope != "read_basic" {
		response.Response(c, http.StatusBadRequest, 400, nil, "不支持的scope")
		return
	}

	// 验证应用是否存在
	var app model.OAuth2App
	if err := db.Where("app_id = ? AND is_active = ?", appId, true).First(&app).Error; err != nil {
		response.Response(c, http.StatusNotFound, 404, nil, "应用不存在或未激活")
		return
	}

	// 验证redirect_uri是否与app_id对应应用一致
	if app.RedirectURI != redirectURI {
		response.Response(c, http.StatusBadRequest, 400, nil, "redirect_uri不匹配")
		return
	}

	// 获取当前用户（需要先通过AuthMiddleware验证）
	userInterface, exists := c.Get("user")
	if !exists {
		response.Response(c, http.StatusUnauthorized, 401, nil, "用户未登录")
		return
	}
	user := userInterface.(model.User)

	// 生成授权码
	code := util.GenerateRandomString(32)

	// 保存授权码到数据库
	authCode := model.OAuth2Code{
		Code:        code,
		AppID:       appId,
		UserID:      user.UserID,
		Scope:       scope,
		RedirectURI: redirectURI,
		State:       state,
		ExpiresAt:   time.Now().Add(5 * time.Minute), // 5分钟有效期
		IsUsed:      false,
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&authCode).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "生成授权码失败")
		return
	}

	// 记录授权信息
	authRecord := model.OAuth2AuthRecord{
		UserID:       user.UserID,
		AppID:        appId,
		AppName:      app.AppName,
		Scope:        scope,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
		Status:       "authorized",
		AuthorizedAt: time.Now(),
		ExpiresAt:    time.Now().Add(5 * time.Minute), // 与访问令牌有效期一致
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(&authRecord).Error; err != nil {
		fmt.Printf("记录OAuth2授权信息失败: %v\n", err)
	}

	// 构建重定向URL
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state)

	response.Response(c, http.StatusOK, 200, gin.H{
		"redirect_url": redirectURL,
		//"code": code,
		//"state": state,
	}, "授权成功")
}

// 三方APP取得临时code后，其服务端带着code来请求，根据scope返回对应的access_token
func GetTokenbyOAuth2(c *gin.Context) {
	db := common.GetDB()

	code := c.PostForm("code")
	appId := c.PostForm("app_id")
	appSecret := c.PostForm("app_secret")

	if code == "" || appId == "" || appSecret == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "缺少必要参数")
		return
	}

	// 验证应用
	var app model.OAuth2App
	if err := db.Where("app_id = ? AND is_active = ?", appId, true).First(&app).Error; err != nil {
		response.Response(c, http.StatusNotFound, 404, nil, "应用不存在或未激活")
		return
	}

	// 验证app_secret是否与app_id对应，避免其他应用冒用临时code
	if app.AppSecret != appSecret {
		response.Response(c, http.StatusUnauthorized, 401, nil, "应用密钥错误")
		return
	}

	// 查找授权码，临时code只允许使用1次
	var authCode model.OAuth2Code
	if err := db.Where("code = ? AND app_id = ? AND is_used = ?", code, appId, false).First(&authCode).Error; err != nil {
		response.Response(c, http.StatusNotFound, 404, nil, "授权码不存在或已使用")
		return
	}

	// 检查授权码是否过期
	if time.Now().After(authCode.ExpiresAt) {
		response.Response(c, http.StatusBadRequest, 400, nil, "授权码已过期")
		return
	}

	// 标记授权码为已使用
	authCode.IsUsed = true
	db.Save(&authCode)

	// 生成访问令牌
	accessToken, err := common.ReleaseOAuth2Token(authCode.UserID, appId, authCode.Scope)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "生成访问令牌失败")
		return
	}

	// 保存令牌到数据库
	oauth2Token := model.OAuth2Token{
		AccessToken:  accessToken,
		RefreshToken: "", // 暂不实现refreshToken
		AppID:        appId,
		UserID:       authCode.UserID,
		Scope:        authCode.Scope,
		ExpiresAt:    time.Now().Add(5 * time.Minute), // 5分钟有效期
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := db.Create(&oauth2Token).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "保存访问令牌失败")
		return
	}

	response.Response(c, http.StatusOK, 200, gin.H{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   300, // 5分钟
		"scope":        authCode.Scope,
	}, "获取访问令牌成功")
}

// 三方APP取得access_token后，其服务端带着access_token来请求，验证token确实有对应权限后（这在oauthMiddleware实现），返回对应的用户信息
func GetInfobyOAuth2(c *gin.Context) {
	// 从OAuth中间件获取用户信息
	userInterface, exists := c.Get("oauth2_user")
	if !exists {
		response.Response(c, http.StatusUnauthorized, 401, nil, "OAuth2认证失败")
		return
	}
	user := userInterface.(model.User)

	// 获取权限范围
	scopeInterface, exists := c.Get("oauth2_scope")
	if !exists {
		response.Response(c, http.StatusUnauthorized, 401, nil, "权限范围获取失败")
		return
	}
	scope := scopeInterface.(string)

	// 根据权限范围返回不同信息
	userInfo := gin.H{}

	if scope == "read_basic" {
		userInfo = gin.H{
			"user_id":    user.UserID,
			"name":       user.Name,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"intro":      user.Intro,
			"identity":   user.Identity,
		}
	} else {
		response.Response(c, http.StatusForbidden, 403, nil, "权限不足")
		return
	}

	response.Response(c, http.StatusOK, 200, userInfo, "获取用户信息成功")
}

// 用户获取OAuth2授权记录
func GetOAuth2AuthRecords(c *gin.Context) {
	// 从Auth中间件获取用户信息
	userInterface, exists := c.Get("user")
	if !exists {
		response.Response(c, http.StatusUnauthorized, 401, nil, "用户未登录")
		return
	}
	user := userInterface.(model.User)

	db := common.GetDB()

	// 获取分页参数
	page := 1
	pageSize := 10
	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsedPageSize, err := strconv.Atoi(ps); err == nil && parsedPageSize > 0 && parsedPageSize <= 100 {
			pageSize = parsedPageSize
		}
	}

	// 查询授权记录
	var records []model.OAuth2AuthRecord
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	db.Model(&model.OAuth2AuthRecord{}).Where("user_id = ?", user.UserID).Count(&total)

	// 获取记录列表
	if err := db.Where("user_id = ?", user.UserID).
		Order("authorized_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "查询授权记录失败")
		return
	}

	// 构建返回数据
	var recordList []gin.H
	for _, record := range records {
		recordList = append(recordList, gin.H{
			"record_id":     record.RecordID,
			"app_id":        record.AppID,
			"app_name":      record.AppName,
			"scope":         record.Scope,
			"ip_address":    record.IPAddress,
			"user_agent":    record.UserAgent,
			"status":        record.Status,
			"authorized_at": record.AuthorizedAt.Format("2006-01-02 15:04:05"),
			"expires_at":    record.ExpiresAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 返回分页数据
	response.Response(c, http.StatusOK, 200, gin.H{
		"records":     recordList,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}, "获取授权记录成功")
}

func ValidateEmail(c *gin.Context) {
	var request requestEmail
	c.Bind(&request)
	email := request.Email
	mode := request.Mode
	fmt.Println("email = ", email)
	fmt.Println("mode = ", mode)

	//arr := strings.Split(email, "@")
	//fmt.Println(arr[1])
	//if arr[1] != "mail2.sysu.edu.cn" && arr[1] != "mail.sysu.edu.cn" {
	//	response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email, "mode": mode}, "请使用中大邮箱！")
	//	return
	//}

	db := common.DB
	if mode == 1 && !isEmailExist(db, email) {
		//fmt.Println(1)
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email, "mode": mode}, "该邮箱未注册，请先完成注册操作")
		return
	}
	if mode == 0 && isEmailExist(db, email) {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email, "mode": mode}, "邮箱已注册")
		return
	}
	if email == "" {
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email}, "邮箱获取错误")
		return
	}
	err := api.SendEmail(email, 0, "")
	if err != nil {
		msg := err.Error()
		if msg == "" {
			msg = "发送邮箱错误"
		}
		response.Response(c, http.StatusBadRequest, 400, gin.H{"data": email}, msg)
		return
	}
	response.Success(c, gin.H{"data": email}, "邮箱发送成功")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件上传失败")
		return
	}

	// db := common.GetDB()
	// var user model.User

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filepath := "public/uploads/" + filename

	// 保存文件到本地
	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件保存失败")
		return
	}
	// 压缩头像
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

	// 对图像进行压缩处理
	resizedImage := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	// 这里似乎只支持绝对路径，这里要根据服务器修改
	resizedPath := "public/resized/" + filename
	// 服务器本地备份
	err = imaging.Save(resizedImage, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	fileURL := ""
	fileURL, err = api.UploadImage("/src/images/resized/"+filename, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 更新用户头像URL
	// 我们存储的是可以通过 HTTP 访问的 URL，而不是服务器本地的文件路径
	// user.AvatarURL = fileURL
	// user.AvatarURL = "https://ssemarket.cn:8080/uploads/" + filename
	// db.Save(&user)

	// 返回成功
	response.Success(c, gin.H{"avatar_url": fileURL}, "上传成功")
}

func UpdateAvatar(c *gin.Context) {
	// 用于移动端
	phone := c.PostForm("phone")
	if len(phone) != 11 {
		response.Response(c, http.StatusBadRequest, 400, nil, "Invalid phone number")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件上传失败")
		return
	}

	// Add a check for the file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" { // Add more formats if needed
		response.Response(c, http.StatusBadRequest, 400, nil, "Invalid file format")
		return
	}

	db := common.GetDB()
	var user model.User

	if db.Where("phone = ?", phone).First(&user).RecordNotFound() {
		response.Response(c, http.StatusNotFound, 404, nil, "User not found")
		return
	}
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename)) // Use base to avoid path traversal
	filepath := "public/uploads/" + filename

	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "文件保存失败")
		return
	}
	// 压缩头像
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

	// 对图像进行压缩处理
	resizedImage := imaging.Thumbnail(img, 100, 100, imaging.Lanczos)
	// 这里似乎只支持绝对路径，这里要根据服务器修改
	resizedPath := "public/resized/" + filename
	// 服务器本地备份
	err = imaging.Save(resizedImage, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	fileURL := ""
	fileURL, err = api.UploadImage("/src/images/resized/"+filename, resizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缩略图生成保存失败" + err.Error()})
		return
	}

	// 更新 User 的 AvatarURL 字段
	// user.AvatarURL = "https://ssemarket.cn:8080/uploads/" + filename
	user.AvatarURL = fileURL
	if err := db.Save(&user).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "Database update failed")
		return
	}

	response.Success(c, gin.H{"phone": user.Phone, "avatar_url": user.AvatarURL}, "上传成功")
}

// GetAvatar 用于处理获取用户头像的请求
func GetAvatar(c *gin.Context) {
	// 连接数据库
	db := common.GetDB()

	// 从前端获取电话号码
	phone := c.PostForm("phone")
	// 在数据库中查找用户
	var user model.User
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		} else {
			response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		}
		return
	}
	// 返回用户的头像URL
	response.Success(c, gin.H{"phone": user.Phone, "AvatarURL": user.AvatarURL}, "获取成功")
}
func GetInfo(c *gin.Context) {
	db := common.GetDB()
	// 从前端获取电话号码
	var requestUser = model.User{}
	c.Bind(&requestUser)
	phone := requestUser.Phone
	userID := requestUser.UserID
	// 在数据库中查找用户
	var user model.User
	if err := db.Where("phone = ? or userID = ?", phone, userID).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		} else {
			response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		}
		return
	}
	info := gin.H{
		"userID":           user.UserID,
		"name":             user.Name,
		"phone":            user.Phone,
		"intro":            user.Intro,
		"score":            user.Score,
		"ban":              user.Banend,
		"punishnum":        user.Punishnum,
		"avatarURL":        user.AvatarURL,
		"hideProfilePosts": user.HideProfilePosts,
	}

	// 判断为本人
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != 0 && tokenUserID == user.UserID {
		info["email"] = user.Email
		info["emailNotificationBlocked"] = user.ISEmailNotificationBlocked
		// emailpush 历史字段名易误解为「已开启推送」，实际同为 is_email_notification_blocked
		info["emailpush"] = user.ISEmailNotificationBlocked
	}

	c.JSON(http.StatusOK, info)
}

type updateUser struct {
	UserID           int    `json:"userID"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Name             string `json:"name"`
	Num              int    `json:"num"`
	Intro            string `json:"intro"`
	AvatarURL        string `json:"avatarURL"`
	HideProfilePosts *bool  `json:"hideProfilePosts,omitempty"`
}

func UpdateUserInfo(c *gin.Context) {
	db := common.GetDB()

	// 解析请求参数
	var userInfo updateUser
	if err := c.Bind(&userInfo); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数解析错误")
		return
	}

	// 根据用户ID查找用户
	var user model.User
	if err := db.Where("userID = ?", userInfo.UserID).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Response(c, http.StatusNotFound, 404, nil, "用户不存在")
		} else {
			response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		}
		return
	}
	// 获取token中的用户标识符
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != user.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 400, nil, "权限不足")
		return
	}

	// 更新用户信息
	user.Name = userInfo.Name
	//user.Num = userInfo.Num
	user.Intro = userInfo.Intro
	user.AvatarURL = userInfo.AvatarURL
	if userInfo.HideProfilePosts != nil {
		user.HideProfilePosts = *userInfo.HideProfilePosts
	}

	if err := db.Save(&user).Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "数据库错误")
		return
	}
	response.Response(c, http.StatusOK, 200, nil, "用户信息更新成功")
}

func CalculateAndSaveScores() {
	log.Println("CalculateAndSaveScores")
	db := common.GetDB()
	var users []model.User
	db.Find(&users)
	for _, user := range users { //对每个用户进行积分统计
		if user.UserID == 0 {
			continue
		}
		// 查询用户的post ID
		var postIDs []int
		db.Model(&model.Post{}).Where("userID = ?", user.UserID).Pluck("postID", &postIDs)
		// 统计post被点赞个数
		totalLikeNum := 0
		for _, postID := range postIDs {
			var likeNum int
			db.Model(&model.Post{}).Where("postID = ?", postID).Select("like_num").Row().Scan(&likeNum)
			totalLikeNum += likeNum
		}
		// 查询用户的 pcomment ID
		var pcommentIDs []int
		db.Model(&model.Pcomment{}).Where("userID = ?", user.UserID).Pluck("pcommentID", &pcommentIDs)
		// 统计pcomment被点赞个数
		for _, pcommentID := range pcommentIDs {
			var likeNum int
			db.Model(&model.Pcomment{}).Where("pcommentID = ?", pcommentID).Select("like_num").Row().Scan(&likeNum)
			totalLikeNum += likeNum
		}
		// 统计被收藏个数
		var psaveCount int
		db.Model(&model.Psave{}).Where("ptargetID IN (?)", postIDs).Count(&psaveCount)
		// 统计帖子被评论个数
		var commentedCount int
		db.Model(&model.Pcomment{}).Where("ptargetID IN (?)", postIDs).Not("userID = ?", user.UserID).Count(&commentedCount)
		// 统计评论被回复个数
		var repliedCount int
		db.Model(&model.Ccomment{}).Where("ctargetID IN (?)", pcommentIDs).Not("userID = ?", user.UserID).Count(&repliedCount)
		// 查询用户的 ccomment ID
		var ccommentIDs []int
		db.Model(&model.Ccomment{}).Where("userID = ?", user.UserID).Pluck("ccommentID", &ccommentIDs)
		// 统计ccomment被点赞个数
		for _, ccommentID := range ccommentIDs {
			var likeNum int
			db.Model(&model.Ccomment{}).Where("ccommentID = ?", ccommentID).Select("like_num").Row().Scan(&likeNum)
			totalLikeNum += likeNum
		}
		// 统计成功举报个数
		var sueCount int
		db.Model(&model.Sue{}).Where("status = ? AND finish = ? AND userID = ?", " ok", 1, user.UserID).
			Select("DISTINCT targettype, targetID").
			Count(&sueCount)
		// 计算总积分并保存
		totalScore := len(postIDs)*10 + len(pcommentIDs)*5 + len(ccommentIDs)*5 + totalLikeNum + psaveCount*3 +
			(commentedCount+repliedCount)*2 + sueCount*20 - (user.Punishnum)*20
		db.Model(&user).Update("Score", totalScore)
	}
	log.Println("Score updated")
}

func ChangeEmailPushStatus(c *gin.Context) {
	db := common.GetDB()
	var user_id UserMsg
	var user model.User
	if err := c.Bind(&user_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Where("userID = ?", user_id.UserID).First(&user)
	if user.ISEmailNotificationBlocked {
		db.Model(&user).Update("is_email_notification_blocked", false)
	} else {
		db.Model(&user).Update("is_email_notification_blocked", true)
	}
}

type ChatUser struct {
	UserID int `json:"userID"`
	// Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarURL"`
	Identity  string `json:"identity"`
	Score     int    `json:"score"`
}

func GetStatistics(c *gin.Context) {
	db := common.GetDB()
	userId := GetTokenUserID(c)

	user := model.User{}
	db.Where("userID = ?", userId).First(&user)

	curTime := time.Now()
	durationTime := curTime.Sub(user.Banend)

	lastYear := curTime.Year() - 1
	statisticsStartTime := time.Date(lastYear, 1, 1, 0, 0, 0, 0, curTime.Location())
	statisticsEndTime := time.Date(curTime.Year(), 1, 1, 0, 0, 0, 0, curTime.Location())

	var posts []model.Post
	var postIDs []int
	db.Model(&model.Post{}).Where("userID = ?", userId).Find(&posts)

	totalPostCnt := len(posts)

	totalPostLikeNum := 0
	thisYearPostCnt := 0
	thisYearLikeNum := 0

	maxPostLikeNum := 0
	maxPostBrowseNum := 0
	maxPostCommentNum := 0
	for _, post := range posts {
		postIDs = append(postIDs, post.PostID)
		totalPostLikeNum += post.LikeNum
		if post.PostTime.After(statisticsStartTime) && post.PostTime.Before(statisticsEndTime) {
			thisYearPostCnt++
			thisYearLikeNum += post.LikeNum

			if post.LikeNum > maxPostLikeNum {
				maxPostLikeNum = post.LikeNum
			}
			if post.BrowseNum > maxPostBrowseNum {
				maxPostBrowseNum = post.BrowseNum
			}
			if post.CommentNum > maxPostCommentNum {
				maxPostCommentNum = post.CommentNum
			}
		}
	}

	// 查询用户的 pcomment ID
	var pcommentIDs []int
	db.Model(&model.Pcomment{}).Where("userID = ?", user.UserID).Pluck("pcommentID", &pcommentIDs)
	pCommentCnt := len(pcommentIDs)
	pCommentLikeNum := 0

	// 统计pcomment被点赞个数
	for _, pcommentID := range pcommentIDs {
		var likeNum int
		db.Model(&model.Pcomment{}).Where("pcommentID = ?", pcommentID).Select("like_num").Row().Scan(&likeNum)
		pCommentLikeNum += likeNum
	}
	// 统计被收藏个数
	var besavedCount int
	db.Model(&model.Psave{}).Where("ptargetID IN (?)", postIDs).Count(&besavedCount)
	// 统计帖子被评论个数
	var commentedCount int
	db.Model(&model.Pcomment{}).Where("ptargetID IN (?)", postIDs).Not("userID = ?", user.UserID).Count(&commentedCount)
	// 统计评论被回复个数
	var repliedCount int
	db.Model(&model.Ccomment{}).Where("ctargetID IN (?)", pcommentIDs).Not("userID = ?", user.UserID).Count(&repliedCount)
	// 查询用户的 ccomment ID
	var ccommentIDs []int
	db.Model(&model.Ccomment{}).Where("userID = ?", user.UserID).Pluck("ccommentID", &ccommentIDs)
	ccommentCnt := len(ccommentIDs)
	ccommentLikeNum := 0
	// 统计ccomment被点赞个数
	for _, ccommentID := range ccommentIDs {
		var likeNum int
		db.Model(&model.Ccomment{}).Where("ccommentID = ?", ccommentID).Select("like_num").Row().Scan(&likeNum)
		ccommentLikeNum += likeNum
	}

	var savedCount int
	db.Model(&model.Psave{}).Where("userID = ?", user.UserID).Count(&savedCount)

	chatMsgs := []model.ChatMsg{}
	err := db.Where("senderUserID = ? OR targetUserID = ?", userId, userId).
		Preload("TargetUser").
		Preload("SenderUser").
		Find(&chatMsgs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chatCount := len(chatMsgs)

	relevantRespUsers := []ChatUser{}
	existsUserIDs := make(map[int]bool) // 去重：key=用户ID，value=是否已存在
	chatUserCnt := make(map[int]int)    // 统计每个用户的聊天句子数量：key=用户ID，value=聊天句子数量
	for _, chatMsg := range chatMsgs {
		if userId != chatMsg.TargetUserID {
			if _, ok := existsUserIDs[chatMsg.TargetUserID]; !ok {
				existsUserIDs[chatMsg.TargetUserID] = true
				relevantRespUsers = append(relevantRespUsers, ChatUser{
					UserID: chatMsg.TargetUserID,
					// Email:     chatMsg.TargetUser.Email,
					Name:      chatMsg.TargetUser.Name,
					AvatarURL: chatMsg.TargetUser.AvatarURL,
					Identity:  chatMsg.TargetUser.Identity,
					Score:     chatMsg.TargetUser.Score,
				})
			}
			chatUserCnt[chatMsg.TargetUserID]++
		}

		if userId != chatMsg.SenderUserID {
			if _, ok := existsUserIDs[chatMsg.SenderUserID]; !ok {
				existsUserIDs[chatMsg.SenderUserID] = true
				relevantRespUsers = append(relevantRespUsers, ChatUser{
					UserID: chatMsg.SenderUserID,
					// Email:     chatMsg.SenderUser.Email,
					Name:      chatMsg.SenderUser.Name,
					AvatarURL: chatMsg.SenderUser.AvatarURL,
					Identity:  chatMsg.SenderUser.Identity,
					Score:     chatMsg.SenderUser.Score,
				})
			}
			chatUserCnt[chatMsg.SenderUserID]++
		}
	}

	maxSayUserID := 0
	maxSayUserCnt := 0
	for userID, cnt := range chatUserCnt {
		if cnt > maxSayUserCnt {
			maxSayUserID = userID
			maxSayUserCnt = cnt
		}
	}

	var maxSayUser ChatUser
	for _, user := range relevantRespUsers {
		if user.UserID == maxSayUserID {
			maxSayUser = user
			break
		}
	}

	statistics := gin.H{
		"name":         user.Name,      // 用户名
		"intro":        user.Intro,     // 用户介绍
		"score":        user.Score,     // 用户积分
		"avatarUrl":    user.AvatarURL, // 用户头像URL
		"durationTime": durationTime,   // 注册时间

		"totalPostCnt":     totalPostCnt,      // 总帖子数量
		"totalPostLikeNum": totalPostLikeNum,  // 总点赞数量
		"thisYearPostCnt":  thisYearPostCnt,   // 这一年发布的帖子数量
		"thisYearLikeNum":  thisYearLikeNum,   // 这一年收到的点赞数量
		"maxLikeNum":       maxPostLikeNum,    // 这一年发布的帖子中点赞数量最多的帖子的点赞数量
		"maxBrowseNum":     maxPostBrowseNum,  // 这一年发布的帖子中浏览数量最多的帖子的浏览数量
		"maxCommentNum":    maxPostCommentNum, // 这一年发布的帖子中评论数量最多的帖子的评论数量

		"pCommentCnt":     pCommentCnt,     // 总评论数量
		"pCommentLikeNum": pCommentLikeNum, // 总评论点赞数量
		"ccommentCnt":     ccommentCnt,     // 总回复数量
		"ccommentLikeNum": ccommentLikeNum, // 总回复点赞数量\

		"beSavedCount":     besavedCount,   // 总被收藏数量
		"beCommentedCount": commentedCount, // 总被评论数量
		"repliedCount":     repliedCount,   // 总被回复数量

		"savedCount":        savedCount,        // 总收藏数量
		"chatCount":         chatCount,         // 总聊天句子数量
		"relevantRespUsers": relevantRespUsers, // 相关聊天用户列表
		"maxSayUser":        maxSayUser,        // 最常联系用户
		"maxSayUserCnt":     maxSayUserCnt,     // 最常联系用户的聊天句子数量
	}
	response.Response(c, http.StatusOK, 200, statistics, "获取用户统计信息成功")
}
