package main

import "gorm.io/gorm"

func mustBootstrap(db *gorm.DB) {
	if err := db.AutoMigrate(&Product{}); err != nil {
		panic("failed to migrate database")
	}
	if err := db.Exec("TRUNCATE TABLE products RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate products table")
	}

	products := []Product{
		{Name: "Widget A", Featured: true, OnSale: false}, // featured only
		{Name: "Widget B", Featured: true, OnSale: true},  // featured AND on sale
		{Name: "Widget C", Featured: false, OnSale: true}, // on sale only
		{Name: "Widget D", Featured: false, OnSale: false}, // neither
	}
	if err := db.Create(&products).Error; err != nil {
		panic("failed to seed products")
	}
}
