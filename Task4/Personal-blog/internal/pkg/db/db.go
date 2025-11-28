package db

import (
	"Personal-blog/configs"
	"Personal-blog/internal/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接（Gorm）
func Init() error {
	// 配置 Gorm 日志（可选：生产环境可关闭）
	gormLogLevel := gormLogger.Silent
	if configs.Config.App.Env == "dev" {
		gormLogLevel = gormLogger.Info
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(configs.Config.Database.DSN), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		logger.Error("数据库连接失败：", err)
		return err
	}

	// 获取底层 sql.DB 连接池，配置连接参数
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取数据库连接池失败：", err)
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(configs.Config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(configs.Config.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(configs.Config.Database.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		logger.Error("数据库 ping 失败：", err)
		return err
	}

	DB = db
	logger.Info("数据库连接成功！")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
