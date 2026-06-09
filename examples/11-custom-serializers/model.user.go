package main

import "time"

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`

	// Go type: key=value string        →  "theme=dark,lang=en,timezone=UTC"
	// DB type: JSONB map               →  {"theme":"dark","lang":"en","timezone":"UTC"}
	//
	// The serializer tag is field-scoped — Name and Email are plain text columns;
	// only Preferences gets the kv_json treatment.
	Preferences string `gorm:"serializer:kv_json"`

	CreatedAt time.Time
	UpdatedAt *time.Time
}
