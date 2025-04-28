package httprouter

import (
	"web3-ecommerce-app/internal/config"

	"github.com/gin-gonic/gin"
)

// NewGinEngine 创建一个新的Gin引擎
func NewGinEngine(cfg *config.ServerConfig) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Mode)

	// 创建默认引擎，包含Logger和Recovery中间件
	router := gin.Default()

	// 设置信任的代理
	router.SetTrustedProxies(nil)

	return router
}

// 注册路由组，并应用通用中间件
func NewRouterGroup(router *gin.Engine, path string, middlewares ...gin.HandlerFunc) *gin.RouterGroup {
	group := router.Group(path)

	if len(middlewares) > 0 {
		group.Use(middlewares...)
	}

	return group
}
