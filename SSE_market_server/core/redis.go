package core

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var MyRedis *redis.Client

func RedisInit() *redis.Client {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	password := viper.GetString("redis.password")
	addr := fmt.Sprintf("%s:%s", host, port)

	rds := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
		// 连接池配置（高并发必备）
		PoolSize:     20,  // 最大连接数
		MinIdleConns: 5,   // 最小空闲连接
		// 超时配置（防止阻塞）
		DialTimeout:  500 * time.Millisecond,  // 连接超时
		ReadTimeout:  300 * time.Millisecond,  // 读超时
		WriteTimeout: 300 * time.Millisecond,  // 写超时
	})

	// 3. 新增：Redis健康检查（启动即校验连接，避免运行时报错）
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := rds.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("❌ Redis 连接失败：%v", err))
	}

	// 4. 赋值全局变量（保留原有逻辑）
	MyRedis = rds
	fmt.Println("✅ Redis 初始化成功（生产级配置）")
	return rds
}