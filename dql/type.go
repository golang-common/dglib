/**
 * @Author: daipengyuan
 * @Description:
 * @File:  type
 * @Version: 1.0.0
 * @Date: 2021/8/11 12:18
 */

package dql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var TypeMap = map[string]Type{}

type Type struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields,omitempty"`
}

type Field struct {
	Name string `json:"name"`
}

func (t Type) Schema() string {
	var (
		preds []string
		r     string
	)
	for _, p := range t.Fields {
		var pr = p.Name
		if strings.HasPrefix(p.Name, "~") {
			pr = fmt.Sprintf("<%s>", p.Name)
		}
		preds = append(preds, pr)
	}
	if len(preds) > 0 {
		r = fmt.Sprintf("type %s{\n\t%s\n}", t.Name, strings.Join(preds, "\n\t"))
	}
	return r
}

func UnmarshalTypeString(s string) (*Type, error) {
	var (
		name      string
		plist     []string
		rplist    []string
		tbody     string
		errFormat = errors.New("error type format")
	)
	if len(s) < 12 {
		return nil, errFormat
	}
	s = strings.Trim(s, " \n")
	bsIndex := strings.Index(s, "{")
	beIndex := strings.Index(s, "}")
	if bsIndex == -1 || beIndex == -1 || bsIndex >= beIndex {
		return nil, errFormat
	}
	name = s[5:bsIndex]
	if name == "" {
		return nil, errFormat
	}
	tbody = s[bsIndex+1 : beIndex]
	tbody = strings.Trim(tbody, " \n\t")
	if len(tbody) == 0 {
		return nil, errFormat
	}
	for _, line := range strings.Split(tbody, "\n") {
		line = strings.Trim(line, " \t")
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "<~") && len(line) > 3 {
			rplist = append(rplist, line[1:len(line)-1])
		}
		plist = append(plist, line)
	}
	if len(plist) == 0 {
		return nil, errFormat
	}
	t := new(Type)
	t.Name = name
	for _, p := range plist {
		t.Fields = append(t.Fields, Field{Name: p})
	}
	return t, nil
}

// UnmarshalSchema 将结构体解析为Schema
func UnmarshalSchema(obj interface{}) (*Type, []Pred, error) {
	var (
		rType     *Type
		rPreds    []Pred
		dtype     string
		pmap      = make(map[string]Pred)
		typeField []Field
	)
	tp := reflect.TypeOf(obj)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tp.Kind() != reflect.Struct {
		return nil, nil, errors.New("UnmarshalSchema must give a struct value")
	}
	for i := 0; i < tp.NumField(); i++ {
		var (
			pred     = Pred{}
			predName string
			predType string
			tokens   []string
		)
		field := tp.Field(i)
		if field.Name == "Uid" {
			dtype = field.Tag.Get(TagDtype)
			continue
		}
		db := field.Tag.Get(TagDb)
		if db == "" {
			return nil, nil, errors.New("empty db tag field")
		}
		dbList := strings.Split(db, ",")
		if len(dbList) != 2 {
			return nil, nil, errors.New("unable to parse db tag,need [predname]:[datatype]")
		}
		predName = dbList[0]
		predType = dbList[1]
		// 不解析面
		if strings.Contains(predName, "|") {
			continue
		}
		if strings.Contains(predName, "@") {
			predName = strings.Split(predName, "@")[0]
		}
		typeField = append(typeField, Field{Name: predName})
		if strings.HasPrefix(predName, "~") {
			continue
		}
		pred.Predicate = predName
		pred.Type = predType
		if _, ok := TypeAttrMap[predType]; !ok {
			return nil, nil, errors.New(fmt.Sprintf("unsupport datatype %s", predType))
		}
		index := field.Tag.Get(TagIndex)
		indexList := strings.Split(index, ",")
		for _, idx := range indexList {
			for _, tl := range TokenList {
				if idx == tl {
					tokens = append(tokens, idx)
					break
				}
			}
			if idx == IndexCount {
				pred.Count = true
				continue
			}
			if idx == IndexList {
				pred.List = true
				continue
			}
			if idx == IndexLang {
				pred.Lang = true
				continue
			}
			if idx == IndexReverse {
				pred.Reverse = true
				continue
			}
			if idx == IndexIndex {
				pred.Index = true
				continue
			}
			if idx == IndexUpsert {
				pred.Upsert = true
			}
		}
		pred.Tokenizer = tokens
		if _, ok := pmap[pred.Predicate]; !ok {
			pmap[pred.Predicate] = pred
		}
	}
	if dtype == "" {
		return nil, nil, errors.New(`obj [Uid] field must have "dtype" tag,or Uid field not exist`)
	}
	rType = &Type{Name: dtype, Fields: typeField}
	for _, v := range pmap {
		rPreds = append(rPreds, v)
	}
	return rType, rPreds, nil
}
