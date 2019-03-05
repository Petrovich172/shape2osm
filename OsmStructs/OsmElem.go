package OsmStructs

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