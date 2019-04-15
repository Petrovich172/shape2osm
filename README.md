# shape2xml tool
- Basically, converts given shaped-format geo data to osm format.
- Can query shaped data from database.
- Creates .xml file, filled with osm data.
- Can upload recieved osm data to DB
- Contrants osm data, using pgr_contractGraph and DB temp tables
- Writes contracted data to gob or csv format

## Using
1. Set up your DB parameters in main.go
```	
	db := pg.Connect(&pg.Options{
			Addr:      "172.20.12.159" + ":" + "5432",
			User:      "postgres",
			Password:  "postgres",
			Database:  "postgres",
		})
```

2. Get shaped geo data from you DB
```	// Querying shaped geo data from DB
	dbData := utils.GetSomeData(db)
```

**Set up your own sql query with correct tables at OsmTools.go**
* sqlString1 is for nodes coordinates with ids
* sqlString2 is for edges with source and target node id and all geo information data (e.g. speedlim, surface, oneway etc.)

```
// Get shaped geo data from DB
func GetSomeData(db *pg.DB) ShapedStructs.ShapeData {
	ret := ShapedStructs.ShapeData{}
	var err error
	sqlString1 := `select id, st_asgeojson(the_geom) as geom from graph.tline_2_noded_vertices_pgr`
	sqlString2 := `SELECT bicyclanes, t_buslanes, f_buslanes, "tline_old".r_weight as r_weight, "tline_old".r_height as r_height, "tline_old".r_width as r_width, 
					btf, snip_ad, tollway, rd_name, speedlim, f_lanes, t_lanes, "tline_old".typ_cod as tline_typ, "gman".typ_cod as gman_typ, "tline".id as id, 
					"source", target, oneway, surface, highway, edge2id, edge3id, edge4id, edge5id
					from graph.tline_2_noded as "tline"
						join graph.tline as "tline_old" on "tline_old".id = "tline".old_id
						left join graph.gman as "gman" on "tline".old_id = "gman".edge1id
							where st_isempty(the_geom) is false and "source" <> target
									`
```

3. Tool, converting shaped to osm format
```	// Converting Shaped data to Osm format
	osmData := utils.Convert(dbData)
```

4. Tools, creating xml file with osm data and inserting data to DB. Comment if you don't need that
```	// Creating output xml file
	utils.Xml2file(osmData)

	// Inserting OSM data to DB
	// contracted := utils.InsertOsm2DB(osmData, db)
```

5. Contracting osm data, using pgr_contractGraph tool. Creates temp tables in your DB
```	// Contracting using temp table and pgr_contractGraph tool
	contracted := utils.OsmContract(osmData, db)
```

6. Then you can write your contracted data either in .gob
```	
	// Writing contracted data to the file contracted.gob
	utils.WriteContracted(contracted)
```
or .csv format
```	
	utils.CsvExport(contracted)	
```

Good luck!
