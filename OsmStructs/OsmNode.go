package OsmStructs

import (
	"encoding/xml"
	// "../utils"
	// "time"
)

// Node structure
type Node struct {
	Elem
	XMLName xml.Name `xml:"node"`
	Lat     float64  `xml:"lat,attr"`
	Lng     float64  `xml:"lon,attr"`
	Tags     []Tag    `xml:"tag"`
}