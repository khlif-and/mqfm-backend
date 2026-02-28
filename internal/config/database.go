package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	adminModel "mqfm-backend/internal/models/auth/admin"
	userModel "mqfm-backend/internal/models/auth/user"
	categoryAdminModel "mqfm-backend/internal/models/category/admin"
	historyModel "mqfm-backend/internal/models/history/user"
	likeModel "mqfm-backend/internal/models/likes/user"
	playlistModel "mqfm-backend/internal/models/playlist/user"
	audioAdminModel "mqfm-backend/internal/models/podcast/audio/admin"
	"mqfm-backend/internal/utils"

)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Log.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}

	database.AutoMigrate(
		&adminModel.Admin{},
		&userModel.User{},
		&categoryAdminModel.Category{},
		&audioAdminModel.Audio{},
		&playlistModel.Playlist{},
		&likeModel.Like{},
		&historyModel.History{},
	)
	DB = database
}