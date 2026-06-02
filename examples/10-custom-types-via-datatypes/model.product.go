package main

import (
	"time"

	"gorm.io/datatypes"
)

type ProductAttributes struct {
	Color    string `json:"color,omitempty"`
	Size     string `json:"size,omitempty"`
	Material string `json:"material,omitempty"`
}

type Product struct {
	ID   uint `gorm:"primaryKey"`
	Name string

	// Stored as JSONB — no custom Valuer/Scanner needed.
	Tags       datatypes.JSONSlice[string]               // ["electronics","sale"]
	Attributes datatypes.JSONType[ProductAttributes]     // {"color":"black","material":"leather"}

	CreatedAt time.Time
	UpdatedAt *time.Time
}
