package main

import (
	"fmt"
	"log"
	"loginTest/common"
	"loginTest/config"
	"loginTest/controller"
	"loginTest/core"
	"loginTest/heat"
	"loginTest/middleware"
	"loginTest/route"
	"loginTest/service"
	"net/http"
	"net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	rpprof "runtime/pprof"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

// 核心变量
var (
	r            *gin.Engine
	pprofRunning bool
	pprofMutex   sync.Mutex
)

func Copy() {
	// 数据库连接信息
	dbHost := viper.GetString("datasource.host")
	dbPort := viper.GetInt("datasource.port")
	dbUser := viper.GetString("datasource.username")
	dbPassword := viper.GetString("datasource.password")
	dbName := viper.GetString("datasource.database")

	// 备份目录
	backupDir := "/app/database"

	c := cron.New()
	c.AddFunc("@every 12h", func() {
		backupFile := fmt.Sprintf("%s/backup_%s.sql", backupDir, time.Now().Format("2006-01-02 15:04:05"))
		cmd := exec.Command("mysqldump", fmt.Sprintf("-h%s", dbHost), fmt.Sprintf("-P%d", dbPort), fmt.Sprintf("-u%s", dbUser), fmt.Sprintf("-p%s", dbPassword), dbName, "--skip-ssl", "--result-file="+backupFile)
		err := cmd.Run()
		if err != nil {
			log.Println("备份失败:", err)
			return
		}
		log.Println("备份成功:", backupFile)
	})
	c.AddFunc("@every 1h", controller.CalculateAndSaveScores)
	c.Start()
}

// 手动触发pprof采集
func startPprof(c *gin.Context) {
	pprofMutex.Lock()
	if pprofRunning {
		pprofMutex.Unlock()
		c.JSON(400, gin.H{"msg": "pprof 正在采集中"})
		return
	}
	pprofRunning = true
	pprofMutex.Unlock()

	// 解析时长
	duration, _ := time.ParseDuration(c.DefaultQuery("duration", "30s"))
	saveDir := "/app/pprof"
	_ = os.MkdirAll(saveDir, 0755)
	now := time.Now().Format("20060102_150405")

	// 异步采集
	go func() {
		defer func() {
			pprofMutex.Lock()
			pprofRunning = false
			pprofMutex.Unlock()
			log.Println("✅ pprof 采集完成")
		}()

		// CPU采集（修复：使用别名rpprof）
		f, _ := os.Create(filepath.Join(saveDir, fmt.Sprintf("cpu_%s.pprof", now)))
		_ = rpprof.StartCPUProfile(f)
		time.Sleep(duration)
		rpprof.StopCPUProfile()
		_ = f.Close()

		// 内存采集（修复：使用别名rpprof）
		f2, _ := os.Create(filepath.Join(saveDir, fmt.Sprintf("heap_%s.pprof", now)))
		_ = rpprof.WriteHeapProfile(f2)
		_ = f2.Close()
	}()

	c.JSON(200, gin.H{
		"msg":      "✅ pprof 已启动",
		"duration": duration.String(),
		"savePath": "/app/pprof",
	})
}

func main() {
	config.InitConfig()
	config.InitEmailConfig()
	go Copy()
	db := common.InitDB()
	service.RebuildRankStats(db)
	rds := core.RedisInit()
	common.InitJWTkey()
	core.InitLocalCache()
	common.InitBloomFilter()
	service.InitHotDataService()

	defer rds.Close()
	defer db.Close()

	r = gin.Default()
	r.Use(middleware.LoggerToFile())
	r.StaticFS("/uploads", http.Dir("./public/uploads"))
	r.StaticFS("/resized", http.Dir("./public/resized"))

	debug := r.Group("/debug/pprof")
	{
		debug.GET("/", gin.WrapF(pprof.Index))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		debug.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
	}
	r.GET("/debug/pprof/start", startPprof)

	route.CollectRoute(r)
	go heat.StartHeat()

	srv := &http.Server{Addr: ":8080", Handler: r}
	log.Println("🚀 服务启动: 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
