package ShapedStructs

// Way basic information
type NodedLine struct {
	Id			int32	`xml:"id"`
	Source		int32	`xml:"source"`
	Target		int32	`xml:"target"`
	F_lanes		string	`xml:"f_lanes"`
	T_lanes		string	`xml:"t_lanes"`
	GmanTyp 	string 	`xml:"gman_typ"`
	TlineTyp	string 	`xml:"tline_typ"`
	Speedlim	string	`xml:"speedlim"`
	RdName		string	`xml:"rd_name"`
	Tollway		string 	`xml:"tollway"`
	SnipAd 		string 	`xml:"snip_ad"`
	Btf 		string 	`xml:"btf"`
	RWeight 	string 	`xml:"r_weight"`
	RHeight 	string 	`xml:"r_height"`
	RWidth 		string 	`xml:"r_width"`
	Bicyclanes 	string 	`xml:"bicyclanes"`
	TBuslanes 	string 	`xml:"t_buslanes"`
	FBuslanes 	string 	`xml:"f_buslanes"`
	// Edge1id		int32	`xml:"ref"`
	Edge2id		int32	`xml:"ref"`
	Edge3id		int32	`xml:"ref"`
	Edge4id		int32	`xml:"ref"`
	Edge5id		int32	`xml:"ref"`
	Oneway		string	`xml:"oneway"`
	Surface		string	`xml:"surface"`
	Highway		string	`xml:"highway"`
}