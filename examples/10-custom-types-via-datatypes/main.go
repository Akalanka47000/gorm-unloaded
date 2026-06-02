package main

import (
	"fmt"
	"gorm-unloaded/db"

	"gorm.io/datatypes"
)

func main() {
	database := db.MustNew()
	mustBootstrap(database)

	// Read all products — JSONB columns deserialise back into typed Go values.
	var products []Product
	database.Find(&products)

	fmt.Println("All products:")
	for _, p := range products {
		fmt.Printf("\n  %s\n", p.Name)
		fmt.Printf("    tags:       %v\n", []string(p.Tags))
		fmt.Printf("    attributes: color=%s size=%s material=%s\n",
			p.Attributes.Data().Color,
			p.Attributes.Data().Size,
			p.Attributes.Data().Material,
		)
	}

	// Query using datatypes.JSONQuery — filter by a specific JSON attribute value.
	// JSONQuery generates: attributes->>'color' = 'black'
	fmt.Println("\nProducts where color = black:")
	var byColor []Product
	err := database.Where(datatypes.JSONQuery("attributes").Equals("black", "color")).Find(&byColor).Error
	if err != nil {
		panic("failed to query products by color")
	}
	for _, p := range byColor {
		fmt.Printf("  %s  color=%s\n", p.Name, p.Attributes.Data().Color)
	}

	// Add a tag to one product and update.
	headphones := products[0]
	headphones.Tags = append(headphones.Tags, "new-arrival")
	err = database.Save(&headphones).Error
	if err != nil {
		panic("failed to update headphones")
	}

	fmt.Println("\nHeadphones tags after update:")
	err = database.First(&headphones, headphones.ID).Error
	if err != nil {
		panic("failed to fetch updated headphones")
	}
	fmt.Printf("  %v\n", []string(headphones.Tags))
}
