package ContractedStructs

// Struct for pgr_contractGraph method
type Contracted struct {
	Seq					int32	`sql:"seq"`
	Type				string	`sql:"type"`
	Id					int32	`sql:"id"`
	ContractedVertices	[]int32	`sql:"contracted_vertices, type:integer[]"`
	Source				int32	`sql:"source"`
	Target				int32	`sql:"target"`
	Cost				float32	`sql:"cost"`
}