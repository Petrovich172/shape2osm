package main

import (
	"log"
	"github.com/go-pg/pg"
	"shape2osm/Utils"
)


func main() {
	log.Println("Heey!")

	db := pg.Connect(&pg.Options{
			Addr:      "172.20.12.159" + ":" + "5432",
			User:      "postgres",
			Password:  "postgres",
			Database:  "postgres",
		})
	defer db.Close()

	// Querying shaped geo from DB
	dbData := utils.GetSomeData(db)

	// Convert Shaped data to Osm
	osmData := utils.Convert(dbData)

	// Creating output xml file
	utils.Xml2file(osmData)
}


