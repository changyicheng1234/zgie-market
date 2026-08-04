package controller

import (
	"loginTest/common"
	"loginTest/config"
	"loginTest/model"
	"loginTest/response"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
)

type NoticeResponse struct {
	NoticeID     int       `json:"noticeID"`
	ReceiverName string    `json:"receiverName"`
	SenderName   string    `json:"senderName"`
	SenderAvatar string    `json:"senderAvatar"`
	SenderID     int       `json:"senderID"`
	Type         string    `json:"type"`
	Content      string    `json:"content"`
	Read         bool      `json:"read"`
	PostID       int       `json:"postID"`
	Target       int       `json:"target"`
	PcommentID   int       `json:"pcommentID"`
	Time         time.Time `json:"time"`
}
type NoticeGet struct {
	NoticeList     []NoticeResponse `json:"noticeList"`
	LastID         int              `json:"lastID"`
	TotalNum       int              `json:"totalNum"`
	UnreadTotalNum int              `json:"unreadTotalNum"`
}
type NoticeNumResponse struct {
	TotalNum       int `json:"totalNum"`
	UnreadTotalNum int `json:"unreadTotalNum"`
	ReadTotalNum   int `json:"readTotalNum"`
}

func GetNotice(c *gin.Context) {
	db := common.GetDB()
	//从中间件存入的user中获取userID
	value, exisits := c.Get("user")
	pageSizeStr := c.Query("pageSize")
	requireIDStr := c.Query("requireID")
	readStr := c.Query("read")
	requireID, _ := strconv.Atoi(requireIDStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	read, _ := strconv.Atoi(readStr)
	var user model.User
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "游客无法访问通知")
		return
	} else {
		user = value.(model.User)
	}
	//执行操作，分页返回通知
	var notices []model.Notice
	var total int
	var unreadTotal int
	// requireID==0说明是首次查询,返回pagesize条通知,计算totalNum
	if requireID == 0 {
		db.Model(&model.Notice{}).Where("receiver =?", user.UserID).Count(&total)
		db.Model(&model.Notice{}).Where("receiver =? AND `read` =?", user.UserID, read).Count(&unreadTotal)
		db.Where("receiver =? AND `read` =?", user.UserID, read).Order("noticeID DESC").Limit(pageSize).Find(&notices)
	} else { //否则查询比requireID小的通知
		db.Where("receiver =? AND `read` =? AND noticeID < ?", user.UserID, read, requireID).Order("noticeID DESC").Limit(pageSize).Find(&notices)
	}
	if len(notices) == 0 {
		response.Response(c, http.StatusOK, 201, nil, "没有更多通知")
		return
	} else if len(notices) < 5 {

	}
	var noticeGet NoticeGet
	noticeGet.TotalNum = total
	noticeGet.LastID = notices[len(notices)-1].NoticeID
	for _, notice := range notices {
		var temuser model.User
		db.Where("userID =?", notice.Sender).First(&temuser)
		var tempcomment model.Pcomment
		if notice.Type == "ccomment" {
			var temccoment model.Ccomment
			db.Where("ccommentID=?", notice.Target).First(&temccoment)
			db.Where("pcommentID =?", temccoment.CtargetID).First(&tempcomment)
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				SenderID:     notice.Sender,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
				PcommentID:   temccoment.CtargetID,
				PostID:       tempcomment.PtargetID,
				Time:         temccoment.Time,
			})
		} else if notice.Type == "pcomment" {
			db.Where("pcommentID =?", notice.Target).First(&tempcomment)
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				SenderID:     notice.Sender,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
				PostID:       tempcomment.PtargetID,
				Time:         tempcomment.Time,
			})
		} else if notice.Type == "post" {
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				SenderID:     notice.Sender,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
				PostID:       notice.Target,
				Time:         notice.Time,
			})
		} else {
			noticeGet.NoticeList = append(noticeGet.NoticeList, NoticeResponse{
				NoticeID:     notice.NoticeID,
				ReceiverName: user.Name,
				SenderName:   temuser.Name,
				SenderAvatar: temuser.AvatarURL,
				SenderID:     notice.Sender,
				Type:         notice.Type,
				Content:      notice.Ntext,
				Read:         notice.Read,
				Target:       notice.Target,
				Time:         notice.Time,
			})
		}

	}
	c.JSON(http.StatusOK, &noticeGet)
}

func GetNoticeNum(c *gin.Context) {
	db := common.GetDB()
	//从中间件存入的user中获取userID
	value, exisits := c.Get("user")
	var user model.User
	if !exisits {
		response.Response(c, http.StatusBadRequest, 400, nil, "游客无法获得通知数量")
		return
	} else {
		user = value.(model.User)
	}
	var noticeNum NoticeNumResponse
	var readNum int
	var unreadNum int
	db.Model(&model.Notice{}).Where("receiver =? AND `read` =0", user.UserID).Count(&unreadNum)
	db.Model(&model.Notice{}).Where("receiver =? AND `read` =1", user.UserID).Count(&readNum)
	noticeNum.TotalNum = readNum + unreadNum
	noticeNum.ReadTotalNum = readNum
	noticeNum.UnreadTotalNum = unreadNum
	c.JSON(http.StatusOK, &noticeNum)
}

func ReadNotice(c *gin.Context) {
	noticeID := c.Param("noticeID")
	db := common.GetDB()
	// Update notice data in database to set 'read' field to true
	err := db.Model(&model.Notice{}).Where("noticeID = ?", noticeID).Update("read", true).Error
	if err != nil {
		c.JSON(500, gin.H{"更新已读出现error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func EmailNotice(c *gin.Context) {
	value, exists := c.Get("admin")
	var admin model.Admin
	if !exists {
		response.Fail(c, nil, "验证错误")
		return
	}
	admin = value.(model.Admin)
	if admin.AdminID == 0 {
		response.Fail(c, nil, "管理员不存在")
	}

	// 检查通知邮件配置是否已初始化
	if config.NoticeEmailUsername == "" || config.NoticePassword == "" || config.NoticeAddr == "" || config.NoticeHost == "" {
		response.Fail(c, nil, "通知邮件配置未正确初始化")
		return
	}

	e := email.NewEmail()
	e.From = "软工集市 <" + config.NoticeEmailUsername + ">"
	e.To = []string{config.NoticeEmailUsername}
	e.Subject = "软工集市version2.0升级部署通知"
	e.HTML = []byte("<!DOCTYPE html>\n<html>\n<head>\n    <title>软工集市平台升级维护通知</title>\n</head>\n<body>\n    <h1>软工集市平台升级维护通知</h1>\n    <p>尊敬的软工集市用户，</p>\n    <p>为了提供更好的服务和功能，我们计划进行软工集市平台的升级维护。这次升级将带来一些重要的改进和新功能，以增强您的使用体验。</p>\n    <p>升级维护计划如下：</p>\n    <ul>\n        <li>维护开始时间：[维护开始日期和时间]</li>\n        <li>维护结束时间：[维护结束日期和时间]</li>\n        <li>期间软工集市将无法访问，为了不影响您的工作，请在维护前备份您的数据。</li>\n    </ul>\n    <p>升级后，您可以期待以下改进：</p>\n    <ul>\n        <li>[列出将引入的新功能和改进]</li>\n        <li>[列出其他相关信息]</li>\n    </ul>\n    <p>我们对给您带来的不便深感抱歉，但我们相信这些改进将使您更加满意。在维护期间，如果您有任何疑问或需要帮助，请随时联系我们的客户支持团队。</p>\n    <p>感谢您一直以来对软工集市的支持和信任。</p>\n    <p>祝您使用愉快！</p>\n    <p>谨启</p>\n    <p>软工集市团队</p>\n</body>\n</html>\n")
	err := e.Send(config.NoticeAddr, smtp.PlainAuth("", config.NoticeEmailUsername, config.NoticePassword, config.NoticeHost))
	if err != nil {
		response.Fail(c, nil, "邮件发送失败: "+err.Error())
		return
	}
	response.Success(c, nil, "通知邮件发送成功")
}
