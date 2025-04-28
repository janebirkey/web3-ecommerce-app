package admin

import (
	"context"
	"time"
	"web3-ecommerce-app/internal/domain/user"
)
// 建立模型主要是用于接收对应的参数
// PaginationParam 分页参数
type PaginationParam struct {
	Page     int
	PageSize int
}

// UserPaginationResult 用户分页结果
type UserPaginationResult struct {
	Total int         `json:"total"`
	Users []user.User `json:"users"`
}

// ProductFilter 产品过滤条件
// 我使用继承继承分页参数，然后添加对应的过滤条件
type ProductFilter struct {
	PaginationParam
	CategoryID uint
	Status     string
	Search     string
}

// OrderFilter 订单过滤条件
type OrderFilter struct {
	PaginationParam
	UserID    uint
	Status    string
	StartDate time.Time
	EndDate   time.Time
}

// TransactionFilter 交易过滤条件
type TransactionFilter struct {
	PaginationParam
	UserID    uint
	Type      string
	Status    string
	StartDate time.Time
	EndDate   time.Time
}

// SystemOverview 系统概览
type SystemOverview struct {
	TotalUsers         int     `json:"total_users"`
	TotalProducts      int     `json:"total_products"`
	TotalOrders        int     `json:"total_orders"`
	TotalTransactions  int     `json:"total_transactions"`
	TotalSales         float64 `json:"total_sales"`
	PendingOrders      int     `json:"pending_orders"`
	PendingWithdrawals int     `json:"pending_withdrawals"`
}

// AdminRepository 管理后台仓库接口
type AdminRepository interface {
	// 用户管理
	FindUsers(ctx context.Context, param PaginationParam) (*UserPaginationResult, error)

	// 统计数据
	GetSystemOverview(ctx context.Context) (*SystemOverview, error)
}
