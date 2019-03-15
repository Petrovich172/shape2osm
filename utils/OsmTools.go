package utils

import (
	"strconv"
	// "fmt"
	"log"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"math/rand"
	"io/ioutil"
	"io"
	"os"
	"encoding/xml"
	"shape2osm/OsmStructs"
	"shape2osm/ShapedStructs"

)


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
	_, err = db.Model().Query(&ret.Edges, sqlString1)
	if err != nil {
		log.Println("Problems with DB:", "\n", err)
		log.Panicln(err)
	}
	_, err = db.Model().Query(&ret.NodedLines, sqlString2)
	if err != nil {
		log.Println("Problems with DB:", "\n", err)
		log.Panicln(err)
	}
	log.Println("query answer first row:","\n",ret.Edges[0], ret.NodedLines[0])	
	return ret
}


func InsertOsm2DB(xmlData OsmStructs.Osm, db *pg.DB) {
	var err error
    err = db.CreateTable(&OsmStructs.Node{}, &orm.CreateTableOptions{
        // Temp:          true, // create temp table
   	})
   	if err != nil {
   		log.Println("can't create nodes table:", err)
   	} else {
   		log.Println("Successfully created Nodes table")
   	}
   	// Inserting nodes
   	err = db.Insert(&xmlData.Nodes)
   	if err != nil {
   		log.Panicln("can't insert nodes", err)
   	} else {
   		log.Println("Nodes successfully inserted")
   	}

   	type Way struct {
		Id 		int32	`sql:"id"`		
		Nds		[]int32 `sql:"nodes, type:integer[]"`
		Tags    []OsmStructs.Tag `sql:"tags"`

	}
    err = db.CreateTable(&OsmStructs.Way{}, &orm.CreateTableOptions{
        // Temp:          true, // create temp table    	
   	})
    if err != nil {
   		log.Println("can't create ways table:", err)
   	} else {
   		log.Println("Successfully created Ways table")
   	}

   	// Iterating []NdId{} to get only values
	NodesOut := func (nds []OsmStructs.NdId) (nodes []int32) {
		for _, id := range nds {
			nodes = append(nodes, id.ID)
		}
		return
	}
   	
	// Inserting ways
   	for _, way := range xmlData.Ways {
		err = db.Insert(&Way{
			Id:		way.Elem.ID,
			Nds:	NodesOut(way.Nds),
			// Nds:	way.Nds,
			Tags:	way.Tags,
		})
		if err != nil {
			log.Panicln("can't insert ways", err, way.Elem.ID)
		}
	}
	// var wway []Way
	// _, err = db.Query(wway, `select * FROM public.ways as ways limit 1`)
	// log.Println("Select ways:",err)
}

// Convert Shaped data to Osm format
func Convert (shaped ShapedStructs.ShapeData) OsmStructs.Osm {

	// Checking if given data is empty
	if shaped.Edges == nil || shaped.NodedLines == nil {
		log.Println("utils.Convert: Shaped data is empty")
		log.Panicln()
	}
	
	// Initiating xmlData — body struct for .xml 
	var xmlData OsmStructs.Osm

	// Initiating random ID for nodes
	generate := rand.New(rand.NewSource(99)).Int31

	// Initiating structs for nodes, ways and relations elements
	var nodeId OsmStructs.Elem
	var wayId OsmStructs.Elem
	var relId OsmStructs.Elem

	// Node Id array in Ways struct
	var nodeIDs []OsmStructs.NdId

	// Tags and members array
	var arrTags []OsmStructs.Tag
	var arrMember []OsmStructs.Member
	var restrictionsArr []OsmStructs.Tag


	// Iterating every node
	for i := 0; i < len(shaped.Edges); i++ {
			nodeId.ID = shaped.Edges[i].Id
			nodeId.Ts = "2019-01-01T00:00:00Z"
			nodeId.Version = 1
			xmlData.Nodes = append(xmlData.Nodes, OsmStructs.Node{
				Elem:	nodeId,
				Lat:	shaped.Edges[i].Geom.Coordinates[1],
				Lng:	shaped.Edges[i].Geom.Coordinates[0],
				// Tags:	arrTags,
				}	)
	}
	
	// Iterating every noded line (way)
	for i := 0; i < len(shaped.NodedLines); i++ {
		wayId.ID = shaped.NodedLines[i].Id
		wayId.Ts = "2019-01-01T00:00:00Z"
		wayId.Version = 1
		nodeIDs = nil
		arrTags = nil
		arrMember = nil

		// Surface types
		switch shaped.NodedLines[i].Surface {
		case "0":
			shaped.NodedLines[i].Surface = "no data"
		case "1":
			shaped.NodedLines[i].Surface = "unpaved"
		case "2":
			shaped.NodedLines[i].Surface = "asphalt"
		case "3":
			shaped.NodedLines[i].Surface = "rails"
		}

		// Road types
		switch shaped.NodedLines[i].SnipAd {
		case "0":
			shaped.NodedLines[i].SnipAd = "road"
		case "1":
			shaped.NodedLines[i].SnipAd = "motorway"
		case "2":
			shaped.NodedLines[i].SnipAd = "trunk"
		case "3":
			shaped.NodedLines[i].SnipAd = "primary"
		case "4":
			shaped.NodedLines[i].SnipAd = "secondary"
		case "5":
			shaped.NodedLines[i].SnipAd = "tertiary"
		case "6":
			shaped.NodedLines[i].SnipAd = "unclassified"
		}

		// TYP_COD types
		switch shaped.NodedLines[i].TlineTyp {
		case "7701":
			shaped.NodedLines[i].SnipAd = "footway"
		case "7702":
			shaped.NodedLines[i].SnipAd = "residential"
		case "7703":
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"parking:lane",
				Value:	"marked",
			})
			shaped.NodedLines[i].SnipAd = "road"
		case "7704":
			shaped.NodedLines[i].SnipAd = "corridor"
		case "7705":
			shaped.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"winter_road",
				Value:	"yes",
			})
		case "7706":
			shaped.NodedLines[i].SnipAd = "cycleway"
		case "7007":
			shaped.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"aerialway",
				Value:	"cable_car",
			})
		case "7730":
			shaped.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"railway",
				Value:	"rail",
			})
		case "7740":
			shaped.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"railway",
				Value:	"tram",
			})
		case "7750":
			shaped.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"railway",
				Value:	"subway",
			})
		case "7760":
			shaped.NodedLines[i].SnipAd = "*"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"railway",
				Value:	"monorail",
			})
		}

		// Construction types
		switch shaped.NodedLines[i].Btf {
		case "1":
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"bridge",
				Value:	"yes",
			})
		case "2":
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"bridge",
				Value:	"pontoon",
			})
		case "3":
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"tunnel",
				Value:	"yes",
			})
		case "4":
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"route",
				Value:	"ferry",
			})
		case "5":
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"railway",
				Value:	"level_crossing",
			})
		case "6":
			shaped.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"bridge",
				Value:	"yes",
			})
		case "7":
			shaped.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"footway",
				Value:	"crossing",
			})
		case "8":
			shaped.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"tunnel",
				Value:	"yes",
			})
		case "9":
			shaped.NodedLines[i].SnipAd = "steps"
		case "10":
			shaped.NodedLines[i].SnipAd = "steps"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"conveying",
				Value:	"yes",
			})
		case "11":
			shaped.NodedLines[i].SnipAd = "footway"
			arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"conveying",
				Value:	"yes",
			})		
		}

		// Bicycle road types
		switch shaped.NodedLines[i].Bicyclanes {
		case "1":
			arrTags = append(arrTags, 
				OsmStructs.Tag{
					Key:	"cycleway",
					Value:	"lane",
				})
		case "2":
			arrTags = append(arrTags, 
				OsmStructs.Tag{
					Key:	"cycleway",
					Value:	"opposite_lane",
				})
		case "3":
			shaped.NodedLines[i].SnipAd = "cycleway"
		}

		// Bus way types
		if shaped.NodedLines[i].FBuslanes == "1" {
			arrTags = append(arrTags, 
				OsmStructs.Tag{
					Key:	"busway:right",
					Value:	"lane",
				})
		} else if shaped.NodedLines[i].TBuslanes == "1" {
			arrTags = append(arrTags, 
				OsmStructs.Tag{
					Key:	"busway:left",
					Value:	"lane",
				})
		}

		// Size restrictions
		rWeight, err := strconv.Atoi(shaped.NodedLines[i].RWeight)
			if err != nil {
				log.Println(err)
			}
		if rWeight >= 2 {
			arrTags = append(arrTags,
				OsmStructs.Tag{
					Key:	"maxweight",
					Value:	shaped.NodedLines[i].RWeight,
					})
		} else if shaped.NodedLines[i].RHeight != "0" {
			arrTags = append(arrTags,
				OsmStructs.Tag{
					Key:	"maxheight",
					Value:	shaped.NodedLines[i].RHeight,
					})
		} else if shaped.NodedLines[i].RWidth != "0" {
			arrTags = append(arrTags,
				OsmStructs.Tag{
					Key:	"maxwidth",
					Value:	shaped.NodedLines[i].RWidth,
				})
		}

		// Filling tags array
		arrTags = append(arrTags, 
			OsmStructs.Tag{
				Key:	"highway",
				Value:	shaped.NodedLines[i].SnipAd,
			}, 
			OsmStructs.Tag{
				Key:	"oneway",
				Value:	shaped.NodedLines[i].Oneway,
			},
			OsmStructs.Tag{
				Key:	"surface",
				Value:	shaped.NodedLines[i].Surface,
			},
			OsmStructs.Tag{
				Key:	"lanes:forward",
				Value:	shaped.NodedLines[i].F_lanes,
			},
			OsmStructs.Tag{
				Key:	"lanes:backward",
				Value:	shaped.NodedLines[i].T_lanes,
			},
			OsmStructs.Tag{
				Key:	"maxspeed",
				Value:	shaped.NodedLines[i].Speedlim,
			},
			OsmStructs.Tag{
				Key:	"name",
				Value:	shaped.NodedLines[i].RdName,
			},
			OsmStructs.Tag{
				Key:	"toll",
				Value:	shaped.NodedLines[i].Tollway,
			}	)
		
		// Filling members array => relations
		if shaped.NodedLines[i].Edge2id > 0 {
			if shaped.NodedLines[i].Edge5id != 0 {
				arrMember = append(arrMember, 
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge2id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge3id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge4id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge5id,
						Role:	"",
					}	)
			} else if shaped.NodedLines[i].Edge5id == 0 && shaped.NodedLines[i].Edge4id != 0 {
				arrMember = append(arrMember, 
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge2id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge3id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge4id,
						Role:	"",
					}	)
			} else if shaped.NodedLines[i].Edge5id == 0 && shaped.NodedLines[i].Edge4id == 0 && shaped.NodedLines[i].Edge3id != 0 {
				arrMember = append(arrMember, 
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge2id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge3id,
						Role:	"",
					}	)
			} else {
				arrMember = append(arrMember, 
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Id,
						Role:	"",
					},
					OsmStructs.Member{
						Type:	"way",
						Ref:	shaped.NodedLines[i].Edge2id,
						Role:	"",
					}	)
			}

			/*// restriction tags
			restrictionsArr = nil
			if shaped.NodedLines[i].GmanTyp == "7980" {
				restrictionsArr = append(restrictionsArr, 
					OsmStructs.Tag{
						Key:	"type",
						Value:	"restriction",
					},
					OsmStructs.Tag{
						Key:	"restriction",
						Value:	"no_entry",
					}	)
			} else if shaped.NodedLines[i].GmanTyp == "7990" {
				restrictionsArr = append(restrictionsArr, 
					OsmStructs.Tag{
						Key:	"type",
						Value:	"restriction",
					},
					OsmStructs.Tag{
						Key:	"restriction",
						Value:	"no_u_turn",
					}	)
			} */

			// Relations
			relId.ID = generate()
			relId.Ts = "2019-01-01T00:00:00Z"
			relId.Version = 1
			xmlData.Relations = append(xmlData.Relations, OsmStructs.Relation{
				Elem:	relId,
				Members:	arrMember,
				Tags:	restrictionsArr,
				})			
		}

		// Filling .xml with ways
		wayId.Version = 1
		var tmpnode1 OsmStructs.NdId
		var tmpnode2 OsmStructs.NdId
		tmpnode1.ID = shaped.NodedLines[i].Source
		tmpnode2.ID = shaped.NodedLines[i].Target
		nodeIDs = append(nodeIDs, tmpnode1, tmpnode2)
		xmlData.Ways = append(xmlData.Ways, OsmStructs.Way{
			Elem:	wayId,
			Nds:	nodeIDs,
			Tags:	arrTags,
			}	)
	}
	xmlData.Ways = UniqWays(xmlData.Ways)
	return xmlData
}

// Cleaning id dublicates
func UniqWays(s []OsmStructs.Way) []OsmStructs.Way {
	seen := make(map[OsmStructs.Elem]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v.Elem]; ok {
			continue
		}
		seen[v.Elem] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

// Creating output xml file
func Xml2file (xmlData OsmStructs.Osm) {	
	xmlData.Version = "0.6"
	xmlData.Ts = "2019-01-28T01:59:52Z"
	f, err := os.Create("out.xml")
	if err != nil { panic(err) }
	defer f.Close()
	newFile := io.Writer(f)
	enc := xml.NewEncoder(newFile)
	f.Write([]byte("<?xml version=\"1\" encoding=\"UTF-8\"?>\n"))
	enc.Indent("", "    ")
    	if err := enc.Encode(&xmlData); err != nil {
				log.Printf("File writing error: %v\n", err, "%v\n", enc)
				log.Panicln(err)
		}
}

// Read xml file
func ReadXml(filename string) []byte {
	// Some stuff to open & read .xml file
	xmlFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	log.Println("Successfully opened", filename)
	defer xmlFile.Close()

	// Read our .xml file as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)
	return byteValue
}