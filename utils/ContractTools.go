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
	// "github.com/go-pg/pg/orm"
	"github.com/go-pg/pg"
    "os"
    "fmt"
    "strconv"
    "strings"
)

type Contracted struct {
	tableName struct{} `sql:"table_name"`
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
	file := "./contracted.gob"
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
	qs := []string{
		`create temp table nodes (id int, lat float, lon float);`,
		`create temp table ways (id int, nodes int[]);`,	
	}

	for _, q := range qs {
		_, err := db.ExecOne(q)
		if err != nil {
			log.Panicln("Can't create table", err)
		}
	}

	log.Println("Tables created")

	for _, node := range xmlData.Nodes {
		sqlString := fmt.Sprintf("INSERT INTO pg_temp.nodes (id, lat, lon) VALUES(%v, %v, %v);", node.Elem.ID, node.Lat, node.Lng)
		_, err := db.ExecOne(sqlString)
		if err != nil {
			log.Panicln("Can't insert nodes", err, sqlString)
		}
	}

	// Iterating []NdId{} to get only values
	NodesOut := func (nds []OsmStructs.NdId) (nodes []int32) {
		for _, id := range nds {
			nodes = append(nodes, id.ID)
		}
		return
	}

	for _, way := range xmlData.Ways {
		arr := NodesOut(way.Nds)
		arrStr := make([]string, len(arr))
		for i := range arr {
			arrStr[i] = strconv.Itoa(int(arr[i]))
		}
		sqlString := fmt.Sprintf("INSERT INTO pg_temp.ways (id, nodes) VALUES(%v, ARRAY[%v]::integer[]);", way.Elem.ID, strings.Join(arrStr, ","))
		_, err := db.ExecOne(sqlString)
		if err != nil {
			log.Panicln("Can't insert ways", err, sqlString)
		}
	}

   	log.Println("Nodes & ways inserted")

	// Contracting with pgr_contractGraph tool
	var res []Contracted
	sqlString := fmt.Sprintf(`
		SELECT seq, type, id, contracted_vertices, source, target, cost FROM pgr_contractGraph(
			'SELECT ways.id, nodes1.id as source, nodes2.id as target, 1 as cost FROM pg_temp.ways as ways 
			join pg_temp.nodes as nodes1 on nodes1.id = ways.nodes[1]
			join pg_temp.nodes as nodes2 on nodes2.id = ways.nodes[2]
			', ARRAY[1, 2]);
	`)
	_, err = db.Query((&res),sqlString)
	// log.Println("res:",res)
	if err != nil {
		// log.Panicln("can't contract:", err)
	} else {
		log.Println("Successfully contracted!")
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
