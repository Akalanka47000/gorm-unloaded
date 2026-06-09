package main

import "gorm.io/gorm"

func mustBootstrap(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database")
	}
	if err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY").Error; err != nil {
		panic("failed to truncate users table")
	}

	users := []User{
		{Name: "Alice", Email: "alice@example.com", Preferences: "theme=dark,lang=en,timezone=UTC"},
		{Name: "Bob", Email: "bob@example.com", Preferences: "theme=light,lang=fr,timezone=CET"},
		{Name: "Carol", Email: "carol@example.com", Preferences: "theme=dark,lang=en,timezone=IST"},
	}
	if err := db.Create(&users).Error; err != nil {
		panic("failed to seed users")
	}
}
