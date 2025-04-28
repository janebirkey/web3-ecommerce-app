package repository

// 这一层是仓库层
import (
	"context"
	"fmt"
	"web3-ecommerce-app/internal/domain/user"
	"web3-ecommerce-app/pkg/apierror"

	"gorm.io/gorm"
)

// UserModel 是GORM用户模型
type UserModel struct {
	gorm.Model
	Username   string `gorm:"type:varchar(50);not null;uniqueIndex:idx_username"`
	Email      string `gorm:"type:varchar(100);not null;uniqueIndex:idx_email"`
	Password   string `gorm:"type:varchar(100);not null"`
	WalletAddr string `gorm:"type:varchar(42);uniqueIndex:idx_wallet_addr"`
	UserType   string `gorm:"type:varchar(20);not null;default:'regular'"`
}

// TableName 指定表名
func (UserModel) TableName() string {
	return "users"
}

// GormUserRepository 是用户仓库的GORM实现
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建一个新的GORM用户仓库
func NewGormUserRepository(db *gorm.DB) user.UserRepository {
	return &GormUserRepository{db: db}
}

// domainToModel 将领域模型转换为GORM模型
func domainToModel(u *user.User) *UserModel {
	return &UserModel{
		Model: gorm.Model{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Username:   u.Username,
		Email:      u.Email,
		Password:   u.Password,
		WalletAddr: u.WalletAddr,
		UserType:   u.UserType,
	}
}

// modelToDomain 将GORM模型转换为领域模型
func modelToDomain(m *UserModel) *user.User {
	return &user.User{
		ID:         m.ID,
		Username:   m.Username,
		Email:      m.Email,
		Password:   m.Password,
		WalletAddr: m.WalletAddr,
		UserType:   m.UserType,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// FindByID 根据ID查找用户
func (r *GormUserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewNotFoundError("用户不存在", fmt.Sprintf("ID: %d", id))
		}
		return nil, fmt.Errorf("查询用户错误: %w", err)
	}
	return modelToDomain(&model), nil
}

// FindByEmail 根据Email查找用户
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewNotFoundError("用户不存在", fmt.Sprintf("Email: %s", email))
		}
		return nil, fmt.Errorf("查询用户错误: %w", err)
	}
	return modelToDomain(&model), nil
}

// FindByWalletAddr 根据钱包地址查找用户
func (r *GormUserRepository) FindByWalletAddr(ctx context.Context, walletAddr string) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("wallet_addr = ?", walletAddr).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierror.NewNotFoundError("用户不存在", fmt.Sprintf("钱包地址: %s", walletAddr))
		}
		return nil, fmt.Errorf("查询用户错误: %w", err)
	}
	return modelToDomain(&model), nil
}

// Create 创建用户
func (r *GormUserRepository) Create(ctx context.Context, u *user.User) error {
	model := domainToModel(u)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if r.db.WithContext(ctx).Where("email = ?", u.Email).First(&UserModel{}).Error == nil {
			return apierror.NewDuplicateEntityError("邮箱已被使用", u.Email)
		}
		if u.WalletAddr != "" && r.db.WithContext(ctx).Where("wallet_addr = ?", u.WalletAddr).First(&UserModel{}).Error == nil {
			return apierror.NewDuplicateEntityError("钱包地址已被绑定", u.WalletAddr)
		}
		if r.db.WithContext(ctx).Where("username = ?", u.Username).First(&UserModel{}).Error == nil {
			return apierror.NewDuplicateEntityError("用户名已被使用", u.Username)
		}
		return fmt.Errorf("创建用户错误: %w", err)
	}

	// 更新领域模型
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt

	return nil
}

// Update 更新用户
func (r *GormUserRepository) Update(ctx context.Context, u *user.User) error {
	model := domainToModel(u)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("更新用户错误: %w", err)
	}

	// 更新领域模型
	u.UpdatedAt = model.UpdatedAt

	return nil
}

// Delete 删除用户
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&UserModel{}, id).Error; err != nil {
		return fmt.Errorf("删除用户错误: %w", err)
	}
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func (r *GormUserRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&UserModel{})
}
