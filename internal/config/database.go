package config

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	adminModel "mqfm-backend/internal/models/auth/admin"
	"mqfm-backend/internal/utils"

)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("mqfm.db"), &gorm.Config{})
	if err != nil {
		utils.Log.Fatal(fmt.Sprintf("Database connection failed: %v", err))
	}

	database.AutoMigrate(&adminModel.Admin{})
	DB = database
}