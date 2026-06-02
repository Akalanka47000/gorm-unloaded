package main

import "time"

type Product struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Featured bool
	OnSale   bool

	CreatedAt time.Time
	UpdatedAt *time.Time
}
