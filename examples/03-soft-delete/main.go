package main

import (
	"fmt"
	"gorm-unloaded/db"
)

func main() {
	db := db.MustNew()

	mustBootstrap(db) // Migrate the database and seed initial data

	// Read the city and see that it exists before deletion
	var city City
	err := db.First(&city, "name = ?", "Maharagama").Error
	if err != nil {
		panic("failed to find city before deletion")
	}
	fmt.Printf("City read before deletion. City name: %s\n", city.Name)

	// Soft delete the city
	err = db.Delete(&city).Error
	if err != nil {
		panic("failed to delete city")
	}

	// Try to read the city again after deletion. This should fail as the record is soft deleted and won't be returned in normal queries.
	err = db.First(&city, "name = ?", "Maharagama").Error
	if err != nil {
		fmt.Printf("\nCity not found after deletion, as expected. Error: %v\n\n", err)
	} else {
		panic("city should not be found after deletion")
	}

	// However, if we want to include soft deleted records in our query, we can use Unscoped() method.
	err = db.Unscoped().First(&city, "name = ?", "Maharagama").Error
	if err != nil {
		panic("failed to find city even with Unscoped")
	}
	fmt.Printf("City read with Unscoped after deletion. City name: %s, DeletedAt: %v\n", city.Name, city.DeletedAt.Time)

	// Permanently delete the city record from the database
	err = db.Unscoped().Delete(&city).Error
	if err != nil {
		panic("failed to permanently delete city")
	}

	// Try to read the city again after permanent deletion. This should fail as the record is now permanently deleted.
	err = db.Unscoped().First(&city, "name = ?", "Maharagama").Error
	if err != nil {
		fmt.Printf("\nCity not found after permanent deletion, as expected. Error: %v\n\n", err)
	} else {
		panic("city should not be found after permanent deletion")
	}
}
