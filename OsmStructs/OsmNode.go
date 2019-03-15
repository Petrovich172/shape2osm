package OsmStructs

import (
	"encoding/xml"
)

// Node structure
type Node struct {
	Elem
	TableName          struct{} `sql:"nodes"`
	XMLName xml.Name `xml:"node"`
	Lat     float64  `xml:"lat,attr" sql:"lat"`
	Lng     float64  `xml:"lon,attr" sql:"lon"`
	Tags     []Tag    `xml:"tag" sql:"tags"`
}