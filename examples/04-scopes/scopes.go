package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ActiveOnly is a plain scope — no parameters needed, so it matches func(*gorm.DB) *gorm.DB directly.
func ActiveOnly(db *gorm.DB) *gorm.DB {
	return db.Where("active = ?", true)
}

// InCountry is a parameterised scope — it returns a func(*gorm.DB) *gorm.DB via closure.
func InCountry(country string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("country = ?", country)
	}
}

// Paginate is a parameterised scope that wraps Offset + Limit into a single reusable unit.
func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
}

type contextKey string

// ViewerCountryKey is the context key for storing the viewer's country.
const ViewerCountryKey contextKey = "viewerCountry"

// ForViewerCountry reads the viewer's country out of db.Statement.Context and scopes the
// query to that country. The caller never passes a country argument — the scope pulls it
// from the context that was attached via db.WithContext(ctx).
// This pattern is useful for tenant isolation or per-user query filtering.
func ForViewerCountry(db *gorm.DB) *gorm.DB {
	country, ok := db.Statement.Context.Value(ViewerCountryKey).(string)
	if !ok || country == "" {
		return db
	}
	return db.Where("country = ?", country)
}

// CountrySummaryRow is the scan target for the CountrySummary scope.
type CountrySummaryRow struct {
	Country     string
	Total       int
	ActiveCount int
}

// CountrySummary builds an aggregation query manully by writing SQL clauses and selected columns
// directly to db.Statement fields instead of going through builder methods.
// 
// The purpose of this is to demonstrate that scopes have full access to the db.Statement struct and can manipulate it directly, which allows for maximum flexibility in how scopes are implemented.
//
// db.Statement is the single struct that accumulates all query state.
// Builder methods (Select, Group, Order, Where, Joins, Omit…) are thin wrappers
// that ultimately write to these fields — you can bypass them entirely.
//
//	db.Statement.Selects   — columns passed to SELECT ([]string)
//	db.Statement.Omits     — columns excluded from SELECT/INSERT/UPDATE ([]string)
//	db.Statement.Table     — resolved table name (string, readable at scope time)
//	db.Statement.Clauses   — map[string]clause.Clause holding WHERE/ORDER/GROUP/LIMIT…
//	db.Statement.Joins     — slice of join descriptors
//	db.Statement.Vars      — positional bind variables matched to placeholders in SQL
//	db.Statement.Context   — the context attached via db.WithContext()
func CountrySummary(db *gorm.DB) *gorm.DB {
	// Directly overwrite the selected columns list.
	db.Statement.Selects = []string{
		"country",
		"COUNT(*) AS total",
		"SUM(CASE WHEN active THEN 1 ELSE 0 END) AS active_count",
	}

	// Directly write a GROUP BY clause into the Clauses map.
	db.Statement.AddClause(clause.GroupBy{
		Columns: []clause.Column{{Name: "country"}},
	})

	// Directly write an ORDER BY clause.
	db.Statement.AddClause(clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Raw: true, Name: "total DESC"}},
		},
	})

	return db
}
