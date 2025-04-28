package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web3-ecommerce-app/internal/config"
	"web3-ecommerce-app/internal/module/admin"
	adminHandler "web3-ecommerce-app/internal/module/admin/handler"
	adminRepo "web3-ecommerce-app/internal/module/admin/repository"
	adminService "web3-ecommerce-app/internal/module/admin/service"
	"web3-ecommerce-app/internal/module/user"
	"web3-ecommerce-app/internal/module/user/handler"
	"web3-ecommerce-app/internal/module/user/repository"
	"web3-ecommerce-app/internal/module/user/service"
	"web3-ecommerce-app/internal/platform/database"
	"web3-ecommerce-app/internal/platform/httprouter"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf(`无法加载配置: %v`, err)
	}

	// 初始化数据库连接
	db, err := database.NewGormDB(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化仓库
	userRepo := repository.NewGormUserRepository(db)
	adminRepo := adminRepo.NewGormAdminRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo, &cfg.JWT, &cfg.Web3)

	// 初始化管理后台服务
	adminSvc := adminService.NewAdminService(adminRepo, userRepo, userService)

	// 初始化处理器
	userHandler := handler.NewUserHTTPHandler(userService)
	adminHandler := adminHandler.NewAdminHTTPHandler(adminSvc)

	// 初始化HTTP路由器,创建对应的gin引擎
	router := httprouter.NewGinEngine(&cfg.Server)

	// 注册路由
	user.RegisterRoutes(router, userHandler, &cfg.JWT)
	admin.RegisterRoutes(router, adminHandler, &cfg.JWT)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  120 * time.Second,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("服务器正在监听端口 %d\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 是 syscall.SIGINT
	// kill -9 是 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("服务器被强制关闭:", err)
	}

	log.Println("服务器已优雅关闭")
}
