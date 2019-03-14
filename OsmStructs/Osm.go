package OsmStructs

import (
	"encoding/xml"
)

// Osm body struct
type Osm struct {
	XMLName xml.Name	 `xml:"osm"`
	Version   string     `xml:"version,attr" sql:"version"`
	Ts        string	 `xml:"timestamp,attr" sql:"ts"`
	Bounds    Bounds	 `sql:"bounds"`
	Nodes     []Node	 `sql:"nodes"`
	Ways      []Way		 `sql:"ways"`
	Relations []Relation `sql:"relations"`
}