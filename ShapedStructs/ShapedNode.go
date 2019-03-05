package ShapedStructs

// Points cutted out from line geometry using pgr_createTopology
type ShapedNode struct {
	Id	int32	`xml:"id"`
	Geom	PointString	`xml:"geom""`	
}

type PointString struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}