package cfg

import (
	"encoding/xml"
	"io"
	"os"
	"strings"
	// "time"
)

// Osm body struct
type Osm struct {
	XMLName xml.Name `xml:"osm"`
	Version   string       `xml:"version,attr"`
	Ts        string `xml:"timestamp,attr"`
	Bounds    Bounds
	Nodes     []Node
	Ways      []Way
	Relations []Relation
}

// Bounds struct
type Bounds struct {
	XMLName xml.Name `xml:"bounds"`
	Minlat  float64  `xml:"minlat,attr"`
	Minlon  float64  `xml:"minlon,attr"`
	Maxlat  float64  `xml:"maxlat,attr"`
	Maxlon  float64  `xml:"maxlon,attr"`
}

// Location struct
type Location struct {
	Type        string
	Coordinates []float64
}

// Tag struct
type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}

// Elem is a Osm base element
type Elem struct {
	ID        int32 `xml:"id,attr"`
	// Loc       Location
	Version   int       `xml:"version,attr"`
	// Ts        time.Time `xml:"timestamp,attr"`
	Ts        string `xml:"timestamp,attr"`
	// UID       int32     `xml:"uid,attr"`
	// User      string    `xml:"user,attr"`
	// ChangeSet int32     `xml:"changeset,attr"`
}

// Node structure
type Node struct {
	Elem
	XMLName xml.Name `xml:"node"`
	Lat     float64  `xml:"lat,attr"`
	Lng     float64  `xml:"lon,attr"`
	Tags     []Tag    `xml:"tag"`
}

// Way struct
type Way struct {
	Elem
	XMLName xml.Name `xml:"way"`
	Nds		[]NdId `xml:"nd"`
	Tags   []Tag `xml:"tag"`

}

// Node Id in Way struct
type NdId struct {
		ID int32 `xml:"ref,attr"`
}

// Member struct
type Member struct {
	Type string `xml:"type,attr"`
	Ref  int32  `xml:"ref,attr"`
	Role string `xml:"role,attr"`
}

// Relation struct
type Relation struct {
	Elem
	XMLName xml.Name `xml:"relation"`
	Visible bool     `xml:"visible,attr"`
	Members []Member `xml:"member"`
	Tags    []Tag    `xml:"tag"`
}



// DecodeFile an Osm file
func DecodeFile(fileName string) (*Osm, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Decode(file)
}

func DecodeString(data string) (*Osm, error) {
	return Decode(strings.NewReader(data))
}

// Decode an reader
func Decode(reader io.Reader) (*Osm, error) {
	var (
		o   = new(Osm)
		err error
	)

	decoder := xml.NewDecoder(reader)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch typedToken := token.(type) {
		case xml.StartElement:
			switch typedToken.Name.Local {
			case "bounds":
				var b Bounds
				err = decoder.DecodeElement(&b, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Bounds = b

			case "node":
				var n Node
				err = decoder.DecodeElement(&n, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Nodes = append(o.Nodes, n)

			case "way":
				var w Way
				err = decoder.DecodeElement(&w, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Ways = append(o.Ways, w)

			case "relation":
				var r Relation
				err = decoder.DecodeElement(&r, &typedToken)
				if err != nil {
					return nil, err
				}
				o.Relations = append(o.Relations, r)
			}
		}
	}
	return o, nil
}