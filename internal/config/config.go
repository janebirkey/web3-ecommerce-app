package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Web3     Web3Config
}

type ServerConfig struct {
	Port    int
	Mode    string
	Timeout time.Duration
}

type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Username        string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type Web3Config struct {
	RPCURL      string `mapstructure:"rpc_url"`
	ChainID     int    `mapstructure:"chain_id"`
	NonceExpire time.Duration
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
