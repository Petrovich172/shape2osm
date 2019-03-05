package OsmStructs

import (
	"encoding/xml"
	// "../utils"
	// "time"
)

// Tag struct
type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}