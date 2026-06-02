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
		{Name: "Widget A", TenantID: "tenant-a"},
		{Name: "Widget B", TenantID: "tenant-a"},
		{Name: "Gadget X", TenantID: "tenant-b"},
		{Name: "Gadget Y", TenantID: "tenant-b"},
		{Name: "Gadget Z", TenantID: "tenant-b"},
	}
	if err := db.Create(&products).Error; err != nil {
		panic("failed to seed products")
	}
}
