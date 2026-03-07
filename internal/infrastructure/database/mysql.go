package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/infrastructure/config"
	"mqfm-backend/internal/shared/logger"
)

func NewMySQL(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		logger.Fatal("database connection failed: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get underlying sql.DB: " + err.Error())
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	db.AutoMigrate(
		&entity.Admin{},
		&entity.User{},
		&entity.Category{},
		&entity.Audio{},
		&entity.Playlist{},
		&entity.Like{},
		&entity.History{},
		&entity.LiveStream{},
	)

	logger.Info("database connected successfully")
	return db
}
