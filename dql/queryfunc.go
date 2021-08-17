/**
 * @Author: daipengyuan
 * @Description:
 * @File:  queryfunc
 * @Version: 1.0.0
 * @Date: 2021/8/12 19:42
 */

package dql

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Filter 开放给前端传入的
// Expr即关系表达式，如：(NOT A OR B) AND (C AND NOT (D OR E))
type Filter struct {
	Expr  string           `json:"expr,omitempty"`
	Funcs map[string]FFunc `json:"funcs,omitempty"`
	Facet bool             `json:"facet,omitempty"` // 是否为面过滤
}

// Parse 解析表达式与方法
func (f Filter) Parse() (string, error) {
	var rplcList []string
	// 解析方法替换列表
	fl := strings.FieldsFunc(f.Expr, func(r rune) bool {
		if r == '(' || r == ')' {
			return true
		}
		return false
	})
	for _, fs := range fl {
		fs = strings.Trim(fs, " ")
		upperFs := strings.ToUpper(fs)
		if upperFs == "AND" || upperFs == "OR" || upperFs == "NOT" {
			continue
		}
		v, ok := f.Funcs[fs]
		if !ok {
			return "", errors.New(fmt.Sprintf("function name [%s] not defined in funcs", fs))
		}
		vparse, err := v.Parse(f.Facet)
		if err != nil {
			return "", err
		}
		rplcList = append(rplcList, fs, vparse)
	}
	replacer := strings.NewReplacer(rplcList...)
	return fmt.Sprintf(`%s`, replacer.Replace(f.Expr)), nil
}

type FFunc struct {
	Key   string      `json:"key,omitempty"`
	Type  string      `json:"type,omitempty"`
	Val   interface{} `json:"val,omitempty"`
	Facet bool        `json:"facet,omitempty"`
}

func (f FFunc) Parse(facet ...bool) (string, error) {
	var fct bool
	if len(facet) > 0 && facet[0] == true {
		fct = true
	}
	if !fct {
		return f.parseFilter()
	}
	return f.parseFacet()
}

func (f FFunc) parseFacet() (string, error) {
	var (
		fs  string
		err error
	)
	switch f.Type {
	case FuncEq:
		fs, err = f.funcEqual()
	case FuncLe, FuncLt, FuncGe, FuncGt:
		fs, err = f.funcInequal()
	case FuncTermAll, FuncTermAny:
		fs, err = f.funcTerm()
	default:
		err = errors.New("unsupport function on facets filter " + f.Type)
	}
	if err != nil {
		return "", err
	}
	return fs, nil
}

func (f FFunc) parseFilter() (string, error) {
	err := f.checkFilterKey()
	if err != nil {
		return "", err
	}
	var fs string
	switch f.Type {
	case FuncEq:
		fs, err = f.funcEqual()
	case FuncLe, FuncLt, FuncGe, FuncGt:
		fs, err = f.funcInequal()
	case FuncTermAll, FuncTermAny:
		fs, err = f.funcTerm()
	case FuncMatch:
		fs, err = f.funcMatch()
	case FuncRegexp:
		fs, err = f.funcReg()
	case FuncTextAny, FuncTextAll:
		fs, err = f.funcFulltext()
	case FuncBetween:
		fs, err = f.funcBetween()
	case FuncUid:
		fs, err = f.funcUid()
	case FuncUidIn:
		fs, err = f.funcUidin()
	case FuncType:
		fs, err = f.funcType()
	case FuncHas:
		fs, err = f.funcHas()
	default:
		err = errors.New("unsupport function on filter " + f.Type)
	}
	if err != nil {
		return "", err
	}
	return fs, nil
}

// checkFilterKey 检查key对应的方法是否合法
func (f FFunc) checkFilterKey() error {
	// 如果为uid方法则忽略key值
	if f.Type == FuncUid {
		return nil
	}
	if f.Type == FuncType {
		if _, ok := TypeMap[f.Key]; !ok {
			return errors.New(fmt.Sprintf("target type does not exist in type func,[%s]", f.Key))
		}
	}
	pred, ok := PredMap[f.Key]
	if !ok {
		return errors.New(fmt.Sprintf("key of predicate [%s] does not exist", f.Key))
	}
	//查询谓词的类型属性
	tattr, ok := TypeAttrMap[pred.Type]
	if !ok {
		return errors.New(fmt.Sprintf("unsupport pred datatype [%s]", pred.Type))
	}
	// 如果属性对应的方法中找到目标方法，则返回
	var funcMatched bool
	for _, tafunc := range tattr.Fs {
		if tafunc == f.Type {
			funcMatched = true
			break
		}
	}
	if funcMatched {
		return nil
	}
	// 如果属性中未找到对应方法，则展开属性对应的索引方法中查找
	var findex []string
	for k, v := range tattr.Ts {
		for _, ifunc := range v.Fs {
			if ifunc == f.Type {
				findex = append(findex, k)
			}
		}
	}
	if len(findex) == 0 {
		return errors.New(fmt.Sprintf("key [%s] with type [%s] does not support func [%s]", f.Key, pred.Type, f.Type))
	}
	for _, fi := range findex {
		for _, tk := range pred.Tokenizer {
			if fi == tk {
				return nil
			}
		}
	}
	return errors.New(fmt.Sprintf("key [%s] with type [%s] token %v does not support func [%s],probobly need index %v", f.Key, pred.Type, pred.Tokenizer, f.Type, findex))
}

// funcEqual 等判断
func (f FFunc) funcEqual() (string, error) {
	switch f.Val.(type) {
	case string:
		return fmt.Sprintf(`eq(%s,"%s")`, f.Key, f.Val.(string)), nil
	case []string:
		return fmt.Sprintf(`eq(%s,["%s"])`, f.Key, strings.Join(f.Val.([]string), `","`)), nil
	case int:
		return fmt.Sprintf(`eq(%s,%d)`, f.Key, f.Val), nil
	case []int:
		var is []string
		for _, i := range f.Val.([]int) {
			is = append(is, strconv.Itoa(i))
		}
		return fmt.Sprintf(`eq(%s,[%s])`, f.Key, strings.Join(is, ",")), nil
	case float32, float64:
		return fmt.Sprintf(`eq(%s,%f)`, f.Key, f.Val), nil
	case []float32, []float64:
		var is []string
		for _, i := range f.Val.([]float64) {
			is = append(is, fmt.Sprintf("%f", i))
		}
		return fmt.Sprintf(`eq(%s,[%s])`, f.Key, strings.Join(is, ",")), nil
	case bool:
		return fmt.Sprintf(`eq(%s,"%t")`, f.Key, f.Val), nil
	case time.Time:
		return fmt.Sprintf(`eq(%s,"%s")`, f.Key, f.Val.(time.Time).String()), nil
	case []time.Time:
		var tm []string
		for _, t := range f.Val.([]time.Time) {
			tm = append(tm, t.String())
		}
		return fmt.Sprintf(`eq(%s,["%s"])`, f.Key, strings.Join(tm, `","`)), nil
	default:
		return "", errors.New(fmt.Sprintf("unsupport datatype on equal func, %s", reflect.TypeOf(f.Val)))
	}
}

// funcInequal 不等判断,判断左右值的相等，不等关系
func (f FFunc) funcInequal() (string, error) {
	switch f.Val.(type) {
	case string:
		return fmt.Sprintf(`%s(%s,"%s")`, f.Type, f.Key, f.Val.(string)), nil
	case int:
		return fmt.Sprintf(`%s(%s,%d)`, f.Type, f.Key, f.Val.(int)), nil
	case float32, float64:
		return fmt.Sprintf(`%s(%s,%f)`, f.Type, f.Key, f.Val), nil
	case time.Time:
		return fmt.Sprintf(`%s(%s,"%s")`, f.Type, f.Key, f.Val.(time.Time).String()), nil
	default:
		return "", errors.New(fmt.Sprintf("unsupport datatype on inequal func, %s", reflect.TypeOf(f.Val)))
	}
}

// funcTerm 查找字符串中的分组
func (f FFunc) funcTerm() (string, error) {
	if v, ok := f.Val.(string); ok {
		return fmt.Sprintf(`%s(%s,"%s")`, f.Type, f.Key, v), nil
	}
	return "", errors.New("unsupport datatype on term func,need string")
}

// funcReg 正则查询
func (f FFunc) funcReg() (string, error) {
	v, ok := f.Val.(string)
	if !ok {
		return "", errors.New("unsupport datatype on regexp func,need string")
	}
	if _, err := regexp.Compile(v); err != nil {
		return "", err
	}
	return fmt.Sprintf(`regexp(%s,/%s/)`, f.Key, v), nil
}

// funcMatch 字符串模糊查询
// 参数示例 {"val":"value","distance":2}
func (f FFunc) funcMatch() (string, error) {
	vmap, ok := f.Val.(map[string]interface{})
	if !ok {
		return "", errors.New(`unsupport datatype on match func,need map,e.g {"val":"value","distance":2}`)
	}
	val, ok := vmap["val"]
	if !ok {
		return "", errors.New(`match func map must have "val" key`)
	}
	valstr, ok := val.(string)
	if !ok {
		return "", errors.New(`match func map 'val' key must be type string`)
	}
	dist, ok := vmap["distance"]
	if !ok {
		return "", errors.New(`match func map must have "distance" key`)
	}
	distInt, ok := dist.(int)
	if !ok {
		return "", errors.New(`match func map 'distance' key must be type int`)
	}
	return fmt.Sprintf(`match(%s,"%s",%d)`, f.Key, valstr, distInt), nil
}

// funcFulltext 全文查找
func (f FFunc) funcFulltext() (string, error) {
	if v, ok := f.Val.(string); ok {
		return fmt.Sprintf(`%s(%s,"%s")`, f.Type, f.Key, v), nil
	}
	return "", errors.New("unsupport datatype on fulltext func,need string")
}

// funcBetween 范围查找
// 参数示例1 {"start":2,"end":5}
// 参数示例2 {"start":"192.168.1.100","end":"192.168.1.200"}
func (f FFunc) funcBetween() (string, error) {
	v, ok := f.Val.(map[string]interface{})
	if !ok {
		return "", errors.New(`unsupport datatype on between func,need map,e.g {"start":2,"end":5}`)
	}
	stt, okstt := v["start"]
	end, okend := v["end"]
	if !okstt || !okend {
		return "", errors.New(`between func map must have "start" and "end" key`)
	}
	if reflect.TypeOf(stt) != reflect.TypeOf(end) {
		return "", errors.New(`between func ,both "start" and "end" must have same datatype`)
	}
	switch stt.(type) {
	case int:
		return fmt.Sprintf("between(%s,%d,%d)", f.Key, stt, end), nil
	case float32, float64:
		return fmt.Sprintf("between(%s,%f,%f)", f.Key, stt, end), nil
	case string:
		return fmt.Sprintf(`between(%s,"%s","%s")`, f.Key, stt, end), nil
	case time.Time:
		return fmt.Sprintf(`between(%s,"%s","%s")`, f.Key, stt.(time.Time).String(), end.(time.Time).String()), nil
	default:
		return "", errors.New(`between func ,wrong data type on "start" and "end"`)
	}
}

// funcUid 查找uid
func (f FFunc) funcUid() (string, error) {
	switch f.Val.(type) {
	case string:
		return fmt.Sprintf(`uid("%s")`, f.Val), nil
	case int, int64:
		return fmt.Sprintf(`uid(%d)`, f.Val), nil
	case []string:
		return fmt.Sprintf(`uid(["%s"])`, strings.Join(f.Val.([]string), `","`)), nil
	case []int, []int64:
		var s []string
		for _, i := range f.Val.([]int64) {
			s = append(s, fmt.Sprintf("%d", i))
		}
		return fmt.Sprintf(`uid([%s])`, strings.Join(s, `,`)), nil
	default:
		return "", errors.New(fmt.Sprintf("unsupport datatype in uid func,type=[%s]", reflect.TypeOf(f.Val)))
	}
}

// funcUidin 查找谓词中的uid
func (f FFunc) funcUidin() (string, error) {
	switch f.Val.(type) {
	case string:
		return fmt.Sprintf(`uid_in(%s,"%s")`, f.Key, f.Val), nil
	case int, int64:
		return fmt.Sprintf(`uid_in(%s,%d)`, f.Key, f.Val), nil
	case []string:
		return fmt.Sprintf(`uid_in(%s,["%s"])`, f.Key, strings.Join(f.Val.([]string), `","`)), nil
	case []int, []int64:
		var s []string
		for _, i := range f.Val.([]int64) {
			s = append(s, fmt.Sprintf("%d", i))
		}
		return fmt.Sprintf(`uid_in(%s,[%s])`, f.Key, strings.Join(s, `,`)), nil
	default:
		return "", errors.New(fmt.Sprintf("unsupport datatype in uid func,type=[%s]", reflect.TypeOf(f.Val)))
	}
}

// funcHas 过滤包含某个谓词的结果
func (f FFunc) funcHas() (string, error) {
	return fmt.Sprintf("has(%s)", f.Key), nil
}

// funcType 过滤出某个类型
func (f FFunc) funcType() (string, error) {
	return fmt.Sprintf("type(%s)", f.Key), nil
}

// TODO:增加Geo相关过滤方法
