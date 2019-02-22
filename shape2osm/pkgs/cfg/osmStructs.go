package cfg

import (
	"encoding/xml"
	"../utils"
	// "time"
)

// OSM FORMAT STRUCTS
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

// Bounds struct
type Bounds struct {
	XMLName xml.Name `xml:"bounds"`
	Minlat  float64  `xml:"minlat,attr"`
	Minlon  float64  `xml:"minlon,attr"`
	Maxlat  float64  `xml:"maxlat,attr"`
	Maxlon  float64  `xml:"maxlon,attr"`
}

// Location struct
type Location struct {
	Type        string
	Coordinates []float64
}

// Tag struct
type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}

// Elem is a Osm base element
type Elem struct {
	ID        int32 `xml:"id,attr"`
	// Loc       Location
	Version   int       `xml:"version,attr"`
	// Ts        time.Time `xml:"timestamp,attr"`
	Ts        string `xml:"timestamp,attr"`
	// UID       int32     `xml:"uid,attr"`
	// User      string    `xml:"user,attr"`
	// ChangeSet int32     `xml:"changeset,attr"`
}

// Node structure
type Node struct {
	Elem
	XMLName xml.Name `xml:"node"`
	Lat     float64  `xml:"lat,attr"`
	Lng     float64  `xml:"lon,attr"`
	Tags     []Tag    `xml:"tag"`
}

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

// Member struct
type Member struct {
	Type string `xml:"type,attr"`
	Ref  int32  `xml:"ref,attr"`
	Role string `xml:"role,attr"`
}

// Relation struct
type Relation struct {
	Elem
	XMLName xml.Name `xml:"relation"`
	Visible bool     `xml:"visible,attr"`
	Members []Member `xml:"member"`
	Tags    []Tag    `xml:"tag"`
}


// SHAPED FORMAT STRUCTS
// Shaped body struct
type ShapeData struct {
	Edges	[]Edge
	NodedLines	[]NodedLine
}

// Way basic information
type NodedLine struct {
	Id			int32	`xml:"id"			sql:"id"`
	Source		int32	`xml:"source"		sql:"source"`
	Target		int32	`xml:"target"		sql:"target"`
	F_lanes		string	`xml:"f_lanes"		sql:f_lanes`
	T_lanes		string	`xml:"t_lanes"		sql:"t_lanes"`
	GmanTyp 	string 	`xml:"gman_typ"		sql:"gman_typ"`
	TlineTyp	string 	`xml:"tline_typ"	sql:"tline_typ"`
	Speedlim	string	`xml:"speedlim" 	sql:"speedlim"`
	RdName		string	`xml:"rd_name" 		sql:"rd_name"`
	Tollway		string 	`xml:"tollway" 		sql:"tollway"`
	SnipAd 		string 	`xml:"snip_ad"		sql:"snip_ad"`
	Btf 		string 	`xml:"btf"			sql:"btf"`
	RWeight 	string 	`xml:"r_weight"		sql:"r_weight"`
	RHeight 	string 	`xml:"r_height"		sql:"r_height"`
	RWidth 		string 	`xml:"r_width"		sql:"r_width"`
	Bicyclanes 	string 	`xml:"bicyclanes"	sql:"bicyclanes"`
	TBuslanes 	string 	`xml:"t_buslanes"	sql:"t_buslanes"`
	FBuslanes 	string 	`xml:"f_buslanes"	sql:"f_buslanes"`
	// Edge1id		int32	`xml:"ref"			sql:"edge1id"`
	Edge2id		int32	`xml:"ref"			sql:"edge2id"`
	Edge3id		int32	`xml:"ref"			sql:"edge3id"`
	Edge4id		int32	`xml:"ref"			sql:"edge4id"`
	Edge5id		int32	`xml:"ref"			sql:"edge5id"`
	Oneway		string	`xml:"oneway"		sql:"oneway"`
	Surface		string	`xml:"surface"		sql:"surface"`
	Highway		string	`xml:"highway"		sql:"highway"`
}

// Points cutted out from line geometry using pgr_createTopology
type Edge struct {
	Id	int32	`xml:"id"	sql:"id"`
	Geom	utils.PointString	`xml:"geom"	sql:"geom"`	
}