package utils

import (
    // "bufio"
    // "encoding/csv"
    "log"
    // "bytes"
	"encoding/gob"    
    // "io"
    // "io/ioutil"
   	"shape2osm/OsmStructs"
	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg"
    "os"
)

type Contracted struct {
	Seq					int32	`sql:"seq"`
	Type				string	`sql:"type"`
	Id					int32	`sql:"id"`
	ContractedVertices	[]int32	`sql:"contracted_vertices"`
	Source				int32	`sql:"source"`
	Target				int32	`sql:"target"`
	Cost				float32	`sql:"cost"`
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ReadBytes() {
	file := "/home/pete/OSRM/data/out.xml.osrm"
	// dat, err := ioutil.ReadFile(file)
	// check(err)
	// log.Print(string(dat))

	f, err := os.Open(file)
	check(err)

	b1 := make([]byte, 2000)
	n1, err := f.Read(b1)
	check(err)
	log.Printf("%d bytes: %s\n", n1, string(b1))
}

func OsmContract(xmlData OsmStructs.Osm, db *pg.DB) []Contracted {
	var err error
    err = db.CreateTable(&OsmStructs.Node{}, &orm.CreateTableOptions{
        Temp:          true,
   	})
   	if err != nil {
   		log.Println("can't create nodes table:", err)
   	} else {
   		log.Println("Successfully created Nodes table")
   	}
   	// Inserting nodes
   	// err = db.TableExpr("pg_temp.nodes AS ways").Insert(&xmlData.Nodes)
   	err = db.Insert(&xmlData.Nodes)
   	if err != nil {
   		log.Println("can't insert nodes", err)
   	} else {
   		log.Println("Nodes successfully inserted")
   	}

    err = db.CreateTable(&OsmStructs.Way{}, &orm.CreateTableOptions{
        Temp:          true,
   	})
    if err != nil {
   		log.Println("can't create ways table:", err)
   	} else {
   		log.Println("Successfully created Ways table")
   	}

   	arr := NodesOut(way.Nds)
    arrStr := make([]string, len(arr))
    for i := range arr {
      arrStr[i] = strconv.Itoa(int(arr[i]))
    }
    sqlString := fmt.Sprintf("INSERT INTO ways (id, nodes) VALUES(%d, ARRAY[%v]::integer[]);", way.Elem.ID, strings.Join(arrStr, ","))
    _, err := db.ExecOne(sqlString)
    

   	// Iterating []NdId{} to get only values
	NodesOut := func (nds []OsmStructs.NdId) (nodes []int32) {
		for _, id := range nds {
			nodes = append(nodes, id.ID)
		}
		return
	}

	type Way struct {
		Id 		int32	`sql:"id"`		
		Nds		[]int32 `sql:"nodes, type:integer[]"`
		Tags    []OsmStructs.Tag `sql:"tags"`

	}
	// Inserting ways
   	for _, way := range xmlData.Ways {
   		err = db.Insert(&Way{
   			Id:		way.Elem.ID,
   			Nds:	NodesOut(way.Nds),
   			Tags:	way.Tags,
   		})
   	   	if err != nil {
   		log.Println("can't insert ways", err)
   		}
   	}
  //  	type nd struct {
  //  		Id 		int32	 `sql:"id"`
		// Lat     float64  `xml:"lat,attr" sql:"lat"`
		// Lng     float64  `xml:"lon,attr" sql:"lon"`
  //  	}

  //  	var ndd *nd = &nd{}
   	// var wwway *Way = &Way{}
	var wway []Way
	_, err = db.Query(&wway, `select id, nodes, tags FROM pg_temp.ways as ways limit 1`)
   	log.Println("Select ways:",err)

   	// Contracting with pgr_contractGraph tool
   	var res []Contracted
   	_, err = db.Query(&res, `
   		SELECT * FROM pgr_contractGraph(
			'SELECT ways.id, nodes1.id as source, nodes2.id as target, 1 as cost FROM pg_temp.ways as ways 
			join pg_temp.nodes as nodes1 on nodes1.id = ways.nodes[1]
			join pg_temp.nodes as nodes2 on nodes2.id = ways.nodes[2]
			', ARRAY[1, 2]);
	`)
	if err != nil {
		log.Panicln("can't contract:", err)
	} else {
		log.Println("Contracted!")
	}
	return res
}

// func CsvExport(contracted []Contracted) error {
// 	var data [][]byte
// 	for _, d := range contracted {
// 		data = append(data, []byte(d))
// 	}
//     file, err := os.Create("result.csv")
//     if err != nil {
//         return err
//     }
//     defer file.Close()

//     writer := csv.NewWriter(file)
//     defer writer.Flush()

//     for _, value := range data {
//         if _, err := file.Write([]byte(value)); err != nil {
//             return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
//         }
//     }
//     return nil
// }

func WriteContracted(contracted []Contracted) error {
	err := writeGob("./contracted.gob",contracted)
	if err != nil{
		log.Println(err)
	}
	return err
}

// Using gob
func writeGob(filePath string,object interface{}) error {
       file, err := os.Create(filePath)
       if err == nil {
              encoder := gob.NewEncoder(file)
              encoder.Encode(object)
       }
       file.Close()
       return err
}
