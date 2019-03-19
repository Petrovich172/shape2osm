package main

import (
	"log"
	"github.com/go-pg/pg"
	// "github.com/go-pg/pg/orm"
	"shape2osm/utils"
	// "os"
)

func main() {
	log.Println("Heey!")

	db := pg.Connect(&pg.Options{
			Addr:      "172.20.12.159" + ":" + "5432",
			User:      "postgres",
			Password:  "postgres",
			Database:  "postgres",
		})
	

	// db.AddQueryHook(dbLogger{})
	defer db.Close()

	// Querying shaped geo data from DB
	dbData := utils.GetSomeData(db)

	// Converting Shaped data to Osm format
	osmData := utils.Convert(dbData)

	// Creating output xml file
	utils.Xml2file(osmData)

	// Inserting OSM data to DB
	// contracted := utils.InsertOsm2DB(osmData, db)
	

	// Contracting using temp table and pgr_contractGraph tool
	contracted := utils.OsmContract(osmData, db)
	
	// Writing contracted data to the file contracted.gob
	utils.WriteContracted(contracted)
	// log.Println(contracted)

	utils.ReadBytes()


}

// type dbLogger struct{}

// func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}
// func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
//   log.Println(q.FormattedQuery())
// }