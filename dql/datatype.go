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
)
const (
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

const (
	FuncEq         = "eq"
	FuncGt         = "gt"
	FuncGe         = "ge"
	FuncLt         = "lt"
	FuncLe         = "le"
	FuncTermAll    = "allofterms"
	FuncTermAny    = "anyofterms"
	FuncRegexp     = "regexp"
	FuncMatch      = "match"
	FuncBetween    = "between"
	FuncTextAll    = "alloftext"
	FuncTextAny    = "anyoftext"
	FuncHas        = "has"
	FuncType       = "type"
	FuncUid        = "uid"
	FuncUidIn      = "uid_in"
	FuncNear       = "near"
	FuncWithin     = "within"
	FuncContain    = "contain"
	FuncIntersects = "intersects"
)

var (
	TypeAttrMap = map[string]typeAttr{
		TypeDefault:  {Fs: []string{FuncHas}, Ts: map[string]tokenAttr{IndexDefault: {}}},
		TypeUid:      {Fs: []string{FuncHas, FuncUidIn, FuncUid}},
		TypeInt:      {Fs: []string{FuncHas}, Ts: map[string]tokenAttr{IndexInt: {true, []string{FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}}}},
		TypeFloat:    {Fs: []string{FuncHas}, Ts: map[string]tokenAttr{IndexFloat: {true, []string{FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}}}},
		TypeString:   {Fs: []string{FuncHas}, Ts: map[string]tokenAttr{IndexHash: {false, []string{FuncEq}}, IndexExact: {true, []string{FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}}, IndexTerm: {false, []string{FuncEq, FuncTermAny, FuncTermAll}}, IndexFulltext: {false, []string{FuncEq, FuncTextAny, FuncTextAll}}, IndexTrigram: {false, []string{FuncRegexp, FuncMatch}}}},
		TypeBool:     {Fs: []string{FuncHas}, Ts: map[string]tokenAttr{IndexBool: {false, []string{FuncEq}}}},
		TypeDateTime: {Fs: []string{FuncHas, FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}, Ts: map[string]tokenAttr{IndexYear: {Stb: true}, IndexMonth: {Stb: true}, IndexDay: {Stb: true}, IndexHour: {Stb: true}}},
		TypeGeo:      {Fs: []string{FuncHas}, Ts: map[string]tokenAttr{IndexGeo: {false, []string{FuncNear, FuncIntersects, FuncWithin, FuncContain}}}},
		TypePassword: {Fs: []string{FuncHas}},
	}
)

type typeAttr struct {
	Fs []string             // func surpport
	Ts map[string]tokenAttr // token surpport
}

type tokenAttr struct {
	Stb bool     // sortable
	Fs  []string // func surpport
}
