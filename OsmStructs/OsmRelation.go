package OsmStructs

import (
	"encoding/xml"
	// "../utils"
	// "time"
)

// Relation struct
type Relation struct {
	Elem
	XMLName xml.Name `xml:"relation"`
	Visible bool     `xml:"visible,attr"`
	Members []Member `xml:"member"`
	Tags    []Tag    `xml:"tag"`
}