package repository

import (
	"context"
	"web3-ecommerce-app/internal/domain/admin"
	"web3-ecommerce-app/internal/domain/user"
	userRepo "web3-ecommerce-app/internal/module/user/repository"

	"gorm.io/gorm"
)

// GormAdminRepository 是管理后台仓库的GORM实现
type GormAdminRepository struct {
	db *gorm.DB
}

// NewGormAdminRepository 创建一个新的GORM管理后台仓库
func NewGormAdminRepository(db *gorm.DB) admin.AdminRepository {
	return &GormAdminRepository{db: db}
}

// FindUsers 查找用户列表（带分页）
func (r *GormAdminRepository) FindUsers(ctx context.Context, param admin.PaginationParam) (*admin.UserPaginationResult, error) {
	var userModels []userRepo.UserModel
	var total int64

	// 计算分页参数
	offset := (param.Page - 1) * param.PageSize

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&userRepo.UserModel{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := r.db.WithContext(ctx).Offset(offset).Limit(param.PageSize).Find(&userModels).Error; err != nil {
		return nil, err
	}

	// 转换为领域模型
	users := make([]user.User, 0, len(userModels))
	for _, m := range userModels {
		users = append(users, user.User{
			ID:         m.ID,
			Username:   m.Username,
			Email:      m.Email,
			WalletAddr: m.WalletAddr,
			UserType:   m.UserType,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
		})
	}

	return &admin.UserPaginationResult{
		Total: int(total),
		Users: users,
	}, nil
}

// GetSystemOverview 获取系统概览数据
func (r *GormAdminRepository) GetSystemOverview(ctx context.Context) (*admin.SystemOverview, error) {
	var totalUsers, totalProducts, totalOrders, totalTransactions int64
	var totalSales float64
	var pendingOrders, pendingWithdrawals int64

	// 查询用户总数
	r.db.WithContext(ctx).Model(&userRepo.UserModel{}).Count(&totalUsers)

	// 其他统计数据查询
	// 注：由于其他模块未实现，这里暂时返回模拟数据

	return &admin.SystemOverview{
		TotalUsers:         int(totalUsers),
		TotalProducts:      int(totalProducts),
		TotalOrders:        int(totalOrders),
		TotalTransactions:  int(totalTransactions),
		TotalSales:         totalSales,
		PendingOrders:      int(pendingOrders),
		PendingWithdrawals: int(pendingWithdrawals),
	}, nil
}
