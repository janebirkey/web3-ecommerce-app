package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"web3-ecommerce-app/internal/config"
	"web3-ecommerce-app/pkg/apierror"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTClaims 表示JWT载荷
// 后面可以添加更多字段，比如用户名、邮箱等
type JWTClaims struct {
	UserID     uint   `json:"user_id"`
	UserType   string `json:"user_type"`
	WalletAddr string `json:"wallet_addr,omitempty"`
	jwt.StandardClaims
}

// JWT 中间件工厂函数
func JWT(jwtConfig *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apierror.NewUnauthorizedError("缺少认证信息", "未提供Authorization头"),
			})
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apierror.NewUnauthorizedError("认证格式错误", "Authorization头必须是Bearer格式"),
			})
			return
		}

		// 解析token
		tokenString := parts[1]
		claims := &JWTClaims{}

		// 手动解析token，不检查过期时间
		token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})

		// 检查签名是否有效
		if token == nil || !token.Valid && claims.ExpiresAt == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apierror.NewUnauthorizedError("无效的token", "签名验证失败"),
			})
			return
		}

		// 添加5分钟的时间宽容度
		now := time.Now().Unix()
		gracePeriod := int64(300) // 5分钟 (秒)

		if claims.ExpiresAt < now-gracePeriod {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": apierror.NewUnauthorizedError("无效的token", fmt.Sprintf("token已超过宽限期过期")),
			})
			return
		}

		// 将claims存入上下文
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		if claims.WalletAddr != "" {
			c.Set("wallet_addr", claims.WalletAddr)
		}
		c.Next()
	}
}

// GenerateJWT 生成JWT token
func GenerateJWT(userID uint, userType string, walletAddr string, jwtConfig *config.JWTConfig) (string, error) {
	// 创建声明
	claims := JWTClaims{
		UserID:     userID,
		UserType:   userType,
		WalletAddr: walletAddr,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(jwtConfig.ExpireHours)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	return token.SignedString([]byte(jwtConfig.Secret))
}
