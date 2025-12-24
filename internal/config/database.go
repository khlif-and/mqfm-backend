package config

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	adminModel "mqfm-backend/internal/models/auth/admin"
	userModel "mqfm-backend/internal/models/auth/user"
	categoryAdminModel "mqfm-backend/internal/models/category/admin"
	audioAdminModel "mqfm-backend/internal/models/podcast/audio/admin"
	"mqfm-backend/internal/utils"

)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("mqfm.db"), &gorm.Config{})
	if err != nil {
		utils.Log.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}

	database.AutoMigrate(
		&adminModel.Admin{},
		&userModel.User{},
		&categoryAdminModel.Category{},
		&audioAdminModel.Audio{},
	)
	DB = database
}