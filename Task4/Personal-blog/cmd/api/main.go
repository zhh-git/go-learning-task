package main

import (
	"Personal-blog/configs"
	"Personal-blog/internal/api/router"
	"Personal-blog/internal/model"
	"Personal-blog/internal/pkg/db"
	"Personal-blog/internal/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 初始化配置
	if err := configs.Init(); err != nil {
		panic(fmt.Sprintf("配置初始化失败：%v", err))
	}

	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("日志初始化失败：%v", err))
	}

	// 3. 初始化数据库
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("数据库初始化失败：%v", err))
	}
	defer db.Close() // 程序退出时关闭数据库连接

	// 4. 自动迁移数据表（Gorm AutoMigrate）
	if err := db.DB.AutoMigrate(&model.User{}); err != nil {
		logger.Fatal("数据表迁移失败：", err)
	}
	if err := db.DB.AutoMigrate(&model.Post{}); err != nil {
		logger.Fatal("数据表迁移失败：", err)
	}
	if err := db.DB.AutoMigrate(&model.Comment{}); err != nil {
		logger.Fatal("数据表迁移失败：", err)
	}
	logger.Info("数据表迁移成功！")

	// 5. 初始化路由
	r := router.InitRouter()

	// 6. 启动服务（非阻塞方式）
	go func() {
		addr := fmt.Sprintf(":%d", configs.Config.App.Port)
		logger.Info(fmt.Sprintf("服务启动成功，监听地址：http://localhost%s", addr))
		if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务启动失败：", err)
		}
	}()

	// 7. 优雅关闭服务（监听信号）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 监听 Ctrl+C 和 kill 信号
	<-quit
	logger.Info("服务开始关闭...")

	// （可选）添加服务关闭前的清理逻辑（如关闭连接池、释放资源等）

	logger.Info("服务已正常关闭！")
}
