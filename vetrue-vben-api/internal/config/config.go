package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	App       App              `mapstructure:"app"`
	Databases []DatabaseConfig `mapstructure:"databases"`
	Redis     RedisConfig      `mapstructure:"redis"`
	Mongo     MongoConfig      `mapstructure:"mongo"`
	Log       LogConfig        `mapstructure:"log"`
	JWT       JWTConfig        `mapstructure:"jwt"`
	Casbin    CasbinConfig     `mapstructure:"casbin"`
	OSS       OSSConfig        `mapstructure:"oss"`
	Telegram  TelegramConfig   `mapstructure:"telegram"`
}

type App struct {
	Name         string          `mapstructure:"name"`
	Env          string          `mapstructure:"env"`
	Addr         string          `mapstructure:"addr"`
	Timeout      int             `mapstructure:"timeout"`
	RouterPrefix string          `mapstructure:"routerPrefix"`
	RateLimit    RateLimitConfig `mapstructure:"rateLimit"`
}

type RateLimitConfig struct {
	Enabled  bool `mapstructure:"enabled"`
	Rate     int  `mapstructure:"rate"`     // 每秒请求数
	Capacity int  `mapstructure:"capacity"` // 令牌桶容量
}

type DatabaseConfig struct {
	Name            string `mapstructure:"name"`
	Driver          string `mapstructure:"driver"`
	DSN             string `mapstructure:"dsn"`
	UseGorm         bool   `mapstructure:"useGorm"`
	Remark          string `mapstructure:"remark"`
	LogLevel        int    `mapstructure:"logLevel"`
	EnableLogWriter bool   `mapstructure:"enableLogWriter"`
	MaxIdleConn     int    `mapstructure:"maxIdleConn"`
	MaxConn         int    `mapstructure:"maxConn"`
	SlowThreshold   int    `mapstructure:"slowThreshold"`
}

type RedisConfig struct {
	Addr      string `mapstructure:"addr"`
	Password  string `mapstructure:"password"`
	DB        int    `mapstructure:"db"`
	IsCluster bool   `mapstructure:"isCluster"`
}

type MongoConfig struct {
	URL string `mapstructure:"url"`
}

type LogConfig struct {
	Remark     string `mapstructure:"remark"`
	Path       string `mapstructure:"path"`
	Mode       string `mapstructure:"mode"`
	LogRotate  bool   `mapstructure:"logrotate"`
	Recover    bool   `mapstructure:"recover"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}

type JWTConfig struct {
	SecretKey string `mapstructure:"secretKey"`
	Expire    int64  `mapstructure:"expire"`
}

type CasbinConfig struct {
	ModePath string `mapstructure:"modePath"`
	Remark   string `mapstructure:"remark"`
	DbName   string `mapstructure:"dbName"`
}

type OSSConfig struct {
	Type       string `mapstructure:"type"`
	SavePath   string `mapstructure:"savePath"`
	URL        string `mapstructure:"url"`
	AccessKey  string `mapstructure:"accessKey"`
	SecretKey  string `mapstructure:"secretKey"`
	BucketName string `mapstructure:"bucketName"`
}

type TelegramConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	BotToken      string `mapstructure:"botToken"`
	BotUsername   string `mapstructure:"botUsername"`
	MiniAppURL    string `mapstructure:"miniAppUrl"`
	SetGlobalMenu bool   `mapstructure:"setGlobalMenu"`
	Debug         bool   `mapstructure:"debug"`
}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}
