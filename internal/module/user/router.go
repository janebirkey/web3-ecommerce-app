package user

import (
	"web3-ecommerce-app/internal/config"
	"web3-ecommerce-app/internal/middleware"
	"web3-ecommerce-app/internal/module/user/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户模块路由
func RegisterRoutes(router *gin.Engine, handler *handler.UserHTTPHandler, jwtConfig *config.JWTConfig) {
	// 创建v1版本API路由组
	v1 := router.Group("/api/v1")

	// 用户认证相关路由(不需要认证)
	authRoutes := v1.Group("/auth")
	{
		// 用户注册
		authRoutes.POST("/register", handler.Register)

		// 用户登录
		authRoutes.POST("/login", handler.Login)

		
	}

	// 用户相关路由(需要认证)
	userRoutes := v1.Group("/users")
	userRoutes.Use(middleware.JWT(jwtConfig))
	{
		// 获取当前用户信息
		userRoutes.GET("/profile", handler.GetProfile)

		// 获取指定用户信息
		userRoutes.GET("/:id", handler.GetUser)
	}
}
