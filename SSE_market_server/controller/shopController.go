package controller

import (
	"fmt"
	"loginTest/common"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

type ProductMsg struct {
	UserID       int      `json:"UserID"`
	Price        int      `json:"Price"`
	Title        string   `json:"Title"`
	Content      string   `json:"Content"`
	Photos       []string `json:"Photos"`
	ISAnonymous  bool     `json:"ISAnonymous"`
	SyncFeedPost *bool    `json:"syncFeedPost"`
}

type ProductReponse struct {
	SellerID    int       `json:"SellerID"`
	ProductID   int       `json:"ProductID"`
	Seller      string    `json:"Seller"`
	Price       int       `json:"Price"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Photos      []string  `json:"Photos"`
	ISSold      bool      `json:"ISSold"`
	ISAnonymous bool      `json:"ISAnonymous"`
	PostTime    time.Time `json:"PostTime"`
}

type ShopBrowseMeg struct {
	UserID     int    `json:"UserID"`
	Searchinfo string `json:"searchinfo"`
	Searchsort string `json:"searchsort"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type ProductIDmsg struct {
	ProductID int `json:"productID"`
}

func PostProduct(c *gin.Context) {
	db := common.GetDB()
	var msg ProductMsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	title := strings.TrimSpace(msg.Title)
	content := strings.TrimSpace(msg.Content)
	if !isUserExist(db, msg.UserID) {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户不存在")
		return
	}
	if title == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题不能为空")
		return
	}
	if utf8.RuneCountInString(title) > 30 {
		response.Response(c, http.StatusBadRequest, 400, nil, "标题最多为30个字")
		return
	}
	if content == "" {
		response.Response(c, http.StatusBadRequest, 400, nil, "内容不能为空")
		return
	}
	if utf8.RuneCountInString(content) > 2000 {
		response.Response(c, http.StatusBadRequest, 400, nil, "描述最多2000字")
		return
	}
	if msg.Price < 0 || msg.Price > 1000000 {
		response.Response(c, http.StatusBadRequest, 400, nil, "价格无效")
		return
	}
	tokenUserID := GetTokenUserID(c)
	if tokenUserID != msg.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 403, nil, "权限不足")
		return
	}
	photos := make([]string, 0, len(msg.Photos))
	for _, p := range msg.Photos {
		if s := strings.TrimSpace(p); s != "" {
			photos = append(photos, s)
		}
	}
	syncFeed := true
	if msg.SyncFeedPost != nil {
		syncFeed = *msg.SyncFeedPost
	}

	var user model.User
	db.Where("userID = ?", msg.UserID).First(&user)

	newProduct := model.Product{
		UserID:      msg.UserID,
		Name:        title,
		Description: content,
		PostTime:    time.Now(),
		Photos:      strings.Join(photos, ","),
		ISSold:      false,
		ISAnonymous: msg.ISAnonymous,
		Price:       msg.Price,
	}

	tx := db.Begin()
	if err := tx.Create(&newProduct).Error; err != nil {
		tx.Rollback()
		response.Response(c, http.StatusInternalServerError, 500, nil, "发布失败")
		return
	}
	if syncFeed {
		if err := createProductFeedPost(tx, user, newProduct); err != nil {
			tx.Rollback()
			response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "发布失败")
		return
	}
	common.BloomFilter.Add([]byte(strconv.FormatUint(uint64(newProduct.ProductID), 10)))
	response.Response(c, http.StatusOK, 200, gin.H{"ProductID": newProduct.ProductID}, "商品发布成功")
}

func GetProducts(c *gin.Context) {
	db := common.GetDB()
	var msg ShopBrowseMeg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	tokenUserID := GetTokenUserID(c)
	limit := msg.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	offset := msg.Offset
	if offset < 0 {
		offset = 0
	}

	query := db.Model(&model.Product{}).Where("is_delete = ?", false)
	searchinfo := strings.TrimSpace(msg.Searchinfo)
	if msg.Searchsort == "history" {
		query = query.Where("userID = ?", tokenUserID)
	} else if searchinfo != "" {
		like := "%" + searchinfo + "%"
		query = query.Where("(product_name LIKE ? OR description LIKE ?)", like, like)
	}

	var products []model.Product
	query.Order("is_sold ASC, productID DESC").Offset(offset).Limit(limit).Find(&products)

	out := make([]ProductReponse, 0, len(products))
	for _, p := range products {
		out = append(out, productToResponse(db, p, tokenUserID, msg.Searchsort == "history"))
	}
	c.JSON(http.StatusOK, out)
}

func GetProductDetail(c *gin.Context) {
	db := common.GetDB()
	var msg ProductIDmsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	var product model.Product
	db.Where("productID = ? AND is_delete = ?", msg.ProductID, false).First(&product)
	if product.ProductID == 0 {
		response.Response(c, http.StatusNotFound, 404, nil, "商品不存在")
		return
	}
	tokenUserID := GetTokenUserID(c)
	c.JSON(http.StatusOK, productToResponse(db, product, tokenUserID, product.UserID == tokenUserID))
}

func GetProductNum(c *gin.Context) {
	db := common.GetDB()
	var msg ShopBrowseMeg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	tokenUserID := GetTokenUserID(c)
	query := db.Model(&model.Product{}).Where("is_delete = ?", false)
	searchinfo := strings.TrimSpace(msg.Searchinfo)
	if msg.Searchsort == "history" {
		query = query.Where("userID = ?", tokenUserID)
	} else if searchinfo != "" {
		like := "%" + searchinfo + "%"
		query = query.Where("(product_name LIKE ? OR description LIKE ?)", like, like)
	}
	var count int
	query.Count(&count)
	c.JSON(http.StatusOK, gin.H{"Productcount": count})
}

func DeleteProduct(c *gin.Context) {
	db := common.GetDB()
	var msg ProductIDmsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	var product model.Product
	db.Where("productID = ?", msg.ProductID).First(&product)
	if product.ProductID == 0 || product.ISDelete {
		response.Response(c, http.StatusBadRequest, 400, nil, "商品不存在")
		return
	}
	if GetTokenUserID(c) != product.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 403, nil, "权限不足")
		return
	}
	db.Model(&product).Update("is_delete", true)
	c.JSON(http.StatusOK, gin.H{"message": "商品删除成功"})
}

func SaleProduct(c *gin.Context) {
	db := common.GetDB()
	var msg ProductIDmsg
	if err := c.ShouldBindBodyWith(&msg, binding.JSON); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	var product model.Product
	db.Where("productID = ?", msg.ProductID).First(&product)
	if product.ProductID == 0 || product.ISDelete {
		response.Response(c, http.StatusBadRequest, 400, nil, "商品不存在")
		return
	}
	if GetTokenUserID(c) != product.UserID {
		response.Response(c, http.StatusUnprocessableEntity, 403, nil, "权限不足")
		return
	}
	db.Model(&product).Update("is_sold", true)
	c.JSON(http.StatusOK, gin.H{"message": "商品出售成功"})
}

func GetCarouselImg(c *gin.Context) {
	db := common.GetDB()
	var products []model.Product
	db.Where("is_delete = ? AND is_sold = ?", false, false).
		Where("photos != ?", "").
		Order("productID DESC").
		Limit(8).
		Find(&products)
	urls := make([]string, 0, len(products))
	for _, p := range products {
		for _, u := range strings.Split(p.Photos, ",") {
			u = strings.TrimSpace(u)
			if u == "" {
				continue
			}
			urls = append(urls, replaceResizedWithURL(u))
			break
		}
		if len(urls) >= 5 {
			break
		}
	}
	c.JSON(http.StatusOK, urls)
}

func productToResponse(db *gorm.DB, product model.Product, viewerID int, viewerIsOwner bool) ProductReponse {
	photos := splitProductPhotos(product.Photos)
	sellerID, sellerName := 0, "匿名卖家"
	if !product.ISAnonymous || viewerIsOwner {
		var user model.User
		db.Where("userID = ?", product.UserID).First(&user)
		if user.UserID != 0 {
			sellerID = user.UserID
			sellerName = user.Name
		}
	}
	return ProductReponse{
		ProductID:   product.ProductID,
		SellerID:    sellerID,
		Seller:      sellerName,
		Price:       product.Price,
		Name:        product.Name,
		Description: product.Description,
		Photos:      photos,
		ISSold:      product.ISSold,
		ISAnonymous: product.ISAnonymous,
		PostTime:    product.PostTime,
	}
}

func splitProductPhotos(raw string) []string {
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, replaceResizedWithURL(p))
	}
	return out
}

func isUserExist(db *gorm.DB, userID int) bool {
	var user model.User
	db.Where("userID = ?", userID).First(&user)
	return user.UserID != 0
}

func replaceResizedWithURL(photoURL string) string {
	return strings.ReplaceAll(photoURL, "/resized/", "/uploads/")
}

func createProductFeedPost(db *gorm.DB, user model.User, product model.Product) error {
	feedTitle := common.FeedTitleWithPrefix("新上架-", product.Name)
	feedContent := productFeedPostContent(product)
	return common.CreateHomeFeedPost(db, user, feedTitle, feedContent)
}

func productFeedPostContent(product model.Product) string {
	link := fmt.Sprintf("%s/shop/productdetail/%d", common.NewFrontendBaseURL, product.ProductID)
	var b strings.Builder
	b.WriteString("商城有新商品上架「")
	b.WriteString(product.Name)
	b.WriteString(fmt.Sprintf("」，标价 ¥%d。\n\n", product.Price))
	b.WriteString("商品链接：\n")
	b.WriteString(link)
	b.WriteString("\n")
	desc := strings.TrimSpace(product.Description)
	if desc != "" {
		if utf8.RuneCountInString(desc) > 400 {
			runes := []rune(desc)
			desc = string(runes[:400]) + "…"
		}
		b.WriteString("\n")
		b.WriteString(desc)
		b.WriteString("\n")
	}
	return b.String()
}
