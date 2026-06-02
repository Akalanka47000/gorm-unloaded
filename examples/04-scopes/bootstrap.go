package main

import "gorm.io/gorm"

func mustBootstrap(db *gorm.DB) {
	err := db.AutoMigrate(&City{})
	if err != nil {
		panic("failed to migrate database")
	}

	if err = db.Exec("TRUNCATE TABLE cities RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate cities table")
	}

	cities := []City{
		{Name: "Colombo", Country: "Sri Lanka", Active: true},
		{Name: "Kandy", Country: "Sri Lanka", Active: true},
		{Name: "Galle", Country: "Sri Lanka", Active: false},
		{Name: "Maharagama", Country: "Sri Lanka", Active: true},
		{Name: "Nugegoda", Country: "Sri Lanka", Active: true},
		{Name: "Negombo", Country: "Sri Lanka", Active: true},
		{Name: "Mumbai", Country: "India", Active: true},
		{Name: "Delhi", Country: "India", Active: true},
		{Name: "Chennai", Country: "India", Active: true},
		{Name: "London", Country: "UK", Active: true},
		{Name: "Manchester", Country: "UK", Active: false},
	}

	err = db.Create(&cities).Error
	if err != nil {
		panic("failed to seed cities")
	}
}
