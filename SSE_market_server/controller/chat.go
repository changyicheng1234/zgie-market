package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"loginTest/api"
	"loginTest/common"
	"loginTest/core"
	"loginTest/model"
	"loginTest/response"
	"loginTest/util"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// 定义响应用户结构体
type ChatRespUser struct {
	UserID             int    `json:"userID"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	AvatarURL          string `json:"avatarURL"`
	Identity           string `json:"identity"`
	Score              int    `json:"score"`
	UnRead             int    `json:"unRead"`
	LastMessageAt      string `json:"lastMessageAt"`
	LastMessagePreview string `json:"lastMessagePreview"`
}

// 定义心跳结构体
type HeartBeat struct {
	UserID int `json:"userID"`
	Beat   int `json:"beat"`
}

// ws连接节点
type ChatNode struct {
	Conn      *websocket.Conn
	DataQueue chan []byte //发送消息的消息队列
}

// 创建用户id与ws连接映射map
var clientMap sync.Map

// 获得聊天记录

type ChatHistoryQuery struct {
	SenderUserID int `form:"senderUserID" binding:"required"` // 绑定 query 参数 senderUserID
	TargetUserID int `form:"targetUserID" binding:"required"` // 绑定 query 参数 targetUserID
}

func GetChatHistory(c *gin.Context) {
	chatQuery := ChatHistoryQuery{}
	if err := c.ShouldBindQuery(&chatQuery); err != nil {
		response.Fail(c, gin.H{}, "无效的查询参数")
		return
	}
	//校验是否本人操作
	ctxUser, ok := c.Get("user")
	if !ok {
		response.Fail(c, gin.H{}, "获得聊天记录失败，请重试")
		return
	}
	currentUser, _ := ctxUser.(model.User)
	if currentUser.UserID != chatQuery.SenderUserID {
		response.Fail(c, gin.H{}, "没有权利查看他人聊天记录")
		return
	}
	_ = removeChatViewInRedis(chatQuery.SenderUserID)

	redisCtx, cancel := context.WithTimeout(c.Request.Context(), 800*time.Millisecond)
	defer cancel()
	chatViewKey := util.GenerateChatViewKey(chatQuery.SenderUserID, chatQuery.TargetUserID)
	//设置redis表示用户存在当前聊天界面
	if err := core.MyRedis.Set(redisCtx, chatViewKey, 1, 24*time.Hour).Err(); err != nil {
		response.Fail(c, gin.H{}, "获得聊天记录失败，请重试")
	}
	_ = core.MyRedis.SAdd(redisCtx, util.GenerateChatViewIndexKey(chatQuery.SenderUserID), chatViewKey).Err()
	//清空redis的未读情况
	unreadKey := util.GenerateUnreadKey(chatQuery.SenderUserID, chatQuery.TargetUserID)
	if err := core.MyRedis.Del(redisCtx, unreadKey).Err(); err != nil {
		response.Fail(c, gin.H{}, "获得聊天记录失败，请重试")
	}
	_ = core.MyRedis.SRem(redisCtx, util.GenerateUnreadIndexKey(chatQuery.SenderUserID), unreadKey).Err()

	//获得聊天信息
	chatList, err := GetChatHistoryService(chatQuery.SenderUserID, chatQuery.TargetUserID)
	if err != nil {
		response.Fail(c, gin.H{}, "获得聊天记录失败，请重试")
	} else {
		response.Success(c, gin.H{"chatHistoryList": chatList}, "获得聊天记录成功")
	}
}

// 获得聊天记录服务
type chatHistoryResp struct {
	ChatMsgID    int       `gorm:"column:chatMsgID" json:"chatMsgID"`
	TargetUserID int       `gorm:"column:targetUserID" json:"targetUserID"`
	SenderUserID int       `gorm:"column:senderUserID" json:"senderUserID"`
	Content      string    `gorm:"column:content" json:"content"`
	CreatedAt    time.Time `gorm:"column:createdAt" json:"createdAt"`
}

func GetChatHistoryService(userIdTarget int, fromUserId int) ([]chatHistoryResp, error) {
	var messageList []chatHistoryResp
	if err := common.DB.Model(model.ChatMsg{}).Where("(targetUserID = ? AND senderUserID = ?) OR (targetUserID = ? AND senderUserID = ?)", userIdTarget, fromUserId, fromUserId, userIdTarget).
		Order("createdAt").
		Scan(&messageList).Error; err != nil {
		return nil, err
	}
	return messageList, nil
}

// 取消用户redis里所有的在线状态
func removeChatViewInRedis(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	indexKey := util.GenerateChatViewIndexKey(userID)
	keys, err := core.MyRedis.SMembers(ctx, indexKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if len(keys) > 0 {
		if _, err = core.MyRedis.Del(ctx, keys...).Result(); err != nil {
			return err
		}
	}
	return core.MyRedis.Del(ctx, indexKey).Err()
}

// 聊天ws
func ChatHandler(c *gin.Context) {
	//获得当前操作用户的userID
	var userID int
	userInterface, ok := c.Get("user")
	if !ok {
		tokenString := c.Query("token")
		_, claim, err := common.ParseToken(tokenString)
		if err != nil {
			response.Fail(c, gin.H{}, "请重新登录")
		}
		userID = claim.UserID
	} else {
		user := userInterface.(model.User)
		userID = user.UserID
	}

	// 避免 Redis 异常阻塞握手：清理逻辑异步执行
	go removeChatViewInRedis(userID)

	chatFunc(c.Request, c.Writer, userID)
}

// 处理聊天逻辑
func chatFunc(request *http.Request, writer http.ResponseWriter, userID int) {
	//升级http请求为ws
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		// 打印错误并返回错误信息
		http.Error(writer, "Failed to upgrade WebSocket connection: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//	创建chatnode节点并且添加映射
	node := &ChatNode{
		conn,
		make(chan []byte, 200),
	}
	clientMap.Store(userID, node)

	//	启动发送消息协程
	go sendMsg2User(node)

	//	启动接收消息协程
	go recvMsgFromUser(node, userID)

	//	响应该用户所有的聊天用户信息
	relevantUserJSON, err := GetRelevantUser(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	sendMsg(0, userID, relevantUserJSON)

}

// 将用户消息队列的消息发送给用户
func sendMsg2User(node *ChatNode) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		}
	}
}

// 从请求读取用户发送的信息，并且将信息派发（考虑到后续群聊功能的实现因此实现dispatch，增加可扩展性）
func recvMsgFromUser(node *ChatNode, userID int) {
	defer func() {
		node.Conn.Close()        // 关闭连接
		clientMap.Delete(userID) // 移除用户的连接
	}()

	// 设置读超时时间，比如 60 秒
	readTimeout := 60 * time.Second

	for {
		// 设置读超时
		node.Conn.SetReadDeadline(time.Now().Add(readTimeout))

		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			// 检查是否是超时错误
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//取消用户redis里所有的在线状态
				removeChatViewInRedis(userID)

				return
			}
			// 其他错误，直接退出
			removeChatViewInRedis(userID)
			return
		}

		// 处理心跳
		heartB := HeartBeat{}
		if err = json.Unmarshal(data, &heartB); err != nil {
			errorMsg := []byte(`{"code": 0, "msg": "` + err.Error() + `"}`)
			node.Conn.WriteMessage(websocket.TextMessage, errorMsg)
			continue
		}

		if heartB.UserID != 0 {
			// 收到心跳包，重新设置读超时
			node.Conn.SetReadDeadline(time.Now().Add(readTimeout))
			continue
		}

		// 将消息分发给目标用户
		err = dispatch(data, userID)
		if err != nil {
			errorMsg := []byte(`{"code": 0, "msg": "` + err.Error() + `"}`)
			node.Conn.WriteMessage(websocket.TextMessage, errorMsg)
		} else {
			successMsg := []byte(`{"code": 1, "msg": "发送成功"}`)
			node.Conn.WriteMessage(websocket.TextMessage, successMsg)
		}
	}
}

// 将信息派发（考虑到后续群聊功能的实现因此实现dispatch，增加可扩展性）
func dispatch(data []byte, userID int) error {
	msg := model.ChatMsg{}
	err := json.Unmarshal(data, &msg)
	if msg.ChatMsgID != 0 {
		return errors.New("非法操作！")
	}
	if msg.SenderUserID != userID {
		return errors.New("非法操作！")
	}

	if err != nil {
		return err
	}
	// 添加聊天记录到数据库
	db := common.GetDB()
	db.Create(&msg)
	chatResp := ChatMsgResp{
		ChatMsgID:    msg.ChatMsgID,
		TargetUserID: msg.TargetUserID,
		SenderUserID: msg.SenderUserID,
		Content:      msg.Content,
		CreatedAt:    msg.CreatedAt,
	}

	data, err = json.Marshal(chatResp)
	if err != nil {
		return err
	}

	if msg.SenderUserID != msg.TargetUserID {
		// 发送消息
		sendMsg(msg.ChatMsgID, msg.TargetUserID, data)
	}
	return nil
}

type ChatMsgResp struct {
	ChatMsgID    int       `json:"chatMsgID"`
	TargetUserID int       `json:"targetUserID"`
	SenderUserID int       `json:"senderUserID"`
	Content      string    `json:"content"`
	Unread       int       `json:"unread"`
	CreatedAt    time.Time `json:"createdAt"`
}

// 发送消息到对应用户的消息队列
func sendMsg(chatMsgID, userID int, msg []byte) {
	chatMsg := ChatMsgResp{}
	//查看目标用户是否存在链接，不存在直接返回，否则发送消息
	val, ok := clientMap.Load(userID)
	//如果不存在表示一定是消息而不是相关用户的获得
	if err := json.Unmarshal(msg, &chatMsg); err != nil {
		return
	}
	if !ok {
		Chat_EmailAlerter(chatMsg.TargetUserID, "")
		redisCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		unreadKey := util.GenerateUnreadKey(userID, chatMsg.SenderUserID)
		core.MyRedis.Incr(redisCtx, unreadKey)
		core.MyRedis.SAdd(redisCtx, util.GenerateUnreadIndexKey(userID), unreadKey)
		return
	}
	node, ok := val.(*ChatNode)
	if ok {
		//如果消息接受者用户确实在ws链接中，判断消息接受者用户是否在聊天界面里
		//如果是聊天信息而不是返回聊天用户列表
		if chatMsgID != 0 {
			//检查接受者是否在聊天界面
			chatMsg.ChatMsgID = chatMsgID
			redisCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			err := core.MyRedis.Get(redisCtx, util.GenerateChatViewKey(chatMsg.TargetUserID, chatMsg.SenderUserID)).Err()
			if err == redis.Nil { //不在聊天界面,添加未读
				unreadKey := util.GenerateUnreadKey(chatMsg.TargetUserID, chatMsg.SenderUserID)
				resUnread := core.MyRedis.Incr(redisCtx, unreadKey)
				core.MyRedis.SAdd(redisCtx, util.GenerateUnreadIndexKey(chatMsg.TargetUserID), unreadKey)
				unread, _ := resUnread.Result()
				chatMsg.Unread = int(unread)
				Chat_EmailAlerter(chatMsg.TargetUserID, "")
			}
			chatMsgByte, _ := json.Marshal(chatMsg)
			node.DataQueue <- chatMsgByte
			return
		}
		//如果是聊天列表则返回
		node.DataQueue <- msg
	}
}

// 获得当前用户相关聊天的用户
func GetRelevantUser(userID int) ([]byte, error) {
	db := common.GetDB()
	chatMsgs := []model.ChatMsg{}
	err := db.Where("senderUserID = ? OR targetUserID = ?", userID, userID).
		Preload("TargetUser").
		Preload("SenderUser").
		Find(&chatMsgs).Error

	if err != nil {
		return []byte{}, err
	}
	type peerState struct {
		user        model.User
		lastAt      time.Time
		lastPreview string
	}
	peerStates := map[int]*peerState{}
	for _, msg := range chatMsgs {
		if msg.SenderUserID == msg.TargetUserID {
			continue
		}
		var peer model.User
		var peerID int
		if msg.SenderUserID == userID {
			peer = msg.TargetUser
			peerID = msg.TargetUserID
		} else {
			peer = msg.SenderUser
			peerID = msg.SenderUserID
		}
		if peerID == 0 || peerID == userID {
			continue
		}
		if st, ok := peerStates[peerID]; !ok || msg.CreatedAt.After(st.lastAt) {
			peerStates[peerID] = &peerState{
				user:        peer,
				lastAt:      msg.CreatedAt,
				lastPreview: msg.Content,
			}
		}
	}

	peerIDs := make([]int, 0, len(peerStates))
	for id := range peerStates {
		peerIDs = append(peerIDs, id)
	}
	sort.Slice(peerIDs, func(i, j int) bool {
		return peerStates[peerIDs[i]].lastAt.After(peerStates[peerIDs[j]].lastAt)
	})

	relevantRespUsers := make([]ChatRespUser, 0, len(peerIDs))
	for _, peerID := range peerIDs {
		st := peerStates[peerID]
		relevantUser := st.user
		redisCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		unreadStr, err := core.MyRedis.Get(redisCtx, util.GenerateUnreadKey(userID, relevantUser.UserID)).Result()
		cancel()
		unreadNum, _ := strconv.Atoi(unreadStr)
		if err != nil && err != redis.Nil {
			return nil, err
		}
		preview := st.lastPreview
		if len([]rune(preview)) > 80 {
			preview = string([]rune(preview)[:80]) + "…"
		}
		relevantRespUsers = append(relevantRespUsers, ChatRespUser{
			UserID:             relevantUser.UserID,
			Email:              relevantUser.Email,
			Name:               relevantUser.Name,
			AvatarURL:          relevantUser.AvatarURL,
			Identity:           relevantUser.Identity,
			Score:              relevantUser.Score,
			UnRead:             unreadNum,
			LastMessageAt:      st.lastAt.Format(time.RFC3339),
			LastMessagePreview: preview,
		})
	}
	jsonData, err := json.Marshal(struct {
		RelevantUsers []ChatRespUser `json:"relevantUsers"`
	}{relevantRespUsers})
	if err != nil {
		return []byte{}, err
	}
	return jsonData, nil
}

// 用户进入与某人聊天
//func IntoChatView(c *gin.Context) {
//
//	ctx := context.Background()
//	core.MyRedis.Set()
//}

// 用户退出与某人的聊天界面
func LeaveChatView(c *gin.Context) {
	chatQuery := ChatHistoryQuery{}
	if err := c.ShouldBind(&chatQuery); err != nil {
		response.Fail(c, gin.H{}, "无效的查询参数")
		return
	}
	//校验是否本人操作
	ctxUser, ok := c.Get("user")
	if !ok {
		response.Fail(c, gin.H{}, "操作失败，请重试")
		return
	}
	currentUser, _ := ctxUser.(model.User)
	if currentUser.UserID != chatQuery.SenderUserID {
		response.Fail(c, gin.H{}, "非法操作")
		return
	}

	//	删除redis存储的键值对
	redisCtx, cancel := context.WithTimeout(c.Request.Context(), 800*time.Millisecond)
	defer cancel()
	chatViewKey := util.GenerateChatViewKey(chatQuery.SenderUserID, chatQuery.TargetUserID)
	err := core.MyRedis.Del(redisCtx, chatViewKey).Err()
	_ = core.MyRedis.SRem(redisCtx, util.GenerateChatViewIndexKey(chatQuery.SenderUserID), chatViewKey).Err()
	if err != nil {
		response.Fail(c, gin.H{}, "操作异常")
	} else {
		response.Success(c, gin.H{}, "操作成功")
	}
}

func GetChatNotice(c *gin.Context) {
	userIDStr := c.Query("userID")
	userID, err := strconv.Atoi(userIDStr)
	userFromCtx, ok := c.Get("user")
	user := userFromCtx.(model.User)
	if err != nil || userID == 0 || !ok || user.UserID != userID {
		response.Fail(c, gin.H{}, "操作异常")
		return
	}

	var totalUnread int // 总的未读数

	// 使用未读索引集合，避免全库 SCAN
	redisCtx, cancel := context.WithTimeout(c.Request.Context(), 800*time.Millisecond)
	defer cancel()
	unreadKeys, err := core.MyRedis.SMembers(redisCtx, util.GenerateUnreadIndexKey(userID)).Result()
	if err != nil && err != redis.Nil {
		response.Fail(c, gin.H{}, "操作异常")
		return
	}
	if len(unreadKeys) > 0 {
		unreadVals, err := core.MyRedis.MGet(redisCtx, unreadKeys...).Result()
		if err != nil {
			response.Fail(c, gin.H{}, "获取未读消息失败")
			return
		}
		for i, v := range unreadVals {
			if v == nil {
				_ = core.MyRedis.SRem(redisCtx, util.GenerateUnreadIndexKey(userID), unreadKeys[i]).Err()
				continue
			}
			unreadStr, ok := v.(string)
			if !ok || unreadStr == "" {
				continue
			}
			unreadCount, convErr := strconv.Atoi(unreadStr)
			if convErr != nil {
				continue
			}
			totalUnread += unreadCount
		}
	}

	// 返回累加后的未读消息总数
	response.Success(c, gin.H{
		"noticeNum": totalUnread,
	}, "获取未读消息成功")
}

func Chat_EmailAlerter(receiverID int, msg string) {
	db := common.GetDB()
	var receiver model.User
	db.Where("userID = ?", receiverID).First(&receiver)
	if receiver.ISEmailNotificationBlocked || receiver.Email == "" {
		return
	}
	emailAddr, body := receiver.Email, msg
	go func() {
		if err := api.SendEmail(emailAddr, 2, body); err != nil {
			fmt.Println("async chat notify email:", err)
		}
	}()
}
