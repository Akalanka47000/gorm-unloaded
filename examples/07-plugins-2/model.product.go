package main

import "time"

type Product struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	TenantID string

	CreatedAt time.Time
	UpdatedAt *time.Time
}
