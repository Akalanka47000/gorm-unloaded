package main

import "gorm.io/gorm"

// This function is responsible for setting up the database schema and seeding initial data required for the example
func mustBootstrap(db *gorm.DB) {
	err := db.AutoMigrate(&City{})
	if err != nil {
		panic("failed to migrate database")
	}
	// Create a city
	err = db.Create(&City{Name: "Maharagama"}).Error
	if err != nil {
		panic("failed to create city")
	}
}
