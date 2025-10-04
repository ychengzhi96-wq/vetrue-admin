package database

import (
	"fmt"
	"log"
	"time"
	"vetrue-vben-api/internal/config"

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

	// 使用 Go 1.25 优化的 switch 表达式
	dialector, err := getDialector(dbConfig.Driver, dbConfig.DSN)
	if err != nil {
		return err
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

	var openErr error
	DB, openErr = gorm.Open(dialector, gormConfig)
	if openErr != nil {
		return fmt.Errorf("连接数据库失败: %w", openErr)
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

// getDialector 根据数据库类型返回对应的 dialector
func getDialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "mysql":
		return mysql.Open(dsn), nil
	case "postgres":
		return postgres.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	case "sqlserver":
		return sqlserver.Open(dsn), nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", driver)
	}
}
