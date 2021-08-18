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
	TokenDefault  string = "default"
	TokenInt      string = "int"
	TokenFloat    string = "float"
	TokenBool     string = "bool"
	TokenGeo      string = "geo"
	TokenYear     string = "year"
	TokenMonth    string = "month"
	TokenDay      string = "day"
	TokenHour     string = "hour"
	TokenHash     string = "hash"
	TokenExact    string = "exact"
	TokenTerm     string = "term"
	TokenFulltext string = "fulltext"
	TokenTrigram  string = "trigram"
)

const (
	IndexCount   string = "count"
	IndexList    string = "list"
	IndexLang    string = "lang"
	IndexReverse string = "reverse"
	IndexIndex   string = "index"
	IndexUpsert  string = "upsert"
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
	TokenList = []string{
		TokenDefault,
		TokenInt,
		TokenFloat,
		TokenBool,
		TokenGeo,
		TokenYear,
		TokenMonth,
		TokenDay,
		TokenHour,
		TokenHash,
		TokenExact,
		TokenTerm,
		TokenFulltext,
		TokenTrigram,
	}

	//IndicesList = []string{
	//	IndexCount,
	//	IndexList,
	//	IndexLang,
	//	IndexReverse,
	//	IndexIndex,
	//	IndexUpsert,
	//}

	TypeAttrMap = map[string]typeAttr{
		TypeDefault:  {Fs: []string{FuncHas}, Is: []string{IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenDefault: {}}},
		TypeUid:      {Fs: []string{FuncHas, FuncUidIn, FuncUid}, Is: []string{IndexList, IndexCount, IndexReverse}},
		TypeInt:      {Fs: []string{FuncHas}, Is: []string{IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenInt: {true, []string{FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}}}},
		TypeFloat:    {Fs: []string{FuncHas}, Is: []string{IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenFloat: {true, []string{FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}}}},
		TypeString:   {Fs: []string{FuncHas}, Is: []string{IndexLang, IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenHash: {false, []string{FuncEq}}, TokenExact: {true, []string{FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}}, TokenTerm: {false, []string{FuncEq, FuncTermAny, FuncTermAll}}, TokenFulltext: {false, []string{FuncEq, FuncTextAny, FuncTextAll}}, TokenTrigram: {false, []string{FuncRegexp, FuncMatch}}}},
		TypeBool:     {Fs: []string{FuncHas}, Is: []string{IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenBool: {false, []string{FuncEq}}}},
		TypeDateTime: {Fs: []string{FuncHas, FuncEq, FuncGe, FuncGt, FuncLt, FuncLe, FuncBetween}, Is: []string{IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenYear: {Stb: true}, TokenMonth: {Stb: true}, TokenDay: {Stb: true}, TokenHour: {Stb: true}}},
		TypeGeo:      {Fs: []string{FuncHas}, Is: []string{IndexList, IndexCount, IndexIndex, IndexUpsert}, Ts: map[string]tokenAttr{TokenGeo: {false, []string{FuncNear, FuncIntersects, FuncWithin, FuncContain}}}},
		TypePassword: {Fs: []string{FuncHas}},
	}
)

type typeAttr struct {
	Fs []string             // func surpport
	Is []string             // index surpport
	Ts map[string]tokenAttr // token surpport
}

type tokenAttr struct {
	Stb bool     // sortable
	Fs  []string // func surpport
}
