package OsmStructs

import (
	"encoding/xml"
	// "../utils"
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