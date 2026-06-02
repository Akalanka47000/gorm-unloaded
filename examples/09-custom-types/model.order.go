package main

import "time"

type Order struct {
	ID        uint   `gorm:"primaryKey"`
	Reference string `gorm:"uniqueIndex"`
	Status    Status

	CreatedAt time.Time
	UpdatedAt *time.Time
}
