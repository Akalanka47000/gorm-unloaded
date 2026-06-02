package main

import "time"

type City struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Country    string
	Active     bool

	CreatedAt time.Time
	UpdatedAt *time.Time
}
