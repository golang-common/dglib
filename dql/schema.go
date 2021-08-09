package dql

import (
	"errors"
	"fmt"
	"strings"
)

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

	AList    DAttr = "@list"
	ACount   DAttr = "@count"
	AIndex   DAttr = "@index"
	AUpsert  DAttr = "@upsert"
	AReverse DAttr = "@reverse"
	ALang    DAttr = "@lang"
)

var (
	// attrMap 属性字典，用于速查
	attrMap = map[DAttr]struct{}{
		AList:    {},
		ACount:   {},
		AIndex:   {},
		AUpsert:  {},
		AReverse: {},
		ALang:    {},
	}
	// typeIndexMap 数据类型支持的索引类型的字典
	typeIndexMap = map[DType]map[DIndex]struct{}{
		TDefault:  {IDefault: {}},
		TUid:      {},
		TInt:      {IInt: {}},
		TFloat:    {IFloat: {}},
		TString:   {IHash: {}, IExact: {}, ITerm: {}, IFulltext: {}, ITrigram: {}},
		TBool:     {IBool: {}},
		TDateTime: {IYear: {}, IMonth: {}, IDay: {}, IHour: {}},
		TGeo:      {IGeo: {}},
		TPassword: {},
	}
	// typeAttrMap 数据类型支持的属性类型的字典
	typeAttrMap = map[DType]map[DAttr]struct{}{
		TDefault:  {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}},
		TUid:      {AList: {}, ACount: {}, AReverse: {}},
		TInt:      {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}},
		TFloat:    {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}},
		TString:   {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}, ALang: {}},
		TBool:     {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}},
		TDateTime: {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}},
		TGeo:      {AList: {}, ACount: {}, AIndex: {}, AUpsert: {}},
		TPassword: {},
	}
)

type Pred interface {
	DType() DType
	Name() string
	Schema(attrs []DAttr, idx []DIndex) Schema
}

func (Default) DType() DType {
	return TDefault
}

func (s Default) Name() string {
	return string(s)
}

func (s Default) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (Uid) DType() DType {
	return TUid
}

func (s Uid) Name() string {
	return string(s)
}

func (s Uid) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (Int) DType() DType {
	return TInt
}

func (s Int) Name() string {
	return string(s)
}

func (s Int) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (Float) DType() DType {
	return TFloat
}

func (s Float) Name() string {
	return string(s)
}

func (s Float) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (String) DType() DType {
	return TString
}

func (s String) Name() string {
	return string(s)
}

func (s String) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (Bool) DType() DType {
	return TBool
}

func (s Bool) Name() string {
	return string(s)
}

func (s Bool) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (DateTime) DType() DType {
	return TDateTime
}

func (s DateTime) Name() string {
	return string(s)
}

func (s DateTime) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (Geo) DType() DType {
	return TGeo
}

func (s Geo) Name() string {
	return string(s)
}

func (s Geo) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

func (Password) DType() DType {
	return TPassword
}

func (s Password) Name() string {
	return string(s)
}

func (s Password) Schema(attrs []DAttr, idx []DIndex) Schema {
	return Schema{
		Name:    s.Name(),
		Type:    s.DType(),
		Attrs:   attrs,
		Indices: idx,
	}
}

type Schema struct {
	Name    string
	Type    DType
	Attrs   []DAttr
	Indices []DIndex
}

func (s Schema) MashalRdf() (string, error) {
	var (
		tp        = string(s.Type)
		result    string
		indexSet  bool
		dAttrList []string
		dIdxList  []string
		model     = "%s: %s %s ."
	)
	err := s.checkSchema()
	if err != nil {
		return "", err
	}
	for _, attr := range s.Attrs {
		if attr == AList {
			tp = fmt.Sprintf("[%s]", tp)
			continue
		}
		var str = string(attr)
		if attr == AIndex {
			indexSet = true
			str += "(%s)"
		}
		dAttrList = append(dAttrList, str)
	}
	result = fmt.Sprintf(model, s.Name, tp, strings.Join(dAttrList, " "))
	if indexSet {
		for _, idx := range s.Indices {
			dIdxList = append(dIdxList, string(idx))
		}
		result = fmt.Sprintf(result, strings.Join(dIdxList, ","))
	}
	return result, nil
}

func (s *Schema) UnmashalRdf(str string) error {
	slist := strings.Split(str, " ")
	for _, v := range slist {
		if strings.HasSuffix(v, ":") {
			s.Name = v[0 : len(v)-1]
			continue
		}
		if !strings.HasPrefix(v, "@") && v != "." {
			s.Type = DType(v)
			continue
		}
		if strings.HasPrefix(v, "@") {
			if _, ok := attrMap[DAttr(v)]; ok {
				s.Attrs = append(s.Attrs, DAttr(v))
				continue
			}
			if strings.HasPrefix(v, "@index") {
				s.Attrs = append(s.Attrs, AIndex)
				idxStr := v[strings.Index(v, "(")+1 : strings.Index(v, ")")]
				for _, idx := range strings.Split(idxStr, ",") {
					s.Indices = append(s.Indices, DIndex(idx))
				}
			}
		}
	}
	err := s.checkSchema()
	if err != nil {
		return err
	}
	return nil
}

func (s Schema) checkSchema() error {
	const ErrBase = "check schema err,"
	var (
		aIndexSet   bool
		aUpsertSet  bool
		iStrBaseSet bool // exact hash term 其中之一则置真
		iDateSet    bool // year month day hour 其中之一则置真
	)
	if s.Name == "" {
		return errors.New(ErrBase + "empty schema name")
	}
	_, ok := typeAttrMap[s.Type]
	if !ok {
		return errors.New(fmt.Sprintf("%s,[%s] unsupport type %s", ErrBase, s.Name, string(s.Type)))
	}
	for _, attr := range s.Attrs {
		if _, ok = typeAttrMap[s.Type][attr]; !ok {
			return errors.New(
				fmt.Sprintf("%s,[%s]'s type %s unsupport attr %s", ErrBase, s.Name, s.Type, attr))
		}
		if attr == AIndex {
			aIndexSet = true
			continue
		}
		if attr == AUpsert {
			aUpsertSet = true
		}
	}
	if aUpsertSet && !aIndexSet {
		return errors.New(
			fmt.Sprintf(
				"%s,@index tokenizer is mandatory for: [%s] when specifying @upsert directive", ErrBase, s.Name))
	}
	if !aIndexSet && len(s.Indices) > 0 {
		return errors.New(
			fmt.Sprintf("%s,[%s]'s some indices given,but @index attr is not set", ErrBase, s.Name))
	}
	if aIndexSet && len(s.Indices) == 0 {
		return errors.New(
			fmt.Sprintf("%s,[%s]'s @index attr set,but no indices given", ErrBase, s.Name))
	}
	_, ok = typeIndexMap[s.Type]
	if !ok {
		return errors.New(fmt.Sprintf("%s,[%s] unsupport type %s", ErrBase, s.Name, string(s.Type)))
	}
	for _, index := range s.Indices {
		if _, ok = typeIndexMap[s.Type]; !ok {
			return errors.New(
				fmt.Sprintf("%s,[%s]'s type %s unsupport index %s", ErrBase, s.Name, s.Type, index))
		}
		if index == IHash || index == IExact || index == ITerm {
			if iStrBaseSet == true {
				return errors.New(
					fmt.Sprintf("%s,[%s]'s type %s can only set hash/exact/term once", ErrBase, s.Name, s.Type))
			}
			iStrBaseSet = true
			continue
		}
		if index == IYear || index == IMonth || index == IDay || index == IHour {
			if iDateSet == true {
				return errors.New(fmt.Sprintf("%s,[%s]'s type %s can only set year/month/day/hour once",
					ErrBase, s.Name, s.Type))
			}
			iDateSet = true
			continue
		}
	}
	return nil
}

type Type struct {
	Name   string
	Preds  []string
	RPreds []string
}

func (s Type) MashalRdf() (string, error) {
	var (
		r     string
		model = "type %s{%s\n}"
		plist []string
	)
	if s.Name == "" {
		return "", errors.New("empty type name")
	}
	if len(s.Preds) == 0 {
		return "", errors.New("type don't have any predicates")
	}
	for _, pred := range s.Preds {
		plist = append(plist, pred)
	}
	for _, rpred := range s.RPreds {
		plist = append(plist, fmt.Sprintf("<~%s>", rpred))
	}
	r = fmt.Sprintf(model, s.Name, strings.Join(plist, "\n\t"))
	return r, nil
}

func (s *Type) UnmashalRdf(str string) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()
	tname := str[5:strings.Index(str, "{")]
	s.Name = tname
	tprestr := str[strings.Index(str, "{")+1 : strings.Index(str, "}")]
	tprestr = strings.Trim(tprestr, "\t\n")
	tprestr = strings.Replace(tprestr, "\t", "", -1)
	preList := strings.Split(tprestr, "\n")
	for _, pre := range preList {
		if !strings.HasPrefix(pre, "<") {
			s.Preds = append(s.Preds, pre)
		} else {
			s.RPreds = append(s.RPreds, pre[2:len(pre)-1])
		}
	}
	return
}
