package service

import (
	"context"
	"fmt"
	"web3-ecommerce-app/internal/config"
	"web3-ecommerce-app/internal/domain/user"
	"web3-ecommerce-app/internal/middleware"
	"web3-ecommerce-app/pkg/apierror"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	// Register 注册用户
	Register(ctx context.Context, input user.CreateUserInput) (*user.UserOutput, error)

	// Login 登录用户
	Login(ctx context.Context, input user.LoginUserInput) (*user.UserOutput, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id uint) (*user.User, error)
}

// DefaultUserService 默认用户服务实现
type DefaultUserService struct {
	userRepo   user.UserRepository
	jwtConfig  *config.JWTConfig
	web3Config *config.Web3Config
}

// NewUserService 创建用户服务
func NewUserService(userRepo user.UserRepository, jwtConfig *config.JWTConfig, web3Config *config.Web3Config) UserService {
	return &DefaultUserService{
		userRepo:   userRepo,
		jwtConfig:  jwtConfig,
		web3Config: web3Config,
	}
}

// Register 注册用户
func (s *DefaultUserService) Register(ctx context.Context, input user.CreateUserInput) (*user.UserOutput, error) {
	// 检查用户是否已存在
	existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, apierror.NewDuplicateEntityError("邮箱已被注册", input.Email)
	}

	// 如果提供了钱包地址，检查是否已被使用
	if input.WalletAddr != "" {
		existingUserWallet, err := s.userRepo.FindByWalletAddr(ctx, input.WalletAddr)
		if err == nil && existingUserWallet != nil {
			return nil, apierror.NewDuplicateEntityError("钱包地址已被绑定", input.WalletAddr)
		}
	}

	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	newUser := &user.User{
		Username:   input.Username,
		Email:      input.Email,
		Password:   string(hashedPassword),
		WalletAddr: input.WalletAddr,
		UserType:   user.UserTypeRegular,
	}

	// 保存用户
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// 生成JWT令牌
	token, err := middleware.GenerateJWT(newUser.ID, newUser.UserType, newUser.WalletAddr, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 返回用户信息和令牌
	return &user.UserOutput{
		ID:         newUser.ID,
		Username:   newUser.Username,
		Email:      newUser.Email,
		WalletAddr: newUser.WalletAddr,
		UserType:   newUser.UserType,
		CreatedAt:  newUser.CreatedAt,
		Token:      token,
	}, nil
}

// Login 登录用户
func (s *DefaultUserService) Login(ctx context.Context, input user.LoginUserInput) (*user.UserOutput, error) {
	var userEntity *user.User
	var err error

	// 根据提供的登录方式处理
	if input.Email != "" && input.Password != "" {
		// 普通邮箱密码登录
		userEntity, err = s.userRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return nil, apierror.NewUnauthorizedError("登录失败", "邮箱或密码错误")
		}

		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(input.Password))
		if err != nil {
			return nil, apierror.NewUnauthorizedError("登录失败", "邮箱或密码错误")
		}
	} else if input.WalletAddr != "" && input.Message != "" && input.Signature != "" {
		// Web3登录

	} else {
		return nil, apierror.NewBadRequestError("登录失败", "请提供有效的登录凭证")
	}

	// 生成JWT令牌
	token, err := middleware.GenerateJWT(userEntity.ID, userEntity.UserType, userEntity.WalletAddr, s.jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	// 返回用户信息和令牌
	return &user.UserOutput{
		ID:         userEntity.ID,
		Username:   userEntity.Username,
		Email:      userEntity.Email,
		WalletAddr: userEntity.WalletAddr,
		UserType:   userEntity.UserType,
		CreatedAt:  userEntity.CreatedAt,
		Token:      token,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *DefaultUserService) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	return s.userRepo.FindByID(ctx, id)
}


