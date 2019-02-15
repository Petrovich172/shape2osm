package main

import (
	"log"
	"github.com/go-pg/pg"
	"math/rand"
	// "encoding/json"
	"io/ioutil"
	"io"
	"os"
	"encoding/xml"
	utils "./pkgs/utils"
	"./pkgs/cfg"
)

// Struct for geo data from DB
type Edge struct {
	Id	int64	`xml:"id"	sql:"id"`
	Edge2id	int64	`xml:"ref"	sql:"edge2id"`
	Edge3id	int64	`xml:"ref"	sql:"edge3id"`
	Edge4id	int64	`xml:"ref"	sql:"edge4id"`
	Edge5id	int64	`xml:"ref"	sql:"edge5id"`
	Oneway	string	`xml:"oneway"	sql:"oneway"`
	Surface	string	`xml:"surface"	sql:"surface"`
	Highway	string	`xml:"highway"	sql:"highway"`
	Geom	utils.MultiLineString	`xml:"geom"	sql:"geom"`
}


func main() {
	log.Println("Heey!")

	db := pg.Connect(&pg.Options{
			Addr:      "172.20.12.159" + ":" + "5432",
			User:      "postgres",
			Password:  "postgres",
			Database:  "postgres",
		})
	defer db.Close()

	// xmlData — body struct for .xml 
	var xmlData cfg.Map

	// initiating random ID for nodes
	generate := rand.New(rand.NewSource(99)).Int63

	// initiating structs for nodes, ways and relations ID
	var nodeId cfg.Elem
	var wayId cfg.Elem
	var relId cfg.Elem

	// node Id array in Ways struct
	var TempNodeId cfg.NdId
	var nodeIDs []cfg.NdId

	// tags and members array
	var arrTags []cfg.Tag
	var arrMember []cfg.Member

	// querying geo from DB
	dbData := getSomeData(db)

	// iterating every row from DB
	for i := 0; i < len(dbData); i++ {

		node := dbData[i].Geom.Coordinates[0]
		wayId.ID = dbData[i].Id
		nodeIDs = nil
		arrTags = nil
		arrMember = nil
		// filling tags array
		arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"Highway",
				Value:	dbData[i].Highway,
			}, 
			cfg.Tag{
				Key:	"Oneway",
				Value:	dbData[i].Oneway,
			},
			cfg.Tag{
				Key:	"Surface",
				Value:	dbData[i].Surface,
			}	)
		// filling members array => relations
		arrMember = append(arrMember, 
			cfg.Member{
				Type:	"way",
				Ref:	dbData[i].Id,
				Role:	"",
				},
			cfg.Member{
				Type:	"way",
				Ref:	dbData[i].Edge2id,
				Role:	"",
			}	)
		relId.ID = generate()
		xmlData.Relations = append(xmlData.Relations, cfg.Relation{
			Elem:	relId,
			Members:	arrMember,
			})

		// iterating every node
		for y := 0; y < len(node); y++ {
			nodeId.ID = generate()
			xmlData.Nodes = append(xmlData.Nodes, cfg.Node{
				Elem:	nodeId,
				Lat:	node[y][0],
				Lng:	node[y][1],
				Tags:	arrTags,
				}	)
			// making array of node ID for ways
			TempNodeId.ID = nodeId.ID
			nodeIDs = append(nodeIDs, TempNodeId)
		}

		// filling .xml with ways
		xmlData.Ways = append(xmlData.Ways, cfg.Way{
			Elem:	wayId,
			Nds:	nodeIDs,
			Tags:	arrTags,
			}	)
	}

		// Encode to XML
	// x, _ := xml.MarshalIndent(WayTags(xmlData.Ways[0].Tags), "", "  ")
	// log.Println(string(x))
	
	// log.Println(xmlData)
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	// xml.Unmarshal(	getSomeData(db), &data	)	

	// creating output xml file
	f, err := os.Create("out.xml")
	if err != nil { panic(err) }
	defer f.Close()
	newFile := io.Writer(f)
	enc := xml.NewEncoder(newFile)
	enc.Indent("  ", "    ")
    	if err := enc.Encode(&xmlData); err != nil {
				log.Printf("error: %v\n", err, "%v\n", enc)
				// log.Println("map:", xmlData)
		}
}


// Get geo data from DB
func getSomeData(db *pg.DB) []Edge {
	var ret []Edge
	var err error
	sqlString := `SELECT "tline".id as id, ST_AsGeoJSON(ST_Transform("tline".geom, 4326)) as geom, oneway, surface, highway, edge2id, edge3id, edge4id, edge5id
				FROM graph.tline as "tline", graph.gman as "gman" 
				where "tline".id = "gman".edge1id
				limit 100`
	_, err = db.Model().Query(&ret, sqlString)
	if err != nil {
		log.Println("some shit happend:", "\n", err)		
	}
	log.Println("query:","\n",ret)
	return ret
}

func ReadXml(filename string) []byte {
	// Some stuff to open & read .xml file
	// Open our .xml file
	xmlFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully opened", filename)
	defer xmlFile.Close()

	// read our .xml file as a byte array.	
	byteValue, _ := ioutil.ReadAll(xmlFile)
	return byteValue
}