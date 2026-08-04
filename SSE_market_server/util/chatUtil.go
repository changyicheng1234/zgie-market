package util

import "fmt"

func GenerateChatViewKey(senderUserID, targetUserID int) string {
	return fmt.Sprintf("%dTo%dIsExist", senderUserID, targetUserID)
}

func GenerateUnreadKey(senderUserID, targetUserID int) string {
	return fmt.Sprintf("%dRecvUnread%d", senderUserID, targetUserID)
}

func GenerateChatViewIndexKey(userID int) string {
	return fmt.Sprintf("chatViewIndex:%d", userID)
}

func GenerateUnreadIndexKey(userID int) string {
	return fmt.Sprintf("chatUnreadIndex:%d", userID)
}
