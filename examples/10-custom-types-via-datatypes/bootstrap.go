package main

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func mustBootstrap(db *gorm.DB) {
	if err := db.AutoMigrate(&Product{}); err != nil {
		panic("failed to migrate database")
	}
	if err := db.Exec("TRUNCATE TABLE products RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate products table")
	}

	products := []Product{
		{
			Name:       "Wireless Headphones",
			Tags:       datatypes.NewJSONSlice([]string{"electronics", "audio", "sale"}),
			Attributes: datatypes.NewJSONType(ProductAttributes{Color: "black"}),
		},
		{
			Name:       "Leather Wallet",
			Tags:       datatypes.NewJSONSlice([]string{"accessories", "sale"}),
			Attributes: datatypes.NewJSONType(ProductAttributes{Color: "brown", Material: "leather"}),
		},
		{
			Name:       "Running Shoes",
			Tags:       datatypes.NewJSONSlice([]string{"footwear", "sports"}),
			Attributes: datatypes.NewJSONType(ProductAttributes{Color: "white", Size: "42"}),
		},
	}
	if err := db.Create(&products).Error; err != nil {
		panic("failed to seed products")
	}
}
