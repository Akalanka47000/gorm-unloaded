package main

import (
	"fmt"

	extraClausePlugin "github.com/WinterYukky/gorm-extra-clause-plugin"
	"github.com/WinterYukky/gorm-extra-clause-plugin/exclause"
	"gorm-unloaded/db"
)

func main() {
	database := db.MustNew()
	database.Use(extraClausePlugin.New())
	mustBootstrap(database)

	var results []Product
	database.Model(&Product{}).
		Where("featured = ?", true).
		Clauses(exclause.NewUnion(database.Model(&Product{}).Where("on_sale = ?", true))).
		Order("name").
		Scan(&results)

	fmt.Println("--- UNION: featured OR on_sale ---")
	for _, p := range results {
		fmt.Printf("  %-10s featured=%-5v on_sale=%v\n", p.Name, p.Featured, p.OnSale)
	}
}
