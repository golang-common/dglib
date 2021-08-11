package dql

const (
	TypeDefault  string = "default"
	TypeUid      string = "uid"
	TypeInt      string = "int"
	TypeFloat    string = "float"
	TypeString   string = "string"
	TypeBool     string = "bool"
	TypeDateTime string = "datetime"
	TypeGeo      string = "geo"
	TypePassword string = "password"
	TypeVar      string = "variable"

	IndexDefault  string = "default"
	IndexInt      string = "int"
	IndexFloat    string = "float"
	IndexBool     string = "bool"
	IndexGeo      string = "geo"
	IndexYear     string = "year"
	IndexMonth    string = "month"
	IndexDay      string = "day"
	IndexHour     string = "hour"
	IndexHash     string = "hash"
	IndexExact    string = "exact"
	IndexTerm     string = "term"
	IndexFulltext string = "fulltext"
	IndexTrigram  string = "trigram"
)

var (
	TypeMap = map[string]struct{}{
		TypeDefault:  {},
		TypeUid:      {},
		TypeInt:      {},
		TypeFloat:    {},
		TypeString:   {},
		TypeBool:     {},
		TypeDateTime: {},
		TypeGeo:      {},
		TypePassword: {},
		TypeVar:      {},
	}

	IndexMap = map[string]struct{}{
		IndexDefault:  {},
		IndexInt:      {},
		IndexFloat:    {},
		IndexBool:     {},
		IndexGeo:      {},
		IndexYear:     {},
		IndexMonth:    {},
		IndexDay:      {},
		IndexHour:     {},
		IndexHash:     {},
		IndexExact:    {},
		IndexTerm:     {},
		IndexFulltext: {},
		IndexTrigram:  {},
	}
)