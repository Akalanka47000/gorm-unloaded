package main

import "gorm.io/gorm"

func mustBootstrap(db *gorm.DB) {
	if err := db.AutoMigrate(&Event{}); err != nil {
		panic("failed to migrate database")
	}
	if err := db.Exec("TRUNCATE TABLE events RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate events table")
	}
}
