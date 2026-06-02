package main

import "gorm.io/gorm"

func mustBootstrap(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database")
	}
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate users table")
	}
}
