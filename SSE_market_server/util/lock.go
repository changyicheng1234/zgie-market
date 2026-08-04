package util

import (
	"context"
	"loginTest/core"
	"time"
)

var ctx = context.Background()

// Lock 分布式锁结构体
type Lock struct {
	key    string
	expire time.Duration
}

// NewLock 创建锁（key=锁名称）
func NewLock(key string) *Lock {
	return &Lock{
		key:    key,
		expire: 5 * time.Second, // 锁过期时间，防止死锁
	}
}

// TryLock 尝试加锁（SETNX原子操作）
func (l *Lock) TryLock() bool {
	// 核心：Redis SETNX 命令，不存在则设置
	ok, err := core.MyRedis.SetNX(ctx, l.key, "1", l.expire).Result()
	return err == nil && ok
}

// UnLock 释放锁
func (l *Lock) UnLock() {
	core.MyRedis.Del(ctx, l.key)
}
