package main

import (
	"fmt"
	"log"
	"unigo/config"
	"unigo/database"
	"unigo/repository"
	"unigo/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置 (YAML 文件)
	cfg := config.Load()

	// 打印配置信息
	log.Println("========================================")
	log.Println("       UniGo 服务启动中...")
	log.Println("========================================")
	fmt.Printf("  [Server]   端口: %s\n", cfg.Server.Port)
	fmt.Printf("  [Database] 驱动: %s | 地址: %s:%d | 库: %s\n",
		cfg.Database.Driver, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Println("========================================")

	// 2. 初始化数据库连接
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("[FATAL] 数据库连接失败: %v", err)
	}

	// 3. 自动迁移表结构 (检查并创建不存在的表)
	if err := repository.AutoMigrate(); err != nil {
		log.Fatalf("[FATAL] 数据库迁移失败: %v", err)
	}

	// 4. 创建 Gin 引擎
	r := gin.Default()

	// 5. 注册路由 (传入 JWT 配置)
	router.SetupRouter(r, cfg.JWT)

	// 6. 启动服务
	addr := ":" + cfg.Server.Port
	fmt.Printf("\n  服务地址: http://localhost%s\n", addr)
	fmt.Printf("  API 文档: http://localhost%s/api/v1/lecturers/info/1\n", addr)
	log.Println("========================================")
	log.Printf("服务启动成功, 监听端口: %s", cfg.Server.Port)

	r.Run(addr)
}
