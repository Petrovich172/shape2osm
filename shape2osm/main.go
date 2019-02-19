package main

import (
	// "time"
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
	Id	int32	`xml:"id"	sql:"id"`
	// Edge1id	int32	`xml:"ref"	sql:"edge1id"`
	Edge2id	int32	`xml:"ref"	sql:"edge2id"`
	Edge3id	int32	`xml:"ref"	sql:"edge3id"`
	Edge4id	int32	`xml:"ref"	sql:"edge4id"`
	Edge5id	int32	`xml:"ref"	sql:"edge5id"`
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
	var xmlData cfg.Osm

	// initiating random ID for nodes
	generate := rand.New(rand.NewSource(99)).Int31

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
		wayId.Ts = "2019-01-01T00:00:00Z"
		wayId.Version = 1
		nodeIDs = nil
		arrTags = nil
		arrMember = nil
		// filling tags array
		switch dbData[i].Surface {
		case "0":
			dbData[i].Surface = "no data"
		case "1":
			dbData[i].Surface = "unpaved"
		case "2":
			dbData[i].Surface = "asphalt"
		case "3":
			dbData[i].Surface = "rails"
		}
		arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"highway",
				Value:	dbData[i].Highway,
			}, 
			cfg.Tag{
				Key:	"oneway",
				Value:	dbData[i].Oneway,
			},
			cfg.Tag{
				Key:	"surface",
				Value:	dbData[i].Surface,
			}	)
		// filling members array => relations
		if dbData[i].Edge2id > 0 {
			if dbData[i].Edge5id != 0 {
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
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData[i].Edge3id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData[i].Edge4id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData[i].Edge5id,
						Role:	"",
					}	)
			} else if dbData[i].Edge5id == 0 && dbData[i].Edge4id != 0 {
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
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData[i].Edge3id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData[i].Edge4id,
						Role:	"",
					}	)
			} else if dbData[i].Edge5id == 0 && dbData[i].Edge4id == 0 && dbData[i].Edge3id != 0 {
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
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData[i].Edge3id,
						Role:	"",
					}	)
			} else {
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
			}

			relId.ID = generate()
			relId.Ts = "2019-01-01T00:00:00Z"
			relId.Version = 1
			xmlData.Relations = append(xmlData.Relations, cfg.Relation{
				Elem:	relId,
				Members:	arrMember,
				})			
		}

		// iterating every node
		for y := 0; y < len(node); y++ {
			nodeId.ID = generate()
			nodeId.Ts = "2019-01-01T00:00:00Z"
			nodeId.Version = 1
			xmlData.Nodes = append(xmlData.Nodes, cfg.Node{
				Elem:	nodeId,
				Lat:	node[y][1],
				Lng:	node[y][0],
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

	xmlData.Version = "0.6"
	xmlData.Ts = "2019-01-28T01:59:52Z"

	// creating output xml file
	f, err := os.Create("out.xml")
	if err != nil { panic(err) }
	defer f.Close()
	newFile := io.Writer(f)
	enc := xml.NewEncoder(newFile)
	f.Write([]byte("<?xml version=\"1\" encoding=\"UTF-8\"?>\n"))
	enc.Indent("", "    ")
    	if err := enc.Encode(&xmlData); err != nil {
				log.Printf("error: %v\n", err, "%v\n", enc)
				// log.Println("Osm:", xmlData)
		}
}


// Get geo data from DB
func getSomeData(db *pg.DB) []Edge {
	var ret []Edge
	var err error
	sqlString := `SELECT "tline".id as id, ST_AsGeoJSON(ST_Transform("tline".geom, 4326)) as geom, oneway, surface, highway, edge2id, edge3id, edge4id, edge5id
				FROM graph.tline as "tline" left join graph.gman as "gman" on "tline".id = "gman".edge1id  
				--where "tline".id = "gman".edge1id
				order by id desc
				limit 1000`
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