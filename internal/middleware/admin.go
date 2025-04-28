package middleware

import (
	"net/http"
	"web3-ecommerce-app/internal/domain/user"
	"web3-ecommerce-app/pkg/apierror"

	"github.com/gin-gonic/gin"
)

// AdminRequired 验证用户是否为管理员
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户类型
		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apierror.NewUnauthorizedError("未授权", "请先登录"),
			})
			return
		}

		// 验证用户是否为管理员
		if userType.(string) != user.UserTypeAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": apierror.NewForbiddenError("权限不足", "需要管理员权限"),
			})
			return
		}

		c.Next()
	}
}
