package api

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/redis/go-redis/v9"
	"loginTest/core" 
	// 	"log"
	"loginTest/config"
)

func formVcode(ctx string) (string, string) {
	vcode := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			// 如果生成出错退回到时间戳拼接
			/*
			vcode += strconv.Itoa(int(time.Now().UnixNano()%10))
			continue
			*/
		}
		vcode += strconv.Itoa(int(n.Int64()))
	}
	ctx = strings.Replace(ctx, "{{vcode}}", vcode, -1)
	return ctx, vcode
}

func saveVcode(vcode, receiver string) {
	rds := core.MyRedis
	ctx := context.Background()
	codeKey := "vcode:" + receiver
	failKey := "vcode_fail:" + receiver
	_, err := rds.Get(ctx, codeKey).Result()
	if err == nil {
		rds.Del(ctx, codeKey)
	}
	// 重置失败次数
	rds.Del(ctx, failKey)
	rds.Set(ctx, codeKey, vcode, 5*time.Minute)
}

func SendEmail(receiver string, status int, msg string) error {
	e := email.NewEmail()
	senderString := config.Sender
	senderString = strings.Replace(senderString, "emailUsername", config.EmailUsername, -1)
	e.From = senderString
	if len([]rune(msg)) > 40 {
		msg = string([]rune(msg)[:40]) + "..."
	}
	e.To = []string{receiver}
	text := ""
	if status == 0 {
		rds := core.MyRedis
		ctx := context.Background()
		codeKey := "vcode:" + receiver
		_, err := rds.Get(ctx, codeKey).Result()
		// 如果 err == nil，说明找到了验证码（key存在），说明验证码还未过期
		// 如果 err == redis.Nil，说明 key 不存在，可以发送新验证码
		if err == nil {
			return errors.New("验证码已发送，请5分钟后再试")
		}
		// 如果是其他错误（如Redis连接错误），也阻止发送，避免异常情况
		if err != redis.Nil {
			return fmt.Errorf("检查验证码状态失败: %v", err)
		}
		e.Subject = config.ValidateTitle
		text = config.ValidateText
		text, vcode := formVcode(text)
		saveVcode(vcode, receiver)
		e.HTML = []byte(text)
	} else if status == 1 {
		e.Subject = config.CommentTitle
		text = strings.Replace(config.CommentText, "{{msg}}", msg, -1)
		e.HTML = []byte(text)
	} else if status == 2 {
		e.Subject = config.ChatTitle
		text = config.ChatText
		e.HTML = []byte(text)
	} else if status == 3 {
		e.Subject = config.ReplyTitle
		text = strings.Replace(config.ReplyText, "{{msg}}", msg, -1)
		e.HTML = []byte(text)
	}

	auth := smtp.PlainAuth("", config.EmailUsername, config.Password, config.Host)
	err := e.SendWithTLS(config.Addr, auth, &tls.Config{ServerName: config.Host})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Send Successfully")
	return nil
}
