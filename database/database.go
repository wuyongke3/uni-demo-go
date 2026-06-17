package database

import (
	"fmt"
	"log"
	"unigo/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库连接并自动迁移
func Init(cfg *config.DatabaseConfig) error {
	var err error
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN())
	case "postgresql":
		dialector = postgres.Open(cfg.DSN())
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.Printf("数据库连接成功 (%s)", cfg.Driver)
	return nil
}
