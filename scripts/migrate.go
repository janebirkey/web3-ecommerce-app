package main

import (
	"fmt"
	"log"
	"web3-ecommerce-app/internal/config"
	"web3-ecommerce-app/internal/module/user/repository"
	"web3-ecommerce-app/internal/platform/database"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewGormDB(&cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 创建用户仓库
	userRepo := repository.NewGormUserRepository(db)

	// 执行迁移
	userRepoWithMigration, ok := userRepo.(*repository.GormUserRepository)
	if !ok {
		log.Fatalf("无法转换仓库类型")
	}

	if err := userRepoWithMigration.AutoMigrate(); err != nil {
		log.Fatalf("迁移失败: %v", err)
	}

	fmt.Println("数据库迁移成功完成")
}
