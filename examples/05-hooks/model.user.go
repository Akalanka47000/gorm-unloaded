package main

import (
	"time"
)

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`
	Role  string

	CreatedAt time.Time
	UpdatedAt *time.Time
}
