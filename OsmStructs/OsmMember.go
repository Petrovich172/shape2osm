package OsmStructs

import (
	// "encoding/xml"
	// "../utils"
	// "time"
)

// Member struct
type Member struct {
	Type string `xml:"type,attr"`
	Ref  int32  `xml:"ref,attr"`
	Role string `xml:"role,attr"`
}