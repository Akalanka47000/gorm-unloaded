package main

import (
	"context"
	"fmt"
	"gorm-unloaded/db"
)

func main() {
	db := db.MustNew()

	mustBootstrap(db)

	// --- 1. Simple scope ---
	// ActiveOnly is a plain func(*gorm.DB) *gorm.DB — passed directly, no call needed.
	var activeCities []City
	err := db.Scopes(ActiveOnly).Find(&activeCities).Error
	if err != nil {
		panic("failed to query active cities")
	}
	fmt.Printf("Found %d active cities:\n", len(activeCities))

	// --- 2. Parameterised scope ---
	// InCountry wraps the WHERE condition in a closure so it accepts arguments.
	// The return value — func(*gorm.DB) *gorm.DB — is what db.Scopes() expects.
	var sriLankaCities []City
	err = db.Scopes(InCountry("Sri Lanka")).Find(&sriLankaCities).Error
	if err != nil {
		panic("failed to query cities in Sri Lanka")
	}
	fmt.Printf("\nFound %d cities in Sri Lanka\n", len(sriLankaCities))

	// --- 3. Multiple scopes ---
	// Scopes can be chained together, so you can combine multiple scopes in a single query.
	var page1 []City
	err = db.Scopes(ActiveOnly, InCountry("Sri Lanka"), Paginate(1, 3)).Find(&page1).Error
	if err != nil {
		panic("failed to query page 1 cities")
	}

	fmt.Printf("\nPage 1 of active cities in Sri Lanka (3 per page):\n")
	for _, city := range page1 {
		fmt.Printf("- %s\n", city.Name)
	}

	// --- 4. Context-aware scope ---
	// Scopes receive *gorm.DB and can read db.Statement.Context directly.
	// Store a value in the context once; the scope picks it up automatically —
	// the caller never has to pass a country argument.
	ctx := context.WithValue(context.Background(), ViewerCountryKey, "India")
	var viewerCities []City
	err = db.WithContext(ctx).Scopes(ForViewerCountry, ActiveOnly).Find(&viewerCities).Error
	if err != nil {
		panic("failed to query viewer cities")
	}
	fmt.Printf("\nActive cities visible to Indian viewer (%d):\n", len(viewerCities))
	for _, c := range viewerCities {
		fmt.Printf("- %s\n", c.Name)
	}

	// --- 5. Direct db.Statement field manipulation ---
	// CountrySummaryDirect writes to db.Statement.Selects, db.Statement.Clauses (GROUP BY,
	// ORDER BY) directly — no builder methods called inside the scope at all.
	// Every GORM builder method is a thin wrapper over these same fields; the scope
	// skips the wrapper and writes to the struct directly.
	var summaries []CountrySummaryRow
	err = db.Model(&City{}).Scopes(CountrySummary).Scan(&summaries).Error
	if err != nil {
		panic("failed to query country summaries")
	}
	fmt.Printf("\nCountry summary via db.Statement manipulation:\n")
	for _, s := range summaries {
		fmt.Printf("- %-12s total: %d  active: %d\n", s.Country, s.Total, s.ActiveCount)
	}
}
