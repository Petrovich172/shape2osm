package OsmStructs

import (
	"encoding/xml"
	// "../utils"
	// "time"
)

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