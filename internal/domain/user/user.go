package user

import (
	"context"
	"time"
)

// User 表示用户领域模型
type User struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"-"` // 密码不返回给客户端
	WalletAddr string    `json:"wallet_addr,omitempty"`
	UserType   string    `json:"user_type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UserType 用户类型常量
const (
	UserTypeRegular = "regular" // 普通用户
	UserTypeAdmin   = "admin"   // 管理员
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*User, error)

	// FindByEmail 根据Email查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByWalletAddr 根据钱包地址查找用户
	FindByWalletAddr(ctx context.Context, walletAddr string) (*User, error)

	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// Update 更新用户
	Update(ctx context.Context, user *User) error

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
}

// CreateUserInput 创建用户的输入参数
type CreateUserInput struct {
	Username   string `json:"username" binding:"required,min=3,max=50"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	WalletAddr string `json:"wallet_addr" binding:"omitempty,eth_addr"`
}

// LoginUserInput 登录用户的输入参数
type LoginUserInput struct {
	Email      string `json:"email" binding:"required_without=WalletAddr,omitempty,email"`
	Password   string `json:"password" binding:"required_without=Signature,omitempty,min=8"`
	WalletAddr string `json:"wallet_addr" binding:"required_without=Email,omitempty,eth_addr"`
	Message    string `json:"message" binding:"required_with=Signature"`
	Signature  string `json:"signature" binding:"required_with=Message"`
}

// UserOutput 用户输出
type UserOutput struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	WalletAddr string    `json:"wallet_addr,omitempty"`
	UserType   string    `json:"user_type"`
	CreatedAt  time.Time `json:"created_at"`
	Token      string    `json:"token,omitempty"` // JWT令牌
}

