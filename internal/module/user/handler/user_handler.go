package handler

import (
	"net/http"
	"strconv"
	"web3-ecommerce-app/internal/domain/user"
	"web3-ecommerce-app/internal/module/user/service"
	"web3-ecommerce-app/pkg/apierror"

	"github.com/gin-gonic/gin"
)

// UserHTTPHandler 用户HTTP处理器
type UserHTTPHandler struct {
	userService service.UserService
}

// NewUserHTTPHandler 创建用户HTTP处理器
func NewUserHTTPHandler(userService service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: userService,
	}
}

// Register 注册用户
func (h *UserHTTPHandler) Register(c *gin.Context) {
	var input user.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("注册失败", err.Error()),
		})
		return
	}

	output, err := h.userService.Register(c.Request.Context(), input)
	if err != nil {
		// 通用错误处理
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// Login 用户登录
func (h *UserHTTPHandler) Login(c *gin.Context) {
	var input user.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("登录失败", err.Error()),
		})
		return
	}

	output, err := h.userService.Login(c.Request.Context(), input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, output)
}


// GetUser 获取用户信息
func (h *UserHTTPHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewBadRequestError("无效的用户ID", err.Error()),
		})
		return
	}

	userEntity, err := h.userService.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          userEntity.ID,
		"username":    userEntity.Username,
		"email":       userEntity.Email,
		"wallet_addr": userEntity.WalletAddr,
		"user_type":   userEntity.UserType,
		"created_at":  userEntity.CreatedAt,
	})
}

// GetProfile 获取当前用户信息
func (h *UserHTTPHandler) GetProfile(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": apierror.NewUnauthorizedError("未授权", "请先登录"),
		})
		return
	}

	userEntity, err := h.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          userEntity.ID,
		"username":    userEntity.Username,
		"email":       userEntity.Email,
		"wallet_addr": userEntity.WalletAddr,
		"user_type":   userEntity.UserType,
		"created_at":  userEntity.CreatedAt,
	})
}

// handleError 处理错误
func (h *UserHTTPHandler) handleError(c *gin.Context, err error) {
	// 检查是否是API错误
	if apiErr, ok := err.(*apierror.APIError); ok {
		c.JSON(apiErr.Status, gin.H{"error": apiErr})
		return
	}

	// 默认为内部服务器错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": apierror.NewInternalServerError("服务器内部错误", err.Error()),
	})
}
