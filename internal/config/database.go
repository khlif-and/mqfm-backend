package config

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	adminModel "mqfm-backend/internal/models/auth/admin"
	userModel "mqfm-backend/internal/models/auth/user"
	categoryAdminModel "mqfm-backend/internal/models/category/admin"
	"mqfm-backend/internal/utils"

)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("mqfm.db"), &gorm.Config{})
	if err != nil {
		utils.Log.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}

	// Tambahkan categoryAdminModel.Category{} ke AutoMigrate
	database.AutoMigrate(
		&adminModel.Admin{},
		&userModel.User{},
		&categoryAdminModel.Category{},
	)
	DB = database
}