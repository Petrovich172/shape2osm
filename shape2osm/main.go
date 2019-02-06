package main

import (
	"log"
	"github.com/go-pg/pg"
	"io/ioutil"
	"io"
	"os"
	"encoding/xml"

	// inits "pjob/pkgs/init"

	// "github.com/gin-gonic/gin"
	// "github.com/go-redis/redis"
)


// our struct which contains the complete
// array of all Users in the file
type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

// the user struct, this contains our
// Type attribute, our user's name and
// a social struct which will contain all
// our social links
type User struct {
	XMLName xml.Name `xml:"user"`
	Type    string   `xml:"type,attr"`
	Name    string   `xml:"name"`
	Social  Social   `xml:"social"`
}

// a simple struct which contains all our
// social links
type Social struct {
	XMLName  xml.Name `xml:"social"`
	Facebook string   `xml:"facebook"`
	Twitter  string   `xml:"twitter"`
	Youtube  string   `xml:"youtube"`
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

	getSomeData(db)

	// Open our xmlFile
	xmlFile, err := os.Open("sample.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully Opened sample.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var users Users
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &users)	

	for i := 0; i < len(users.Users); i++ {
		log.Println("User Type: " + users.Users[i].Type)
		log.Println("User Name: " + users.Users[i].Name)
		log.Println("Facebook Url: " + users.Users[i].Social.Facebook)
	}

	// creating output xml file
	f, err := os.Create("out.xml")
	if err != nil { panic(err) }
	defer f.Close()
	newFile := io.Writer(f)
	enc := xml.NewEncoder(newFile)
	enc.Indent("  ", "    ")
    	if err := enc.Encode(&users); err != nil {
				log.Printf("error: %v\n", err)
		}
}

func getSomeData(db *pg.DB) {
	var ret string
	var err error
	sqlString := "SELECT geom FROM graph.jytomir limit 1"
	_, err = db.Model().Query(&ret, sqlString)
	if err != nil {
		log.Println("some shit happend:", "\n", err)		
	}
	log.Println("query:","\n",ret)
}