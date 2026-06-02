package main

import "gorm.io/gorm"

func mustBootstrap(db *gorm.DB) {
	if err := db.AutoMigrate(&Order{}); err != nil {
		panic("failed to migrate database")
	}
	if err := db.Exec("TRUNCATE TABLE orders RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate orders table")
	}

	orders := []Order{
		{Reference: "ORD-001", Status: StatusPending},
		{Reference: "ORD-002", Status: StatusActive},
		{Reference: "ORD-003", Status: StatusActive},
		{Reference: "ORD-004", Status: StatusCancelled},
		{Reference: "ORD-005", Status: StatusPending},
	}
	if err := db.Create(&orders).Error; err != nil {
		panic("failed to seed orders")
	}
}
