package ShapedStructs

import (
	// "encoding/xml"
	"../utils"
	// "time"
)

// Points cutted out from line geometry using pgr_createTopology
type ShapedNode struct {
	Id	int32	`xml:"id"`
	Geom	utils.PointString	`xml:"geom""`	
}