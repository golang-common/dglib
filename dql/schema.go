package dql

type (
	DIndex string // dgraph索引类型
	DType  string // dgraph类型
	DAttr  string // dgraph谓词属性

	Default  string // string
	Uid      string // int64
	Int      string // int64
	Float    string // float64
	String   string // string
	Bool     string // bool
	DateTime string // time.Time
	Geo      string // geom.T
	Password string // string
)

const (
	TDefault  DType = "default"
	TUid      DType = "uid"
	TInt      DType = "int"
	TFloat    DType = "float"
	TString   DType = "string"
	TBool     DType = "bool"
	TDateTime DType = "datetime"
	TGeo      DType = "geo"
	TPassword DType = "password"
	TVar      DType = "variable"

	IDefault  DIndex = "default"
	IInt      DIndex = "int"
	IFloat    DIndex = "float"
	IBool     DIndex = "bool"
	IGeo      DIndex = "geo"
	IYear     DIndex = "year"
	IMonth    DIndex = "month"
	IDay      DIndex = "day"
	IHour     DIndex = "hour"
	IHash     DIndex = "hash"
	IExact    DIndex = "exact"
	ITerm     DIndex = "term"
	IFulltext DIndex = "fulltext"
	ITrigram  DIndex = "trigram"

	IList    DIndex = "@list"
	ICount   DIndex = "@count"
	IIndex   DIndex = "@index"
	IUpsert  DIndex = "@upsert"
	IReverse DIndex = "@reverse"
	ILang    DIndex = "@lang"
)

var (
	DTypeMap = map[DType]struct{}{
		TDefault:  {},
		TUid:      {},
		TInt:      {},
		TFloat:    {},
		TString:   {},
		TBool:     {},
		TDateTime: {},
		TGeo:      {},
		TPassword: {},
		TVar:      {},
	}

	IGIndexMap = map[DIndex]struct{}{
		IList:    {},
		ICount:   {},
		IIndex:   {},
		IUpsert:  {},
		IReverse: {},
		ILang:    {},
	}

	IIndexMap = map[DIndex]struct{}{
		IDefault:  {},
		IInt:      {},
		IFloat:    {},
		IBool:     {},
		IGeo:      {},
		IYear:     {},
		IMonth:    {},
		IDay:      {},
		IHour:     {},
		IHash:     {},
		IExact:    {},
		ITerm:     {},
		IFulltext: {},
		ITrigram:  {},
	}
)
