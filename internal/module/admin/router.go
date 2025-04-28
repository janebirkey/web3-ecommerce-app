package admin

import (
	"web3-ecommerce-app/internal/config"
	"web3-ecommerce-app/internal/middleware"
	"web3-ecommerce-app/internal/module/admin/handler"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册管理后台模块路由
func RegisterRoutes(
	router *gin.Engine,
	adminHandler *handler.AdminHTTPHandler,
	jwtConfig *config.JWTConfig,
) {
	// 创建管理后台API路由组
	adminRoutes := router.Group("/api/v1/admin")

	// 管理后台需要JWT认证和管理员权限验证
	adminRoutes.Use(middleware.JWT(jwtConfig))
	adminRoutes.Use(middleware.AdminRequired())

	// 用户管理
	{
		// 获取用户列表
		adminRoutes.GET("/users", adminHandler.ListUsers)

		// 获取单个用户详情
		adminRoutes.GET("/users/:id", adminHandler.GetUser)

		// 更新用户信息
		adminRoutes.PUT("/users/:id", adminHandler.UpdateUser)

		// 删除用户
		adminRoutes.DELETE("/users/:id", adminHandler.DeleteUser)
	}

	// 产品管理
	{
		// 创建产品
		adminRoutes.POST("/products", adminHandler.CreateProduct)

		// 获取产品列表
		adminRoutes.GET("/products", adminHandler.ListProducts)

		// 获取单个产品详情
		adminRoutes.GET("/products/:id", adminHandler.GetProduct)

		// 更新产品
		adminRoutes.PUT("/products/:id", adminHandler.UpdateProduct)

		// 删除产品
		adminRoutes.DELETE("/products/:id", adminHandler.DeleteProduct)

		// 更改产品状态
		adminRoutes.PATCH("/products/:id/status", adminHandler.UpdateProductStatus)
	}

	// 订单管理
	{
		// 获取订单列表
		adminRoutes.GET("/orders", adminHandler.ListOrders)

		// 获取单个订单详情
		adminRoutes.GET("/orders/:id", adminHandler.GetOrder)

		// 更新订单状态
		adminRoutes.PATCH("/orders/:id/status", adminHandler.UpdateOrderStatus)
	}

	// 支付管理
	{
		// 获取交易列表
		adminRoutes.GET("/transactions", adminHandler.ListTransactions)

		// 处理提现请求
		adminRoutes.POST("/withdrawals/:id/process", adminHandler.ProcessWithdrawal)
	}

	// 统计数据
	{
		// 获取系统概览统计
		adminRoutes.GET("/stats/overview", adminHandler.GetSystemOverview)
	}
}
