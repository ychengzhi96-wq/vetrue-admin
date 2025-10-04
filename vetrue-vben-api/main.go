package main

import (
	"context"
	"vetrue-vben-api/internal/config"
	"vetrue-vben-api/internal/database"
	"vetrue-vben-api/internal/router"
	"vetrue-vben-api/pkg/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化Casbin
	if err := middleware.InitCasbin(); err != nil {
		log.Fatalf("初始化Casbin失败: %v", err)
	}

	// 初始化路由
	r := router.InitRouter()

	// 启动HTTP服务器
	addr := config.AppConfig.App.Addr
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(config.AppConfig.App.Timeout) * time.Second,
		WriteTimeout: time.Duration(config.AppConfig.App.Timeout) * time.Second,
	}

	// 优雅启动
	go func() {
		log.Printf("服务器启动在: %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭异常: %v", err)
	}

	log.Println("服务器已关闭")
}
