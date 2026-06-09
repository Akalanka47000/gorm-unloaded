package main

import "time"

type Event struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	CreatedAt time.Time
}
