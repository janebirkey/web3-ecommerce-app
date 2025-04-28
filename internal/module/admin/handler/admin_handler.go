package handler

import (
	"net/http"
	"strconv"
	"web3-ecommerce-app/internal/domain/admin"
	"web3-ecommerce-app/internal/module/admin/service"
	"web3-ecommerce-app/pkg/apierror"

	"github.com/gin-gonic/gin"
)

// AdminHTTPHandler 管理后台HTTP处理器
type AdminHTTPHandler struct {
	adminService service.AdminService
}

// NewAdminHTTPHandler 创建管理后台HTTP处理器
func NewAdminHTTPHandler(adminService service.AdminService) *AdminHTTPHandler {
	return &AdminHTTPHandler{
		adminService: adminService,
	}
}

// getIDFromParam 从URL参数中获取ID
func getIDFromParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, apierror.NewBadRequestError("无效的ID", err.Error())
	}
	return uint(id), nil
}

// handleError 处理错误
func handleError(c *gin.Context, err error) {
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

// 用户管理
// ListUsers 获取用户列表
func (h *AdminHTTPHandler) ListUsers(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	result, err := h.adminService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUser 获取用户详情
func (h *AdminHTTPHandler) GetUser(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	user, err := h.adminService.GetUser(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户信息
func (h *AdminHTTPHandler) UpdateUser(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var userData map[string]interface{}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", err.Error()),
		})
		return
	}

	user, err := h.adminService.UpdateUser(c.Request.Context(), id, userData)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser 删除用户
func (h *AdminHTTPHandler) DeleteUser(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.adminService.DeleteUser(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// 产品管理
// CreateProduct 创建产品
func (h *AdminHTTPHandler) CreateProduct(c *gin.Context) {
	var productData map[string]interface{}
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", err.Error()),
		})
		return
	}

	product, err := h.adminService.CreateProduct(c.Request.Context(), productData)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// ListProducts 获取产品列表
func (h *AdminHTTPHandler) ListProducts(c *gin.Context) {
	// 获取过滤参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.ParseUint(c.DefaultQuery("category_id", "0"), 10, 32)
	status := c.DefaultQuery("status", "")
	search := c.DefaultQuery("search", "")

	filter := admin.ProductFilter{
		PaginationParam: admin.PaginationParam{
			Page:     page,
			PageSize: pageSize,
		},
		CategoryID: uint(categoryID),
		Status:     status,
		Search:     search,
	}

	result, err := h.adminService.ListProducts(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProduct 获取产品详情
func (h *AdminHTTPHandler) GetProduct(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	product, err := h.adminService.GetProduct(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct 更新产品
func (h *AdminHTTPHandler) UpdateProduct(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var productData map[string]interface{}
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", err.Error()),
		})
		return
	}

	product, err := h.adminService.UpdateProduct(c.Request.Context(), id, productData)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct 删除产品
func (h *AdminHTTPHandler) DeleteProduct(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.adminService.DeleteProduct(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品删除成功"})
}

// UpdateProductStatus 更新产品状态
func (h *AdminHTTPHandler) UpdateProductStatus(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var statusData map[string]string
	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", err.Error()),
		})
		return
	}

	status, ok := statusData["status"]
	if !ok || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", "状态不能为空"),
		})
		return
	}

	if err := h.adminService.UpdateProductStatus(c.Request.Context(), id, status); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "产品状态更新成功"})
}

// 订单管理
// ListOrders 获取订单列表
func (h *AdminHTTPHandler) ListOrders(c *gin.Context) {
	// 获取过滤参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID, _ := strconv.ParseUint(c.DefaultQuery("user_id", "0"), 10, 32)
	status := c.DefaultQuery("status", "")

	// 时间参数可能需要更复杂的处理，这里简化处理
	filter := admin.OrderFilter{
		PaginationParam: admin.PaginationParam{
			Page:     page,
			PageSize: pageSize,
		},
		UserID: uint(userID),
		Status: status,
	}

	result, err := h.adminService.ListOrders(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetOrder 获取订单详情
func (h *AdminHTTPHandler) GetOrder(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	order, err := h.adminService.GetOrder(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus 更新订单状态
func (h *AdminHTTPHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var statusData map[string]string
	if err := c.ShouldBindJSON(&statusData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", err.Error()),
		})
		return
	}

	status, ok := statusData["status"]
	if !ok || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": apierror.NewValidationError("无效的请求数据", "状态不能为空"),
		})
		return
	}

	if err := h.adminService.UpdateOrderStatus(c.Request.Context(), id, status); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单状态更新成功"})
}

// 支付管理
// ListTransactions 获取交易列表
func (h *AdminHTTPHandler) ListTransactions(c *gin.Context) {
	// 获取过滤参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	userID, _ := strconv.ParseUint(c.DefaultQuery("user_id", "0"), 10, 32)
	txType := c.DefaultQuery("type", "")
	status := c.DefaultQuery("status", "")

	filter := admin.TransactionFilter{
		PaginationParam: admin.PaginationParam{
			Page:     page,
			PageSize: pageSize,
		},
		UserID: uint(userID),
		Type:   txType,
		Status: status,
	}

	result, err := h.adminService.ListTransactions(c.Request.Context(), filter)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ProcessWithdrawal 处理提现
func (h *AdminHTTPHandler) ProcessWithdrawal(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.adminService.ProcessWithdrawal(c.Request.Context(), id); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "提现处理成功"})
}

// 统计数据
// GetSystemOverview 获取系统概览
func (h *AdminHTTPHandler) GetSystemOverview(c *gin.Context) {
	result, err := h.adminService.GetSystemOverview(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
