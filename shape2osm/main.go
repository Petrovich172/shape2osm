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
type NodedLine struct {
	Id			int32	`xml:"id"			sql:"id"`
	Source		int32	`xml:"source"		sql:"source"`
	Target		int32	`xml:"target"		sql:"target"`
	F_lanes		string	`xml:"f_lanes"		sql:f_lanes`
	T_lanes		string	`xml:"t_lanes"		sql:"t_lanes"`
	GmanTyp 	string 	`xml:"gman_typ"		sql:"gman_typ"`
	TlineTyp	string 	`xml:"tline_typ"	sql:"tline_typ"`
	Speedlim	string	`xml:"speedlim" 	sql:"speedlim"`
	RdName		string	`xml:"rd_name" 		sql:"rd_name"`
	Tollway		string 	`xml:"tollway" 		sql:"tollway"`
	SnipAd 		string 	`xml:"snip_ad"		sql:"snip_ad"`
	Btf 		string 	`xml:"btf"			sql:"btf"`
	// Edge1id		int32	`xml:"ref"			sql:"edge1id"`
	Edge2id		int32	`xml:"ref"			sql:"edge2id"`
	Edge3id		int32	`xml:"ref"			sql:"edge3id"`
	Edge4id		int32	`xml:"ref"			sql:"edge4id"`
	Edge5id		int32	`xml:"ref"			sql:"edge5id"`
	Oneway		string	`xml:"oneway"		sql:"oneway"`
	Surface		string	`xml:"surface"		sql:"surface"`
	Highway		string	`xml:"highway"		sql:"highway"`
}

type Edge struct {
	Id	int32	`xml:"id"	sql:"id"`
	Geom	utils.PointString	`xml:"geom"	sql:"geom"`	
}

type DbGeom struct {
	Edges	[]Edge
	NodedLines	[]NodedLine
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
	// var TempNodeId cfg.NdId
	var nodeIDs []cfg.NdId

	// tags and members array
	var arrTags []cfg.Tag
	var arrMember []cfg.Member
	var restrictionsArr []cfg.Tag

	// querying geo from DB
	dbData := getSomeData(db)

	// iterating every node
	for i := 0; i < len(dbData.Edges); i++ {
			nodeId.ID = dbData.Edges[i].Id
			nodeId.Ts = "2019-01-01T00:00:00Z"
			nodeId.Version = 1
			xmlData.Nodes = append(xmlData.Nodes, cfg.Node{
				Elem:	nodeId,
				Lat:	dbData.Edges[i].Geom.Coordinates[1],
				Lng:	dbData.Edges[i].Geom.Coordinates[0],
				// Tags:	arrTags,
				}	)
	}
	
	// iteration every noded line (way)
	for i := 0; i < len(dbData.NodedLines); i++ {
		wayId.ID = dbData.NodedLines[i].Id
		wayId.Ts = "2019-01-01T00:00:00Z"
		wayId.Version = 1
		nodeIDs = nil
		arrTags = nil
		arrMember = nil
		// filling tags array

		// Surface types
		switch dbData.NodedLines[i].Surface {
		case "0":
			dbData.NodedLines[i].Surface = "no data"
		case "1":
			dbData.NodedLines[i].Surface = "unpaved"
		case "2":
			dbData.NodedLines[i].Surface = "asphalt"
		case "3":
			dbData.NodedLines[i].Surface = "rails"
		}

		// road types
		switch dbData.NodedLines[i].SnipAd {
		case "0":
			dbData.NodedLines[i].SnipAd = "road"
		case "1":
			dbData.NodedLines[i].SnipAd = "motorway"
		case "2":
			dbData.NodedLines[i].SnipAd = "trunk"
		case "3":
			dbData.NodedLines[i].SnipAd = "primary"
		case "4":
			dbData.NodedLines[i].SnipAd = "secondary"
		case "5":
			dbData.NodedLines[i].SnipAd = "tertiary"
		case "6":
			dbData.NodedLines[i].SnipAd = "unclassified"
		}

		// TYP_COD types
		switch dbData.NodedLines[i].TlineTyp {
		case "7701":
			dbData.NodedLines[i].SnipAd = "footway"
		case "7702":
			dbData.NodedLines[i].SnipAd = "residential"
		case "7703":
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"parking:lane",
				Value:	"marked",
			})
			dbData.NodedLines[i].SnipAd = "road"
		case "7704":
			dbData.NodedLines[i].SnipAd = "corridor"
		case "7705":
			dbData.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"winter_road",
				Value:	"yes",
			})
		case "7706":
			dbData.NodedLines[i].SnipAd = "cycleway"
		case "7007":
			dbData.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"aerialway",
				Value:	"cable_car",
			})
		case "7730":
			dbData.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"railway",
				Value:	"rail",
			})
		case "7740":
			dbData.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"railway",
				Value:	"tram",
			})
		case "7750":
			dbData.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"railway",
				Value:	"subway",
			})
		case "7760":
			dbData.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"railway",
				Value:	"monorail",
			})
		}

		// construction types
		switch dbData.NodedLines[i].Btf {
		case "1":
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"bridge",
				Value:	"yes",
			})
		case "2":
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"bridge",
				Value:	"pontoon",
			})
		case "3":
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"tunnel",
				Value:	"yes",
			})
		case "4":
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"route",
				Value:	"ferry",
			})
		case "5":
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"railway",
				Value:	"level_crossing",
			})
		case "6":
			dbData.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"bridge",
				Value:	"yes",
			})
		case "7":
			dbData.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"footway",
				Value:	"crossing",
			})
		case "8":
			dbData.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"tunnel",
				Value:	"yes",
			})
		case "9":
			dbData.NodedLines[i].SnipAd = "steps"
		case "10":
			dbData.NodedLines[i].SnipAd = "steps"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"conveying",
				Value:	"yes",
			})
		case "11":
			dbData.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"conveying",
				Value:	"yes",
			})			
		}

		// filling tags array
		arrTags = append(arrTags, 
			cfg.Tag{
				Key:	"highway",
				Value:	dbData.NodedLines[i].SnipAd,
			}, 
			cfg.Tag{
				Key:	"oneway",
				Value:	dbData.NodedLines[i].Oneway,
			},
			cfg.Tag{
				Key:	"surface",
				Value:	dbData.NodedLines[i].Surface,
			},
			cfg.Tag{
				Key:	"lanes:forward",
				Value:	dbData.NodedLines[i].F_lanes,
			},
			cfg.Tag{
				Key:	"lanes:backward",
				Value:	dbData.NodedLines[i].T_lanes,
			},
			cfg.Tag{
				Key:	"maxspeed",
				Value:	dbData.NodedLines[i].Speedlim,
			},
			cfg.Tag{
				Key:	"name",
				Value:	dbData.NodedLines[i].RdName,
			},
			cfg.Tag{
				Key:	"toll",
				Value:	dbData.NodedLines[i].Tollway,
			}	)
		
		// filling members array => relations
		if dbData.NodedLines[i].Edge2id > 0 {
			if dbData.NodedLines[i].Edge5id != 0 {
				arrMember = append(arrMember, 
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge2id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge3id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge4id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge5id,
						Role:	"",
					}	)
			} else if dbData.NodedLines[i].Edge5id == 0 && dbData.NodedLines[i].Edge4id != 0 {
				arrMember = append(arrMember, 
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge2id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge3id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge4id,
						Role:	"",
					}	)
			} else if dbData.NodedLines[i].Edge5id == 0 && dbData.NodedLines[i].Edge4id == 0 && dbData.NodedLines[i].Edge3id != 0 {
				arrMember = append(arrMember, 
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge2id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge3id,
						Role:	"",
					}	)
			} else {
				arrMember = append(arrMember, 
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Id,
						Role:	"",
					},
					cfg.Member{
						Type:	"way",
						Ref:	dbData.NodedLines[i].Edge2id,
						Role:	"",
					}	)
			}

			/*// restriction tags
			restrictionsArr = nil
			if dbData.NodedLines[i].GmanTyp == "7980" {
				restrictionsArr = append(restrictionsArr, 
					cfg.Tag{
						Key:	"type",
						Value:	"restriction",
					},
					cfg.Tag{
						Key:	"restriction",
						Value:	"no_entry",
					}	)
			} else if dbData.NodedLines[i].GmanTyp == "7990" {
				restrictionsArr = append(restrictionsArr, 
					cfg.Tag{
						Key:	"type",
						Value:	"restriction",
					},
					cfg.Tag{
						Key:	"restriction",
						Value:	"no_u_turn",
					}	)
			} */

			// relations
			relId.ID = generate()
			relId.Ts = "2019-01-01T00:00:00Z"
			relId.Version = 1
			xmlData.Relations = append(xmlData.Relations, cfg.Relation{
				Elem:	relId,
				Members:	arrMember,
				Tags:	restrictionsArr,
				})			
		}

		// filling .xml with ways
		wayId.Version = 1
		var tmpnode1 cfg.NdId
		var tmpnode2 cfg.NdId
		tmpnode1.ID = dbData.NodedLines[i].Source
		tmpnode2.ID = dbData.NodedLines[i].Target
		nodeIDs = append(nodeIDs, tmpnode1, tmpnode2)
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
		}
}


// Get geo data from DB
func getSomeData(db *pg.DB) DbGeom {
	var ret DbGeom
	var err error
	sqlString1 := `select id, st_asgeojson(the_geom) as geom from graph.tline_2_noded_vertices_pgr`
	sqlString2 := `SELECT btf, snip_ad, tollway, rd_name, speedlim, f_lanes, t_lanes, "tline_old".typ_cod as tline_typ, "gman".typ_cod as gman_typ, "tline".id as id, "source", target, oneway, surface, highway, edge2id, edge3id, edge4id, edge5id
					from graph.tline_2_noded as "tline"
						join graph.tline as "tline_old" on "tline_old".id = "tline".old_id
						left join graph.gman as "gman" on "tline".old_id = "gman".edge1id
									--where "gman".edge2id >= 0 and
									where st_isempty(the_geom) is false and "source" <> target
									`
	_, err = db.Model().Query(&ret.Edges, sqlString1)
	if err != nil {
		log.Println("some shit happend:", "\n", err)		
	}
	_, err = db.Model().Query(&ret.NodedLines, sqlString2)
	if err != nil {
		log.Println("some shit happend:", "\n", err)		
	}
	log.Println("query:","\n",ret.Edges[0], ret.NodedLines[0])	
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