package main

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// This is the key field that enables soft delete functionality in GORM.
	// When a record is soft deleted, this field will be set to the current timestamp,
	// indicating that the record is deleted without actually removing it from the database.
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
