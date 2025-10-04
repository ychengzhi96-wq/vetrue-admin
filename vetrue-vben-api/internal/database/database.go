package database

import (
	"fmt"
	"vetrue-vben-api/internal/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() error {
	if len(config.AppConfig.Databases) == 0 {
		return fmt.Errorf("未配置数据库")
	}

	// 使用第一个数据库作为主数据库
	dbConfig := config.AppConfig.Databases[0]

	var dialector gorm.Dialector

	switch dbConfig.Driver {
	case "mysql":
		dialector = mysql.Open(dbConfig.DSN)
	case "postgres":
		dialector = postgres.Open(dbConfig.DSN)
	case "sqlite":
		dialector = sqlite.Open(dbConfig.DSN)
	case "sqlserver":
		dialector = sqlserver.Open(dbConfig.DSN)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbConfig.Driver)
	}

	// 配置日志级别
	logLevel := logger.Silent
	switch dbConfig.LogLevel {
	case 1:
		logLevel = logger.Silent
	case 2:
		logLevel = logger.Error
	case 3:
		logLevel = logger.Warn
	case 4:
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var err error
	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	if dbConfig.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	}
	if dbConfig.MaxConn > 0 {
		sqlDB.SetMaxOpenConns(dbConfig.MaxConn)
	}
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("数据库连接成功: %s", dbConfig.Name)
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
