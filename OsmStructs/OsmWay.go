package OsmStructs

import (
	"encoding/xml"
)

// Way struct
type Way struct {
	Elem
	XMLName xml.Name `xml:"way"`
	Nds		[]NdId `xml:"nd" sql:"nodes, type:integer[]"`
	// Nds		[]NdId `xml:"nd" sql:"nodes"`
	Tags   []Tag `xml:"tag" sql:"tags"`

}

// Node Id in Way struct
type NdId struct {
		ID int32 `xml:"ref,attr" sql:"nodeId"`
}